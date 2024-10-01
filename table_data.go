package nestcsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	TableFieldNameRow = 0
	TableFieldTypeRow = 1
	TableDataStartRow = 3

	TableFieldIndexCol = 0
)

var ErrSkipTable = fmt.Errorf("skip table")

type TableData struct {
	Name       string
	Columns    int
	FieldNames []string
	FieldTypes []string
	DataRows   [][]string
}

func ParseTableData(name string, csvData [][]string) (*TableData, error) {
	tableName := strings.TrimSuffix(filepath.Base(name), ".csv")
	if strings.HasPrefix(tableName, "#") {
		return nil, ErrSkipTable
	}

	if len(csvData) < TableDataStartRow {
		return nil, fmt.Errorf("invalid table data: %s", name)
	}

	var (
		columns     = len(csvData[TableFieldNameRow])
		dropColumns = make([]int, 0)
		fieldNames  = make([]string, 0, columns)
		fieldTypes  = make([]string, 0, columns)
	)

	for col := 0; col < columns; col++ {
		fieldName := csvData[TableFieldNameRow][col]
		if fieldName == "" || strings.HasPrefix(fieldName, "#") {
			dropColumns = append(dropColumns, col)
		} else {
			fieldNames = append(fieldNames, fieldName)
			fieldTypes = append(fieldTypes, csvData[TableFieldTypeRow][col])
		}
	}

	columns -= len(dropColumns)
	if columns == 0 {
		return nil, fmt.Errorf("no columns in the csv file: %s", name)
	}

	dataRows := make([][]string, 0, len(csvData))
	for _, row := range csvData[TableDataStartRow:] {
		id := row[TableFieldIndexCol]
		if id == "" || strings.HasPrefix(id, "#") {
			continue
		}

		dataRow := make([]string, columns)
		i := 0
		for col, cell := range row {
			if !slices.Contains(dropColumns, col) {
				dataRow[i] = cell
				i++
			}
		}
		dataRows = append(dataRows, dataRow)
	}

	idxName := fieldNames[TableFieldIndexCol]
	if strings.Contains(idxName, ".") {
		return nil, fmt.Errorf("invalid index field: %s, %s", tableName, idxName)
	}
	idxType := fieldTypes[TableFieldIndexCol]
	if idxType != "int" && idxType != "long" && idxType != "string" {
		return nil, fmt.Errorf("invalid index field type: %s, %s, %s", tableName, idxName, idxType)
	}

	return &TableData{
		Name:       tableName,
		Columns:    columns,
		FieldNames: fieldNames,
		FieldTypes: fieldTypes,
		DataRows:   dataRows,
	}, nil
}

func (c *TableData) SaveAsCSV(rootDir string) error {
	if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create the directory: %s, %w", rootDir, err)
	}
	file, err := os.Create(filepath.Join(rootDir, c.Name+".csv"))
	if err != nil {
		return fmt.Errorf("failed to create the file: %s, %w", c.Name, err)
	}
	defer file.Close()

	rows := append([][]string{c.FieldNames, c.FieldTypes}, c.DataRows...)
	if err := csv.NewWriter(file).WriteAll(rows); err != nil {
		return fmt.Errorf("failed to write the file: %s, %w", c.Name, err)
	}
	return nil
}
