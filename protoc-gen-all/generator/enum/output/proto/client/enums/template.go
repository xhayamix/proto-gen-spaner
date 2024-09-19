package enums

import (
	"bytes"
	_ "embed"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/enum/input"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/enum/output"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

type Element struct {
	PascalName string
	Value      int32
	Comment    string
}
type Enum struct {
	PascalName string
	Comment    string
	Elements   []*Element
}

//go:embed enum.gen.proto.tpl
var enumTemplateFileBytes []byte

type EachCreator struct{}

func (c *EachCreator) Create(enum *input.Enum) (*output.TemplateInfo, error) {
	if enum.AccessorType != input.AccessorType_ServerAndClient {
		return nil, nil
	}

	e := &Enum{
		PascalName: core.ToPascalCase(enum.SnakeName),
		Comment:    "",
		Elements:   make([]*Element, 0, len(enum.Elements)),
	}
	for _, element := range enum.Elements {
		e.Elements = append(e.Elements, &Element{
			PascalName: element.RawName,
			Value:      element.Value,
			Comment:    element.Comment,
		})
	}

	tpl, err := core.GetBaseTemplate().Parse(string(enumTemplateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, e); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("proto/server/enums", enum.SnakeName+"_gen.proto"),
	}, nil
}
