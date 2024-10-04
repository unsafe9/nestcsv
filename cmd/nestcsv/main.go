package main

import (
	"flag"
	"github.com/unsafe9/nestcsv"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "nestcsv.yaml", "config file path")
	flag.Parse()

	config, err := nestcsv.ParseConfig(configPath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	out := make(chan *nestcsv.TableData, 1000)

	go func() {
		defer close(out)

		var wg errgroup.Group
		for _, datasource := range config.Datasource.List() {
			wg.Go(func() error {
				return datasource.Collect(out)
			})
		}
		if err := wg.Wait(); err != nil {
			log.Fatalf("failed to collect table: %v", err)
		}
	}()

	var (
		tables []*nestcsv.Table
		mu     sync.Mutex
	)

	var wg errgroup.Group
	for tableData := range out {
		wg.Go(func() error {
			table, err := nestcsv.ParseTable(tableData)
			if err != nil {
				return err
			}
			if err := config.Output.Encode(table); err != nil {
				return err
			}

			mu.Lock()
			tables = append(tables, table)
			mu.Unlock()
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		log.Fatalf("failed to save table: %v", err)
	}

	code, err := nestcsv.AnalyzeTableCode(tables)
	if err != nil {
		log.Fatalf("failed to analyze tables: %v", err)
	}

	for _, codegen := range config.Codegen.List() {
		wg.Go(func() error {
			return codegen.Generate(code)
		})
	}
	if err := wg.Wait(); err != nil {
		log.Fatalf("failed to generate code: %v", err)
	}
}
