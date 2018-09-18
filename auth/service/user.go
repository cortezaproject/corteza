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

		FindByID(id uint64) (*types.User, error)
		Find(filter *types.UserFilter) ([]*types.User, error)

		Create(input *types.User) (*types.User, error)
		Update(mod *types.User) (*types.User, error)

		FindOrCreate(*types.User) (*types.User, error)
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

// Finds if user with a specific email exists and returns it otherwise it creates a fresh one
func (svc *user) FindOrCreate(user *types.User) (out *types.User, err error) {
	//return out, svc.repository.DB().Transaction(func() error {
	if out, err = svc.repository.FindUserByEmail(user.Email); err != repository.ErrUserNotFound {
		return out, err
	} else if out, err = svc.repository.CreateUser(user); err != nil {
		return out, err
	}

	return out, nil
	//})
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
