{{ template "autogen_comment" }}
package converter

import (
    "context"

    "github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
    "github.com/xhayamix/proto-gen-spanner/pkg/domain/dto/payment"
    "github.com/xhayamix/proto-gen-spanner/pkg/domain/entity/transaction"
    "github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/enums"
    prototransaction "github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/transaction"
    "github.com/xhayamix/proto-gen-spanner/pkg/util/time"
)

type PKConverter interface {
{{- range . }}
    ToProto{{ .PascalName }}{{ if .IsCompositePK }}List{{ end }}(ctx context.Context, {{ if .IsCompositePK }}rows {{ .PkgName }}.{{ .PascalName }}Slice{{ else }}row *{{ .PkgName }}.{{ .PascalName }}{{ end }}) ({{ if .IsCompositePK }}[]{{ end }}*prototransaction.{{ .PascalName }}, error)
{{- end }}
}

type pkConverter struct{}

func NewPKConverter() PKConverter {
    return &pkConverter{}
}
{{ range . }}
func (c *pkConverter) toProto{{ .PascalName }}(_ context.Context, row *{{ .PkgName }}.{{ .PascalName }}) (*prototransaction.{{ .PascalName }}, error) {
    return &prototransaction.{{ .PascalName }}{
        {{- range .Fields }}
        {{ template "row" . }}
        {{- end }}
    }, nil
}
{{ if .IsCompositePK }}
func (c *pkConverter) ToProto{{ .PascalName }}List(ctx context.Context, rows {{ .PkgName }}.{{ .PascalName }}Slice) ([]*prototransaction.{{ .PascalName }}, error) {
    results := make([]*prototransaction.{{ .PascalName }}, 0, len(rows))

    for _, row := range rows {
        result, err := c.toProto{{ .PascalName }}(ctx, row)
        if err != nil {
            return nil, cerrors.Stack(err)
        }
        results = append(results, result)
    }

    return results, nil
}
{{- else }}
func (c *pkConverter) ToProto{{ .PascalName }}(ctx context.Context, row *{{ .PkgName }}.{{ .PascalName }}) (*prototransaction.{{ .PascalName }}, error) {
    if row == nil {
        return nil, nil
    }

    return c.toProto{{ .PascalName }}(ctx, row)
}
{{- end }}
{{ end -}}
