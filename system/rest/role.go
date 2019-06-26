package rest

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/internal/payload"
	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/types"
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
	role := &types.Role{
		Name: r.Name,
	}

	role, err := ctrl.svc.role.With(ctx).Create(role)
	if err != nil {
		return nil, err
	}

	for _, userID := range payload.ParseUInt64s(r.Members) {
		err := ctrl.svc.role.With(ctx).MemberAdd(role.ID, userID)
		if err != nil {
			return nil, err
		}
	}
	return role, nil
}

func (ctrl *Role) Update(ctx context.Context, r *request.RoleUpdate) (interface{}, error) {
	role := &types.Role{
		ID:   r.RoleID,
		Name: r.Name,
	}

	role, err := ctrl.svc.role.With(ctx).Update(role)
	if err != nil {
		return nil, err
	}

	if len(r.Members) > 0 {
		members, err := ctrl.svc.role.With(ctx).MemberList(r.RoleID)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			err := ctrl.svc.role.With(ctx).MemberRemove(role.ID, member.UserID)
			if err != nil {
				return nil, err
			}
		}

		for _, userID := range payload.ParseUInt64s(r.Members) {
			err := ctrl.svc.role.With(ctx).MemberAdd(role.ID, userID)
			if err != nil {
				return nil, err
			}
		}
	}
	return role, nil
}

func (ctrl *Role) Delete(ctx context.Context, r *request.RoleDelete) (interface{}, error) {
	return resputil.OK(), ctrl.svc.role.With(ctx).Delete(r.RoleID)
}

func (ctrl *Role) Archive(ctx context.Context, r *request.RoleArchive) (interface{}, error) {
	return resputil.OK(), ctrl.svc.role.With(ctx).Archive(r.RoleID)
}

func (ctrl *Role) Merge(ctx context.Context, r *request.RoleMerge) (interface{}, error) {
	return resputil.OK(), ctrl.svc.role.With(ctx).Merge(r.RoleID, r.Destination)
}

func (ctrl *Role) Move(ctx context.Context, r *request.RoleMove) (interface{}, error) {
	return resputil.OK(), ctrl.svc.role.With(ctx).Move(r.RoleID, r.OrganisationID)
}

func (ctrl *Role) MemberList(ctx context.Context, r *request.RoleMemberList) (interface{}, error) {
	if mm, err := ctrl.svc.role.With(ctx).MemberList(r.RoleID); err != nil {
		return nil, err
	} else {
		rval := make([]string, len(mm))
		for i := range mm {
			rval[i] = payload.Uint64toa(mm[i].UserID)
		}
		return rval, nil
	}
}

func (ctrl *Role) MemberAdd(ctx context.Context, r *request.RoleMemberAdd) (interface{}, error) {
	return resputil.OK(), ctrl.svc.role.With(ctx).MemberAdd(r.RoleID, r.UserID)
}

func (ctrl *Role) MemberRemove(ctx context.Context, r *request.RoleMemberRemove) (interface{}, error) {
	return resputil.OK(), ctrl.svc.role.With(ctx).MemberRemove(r.RoleID, r.UserID)
}
