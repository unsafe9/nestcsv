package nestcsv

type Datasource interface {
	Collect(out chan<- *TableData) error
}

type DatasourceConfig struct {
	SpreadsheetGAS *DatasourceSpreadsheetGAS `yaml:"spreadsheet_gas,omitempty"`
	Excel          *DatasourceExcel          `yaml:"excel,omitempty"`
	CSV            *DatasourceCSV            `yaml:"csv,omitempty"`
}

func (c *DatasourceConfig) List() []Datasource {
	return collectStructFieldsImplementing[Datasource](c)
}
