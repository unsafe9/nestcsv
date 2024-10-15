package nestcsv

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type CodegenGo struct {
	RootDir     string `yaml:"root_dir"`
	PackageName string `yaml:"package_name"`
	Singleton   bool   `yaml:"singleton"`
	Context     bool   `yaml:"context"`
}

func (c *CodegenGo) Generate(code *Code) error {
	if c.PackageName == "" {
		c.PackageName = filepath.Base(c.RootDir)
	}

	if err := c.template("table_base.go", "table_base.go.tpl", nil); err != nil {
		return err
	}

	for file := range code.Files() {
		values := map[string]any{
			"File": file,
		}
		if err := c.template(file.Name+".go", "file.go.tpl", values); err != nil {
			return err
		}
	}

	values := map[string]any{
		"Tables": code.Tables,
	}
	return c.template("loader.go", "loader.go.tpl", values)
}

var goEmptyImportRegexp = regexp.MustCompile(`import \(\s*\n\s*\)`)

func (c *CodegenGo) template(fileName, templateName string, values map[string]any) error {
	tmpl, err := template.
		New(filepath.Base(templateName)).
		Funcs(templateFuncMap).
		Funcs(template.FuncMap{
			"fieldType":          c.fieldType,
			"fieldElemType":      c.fieldElemType,
			"fieldPrimitiveType": c.fieldPrimitiveType,
		}).
		ParseFS(templateFS, "templates/go/"+templateName)
	if err != nil {
		return fmt.Errorf("error parsing template: %s, %w", templateName, err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(
		&buf,
		extendMap(
			values,
			map[string]any{
				"PackageName": c.PackageName,
				"Singleton":   c.Singleton,
				"Context":     c.Context,
			},
		),
	)
	if err != nil {
		return fmt.Errorf("error executing template: %s, %w", fileName, err)
	}

	fileBytes := goEmptyImportRegexp.ReplaceAll(buf.Bytes(), nil)

	fileBytes, err = format.Source(fileBytes)
	if err != nil {
		log.Printf("%s", string(fileBytes))
		return fmt.Errorf("error formatting source: %s, %w", fileName, err)
	}

	file, err := createFile(c.RootDir, strings.ToLower(fileName), "go")
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(fileBytes); err != nil {
		return err
	}
	return nil
}

func (c *CodegenGo) fieldType(f *CodeStructField) string {
	if f.IsArray {
		return "[]" + c.fieldElemType(f)
	}
	return c.fieldElemType(f)
}

func (c *CodegenGo) fieldElemType(f *CodeStructField) string {
	if f.Type == FieldTypeStruct {
		return pascal(f.StructRef.Name)
	}
	return c.fieldPrimitiveType(f.Type)
}

func (c *CodegenGo) fieldPrimitiveType(typ FieldType) string {
	switch typ {
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
	default:
		panic("unknown type: " + typ)
	}
}
