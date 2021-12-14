package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/pkg/errors"
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
		CanManageMembersOnRole(context.Context, *types.Role) bool
	}

	rolePayload struct {
		*types.Role

		IsSystem bool `json:"isSystem"`
		IsBypass bool `json:"isBypass"`
		IsClosed bool `json:"isClosed"`

		CanGrant               bool `json:"canGrant"`
		CanUpdateRole          bool `json:"canUpdateRole"`
		CanDeleteRole          bool `json:"canDeleteRole"`
		CanManageMembersOnRole bool `json:"canManageMembersOnRole"`
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
	role, err := ctrl.role.FindByID(ctx, r.RoleID)
	return ctrl.makePayload(ctx, role, err)
}

func (ctrl Role) List(ctx context.Context, r *request.RoleList) (interface{}, error) {
	var (
		err error
		f   = types.RoleFilter{
			Query:    r.Query,
			Labels:   r.Labels,
			MemberID: r.MemberID,

			Archived: filter.State(r.Archived),
			Deleted:  filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.role.Find(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl Role) Create(ctx context.Context, r *request.RoleCreate) (interface{}, error) {
	var (
		err  error
		role = &types.Role{
			Name:   r.Name,
			Handle: r.Handle,
			Labels: r.Labels,
			Meta:   r.Meta,
		}
	)

	role, err = ctrl.role.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	for _, userID := range payload.ParseUint64s(r.Members) {
		err := ctrl.role.MemberAdd(ctx, role.ID, userID)
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
			Labels: r.Labels,
			Meta:   r.Meta,
		}
	)

	role, err = ctrl.role.Update(ctx, role)
	if err != nil {
		return nil, err
	}

	if len(r.Members) > 0 {
		members, err := ctrl.role.MemberList(ctx, r.RoleID)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			err := ctrl.role.MemberRemove(ctx, role.ID, member.UserID)
			if err != nil {
				return nil, err
			}
		}

		for _, userID := range payload.ParseUint64s(r.Members) {
			err := ctrl.role.MemberAdd(ctx, role.ID, userID)
			if err != nil {
				return nil, err
			}
		}
	}

	return ctrl.makePayload(ctx, role, err)
}

func (ctrl Role) Delete(ctx context.Context, r *request.RoleDelete) (interface{}, error) {
	return api.OK(), ctrl.role.Delete(ctx, r.RoleID)
}

func (ctrl Role) Undelete(ctx context.Context, r *request.RoleUndelete) (interface{}, error) {
	return api.OK(), ctrl.role.Undelete(ctx, r.RoleID)
}

func (ctrl Role) Archive(ctx context.Context, r *request.RoleArchive) (interface{}, error) {
	return api.OK(), ctrl.role.Archive(ctx, r.RoleID)
}

func (ctrl Role) Unarchive(ctx context.Context, r *request.RoleUnarchive) (interface{}, error) {
	return api.OK(), ctrl.role.Unarchive(ctx, r.RoleID)
}

// deprecated
func (ctrl Role) Merge(ctx context.Context, r *request.RoleMerge) (interface{}, error) {
	return api.OK(), nil
}

// deprecated
func (ctrl Role) Move(ctx context.Context, r *request.RoleMove) (interface{}, error) {
	return api.OK(), nil
}

func (ctrl Role) MemberList(ctx context.Context, r *request.RoleMemberList) (interface{}, error) {
	if mm, err := ctrl.role.MemberList(ctx, r.RoleID); err != nil {
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
	return api.OK(), ctrl.role.MemberAdd(ctx, r.RoleID, r.UserID)
}

func (ctrl Role) MemberRemove(ctx context.Context, r *request.RoleMemberRemove) (interface{}, error) {
	return api.OK(), ctrl.role.MemberRemove(ctx, r.RoleID, r.UserID)
}

func (ctrl *Role) TriggerScript(ctx context.Context, r *request.RoleTriggerScript) (rsp interface{}, err error) {
	var (
		role *types.Role
	)

	if role, err = ctrl.role.FindByID(ctx, r.RoleID); err != nil {
		return
	}

	// @todo implement same behaviour as we have on record - role+oldRole
	err = corredor.Service().Exec(ctx, r.Script, corredor.ExtendScriptArgs(event.RoleOnManual(role, role), r.Args))
	return role, err
}

func (ctrl Role) makePayload(ctx context.Context, r *types.Role, err error) (*rolePayload, error) {
	if err != nil || r == nil {
		return nil, err
	}

	return &rolePayload{
		Role: r,

		CanGrant:               ctrl.ac.CanGrant(ctx),
		CanUpdateRole:          ctrl.ac.CanUpdateRole(ctx, r),
		CanDeleteRole:          ctrl.ac.CanDeleteRole(ctx, r),
		CanManageMembersOnRole: ctrl.ac.CanManageMembersOnRole(ctx, r),

		IsSystem: ctrl.role.IsSystem(r),
		IsClosed: ctrl.role.IsClosed(r),
		IsBypass: auth.BypassRoles().FindByID(r.ID) != nil,
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
