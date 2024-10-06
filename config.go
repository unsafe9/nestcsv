package nestcsv

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"slices"
	"strings"
)

var commandArgs []string

type Config struct {
	Datasources []DatasourceConfig `yaml:"datasources"`
	Outputs     []TableEncoder     `yaml:"outputs"`
	Codegens    []CodegenConfig    `yaml:"codegens"`
}

func ParseConfig(configPath string) (*Config, error) {
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
		return d.When == nil || d.When.Match()
	})
	config.Outputs = filter(config.Outputs, func(e TableEncoder) bool {
		return e.When == nil || e.When.Match()
	})
	config.Codegens = filter(config.Codegens, func(c CodegenConfig) bool {
		return c.When == nil || c.When.Match()
	})

	return &config, nil
}

func SetCommandArgs(argsStr string) {
	commandArgs = strings.Split(argsStr, " ")
	for i := 0; i < len(commandArgs); i++ {
		commandArgs[i] = strings.TrimSpace(commandArgs[i])
	}
}

type When struct {
	Env  map[string]string `yaml:"env,omitempty"`
	Args []string          `yaml:"args,omitempty"`
}

func (w *When) Match() bool {
	for key, value := range w.Env {
		if os.Getenv(key) != value {
			return false
		}
	}
	for _, arg := range w.Args {
		if !slices.Contains(commandArgs, arg) {
			return false
		}
	}
	return true
}
