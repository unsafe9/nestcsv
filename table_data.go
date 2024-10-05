package nestcsv

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

const (
	TableMetadataRow  = 0
	TableFieldTagRow  = 1
	TableFieldNameRow = 2
	TableFieldTypeRow = 3
	TableDataStartRow = 5

	TableFieldIndexCol = 0
)

var ErrSkipTable = fmt.Errorf("skip table")

type TableData struct {
	Name       string
	Metadata   *TableMetadata
	Columns    int
	FieldTags  [][]string
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
		fieldTags   = make([][]string, 0, columns)
		fieldNames  = make([]string, 0, columns)
		fieldTypes  = make([]string, 0, columns)
	)

	for col := 0; col < columns; col++ {
		fieldName := csvData[TableFieldNameRow][col]
		if fieldName == "" || strings.HasPrefix(fieldName, "#") {
			dropColumns = append(dropColumns, col)
		} else {
			tags := make([]string, 0)
			if tagsCell := csvData[TableFieldTagRow][col]; tagsCell != "" {
				tags = strings.Split(tagsCell, ",")
			}
			fieldTags = append(fieldTags, tags)
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
	idxType := FieldType(fieldTypes[TableFieldIndexCol])
	if !idxType.isValidIndexType() {
		return nil, fmt.Errorf("invalid index field type: %s, %s, %s", tableName, idxName, idxType)
	}

	table := &TableData{
		Name:       tableName,
		Columns:    columns,
		FieldTags:  fieldTags,
		FieldNames: fieldNames,
		FieldTypes: fieldTypes,
		DataRows:   dataRows,
	}

	metadataQuery := TableMetadataQuery(csvData[TableMetadataRow][0])
	metadata, err := metadataQuery.Decode()
	if err != nil {
		return nil, fmt.Errorf("failed to decode table metadata query: %s, %w", name, err)
	}
	if err := metadata.Validate(table); err != nil {
		return nil, fmt.Errorf("invalid table metadata: %s, %w", name, err)
	}
	table.Metadata = metadata
	return table, nil
}

func (d *TableData) CSV() [][]string {
	return append([][]string{d.FieldNames, d.FieldTypes}, d.DataRows...)
}
