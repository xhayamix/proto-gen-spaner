package user

import (
	"context"
	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"

	"github.com/xhayamix/proto-gen-spanner/pkg/domain/database"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/dto"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/service/user"
)

type Interactor interface {
	GetProfile(ctx context.Context, userID string) (*dto.User, error)
	SaveUser(ctx context.Context, userID, publicUserID string) error
}

type interactor struct {
	txManager   database.TxManager
	userService user.Service
}

func New(
	txManager database.TxManager,
	userService user.Service,
) Interactor {
	return &interactor{
		txManager:   txManager,
		userService: userService,
	}
}

func (i *interactor) GetProfile(ctx context.Context, userID string) (*dto.User, error) {
	var user *dto.User
	if err := i.txManager.ReadOnlyTransaction(ctx, func(ctx context.Context, tx database.ROTx) error {
		var err error
		user, err = i.userService.GetProfile(ctx, tx, userID)
		if err != nil {
			return cerrors.Stack(err)
		}
		return nil
	}); err != nil {
		return nil, cerrors.Stack(err)
	}

	return user, nil
}

func (i *interactor) SaveUser(ctx context.Context, userID, publicUserID string) error {
	return i.txManager.Transaction(ctx, func(ctx context.Context, tx database.RWTx) error {
		return i.userService.SaveUser(ctx, tx, userID, publicUserID)
	})
}
