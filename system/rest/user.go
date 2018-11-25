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
	return ctrl.user.With(ctx).Find(&types.UserFilter{Query: r.Query})
}

func (ctrl *User) Create(ctx context.Context, r *request.UserCreate) (interface{}, error) {
	user := &types.User{
		Email:          r.Email,
		Username:       r.Username,
		Name:           r.Name,
		Handle:         r.Handle,
		Meta:           r.Meta,
		SatosaID:       r.SatosaID,
		OrganisationID: r.OrganisationID,
	}
	if err := user.GeneratePassword(r.Password); err != nil {
		return nil, err
	}
	return ctrl.user.With(ctx).Create(user)
}

func (ctrl *User) Edit(ctx context.Context, r *request.UserEdit) (interface{}, error) {
	user := &types.User{
		ID:             r.UserID,
		Email:          r.Email,
		Username:       r.Username,
		Name:           r.Name,
		Handle:         r.Handle,
		Meta:           r.Meta,
		SatosaID:       r.SatosaID,
		OrganisationID: r.OrganisationID,
	}
	if err := user.GeneratePassword(r.Password); err != nil {
		return nil, err
	}
	return ctrl.user.With(ctx).Update(user)
}

func (ctrl *User) Read(ctx context.Context, r *request.UserRead) (interface{}, error) {
	return ctrl.user.With(ctx).FindByID(r.UserID)
}

func (ctrl *User) Remove(ctx context.Context, r *request.UserRemove) (interface{}, error) {
	return nil, ctrl.user.With(ctx).Delete(r.UserID)
}

func (ctrl *User) Suspend(ctx context.Context, r *request.UserSuspend) (interface{}, error) {
	return nil, ctrl.user.With(ctx).Suspend(r.UserID)
}

func (ctrl *User) Unsuspend(ctx context.Context, r *request.UserUnsuspend) (interface{}, error) {
	return nil, ctrl.user.With(ctx).Unsuspend(r.UserID)
}
