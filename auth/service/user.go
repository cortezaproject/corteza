package service

import (
	"context"
	"github.com/crusttech/crust/auth/repository"
	"github.com/crusttech/crust/auth/types"
)

const (
	ErrUserInvalidCredentials = serviceError("UserInvalidCredentials")
	ErrUserLocked             = serviceError("UserLocked")
)

type (
	user struct {
		repository repository.User
	}

	UserService interface {
		With(ctx context.Context) UserService

		Create(input *types.User) (*types.User, error)
		ValidateCredentials(username, password string) (*types.User, error)
	}
)

func User() UserService {
	return &user{
		repository.NewUser(context.Background()),
	}
}

func (svc *user) With(ctx context.Context) UserService {
	return &user{
		svc.repository.With(ctx),
	}
}

func (svc *user) ValidateCredentials(username, password string) (*types.User, error) {
	user, err := svc.repository.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if !user.ValidatePassword(password) {
		return nil, ErrUserInvalidCredentials
	}

	if !svc.canLogin(user) {
		return nil, ErrUserLocked
	}

	return user, nil
}

func (svc *user) FindByID(id uint64) (*types.User, error) {
	return svc.repository.FindUserByID(id)
}

func (svc *user) Find(filter *types.UserFilter) ([]*types.User, error) {
	return svc.repository.FindUsers(filter)
}

func (svc *user) Create(input *types.User) (user *types.User, err error) {
	return user, svc.repository.DB().Transaction(func() error {
		// Encrypt user password
		if user, err = svc.repository.CreateUser(input); err != nil {
			return err
		}
		return nil
	})
}

func (svc *user) Update(mod *types.User) (*types.User, error) {
	return svc.repository.UpdateUser(mod)
}

func (svc *user) canLogin(u *types.User) bool {
	return u != nil && u.ID > 0 && u.SuspendedAt == nil && u.DeletedAt == nil
}

func (svc *user) Delete(id uint64) error {
	// @todo: permissions check if current user can delete this user
	// @todo: notify users that user has been removed (remove from web UI)
	return svc.repository.DeleteUserByID(id)
}

func (svc *user) Suspend(id uint64) error {
	// @todo: permissions check if current user can suspend this user
	// @todo: notify users that user has been supsended (remove from web UI)
	return svc.repository.SuspendUserByID(id)
}

func (svc *user) Unsuspend(id uint64) error {
	// @todo: permissions check if current user can unsuspend this user
	// @todo: notify users that user has been unsuspended
	return svc.repository.UnsuspendUserByID(id)
}
