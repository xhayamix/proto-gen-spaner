package enum

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"text/template"

	"github.com/huandu/xstrings"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/enum/input"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/enum/output"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

//go:embed init.gen.schema.tpl
var templateFileBytes []byte

type Creator struct{}

type EnumElementSchema struct {
	RawName  string `json:"rawName,omitempty"`
	Value    int32  `json:"value,omitempty"`
	Comment  string `json:"comment,omitempty"`
	EnumName string `json:"enumName,omitempty"`
}

type EnumSchema struct {
	Name     string               `json:"name,omitempty"`
	Comment  string               `json:"comment,omitempty"`
	Elements []*EnumElementSchema `json:"elements,omitempty"`
}

type Data struct {
	SchemaJSONString string
}

func (c *Creator) Create(enums []*input.Enum) (*output.TemplateInfo, error) {
	schemas := make([]*EnumSchema, 0, len(enums))
	for _, enum := range enums {
		elements := make([]*EnumElementSchema, 0, len(enum.Elements))
		for _, element := range enum.Elements {
			elements = append(elements, &EnumElementSchema{
				RawName: element.RawName,
				Value:   element.Value,
				Comment: element.Comment,
			})
		}
		schemas = append(schemas, &EnumSchema{
			Name:     xstrings.ToCamelCase(enum.SnakeName),
			Comment:  enum.Comment,
			Elements: elements,
		})
	}

	b, err := json.MarshalIndent(schemas, "", " ")
	if err != nil {
		return nil, perrors.Wrapf(err, "schemaをマーシャルできませんでした")
	}

	tpl := template.Must(core.GetBaseTemplate().Parse(string(templateFileBytes)))

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, &Data{
		SchemaJSONString: string(b),
	}); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: "db/ddl/enum/init.gen.schema",
	}, nil
}
