package nestcsv

type Table struct {
	Name     string
	Metadata *TableMetadata
	Fields   []*TableField
	Values   []map[string]any
}
