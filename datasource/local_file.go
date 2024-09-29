package datasource

import (
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
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			out <- NewCSV(path, data)
		}
		return nil
	})
}
