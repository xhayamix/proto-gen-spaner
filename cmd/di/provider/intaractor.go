package provider

import (
	"github.com/google/wire"

	"github.com/xhayamix/proto-gen-spanner/pkg/cmd/api/intaractor/user"
)

var DefaultIntaractorSet = wire.NewSet(
	user.New,
)
