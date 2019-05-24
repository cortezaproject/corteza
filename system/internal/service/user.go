package service

import (
	"context"
	"io"

	"github.com/titpetric/factory"
	"go.uber.org/zap"

	internalAuth "github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/system/internal/repository"
	"github.com/cortezaproject/corteza-server/system/types"
)

const (
	ErrUserInvalidCredentials = serviceError("UserInvalidCredentials")
	ErrUserLocked             = serviceError("UserLocked")
)

type (
	user struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac   userAccessController
		user repository.UserRepository
	}

	userAccessController interface {
		CanCreateUser(context.Context) bool
		CanUpdateUser(context.Context, *types.User) bool
		CanDeleteUser(context.Context, *types.User) bool
		CanSuspendUser(context.Context, *types.User) bool
		CanUnsuspendUser(context.Context, *types.User) bool
	}

	UserService interface {
		With(ctx context.Context) UserService

		FindByUsername(username string) (*types.User, error)
		FindByEmail(email string) (*types.User, error)
		FindByID(id uint64) (*types.User, error)
		FindByIDs(id ...uint64) (types.UserSet, error)
		Find(filter *types.UserFilter) (types.UserSet, error)

		Create(input *types.User) (*types.User, error)
		Update(mod *types.User) (*types.User, error)

		CreateWithAvatar(input *types.User, avatar io.Reader) (*types.User, error)
		UpdateWithAvatar(mod *types.User, avatar io.Reader) (*types.User, error)

		Delete(id uint64) error
		Suspend(id uint64) error
		Unsuspend(id uint64) error

		// ValidateCredentials(username, password string) (*types.User, error)
	}
)

func User(ctx context.Context) UserService {
	return (&user{
		logger: DefaultLogger.Named("user"),
		ac:     DefaultAccessControl,
	}).With(ctx)
}

func (svc user) With(ctx context.Context) UserService {
	db := repository.DB(ctx)

	return &user{
		ctx:    ctx,
		db:     db,
		ac:     svc.ac,
		logger: svc.logger,

		user: repository.User(ctx, db),
	}
}

// log() returns zap's logger with requestID from current context and fields.
// func (svc user) log(fields ...zapcore.Field) *zap.Logger {
// 	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
// }

func (svc user) FindByID(ID uint64) (*types.User, error) {
	if ID == 0 {
		return nil, ErrInvalidID
	}

	return svc.user.FindByID(ID)
}

func (svc user) FindByIDs(userIDs ...uint64) (types.UserSet, error) {
	return svc.user.FindByIDs(userIDs...)
}

func (svc user) FindByEmail(email string) (*types.User, error) {
	return svc.user.FindByEmail(email)
}

func (svc user) FindByUsername(username string) (*types.User, error) {
	return svc.user.FindByUsername(username)
}

func (svc user) Find(filter *types.UserFilter) (types.UserSet, error) {
	return svc.user.Find(filter)
}

func (svc user) Create(input *types.User) (out *types.User, err error) {
	if !svc.ac.CanCreateUser(svc.ctx) {
		return nil, ErrNoPermissions.withStack()
	}

	return out, svc.db.Transaction(func() (err error) {
		out, err = svc.user.Create(input)
		return
	})
}

func (svc user) CreateWithAvatar(input *types.User, avatar io.Reader) (out *types.User, err error) {
	// @todo: avatar
	return svc.Create(input)
}

func (svc user) Update(mod *types.User) (u *types.User, err error) {
	if mod.ID == 0 {
		return nil, ErrInvalidID
	}

	if u, err = svc.user.FindByID(mod.ID); err != nil {
		return
	}

	if mod.ID != internalAuth.GetIdentityFromContext(svc.ctx).Identity() && !svc.ac.CanUpdateUser(svc.ctx, u) {
		return nil, ErrNoPermissions.withStack()
	}

	// Assign changed values
	u.Email = mod.Email
	u.Username = mod.Username
	u.Name = mod.Name
	u.Handle = mod.Handle
	u.Kind = mod.Kind

	return u, svc.db.Transaction(func() (err error) {
		u, err = svc.user.Update(u)
		return
	})
}

func (svc user) UpdateWithAvatar(mod *types.User, avatar io.Reader) (out *types.User, err error) {
	// @todo: avatar
	return svc.Create(mod)
}

func (svc user) Delete(ID uint64) (err error) {
	if ID == 0 {
		return ErrInvalidID
	}

	var u *types.User
	if u, err = svc.user.FindByID(ID); err != nil {
		return
	}

	if !svc.ac.CanDeleteUser(svc.ctx, u) {
		return ErrNoPermissions.withStack()
	}

	return svc.db.Transaction(func() (err error) {
		return svc.user.DeleteByID(ID)
	})
}

func (svc user) Suspend(ID uint64) (err error) {
	if ID == 0 {
		return ErrInvalidID
	}

	var u *types.User
	if u, err = svc.user.FindByID(ID); err != nil {
		return
	}

	if !svc.ac.CanSuspendUser(svc.ctx, u) {
		return ErrNoPermissions.withStack()
	}

	return svc.db.Transaction(func() (err error) {
		return svc.user.SuspendByID(ID)
	})
}

func (svc user) Unsuspend(ID uint64) (err error) {
	if ID == 0 {
		return ErrInvalidID
	}

	var u *types.User
	if u, err = svc.user.FindByID(ID); err != nil {
		return
	}

	if !svc.ac.CanUnsuspendUser(svc.ctx, u) {
		return ErrNoPermissions.withStack()
	}

	return svc.db.Transaction(func() (err error) {
		return svc.user.UnsuspendByID(ID)
	})
}
