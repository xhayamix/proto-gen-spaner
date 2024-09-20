package provider

import (
	"cloud.google.com/go/spanner"
	"github.com/google/wire"
	pb "github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/api"
)

var DefaultAppSet = wire.NewSet(
	DefaultHandlerSet,
	DefaultIntaractorSet,
	DefaultServiceSet,
	DefaultInfraSet,
	NewApp,
)

type App struct {
	UserHandler pb.UserServer
	SpannerDB   *spanner.Client
}

func NewApp(
	UserHandler pb.UserServer,
	SpannerDB *spanner.Client,
) *App {
	return &App{
		UserHandler: UserHandler,
		SpannerDB:   SpannerDB,
	}
}
