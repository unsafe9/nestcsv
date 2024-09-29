package datasource

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type ExcelOption struct {
	Directories  []string `yaml:"directories"`
	Extensions   []string `yaml:"extensions"`
	Files        []string `yaml:"files"`
	DebugSaveDir *string  `yaml:"debug_save_dir,omitempty"`
}

func CollectExcelFiles(out chan<- CSV, option *ExcelOption) error {
	if len(option.Extensions) == 0 {
		option.Extensions = []string{"xlsx", "xlsm", "xlsb", "xls"}
	}

	ch := make(chan string, 1000)
	go func() {
		visited := make(map[string]struct{})

		for _, dir := range option.Directories {
			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				if len(option.Extensions) > 0 {
					ext := filepath.Ext(path)[1:]
					if !slices.Contains(option.Extensions, ext) {
						return nil
					}
				}

				if _, ok := visited[path]; ok {
					return nil
				}
				visited[path] = struct{}{}

				ch <- path
				return nil
			})
			if err != nil {
				return
			}
		}

		for _, file := range option.Files {
			file = filepath.Clean(file)
			if _, ok := visited[file]; ok {
				continue
			}
			visited[file] = struct{}{}
			ch <- file
		}

		close(ch)
	}()

	var excelParseWaitGroup errgroup.Group
	for path := range ch {
		excelParseWaitGroup.Go(func() error {
			file, err := excelize.OpenFile(path)
			if err != nil {
				return err
			}
			defer file.Close()

			for _, sheet := range file.GetSheetList() {
				if strings.HasPrefix(sheet, "#") {
					continue
				}

				rows, err := file.GetRows(sheet)
				if err != nil {
					return err
				}
				if len(rows) == 0 {
					return fmt.Errorf("no rows in the sheet: %s", sheet)
				}

				// Ensure that all rows have the same length as the header row
				headerLen := len(rows[0])
				for i, row := range rows {
					if len(row) < headerLen {
						rows[i] = append(row, make([]string, headerLen-len(row))...)
					}
				}
			
				csv := NewCSV(sheet, rows)
				if option.DebugSaveDir != nil {
					if err := csv.Save(*option.DebugSaveDir); err != nil {
						return err
					}
				}
				out <- csv
			}
			return nil
		})
	}
	return excelParseWaitGroup.Wait()
}
