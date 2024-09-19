package repository

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

//go:embed repository.gen.go.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(message *input.Message) (*output.TemplateInfo, error) {
	if !input.ServerMessageAccessorSet.Contains(message.Option.AccessorType) {
		return nil, nil
	}

	type Column struct {
		GoName    string
		LocalName string
		Type      string
		PK        bool
		IsList    bool
		IsEnum    bool
	}
	type Method struct {
		Name         string
		Args         string
		SliceArgName string
		ReturnName   string
		SelectType   string
		Wheres       string
		UseCache     bool
	}
	type Type struct {
		GoName  string
		Key     string
		Columns []*Column
	}
	type Table struct {
		PkgName            string
		GoName             string
		CamelName          string
		PKColumns          []*Column
		Methods            []*Method
		Types              []*Type
		NeedCommonResponse bool
	}
	data := &Table{
		PkgName:            core.ToPkgName(message.SnakeName),
		GoName:             core.ToGolangPascalCase(message.SnakeName),
		CamelName:          core.ToCamelCase(message.SnakeName),
		PKColumns:          make([]*Column, 0),
		Methods:            make([]*Method, 0),
		Types:              make([]*Type, 0),
		NeedCommonResponse: message.Option.AccessorType == input.MessageAccessorType_AllWithCommonResponse,
	}

	columns := make([]*Column, 0, len(message.Fields))
	columnMap := make(map[string]*Column)
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

		// type予約語との衝突回避
		goName := core.ToGolangPascalCase(field.SnakeName)
		localName := core.ToGolangCamelCase(field.SnakeName)
		if localName == "type" {
			localName = "typ"
		}

		column := &Column{
			GoName:    goName,
			LocalName: localName,
			Type:      typeName,
			PK:        field.Option.DDL.PK,
			IsList:    field.IsList,
			IsEnum:    field.TypeKind == input.TypeKind_Enum,
		}

		if field.Option.DDL.PK {
			data.PKColumns = append(data.PKColumns, column)
		}
		columns = append(columns, column)
		columnMap[field.SnakeName] = column
	}
	data.Types = append(data.Types, &Type{
		GoName:  data.GoName,
		Key:     "All",
		Columns: columns,
	})

	var methodName string
	var methodArgs string
	var methodWheres string
	last := len(data.PKColumns) - 1
	for i, column := range data.PKColumns {
		if i == last {
			break
		}
		argName := core.ToGolangCamelCase(column.GoName)
		methodName += column.GoName
		methodArgs += argName
		methodWheres += column.GoName
		data.Methods = append(data.Methods,
			&Method{
				Name:       methodName,
				Args:       methodArgs + " " + column.Type,
				ReturnName: data.GoName,
				SelectType: "All",
				Wheres:     methodWheres + "Eq(" + argName + ")",
				UseCache:   true,
			},
			&Method{
				Name:         methodName + "s",
				Args:         methodArgs + "s []" + column.Type,
				SliceArgName: argName + "s",
				ReturnName:   data.GoName,
				SelectType:   "All",
				Wheres:       methodWheres + "In(" + argName + "s)",
				UseCache:     true,
			},
		)
		methodName += "And"
		methodArgs += " " + column.Type + ", "
		methodWheres += "Eq(" + argName + ").And()."
	}

	for _, index := range message.Option.DDL.Indexes {
		typ := &Type{}
		keySet := strset.NewWithSize(len(index.Keys) + len(index.SnakeStoring))
		goNames := make([]string, 0, len(index.Keys))
		for _, key := range index.Keys {
			if column, ok := columnMap[key.SnakeName]; ok {
				goNames = append(goNames, column.GoName)
				keySet.Add(column.GoName)
			}
		}
		typ.Key = "Idx" + strings.Join(goNames, "")
		typ.GoName = data.GoName + typ.Key
		for _, key := range index.SnakeStoring {
			keySet.Add(core.ToGolangPascalCase(key))
		}

		// columns定義順にindex内のcolumnを整えるためcolumnsでforを回す
		for _, column := range columns {
			if column.PK || keySet.Has(column.GoName) {
				typ.Columns = append(typ.Columns, column)
			}
		}
		data.Types = append(data.Types, typ)

		methodName = typ.Key + "With"
		methodArgs = ""
		methodWheres = ""
		for _, key := range index.Keys {
			column, ok := columnMap[key.SnakeName]
			if !ok {
				continue
			}
			methodName += column.GoName
			methodArgs += column.GoName
			methodWheres += column.GoName
			data.Methods = append(data.Methods,
				&Method{
					Name:       methodName,
					Args:       methodArgs + " " + column.Type,
					ReturnName: data.GoName + typ.Key,
					SelectType: typ.Key,
					Wheres:     methodWheres + "Eq(" + column.GoName + ")",
					UseCache:   false,
				},
				&Method{
					Name:       methodName + "s",
					Args:       methodArgs + "s []" + column.Type,
					ReturnName: data.GoName + typ.Key,
					SelectType: typ.Key,
					Wheres:     methodWheres + "In(" + column.GoName + "s)",
					UseCache:   false,
				},
			)
			methodName += "And"
			methodArgs += " " + column.Type + ", "
			methodWheres += "Eq(" + column.GoName + ").And()."
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
		FilePath: core.JoinPath("pkg/infra/transaction/repository", message.SnakeName+"_repository.gen.go"),
	}, nil
}
