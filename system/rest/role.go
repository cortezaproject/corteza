package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/system/rest/request"
	"github.com/crusttech/crust/system/service"
	"github.com/crusttech/crust/system/types"
)

var _ = errors.Wrap

type (
	Role struct {
		svc struct {
			role service.RoleService
		}
	}
)

func (Role) New() *Role {
	ctrl := &Role{}
	ctrl.svc.role = service.DefaultRole
	return ctrl
}

func (ctrl *Role) Read(ctx context.Context, r *request.RoleRead) (interface{}, error) {
	return ctrl.svc.role.With(ctx).FindByID(r.RoleID)
}

func (ctrl *Role) List(ctx context.Context, r *request.RoleList) (interface{}, error) {
	return ctrl.svc.role.With(ctx).Find(&types.RoleFilter{Query: r.Query})
}

func (ctrl *Role) Create(ctx context.Context, r *request.RoleCreate) (interface{}, error) {
	org := &types.Role{
		Name: r.Name,
	}

	return ctrl.svc.role.With(ctx).Create(org)
}

func (ctrl *Role) Update(ctx context.Context, r *request.RoleUpdate) (interface{}, error) {
	org := &types.Role{
		ID:   r.RoleID,
		Name: r.Name,
	}

	return ctrl.svc.role.With(ctx).Update(org)
}

func (ctrl *Role) Remove(ctx context.Context, r *request.RoleRemove) (interface{}, error) {
	return nil, ctrl.svc.role.With(ctx).Delete(r.RoleID)
}

func (ctrl *Role) Archive(ctx context.Context, r *request.RoleArchive) (interface{}, error) {
	return nil, ctrl.svc.role.With(ctx).Archive(r.RoleID)
}

func (ctrl *Role) Merge(ctx context.Context, r *request.RoleMerge) (interface{}, error) {
	return nil, ctrl.svc.role.With(ctx).Merge(r.RoleID, r.Destination)
}

func (ctrl *Role) Move(ctx context.Context, r *request.RoleMove) (interface{}, error) {
	return nil, ctrl.svc.role.With(ctx).Move(r.RoleID, r.OrganisationID)
}

func (ctrl *Role) MemberAdd(ctx context.Context, r *request.RoleMemberAdd) (interface{}, error) {
	return nil, ctrl.svc.role.With(ctx).MemberAdd(r.RoleID, r.UserID)
}

func (ctrl *Role) MemberRemove(ctx context.Context, r *request.RoleMemberRemove) (interface{}, error) {
	return nil, ctrl.svc.role.With(ctx).MemberRemove(r.RoleID, r.UserID)
}
