// Code generated by "nestcsv"; DO NOT EDIT.

package {{ .PackageName }}

type Tables struct{
{{- range sortBy .Tables "Name" }}
    {{ pascal .Struct.Name }} {{ pascal .Struct.Name }}Table
{{- end }}
}

func LoadTables() (*Tables, error) {
    var t Tables
{{- range sortBy .Tables "Name" }}
    if err := t.{{ pascal .Struct.Name }}.Load(); err != nil {
        return nil, err
    }
{{- end }}
    return &t, nil
}
