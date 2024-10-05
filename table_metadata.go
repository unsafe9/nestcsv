package nestcsv

import (
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type TableMetadata struct {
	AsMap       bool              `query:"as_map"`
	SortAscBy   string            `query:"sort_asc_by"`
	SortDescBy  string            `query:"sort_desc_by"`
	StructTypes map[string]string `query:"struct_type"`
}

func (m *TableMetadata) Validate(td *TableData) error {
	if m.AsMap && (m.SortAscBy != "" || m.SortDescBy != "") {
		return fmt.Errorf("as_map and sort_by are mutually exclusive")
	}

	if m.SortAscBy != "" && m.SortDescBy != "" {
		return fmt.Errorf("both sort_asc_by and sort_desc_by are set")
	}
	if m.SortAscBy != "" {
		if err := m.validateSortByField(td, m.SortAscBy); err != nil {
			return err
		}
	}
	if m.SortDescBy != "" {
		if err := m.validateSortByField(td, m.SortDescBy); err != nil {
			return err
		}
	}

	return nil
}

func (m *TableMetadata) validateSortByField(td *TableData, field string) error {
	col := slices.Index(td.FieldNames, field)
	if col == -1 {
		return fmt.Errorf("sort_by: field not found: %s", field)
	}
	fieldType := td.FieldTypes[col]
	if strings.Contains(field, "[]") || strings.Contains(fieldType, "[]") {
		return fmt.Errorf("sort_by: field is array: %s", field)
	}
	if fieldType == "json" || fieldType == "bool" {
		return fmt.Errorf("sort_by: invalid field type: %s, %s", field, fieldType)
	}
	return nil
}

type TableMetadataQuery string

func (q TableMetadataQuery) Decode() (*TableMetadata, error) {
	values, err := url.ParseQuery(string(q))
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}

	var metadata TableMetadata
	v := reflect.ValueOf(&metadata).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("query")
		if tag == "" || tag == "-" {
			continue
		}
		if val, ok := values[tag]; ok {
			if err := q.parseStringSliceInto(v.Field(i), val); err != nil {
				return nil, err
			}
		}
	}
	return &metadata, nil
}

func (q TableMetadataQuery) parseStringSliceInto(field reflect.Value, val []string) error {
	if len(val) == 0 {
		panic("empty value")
	}

	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	switch field.Kind() {
	case reflect.Slice:
		slice := reflect.MakeSlice(field.Type(), len(val), len(val))
		for i, v := range val {
			if err := q.parseStringInto(slice.Index(i), v); err != nil {
				return err
			}
		}
		field.Set(slice)

	case reflect.Map:
		m := reflect.MakeMap(field.Type())
		for _, v := range val {
			parts := strings.Split(v, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid map value: %s", v)
			}
			key := reflect.New(field.Type().Key()).Elem()
			if err := q.parseStringInto(key, parts[0]); err != nil {
				return err
			}
			value := reflect.New(field.Type().Elem()).Elem()
			if err := q.parseStringInto(value, parts[1]); err != nil {
				return err
			}
			m.SetMapIndex(key, value)
		}
		field.Set(m)

	default:
		return q.parseStringInto(field, val[0])
	}

	return nil
}

func (q TableMetadataQuery) parseStringInto(field reflect.Value, val string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	default:
		return fmt.Errorf("unsupported type: %v", field.Kind())
	}
	return nil
}
