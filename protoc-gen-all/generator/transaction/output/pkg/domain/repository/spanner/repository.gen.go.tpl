{{ template "autogen_comment" }}
//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_{{ .SnakeName }}.go
//go:generate goimports -w --local "github.com/xhayamix/proto-gen-spanner" mock_$GOPACKAGE/mock_{{ .SnakeName }}.go
{{ $name := .GoName }}
package transaction

import (
	"context"

	"github.com/xhayamix/proto-gen-spanner/pkg/domain/database"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/entity/transaction"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/enum"
)

type {{ .GoName }}Repository interface {
	LoadByPK(ctx context.Context, tx database.ROTx, key *transaction.{{ .GoName }}PK) (*transaction.{{ .GoName }}, error)
	LoadByPKs(ctx context.Context, tx database.ROTx, keys transaction.{{ .GoName }}PKs) (transaction.{{ .GoName }}Slice, error)
	SelectByPK(ctx context.Context, tx database.ROTx, key *transaction.{{ .GoName }}PK) (*transaction.{{ .GoName }}, error)
	SelectByPKs(ctx context.Context, tx database.ROTx, keys transaction.{{ .GoName }}PKs) (transaction.{{ .GoName }}Slice, error)
	// SelectAll not use cache. ReadOnlyTx時のみ使用可能
	SelectAll(ctx context.Context, tx database.ROTx, limit, offset int32) (transaction.{{ .GoName }}Slice, error)
	{{ range .Methods -}}
	SelectBy{{ .Name }}(ctx context.Context, tx database.ROTx, {{ .Args }}) (transaction.{{ .ReturnName }}Slice, error)
	{{ end -}}
	// Search not use cache. ReadOnlyTx時のみ使用可能
	Search(ctx context.Context, tx database.ROTx, search string, limit, offset int32) (transaction.{{ .GoName }}Slice, error)
	Insert(ctx context.Context, tx database.RWTx, entity *transaction.{{ .GoName }}) error
	BulkInsert(ctx context.Context, tx database.RWTx, entities transaction.{{ .GoName }}Slice) error
	Update(ctx context.Context, tx database.RWTx, entity *transaction.{{ .GoName }}) error
	Save(ctx context.Context, tx database.RWTx, entity *transaction.{{ .GoName }}) error
	Delete(ctx context.Context, tx database.RWTx, key *transaction.{{ .GoName }}PK) error
	BulkDelete(ctx context.Context, tx database.RWTx, keys transaction.{{ .GoName }}PKs) error
}
