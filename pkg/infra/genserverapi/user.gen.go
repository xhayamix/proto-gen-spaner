// Code generated by protoc-gen-all. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package genserverapi

import (
	"context"

	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/proto/client/api"
)

func (c *client) UserGetProfile(ctx context.Context, req *api.GetProfileRequest) (*api.GetProfileResponse, error) {
	conn, err := c.getConn()
	if err != nil {
		return nil, cerrors.Stack(err)
	}
	cli := api.NewUserClient(conn)
	result, err := cli.GetProfile(ctx, req)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	return result, nil
}
