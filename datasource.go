package nestcsv

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Datasource interface {
	Collect(out chan<- *TableData) error
}

type DatasourceConfig struct {
	DatasourceUnion `yaml:",inline"`
}

func (c *DatasourceConfig) Collect(out chan<- *TableData) error {
	return c.loaded.Collect(out)
}

type DatasourceUnion struct {
	loaded         Datasource                `yaml:"-"`
	SpreadsheetGAS *DatasourceSpreadsheetGAS `yaml:"spreadsheet_gas,omitempty"`
	Excel          *DatasourceExcel          `yaml:"excel,omitempty"`
	CSV            *DatasourceCSV            `yaml:"csv,omitempty"`
}

func (u *DatasourceUnion) UnmarshalYAML(node *yaml.Node) error {
	type wrapped DatasourceUnion
	if err := node.Decode((*wrapped)(u)); err != nil {
		return err
	}
	list := collectStructFieldsImplementing[Datasource](u)
	if len(list) != 1 {
		return fmt.Errorf("expected exactly one datasource, got %d", len(list))
	}
	u.loaded = list[0]
	return nil
}
