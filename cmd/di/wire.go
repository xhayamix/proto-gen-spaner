//go:build wireinject
// +build wireinject

package di

import (
	"context"

	"github.com/google/wire"

	"github.com/xhayamix/proto-gen-spanner/cmd/di/provider"
)

func InitializeApp(ctx context.Context) (*provider.App, error) {
	wire.Build(provider.DefaultAppSet)
	return nil, nil
}
