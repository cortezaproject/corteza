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
	User struct {
		user service.UserService
		role service.RoleService
	}

	userSetPayload struct {
		Filter types.UserFilter `json:"filter"`
		Set    types.UserSet    `json:"set"`
	}
)

func (User) New() *User {
	ctrl := &User{}
	ctrl.user = service.DefaultUser
	ctrl.role = service.DefaultRole
	return ctrl
}

func (ctrl User) List(ctx context.Context, r *request.UserList) (interface{}, error) {
	f := types.UserFilter{
		Query:        r.Query,
		Email:        r.Email,
		Username:     r.Username,
		Kind:         r.Kind,
		IncSuspended: r.IncSuspended,
		IncDeleted:   r.IncDeleted,
		Page:         r.Page,
		PerPage:      r.PerPage,
	}

	set, filter, err := ctrl.user.With(ctx).Find(f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl User) Create(ctx context.Context, r *request.UserCreate) (interface{}, error) {
	user := &types.User{
		Email:  r.Email,
		Name:   r.Name,
		Handle: r.Handle,
		Kind:   r.Kind,
	}

	return ctrl.user.With(ctx).Create(user)
}

func (ctrl User) Update(ctx context.Context, r *request.UserUpdate) (interface{}, error) {
	user := &types.User{
		ID:     r.UserID,
		Email:  r.Email,
		Name:   r.Name,
		Handle: r.Handle,
		Kind:   r.Kind,
	}

	return ctrl.user.With(ctx).Update(user)
}

func (ctrl User) Read(ctx context.Context, r *request.UserRead) (interface{}, error) {
	return ctrl.user.With(ctx).FindByID(r.UserID)
}

func (ctrl User) Delete(ctx context.Context, r *request.UserDelete) (interface{}, error) {
	return resputil.OK(), ctrl.user.With(ctx).Delete(r.UserID)
}

func (ctrl User) Suspend(ctx context.Context, r *request.UserSuspend) (interface{}, error) {
	return resputil.OK(), ctrl.user.With(ctx).Suspend(r.UserID)
}

func (ctrl User) Unsuspend(ctx context.Context, r *request.UserUnsuspend) (interface{}, error) {
	return resputil.OK(), ctrl.user.With(ctx).Unsuspend(r.UserID)
}

func (ctrl User) SetPassword(ctx context.Context, r *request.UserSetPassword) (interface{}, error) {
	return resputil.OK(), ctrl.user.With(ctx).SetPassword(r.UserID, r.Password)
}

func (ctrl User) MembershipList(ctx context.Context, r *request.UserMembershipList) (interface{}, error) {
	if mm, err := ctrl.role.With(ctx).Membership(r.UserID); err != nil {
		return nil, err
	} else {
		rval := make([]string, len(mm))
		for i := range mm {
			rval[i] = payload.Uint64toa(mm[i].RoleID)
		}
		return rval, nil
	}
}

func (ctrl User) MembershipAdd(ctx context.Context, r *request.UserMembershipAdd) (interface{}, error) {
	return resputil.OK(), ctrl.role.With(ctx).MemberAdd(r.RoleID, r.UserID)
}

func (ctrl User) MembershipRemove(ctx context.Context, r *request.UserMembershipRemove) (interface{}, error) {
	return resputil.OK(), ctrl.role.With(ctx).MemberRemove(r.RoleID, r.UserID)
}

func (ctrl User) makeFilterPayload(ctx context.Context, uu types.UserSet, f types.UserFilter, err error) (*userSetPayload, error) {
	if err != nil {
		return nil, err
	}

	return &userSetPayload{Filter: f, Set: uu}, nil
}
