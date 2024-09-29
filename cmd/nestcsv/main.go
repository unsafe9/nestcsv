package main

import (
	"flag"
	"fmt"
	"github.com/unsafe9/nestcsv"
	"github.com/unsafe9/nestcsv/datasource"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Datasource struct {
		SpreadsheetGAS *datasource.GASOption       `yaml:"spreadsheet_gas,omitempty"`
		Excel          *datasource.ExcelOption     `yaml:"excel,omitempty"`
		LocalFile      *datasource.LocalFileOption `yaml:"local_file,omitempty"`
	} `yaml:"datasource"`

	Output *nestcsv.TableSaveOption `yaml:"output"`

	Codegen struct {
	} `yaml:"codegen"`
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "nestcsv.yaml", "config file path")
	flag.Parse()

	config, err := parseConfig(configPath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	out := make(chan datasource.CSV, 1000)

	go func() {
		defer close(out)

		var datasourceWaitGroup errgroup.Group
		if config.Datasource.SpreadsheetGAS != nil {
			datasourceWaitGroup.Go(func() error {
				return datasource.CollectSpreadsheetsThroughGAS(out, config.Datasource.SpreadsheetGAS)
			})
		}
		if config.Datasource.Excel != nil {
			datasourceWaitGroup.Go(func() error {
				return datasource.CollectExcelFiles(out, config.Datasource.Excel)
			})
		}
		if config.Datasource.LocalFile != nil {
			datasourceWaitGroup.Go(func() error {
				return datasource.CollectLocalFiles(out, config.Datasource.LocalFile)
			})
		}
		if err := datasourceWaitGroup.Wait(); err != nil {
			log.Fatalf("failed to collect data: %v", err)
		}
	}()

	var tableWaitGroup errgroup.Group
	for csv := range out {
		tableWaitGroup.Go(func() error {
			table, err := nestcsv.ParseTable(csv.Name, csv.Rows)
			if err != nil {
				return err
			}
			return table.SaveAsJson(config.Output)
		})
	}
	if err := tableWaitGroup.Wait(); err != nil {
		log.Fatalf("failed to save table: %v", err)
	}
}

func parseConfig(configPath string) (*Config, error) {
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
