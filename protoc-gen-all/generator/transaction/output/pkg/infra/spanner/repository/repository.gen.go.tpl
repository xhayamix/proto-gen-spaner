{{ template "autogen_comment" }}
{{ $name := .GoName -}}
{{ $camelName := .CamelName -}}
package repository

import (
	"context"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/scylladb/go-set/strset"
	"google.golang.org/grpc/codes"

	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/database"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/entity/transaction"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/enum"
	repository "github.com/xhayamix/proto-gen-spanner/pkg/domain/repository/transaction"
	cspanner "github.com/xhayamix/proto-gen-spanner/pkg/infra/transaction"
	"github.com/xhayamix/proto-gen-spanner/pkg/infra/transaction/repository/base"
)

type {{ .CamelName }}Repository struct {}

func New{{ .GoName }}Repository() repository.{{ .GoName }}Repository {
	return &{{ .CamelName }}Repository{}
}

/*
func (r *{{ .CamelName }}Repository) extractQueryCache(ctx context.Context) (base.{{ .GoName }}SearchResultCache, base.{{ .GoName }}MutationWaitBuffer) {
	return base.Extract{{ .GoName }}SearchResultCache(ctx), base.Extract{{ .GoName }}MutationWaitBuffer(ctx)
}
*/

func (r *{{ .CamelName }}Repository) LoadByPK(ctx context.Context, tx database.ROTx, pk *transaction.{{ .GoName }}PK) (*transaction.{{ .GoName }}, error) {
	row, err := r.SelectByPK(ctx, tx, pk)
	if err != nil {
		return nil, cerrors.Stack(err)
	}
	if row == nil {
		return nil, cerrors.Newf(cerrors.InvalidArgument, "{{ .GoName }}が見つかりません。 pk = %s", pk)
	}

	return row, nil
}

func (r *{{ .CamelName }}Repository) LoadByPKs(ctx context.Context, tx database.ROTx, pks transaction.{{ .GoName }}PKs) (transaction.{{ .GoName }}Slice, error) {
	rows, err := r.SelectByPKs(ctx, tx, pks)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	set := strset.NewWithSize(len(rows))
	for _, row := range rows {
		set.Add(row.GetPK().Key())
	}

	notFoundPKs := make(transaction.{{ .GoName }}PKs, 0, len(pks))
	for _, pk := range pks {
		if !set.Has(pk.Key()) {
			notFoundPKs = append(notFoundPKs, pk)
		}
	}
	if len(notFoundPKs) > 0 {
		return nil, cerrors.Newf(cerrors.InvalidArgument, "{{ .GoName }}が見つかりません。 pks = %s", notFoundPKs)
	}

	return rows, nil
}

func (r *{{ .CamelName }}Repository) SelectByPK(ctx context.Context, tx database.ROTx, pk *transaction.{{ .GoName }}PK) (entity *transaction.{{ .GoName }}, err error) {
	rtx, err := cspanner.ExtractROTx(tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	row, err := rtx.ReadRow(ctx, transaction.{{ .GoName }}TableName, spanner.Key(pk.Generate()), transaction.{{ .GoName }}ColumnNameSlice)
	if err != nil {
		// PKでのSelect結果がなかった場合
		if spanner.ErrCode(err) == codes.NotFound {
			return nil, nil
		}
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	entity, err = r.decodeAllColumns(row)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	return entity, nil
}

func (r *{{ .CamelName }}Repository) SelectByPKs(ctx context.Context, tx database.ROTx, pks transaction.{{ .GoName }}PKs) (rows transaction.{{ .GoName }}Slice, err error) {
	if len(pks) == 0 {
		return transaction.{{ .GoName }}Slice{}, nil
	}

	rtx, err := cspanner.ExtractROTx(tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	keySets := make([]spanner.KeySet, 0, len(pks))
	for _, pk := range pks {
		keySets = append(keySets, spanner.Key(pk.Generate()))
	}
	ri := rtx.Read(ctx, transaction.{{ .GoName }}TableName, spanner.KeySets(keySets...), transaction.{{ .GoName }}ColumnNameSlice)
	rows = make(transaction.{{ .GoName }}Slice, 0)
	keySet := strset.New()
	if err := ri.Do(func(row *spanner.Row) error {
		if len(rows) == 0 {
			rows = make(transaction.{{ .GoName }}Slice, 0, ri.RowCount)
			keySet = strset.NewWithSize(int(ri.RowCount))
		}
		entity, err := r.decodeAllColumns(row)
		if err != nil {
			return cerrors.Stack(err)
		}
		rows = append(rows, entity)
		keySet.Add(entity.GetPK().Key())
		return nil
	}); err != nil {
		if err, ok := cerrors.As(err); ok {
			return nil, cerrors.Stack(err)
		}
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	return rows, nil
}

func (r *{{ .CamelName }}Repository) SelectAll(ctx context.Context, tx database.ROTx, limit, offset int32) (rows transaction.{{ .GoName }}Slice, err error) {
    roTx, err := cspanner.ExtractROTx(tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	sql, params := base.New{{ .GoName }}QueryBuilder().
		SelectAllFrom{{ .GoName }}().
		OrderBy(base.OrderPairs{ {{- range $i, $col := .PKColumns }}{{ if $i }}, {{ end }}{"{{ $col.GoName }}", base.OrderTypeASC}{{ end -}} }).
		Limit(limit).
		Offset(offset).
		GetQuery()
	stmt := spanner.Statement{
		SQL: sql,
		Params: params,
	}
	ri := roTx.Query(ctx, stmt)

	rows = make(transaction.{{ .GoName }}Slice, 0)
	if err := ri.Do(func(row *spanner.Row) error {
		if len(rows) == 0 {
			rows = make(transaction.{{ .GoName }}Slice, 0, ri.RowCount)
		}
		entity, err := r.decodeAllColumns(row)
		if err != nil {
			return cerrors.Stack(err)
		}
		rows = append(rows, entity)
		return nil
	}); err != nil {
		if err, ok := cerrors.As(err); ok {
			return nil, cerrors.Stack(err)
		}
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	return rows, nil
}
{{- range .Methods }}

func (r *{{ $camelName }}Repository) SelectBy{{ .Name }}(ctx context.Context, tx database.ROTx, {{ .Args }}) (rows transaction.{{ .ReturnName }}Slice, err error) {
	{{ if .SliceArgName }}
	if len({{ .SliceArgName }}) == 0 {
		return transaction.{{ .ReturnName }}Slice{}, nil
	}

	{{ end -}}

	roTx, err := cspanner.ExtractROTx(tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	qb := base.New{{ $name }}QueryBuilder().
		Select{{ .SelectType }}From{{ $name }}().
		Where().{{ .Wheres }}

	sql, params := qb.GetQuery()
	stmt := spanner.Statement{
		SQL: sql,
		Params: params,
	}
	ri := roTx.Query(ctx, stmt)

	rows = make(transaction.{{ .ReturnName }}Slice, 0)
	{{ if .UseCache }}keySet := strset.New(){{ end }}
	if err := ri.Do(func(row *spanner.Row) error {
		if len(rows) == 0 {
			rows = make(transaction.{{ .ReturnName }}Slice, 0, ri.RowCount)
			{{ if .UseCache }}keySet = strset.NewWithSize(int(ri.RowCount)){{ end }}
		}
		entity, err := r.decode{{ .SelectType }}Columns(row)
		if err != nil {
			return cerrors.Stack(err)
		}
		rows = append(rows, entity)
		{{ if .UseCache }}keySet.Add(entity.GetPK().Key()){{ end }}
		return nil
	}); err != nil {
		if err, ok := cerrors.As(err); ok {
			return nil, cerrors.Stack(err)
		}
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	return rows, nil
}
{{- end }}

func (r *{{ .CamelName }}Repository) Search(ctx context.Context, tx database.ROTx, search string, limit, offset int32) (rows transaction.{{ .GoName }}Slice, err error) {
	{{ $firstPK := mustFirst .PKColumns -}}
	roTx, err := cspanner.ExtractROTx(tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

    {{ if $firstPK.IsEnum -}}
    searchTypeInt, ok := {{ $firstPK.Type }}Map[search]
    if !ok {
       return nil, cerrors.Newf(cerrors.InvalidArgument, "存在しないenumです。 search = %s", search)
    }
    searchParam := {{ $firstPK.Type }}(searchTypeInt)
    {{ else }}
    searchParam := search
    {{ end }}

	sql, params := base.New{{ $name }}QueryBuilder().
		SelectAllFrom{{ $name }}().
		Where().{{ $firstPK.GoName }}Eq(searchParam).
		OrderBy(base.OrderPairs{ {{- range $i, $col := .PKColumns }}{{ if $i }}, {{ end }}{"{{ $col.GoName }}", base.OrderTypeASC}{{ end -}} }).
		Limit(limit).
		Offset(offset).
		GetQuery()
	stmt := spanner.Statement{
		SQL: sql,
		Params: params,
	}
	ri := roTx.Query(ctx, stmt)

	rows = make(transaction.{{ .GoName }}Slice, 0)
	if err := ri.Do(func(row *spanner.Row) error {
		if len(rows) == 0 {
			rows = make(transaction.{{ .GoName }}Slice, 0, ri.RowCount)
		}
		entity, err := r.decodeAllColumns(row)
		if err != nil {
			return cerrors.Stack(err)
		}
		rows = append(rows, entity)
		return nil
	}); err != nil {
		if err, ok := cerrors.As(err); ok {
			return nil, cerrors.Stack(err)
		}
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	return rows, nil
}

func (r *{{ .CamelName }}Repository) Insert(ctx context.Context, tx database.RWTx, entity *transaction.{{ .GoName }}) (err error) {
    now := time.Now()
    entity.CreatedTime = now
    entity.UpdatedTime = now

    rwTx, err := cspanner.ExtractRWTx(tx)
    if err != nil {
        return cerrors.Stack(err)
    }

    {{- range .Types }}
    mutation := spanner.Insert(transaction.{{ .GoName }}TableName,
        []string{
           {{- range .Columns }}
           "{{ .GoName }}",
           {{- end }}
        },
        []interface{}{
           {{- range .Columns }}
           entity.{{ .GoName }},
           {{- end }}
        },
    )
    {{- end }}
    if err := rwTx.BufferWrite([]*spanner.Mutation{mutation}); err != nil {
        return err
    }

    return nil
}

// 気が向いたら作る(必要になったら)
func (r *{{ .CamelName }}Repository) BulkInsert(ctx context.Context, tx database.RWTx, entities transaction.{{ .GoName }}Slice) (err error) {
	if len(entities) == 0 {
		return nil
	}

	return nil
}

func (r *{{ .CamelName }}Repository) Update(ctx context.Context, tx database.RWTx, entity *transaction.{{ .GoName }}) (err error) {
    now := time.Now()
    entity.UpdatedTime = now

    rwTx, err := cspanner.ExtractRWTx(tx)
    if err != nil {
        return cerrors.Stack(err)
    }

    {{- range .Types }}
    mutation := spanner.Update(transaction.{{ .GoName }}TableName,
        []string{
           {{- range .Columns }}
           "{{ .GoName }}",
           {{- end }}
        },
        []interface{}{
           {{- range .Columns }}
           entity.{{ .GoName }},
           {{- end }}
        },
    )
    {{- end }}
    if err := rwTx.BufferWrite([]*spanner.Mutation{mutation}); err != nil {
        return err
    }

	return nil
}

func (r *{{ .CamelName }}Repository) Save(ctx context.Context, tx database.RWTx, entity *transaction.{{ .GoName }}) error {
	var err error
	if entity.CreatedTime.IsZero() {
		err = r.Insert(ctx, tx, entity)
	} else {
		err = r.Update(ctx, tx, entity)
	}
	if err != nil {
		return cerrors.Stack(err)
	}

	return nil
}

// 気が向いたら作る(必要になったら)
func (r *{{ .CamelName }}Repository) Delete(ctx context.Context, tx database.RWTx, pk *transaction.{{ .GoName }}PK) (err error) {

	return nil
}

// 気が向いたら作る(必要になったら)
func (r *{{ .CamelName }}Repository) BulkDelete(ctx context.Context, tx database.RWTx, pks transaction.{{ .GoName }}PKs) (err error) {

	return nil
}
{{- range .Types }}

func (r *{{ $camelName }}Repository) decode{{ .Key }}Columns(row *spanner.Row) (*transaction.{{ .GoName }}, error) {
	{{- range .Columns }}
	{{ $slice := hasPrefix "[]" .Type -}}
	{{ if eq .Type "time.Time" -}}
	var {{ .LocalName }} spanner.NullTime
	{{- else if or $slice -}}
	var {{ .LocalName }} {{ .Type }}
	{{- else if hasPrefix "enum." .Type -}}
	var {{ .LocalName }} {{ if and .IsList .IsEnum }}[]{{ end }}spanner.NullInt64
	{{- else -}}
	var {{ .LocalName }} spanner.Null{{ title .Type }}
	{{- end -}}
	{{- end }}

	if err := row.Columns(
		{{ range .Columns -}}
		&{{ .LocalName }},
		{{ end -}}
	); err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	var result transaction.{{ .GoName }}
	{{- range .Columns }}
	{{ $slice := hasPrefix "[]" .Type -}}
	{{ if or $slice -}}
	result.{{ .GoName }} = {{ .LocalName }}
	{{- else if and .IsList .IsEnum -}}
	{
		s := make({{ .Type }}, 0, len({{ .LocalName }}))
		for _, v := range {{ .LocalName }} {
			s = append(s, {{ trimSuffix "Slice" .Type }}(v.Int64))
		}
		result.{{ .GoName }} = s
	}
	{{- else -}}
	if {{ .LocalName }}.Valid {
		{{- if eq .Type "string" }}
		result.{{ .GoName }} = {{ .LocalName }}.StringVal
		{{- else if eq .Type "int64" -}}
		result.{{ .GoName }} = {{ .LocalName }}.Int64
		{{- else if eq .Type "bool" -}}
		result.{{ .GoName }} = {{ .LocalName }}.Bool
		{{- else if eq .Type "time.Time" -}}
		result.{{ .GoName }} = {{ .LocalName }}.Time.In(time.Local)
		{{- else if hasPrefix "enum." .Type -}}
		result.{{ .GoName }} = {{ .Type }}({{ .LocalName }}.Int64)
		{{- end }}
	}
	{{- end -}}
	{{- end }}
	return &result, nil
}
{{- end }}

func (r *{{ $camelName }}Repository) diffEntity(source, target *transaction.{{ .GoName }}) map[string]any {
	result := make(map[string]any)

	// PKの差分は取らない
	{{ range (index .Types 0).Columns -}}
	{{- if .PK }}{{ continue }}
	{{ else if eq .Type "time.Time" -}}
	if !source.{{ .GoName }}.Equal(target.{{ .GoName }}) {
		result["{{ .GoName }}"] = target.{{ .GoName }}
	}
	{{ else if eq .Type "[]byte" -}}
	if !bytes.Equal(source.{{ .GoName }}, target.{{ .GoName }}) {
		result["{{ .GoName }}"] = target.{{ .GoName }}
	}
	{{ else if .IsList -}}
	if !slices.Equal(source.{{ .GoName }}, target.{{ .GoName }}) {
		result["{{ .GoName }}"] = target.{{ .GoName }}
	}
	{{ else -}}
	if source.{{ .GoName }} != target.{{ .GoName }} {
		result["{{ .GoName }}"] = target.{{ .GoName }}
	}
	{{ end -}}
	{{ end }}

	return result
}
