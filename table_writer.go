package nestcsv

import "gopkg.in/yaml.v3"

type TableWriter interface {
	Write(name string, value any) error
}

type OutputConfig struct {
	When *When    `yaml:"when,omitempty"`
	Tags []string `yaml:"tags"`

	exclusiveConfigGroup[TableWriter]
	JSON *TableWriterJSON `yaml:"json,omitempty"`
	Bin  *TableWriterBin  `yaml:"bin,omitempty"`
}

func (c *OutputConfig) Write(tableData *TableData) error {
	tableParser := NewTableParser(tableData)
	tableFields, err := tableParser.ParseTableFields(c.Tags)
	if err != nil {
		return err
	}
	if len(tableFields) == 0 {
		return nil
	}
	value, err := tableParser.Marshal(tableFields)
	if err != nil {
		return err
	}
	return c.loaded.Write(tableData.Name, value)
}

func (c *OutputConfig) UnmarshalYAML(node *yaml.Node) error {
	type wrapped OutputConfig
	if err := node.Decode((*wrapped)(c)); err != nil {
		return err
	}
	return c.postUnmarshalYAML(c)
}
