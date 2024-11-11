package nestcsv

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type CodegenUE5 struct {
	RootDir string `yaml:"root_dir"`
	Prefix  string `yaml:"prefix"`
}

func (c *CodegenUE5) Generate(code *Code) error {
	if err := c.template("TableDataBase.h", "TableDataBase.h.tpl", false, nil); err != nil {
		return err
	}
	if err := c.template("TableBase.h", "TableBase.h.tpl", false, nil); err != nil {
		return err
	}

	values := map[string]any{
		"Tables": code.Tables,
	}
	if err := c.template("TableHolder.h", "TableHolder.h.tpl", false, values); err != nil {
		return err
	}

	for file := range code.Files {
		values = map[string]any{
			"File": file,
		}
		if err := c.template(pascal(file.Name)+".h", "file.h.tpl", true, values); err != nil {
			return err
		}

		if file.IsTable {
			if err := c.template(pascal(file.Name)+"Table.h", "table.h.tpl", true, values); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *CodegenUE5) readExistingFileContent(fileName string) (map[string]any, error) {
	filePath := makeFilePath(c.RootDir, c.Prefix+fileName, filepath.Ext(fileName))
	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var (
		lineBreak         = []byte("\n")
		regionStartRegexp = regexp.MustCompile(`\s*//\s*NESTCSV:(\w+)_START`)
		regionEndRegexp   *regexp.Regexp
		regions           = make(map[string]any)
		region            string
		builder           strings.Builder
		lines             = bytes.Split(file, lineBreak)
	)
	for _, line := range lines {
		if matches := regionStartRegexp.FindSubmatch(line); len(matches) > 1 {
			region = string(matches[1])
			regionEndRegexp = regexp.MustCompile(`\s*//\s*NESTCSV:` + region + `_END`)
			builder.Reset()

		} else if region != "" {
			if regionEndRegexp.Match(line) {
				regions[region] = strings.TrimSpace(builder.String())
				region = ""

			} else {
				builder.Write(line)
				builder.Write(lineBreak)
			}
		}
	}
	return regions, nil
}

func (c *CodegenUE5) template(fileName, templateName string, withRegions bool, values map[string]any) error {
	if withRegions {
		if existingContent, err := c.readExistingFileContent(fileName); err != nil {
			return err
		} else {
			values = extendMap(values, map[string]any{
				"ExistingContent": existingContent,
			})
		}
	}

	file, err := createFile(c.RootDir, c.Prefix+fileName, filepath.Ext(fileName))
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
		ParseFS(templateFS, "templates/ue5/"+templateName)
	if err != nil {
		return err
	}

	return tmpl.Execute(
		file,
		extendMap(
			values,
			map[string]any{
				"Prefix": c.Prefix,
			},
		),
	)
}

func (c *CodegenUE5) fieldType(f *CodeStructField) string {
	if f.IsArray {
		return "TArray<" + c.fieldElemType(f) + ">"
	}
	return c.fieldElemType(f)
}

func (c *CodegenUE5) fieldElemType(f *CodeStructField) string {
	if f.Type == FieldTypeStruct {
		return "F" + c.Prefix + pascal(f.StructRef.Name)
	}
	return c.fieldPrimitiveType(f.Type)
}

func (c *CodegenUE5) fieldPrimitiveType(typ FieldType) string {
	switch typ {
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
	default:
		panic("unknown type: " + typ)
	}
}
