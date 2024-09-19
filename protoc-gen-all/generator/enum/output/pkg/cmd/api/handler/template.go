package interactor

import (
	"bytes"
	_ "embed"

	"github.com/Masterminds/sprig/v3"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/enum/input"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/enum/output"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

type Element struct {
	PascalName string
	Type       string
	IsEnum     bool
}
type Enum struct {
	Elements []*Element
}

//go:embed handler.gen.go.tpl
var templateFileBytes []byte

type EachCreator struct{}

func (c *EachCreator) Create(enum *input.Enum) (*output.TemplateInfo, error) {
	if enum.SnakeName != "preference_type" {
		return nil, nil
	}

	e := &Enum{
		Elements: make([]*Element, 0, len(enum.Elements)),
	}
	for _, el := range enum.Elements {
		var elType string
		var isEnum bool
		switch el.SettingType {
		case input.SettingType_Bool:
			elType = "bool"
		case input.SettingType_Int32:
			elType = "int32"
		case input.SettingType_Int64:
			elType = "int64"
		case input.SettingType_String:
			elType = "string"
		case input.SettingType_Int32List:
			elType = "[]int32"
		case input.SettingType_Int64List:
			elType = "[]int64"
		case input.SettingType_StringList:
			elType = "[]string"
		}

		e.Elements = append(e.Elements, &Element{
			PascalName: core.ToPascalCase(el.RawName),
			Type:       elType,
			IsEnum:     isEnum,
		})
	}

	tpl, err := core.GetBaseTemplate().Funcs(sprig.TxtFuncMap()).Parse(string(templateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, e); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/cmd/api/handler/preference", "handler.gen.go"),
	}, nil
}
