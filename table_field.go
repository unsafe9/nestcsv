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
	FieldTypeJSON   FieldType = "json"
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
	if t == FieldTypeJSON {
		panic("unsupported type: json array")
	}
	return FieldType("[]" + t.String())
}

func (t FieldType) isValidIndexType() bool {
	switch t {
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

func (f *TableField) Identifier() string {
	if f.ParentField != nil {
		return f.ParentField.Identifier() + "." + f.Name
	}
	return f.Name
}

func (f *TableField) IsInMultiLineArray() bool {
	for p := f; p != nil; p = p.ParentField {
		if p.IsMultiLineArray {
			return true
		}
	}
	return false
}

func (f *TableField) printDebug(depth int) {
	println(strings.Repeat("  ", depth), f.Name, f.Type)
	for _, sf := range f.StructFields {
		sf.printDebug(depth + 1)
	}
}
