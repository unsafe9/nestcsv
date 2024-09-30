package datasource

import (
	"fmt"
	"github.com/unsafe9/nestcsv/internal"
	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/errgroup"
	"path/filepath"
	"strings"
)

type ExcelOption struct {
	Directories  []string `yaml:"directories"`
	Files        []string `yaml:"files"`
	Extensions   []string `yaml:"extensions"`
	DebugSaveDir *string  `yaml:"debug_save_dir,omitempty"`
}

func CollectExcelFiles(out chan<- TableData, option *ExcelOption) error {
	if len(option.Extensions) == 0 {
		option.Extensions = []string{"xlsx", "xlsm", "xlsb", "xls"}
	}

	ch := make(chan string, 1000)
	go func() {
		for path := range internal.WalkFiles(option.Directories, option.Files, option.Extensions) {
			if strings.HasPrefix(filepath.Base(path), "#") {
				continue
			}
			ch <- path
		}
		close(ch)
	}()

	var wg errgroup.Group
	for path := range ch {
		if strings.HasPrefix(filepath.Base(path), "#") {
			continue
		}

		wg.Go(func() error {
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

				csv := NewTableData(sheet, rows)
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
	return wg.Wait()
}
