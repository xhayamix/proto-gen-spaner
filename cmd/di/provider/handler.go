package provider

import (
	"github.com/google/wire"

	"github.com/xhayamix/proto-gen-spanner/pkg/cmd/api/handler/user"
)

var DefaultHandlerSet = wire.NewSet(
	user.New,
)
