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

func ParseTable(fileName string, rows [][]string) (*Table, error) {
	const headerRows = 3
	if len(rows) < headerRows {
		return nil, fmt.Errorf("invalid csv data: %s", fileName)
	}

	var (
		names   = rows[0]
		types   = rows[1]
		colLen  = len(names)
		dataLen = len(rows) - headerRows
		table   = Table{
			Name:   fileName,
			Fields: make([]*TableField, 0, colLen),
			Values: make([]map[string]any, 0, dataLen),
		}
		rowMap                 = make(map[string]map[string]any)
		multiLineArrayRowCount = make(map[string]int)
		multiLineArrayIdxMap   = make(map[int]int)
	)
	rows = rows[headerRows:]

	// pre-process ID field. The first field must be ID.
	if strings.Contains(names[0], ".") || (types[0] != "int" && types[0] != "long" && types[0] != "string") {
		return nil, fmt.Errorf("invalid index field: %s, %s, %s", table.Name, names[0], types[0])
	}
	table.Fields = append(table.Fields, &TableField{
		Name: names[0],
		Type: types[0],
	})
	for i := 0; i < dataLen; i++ {
		id := rows[i][0]

		// drop the row if the id is empty or starts with #
		if id == "" || strings.HasPrefix(id, "#") {
			rows = append(rows[:i], rows[i+1:]...)
			i--
			dataLen--
			continue
		}

		if value, ok := rowMap[id]; ok {
			if _, ok := multiLineArrayRowCount[id]; !ok {
				multiLineArrayRowCount[id] = 2
			} else {
				multiLineArrayRowCount[id]++
			}
			multiLineArrayIdxMap[i] = multiLineArrayRowCount[id] - 1
		} else {
			idValue, err := parseGoValue(types[0], id)
			if err != nil {
				return nil, fmt.Errorf("failed to parse id value: %s, %s, %d, %s, %w", table.Name, names[0], i, id, err)
			}
			value = map[string]any{
				names[0]: idValue,
			}
			table.Values = append(table.Values, value)
			rowMap[id] = value
		}
	}

	for col := 1; col < colLen; col++ {
		if strings.HasPrefix(names[col], "#") {
			continue
		}

		nameTokens := strings.Split(names[col], ".")
		tokenLen := len(nameTokens)

		isCellArray := strings.HasPrefix(types[col], "[]")
		if isCellArray {
			types[col] = types[col][len("[]"):]
			if types[col] == "json" {
				return nil, fmt.Errorf("json type is not allowed for cell array: %s, %s", table.Name, names[col])
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
				field.Type = types[col]
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

	for rowIdx, row := range rows {
		container := rowMap[row[0]]
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
