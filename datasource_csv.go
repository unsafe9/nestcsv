package nestcsv

import (
	"encoding/csv"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
	"strings"
)

type CSVOption struct {
	Directories []string `yaml:"directories"`
	Files       []string `yaml:"files"`
}

func CollectCSVFiles(out chan<- TableData, option *CSVOption) error {
	ch := make(chan string, 1000)
	go func() {
		for path := range walkFiles(option.Directories, option.Files, []string{"csv"}) {
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

			out <- NewTableData(path, rows)
			return nil
		})
	}
	return wg.Wait()
}
