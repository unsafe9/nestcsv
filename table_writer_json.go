package nestcsv

import (
	"encoding/json"
)

type TableWriterJSON struct {
	RootDir string `yaml:"root_dir"`
	Indent  string `yaml:"indent"`
}

func (e *TableWriterJSON) Write(name string, value any) error {
	if e.RootDir == "" {
		e.RootDir = "."
	}
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
