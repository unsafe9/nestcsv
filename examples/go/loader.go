// Code generated by "nestcsv"; DO NOT EDIT.

package table

type Tables struct {
	Complex ComplexTable
	Types   TypesTable
}

func LoadTables() (*Tables, error) {
	var t Tables
	if err := t.Complex.Load(); err != nil {
		return nil, err
	}
	if err := t.Types.Load(); err != nil {
		return nil, err
	}
	return &t, nil
}
