package datasource

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CSV struct {
	Name string
	Data []byte
}

func NewCSV(name string, data []byte) CSV {
	return CSV{
		Name: strings.TrimSuffix(filepath.Base(name), ".csv"),
		Data: data,
	}
}

func (c *CSV) Save(rootDir string) error {
	if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create the directory: %s, %w", rootDir, err)
	}
	file, err := os.Create(filepath.Join(rootDir, c.Name+".csv"))
	if err != nil {
		return fmt.Errorf("failed to create the file: %s, %w", c.Name, err)
	}
	defer file.Close()

	if _, err := file.Write(c.Data); err != nil {
		return fmt.Errorf("failed to write the file: %s, %w", c.Name, err)
	}
	return nil
}
