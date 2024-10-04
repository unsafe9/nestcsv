package nestcsv

import (
	"bytes"
	"fmt"
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

func (c *CodegenGo) Generate(code *Code) error {
	if c.PackageName == "" {
		c.PackageName = filepath.Base(c.RootDir)
	}

	for _, file := range code.NamedStructs {
		values := map[string]any{
			"PackageName": c.PackageName,
			"File":        file,
		}
		if err := c.template(file.Name+".go", "named_struct.go.tpl", values); err != nil {
			return err
		}
	}

	for _, file := range code.Tables {
		values := map[string]any{
			"PackageName":  c.PackageName,
			"File":         file,
			"DataLoadPath": filepath.Join(c.DataLoadPath, file.Name+".json"),
		}
		if err := c.template(file.Name+".go", "table.go.tpl", values); err != nil {
			return err
		}
	}

	values := map[string]any{
		"PackageName": c.PackageName,
		"Tables":      code.Tables,
	}
	err := c.template("loader.go", "loader.go.tpl", values)
	if err != nil {
		return err
	}

	return nil
}

func (c *CodegenGo) template(fileName, templateName string, values any) error {
	tmpl, err := template.
		New(filepath.Base(templateName)).
		Funcs(templateFuncMap).
		ParseFS(templateFS, "templates/go/"+templateName, "templates/go/include.tpl")
	if err != nil {
		return fmt.Errorf("error parsing template: %s, %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, values); err != nil {
		return fmt.Errorf("error executing template: %s, %w", fileName, err)
	}

	fileBytes, err := format.Source(buf.Bytes())
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
