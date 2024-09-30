package datasource

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type TableData struct {
	Name string
	Rows [][]string
}

func NewTableData(name string, rows [][]string) TableData {
	return TableData{
		Name: strings.TrimSuffix(filepath.Base(name), ".csv"),
		Rows: rows,
	}
}

func (c *TableData) Save(rootDir string) error {
	if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create the directory: %s, %w", rootDir, err)
	}
	file, err := os.Create(filepath.Join(rootDir, c.Name+".csv"))
	if err != nil {
		return fmt.Errorf("failed to create the file: %s, %w", c.Name, err)
	}
	defer file.Close()

	if err := csv.NewWriter(file).WriteAll(c.Rows); err != nil {
		return fmt.Errorf("failed to write the file: %s, %w", c.Name, err)
	}
	return nil
}
