package service

import (
	"context"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
	"golang.org/x/crypto/bcrypt"
)

const (
	ErrUserInvalidCredentials = serviceError("UserInvalidCredentials")
	ErrUserLocked             = serviceError("UserLocked")
)

type (
	user struct {
		repository userRepository
	}

	userRepository interface {
		FindByUsername(ctx context.Context, username string) (*types.User, error)
		FindByID(ctx context.Context, userID uint64) (*types.User, error)
		Find(ctx context.Context, filter *types.UserFilter) ([]*types.User, error)

		Create(ctx context.Context, user *types.User) (*types.User, error)
		Update(ctx context.Context, user *types.User) (*types.User, error)

		deleter
		suspender
	}
)

func User() *user {
	return &user{repository: repository.User()}
}

func (svc user) ValidateCredentials(ctx context.Context, username, password string) (*types.User, error) {
	user, err := svc.repository.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if !svc.validatePassword(user, password) {
		return nil, ErrUserInvalidCredentials
	}

	if !svc.canLogin(user) {
		return nil, ErrUserLocked
	}

	return user, nil
}

func (svc user) FindByID(ctx context.Context, id uint64) (*types.User, error) {
	return svc.repository.FindByID(ctx, id)
}

func (svc user) Find(ctx context.Context, filter *types.UserFilter) ([]*types.User, error) {
	return svc.repository.Find(ctx, filter)
}

func (svc user) Create(ctx context.Context, mod *types.User) (*types.User, error) {
	return svc.repository.Create(ctx, mod)
}

func (svc user) Update(ctx context.Context, mod *types.User) (*types.User, error) {
	return svc.repository.Update(ctx, mod)
}

func (svc user) validatePassword(user *types.User, password string) bool {
	return user != nil &&
		bcrypt.CompareHashAndPassword(user.Password, []byte(password)) == nil
}

func (svc user) generatePassword(user *types.User, password string) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.SetPassword(pwd)
	return nil
}

func (svc user) canLogin(u *types.User) bool {
	return u != nil && u.ID > 0 && u.SuspendedAt == nil && u.DeletedAt == nil
}

func (svc user) Delete(ctx context.Context, id uint64) error {
	// @todo: permissions check if current user can delete this user
	// @todo: notify users that user has been removed (remove from web UI)
	return svc.repository.Delete(ctx, id)
}

func (svc user) Suspend(ctx context.Context, id uint64) error {
	// @todo: permissions check if current user can suspend this user
	// @todo: notify users that user has been supsended (remove from web UI)
	return svc.repository.Suspend(ctx, id)
}

func (svc user) Unsuspend(ctx context.Context, id uint64) error {
	// @todo: permissions check if current user can unsuspend this user
	// @todo: notify users that user has been unsuspended
	return svc.repository.Unsuspend(ctx, id)
}
