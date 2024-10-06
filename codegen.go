package nestcsv

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Codegen interface {
	Generate(code *Code) error
}

type CodegenConfig struct {
	When *When    `yaml:"when,omitempty"`
	Tags []string `yaml:"tags"`

	loaded Codegen     `yaml:"-"`
	Go     *CodegenGo  `yaml:"go,omitempty"`
	UE5    *CodegenUE5 `yaml:"ue5,omitempty"`
}

func (c *CodegenConfig) Generate(tableDatas []*TableData) error {
	code, err := AnalyzeTableCode(tableDatas, c.Tags)
	if err != nil {
		return err
	}
	return c.loaded.Generate(code)
}

func (c *CodegenConfig) UnmarshalYAML(node *yaml.Node) error {
	type wrapped CodegenConfig
	if err := node.Decode((*wrapped)(c)); err != nil {
		return err
	}
	list := collectStructFieldsImplementing[Codegen](c)
	if len(list) != 1 {
		return fmt.Errorf("expected exactly one codegen, got %d", len(list))
	}
	c.loaded = list[0]
	return nil
}
