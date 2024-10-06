package nestcsv

import (
	"path/filepath"
	"text/template"
)

type CodegenUE5 struct {
	RootDir string `yaml:"root_dir"`
	Prefix  string `yaml:"prefix"`
}

func (c *CodegenUE5) Generate(code *Code) error {
	values := map[string]any{
		"Prefix": c.Prefix,
	}
	if err := c.template("TableDataBase.h", "TableDataBase.h.tpl", values); err != nil {
		return err
	}
	if err := c.template("TableBase.h", "TableBase.h.tpl", values); err != nil {
		return err
	}
	values["Tables"] = code.Tables
	if err := c.template("TableHolder.h", "TableHolder.h.tpl", values); err != nil {
		return err
	}

	for file := range code.Files() {
		values := map[string]any{
			"File":   file,
			"Prefix": c.Prefix,
		}
		if err := c.template(pascal(file.Name)+".h", "file.h.tpl", values); err != nil {
			return err
		}

		if file.IsTable {
			if err := c.template(pascal(file.Name)+"Table.h", "table.h.tpl", values); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *CodegenUE5) template(fileName, templateName string, values any) error {
	file, err := createFile(c.RootDir, c.Prefix+fileName, filepath.Ext(fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.
		New(filepath.Base(templateName)).
		Funcs(templateFuncMap).
		Funcs(template.FuncMap{
			"fieldType":     c.fieldType,
			"fieldElemType": c.fieldElemType,
		}).
		ParseFS(templateFS, "templates/ue5/"+templateName)
	if err != nil {
		return err
	}
	return tmpl.Execute(file, values)
}

func (c *CodegenUE5) fieldType(f *CodeStructField) string {
	if f.IsArray {
		return "TArray<" + c.fieldElemType(f) + ">"
	}
	return c.fieldElemType(f)
}

func (c *CodegenUE5) fieldElemType(f *CodeStructField) string {
	switch f.Type {
	case FieldTypeInt:
		return "int32"
	case FieldTypeLong:
		return "int64"
	case FieldTypeFloat:
		return "double"
	case FieldTypeBool:
		return "bool"
	case FieldTypeString:
		return "FString"
	case FieldTypeTime:
		return "FDateTime"
	case FieldTypeJSON:
		return "TSharedPtr<FJsonValue>"
	case FieldTypeStruct:
		return "F" + c.Prefix + pascal(f.StructRef.Name)
	default:
		panic("unknown type: " + f.Type)
	}
}
