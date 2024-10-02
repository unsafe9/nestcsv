package nestcsv

type Codegen interface {
	Generate([]*Table) error
}

type CodegenConfig struct {
	Go *CodegenGo `yaml:"go,omitempty"`
}

func (c *CodegenConfig) List() []Codegen {
	return collectStructFieldsImplementing[Codegen](c)
}
