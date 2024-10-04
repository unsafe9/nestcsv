package nestcsv

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

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

	return &config, nil
}
