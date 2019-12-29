package rest

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = errors.Wrap

type (
	Role struct {
		role service.RoleService
		ac   roleAccessController
	}

	roleAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateRole(context.Context, *types.Role) bool
		CanDeleteRole(context.Context, *types.Role) bool
	}

	rolePayload struct {
		*types.Role

		CanGrant      bool `json:"canGrant"`
		CanUpdateRole bool `json:"canUpdateRole"`
		CanDeleteRole bool `json:"canDeleteRole"`
	}

	roleSetPayload struct {
		Filter types.RoleFilter `json:"filter"`
		Set    []*rolePayload   `json:"set"`
	}
)

func (Role) New() *Role {
	return &Role{
		role: service.DefaultRole,
		ac:   service.DefaultAccessControl,
	}
}

func (ctrl Role) Read(ctx context.Context, r *request.RoleRead) (interface{}, error) {
	role, err := ctrl.role.With(ctx).FindByID(r.RoleID)
	return ctrl.makePayload(ctx, role, err)
}

func (ctrl Role) List(ctx context.Context, r *request.RoleList) (interface{}, error) {
	f := types.RoleFilter{
		Query: r.Query,

		Archived: rh.FilterState(r.Archived),
		Deleted:  rh.FilterState(r.Deleted),

		Sort:       rh.NormalizeSortColumns(r.Sort),
		PageFilter: rh.Paging(r.Page, r.PerPage),
	}

	set, filter, err := ctrl.role.With(ctx).Find(f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl Role) Create(ctx context.Context, r *request.RoleCreate) (interface{}, error) {
	var (
		err  error
		role = &types.Role{
			Name:   r.Name,
			Handle: r.Handle,
		}
	)

	role, err = ctrl.role.With(ctx).Create(role)
	if err != nil {
		return nil, err
	}

	for _, userID := range payload.ParseUInt64s(r.Members) {
		err := ctrl.role.With(ctx).MemberAdd(role.ID, userID)
		if err != nil {
			return nil, err
		}
	}
	return ctrl.makePayload(ctx, role, err)
}

func (ctrl Role) Update(ctx context.Context, r *request.RoleUpdate) (interface{}, error) {
	var (
		err  error
		role = &types.Role{
			ID:     r.RoleID,
			Name:   r.Name,
			Handle: r.Handle,
		}
	)

	role, err = ctrl.role.With(ctx).Update(role)
	if err != nil {
		return nil, err
	}

	if len(r.Members) > 0 {
		members, err := ctrl.role.With(ctx).MemberList(r.RoleID)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			err := ctrl.role.With(ctx).MemberRemove(role.ID, member.UserID)
			if err != nil {
				return nil, err
			}
		}

		for _, userID := range payload.ParseUInt64s(r.Members) {
			err := ctrl.role.With(ctx).MemberAdd(role.ID, userID)
			if err != nil {
				return nil, err
			}
		}
	}

	return ctrl.makePayload(ctx, role, err)
}

func (ctrl Role) Delete(ctx context.Context, r *request.RoleDelete) (interface{}, error) {
	return resputil.OK(), ctrl.role.With(ctx).Delete(r.RoleID)
}

func (ctrl Role) Undelete(ctx context.Context, r *request.RoleUndelete) (interface{}, error) {
	return resputil.OK(), ctrl.role.With(ctx).Undelete(r.RoleID)
}

func (ctrl Role) Archive(ctx context.Context, r *request.RoleArchive) (interface{}, error) {
	return resputil.OK(), ctrl.role.With(ctx).Archive(r.RoleID)
}

func (ctrl Role) Unarchive(ctx context.Context, r *request.RoleUnarchive) (interface{}, error) {
	return resputil.OK(), ctrl.role.With(ctx).Unarchive(r.RoleID)
}

func (ctrl Role) Merge(ctx context.Context, r *request.RoleMerge) (interface{}, error) {
	return resputil.OK(), ctrl.role.With(ctx).Merge(r.RoleID, r.Destination)
}

func (ctrl Role) Move(ctx context.Context, r *request.RoleMove) (interface{}, error) {
	return resputil.OK(), ctrl.role.With(ctx).Move(r.RoleID, r.OrganisationID)
}

func (ctrl Role) MemberList(ctx context.Context, r *request.RoleMemberList) (interface{}, error) {
	if mm, err := ctrl.role.With(ctx).MemberList(r.RoleID); err != nil {
		return nil, err
	} else {
		rval := make([]string, len(mm))
		for i := range mm {
			rval[i] = payload.Uint64toa(mm[i].UserID)
		}
		return rval, nil
	}
}

func (ctrl Role) MemberAdd(ctx context.Context, r *request.RoleMemberAdd) (interface{}, error) {
	return resputil.OK(), ctrl.role.With(ctx).MemberAdd(r.RoleID, r.UserID)
}

func (ctrl Role) MemberRemove(ctx context.Context, r *request.RoleMemberRemove) (interface{}, error) {
	return resputil.OK(), ctrl.role.With(ctx).MemberRemove(r.RoleID, r.UserID)
}

func (ctrl *Role) FireTrigger(ctx context.Context, r *request.RoleFireTrigger) (rsp interface{}, err error) {
	var (
		role *types.Role
	)

	if role, err = ctrl.role.With(ctx).FindByID(r.RoleID); err != nil {
		return
	}

	return resputil.OK(), corredor.Service().ExecOnManual(ctx, r.Script, event.RoleOnManual(role, nil))
}

func (ctrl Role) makePayload(ctx context.Context, m *types.Role, err error) (*rolePayload, error) {
	if err != nil || m == nil {
		return nil, err
	}

	return &rolePayload{
		Role: m,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateRole: ctrl.ac.CanUpdateRole(ctx, m),
		CanDeleteRole: ctrl.ac.CanDeleteRole(ctx, m),
	}, nil
}

func (ctrl Role) makeFilterPayload(ctx context.Context, nn types.RoleSet, f types.RoleFilter, err error) (*roleSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &roleSetPayload{Filter: f, Set: make([]*rolePayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
