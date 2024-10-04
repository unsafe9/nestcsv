package nestcsv

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/errgroup"
	"path/filepath"
	"strings"
)

type DatasourceExcel struct {
	Patterns     []string `yaml:"patterns"`
	DebugSaveDir *string  `yaml:"debug_save_dir,omitempty"`
}

func (d *DatasourceExcel) Collect(out chan<- *TableData) error {
	ch := make(chan string, 1000)
	go func() {
		for path := range glob(d.Patterns) {
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
				if d.DebugSaveDir != nil {
					if err := saveCSVFile(*d.DebugSaveDir, sheet, rows); err != nil {
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
