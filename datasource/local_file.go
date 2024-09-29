package datasource

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

type LocalFileOption struct {
	RootDir string `yaml:"root_dir"`
}

func CollectLocalFiles(out chan<- CSV, option *LocalFileOption) error {
	return filepath.Walk(option.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".csv" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			rows, err := csv.NewReader(file).ReadAll()
			if err != nil {
				return err
			}

			out <- NewCSV(path, rows)
		}
		return nil
	})
}
