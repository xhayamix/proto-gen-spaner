package spanner

import (
	"context"

	"cloud.google.com/go/spanner"

	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/database"
)

type roTx struct {
	value *spanner.ReadOnlyTransaction
}

func (tx *roTx) GetTxObject() any {
	return tx.value
}

func (tx *roTx) ReadOnlyImpl() {}

type rwTx struct {
	value *spanner.ReadWriteTransaction
}

func (tx *rwTx) GetTxObject() any {
	return tx.value
}

func (tx *rwTx) ReadOnlyImpl() {}

func (tx *rwTx) ReadWriteImpl() {}

// ReadTransaction spanner.ReadWriteTransactionとspanner.ReadOnlyTransactionはtxReadOnlyがそれぞれに埋め込まれているがinterfaceはなかったので定義
type ReadTransaction interface {
	ReadRow(ctx context.Context, table string, key spanner.Key, columns []string) (*spanner.Row, error)
	Read(ctx context.Context, table string, keys spanner.KeySet, columns []string) *spanner.RowIterator
	Query(ctx context.Context, statement spanner.Statement) *spanner.RowIterator
}

func ExtractROTx(tx database.ROTx) (ReadTransaction, error) {
	switch txObject := tx.GetTxObject().(type) {
	case *spanner.ReadOnlyTransaction:
		return txObject, nil
	case *spanner.ReadWriteTransaction:
		return txObject, nil
	}
	return nil, cerrors.New(cerrors.Internal)
}

func ExtractRWTx(tx database.RWTx) (*spanner.ReadWriteTransaction, error) {
	txObject, ok := tx.GetTxObject().(*spanner.ReadWriteTransaction)
	if !ok {
		return nil, cerrors.New(cerrors.Internal)
	}
	return txObject, nil
}
