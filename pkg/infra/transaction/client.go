package spanner

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"

	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
)

func New(ctx context.Context, projectID, instance, db string) (*spanner.Client, error) {
	dbPath := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instance, db)
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
