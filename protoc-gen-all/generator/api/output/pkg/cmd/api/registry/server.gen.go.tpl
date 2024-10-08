{{ template "autogen_comment" }}
package registry

import (
	"google.golang.org/grpc"

	pb "github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/api"
)

func registerServer(s *grpc.Server) func(
{{- range . }}
	{{ .CamelName }}Server pb.{{ .PascalName }}Server,
{{- end }}
) {
	return func(
{{- range . }}
	{{ .CamelName }}Server pb.{{ .PascalName }}Server,
{{- end }}
	) {
{{- range . }}
	pb.Register{{ .PascalName }}Server(s, {{ .CamelName }}Server)
{{- end }}
	}
}
