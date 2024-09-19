package converter

import (
	"bytes"
	_ "embed"
	"sort"

	"github.com/scylladb/go-set/strset"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/input"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

//go:embed row.tpl
var rowTemplateFileBytes []byte

type Table struct {
	PascalName    string
	PkgName       string
	IsCompositePK bool
	Fields        []*Field
}

type Field struct {
	PascalName  string
	GoName      string
	Type        string
	CastFunc    string
	CastWithPtr bool
}

func convertTable(message *input.Message, onlyPK bool) *Table {
	if !input.ClientMessageCommonResponseAccessorSet.Contains(message.Option.AccessorType) {
		return nil
	}

	var pkCount int32
	fields := make([]*Field, 0, len(message.Fields))
	for _, field := range message.Fields {
		if field.Option.DDL.PK {
			pkCount++
		} else if onlyPK {
			continue
		}

		if !input.ClientFieldAccessorSet.Contains(field.Option.AccessorType) {
			continue
		}

		// クライアントだけにしか実装されていないフィールドの場合は自動生成しない（そのままだとそのフィールドにはデフォルト値が入る）
		if field.Option.AccessorType == input.FieldAccessorType_OnlyClient {
			continue
		}

		var castFunc string
		var castWithPtr bool
		switch field.RawTypeKind {
		case input.TypeKind_Int32:
			if field.IsList {
				castFunc = "toInt32Slice"
			} else {
				castFunc = "int32"
			}
		case input.TypeKind_Enum:
			if field.IsList {
				castFunc = "toProto" + field.Type + "Slice"
			} else {
				castFunc = "enums." + field.Type
			}
		}
		if core.IsTimeField(field.SnakeName) {
			castFunc = "time.ToUnixMilli"
			castWithPtr = true
		}

		fields = append(fields, &Field{
			PascalName:  core.ToPascalCase(field.SnakeName),
			GoName:      core.ToGolangPascalCase(field.SnakeName),
			Type:        field.Type,
			CastFunc:    castFunc,
			CastWithPtr: castWithPtr,
		})
	}

	var pkgName string
	if message.SnakeName == "user_balance" {
		pkgName = "payment"
	} else {
		pkgName = "transaction"
	}

	return &Table{
		PascalName:    core.ToPascalCase(message.SnakeName),
		PkgName:       pkgName,
		IsCompositePK: pkCount >= 2,
		Fields:        fields,
	}
}

//go:embed entity_converter.gen.go.tpl
var entityConverterTemplateFileBytes []byte

type EntityConverterCreator struct{}

func (c *EntityConverterCreator) Create(messages []*input.Message, _ map[string]*input.Enum) (*output.TemplateInfo, error) {
	data := make([]*Table, 0, len(messages))

	for _, message := range messages {
		table := convertTable(message, false)
		if table == nil {
			continue
		}

		data = append(data, table)
	}

	rowTpl, err := core.GetBaseTemplate().Parse(string(rowTemplateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	tpl, err := rowTpl.Parse(string(entityConverterTemplateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/domain/proto/converter", "entity_converter.gen.go"),
	}, nil
}

//go:embed pk_converter.gen.go.tpl
var pkConverterTemplateFileBytes []byte

type PKConverterCreator struct{}

func (c *PKConverterCreator) Create(messages []*input.Message, _ map[string]*input.Enum) (*output.TemplateInfo, error) {
	data := make([]*Table, 0, len(messages))

	for _, message := range messages {
		table := convertTable(message, true)
		if table == nil {
			continue
		}

		data = append(data, table)
	}

	rowTpl, err := core.GetBaseTemplate().Parse(string(rowTemplateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	tpl, err := rowTpl.Parse(string(pkConverterTemplateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/domain/proto/converter", "pk_converter.gen.go"),
	}, nil
}

//go:embed util.gen.go.tpl
var utilTemplateFileBytes []byte

type UtilCreator struct{}

func (c *UtilCreator) Create(messages []*input.Message, _ map[string]*input.Enum) (*output.TemplateInfo, error) {
	type Data struct {
		Enums []string
	}

	set := strset.New()
	for _, message := range messages {
		if !input.ClientMessageCommonResponseAccessorSet.Contains(message.Option.AccessorType) {
			continue
		}

		for _, field := range message.Fields {
			if !input.ClientFieldAccessorSet.Contains(field.Option.AccessorType) {
				continue
			}

			if field.TypeKind == input.TypeKind_Enum {
				set.Add(field.Type)
			}
		}
	}

	l := set.List()
	sort.Strings(l)
	data := &Data{
		Enums: l,
	}

	tpl, err := core.GetBaseTemplate().Parse(string(utilTemplateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/domain/proto/converter", "util.gen.go"),
	}, nil
}
