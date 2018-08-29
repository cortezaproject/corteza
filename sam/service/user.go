package service

import (
	"context"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

const (
	ErrUserInvalidCredentials = serviceError("UserInvalidCredentials")
	ErrUserLocked             = serviceError("UserLocked")
)

type (
	user struct {
		rpo userRepository
	}

	UserService interface {
		Find(ctx context.Context, filter *types.UserFilter) ([]*types.User, error)
	}

	userRepository interface {
		repository.Transactionable
		repository.Contextable
		repository.User
	}
)

func User() *user {
	return &user{rpo: repository.New()}
}

func (svc user) ValidateCredentials(ctx context.Context, username, password string) (*types.User, error) {
	user, err := svc.rpo.FindUserByUsername(username)
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

func (svc user) FindByID(ctx context.Context, id uint64) (*types.User, error) {
	return svc.rpo.WithCtx(ctx).FindUserByID(id)
}

func (svc user) Find(ctx context.Context, filter *types.UserFilter) ([]*types.User, error) {
	return svc.rpo.FindUsers(filter)
}

func (svc user) Create(ctx context.Context, input *types.User) (user *types.User, err error) {
	return user, svc.rpo.BeginWith(ctx, func(r repository.Interfaces) error {
		// Encrypt user password
		if user, err = r.CreateUser(input); err != nil {
			return err
		}

		return nil
	})
}

func (svc user) Update(ctx context.Context, mod *types.User) (*types.User, error) {
	return svc.rpo.UpdateUser(mod)
}

func (svc user) canLogin(u *types.User) bool {
	return u != nil && u.ID > 0 && u.SuspendedAt == nil && u.DeletedAt == nil
}

func (svc user) Delete(ctx context.Context, id uint64) error {
	// @todo: permissions check if current user can delete this user
	// @todo: notify users that user has been removed (remove from web UI)
	return svc.rpo.DeleteUserByID(id)
}

func (svc user) Suspend(ctx context.Context, id uint64) error {
	// @todo: permissions check if current user can suspend this user
	// @todo: notify users that user has been supsended (remove from web UI)
	return svc.rpo.SuspendUserByID(id)
}

func (svc user) Unsuspend(ctx context.Context, id uint64) error {
	// @todo: permissions check if current user can unsuspend this user
	// @todo: notify users that user has been unsuspended
	return svc.rpo.UnsuspendUserByID(id)
}

var _ UserService = &user{}
