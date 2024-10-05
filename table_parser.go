package nestcsv

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TableParser struct {
	td *TableData
}

func NewTableParser(td *TableData) *TableParser {
	return &TableParser{td: td}
}

func (p *TableParser) ParseTableFields(tags []string) ([]*TableField, error) {
	var (
		td     = p.td
		fields = make([]*TableField, 0, td.Columns)
	)

	for col := 0; col < td.Columns; col++ {
		if !containsAny(tags, td.FieldTags[col]...) {
			continue
		}

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
					return nil, fmt.Errorf("nested multi-line array is not allowed: %s, %s", td.Name, field.Name)
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
				parentField = findPtr(fields, func(f *TableField) bool {
					return f.Name == field.Name
				})
				if parentField == nil {
					parentField = field
					fields = append(fields, field)
				}
			}
		}
	}
	return fields, nil
}

func (p *TableParser) Marshal(fields []*TableField) (any, error) {
	var (
		td         = p.td
		rowMap     = make(map[string]map[string]any)
		rows       = make([]map[string]any, 0, len(rowMap))
		rowIndices = make([]int, 0, len(rowMap))

		multiLineArrayRowCount = make(map[string]int)
	)

	for rowIdx, row := range td.DataRows {
		id := row[TableFieldIndexCol]
		rowContainer, isMultiLineRow := rowMap[id]
		if !isMultiLineRow {
			rowContainer = make(map[string]any)
			rowMap[id] = rowContainer

			if !td.Metadata.AsMap {
				rows = append(rows, rowContainer)
				rowIndices = append(rowIndices, rowIdx)
			}
		}

		for _, topField := range fields {
			container := rowContainer
			for field := range topField.Iterate() {
				var multiLineArrayIdx int
				multiLineArrayField := field.GetMultiLineArrayField()
				if multiLineArrayField != nil {
					rowCountIdx := id + "_" + multiLineArrayField.Identifier()
					if multiLineArrayField == field {
						if p.checkAllCellsEmpty(field, row) {
							break
						}
						multiLineArrayRowCount[rowCountIdx]++
					}
					multiLineArrayIdx = multiLineArrayRowCount[rowCountIdx] - 1

				} else if isMultiLineRow {
					// skip if it's non array fields
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
								return nil, fmt.Errorf("failed to parse array value: %s, %s, %d, %s, %w", td.Name, field.Name, rowIdx, cell, err)
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
						return nil, fmt.Errorf("failed to parse value: %s, %s, %d, %s, %w", td.Name, field.Name, rowIdx, cell, err)
					}
					container[field.Name] = v
				}
			}
		}
	}

	if td.Metadata.AsMap {
		m := make(map[string]any)
		for id, row := range rowMap {
			m[id] = row
		}
		return m, nil

	} else {
		if td.Metadata.SortAscBy != "" {
			p.sortValues(rows, rowIndices, td.Metadata.SortAscBy, false)
		} else if td.Metadata.SortDescBy != "" {
			p.sortValues(rows, rowIndices, td.Metadata.SortDescBy, true)
		}
		return rows, nil
	}
}

func (p *TableParser) parseGoValue(typ FieldType, cell string) (any, error) {
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

func (p *TableParser) checkAllCellsEmpty(field *TableField, row []string) bool {
	for f := range field.Iterate() {
		if row[f.column] != "" {
			return false
		}
	}
	return true
}

func (p *TableParser) sortValues(values []map[string]any, rowIndices []int, field string, desc bool) {
	fieldCol := slices.Index(p.td.FieldNames, field)
	fieldType, _ := newFieldType(p.td.FieldTypes[fieldCol])
	sort.SliceStable(values, func(i, j int) bool {
		av, _ := p.parseGoValue(fieldType, p.td.DataRows[rowIndices[i]][fieldCol])
		bv, _ := p.parseGoValue(fieldType, p.td.DataRows[rowIndices[j]][fieldCol])
		compareAsc := p.sortCompareAsc(av, bv)
		if desc {
			return !compareAsc
		} else {
			return compareAsc
		}
	})
}

func (p *TableParser) sortCompareAsc(a, b any) bool {
	switch a.(type) {
	case int:
		return a.(int) < b.(int)
	case int64:
		return a.(int64) < b.(int64)
	case string:
		return strings.Compare(a.(string), b.(string)) < 0
	case float64:
		return a.(float64) < b.(float64)
	case time.Time:
		return a.(time.Time).Before(b.(time.Time))
	default:
		log.Panicf("unsupported type: %T", a)
		return false
	}
}
