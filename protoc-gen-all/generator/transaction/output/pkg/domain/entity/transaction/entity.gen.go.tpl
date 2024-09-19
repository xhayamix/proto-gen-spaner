{{ template "autogen_comment" }}
package transaction

var AllTableNames = []string{
{{- range . }}
	{{ . }}TableName,
{{- end }}
}