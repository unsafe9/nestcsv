package datasource

import (
	"os"
	"path/filepath"
)

func CollectLocalFiles(out chan<- CSV, rootDir string) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".csv" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			out <- CSV{
				Name: filepath.Base(path),
				Data: data,
			}
		}
		return nil
	})
}
