package nestcsv

import (
	"errors"
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

func CollectExcelFiles(out chan<- *TableData, option *ExcelOption) error {
	if len(option.Extensions) == 0 {
		option.Extensions = []string{"xlsx", "xlsm", "xlsb", "xls"}
	}

	ch := make(chan string, 1000)
	go func() {
		for path := range walkFiles(option.Directories, option.Files, option.Extensions) {
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

				tableData, err := ParseTableData(sheet, rows)
				if err != nil {
					if errors.Is(err, ErrSkipTable) {
						return nil
					}
					return err
				}
				if option.DebugSaveDir != nil {
					if err := tableData.SaveAsCSV(*option.DebugSaveDir); err != nil {
						return err
					}
				}
				out <- tableData
			}
			return nil
		})
	}
	return wg.Wait()
}
