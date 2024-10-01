package nestcsv

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

const (
	TableMetadataRow  = 0
	TableFieldNameRow = 1
	TableFieldTypeRow = 2
	TableDataStartRow = 4

	TableFieldIndexCol = 0
)

var ErrSkipTable = fmt.Errorf("skip table")

type TableData struct {
	Name       string
	Metadata   *TableMetadata
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

	table := &TableData{
		Name:       tableName,
		Columns:    columns,
		FieldNames: fieldNames,
		FieldTypes: fieldTypes,
		DataRows:   dataRows,
	}

	metadata, err := ParseTableMetadata(csvData[TableMetadataRow][0], table)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table metadata: %s, %w", name, err)
	}

	table.Metadata = metadata
	return table, nil
}

func (c *TableData) CSV() [][]string {
	return append([][]string{c.FieldNames, c.FieldTypes}, c.DataRows...)
}
