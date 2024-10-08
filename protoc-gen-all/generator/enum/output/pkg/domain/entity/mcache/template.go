package mcache

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/enum/input"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/enum/output"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

//go:embed setter.gen.go.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(enum *input.Enum) (*output.TemplateInfo, error) {
	if !strings.HasSuffix(enum.SnakeName, "setting_type") {
		return nil, nil
	}

	type Element struct {
		PascalName        string
		PascalSettingType string
		HasServer         bool
	}
	type Enum struct {
		SnakeName  string
		PascalName string
		LowerName  string
		Elements   []*Element
	}

	name := strings.ReplaceAll(enum.SnakeName, "_type", "")
	data := &Enum{
		SnakeName:  name,
		PascalName: core.ToPascalCase(name),
		LowerName:  strings.ReplaceAll(name, "_", ""),
		Elements:   make([]*Element, 0, len(enum.Elements)),
	}
	for _, element := range enum.Elements {
		hasServer := false
		if element.SettingAccessorType != input.SettingAccessorType_OnlyClient {
			hasServer = true
		}

		isTime := core.IsTimeField(element.SnakeName)

		var pascalSettingType string
		switch element.SettingType {
		case input.SettingType_Bool:
			pascalSettingType = "Bool"
		case input.SettingType_Int32:
			pascalSettingType = "Int32"
		case input.SettingType_Int64:
			pascalSettingType = "Int64"
			if isTime {
				pascalSettingType = "Time"
			}
		case input.SettingType_String:
			pascalSettingType = "String"
		case input.SettingType_Int32List:
			pascalSettingType = "Int32Slice"
		case input.SettingType_Int64List:
			pascalSettingType = "Int64Slice"
		case input.SettingType_StringList:
			pascalSettingType = "StringSlice"
		default:
			return nil, perrors.Newf("サポートされていないSettingTypeです。 SettingType = %v", element.SettingType)
		}

		data.Elements = append(data.Elements, &Element{
			PascalName:        element.RawName,
			PascalSettingType: pascalSettingType,
			HasServer:         hasServer,
		})
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
		FilePath: core.JoinPath("pkg/domain/entity/mcache", data.SnakeName+".setter.gen.go"),
	}, nil
}
