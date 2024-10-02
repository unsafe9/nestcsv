package nestcsv

import (
	"encoding/csv"
	"errors"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
	"strings"
)

type DatasourceCSV struct {
	Directories []string `yaml:"directories"`
	Files       []string `yaml:"files"`
}

func (d *DatasourceCSV) Collect(out chan<- *TableData) error {
	ch := make(chan string, 1000)
	go func() {
		for path := range walkFiles(d.Directories, d.Files, []string{"csv"}) {
			if strings.HasPrefix(filepath.Base(path), "#") {
				continue
			}
			ch <- path
		}
		close(ch)
	}()

	var wg errgroup.Group
	for path := range ch {
		wg.Go(func() error {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			rows, err := csv.NewReader(file).ReadAll()
			if err != nil {
				return err
			}

			tableData, err := ParseTableData(path, rows)
			if err != nil {
				if errors.Is(err, ErrSkipTable) {
					return nil
				}
				return err
			}

			out <- tableData
			return nil
		})
	}
	return wg.Wait()
}
