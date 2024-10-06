// Code generated by "nestcsv"; DO NOT EDIT.

package table

import (
	"encoding/json"
	"os"
)

type ComplexSKU struct {
	Type string `json:"Type"`
	ID   string `json:"ID"`
}

type Complex struct {
	ID      int          `json:"ID"`
	SKU     []ComplexSKU `json:"SKU"`
	Rewards []Reward     `json:"Rewards"`
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

func (t *ComplexTable) Load() error {
	file, err := os.Open("../json/server/complex.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&t.Rows)
}
