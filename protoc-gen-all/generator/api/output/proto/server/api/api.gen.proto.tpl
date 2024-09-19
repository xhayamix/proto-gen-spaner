{{ template "autogen_comment" }}
syntax = "proto3";

package client.api{{ if .IsCommon }}.common{{ end }};
option go_package = "github.com/xhayamix/proto-gen-golang/pkg/domain/proto/client/api{{ if .IsCommon }}/common{{ end }}";

{{ range .ImportPaths -}}
import "{{ . }}";
{{ end }}
{{- if .Service }}
// {{ .Service.Comment }}
{{- $serviceName := .Service.PascalName }}
service {{ $serviceName }} {
{{- range .Service.Methods }}
  // {{ .Comment }}
  rpc {{ .PascalName }}({{ .InputType }}) returns ({{ .OutputType }}) {
    // grpc-gateway使う場合は追加
    // option (google.api.http) = {
    //   post: "/{{ $serviceName }}/{{ .PascalName }}"
    //   body: "*"
    // };
  }
{{- end }}
}
{{ end }}
{{ range .Messages -}}
  {{- template "message" . }}
{{ end -}}
