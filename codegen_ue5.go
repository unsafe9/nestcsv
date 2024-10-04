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
	for file := range code.Files() {
		values := map[string]any{
			"File":   file,
			"Prefix": c.Prefix,
		}
		if err := c.template(c.Prefix+pascal(file.Name)+".hpp", "file.hpp.tpl", values); err != nil {
			return err
		}
	}
	return nil
}

func (c *CodegenUE5) template(fileName, templateName string, values any) error {
	file, err := createFile(c.RootDir, fileName, filepath.Ext(fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.
		New(filepath.Base(templateName)).
		Funcs(templateFuncMap).
		Funcs(template.FuncMap{
			"fieldType": c.fieldType,
		}).
		ParseFS(templateFS, "templates/ue5/"+templateName)
	if err != nil {
		return err
	}
	return tmpl.Execute(file, values)
}

func (c *CodegenUE5) fieldType(f *CodeStructField) string {
	ret := ""
	switch f.Type {
	case FieldTypeInt:
		ret = "int32"
	case FieldTypeLong:
		ret = "int64"
	case FieldTypeFloat:
		ret = "double"
	case FieldTypeBool:
		ret = "bool"
	case FieldTypeString:
		ret = "FString"
	case FieldTypeTime:
		ret = "FDateTime"
	case FieldTypeJSON:
		ret = "TSharedPtr<FJsonValue>"
	case FieldTypeStruct:
		ret = c.Prefix + pascal(f.StructRef.Name)
	default:
		panic("unknown type: " + f.Type)
	}
	if f.IsArray {
		ret = "TArray<" + ret + ">"
	}
	return ret
}
