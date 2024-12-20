package nestcsv

import (
	"fmt"
	"sort"
	"strings"
)

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
	IsTable          bool
	IsMap            bool
	Name             string
	Struct           *CodeStruct
	AnonymousStructs []*CodeStruct
	FileRefs         []*CodeFile
	FieldTypes       []FieldType
	IDField          *CodeStructField
	IDFieldType      FieldType // this will be set even if IDField is nil
}

type Code struct {
	Tables       []*CodeFile
	NamedStructs []*CodeFile
}

func (c *Code) Files(yield func(*CodeFile) bool) {
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

type codeAnalyzer struct {
	namedStructFileFields map[string]*TableField
	namedStructFiles      map[string]*CodeFile
	tableFiles            map[string]*CodeFile
}

type codeAnalyzerTable struct {
	name        string
	metadata    *TableMetadata
	fields      []*TableField
	idField     *TableField
	idFieldType FieldType
}

func (a *codeAnalyzer) buildStruct(file *CodeFile, table *codeAnalyzerTable, name string, fields []*TableField) (*CodeStruct, error) {
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
			if name, ok := table.metadata.Structs.Get(id); ok {
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

func (a *codeAnalyzer) getOrAddAnonymousStruct(file *CodeFile, table *codeAnalyzerTable, field *TableField) (*CodeStruct, error) {
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

func (a *codeAnalyzer) getOrAddNamedStructFile(table *codeAnalyzerTable, name string, field *TableField) (*CodeFile, error) {
	field = field.Clone()
	if fileField, ok := a.namedStructFileFields[name]; ok {
		if !field.StructEqual(fileField) {
			return nil, fmt.Errorf("named struct %q has different fields", name)
		}
		return a.namedStructFiles[name], nil
	}

	file := &CodeFile{
		Name: name,
	}
	for _, f := range field.StructFields {
		f.ParentField = nil
	}
	fileStruct, err := a.buildStruct(file, table, name, field.StructFields)
	if err != nil {
		return nil, err
	}
	file.Struct = fileStruct

	a.namedStructFileFields[name] = field
	a.namedStructFiles[name] = file
	return file, nil
}

func (a *codeAnalyzer) addTableFile(table *codeAnalyzerTable) (*CodeFile, error) {
	file := &CodeFile{
		IsTable:     true,
		IsMap:       table.metadata.AsMap,
		Name:        table.name,
		IDFieldType: table.idFieldType,
	}
	fileStruct, err := a.buildStruct(file, table, table.name, table.fields)
	if err != nil {
		return nil, err
	}
	file.Struct = fileStruct

	if table.idField != nil {
		file.IDField = fileStruct.Fields[0]
	}

	a.tableFiles[table.name] = file
	return file, nil
}

func AnalyzeTableCode(tableDatas []*TableData, tags []string) (*Code, error) {
	tables := make([]*codeAnalyzerTable, 0, len(tableDatas))
	for _, tableData := range tableDatas {
		fields, err := NewTableParser(tableData).ParseTableFields(tags)
		if err != nil {
			return nil, err
		}
		if len(fields) == 0 {
			continue
		}
		var (
			idField     *TableField
			idFieldType FieldType
		)
		if tableData.FieldNames[TableFieldIndexCol] == fields[TableFieldIndexCol].Name {
			idField = fields[TableFieldIndexCol]
			idFieldType = idField.Type
		} else {
			idFieldType, _ = newFieldType(tableData.FieldTypes[TableFieldIndexCol])
		}
		tables = append(tables, &codeAnalyzerTable{
			name:        tableData.Name,
			metadata:    tableData.Metadata,
			fields:      fields,
			idField:     idField,
			idFieldType: idFieldType,
		})
	}

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
