package nestcsv

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"slices"
)

type Config struct {
	Datasources []DatasourceConfig `yaml:"datasources"`
	Outputs     []OutputConfig     `yaml:"outputs"`
	Codegens    []CodegenConfig    `yaml:"codegens"`
}

func ParseConfig(configPath string, args []string) (*Config, error) {

	var config Config
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s, %w", configPath, err)
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode yaml: %s, %w", configPath, err)
	}

	config.Datasources = filter(config.Datasources, func(d DatasourceConfig) bool {
		return d.When == nil || d.When.Match(args)
	})
	config.Outputs = filter(config.Outputs, func(e OutputConfig) bool {
		return e.When == nil || e.When.Match(args)
	})
	config.Codegens = filter(config.Codegens, func(c CodegenConfig) bool {
		return c.When == nil || c.When.Match(args)
	})

	return &config, nil
}

type When struct {
	// Not - if true, the condition is negated
	Not  bool              `yaml:"not"`
	Env  map[string]string `yaml:"env,omitempty"`
	Args []string          `yaml:"args,omitempty"`
}

func (w *When) Match(args []string) bool {
	return w.match(args) != w.Not
}

func (w *When) match(args []string) bool {
	for key, value := range w.Env {
		if os.Getenv(key) != value {
			return false
		}
	}
	for _, arg := range w.Args {
		if !slices.Contains(args, arg) {
			return false
		}
	}
	return true
}

type exclusiveConfigGroup[T any] struct {
	loaded T
}

func (c *exclusiveConfigGroup[T]) postUnmarshalYAML(parent any) error {
	list := collectStructFieldsImplementing[T](parent)
	if len(list) != 1 {
		return fmt.Errorf("expected exactly one %T, got %d", c.loaded, len(list))
	}
	c.loaded = list[0]
	return nil
}
