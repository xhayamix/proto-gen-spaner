//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
//go:generate goimports -w --local "github.com/xhayamix/proto-gen-spanner" mock_$GOPACKAGE/mock_$GOFILE
package database

import (
	"context"
)

type ROTx interface {
	// GetTxObject トランザクションの実態を取得する
	GetTxObject() any
	// ReadOnlyImpl 型パズル用
	ReadOnlyImpl()
}

type RWTx interface {
	ROTx
	// ReadWriteImpl 型パズル用
	ReadWriteImpl()
}

type TxManager interface {
	ReadOnlyTransaction(ctx context.Context, f func(ctx context.Context, tx ROTx) error) error
	Transaction(ctx context.Context, f func(ctx context.Context, tx RWTx) error) error
}

type UserTxManager TxManager
