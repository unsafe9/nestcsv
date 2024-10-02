package nestcsv

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type tableParser struct {
	tableData *TableData
}

func ParseTable(td *TableData) (*Table, error) {
	p := &tableParser{
		tableData: td,
	}
	return p.parse()
}

func (p *tableParser) parse() (*Table, error) {
	var (
		td    = p.tableData
		table = Table{
			Name:     td.Name,
			Metadata: td.Metadata,
			Fields:   make([]*TableField, 0, td.Columns),
			Values:   make([]map[string]any, 0, len(td.DataRows)),
		}
		rowMap                 = make(map[string]map[string]any)
		multiLineArrayRowCount = make(map[string]int)
		multiLineArrayIdxMap   = make(map[int]int)
	)

	idField := &TableField{
		Name: td.FieldNames[TableFieldIndexCol],
		Type: FieldType(td.FieldTypes[TableFieldIndexCol]),
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
			idValue, err := p.parseGoValue(idField.Type, id)
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
		var (
			nameTokens             = strings.Split(td.FieldNames[col], ".")
			tokenLen               = len(nameTokens)
			fieldType, isCellArray = newFieldType(td.FieldTypes[col])
			multiLineArrayField    *TableField
			parentField            *TableField
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
				field.Type = fieldType
				field.IsCellArray = isCellArray
			} else {
				field.Type = FieldTypeStruct
			}

			if parentField != nil {
				field.ParentField = parentField
				existingField := findPtr(parentField.StructFields, func(f *TableField) bool {
					return f.Name == field.Name
				})
				if existingField != nil {
					field = existingField
				} else {
					parentField.StructFields = append(parentField.StructFields, field)
				}
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

	for rowIdx, row := range td.DataRows {
		for _, topField := range table.Fields {
			container := rowMap[row[TableFieldIndexCol]]
			for field := range topField.Iterate() {
				multiLineArrayIdx := multiLineArrayIdxMap[rowIdx]

				// skip if it's non array fields
				if multiLineArrayIdx > 0 && !field.IsInMultiLineArray() {
					continue
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
							if p.checkAllCellsEmpty(field, row) {
								defer func(container map[string]any, name string) {
									container[name] = removeOne(
										container[name].([]map[string]any),
										func(m map[string]any) bool {
											return equalPtr(m, v)
										},
									)
								}(container, field.Name)
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

				} else if field.IsCellArray {
					// fill array value
					cell := row[field.column]
					var arr []any
					if len(cell) > 0 {
						cells := strings.Split(cell, ",")
						for _, elem := range cells {
							v, err := p.parseGoValue(field.Type, elem)
							if err != nil {
								return nil, fmt.Errorf("failed to parse array value: %s, %s, %d, %s, %w", table.Name, field.Name, rowIdx, cell, err)
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
					v, err := p.parseGoValue(field.Type, cell)
					if err != nil {
						return nil, fmt.Errorf("failed to parse value: %s, %s, %d, %s, %w", table.Name, field.Name, rowIdx, cell, err)
					}
					container[field.Name] = v
				}
			}
		}
	}

	return &table, nil
}

func (p *tableParser) parseGoValue(typ FieldType, cell string) (any, error) {
	switch typ {
	case FieldTypeInt:
		if cell == "" {
			return 0, nil
		}
		return strconv.Atoi(cell)
	case FieldTypeLong:
		if cell == "" {
			return int64(0), nil
		}
		return strconv.ParseInt(cell, 10, 64)
	case FieldTypeFloat:
		if cell == "" {
			return float64(0), nil
		}
		return strconv.ParseFloat(cell, 64)
	case FieldTypeBool:
		if cell == "" {
			return false, nil
		}
		return strconv.ParseBool(cell)
	case FieldTypeString:
		return cell, nil
	case FieldTypeTime:
		if cell == "" {
			return time.Time{}, nil
		}
		return time.Parse(time.DateTime, cell)
	case FieldTypeJSON:
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

func (p *tableParser) checkAllCellsEmpty(field *TableField, row []string) bool {
	var cells []string
	for f := range field.Iterate() {
		cells = append(cells, row[f.column])
	}
	return isAllStringEmpty(cells)
}
