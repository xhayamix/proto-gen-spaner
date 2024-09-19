package base

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/scylladb/go-set/strset"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/input"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

//go:embed query_builder.gen.go.tpl
var queryBuilderTemplateFileBytes []byte

type QueryBuilderCreator struct{}

func (c *QueryBuilderCreator) Create(message *input.Message) (*output.TemplateInfo, error) {
	if !input.ServerMessageAccessorSet.Contains(message.Option.AccessorType) {
		return nil, nil
	}

	type Column struct {
		GoName  string
		Type    string
		SetType string
		PK      bool
		IsList  bool
		IsEnum  bool
	}
	type Type struct {
		GoName  string
		Key     string
		Columns []*Column
	}
	type Table struct {
		PkgName   string
		GoName    string
		CamelName string
		Columns   []*Column
		Types     []*Type
	}
	data := &Table{
		PkgName:   core.ToPkgName(message.SnakeName),
		GoName:    core.ToGolangPascalCase(message.SnakeName),
		CamelName: core.ToCamelCase(message.SnakeName),
		Columns:   make([]*Column, 0, len(message.Fields)),
		Types:     make([]*Type, 0),
	}

	for _, field := range message.Fields {
		if !input.AdminFieldAccessorSet.Contains(field.Option.AccessorType) {
			continue
		}

		typeName := field.Type
		if core.IsTimeField(field.SnakeName) {
			typeName = "time.Time"
		}
		if field.TypeKind == input.TypeKind_Enum && field.IsList {
			typeName = "enum." + typeName + "Slice"
		} else {
			if field.TypeKind == input.TypeKind_Enum {
				typeName = "enum." + typeName
			}
			if field.IsList {
				typeName = "[]" + typeName
			}
		}

		var setType string
		switch field.TypeKind {
		case input.TypeKind_Int32, input.TypeKind_Enum:
			setType = "i32set"
		case input.TypeKind_Int64:
			setType = "i64set"
		case input.TypeKind_String:
			setType = "strset"
		case input.TypeKind_Bool, input.TypeKind_Bytes:
		default:
			return nil, perrors.Newf("サポートされていないTypeKindです。 TypeKind = %v", field.TypeKind)
		}

		column := &Column{
			GoName:  core.ToGolangPascalCase(field.SnakeName),
			Type:    typeName,
			SetType: setType,
			PK:      field.Option.DDL.PK,
			IsList:  field.IsList,
			IsEnum:  field.TypeKind == input.TypeKind_Enum,
		}
		data.Columns = append(data.Columns, column)
	}
	data.Types = append(data.Types, &Type{
		GoName:  "All",
		Key:     "All",
		Columns: data.Columns,
	})

	for _, index := range message.Option.DDL.Indexes {
		typ := &Type{}
		keySet := strset.NewWithSize(len(index.Keys) + len(index.SnakeStoring))
		goNames := make([]string, 0, len(index.Keys))
		for _, key := range index.Keys {
			gn := core.ToGolangPascalCase(key.SnakeName)
			goNames = append(goNames, gn)
			keySet.Add(gn)
		}
		typ.Key = "Idx" + data.GoName + "By" + strings.Join(goNames, "")
		typ.GoName = "Idx" + strings.Join(goNames, "")
		for _, key := range index.SnakeStoring {
			keySet.Add(core.ToGolangPascalCase(key))
		}

		// columns定義順にindex内のcolumnを整えるためcolumnsでforを回す
		for _, column := range data.Columns {
			if column.PK || keySet.Has(column.GoName) {
				typ.Columns = append(typ.Columns, column)
			}
		}
		data.Types = append(data.Types, typ)
	}

	tpl, err := core.GetBaseTemplate().Funcs(sprig.TxtFuncMap()).Parse(string(queryBuilderTemplateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/infra/transaction/repository/base", message.SnakeName+"-query_builder.gen.go"),
	}, nil
}
