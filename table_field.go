package nestcsv

import "strings"

type FieldType string

const (
	FieldTypeInt    FieldType = "int"
	FieldTypeLong   FieldType = "long"
	FieldTypeFloat  FieldType = "float"
	FieldTypeBool   FieldType = "bool"
	FieldTypeString FieldType = "string"
	FieldTypeTime   FieldType = "time"
	FieldTypeJson   FieldType = "json"
	FieldTypeStruct FieldType = "struct"
)

func newFieldType(s string) (FieldType, bool) {
	if strings.HasPrefix(s, "[]") {
		return FieldType(s[len("[]"):]), true
	}
	return FieldType(s), false
}

func (t FieldType) String() string {
	return string(t)
}

func (t FieldType) Array() FieldType {
	if t == FieldTypeJson {
		panic("unsupported type: json array")
	}
	return FieldType("[]" + t.String())
}

func isValidIndexType(t string) bool {
	switch FieldType(t) {
	case FieldTypeInt, FieldTypeLong, FieldTypeString:
		return true
	}
	return false
}

type TableField struct {
	Name             string
	Type             FieldType
	IsMultiLineArray bool
	IsCellArray      bool
	StructFields     []*TableField
	ParentField      *TableField
	column           int
}

func (f *TableField) IsArray() bool {
	if len(f.StructFields) > 0 {
		return f.IsMultiLineArray
	}
	return f.IsMultiLineArray || f.IsCellArray
}
