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
	Tags     []string `yaml:"tags"`
	RootDir  string   `yaml:"root_dir"`
	Indent   string   `yaml:"indent"`
	FileType string   `yaml:"file_type"`
}

func (e *TableEncoder) Encode(tableData *TableData) error {
	tableParser := NewTableParser(tableData)
	tableFields, err := tableParser.ParseTableFields(e.Tags)
	if err != nil {
		return err
	}
	if len(tableFields) == 0 {
		return nil
	}
	value, err := tableParser.Marshal(tableFields)
	if err != nil {
		return err
	}

	if e.RootDir == "" {
		e.RootDir = "."
	}
	if e.FileType == "" {
		e.FileType = FileTypeJSON
	}

	if e.FileType == FileTypeJSON {
		return e.saveAsJson(tableData.Name, value)

	} else if e.FileType == FileTypeBin {
		return e.saveAsBin(tableData.Name, value)
	}

	return fmt.Errorf("invalid file type: %s", e.FileType)
}

func (e *TableEncoder) saveAsJson(name string, value any) error {
	var (
		jsonBytes []byte
		err       error
	)
	if e.Indent == "" {
		jsonBytes, err = json.Marshal(value)
	} else {
		jsonBytes, err = json.MarshalIndent(value, "", e.Indent)
	}
	if err != nil {
		return err
	}

	file, err := createFile(e.RootDir, name, "json")
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(jsonBytes); err != nil {
		return err
	}
	return nil
}

func (e *TableEncoder) saveAsBin(name string, value any) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	file, err := createFile(e.RootDir, name, "bin")
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
