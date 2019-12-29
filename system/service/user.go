package service

import (
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

const (
	ErrUserInvalidCredentials = serviceError("UserInvalidCredentials")
	ErrUserHandleNotUnique    = serviceError("UserHandleNotUnique")
	ErrUserUsernameNotUnique  = serviceError("UserUsernameNotUnique")
	ErrUserEmailNotUnique     = serviceError("UserEmailNotUnique")
	ErrUserLocked             = serviceError("UserLocked")

	maskPrivateDataEmail = "####.#######@######.###"
	maskPrivateDataName  = "##### ##########"
)

type (
	user struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		settings *types.Settings

		auth         userAuth
		subscription userSubscriptionChecker

		ac       userAccessController
		eventbus eventDispatcher

		user        repository.UserRepository
		credentials repository.CredentialsRepository

		// @todo wire this with settings (privacy.mask.email)
		privacyMaskEmail bool
		// @todo wire this with settings (privacy.mask.name)
		privacyMaskName bool
	}

	userAuth interface {
		checkPasswordStrength(string) error
		changePassword(uint64, string) error
	}

	userSubscriptionChecker interface {
		CanCreateUser(uint) error
	}

	userAccessController interface {
		CanAccess(context.Context) bool
		CanCreateUser(context.Context) bool
		CanUpdateUser(context.Context, *types.User) bool
		CanDeleteUser(context.Context, *types.User) bool
		CanSuspendUser(context.Context, *types.User) bool
		CanUnsuspendUser(context.Context, *types.User) bool
		CanUnmaskEmail(context.Context, *types.User) bool
		CanUnmaskName(context.Context, *types.User) bool

		FilterReadableUsers(ctx context.Context) *permissions.ResourceFilter
		FilterUsersWithUnmaskableEmail(ctx context.Context) *permissions.ResourceFilter
		FilterUsersWithUnmaskableName(ctx context.Context) *permissions.ResourceFilter
	}

	UserService interface {
		With(ctx context.Context) UserService

		FindByUsername(username string) (*types.User, error)
		FindByEmail(email string) (*types.User, error)
		FindByHandle(handle string) (*types.User, error)
		FindByID(id uint64) (*types.User, error)
		FindByAny(any string) (*types.User, error)
		Find(types.UserFilter) (types.UserSet, types.UserFilter, error)

		Create(input *types.User) (*types.User, error)
		Update(mod *types.User) (*types.User, error)

		CreateWithAvatar(input *types.User, avatar io.Reader) (*types.User, error)
		UpdateWithAvatar(mod *types.User, avatar io.Reader) (*types.User, error)

		Delete(id uint64) error
		Suspend(id uint64) error
		Unsuspend(id uint64) error
		Undelete(id uint64) error

		SetPassword(userID uint64, password string) error
	}
)

func User(ctx context.Context) UserService {
	return (&user{
		logger:   DefaultLogger.Named("user"),
		eventbus: eventbus.Service(),
		ac:       DefaultAccessControl,
		settings: CurrentSettings,
		auth:     DefaultAuth,

		subscription: CurrentSubscription,

		// @todo wire this with settings (privacy.mask.email)
		//       new default value will be true!
		privacyMaskEmail: false,

		// @todo wire this with settings (privacy.mask.name)
		//       new default value will be true!
		privacyMaskName: false,
	}).With(ctx)
}

// log() returns zap's logger with requestID from current context and fields.
func (svc user) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc user) With(ctx context.Context) UserService {
	db := repository.DB(ctx)

	return &user{
		ctx:    ctx,
		db:     db,
		logger: svc.logger,

		ac:           svc.ac,
		eventbus:     svc.eventbus,
		settings:     svc.settings,
		auth:         svc.auth,
		subscription: svc.subscription,

		user:        repository.User(ctx, db),
		credentials: repository.Credentials(ctx, db),

		privacyMaskEmail: svc.privacyMaskEmail,
		privacyMaskName:  svc.privacyMaskName,
	}
}

func (svc user) FindByID(ID uint64) (*types.User, error) {
	if ID == 0 {
		return nil, ErrInvalidID
	}

	return svc.proc(svc.user.FindByID(ID))
}

func (svc user) FindByEmail(email string) (*types.User, error) {
	return svc.proc(svc.user.FindByEmail(email))
}

func (svc user) FindByUsername(username string) (*types.User, error) {
	return svc.proc(svc.user.FindByUsername(username))
}

func (svc user) FindByHandle(handle string) (*types.User, error) {
	return svc.proc(svc.user.FindByHandle(handle))
}

func (svc user) FindByAny(any string) (*types.User, error) {
	return svc.proc(func() (*types.User, error) {
		if id, _ := strconv.ParseUint(any, 10, 64); id > 0 {
			return svc.user.FindByID(id)
		}

		if strings.Contains(any, "@") {
			return svc.user.FindByEmail(any)
		}

		return svc.user.FindByHandle(any)
	}())
}

func (svc user) proc(u *types.User, err error) (*types.User, error) {
	if err != nil {
		return nil, err
	}

	svc.handlePrivateData(u)

	return u, nil
}

func (svc user) Find(f types.UserFilter) (types.UserSet, types.UserFilter, error) {
	if f.Deleted > 0 {
		// If list with deleted users is requested
		// user must have access permissions to system (ie: is admin)
		//
		// not the best solution but ATM it allows us to have at least
		// some kind of control over who can see deleted users
		if !svc.ac.CanAccess(svc.ctx) {
			return nil, f, ErrNoPermissions.withStack()
		}
	}

	if svc.privacyMaskEmail {
		// Prepare filter for email unmasking check
		f.IsEmailUnmaskable = svc.ac.FilterUsersWithUnmaskableEmail(svc.ctx)

	}

	if svc.privacyMaskName {
		// Prepare filter for name unmasking check
		f.IsNameUnmaskable = svc.ac.FilterUsersWithUnmaskableName(svc.ctx)
	}

	f.IsReadable = svc.ac.FilterReadableUsers(svc.ctx)

	return svc.procSet(svc.user.Find(f))
}

func (svc user) procSet(u types.UserSet, f types.UserFilter, err error) (types.UserSet, types.UserFilter, error) {
	if err != nil {
		return nil, f, err
	}

	_ = u.Walk(func(u *types.User) error {
		svc.handlePrivateData(u)
		return nil
	})

	return u, f, nil
}

func (svc user) Create(new *types.User) (u *types.User, err error) {
	if !svc.ac.CanCreateUser(svc.ctx) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	if svc.subscription != nil {
		// When we have an active subscription, we need to check
		// if users can be creare or did this deployment hit
		// it's user-limit
		err = svc.subscription.CanCreateUser(svc.user.Total())
		if err != nil {
			return nil, err
		}
	}

	if err = svc.eventbus.WaitFor(svc.ctx, event.UserBeforeCreate(new, u)); err != nil {
		return
	}

	return u, svc.db.Transaction(func() (err error) {

		if err = svc.UniqueCheck(new); err != nil {
			return
		}

		if u, err = svc.user.Create(new); err != nil {
			return
		}

		defer svc.eventbus.Dispatch(svc.ctx, event.UserAfterCreate(new, u))
		return
	})
}

func (svc user) CreateWithAvatar(input *types.User, avatar io.Reader) (out *types.User, err error) {
	// @todo: avatar
	return svc.Create(input)
}

func (svc user) Update(upd *types.User) (u *types.User, err error) {
	if upd.ID == 0 {
		return nil, ErrInvalidID
	}

	if u, err = svc.user.FindByID(upd.ID); err != nil {
		return
	}

	if upd.ID != internalAuth.GetIdentityFromContext(svc.ctx).Identity() {
		if !svc.ac.CanUpdateUser(svc.ctx, u) {
			return nil, ErrNoUpdatePermissions.withStack()
		}
	}

	// Assign changed values
	u.Email = upd.Email
	u.Username = upd.Username
	u.Name = upd.Name
	u.Handle = upd.Handle
	u.Kind = upd.Kind

	if err = svc.eventbus.WaitFor(svc.ctx, event.UserBeforeUpdate(upd, u)); err != nil {
		return
	}

	return u, svc.db.Transaction(func() (err error) {
		if err = svc.UniqueCheck(u); err != nil {
			return
		}

		if u, err = svc.user.Update(u); err != nil {
			return
		}

		defer svc.eventbus.Dispatch(svc.ctx, event.UserAfterUpdate(upd, u))
		return
	})
}

func (svc user) UniqueCheck(u *types.User) (err error) {
	if u.Email != "" {
		if ex, _ := svc.user.FindByEmail(u.Email); ex != nil && ex.ID > 0 && ex.ID != u.ID {
			return ErrUserEmailNotUnique
		}
	}

	if u.Username != "" {
		if ex, _ := svc.user.FindByUsername(u.Username); ex != nil && ex.ID > 0 && ex.ID != u.ID {
			return ErrUserUsernameNotUnique
		}
	}

	if u.Handle != "" {
		if ex, _ := svc.user.FindByHandle(u.Handle); ex != nil && ex.ID > 0 && ex.ID != u.ID {
			return ErrUserHandleNotUnique
		}
	}

	return nil
}

func (svc user) UpdateWithAvatar(mod *types.User, avatar io.Reader) (out *types.User, err error) {
	// @todo: avatar
	return svc.Create(mod)
}

func (svc user) Delete(ID uint64) (err error) {
	var (
		del *types.User
	)

	if ID == 0 {
		return ErrInvalidID
	}

	if del, err = svc.user.FindByID(ID); err != nil {
		return
	}

	if !svc.ac.CanDeleteUser(svc.ctx, del) {
		return ErrNoPermissions.withStack()
	}

	if err = svc.eventbus.WaitFor(svc.ctx, event.UserBeforeUpdate(nil, del)); err != nil {
		return
	}

	if err = svc.user.DeleteByID(ID); err != nil {
		return
	}

	defer svc.eventbus.Dispatch(svc.ctx, event.UserAfterDelete(nil, del))
	return
}

func (svc user) Undelete(ID uint64) (err error) {
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
		return svc.user.UndeleteByID(ID)
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

// SetPassword sets new password for a user
//
// Expecting setter to have permissions to update modify users and internal authentication enabled
func (svc user) SetPassword(userID uint64, newPassword string) (err error) {
	log := svc.log(svc.ctx, zap.Uint64("userID", userID))

	if !svc.settings.Auth.Internal.Enabled {
		return errors.New("internal authentication disabled")
	}

	var u *types.User
	if u, err = svc.user.FindByID(userID); err != nil {
		return
	}

	if !svc.ac.CanUpdateUser(svc.ctx, u) {
		return ErrNoPermissions.withStack()
	}

	if err = svc.auth.checkPasswordStrength(newPassword); err != nil {
		return
	}

	return svc.db.Transaction(func() error {
		if err := svc.auth.changePassword(userID, newPassword); err != nil {
			return err
		}

		log.Info("password changed")

		return nil
	})
}

// Masks (or leaves as-is) private data on user
func (svc user) handlePrivateData(u *types.User) {
	if svc.privacyMaskEmail && !svc.ac.CanUnmaskEmail(svc.ctx, u) {
		u.Email = maskPrivateDataEmail
	}

	if svc.privacyMaskName && !svc.ac.CanUnmaskEmail(svc.ctx, u) {
		u.Name = maskPrivateDataName
	}
}
