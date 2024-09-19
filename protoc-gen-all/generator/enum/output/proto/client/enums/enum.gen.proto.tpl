{{ template "autogen_comment" }}
syntax = "proto3";

package client.enums;

option go_package = "github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/enums";

{{ $Name := .PascalName -}}
enum {{ .PascalName }} {
{{ .PascalName }}_Unknown = 0;
{{- range .Elements }}
    {{ if .Comment }}// {{ .Comment }}{{ end }}
    {{ $Name }}_{{ .PascalName }} = {{ .Value }};
{{- end }}
}
