package nestcsv

import (
	"iter"
	"slices"
	"sort"
	"strings"
)

type Codegen interface {
	Generate(*Code) error
}

type CodegenConfig struct {
	Go  *CodegenGo  `yaml:"go,omitempty"`
	UE5 *CodegenUE5 `yaml:"ue5,omitempty"`
}

func (c *CodegenConfig) List() []Codegen {
	return collectStructFieldsImplementing[Codegen](c)
}

type CodeStructField struct {
	Name      string
	Type      FieldType
	IsArray   bool
	StructRef *CodeStruct
}

type CodeStruct struct {
	Name   string
	Fields []*CodeStructField
}

type CodeFile struct {
	Table            *Table // can be nil (for named structs)
	Name             string
	Struct           *CodeStruct
	AnonymousStructs []*CodeStruct
	FileRefs         []*CodeFile
	FieldTypes       []FieldType
}

type Code struct {
	Tables       []*CodeFile
	NamedStructs []*CodeFile
}

func (c *Code) Files() iter.Seq[*CodeFile] {
	return func(yield func(*CodeFile) bool) {
		for _, file := range c.NamedStructs {
			if !yield(file) {
				return
			}
		}
		for _, file := range c.Tables {
			if !yield(file) {
				return
			}
		}
	}
}

type codeAnalyzer struct {
	namedStructFileFields map[string]*TableField
	namedStructFiles      map[string]*CodeFile
	tableFiles            map[string]*CodeFile
}

func (a *codeAnalyzer) buildStruct(file *CodeFile, table *Table, name string, fields []*TableField) *CodeStruct {
	codeStruct := &CodeStruct{
		Name:   name,
		Fields: make([]*CodeStructField, 0, len(fields)),
	}

	for _, field := range fields {
		codeField := &CodeStructField{
			Name:      field.Name,
			Type:      field.Type,
			IsArray:   field.IsArray(),
			StructRef: nil,
		}
		codeStruct.Fields = append(codeStruct.Fields, codeField)
		if !slices.Contains(file.FieldTypes, field.Type) {
			file.FieldTypes = append(file.FieldTypes, field.Type)
		}

		if field.Type == FieldTypeStruct {
			id := field.Identifier()
			if name, ok := table.Metadata.StructTypes[id]; ok {
				refFile := a.getOrAddNamedStructFile(table, name, field)
				if !slices.Contains(file.FileRefs, refFile) {
					file.FileRefs = append(file.FileRefs, refFile)
				}
				codeField.StructRef = refFile.Struct

			} else {
				codeField.StructRef = a.getOrAddAnonymousStruct(file, table, field)
			}
		}
	}

	return codeStruct
}

func (a *codeAnalyzer) getOrAddAnonymousStruct(file *CodeFile, table *Table, field *TableField) *CodeStruct {
	name := file.Name + "_" + strings.ReplaceAll(field.Identifier(), ".", "_")
	if field.IsArray() {
		name = singular(name)
	}

	existing := findPtr(file.AnonymousStructs, func(s *CodeStruct) bool {
		return s.Name == name
	})
	if existing != nil {
		return existing
	}

	codeStruct := a.buildStruct(file, table, name, field.StructFields)
	file.AnonymousStructs = append(file.AnonymousStructs, codeStruct)
	return codeStruct
}

func (a *codeAnalyzer) getOrAddNamedStructFile(table *Table, name string, field *TableField) *CodeFile {
	field = field.Clone()
	if fileField, ok := a.namedStructFileFields[name]; ok {
		if !field.Equal(fileField) {
			// TODO : better error handling
			panic("named struct field mismatch: " + name)
		}
		return a.namedStructFiles[name]
	}

	file := &CodeFile{
		Name: name,
	}
	for _, f := range field.StructFields {
		f.ParentField = nil
	}
	file.Struct = a.buildStruct(file, table, name, field.StructFields)

	a.namedStructFileFields[name] = field
	a.namedStructFiles[name] = file
	return file
}

func (a *codeAnalyzer) addTableFile(table *Table) *CodeFile {
	if _, ok := a.tableFiles[table.Name]; ok {
		panic("table file already exists")
	}

	file := &CodeFile{
		Table: table,
		Name:  table.Name,
	}
	file.Struct = a.buildStruct(file, table, table.Name, table.Fields)
	a.tableFiles[table.Name] = file
	return file
}

func AnalyzeTableCode(tables []*Table) (*Code, error) {
	sort.Slice(tables, func(i, j int) bool {
		return strings.Compare(tables[i].Name, tables[j].Name) < 0
	})

	a := &codeAnalyzer{
		namedStructFileFields: make(map[string]*TableField),
		namedStructFiles:      make(map[string]*CodeFile),
		tableFiles:            make(map[string]*CodeFile),
	}
	for _, table := range tables {
		a.addTableFile(table)
	}
	code := &Code{
		Tables:       make([]*CodeFile, 0, len(a.tableFiles)),
		NamedStructs: make([]*CodeFile, 0, len(a.namedStructFiles)),
	}
	for _, file := range a.tableFiles {
		code.Tables = append(code.Tables, file)
	}
	for _, file := range a.namedStructFiles {
		code.NamedStructs = append(code.NamedStructs, file)
	}
	return code, nil
}
