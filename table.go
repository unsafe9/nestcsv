package nestcsv

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type TableField struct {
	Name         string
	Type         string
	IsArray      bool
	StructFields []*TableField
	ParentField  *TableField
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

func ParseTable(fileName string, csvData []byte) (*Table, error) {
	reader := csv.NewReader(bytes.NewReader(csvData))
	reader.Comment = '#'
	reader.TrimLeadingSpace = true
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read the csv: %s, %w", fileName, err)
	}

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
			Values: make([]map[string]any, dataLen),
		}
		rowMap               = make(map[string]map[string]any)
		multiLineArrayIdxMap = make(map[int]int)
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
	for i := range rows {
		id := rows[i][0]
		if value, ok := rowMap[id]; ok {
			table.Values[i] = value
			multiLineArrayIdxMap[i]++
		} else {
			idValue, err := parseGoValue(types[0], id)
			if err != nil {
				return nil, fmt.Errorf("failed to parse id value: %s, %s, %d, %s, %w", table.Name, names[0], i, id, err)
			}
			value = map[string]any{
				names[0]: idValue,
			}
			table.Values[i] = value
			rowMap[id] = value
		}
	}

	for col := 1; col < colLen; col++ {
		if strings.HasPrefix(names[col], "#") {
			continue
		}

		nameTokens := strings.Split(names[col], ".")
		tokenLen := len(nameTokens)
		tokenFields := make([]*TableField, tokenLen)

		isCellArray := strings.HasPrefix(types[col], "[]")
		if isCellArray {
			types[col] = types[col][len("[]"):]
			if types[col] == "json" {
				return nil, fmt.Errorf("json type is not allowed for cell array: %s, %s", table.Name, names[col])
			}
		}

		isInMultiLineArray := false

		for i := 0; i < tokenLen; i++ {
			field := &TableField{
				Name: nameTokens[i],
			}
			tokenFields[i] = field

			if i == 0 {
				table.Fields = append(table.Fields, field)
			} else {
				field.ParentField = tokenFields[i-1]
				field.ParentField.StructFields = append(field.ParentField.StructFields, field)
			}

			isMultiLineArray := strings.HasPrefix(field.Name, "[]")
			if isMultiLineArray {
				if isInMultiLineArray {
					return nil, fmt.Errorf("nested multi-line array is not allowed: %s, %s", table.Name, field.Name)
				}
				isInMultiLineArray = true
				field.Name = field.Name[len("[]"):]
				field.IsArray = true
			}

			if i == tokenLen-1 {
				field.Type = types[col]
				field.IsArray = field.IsArray || isCellArray
			} else {
				field.Type = "struct"
			}

			for j, row := range rows {
				if multiLineArrayIdxMap[j] > 0 && !isInMultiLineArray {
					continue
				}

				container := table.Values[j]
				for _, tokenField := range tokenFields[:i] {
					switch v := container[tokenField.Name].(type) {
					case map[string]any:
						container = v
					case []map[string]any:
						idx := multiLineArrayIdxMap[j]
						if len(v) <= idx {
							v = append(v, make(map[string]any))
							container[tokenField.Name] = v
						}
						container = v[idx]
					default:
						return nil, fmt.Errorf("invalid container type: %s, %s, %d", table.Name, tokenField.Name, j)
					}
				}

				if i < tokenLen-1 {
					if field.IsArray {
						if _, ok := container[field.Name]; !ok {
							container[field.Name] = make([]map[string]any, 0)
						}
					} else {
						if _, ok := container[field.Name]; !ok {
							container[field.Name] = make(map[string]any)
						}
					}
					continue
				}

				if field.IsArray {
					if _, ok := container[field.Name]; !ok {
						container[field.Name] = make([]any, 0)
					}
				}

				cell := row[col]
				if isCellArray {
					var arr []any
					if len(cell) > 0 {
						cells := strings.Split(cell, ",")
						for _, elem := range cells {
							v, err := parseGoValue(field.Type, elem)
							if err != nil {
								return nil, fmt.Errorf("failed to parse array value: %s, %s, %d, %s, %w", table.Name, field.Name, j, cell, err)
							}
							arr = append(arr, v)
						}
					} else {
						arr = make([]any, 0)
					}
					container[field.Name] = arr

				} else {
					v, err := parseGoValue(field.Type, cell)
					if err != nil {
						return nil, fmt.Errorf("failed to parse value: %s, %s, %d, %s, %w", table.Name, field.Name, j, cell, err)
					}
					if field.IsArray {
						container[field.Name] = append(container[field.Name].([]any), v)
					} else {
						container[field.Name] = v
					}
				}
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
