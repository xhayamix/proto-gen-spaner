{{- define "row" }}
{{- .PascalName }}: {{ if .CastFunc }}{{ .CastFunc }}({{ if .CastWithPtr }}&{{ end }}row.{{ .GoName }}){{ else }}row.{{ .GoName }}{{ end }},
{{- end -}}
