// Code generated by "nestcsv"; DO NOT EDIT.

package table

import (
	"encoding/json"
	"os"
	"time"
)

type Types struct {
	Int         int         `json:"Int"`
	Long        int64       `json:"Long"`
	Float       float64     `json:"Float"`
	String      string      `json:"String"`
	Time        time.Time   `json:"Time"`
	Json        interface{} `json:"Json"`
	IntArray    []int       `json:"IntArray"`
	LongArray   []int64     `json:"LongArray"`
	FloatArray  []float64   `json:"FloatArray"`
	StringArray []string    `json:"StringArray"`
	TimeArray   []time.Time `json:"TimeArray"`
}

type TypesTable struct {
	Rows map[string]Types
}

func (t *TypesTable) SheetName() string {
	return "types"
}

func (t *TypesTable) Load() error {
	file, err := os.Open("../json/server/types.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&t.Rows)
}
