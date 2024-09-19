{{ range .Tables -}}
-- {{ .Comment }} {{ .CommentInfo }}
CREATE TABLE {{ .GoName }} (
  {{- range .Columns }}
  -- {{ .Comment }} {{ .CommentInfo }}
  {{ .GoName }} {{ .Type }}
    {{- if .PK }} NOT NULL{{ end -}}
    ,
  {{- end }}
) PRIMARY KEY (
  {{- range $i, $pk := .PKColumns -}}
    {{- if $i }}, {{ end }}{{ $pk.GoName }}
  {{- end -}}
)
{{- if .InterleaveTable -}}
  ,
  INTERLEAVE IN PARENT {{ .InterleaveTable }} ON DELETE CASCADE
{{- end -}}
{{- if .DeletionPolicy -}}
  ,
  ROW DELETION POLICY (OLDER_THAN({{ .DeletionPolicy.TimestampColumn }}, INTERVAL {{ .DeletionPolicy.Days }} DAY))
{{- end -}}
;
{{- $name := .GoName }}
{{ range .Indexes -}}
CREATE {{ if .Unique }}UNIQUE {{ end }}INDEX Idx{{ $name }}By{{ range .Keys }}{{ .GoName }}{{ end }} ON {{ $name }}(
  {{- range $i, $col := .Keys }}{{ if $i }}, {{ end }}{{ $col.GoName }}{{ if $col.Desc }} Desc{{ end }}{{ end -}}
)
{{- if len .Storing }} STORING ({{ range $i, $s := .Storing }}{{ if $i }}, {{ end }}{{ $s }}{{ end }}){{ end -}}
;
{{ end }}
{{ end -}}
