package user

import (
	"context"
	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/dto"

	"github.com/xhayamix/proto-gen-spanner/pkg/domain/database"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/entity/transaction"
	transactionrepository "github.com/xhayamix/proto-gen-spanner/pkg/domain/repository/transaction"
)

type Service interface {
	GetProfile(ctx context.Context, tx database.ROTx, userID string) (*dto.User, error)
}

type service struct {
	userRepository transactionrepository.UserRepository
}

func New(
	userRepository transactionrepository.UserRepository,
) Service {
	return &service{
		userRepository: userRepository,
	}
}

func (s *service) GetProfile(ctx context.Context, tx database.ROTx, userID string) (*dto.User, error) {
	user, err := s.userRepository.LoadByPK(ctx, tx, &transaction.UserPK{
		UserID: userID,
	})
	if err != nil {
		return nil, cerrors.Stack(err)
	}
	return &dto.User{
		UserID:       user.UserID,
		PublicUserID: user.PublicUserID,
	}, nil
}
