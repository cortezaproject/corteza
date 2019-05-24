package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = errors.Wrap

type (
	User struct {
		user service.UserService
	}
)

func (User) New() *User {
	ctrl := &User{}
	ctrl.user = service.DefaultUser
	return ctrl
}

// Searches the users table in the database to find users by matching (by-prefix) their.Username
func (ctrl *User) List(ctx context.Context, r *request.UserList) (interface{}, error) {
	return ctrl.user.With(ctx).Find(&types.UserFilter{
		Query:    r.Query,
		Email:    r.Email,
		Username: r.Username,
	})
}

func (ctrl *User) Create(ctx context.Context, r *request.UserCreate) (interface{}, error) {
	user := &types.User{
		Email:  r.Email,
		Name:   r.Name,
		Handle: r.Handle,
		Kind:   types.UserKind(r.Kind),
	}

	return ctrl.user.With(ctx).Create(user)
}

func (ctrl *User) Update(ctx context.Context, r *request.UserUpdate) (interface{}, error) {
	user := &types.User{
		ID:     r.UserID,
		Email:  r.Email,
		Name:   r.Name,
		Handle: r.Handle,
		Kind:   types.UserKind(r.Kind),
	}

	return ctrl.user.With(ctx).Update(user)
}

func (ctrl *User) Read(ctx context.Context, r *request.UserRead) (interface{}, error) {
	return ctrl.user.With(ctx).FindByID(r.UserID)
}

func (ctrl *User) Delete(ctx context.Context, r *request.UserDelete) (interface{}, error) {
	return nil, ctrl.user.With(ctx).Delete(r.UserID)
}

func (ctrl *User) Suspend(ctx context.Context, r *request.UserSuspend) (interface{}, error) {
	return nil, ctrl.user.With(ctx).Suspend(r.UserID)
}

func (ctrl *User) Unsuspend(ctx context.Context, r *request.UserUnsuspend) (interface{}, error) {
	return nil, ctrl.user.With(ctx).Unsuspend(r.UserID)
}
