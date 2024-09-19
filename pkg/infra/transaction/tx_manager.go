package spanner

import (
	"cloud.google.com/go/spanner"
	"context"

	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/database"
)

func NewSpannerUserTxManager(c *spanner.Client) database.TxManager {
	return &txManager{
		client: c,
	}
}

type txManager struct {
	client *spanner.Client
}

func (tm *txManager) ReadOnlyTransaction(ctx context.Context, f func(ctx context.Context, tx database.ROTx) error) error {
	transaction := tm.client.ReadOnlyTransaction()
	defer transaction.Close()

	tx := &roTx{value: transaction}
	err := f(ctx, tx)
	if err != nil {
		return cerrors.Stack(err)
	}

	return nil
}

func (tm *txManager) Transaction(ctx context.Context, f func(ctx context.Context, tx database.RWTx) error) error {
	_, err := tm.client.ReadWriteTransaction(ctx, func(ctx context.Context, spannerRWTxn *spanner.ReadWriteTransaction) error {
		tx := &rwTx{value: spannerRWTxn}
		err := f(ctx, tx)
		if err != nil {
			return cerrors.Stack(err)
		}
		return nil
	})
	return err
}
