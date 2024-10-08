{{ template "autogen_comment" }}

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_api.go
//go:generate goimports -w --local "github.com/xhayamix/proto-gen-spanner" mock_$GOPACKAGE/mock_api.go
package api

import (
	"context"

	"github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/api"
)

type Method string

var (
{{- range $service := . }}
    // {{ $service.PascalName }}
{{ range $method := $service.Methods -}}
	{{ $service.PascalName }}{{ $method.PascalName }} Method = "{{ $service.PascalName }}{{ $method.PascalName }}"
{{ end -}}
{{ end -}}
)

type API interface {
{{ range $service := . -}}
{{ range $method := $service.Methods -}}
	// {{ $service.PascalName }}{{ $method.PascalName }} {{ $method.Description }}
	{{ $service.PascalName }}{{ $method.PascalName }}(ctx context.Context {{ if not $method.IsRequestEmpty }}, req *api.{{ $method.PascalName }}Request {{ end }} ) (*api.{{ $method.PascalName }}Response, error)
{{ end -}}
{{ end -}}
	// Close クローズ処理
	Close() error
}
