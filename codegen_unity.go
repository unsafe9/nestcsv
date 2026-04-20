package nestcsv

import (
	"path/filepath"
	"text/template"
)

type CodegenUnity struct {
	RootDir        string `yaml:"root_dir"`
	Prefix         string `yaml:"prefix"`
	Namespace      string `yaml:"namespace"`
	Singleton      bool   `yaml:"singleton"`
	DataSuffix     string `yaml:"data_suffix"`
	TableSuffix    string `yaml:"table_suffix"`
	ResourceFolder string `yaml:"resource_folder"`
	FileSuffix     string `yaml:"file_suffix"`
}

func (c *CodegenUnity) Generate(code *Code) error {
	if c.TableSuffix == "" {
		c.TableSuffix = "Table"
	}
	if c.FileSuffix == "" {
		c.FileSuffix = ".cs"
	}

	baseValues := map[string]any{}
	if err := c.template(c.Prefix+"TableDataBase", "TableDataBase.cs.tpl", baseValues); err != nil {
		return err
	}
	if err := c.template(c.Prefix+"TableBase", "TableBase.cs.tpl", baseValues); err != nil {
		return err
	}

	holderValues := map[string]any{
		"Tables": code.Tables,
	}
	if err := c.template(c.Prefix+"TableHolder", "TableHolder.cs.tpl", holderValues); err != nil {
		return err
	}

	for file := range code.Files {
		fileValues := map[string]any{
			"File": file,
		}
		className := c.Prefix + pascal(file.Name) + c.DataSuffix
		if err := c.template(className, "file.cs.tpl", fileValues); err != nil {
			return err
		}
	}
	return nil
}

func (c *CodegenUnity) template(fileName, templateName string, values map[string]any) error {
	file, err := createFile(c.RootDir, fileName, c.FileSuffix)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.
		New(filepath.Base(templateName)).
		Funcs(templateFuncMap).
		Funcs(template.FuncMap{
			"fieldType":          c.fieldType,
			"fieldElemType":      c.fieldElemType,
			"fieldPrimitiveType": c.fieldPrimitiveType,
		}).
		ParseFS(templateFS, "templates/unity/"+templateName)
	if err != nil {
		return err
	}

	return tmpl.Execute(
		file,
		extendMap(
			values,
			map[string]any{
				"Prefix":         c.Prefix,
				"Namespace":      c.Namespace,
				"Singleton":      c.Singleton,
				"DataSuffix":     c.DataSuffix,
				"TableSuffix":    c.TableSuffix,
				"ResourceFolder": c.ResourceFolder,
			},
		),
	)
}

func (c *CodegenUnity) fieldType(f *CodeStructField) string {
	if f.IsArray {
		return "List<" + c.fieldElemType(f) + ">"
	}
	return c.fieldElemType(f)
}

func (c *CodegenUnity) fieldElemType(f *CodeStructField) string {
	if f.Type == FieldTypeStruct {
		return c.Prefix + pascal(f.StructRef.Name) + c.DataSuffix
	}
	return c.fieldPrimitiveType(f.Type)
}

func (c *CodegenUnity) fieldPrimitiveType(typ FieldType) string {
	switch typ {
	case FieldTypeInt:
		return "int"
	case FieldTypeLong:
		return "long"
	case FieldTypeFloat:
		return "double"
	case FieldTypeBool:
		return "bool"
	case FieldTypeString:
		return "string"
	case FieldTypeTime:
		return "DateTime"
	case FieldTypeJSON:
		return "JToken"
	default:
		panic("unknown type: " + typ)
	}
}
