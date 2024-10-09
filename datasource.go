package nestcsv

import (
	"gopkg.in/yaml.v3"
)

type Datasource interface {
	Collect(out chan<- *TableData) error
}

type DatasourceConfig struct {
	When *When `yaml:"when,omitempty"`

	exclusiveConfigGroup[Datasource]
	SpreadsheetGAS *DatasourceSpreadsheetGAS `yaml:"spreadsheet_gas,omitempty"`
	Excel          *DatasourceExcel          `yaml:"excel,omitempty"`
	CSV            *DatasourceCSV            `yaml:"csv,omitempty"`
}

func (c *DatasourceConfig) Collect(out chan<- *TableData) error {
	return c.loaded.Collect(out)
}

func (c *DatasourceConfig) UnmarshalYAML(node *yaml.Node) error {
	type wrapped DatasourceConfig
	if err := node.Decode((*wrapped)(c)); err != nil {
		return err
	}
	return c.postUnmarshalYAML(c)
}
