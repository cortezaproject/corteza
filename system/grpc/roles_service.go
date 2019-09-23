package grpc

import (
	"context"

	"github.com/cortezaproject/corteza-server/system/proto"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	roleService struct {
		roles service.RoleService
	}
)

func NewRoleService(roles service.RoleService) *roleService {
	return &roleService{
		roles: roles,
	}
}

func (gs roleService) Find(ctx context.Context, req *proto.FindRoleRequest) (rsp *proto.FindRoleResponse, err error) {
	var (
		rr types.RoleSet
	)

	if rr, err = gs.roles.Find(&types.RoleFilter{}); err != nil {
		return
	}

	rsp = &proto.FindRoleResponse{
		Roles: make([]*proto.Role, len(rr)),
	}

	for i := range rr {
		rsp.Roles[i] = &proto.Role{
			ID:     rr[i].ID,
			Handle: rr[i].Handle,
			Name:   rr[i].Name,
		}
	}

	return
}
