package nestcsv

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"
)

type TableSaveOption struct {
	RootDir  string `yaml:"root_dir"`
	Indent   string `yaml:"indent"`
	FileType string `yaml:"file_type"`
}

func (t *Table) Save(option *TableSaveOption) error {
	if option.RootDir == "" {
		option.RootDir = "."
	}
	if option.FileType == "" {
		option.FileType = "json"
	}

	value := t.marshalValue()

	if option.FileType == "json" {
		jsonBytes, err := json.MarshalIndent(value, "", option.Indent)
		if err != nil {
			return fmt.Errorf("failed to marshal the json: %s, %w", t.Name, err)
		}
		if err := saveJSONFile(option.RootDir, t.Name, jsonBytes); err != nil {
			return err
		}
	} else if option.FileType == "bin" {
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal the json: %s, %w", t.Name, err)
		}
		if err := saveBinFile(option.RootDir, t.Name, jsonBytes); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid file type: %s", option.FileType)
	}

	return nil
}

func (t *Table) marshalValue() any {
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
		if t.Metadata.DropID {
			arr := make([]map[string]any, len(t.Values))
			for i, v := range t.Values {
				arr[i] = mapWithoutKey(v, idKey)
			}
			values = arr
		} else {
			values = t.Values
		}
		if t.Metadata.SortAscBy != "" {
			sortValuesAsc(values.([]map[string]any), t.Metadata.SortAscBy)
		} else if t.Metadata.SortDescBy != "" {
			sortValuesDesc(values.([]map[string]any), t.Metadata.SortDescBy)
		}
	}

	return values
}

func getNestedValue(value map[string]any, fieldName string) any {
	keys := strings.Split(fieldName, ".")
	lastIdx := len(keys) - 1
	for _, key := range keys[:lastIdx] {
		value = value[key].(map[string]any)
	}
	return value[keys[lastIdx]]
}

func sortValuesAsc(values []map[string]any, field string) {
	slices.SortStableFunc(values, func(a, b map[string]any) int {
		av := getNestedValue(a, field)
		bv := getNestedValue(b, field)
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

func sortValuesDesc(values []map[string]any, field string) {
	slices.SortStableFunc(values, func(a, b map[string]any) int {
		av := getNestedValue(a, field)
		bv := getNestedValue(b, field)
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
