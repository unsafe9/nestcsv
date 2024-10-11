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

func (f *TableField) GetMultiLineArrayField() *TableField {
	for p := f; p != nil; p = p.ParentField {
		if p.IsMultiLineArray {
			return p
		}
	}
	return nil
}

func (f *TableField) printDebug(depth int) {
	println(strings.Repeat("  ", depth), f.Name, f.Type)
	for _, sf := range f.StructFields {
		sf.printDebug(depth + 1)
	}
}

func (f *TableField) StructEqual(other *TableField) bool {
	return f.structEqual(other, true)
}

func (f *TableField) structEqual(other *TableField, top bool) bool {
	if !top && f.Name != other.Name {
		return false
	}
	if f.Type != other.Type || f.IsCellArray != other.IsCellArray {
		return false
	}
	if len(f.StructFields) != len(other.StructFields) {
		return false
	}
	for i, sf := range f.StructFields {
		if !sf.structEqual(other.StructFields[i], false) {
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

func (f *TableField) Clone() *TableField {
	clone := &TableField{
		Name:             f.Name,
		Type:             f.Type,
		IsMultiLineArray: f.IsMultiLineArray,
		IsCellArray:      f.IsCellArray,
		column:           f.column,
	}
	for _, sf := range f.StructFields {
		sfClone := sf.Clone()
		sfClone.ParentField = clone
		clone.StructFields = append(clone.StructFields, sfClone)
	}
	return clone
}
