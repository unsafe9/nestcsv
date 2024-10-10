// Code generated by "nestcsv"; DO NOT EDIT.

package table

import (
	"context"
)

type Tables struct {
	Complex ComplexTable
	Types   TypesTable
}

var tables *Tables

func Get() *Tables {
	return tables
}

type tablesContextKey struct{}

var tablesContextKeyInstance = tablesContextKey{}

func WithTables(ctx context.Context, t *Tables) context.Context {
	return context.WithValue(ctx, tablesContextKeyInstance, t)
}

func TablesFromContext(ctx context.Context) *Tables {
	t, _ := ctx.Value(tablesContextKeyInstance).(*Tables)
	return t
}

func LoadTablesFromFile(basePath string) (*Tables, error) {
	var t Tables
	if err := t.Complex.LoadFromFile(basePath); err != nil {
		return nil, err
	}
	if err := t.Types.LoadFromFile(basePath); err != nil {
		return nil, err
	}
	tables = &t
	return &t, nil
}

func (t *Tables) GetTables() []TableBase {
	return []TableBase{
		&t.Complex,
		&t.Types,
	}
}

func (t *Tables) GetBySheetName(sheetName string) TableBase {
	switch sheetName {
	case "complex":
		return &t.Complex
	case "types":
		return &t.Types
	default:
		return nil
	}
}