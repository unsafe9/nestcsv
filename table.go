package nestcsv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type TableField struct {
	Name             string
	Type             string
	IsMultiLineArray bool
	IsCellArray      bool
	StructFields     []*TableField
	ParentField      *TableField
	column           int
}

type Table struct {
	Name   string
	Fields []*TableField
	Values []map[string]any
}

type TableSaveOption struct {
	RootDir string `yaml:"root_dir"`
	Indent  string `yaml:"indent"`
	AsMap   bool   `yaml:"as_map"`
	DropID  bool   `yaml:"drop_id"`
}

func (t *Table) SaveAsJson(option *TableSaveOption) error {
	if option.RootDir == "" {
		option.RootDir = "."
	}
	if err := os.MkdirAll(option.RootDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create the directory: %s, %w", option.RootDir, err)
	}

	file, err := os.Create(filepath.Join(option.RootDir, t.Name+".json"))
	if err != nil {
		return fmt.Errorf("failed to create the file: %s, %w", t.Name, err)
	}
	defer file.Close()

	idKey := t.Fields[0].Name

	var values any
	if option.AsMap {
		m := make(map[string]any)
		for _, v := range t.Values {
			idStr := fmt.Sprint(v[idKey])
			if option.DropID {
				v = shallowCopyMap(v)
				delete(v, idKey)
			}
			m[idStr] = v
		}
		values = m

	} else {
		if option.DropID {
			arr := make([]map[string]any, len(t.Values))
			for i, v := range t.Values {
				v = shallowCopyMap(v)
				delete(v, idKey)
				arr[i] = v
			}
			values = arr
		} else {
			values = t.Values
		}
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", option.Indent)

	if err := encoder.Encode(values); err != nil {
		return fmt.Errorf("failed to encode the json: %s, %w", t.Name, err)
	}

	return nil
}

func checkAllCellsEmpty(field *TableField, row []string) bool {
	var cells []string
	var visitField func(f *TableField)
	visitField = func(f *TableField) {
		cells = append(cells, row[f.column])
		for _, structField := range f.StructFields {
			visitField(structField)
		}
	}
	visitField(field)
	return isAllEmpty(cells)
}

func ParseTable(td *TableData) (*Table, error) {
	var (
		table = Table{
			Name:   td.Name,
			Fields: make([]*TableField, 0, td.Columns),
			Values: make([]map[string]any, 0, len(td.DataRows)),
		}
		rowMap                 = make(map[string]map[string]any)
		multiLineArrayRowCount = make(map[string]int)
		multiLineArrayIdxMap   = make(map[int]int)
	)

	idField := &TableField{
		Name: td.FieldNames[TableFieldIndexCol],
		Type: td.FieldTypes[TableFieldIndexCol],
	}
	table.Fields = append(table.Fields, idField)

	for i := 0; i < len(td.DataRows); i++ {
		id := td.DataRows[i][TableFieldIndexCol]

		if value, ok := rowMap[id]; ok {
			if _, ok := multiLineArrayRowCount[id]; !ok {
				multiLineArrayRowCount[id] = 2
			} else {
				multiLineArrayRowCount[id]++
			}
			multiLineArrayIdxMap[i] = multiLineArrayRowCount[id] - 1
		} else {
			idValue, err := parseGoValue(idField.Type, id)
			if err != nil {
				return nil, fmt.Errorf("failed to parse id value: %s, %s, %d, %s, %w", table.Name, td.FieldNames[0], i, id, err)
			}
			value = map[string]any{
				idField.Name: idValue,
			}
			table.Values = append(table.Values, value)
			rowMap[id] = value
		}
	}

	for col := TableFieldIndexCol + 1; col < td.Columns; col++ {
		nameTokens := strings.Split(td.FieldNames[col], ".")
		tokenLen := len(nameTokens)
		valueType := td.FieldTypes[col]

		isCellArray := strings.HasPrefix(valueType, "[]")
		if isCellArray {
			valueType = valueType[len("[]"):]
			if valueType == "json" {
				return nil, fmt.Errorf("json type is not allowed for cell array: %s, %s", table.Name, td.FieldNames[col])
			}
		}

		var (
			multiLineArrayField *TableField
			parentField         *TableField
		)

		for i := 0; i < tokenLen; i++ {
			field := &TableField{
				Name:   nameTokens[i],
				column: col,
			}

			isMultiLineArray := strings.HasPrefix(field.Name, "[]")
			if isMultiLineArray {
				if multiLineArrayField != nil {
					return nil, fmt.Errorf("nested multi-line array is not allowed: %s, %s", table.Name, field.Name)
				}
				multiLineArrayField = field
				field.Name = field.Name[len("[]"):]
				field.IsMultiLineArray = true
			}

			if i == tokenLen-1 {
				field.Type = valueType
				field.IsCellArray = isCellArray
			} else {
				field.Type = "struct"
			}

			if parentField != nil {
				field.ParentField = parentField
				parentField.StructFields = append(parentField.StructFields, field)
				parentField = field
			} else {
				parentField = findPtr(table.Fields, func(f *TableField) bool {
					return f.Name == field.Name
				})
				if parentField == nil {
					parentField = field
					table.Fields = append(table.Fields, field)
				}
			}
		}
	}

	var visitField func(container map[string]any, field *TableField, rowIdx int, row []string) error
	visitField = func(container map[string]any, field *TableField, rowIdx int, row []string) error {
		multiLineArrayIdx := multiLineArrayIdxMap[rowIdx]

		// skip if it's non array fields
		if multiLineArrayIdx > 0 {
			isInMultiLineArray := false
			for f := field; f != nil; f = f.ParentField {
				if f.IsMultiLineArray {
					isInMultiLineArray = true
					break
				}
			}
			if !isInMultiLineArray {
				return nil
			}
		}

		if len(field.StructFields) > 0 {
			if field.IsMultiLineArray {
				// fill struct array container
				objectArrayValue, ok := container[field.Name]
				if !ok {
					objectArrayValue = make([]map[string]any, 0)
				}
				objectArray := objectArrayValue.([]map[string]any)
				if len(objectArray) <= multiLineArrayIdx {
					v := make(map[string]any)
					objectArray = append(objectArray, v)

					// remove the object if all cells are empty
					if checkAllCellsEmpty(field, row) {
						defer func(container map[string]any) {
							container[field.Name] = removeOne(
								container[field.Name].([]map[string]any),
								func(m map[string]any) bool {
									return equalPtr(m, v)
								},
							)
						}(container)
					}
				}
				container[field.Name] = objectArray
				container = objectArray[multiLineArrayIdx]

			} else {
				// fill struct container
				objectValue, ok := container[field.Name]
				if !ok {
					objectValue = make(map[string]any)
					container[field.Name] = objectValue
				}
				container = objectValue.(map[string]any)
			}

			// fill struct fields recursively
			for _, structField := range field.StructFields {
				if err := visitField(container, structField, rowIdx, row); err != nil {
					return err
				}
			}

		} else if field.IsCellArray {
			// fill array value
			cell := row[field.column]
			var arr []any
			if len(cell) > 0 {
				cells := strings.Split(cell, ",")
				for _, elem := range cells {
					v, err := parseGoValue(field.Type, elem)
					if err != nil {
						return fmt.Errorf("failed to parse array value: %s, %s, %d, %s, %w", table.Name, field.Name, rowIdx, cell, err)
					}
					arr = append(arr, v)
				}
			} else {
				arr = make([]any, 0)
			}
			container[field.Name] = arr

		} else {
			// fill single value
			cell := row[field.column]
			v, err := parseGoValue(field.Type, cell)
			if err != nil {
				return fmt.Errorf("failed to parse value: %s, %s, %d, %s, %w", table.Name, field.Name, rowIdx, cell, err)
			}
			container[field.Name] = v
		}

		return nil
	}

	for rowIdx, row := range td.DataRows {
		container := rowMap[row[TableFieldIndexCol]]
		for _, field := range table.Fields {
			if err := visitField(container, field, rowIdx, row); err != nil {
				return nil, err
			}
		}
	}

	return &table, nil
}

func parseGoValue(typ, cell string) (any, error) {
	switch typ {
	case "int":
		if cell == "" {
			return 0, nil
		}
		return strconv.Atoi(cell)
	case "long":
		if cell == "" {
			return int64(0), nil
		}
		return strconv.ParseInt(cell, 10, 64)
	case "float":
		if cell == "" {
			return float64(0), nil
		}
		return strconv.ParseFloat(cell, 64)
	case "time":
		if cell == "" {
			return time.Time{}, nil
		}
		return time.Parse(time.DateTime, cell)
	case "string":
		return cell, nil
	case "bool":
		if cell == "" {
			return false, nil
		}
		return strconv.ParseBool(cell)
	case "json":
		if cell == "" {
			return nil, nil
		}
		v := new(any)
		if err := json.Unmarshal([]byte(cell), v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal json: %s, %w", cell, err)
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unknown type: %s", typ)
	}
}
