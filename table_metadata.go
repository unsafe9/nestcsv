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
	DropID       bool              `query:"drop_id"`
	AsMap        bool              `query:"as_map"`
	SortAscBy    string            `query:"sort_asc_by"`
	SortDescBy   string            `query:"sort_desc_by"`
	TypeMappings map[string]string `query:"type_mapping"`
}

func ParseTableMetadata(a1 string, td *TableData) (*TableMetadata, error) {
	values, err := url.ParseQuery(a1)
	if err != nil {
		return nil, fmt.Errorf("failed to parse a1 cell: %w", err)
	}

	var metadata TableMetadata
	if err := decodeQuery(values, &metadata); err != nil {
		return nil, err
	}

	if metadata.AsMap && (metadata.SortAscBy != "" || metadata.SortDescBy != "") {
		return nil, fmt.Errorf("as_map and sort_by are mutually exclusive")
	}

	if metadata.SortAscBy != "" && metadata.SortDescBy != "" {
		return nil, fmt.Errorf("both sort_asc_by and sort_desc_by are set")
	}
	if metadata.SortAscBy != "" {
		if err := validateSortByField(td, metadata.SortAscBy); err != nil {
			return nil, err
		}
	}
	if metadata.SortDescBy != "" {
		if err := validateSortByField(td, metadata.SortDescBy); err != nil {
			return nil, err
		}
	}

	return &metadata, nil
}

func validateSortByField(td *TableData, field string) error {
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

func decodeQuery(values url.Values, structPtr any) error {
	v := reflect.ValueOf(structPtr).Elem()
	t := v.Type()
	if t.Kind() != reflect.Struct {
		panic("not a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("query")
		if tag == "" || tag == "-" {
			continue
		}
		if val, ok := values[tag]; ok {
			if err := parseStringSliceInto(v.Field(i), val); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseStringSliceInto(field reflect.Value, val []string) error {
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
			if err := parseStringInto(slice.Index(i), v); err != nil {
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
			if err := parseStringInto(key, parts[0]); err != nil {
				return err
			}
			value := reflect.New(field.Type().Elem()).Elem()
			if err := parseStringInto(value, parts[1]); err != nil {
				return err
			}
			m.SetMapIndex(key, value)
		}
		field.Set(m)

	default:
		return parseStringInto(field, val[0])
	}

	return nil
}

func parseStringInto(field reflect.Value, val string) error {
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
