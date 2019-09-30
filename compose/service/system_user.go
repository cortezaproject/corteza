package service

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/system/proto"
	"github.com/cortezaproject/corteza-server/system/types"
)

// gRPC client for

type (
	systemUser struct {
		client proto.UsersClient
	}
)

func SystemUser(c proto.UsersClient) *systemUser {
	return &systemUser{
		client: c,
	}
}

func (svc systemUser) MakeJWT(ctx context.Context, ID uint64) (string, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
		"jwt": []string{auth.GetJwtFromContext(ctx)},
	})

	rsp, err := svc.client.MakeJWT(ctx, &proto.MakeJWTUserRequest{UserID: ID}, grpc.WaitForReady(true))
	if err != nil {
		return "", err
	}

	return rsp.JWT, nil
}

func (svc systemUser) FindByID(ctx context.Context, ID uint64) (*types.User, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
		"jwt": []string{auth.GetJwtFromContext(ctx)},
	})

	rsp, err := svc.client.FindByID(ctx, &proto.FindByIDUserRequest{UserID: ID})
	if err != nil {
		return nil, err
	}

	return &types.User{
		ID:     rsp.User.ID,
		Email:  rsp.User.Email,
		Name:   rsp.User.Name,
		Handle: rsp.User.Handle,
		Kind:   types.UserKind(rsp.User.Kind),
	}, nil
}
