package spanner

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

//go:embed entity.gen.go.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(message *input.Message) (*output.TemplateInfo, error) {
	if !input.ServerMessageAccessorSet.Contains(message.Option.AccessorType) {
		return nil, nil
	}

	type Column struct {
		GoName       string
		CamelName    string
		Type         string
		DatabaseType string
		Comment      string
		PK           bool
		IsList       bool
		IsEnum       bool
		IsOnlyServer bool
	}
	type Index struct {
		GoName  string
		Comment string
		Columns []*Column
	}
	type Table struct {
		PkgName      string
		GoName       string
		SnakeName    string
		Comment      string
		Columns      []*Column
		PKColumns    []*Column
		Indexes      []*Index
		HasUserID    bool
		InsertTiming string
	}

	data := &Table{
		PkgName:      core.ToPkgName(message.SnakeName),
		GoName:       core.ToGolangPascalCase(message.SnakeName),
		SnakeName:    message.SnakeName,
		Comment:      message.Comment,
		Columns:      make([]*Column, 0, len(message.Fields)),
		PKColumns:    make([]*Column, 0),
		Indexes:      make([]*Index, 0, len(message.Option.DDL.Indexes)),
		HasUserID:    false,
		InsertTiming: message.Option.InsertTiming,
	}

	for _, field := range message.Fields {
		if !input.ServerFieldAccessorSet.Contains(field.Option.AccessorType) {
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

		databaseTypeName, err := input.GetSpannerType(field.TypeKind, field.SnakeName, field.IsList)
		if err != nil {
			return nil, perrors.Stack(err)
		}

		column := &Column{
			GoName:       core.ToGolangPascalCase(field.SnakeName),
			CamelName:    core.ToCamelCase(field.SnakeName),
			Type:         typeName,
			DatabaseType: databaseTypeName,
			Comment:      field.Comment,
			PK:           field.Option.DDL.PK,
			IsList:       field.IsList,
			IsEnum:       field.TypeKind == input.TypeKind_Enum,
			IsOnlyServer: field.Option.AccessorType == input.FieldAccessorType_OnlyServer || field.Option.AccessorType == input.FieldAccessorType_ServerAndClient,
		}

		if field.Option.DDL.PK {
			if field.SnakeName == "user_id" {
				data.HasUserID = true
			}

			data.PKColumns = append(data.PKColumns, column)
		}
		data.Columns = append(data.Columns, column)
	}

	for _, index := range message.Option.DDL.Indexes {
		i := &Index{
			GoName:  data.GoName + "Idx",
			Comment: data.GoName + " Index(",
			Columns: make([]*Column, 0, len(data.PKColumns)+len(index.Keys)+len(index.SnakeStoring)),
		}
		keySet := strset.NewWithSize(len(index.Keys) + len(index.SnakeStoring))
		goNames := make([]string, 0, len(index.Keys))
		for _, key := range index.Keys {
			gn := core.ToGolangPascalCase(key.SnakeName)
			goNames = append(goNames, gn)
			keySet.Add(gn)
		}
		i.GoName += strings.Join(goNames, "")
		i.Comment += strings.Join(goNames, ", ") + ")"
		for _, key := range index.SnakeStoring {
			keySet.Add(core.ToGolangPascalCase(key))
		}

		// columns定義順にindex内のcolumnを整えるためcolumnsでforを回す
		for _, column := range data.Columns {
			if column.PK || keySet.Has(column.GoName) {
				i.Columns = append(i.Columns, column)
			}
		}
		data.Indexes = append(data.Indexes, i)
	}

	tpl, err := core.GetBaseTemplate().Funcs(sprig.TxtFuncMap()).Parse(string(templateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("pkg/domain/entity/transaction", data.SnakeName+".gen.go"),
	}, nil
}
