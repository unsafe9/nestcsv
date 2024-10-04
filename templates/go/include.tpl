{{- define "goFieldType" -}}
{{- if eq .Type "int" -}}
int
{{- else if eq .Type "long" -}}
int64
{{- else if eq .Type "float" -}}
float64
{{- else if eq .Type "bool" -}}
bool
{{- else if eq .Type "string" -}}
string
{{- else if eq .Type "time" -}}
time.Time
{{- else if eq .Type "json" -}}
any
{{- else if eq .Type "struct" -}}
{{ pascal .StructRef.Name }}
{{- end -}}
{{- end -}}

{{- define "goField" -}}
{{- if .IsArray -}}[]{{- end -}}{{ template "goFieldType" . }}
{{- end -}}

{{- define "goStruct" -}}
type {{ pascal .Name }} struct {
{{- range .Fields }}
    {{ pascal .Name }} {{ template "goField" . }} `json:"{{ .Name }}"`
{{- end }}
}
{{- end -}}
