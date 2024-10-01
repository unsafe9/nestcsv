package main

import (
	"flag"
	"fmt"
	"github.com/unsafe9/nestcsv"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Datasource struct {
		SpreadsheetGAS *nestcsv.GASOption   `yaml:"spreadsheet_gas,omitempty"`
		Excel          *nestcsv.ExcelOption `yaml:"excel,omitempty"`
		CSV            *nestcsv.CSVOption   `yaml:"csv,omitempty"`
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

	out := make(chan *nestcsv.TableData, 1000)

	go func() {
		defer close(out)

		var wg errgroup.Group
		if config.Datasource.SpreadsheetGAS != nil {
			wg.Go(func() error {
				return nestcsv.CollectSpreadsheetsThroughGAS(out, config.Datasource.SpreadsheetGAS)
			})
		}
		if config.Datasource.Excel != nil {
			wg.Go(func() error {
				return nestcsv.CollectExcelFiles(out, config.Datasource.Excel)
			})
		}
		if config.Datasource.CSV != nil {
			wg.Go(func() error {
				return nestcsv.CollectCSVFiles(out, config.Datasource.CSV)
			})
		}
		if err := wg.Wait(); err != nil {
			log.Fatalf("failed to collect data: %v", err)
		}
	}()

	var wg errgroup.Group
	for tableData := range out {
		wg.Go(func() error {
			table, err := nestcsv.ParseTable(tableData)
			if err != nil {
				return err
			}
			return table.SaveAsJson(config.Output)
		})
	}
	if err := wg.Wait(); err != nil {
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
