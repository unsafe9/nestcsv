package nestcsv

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Datasource interface {
	Collect(out chan<- *TableData) error
}

type DatasourceConfig struct {
	loaded         Datasource                `yaml:"-"`
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
	list := collectStructFieldsImplementing[Datasource](c)
	if len(list) != 1 {
		return fmt.Errorf("expected exactly one datasource, got %d", len(list))
	}
	c.loaded = list[0]
	return nil
}
