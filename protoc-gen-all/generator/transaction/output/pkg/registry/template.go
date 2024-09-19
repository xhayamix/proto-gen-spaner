package registry

import (
	"bytes"
	_ "embed"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/input"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

//go:embed transaction_repository.gen.go.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(messages []*input.Message, _ map[string]*input.Enum) (*output.TemplateInfo, error) {
	type Table struct {
		GoName string
	}
	data := make([]*Table, 0, len(messages))

	for _, message := range messages {
		if !input.ServerMessageAccessorSet.Contains(message.Option.AccessorType) {
			continue
		}

		table := &Table{
			GoName: core.ToGolangPascalCase(message.SnakeName),
		}

		data = append(data, table)
	}

	tpl, err := core.GetBaseTemplate().Parse(string(templateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/registry", "transaction_repository.gen.go"),
	}, nil
}
