package user

import (
	"context"

	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
	"github.com/xhayamix/proto-gen-spanner/pkg/cmd/api/intaractor/user"
	pb "github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/api"
)

type Handler struct {
	interactor user.Interactor
}

func New(
	interactor user.Interactor,
) pb.UserServer {
	return &Handler{
		interactor: interactor,
	}
}

func (h *Handler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, cerrors.Wrapf(err, cerrors.InvalidArgument, "バリデーションエラーが発生しました")
	}

	profile, err := h.interactor.GetProfile(ctx, req.GetUserId())
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	return &pb.GetProfileResponse{
		Profile: &pb.Profile{
			UserId:       profile.UserID,
			PublicUserId: profile.PublicUserID,
		},
	}, nil
}
