package spanner

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/Masterminds/sprig/v3"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/input"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

//go:embed repository.gen.go.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(message *input.Message) (*output.TemplateInfo, error) {
	if !input.ServerMessageAccessorSet.Contains(message.Option.AccessorType) {
		return nil, nil
	}

	type Column struct {
		GoName string
		Type   string
	}
	type Method struct {
		Name       string
		Args       string
		ReturnName string
	}
	type Table struct {
		SnakeName string
		GoName    string
		Methods   []*Method
	}

	data := &Table{
		SnakeName: message.SnakeName,
		GoName:    core.ToGolangPascalCase(message.SnakeName),
		Methods:   make([]*Method, 0),
	}

	pkColumns := make([]*Column, 0)
	columnMap := make(map[string]*Column, len(message.Fields))
	for _, field := range message.Fields {
		if !input.ServerFieldAccessorSet.Contains(field.Option.AccessorType) {
			continue
		}

		typeName := field.Type

		if core.IsTimeField(field.SnakeName) {
			typeName = "time.Time"
		}
		if field.TypeKind == input.TypeKind_Enum {
			typeName = "enum." + typeName
		}

		column := &Column{
			GoName: core.ToGolangPascalCase(field.SnakeName),
			Type:   typeName,
		}
		if field.Option.DDL.PK {
			pkColumns = append(pkColumns, column)
		}
		columnMap[field.SnakeName] = column
	}

	var methodName string
	var methodArgs string
	last := len(pkColumns) - 1
	for i, column := range pkColumns {
		if i == last {
			break
		}
		methodName += column.GoName
		methodArgs += column.GoName
		data.Methods = append(data.Methods,
			&Method{
				Name:       methodName,
				Args:       methodArgs + " " + column.Type,
				ReturnName: data.GoName,
			},
			&Method{
				Name:       methodName + "s",
				Args:       methodArgs + "s []" + column.Type,
				ReturnName: data.GoName,
			},
		)
		methodName += "And"
		methodArgs += " " + column.Type + ", "
	}

	for _, index := range message.Option.DDL.Indexes {
		goNames := make([]string, 0, len(index.Keys))
		for _, key := range index.Keys {
			if column, ok := columnMap[key.SnakeName]; ok {
				goNames = append(goNames, column.GoName)
			}
		}
		goName := "Idx" + strings.Join(goNames, "")
		methodName = goName + "With"
		methodArgs = ""
		for _, key := range index.Keys {
			column, ok := columnMap[key.SnakeName]
			if !ok {
				continue
			}
			methodName += column.GoName
			methodArgs += column.GoName
			data.Methods = append(data.Methods,
				&Method{
					Name:       methodName,
					Args:       methodArgs + " " + column.Type,
					ReturnName: data.GoName + goName,
				},
				&Method{
					Name:       methodName + "s",
					Args:       methodArgs + "s []" + column.Type,
					ReturnName: data.GoName + goName,
				},
			)
			methodName += "And"
			methodArgs += " " + column.Type + ", "
		}
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
		FilePath: core.JoinPath("pkg/domain/repository/transaction", data.SnakeName+".gen.go"),
	}, nil
}
