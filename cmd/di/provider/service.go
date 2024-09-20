package provider

import (
	"github.com/google/wire"

	"github.com/xhayamix/proto-gen-spanner/pkg/domain/service/user"
)

var DefaultServiceSet = wire.NewSet(
	user.New,
)
