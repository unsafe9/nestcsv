package nestcsv

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

const (
	FileTypeJSON = "json"
	FileTypeBin  = "bin"
)

type TableEncoder struct {
	RootDir  string `yaml:"root_dir"`
	Indent   string `yaml:"indent"`
	FileType string `yaml:"file_type"`
}

func (e *TableEncoder) Encode(table *Table) error {
	if e.RootDir == "" {
		e.RootDir = "."
	}
	if e.FileType == "" {
		e.FileType = FileTypeJSON
	}

	if e.FileType == FileTypeJSON {
		return e.saveAsJson(table)

	} else if e.FileType == FileTypeBin {
		return e.saveAsBin(table)
	}

	return fmt.Errorf("invalid file type: %s", e.FileType)
}

func (e *TableEncoder) saveAsJson(table *Table) error {
	jsonBytes, err := json.MarshalIndent(table.Marshal(), "", e.Indent)
	if err != nil {
		return err
	}

	file, err := createFile(e.RootDir, table.Name, "json")
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(jsonBytes); err != nil {
		return err
	}
	return nil
}

func (e *TableEncoder) saveAsBin(table *Table) error {
	jsonBytes, err := json.Marshal(table.Marshal())
	if err != nil {
		return err
	}

	file, err := createFile(e.RootDir, table.Name, "bin")
	if err != nil {
		return err
	}
	defer file.Close()

	binHeader := make([]byte, 4)
	binary.BigEndian.PutUint32(binHeader, uint32(len(jsonBytes)))
	if _, err := file.Write(binHeader); err != nil {
		return err
	}
	if _, err := file.Write(jsonBytes); err != nil {
		return err
	}
	return nil
}
