package spanner

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"

	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/entity"
)

func New(ctx context.Context) (*spanner.Client, error) {
	conf := entity.GetConfig()
	dbPath := fmt.Sprintf("projects/%s/instances/%s/databases/%s", conf.ProjectID, conf.Instance, conf.Instance)
	// TODO: config
	config := spanner.ClientConfig{
		SessionPoolConfig: spanner.DefaultSessionPoolConfig,
	}
	client, err := spanner.NewClientWithConfig(ctx, dbPath, config)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	return client, nil
}
