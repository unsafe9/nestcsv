package nestcsv

import (
	"iter"
	"strings"
)

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

func (f *TableField) Equal(other *TableField) bool {
	if f.Name != other.Name || f.Type != other.Type || f.IsMultiLineArray != other.IsMultiLineArray || f.IsCellArray != other.IsCellArray {
		return false
	}
	if len(f.StructFields) != len(other.StructFields) {
		return false
	}
	for i, sf := range f.StructFields {
		if !sf.Equal(other.StructFields[i]) {
			return false
		}
	}
	return true
}

func (f *TableField) Iterate() iter.Seq[*TableField] {
	return func(yield func(*TableField) bool) {
		var iterate func(f *TableField) bool
		iterate = func(f *TableField) bool {
			if yield(f) {
				for _, structField := range f.StructFields {
					if !iterate(structField) {
						return false
					}
				}
				return true
			}
			return false
		}
		iterate(f)
	}
}
