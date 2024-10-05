package nestcsv

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"iter"
	"sort"
	"strings"
)

type Codegen interface {
	Generate(*Code) error
}

type CodegenConfig struct {
	CodegenUnion `yaml:",inline"`
}

func (c *CodegenConfig) Generate(code *Code) error {
	return c.loaded.Generate(code)
}

type CodegenUnion struct {
	loaded Codegen     `yaml:"-"`
	Go     *CodegenGo  `yaml:"go,omitempty"`
	UE5    *CodegenUE5 `yaml:"ue5,omitempty"`
}

func (u *CodegenUnion) UnmarshalYAML(node *yaml.Node) error {
	type wrapped CodegenUnion
	if err := node.Decode((*wrapped)(u)); err != nil {
		return err
	}
	list := collectStructFieldsImplementing[Codegen](u)
	if len(list) != 1 {
		return fmt.Errorf("expected exactly one codegen, got %d", len(list))
	}
	u.loaded = list[0]
	return nil
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

func (a *codeAnalyzer) buildStruct(file *CodeFile, table *Table, name string, fields []*TableField) (*CodeStruct, error) {
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
		file.FieldTypes = appendUnique(file.FieldTypes, field.Type)

		if field.Type == FieldTypeStruct {
			id := field.Identifier()
			if name, ok := table.Metadata.StructTypes[id]; ok {
				refFile, err := a.getOrAddNamedStructFile(table, name, field)
				if err != nil {
					return nil, err
				}
				file.FileRefs = appendUnique(file.FileRefs, refFile)
				codeField.StructRef = refFile.Struct

			} else {
				structRef, err := a.getOrAddAnonymousStruct(file, table, field)
				if err != nil {
					return nil, err
				}
				codeField.StructRef = structRef
			}
		}
	}

	return codeStruct, nil
}

func (a *codeAnalyzer) getOrAddAnonymousStruct(file *CodeFile, table *Table, field *TableField) (*CodeStruct, error) {
	name := file.Name + "_" + strings.ReplaceAll(field.Identifier(), ".", "_")
	if field.IsArray() {
		name = singular(name)
	}

	existing := findPtr(file.AnonymousStructs, func(s *CodeStruct) bool {
		return s.Name == name
	})
	if existing != nil {
		return existing, nil
	}

	codeStruct, err := a.buildStruct(file, table, name, field.StructFields)
	if err != nil {
		return nil, err
	}

	file.AnonymousStructs = append(file.AnonymousStructs, codeStruct)
	return codeStruct, nil
}

func (a *codeAnalyzer) getOrAddNamedStructFile(table *Table, name string, field *TableField) (*CodeFile, error) {
	field = field.Clone()
	if fileField, ok := a.namedStructFileFields[name]; ok {
		if !field.StructFieldsEqual(fileField) {
			return nil, fmt.Errorf("named struct %q has different fields", name)
		}
		return a.namedStructFiles[name], nil
	}

	for _, f := range field.StructFields {
		f.ParentField = nil
	}

	file, err := a.buildFile(table, name, field.StructFields)
	if err != nil {
		return nil, err
	}

	a.namedStructFileFields[name] = field
	a.namedStructFiles[name] = file
	return file, nil
}

func (a *codeAnalyzer) addTableFile(table *Table) (*CodeFile, error) {
	file, err := a.buildFile(table, table.Name, table.Fields)
	if err != nil {
		return nil, err
	}

	a.tableFiles[table.Name] = file
	return file, nil
}

func (a *codeAnalyzer) buildFile(table *Table, name string, fields []*TableField) (*CodeFile, error) {
	file := &CodeFile{
		Table: table,
		Name:  name,
	}
	fileStruct, err := a.buildStruct(file, table, name, fields)
	if err != nil {
		return nil, err
	}
	file.Struct = fileStruct

	return file, nil
}

func AnalyzeTableCode(tables []*Table) (*Code, error) {
	a := &codeAnalyzer{
		namedStructFileFields: make(map[string]*TableField),
		namedStructFiles:      make(map[string]*CodeFile),
		tableFiles:            make(map[string]*CodeFile),
	}
	for _, table := range tables {
		if _, err := a.addTableFile(table); err != nil {
			return nil, err
		}
	}
	code := &Code{
		Tables:       make([]*CodeFile, 0, len(a.tableFiles)),
		NamedStructs: make([]*CodeFile, 0, len(a.namedStructFiles)),
	}
	for _, file := range a.tableFiles {
		code.Tables = append(code.Tables, file)
	}
	sort.Slice(code.Tables, func(i, j int) bool {
		return code.Tables[i].Name < code.Tables[j].Name
	})
	for _, file := range a.namedStructFiles {
		code.NamedStructs = append(code.NamedStructs, file)
	}
	sort.Slice(code.NamedStructs, func(i, j int) bool {
		return code.NamedStructs[i].Name < code.NamedStructs[j].Name
	})
	return code, nil
}
