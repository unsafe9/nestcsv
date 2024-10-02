package nestcsv

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type tableMarshaler struct {
	table *Table
}

func (t *Table) Marshal() any {
	m := &tableMarshaler{
		table: t,
	}
	return m.marshal()
}

func (m *tableMarshaler) marshal() any {
	t := m.table
	idKey := t.Fields[0].Name

	var values any
	if t.Metadata.AsMap {
		m := make(map[string]any)
		for _, v := range t.Values {
			idStr := fmt.Sprint(v[idKey])
			if t.Metadata.DropID {
				m[idStr] = mapWithoutKey(v, idKey)
			} else {
				m[idStr] = v
			}
		}
		values = m

	} else {
		arr := slices.Clone(t.Values)
		if t.Metadata.SortAscBy != "" {
			m.sortValuesAsc(arr, t.Metadata.SortAscBy)
		} else if t.Metadata.SortDescBy != "" {
			m.sortValuesDesc(arr, t.Metadata.SortDescBy)
		}
		if t.Metadata.DropID {
			for i, v := range arr {
				arr[i] = mapWithoutKey(v, idKey)
			}
		}
		values = arr
	}

	return values
}

func (m *tableMarshaler) getNestedValue(value map[string]any, fieldName string) any {
	keys := strings.Split(fieldName, ".")
	lastIdx := len(keys) - 1
	for _, key := range keys[:lastIdx] {
		value = value[key].(map[string]any)
	}
	return value[keys[lastIdx]]
}

func (m *tableMarshaler) sortValuesAsc(values []map[string]any, field string) {
	slices.SortStableFunc(values, func(a, b map[string]any) int {
		av := m.getNestedValue(a, field)
		bv := m.getNestedValue(b, field)
		switch av.(type) {
		case int:
			return av.(int) - bv.(int)
		case int64:
			return int(av.(int64) - bv.(int64))
		case string:
			return strings.Compare(av.(string), bv.(string))
		case float64:
			return int(av.(float64) - bv.(float64))
		case time.Time:
			return int(av.(time.Time).Sub(bv.(time.Time)))
		default:
			panic("unsupported type")
		}
	})
}

func (m *tableMarshaler) sortValuesDesc(values []map[string]any, field string) {
	slices.SortStableFunc(values, func(a, b map[string]any) int {
		av := m.getNestedValue(a, field)
		bv := m.getNestedValue(b, field)
		switch av.(type) {
		case int:
			return bv.(int) - av.(int)
		case int64:
			return int(bv.(int64) - av.(int64))
		case string:
			return strings.Compare(bv.(string), av.(string))
		case float64:
			return int(bv.(float64) - av.(float64))
		case time.Time:
			return int(bv.(time.Time).Sub(av.(time.Time)))
		default:
			panic("unsupported type")
		}
	})
}
