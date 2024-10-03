package nestcsv

import (
	"bytes"
	"go/format"
	"path/filepath"
	"strings"
	"text/template"
)

type CodegenGo struct {
	RootDir      string `yaml:"root_dir"`
	PackageName  string `yaml:"package_name"`
	DataLoadPath string `yaml:"data_load_path"`
}

type codegenGoStruct struct {
	name         string
	fields       []*TableField
	namedStructs map[string]string
}

func (c *CodegenGo) Generate(tables []*Table) error {
	if c.PackageName == "" {
		c.PackageName = filepath.Base(c.RootDir)
	}

	namedStructs := make(map[string]*codegenGoStruct)
	for _, table := range tables {
		cs := &codegenGoStruct{
			name:         table.Name,
			fields:       table.Fields,
			namedStructs: table.Metadata.StructTypes,
		}
		cs.collectNamedStructs(namedStructs)
		imports := make([]string, 0)
		if cs.includeTime() {
			imports = append(imports, "time")
		}
		funcs := template.FuncMap{
			"renderStruct": cs.renderStruct,
		}
		values := map[string]any{
			"PackageName":  c.PackageName,
			"Table":        table,
			"Imports":      imports,
			"DataLoadPath": filepath.Join(c.DataLoadPath, table.Name+".json"),
		}
		err := c.template(table.Name+".go", "table.go.tpl", funcs, values)
		if err != nil {
			return err
		}
	}

	for name, cs := range namedStructs {
		funcs := template.FuncMap{
			"renderStruct": cs.renderStruct,
		}
		values := map[string]any{
			"PackageName": c.PackageName,
			"Name":        name,
			"Fields":      cs.fields,
		}
		err := c.template(name+".go", "struct.go.tpl", funcs, values)
		if err != nil {
			return err
		}
	}

	values := map[string]any{
		"PackageName": c.PackageName,
		"Tables":      tables,
	}
	err := c.template("loader.go", "loader.go.tpl", nil, values)
	if err != nil {
		return err
	}

	return nil
}

func (c *CodegenGo) template(fileName, templateName string, funcs template.FuncMap, values any) error {
	var buf bytes.Buffer
	err := templateFile(&buf, "go/"+templateName, funcs, values)
	if err != nil {
		return err
	}

	fileBytes, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	file, err := createFile(c.RootDir, strings.ToLower(fileName), "go")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(fileBytes)
	if err != nil {
		return err
	}
	return nil
}

func (c *codegenGoStruct) includeTime() bool {
	for _, field := range c.fields {
		for f := range field.Iterate() {
			if f.Type == FieldTypeTime {
				return true
			}
		}
	}
	return false
}

func (c *codegenGoStruct) collectNamedStructs(m map[string]*codegenGoStruct) {
	for _, topField := range c.fields {
		for f := range topField.Iterate() {
			if f.Type == FieldTypeStruct {
				id := f.Identifier()
				if name, ok := c.namedStructs[id]; ok {
					m[name] = &codegenGoStruct{
						name:         name,
						fields:       f.StructFields,
						namedStructs: c.namedStructs,
					}
				}
			}
		}
	}
}

func (c *codegenGoStruct) renderStruct(name string, fields []*TableField) string {
	var buf bytes.Buffer
	if name != "" {
		buf.WriteString("type " + pascal(name) + " ")
	}
	buf.WriteString("struct {\n")
	for _, f := range fields {
		buf.WriteString(c.renderField(f) + "\n")
	}
	buf.WriteString("}")
	return buf.String()
}

func (c *codegenGoStruct) renderFieldType(f *TableField) string {
	switch f.Type {
	case FieldTypeInt:
		return "int"
	case FieldTypeLong:
		return "int64"
	case FieldTypeFloat:
		return "float64"
	case FieldTypeBool:
		return "bool"
	case FieldTypeString:
		return "string"
	case FieldTypeTime:
		return "time.Time"
	case FieldTypeJSON:
		return "interface{}"
	case FieldTypeStruct:
		if name, ok := c.namedStructs[f.Identifier()]; ok {
			return pascal(name)
		} else {
			return c.renderStruct("", f.StructFields)
		}

	}
	panic("unsupported type " + f.Type)
}

func (c *codegenGoStruct) renderField(f *TableField) string {
	fieldType := c.renderFieldType(f)
	if f.IsArray() {
		if f.Type == "json" {
			panic("unsupported type: json array")
		}
		fieldType = "[]" + fieldType
	}
	return pascal(f.Name) + " " + fieldType + "`json:\"" + f.Name + "\"`"
}