// Code generated by "nestcsv"; DO NOT EDIT.

package table

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ComplexA struct {
	SKU2 SKU `json:"SKU2"`
}

type Complex struct {
	ID      int      `json:"ID"`
	SKU     []SKU    `json:"SKU"`
	Rewards []Reward `json:"Rewards"`
	A       ComplexA `json:"A"`
}

type ComplexTable struct {
	Rows []Complex
}

func (t *ComplexTable) SheetName() string {
	return "complex"
}

func (t *ComplexTable) GetRows() interface{} {
	return t.Rows
}

func (t *ComplexTable) Find(id int) (*Complex, bool) {
	for _, row := range t.Rows {
		if row.ID == id {
			return &row, true
		}
	}
	return nil, false
}

func (t *ComplexTable) Load(data []byte) error {
	return json.Unmarshal(data, &t.Rows)
}

func (t *ComplexTable) LoadFromString(jsonString string) error {
	return t.Load([]byte(jsonString))
}

func (t *ComplexTable) LoadFromFile(basePath string) error {
	file, err := os.Open(filepath.Join(basePath, "complex.json"))
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&t.Rows)
}
