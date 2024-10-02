package nestcsv

import (
	"embed"
	"github.com/Masterminds/sprig/v3"
	"io"
	"path/filepath"
	"text/template"
)

var templateFuncMap template.FuncMap

func init() {
	templateFuncMap = sprig.TxtFuncMap()

	templateFuncMap["sortBy"] = sortBy
	templateFuncMap["anyBy"] = anyBy
	templateFuncMap["pascal"] = pascal
}

//go:embed templates/*
var templateFS embed.FS

func templateFile(w io.Writer, templateName string, funcs template.FuncMap, data any) error {
	tmpl, err := template.
		New(filepath.Base(templateName)).
		Funcs(templateFuncMap).
		Funcs(funcs).
		ParseFS(templateFS, "templates/"+templateName)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
