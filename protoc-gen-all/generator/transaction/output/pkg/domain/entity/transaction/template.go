package transaction

import (
	"bytes"
	_ "embed"

	"github.com/Masterminds/sprig/v3"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/input"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

//go:embed entity.gen.go.tpl
var entityTemplateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(messages []*input.Message, _ map[string]*input.Enum) (*output.TemplateInfo, error) {
	data := make([]string, 0, len(messages))

	for _, message := range messages {
		if !input.ServerMessageAccessorSet.Contains(message.Option.AccessorType) {
			continue
		}

		data = append(data, core.ToPascalCase(message.SnakeName))
	}
	tpl, err := core.GetBaseTemplate().Funcs(sprig.TxtFuncMap()).Parse(string(entityTemplateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/domain/entity/transaction", "entity.gen.go"),
	}, nil
}
