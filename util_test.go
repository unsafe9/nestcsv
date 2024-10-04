package nestcsv

import "testing"

func TestHas(t *testing.T) {
	a := []FieldType{FieldTypeTime, FieldTypeInt, FieldTypeString}
	if !has(a, "time") {
		t.Error("expected true")
	}
	if has(a, "float") {
		t.Error("expected false")
	}
}
