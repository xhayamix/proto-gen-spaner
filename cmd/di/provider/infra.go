package provider

import (
	"github.com/google/wire"

	"github.com/xhayamix/proto-gen-spanner/pkg/infra/transaction"
	"github.com/xhayamix/proto-gen-spanner/pkg/infra/transaction/repository"
)

var DefaultInfraSet = wire.NewSet(
	spanner.NewSpannerUserTxManager,
	repository.NewUserRepository,
	spanner.New,
)
