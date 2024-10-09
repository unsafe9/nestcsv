package nestcsv

import (
	"encoding/binary"
	"encoding/json"
)

type TableWriterBin struct {
	RootDir string `yaml:"root_dir"`
}

func (e *TableWriterBin) Write(name string, value any) error {
	if e.RootDir == "" {
		e.RootDir = "."
	}
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
