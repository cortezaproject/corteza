package service

import (
	"context"
	"io"
	"net/mail"
	"regexp"
	"strconv"
	"strings"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

const (
	maskPrivateDataEmail = "####.#######@######.###"
	maskPrivateDataName  = "##### ##########"
)

type (
	user struct {
		db  *factory.DB
		ctx context.Context

		actionlog actionlog.Recorder

		settings *types.Settings

		auth         userAuth
		subscription userSubscriptionChecker

		ac       userAccessController
		eventbus eventDispatcher

		user        repository.UserRepository
		role        repository.RoleRepository
		credentials repository.CredentialsRepository
	}

	userAuth interface {
		checkPasswordStrength(string) bool
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

	// Temp types to support user.Preloader
	userIdGetter func(chan uint64)
	userSetter   func(*types.User) error

	UserService interface {
		With(ctx context.Context) UserService

		FindByUsername(username string) (*types.User, error)
		FindByEmail(email string) (*types.User, error)
		FindByHandle(handle string) (*types.User, error)
		FindByID(id uint64) (*types.User, error)
		FindByAny(ctx context.Context, identifier interface{}) (*types.User, error)
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

		Preloader(userIdGetter, types.UserFilter, userSetter) error
	}
)

func User(ctx context.Context) UserService {
	return (&user{
		eventbus: eventbus.Service(),
		ac:       DefaultAccessControl,
		settings: CurrentSettings,
		auth:     DefaultAuth,

		actionlog: DefaultActionlog,

		subscription: CurrentSubscription,
	}).With(ctx)
}

func (svc user) With(ctx context.Context) UserService {
	db := repository.DB(ctx)

	return &user{
		ctx: ctx,
		db:  db,

		actionlog: svc.actionlog,

		ac:           svc.ac,
		eventbus:     svc.eventbus,
		settings:     svc.settings,
		auth:         svc.auth,
		subscription: svc.subscription,

		user:        repository.User(ctx, db),
		role:        repository.Role(ctx, db),
		credentials: repository.Credentials(ctx, db),
	}
}

func (svc user) FindByID(userID uint64) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = svc.db.Transaction(func() error {
		if userID == 0 {
			return UserErrInvalidID()
		}

		su := internalAuth.NewIdentity(userID)
		if internalAuth.IsSuperUser(su) {
			// Handling case when looking for a super-user
			//
			// Currently, superuser is a virtual entity
			// and does not exists in the db
			u = &types.User{ID: userID}
			return nil
		}

		u, err = svc.proc(svc.user.FindByID(userID))
		return err
	})

	return u, svc.recordAction(svc.ctx, uaProps, UserActionLookup, err)
}

func (svc user) FindByEmail(email string) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{user: &types.User{Email: email}}
	)

	err = svc.db.Transaction(func() error {
		u, err = svc.proc(svc.user.FindByEmail(email))
		return err
	})

	return u, svc.recordAction(svc.ctx, uaProps, UserActionLookup, err)
}

func (svc user) FindByUsername(username string) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{user: &types.User{Username: username}}
	)

	err = svc.db.Transaction(func() error {
		u, err = svc.proc(svc.user.FindByUsername(username))
		return err
	})

	return u, svc.recordAction(svc.ctx, uaProps, UserActionLookup, err)
}

func (svc user) FindByHandle(handle string) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{user: &types.User{Handle: handle}}
	)

	err = svc.db.Transaction(func() error {
		u, err = svc.proc(svc.user.FindByHandle(handle))
		return err
	})

	return u, svc.recordAction(svc.ctx, uaProps, UserActionLookup, err)
}

// FindByAny finds user by given identifier (context, id, handle, email)
//
// This function goes against the context anti (!!!) pattern we're using
// (and trying to get rid of)
//
// Main reason to push ctx here as the 1st arg is allow (simple) interface definition
// in the consumers that reside under the pkg/
func (svc user) FindByAny(ctx context.Context, identifier interface{}) (u *types.User, err error) {
	if ctx, ok := identifier.(context.Context); ok {
		identifier = internalAuth.GetIdentityFromContext(ctx).Identity()
	}

	if ID, ok := identifier.(uint64); ok {
		u, err = svc.With(ctx).FindByID(ID)
	} else if identity, ok := identifier.(internalAuth.Identifiable); ok {
		u, err = svc.With(ctx).FindByID(identity.Identity())
	} else if strIdentifier, ok := identifier.(string); ok {
		if ID, _ := strconv.ParseUint(strIdentifier, 10, 64); ID > 0 {
			u, err = svc.With(ctx).FindByID(ID)
		} else if strings.Contains(strIdentifier, "@") {
			u, err = svc.With(ctx).FindByEmail(strIdentifier)
		} else {
			u, err = svc.With(ctx).FindByHandle(strIdentifier)
		}
	} else {
		err = UserErrInvalidID()
	}

	if err != nil {
		return
	}

	rr, _, err := svc.role.With(ctx, svc.db).Find(types.RoleFilter{MemberID: u.ID})
	if err != nil {
		return nil, err
	}

	u.SetRoles(rr.IDs())
	return
}

func (svc user) proc(u *types.User, err error) (*types.User, error) {
	if err != nil {
		if repository.ErrUserNotFound.Eq(err) {
			return nil, UserErrNotFound()
		}

		return nil, err
	}

	svc.handlePrivateData(u)

	return u, nil
}

func (svc user) Find(filter types.UserFilter) (uu types.UserSet, f types.UserFilter, err error) {
	var (
		uaProps = &userActionProps{filter: &filter}
	)

	err = svc.db.Transaction(func() error {
		if filter.Deleted > 0 {
			// If list with deleted users is requested
			// user must have access permissions to system (ie: is admin)
			//
			// not the best solution but ATM it allows us to have at least
			// some kind of control over who can see deleted users
			if !svc.ac.CanAccess(svc.ctx) {
				return UserErrNotAllowedToListUsers()
			}
		}

		// Prepare filter for email unmasking check
		filter.IsEmailUnmaskable = svc.ac.FilterUsersWithUnmaskableEmail(svc.ctx)

		// Prepare filter for name unmasking check
		filter.IsNameUnmaskable = svc.ac.FilterUsersWithUnmaskableName(svc.ctx)

		filter.IsReadable = svc.ac.FilterReadableUsers(svc.ctx)

		uu, f, err = svc.user.Find(filter)
		if err != nil {
			return err
		}

		return uu.Walk(func(u *types.User) error {
			svc.handlePrivateData(u)
			return nil
		})
	})

	return uu, f, svc.recordAction(svc.ctx, uaProps, UserActionSearch, err)
}

func (svc user) Create(new *types.User) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{new: new}
	)

	err = svc.db.Transaction(func() (err error) {
		if !svc.ac.CanCreateUser(svc.ctx) {
			return UserErrNotAllowedToCreate()
		}

		if !handle.IsValid(new.Handle) {
			return UserErrInvalidHandle()
		}

		if _, err := mail.ParseAddress(new.Email); err != nil {
			return UserErrInvalidEmail()
		}

		if svc.subscription != nil {
			// When we have an active subscription, we need to check
			// if users can be create or did this deployment hit
			// it's user-limit
			err = svc.subscription.CanCreateUser(svc.user.Total())
			if err != nil {
				return err
			}
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.UserBeforeCreate(new, u)); err != nil {
			return
		}

		if new.Handle == "" {
			createHandle(svc.user, new)
		}

		if err = svc.UniqueCheck(new); err != nil {
			return
		}

		if u, err = svc.user.Create(new); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.UserAfterCreate(new, u))
		return
	})

	return u, svc.recordAction(svc.ctx, uaProps, UserActionCreate, err)
}

func (svc user) CreateWithAvatar(input *types.User, avatar io.Reader) (out *types.User, err error) {
	// @todo: avatar
	return svc.Create(input)
}

func (svc user) Update(upd *types.User) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{update: upd}
	)

	err = svc.db.Transaction(func() (err error) {
		if upd.ID == 0 {
			return UserErrInvalidID()
		}

		if !handle.IsValid(upd.Handle) {
			return UserErrInvalidHandle()
		}

		if _, err := mail.ParseAddress(upd.Email); err != nil {
			return UserErrInvalidEmail()
		}

		if u, err = svc.user.FindByID(upd.ID); err != nil {
			return
		}

		uaProps.setUser(u)

		if upd.ID != internalAuth.GetIdentityFromContext(svc.ctx).Identity() {
			if !svc.ac.CanUpdateUser(svc.ctx, u) {
				return UserErrNotAllowedToUpdate()
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

		if err = svc.UniqueCheck(u); err != nil {
			return
		}

		if u, err = svc.user.Update(u); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.UserAfterUpdate(upd, u))
		return
	})

	return u, svc.recordAction(svc.ctx, uaProps, UserActionUpdate, err)
}

// UniqueCheck verifies user's email, username and handle
func (svc user) UniqueCheck(u *types.User) (err error) {
	isUnique := func(field string) bool {
		f := types.UserFilter{
			// If user exists and is deleted -- not a dup
			Deleted: rh.FilterStateExcluded,

			// If user exists and is suspended -- duplicate
			Suspended: rh.FilterStateInclusive,
		}

		switch field {
		case "email":
			if u.Email == "" {
				return true
			}

			f.Email = u.Email

		case "username":
			if u.Username == "" {
				return true
			}

			f.Username = u.Username
		case "handle":
			if u.Handle == "" {
				return true
			}

			f.Handle = u.Handle
		}

		set, _, err := svc.user.Find(f)
		if err != nil || len(set) > 1 {
			// In case of error or multiple users returned
			return false
		}

		return len(set) == 0 || set[0].ID == u.ID
	}

	if !isUnique("email") {
		return UserErrEmailNotUnique()
	}

	if !isUnique("username") {
		return UserErrUsernameNotUnique()
	}

	if !isUnique("handle") {
		return UserErrHandleNotUnique()
	}

	return nil
}

func (svc user) UpdateWithAvatar(mod *types.User, avatar io.Reader) (out *types.User, err error) {
	// @todo: avatar
	return svc.Create(mod)
}

func (svc user) Delete(userID uint64) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = svc.db.Transaction(func() (err error) {
		if userID == 0 {
			return UserErrInvalidID()
		}

		if u, err = svc.user.FindByID(userID); err != nil {
			return
		}

		if !svc.ac.CanDeleteUser(svc.ctx, u) {
			return UserErrNotAllowedToDelete()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.UserBeforeUpdate(nil, u)); err != nil {
			return
		}

		if err = svc.user.DeleteByID(userID); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.UserAfterDelete(nil, u))
		return nil
	})

	return svc.recordAction(svc.ctx, uaProps, UserActionDelete, err)
}

func (svc user) Undelete(userID uint64) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = svc.db.Transaction(func() (err error) {
		if userID == 0 {
			return UserErrInvalidID()
		}

		if u, err = svc.user.FindByID(userID); err != nil {
			return
		}

		uaProps.setUser(u)

		if err = svc.UniqueCheck(u); err != nil {
			return err
		}

		if !svc.ac.CanDeleteUser(svc.ctx, u) {
			return UserErrNotAllowedToDelete()
		}

		if err = svc.user.UndeleteByID(userID); err != nil {
			return err
		}

		return nil
	})

	return svc.recordAction(svc.ctx, uaProps, UserActionUndelete, err)

}

func (svc user) Suspend(userID uint64) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = svc.db.Transaction(func() (err error) {
		if userID == 0 {
			return UserErrInvalidID()
		}

		if u, err = svc.user.FindByID(userID); err != nil {
			return
		}

		uaProps.setUser(u)

		if !svc.ac.CanSuspendUser(svc.ctx, u) {
			return UserErrNotAllowedToSuspend()
		}

		if err = svc.user.SuspendByID(userID); err != nil {
			return err
		}

		return nil
	})

	return svc.recordAction(svc.ctx, uaProps, UserActionSuspend, err)

}

func (svc user) Unsuspend(userID uint64) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = svc.db.Transaction(func() (err error) {
		if userID == 0 {
			return UserErrInvalidID()
		}

		if u, err = svc.user.FindByID(userID); err != nil {
			return
		}

		uaProps.setUser(u)

		if !svc.ac.CanUnsuspendUser(svc.ctx, u) {
			return UserErrNotAllowedToUnsuspend()
		}

		if err = svc.user.UnsuspendByID(userID); err != nil {
			return err
		}

		return nil
	})

	return svc.recordAction(svc.ctx, uaProps, UserActionUnsuspend, err)

}

// SetPassword sets new password for a user
//
// Expecting setter to have permissions to update modify users and internal authentication enabled
func (svc user) SetPassword(userID uint64, newPassword string) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = svc.db.Transaction(func() error {
		if u, err = svc.user.FindByID(userID); err != nil {
			return err
		}

		uaProps.setUser(u)

		if !svc.ac.CanUpdateUser(svc.ctx, u) {
			return UserErrNotAllowedToUpdate()
		}

		if !svc.auth.checkPasswordStrength(newPassword) {
			return UserErrPasswordNotSecure()
		}

		if err := svc.auth.changePassword(userID, newPassword); err != nil {
			return err
		}

		return nil
	})

	return svc.recordAction(svc.ctx, uaProps, UserActionSetPassword, err)

}

// Masks (or leaves as-is) private data on user
func (svc user) handlePrivateData(u *types.User) {
	if !svc.ac.CanUnmaskEmail(svc.ctx, u) {
		u.Email = maskPrivateDataEmail
	}

	if !svc.ac.CanUnmaskName(svc.ctx, u) {
		u.Name = maskPrivateDataName
	}
}

// Preloader collects all ids of users, loads them and sets them back
//
//
// @todo this kind of preloader is useful and can be implemented in bunch
//       of places and replace old code
func (svc user) Preloader(g userIdGetter, f types.UserFilter, s userSetter) error {
	var (
		// channel that will collect the IDs in the getter
		ch = make(chan uint64, 0)

		// unique index for IDs
		unq = make(map[uint64]bool)
	)

	// Reset the collection of the IDs
	f.UserID = make([]uint64, 0)

	// Call getter and collect the IDs
	go g(ch)

rangeLoop:
	for {
		select {
		case <-svc.ctx.Done():
			close(ch)
			break rangeLoop
		case id, ok := <-ch:
			if !ok {
				// Channel closed
				break rangeLoop
			}

			if !unq[id] {
				unq[id] = true
				f.UserID = append(f.UserID, id)
			}
		}

	}

	// Load all users (even if deleted, suspended) from the given list of IDs
	uu, _, err := svc.Find(f)

	if err != nil {
		return err
	}

	return uu.Walk(s)
}

func createHandle(r repository.UserRepository, u *types.User) {
	if u.Handle == "" {
		u.Handle, _ = handle.Cast(
			// Must not exist before
			func(s string) bool {
				e, err := r.FindByHandle(s)
				return err == repository.ErrUserNotFound && (e == nil || e.ID == u.ID)
			},
			// use name or username
			u.Name,
			u.Username,
			// use email w/o domain
			regexp.
				MustCompile("(@.*)$").
				ReplaceAllString(u.Email, ""),
			//
		)
	}
}
