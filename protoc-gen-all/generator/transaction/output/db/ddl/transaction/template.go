package transaction

import (
	"bytes"
	_ "embed"
	"sort"
	"text/template"

	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/core"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/input"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/generator/transaction/output"
	"github.com/xhayamix/proto-gen-spanner/protoc-gen-all/perrors"
)

//go:embed init.gen.sql.tpl
var templateFileBytes []byte

type Creator struct{}

func (c *Creator) Create(messages []*input.Message, enumMap map[string]*input.Enum) (*output.TemplateInfo, error) {
	type Column struct {
		GoName      string
		Type        string
		PK          bool
		Comment     string
		CommentInfo string
	}
	type IndexKey struct {
		GoName string
		Desc   bool
	}
	type Index struct {
		Keys    []*IndexKey
		Unique  bool
		Storing []string
	}
	type DeletionPolicy struct {
		TimestampColumn string
		Days            int32
	}
	type Table struct {
		GoName          string
		Columns         []*Column
		PKColumns       []*Column
		Indexes         []*Index
		InterleaveTable string
		DeletionPolicy  *DeletionPolicy
		Comment         string
		CommentInfo     string
	}
	type Data struct {
		Tables []*Table
	}

	childrenMap := make(map[string][]*Table)
	for _, message := range messages {
		if !input.ServerMessageAccessorSet.Contains(message.Option.AccessorType) {
			continue
		}

		commentInfo, err := c.getTableCommentInfo(message.Option.InsertTiming)
		if err != nil {
			return nil, perrors.Stack(err)
		}

		table := &Table{
			GoName:      core.ToGolangPascalCase(message.SnakeName),
			Columns:     make([]*Column, 0, len(message.Fields)),
			PKColumns:   make([]*Column, 0),
			Indexes:     make([]*Index, 0, len(message.Option.DDL.Indexes)),
			Comment:     message.Comment,
			CommentInfo: commentInfo,
		}
		indexColumnMap := make(map[string]*Column)

		for _, field := range message.Fields {
			if !input.AdminFieldAccessorSet.Contains(field.Option.AccessorType) {
				continue
			}

			/* 型の整形 */
			typeName, err := input.GetSpannerType(field.TypeKind, field.SnakeName, field.IsList)
			if err != nil {
				return nil, perrors.Stack(err)
			}

			ci, err := c.getColumnCommentInfo(enumMap, field.TypeKind, field.Type)
			if err != nil {
				return nil, perrors.Stack(err)
			}
			pk := field.Option.DDL.PK
			column := &Column{
				GoName:      core.ToGolangPascalCase(field.SnakeName),
				Type:        typeName,
				PK:          pk,
				Comment:     field.Comment,
				CommentInfo: ci,
			}

			table.Columns = append(table.Columns, column)
			if pk {
				table.PKColumns = append(table.PKColumns, column)
			}

			indexColumnMap[field.SnakeName] = column
		}

		for _, index := range message.Option.DDL.Indexes {
			i := &Index{
				Keys:    make([]*IndexKey, 0, len(index.Keys)),
				Unique:  index.Unique,
				Storing: make([]string, 0, len(index.SnakeStoring)),
			}
			for _, key := range index.Keys {
				if column, ok := indexColumnMap[key.SnakeName]; ok {
					i.Keys = append(i.Keys, &IndexKey{
						GoName: column.GoName,
						Desc:   key.Desc,
					})
				}
			}
			for _, s := range index.SnakeStoring {
				if column, ok := indexColumnMap[s]; ok {
					i.Storing = append(i.Storing, column.GoName)
				}
			}
			table.Indexes = append(table.Indexes, i)
		}
		if message.Option.DDL.Interleave != nil {
			table.InterleaveTable = core.ToGolangPascalCase(message.Option.DDL.Interleave.TableSnakeName)
		}
		if message.Option.DDL.TTL != nil {
			table.DeletionPolicy = &DeletionPolicy{
				TimestampColumn: core.ToGolangPascalCase(message.Option.DDL.TTL.TimestampColumnSnakeName),
				Days:            message.Option.DDL.TTL.Days,
			}
		}

		childrenMap[table.InterleaveTable] = append(childrenMap[table.InterleaveTable], table)
	}

	tables := make([]*Table, 0, len(messages))
	// 再起的にツリー状にappendしていく
	var appendTree func(table *Table)
	appendTree = func(table *Table) {
		tables = append(tables, table)
		for _, child := range childrenMap[table.GoName] {
			appendTree(child)
		}
	}

	for _, children := range childrenMap {
		sort.Slice(children, func(i, j int) bool {
			return children[i].GoName < children[j].GoName
		})
	}
	// keyが空文字の配列は最上位の親
	for _, parent := range childrenMap[""] {
		appendTree(parent)
	}

	tpl, err := template.New("").Parse(string(templateFileBytes))
	if err != nil {
		return nil, perrors.Stack(err)
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, &Data{
		Tables: tables,
	}); err != nil {
		return nil, perrors.Stack(err)
	}

	return &output.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: core.JoinPath("db/ddl/transaction", "init.gen.sql"),
	}, nil
}

func (c *Creator) getTableCommentInfo(insertTiming string) (string, error) {
	commentInfo, err := (&core.DDLTableCommentInfo{
		InsertTiming: insertTiming,
	}).String()
	if err != nil {
		return "", perrors.Stack(err)
	}

	return commentInfo, nil
}

func (c *Creator) getColumnCommentInfo(enumMap map[string]*input.Enum, fieldTypeKind input.TypeKind, fieldType input.FieldType) (string, error) {
	if fieldTypeKind != input.TypeKind_Enum {
		return "", nil
	}

	e, ok := enumMap[fieldType]
	if !ok {
		return "", nil
	}

	commentInfo, err := (&core.DDLColumnCommentInfo{
		EnumName:    e.Name,
		EnumInfoMap: e.CommentMapByValue,
	}).String()
	if err != nil {
		return "", perrors.Stack(err)
	}

	return commentInfo, nil
}
