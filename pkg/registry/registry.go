package registry

/*

気が向いたらcloserNewを実装します

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
	cspanner "github.com/xhayamix/proto-gen-spanner/pkg/infra/transaction"
	"github.com/xhayamix/proto-gen-spanner/pkg/util/closure"
)

type SpannerConfig interface {
    GetProjectID() string
	GetSpannerInstance() string
	GetSpannerDB() string
}

func NewSpannerConstructor(
	ctx context.Context,
	closeListener *closure.CloseListener,
	conf SpannerConfig,
) func() (*spanner.Client, error) {
	return func() (*spanner.Client, error) {
		client, err := cspanner.New(ctx, conf.GetProjectID(), conf.GetSpannerInstance(), conf.GetSpannerDB())
		if err != nil {
			return nil, cerrors.Stack(err)
		}
		closeListener.Add(func() {
			client.Close()
		})
		return client, nil
	}
}

*/
