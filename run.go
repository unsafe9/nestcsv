package nestcsv

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
)

func Generate(config *Config) error {
	out := make(chan *TableData, 1000)
	errStop := make(chan error, 1)

	go func() {
		defer close(out)

		var wg errgroup.Group
		for _, datasource := range config.Datasources {
			wg.Go(func() error {
				return datasource.Collect(out)
			})
		}
		if err := wg.Wait(); err != nil {
			errStop <- fmt.Errorf("collect datasource: %w", err)
			return
		}
	}()

	go func() {
		var (
			tableDatas []*TableData
			mu         sync.Mutex
			wg         errgroup.Group
		)
		for tableData := range out {
			wg.Go(func() error {
				for _, output := range config.Outputs {
					if err := output.Write(tableData); err != nil {
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
			errStop <- fmt.Errorf("failed to write output: %w", err)
			return
		}

		if len(tableDatas) > 0 {
			var wg errgroup.Group
			for _, codegen := range config.Codegens {
				wg.Go(func() error {
					return codegen.Generate(tableDatas)
				})
			}
			if err := wg.Wait(); err != nil {
				errStop <- fmt.Errorf("failed to generate code: %w", err)
				return
			}
		}
		errStop <- nil
	}()
	return <-errStop
}
