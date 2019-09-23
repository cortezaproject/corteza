package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/system/proto"
	"github.com/cortezaproject/corteza-server/system/types"
)

// gRPC client for

type (
	systemRole struct {
		client proto.RolesClient
	}
)

func SystemRole(c proto.RolesClient) *systemRole {
	return &systemRole{
		client: c,
	}
}

func (svc systemRole) Find(ctx context.Context, ID uint64) (rr types.RoleSet, err error) {
	rsp, err := svc.client.Find(ctx, &proto.FindRoleRequest{})
	if err != nil {
		return nil, err
	}

	rr = make([]*types.Role, len(rsp.Roles))

	for i := range rsp.Roles {
		rr[i] = &types.Role{
			ID:     rsp.Roles[i].ID,
			Name:   rsp.Roles[i].Name,
			Handle: rsp.Roles[i].Handle,
		}
	}

	return rr, nil
}
