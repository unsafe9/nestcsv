package main

import (
	"flag"
	"github.com/unsafe9/nestcsv"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
)

func main() {
	var (
		configPath  string
		commandArgs string
	)
	flag.StringVar(&configPath, "c", "nestcsv.yaml", "config file path")
	flag.StringVar(&commandArgs, "a", "", "command arguments")
	flag.Parse()

	nestcsv.SetCommandArgs(commandArgs)

	config, err := nestcsv.ParseConfig(configPath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	out := make(chan *nestcsv.TableData, 1000)

	go func() {
		defer close(out)

		var wg errgroup.Group
		for _, datasource := range config.Datasources {
			wg.Go(func() error {
				return datasource.Collect(out)
			})
		}
		if err := wg.Wait(); err != nil {
			log.Fatalf("failed to collect table: %v", err)
		}
	}()

	var (
		tableDatas []*nestcsv.TableData
		mu         sync.Mutex
	)

	var wg errgroup.Group
	for tableData := range out {
		wg.Go(func() error {
			for _, output := range config.Outputs {
				if err := output.Encode(tableData); err != nil {
					return err
				}
			}

			mu.Lock()
			tableDatas = append(tableDatas, tableData)
			mu.Unlock()
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		log.Fatalf("failed to save table: %v", err)
	}

	for _, codegen := range config.Codegens {
		wg.Go(func() error {
			return codegen.Generate(tableDatas)
		})
	}
	if err := wg.Wait(); err != nil {
		log.Fatalf("failed to generate code: %v", err)
	}
}
