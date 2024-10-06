package nestcsv

import (
	"bytes"
	"fmt"
	"go/format"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var goEmptyImportRegexp = regexp.MustCompile(`import \(\s*\n\s*\)`)

type CodegenGo struct {
	RootDir      string `yaml:"root_dir"`
	PackageName  string `yaml:"package_name"`
	Loader       bool   `yaml:"loader"`
	DataLoadPath string `yaml:"data_load_path"`
}

func (c *CodegenGo) Generate(code *Code) error {
	if c.PackageName == "" {
		c.PackageName = filepath.Base(c.RootDir)
	}

	values := map[string]any{
		"PackageName": c.PackageName,
	}
	if err := c.template("table_base.go", "table_base.go.tpl", values); err != nil {
		return err
	}

	for file := range code.Files() {
		values := map[string]any{
			"PackageName": c.PackageName,
			"File":        file,
		}
		if file.IsTable {
			values["DataLoadPath"] = filepath.Join(c.DataLoadPath, file.Name+".json")
		}
		if err := c.template(file.Name+".go", "file.go.tpl", values); err != nil {
			return err
		}
	}

	if c.Loader {
		values := map[string]any{
			"PackageName": c.PackageName,
			"Tables":      code.Tables,
		}
		err := c.template("loader.go", "loader.go.tpl", values)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CodegenGo) template(fileName, templateName string, values any) error {
	tmpl, err := template.
		New(filepath.Base(templateName)).
		Funcs(templateFuncMap).
		Funcs(template.FuncMap{
			"fieldType":     c.fieldType,
			"fieldElemType": c.fieldElemType,
		}).
		ParseFS(templateFS, "templates/go/"+templateName)
	if err != nil {
		return fmt.Errorf("error parsing template: %s, %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, values); err != nil {
		return fmt.Errorf("error executing template: %s, %w", fileName, err)
	}

	fileBytes := goEmptyImportRegexp.ReplaceAll(buf.Bytes(), nil)

	fileBytes, err = format.Source(fileBytes)
	if err != nil {
		return fmt.Errorf("error formatting source: %s, %w", fileName, err)
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

func (c *CodegenGo) fieldType(f *CodeStructField) string {
	if f.IsArray {
		return "[]" + c.fieldElemType(f)
	}
	return c.fieldElemType(f)
}

func (c *CodegenGo) fieldElemType(f *CodeStructField) string {
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
	case FileTypeJSON:
		return "interface{}"
	case FieldTypeStruct:
		return pascal(f.StructRef.Name)
	default:
		panic("unknown type: " + f.Type)
	}
}
