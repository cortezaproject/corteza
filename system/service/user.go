package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"io"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
)

const (
	maskPrivateDataEmail = "####.#######@######.###"
	maskPrivateDataName  = "##### ##########"
)

type (
	user struct {
		ctx context.Context

		actionlog actionlog.Recorder

		settings *types.AppSettings

		auth         userAuth
		subscription userSubscriptionChecker

		ac       userAccessController
		eventbus eventDispatcher

		store store.Storer
	}

	userAuth interface {
		CheckPasswordStrength(string) bool
		SetPasswordCredentials(context.Context, uint64, string) error
	}

	userSubscriptionChecker interface {
		CanCreateUser(uint) error
	}

	userAccessController interface {
		CanAccess(context.Context) bool
		CanCreateUser(context.Context) bool
		CanReadUser(context.Context, *types.User) bool
		CanUpdateUser(context.Context, *types.User) bool
		CanDeleteUser(context.Context, *types.User) bool
		CanSuspendUser(context.Context, *types.User) bool
		CanUnsuspendUser(context.Context, *types.User) bool
		CanUnmaskEmail(context.Context, *types.User) bool
		CanUnmaskName(context.Context, *types.User) bool
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

		store: DefaultStore,

		actionlog: DefaultActionlog,

		subscription: CurrentSubscription,
	}).With(ctx)
}

func (svc user) With(ctx context.Context) UserService {
	return &user{
		ctx: ctx,

		actionlog: svc.actionlog,

		ac:           svc.ac,
		eventbus:     svc.eventbus,
		settings:     svc.settings,
		auth:         svc.auth,
		subscription: svc.subscription,
		store:        svc.store,
	}
}

func (svc user) FindByID(userID uint64) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() error {
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

		if u, err = svc.proc(store.LookupUserByID(svc.ctx, svc.store, userID)); err != nil {
			return err
		}

		uaProps.setUser(u)

		if err = label.Load(svc.ctx, svc.store, u); err != nil {
			return err
		}

		return nil
	}()

	return u, svc.recordAction(svc.ctx, uaProps, UserActionLookup, err)
}

func (svc user) FindByEmail(email string) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{user: &types.User{Email: email}}
	)

	err = func() error {
		if u, err = svc.proc(store.LookupUserByEmail(svc.ctx, svc.store, email)); err != nil {
			return err
		}

		uaProps.setUser(u)

		if err = label.Load(svc.ctx, svc.store, u); err != nil {
			return err
		}

		return nil
	}()

	return u, svc.recordAction(svc.ctx, uaProps, UserActionLookup, err)
}

func (svc user) FindByUsername(username string) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{user: &types.User{Username: username}}
	)

	err = func() error {
		if u, err = svc.proc(store.LookupUserByUsername(svc.ctx, svc.store, username)); err != nil {
			return err
		}

		uaProps.setUser(u)

		if err = label.Load(svc.ctx, svc.store, u); err != nil {
			return err
		}

		return nil
	}()

	return u, svc.recordAction(svc.ctx, uaProps, UserActionLookup, err)
}

func (svc user) FindByHandle(handle string) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{user: &types.User{Handle: handle}}
	)

	err = func() error {
		if u, err = svc.proc(store.LookupUserByHandle(svc.ctx, svc.store, handle)); err != nil {
			return err
		}

		uaProps.setUser(u)

		if err = label.Load(svc.ctx, svc.store, u); err != nil {
			return err
		}

		return nil
	}()

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

	rr, _, err := store.SearchRoles(svc.ctx, svc.store, types.RoleFilter{MemberID: u.ID})
	if err != nil {
		return nil, err
	}

	u.SetRoles(rr.IDs())
	return
}

func (svc user) proc(u *types.User, err error) (*types.User, error) {
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, UserErrNotFound()
		}

		return nil, err
	}

	svc.handlePrivateData(svc.ctx, u)

	return u, nil
}

// Find interacts with backend storage and
//
// @todo rename to Search() for consistency
func (svc user) Find(filter types.UserFilter) (uu types.UserSet, f types.UserFilter, err error) {
	var (
		uaProps = &userActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.MaskedEmailsEnabled = svc.settings.Privacy.Mask.Email
	filter.MaskedNamesEnabled = svc.settings.Privacy.Mask.Name
	filter.Check = func(res *types.User) (bool, error) {
		if !svc.ac.CanReadUser(svc.ctx, res) {
			return false, nil
		}

		if svc.maskEmail(svc.ctx, res) && ((len(filter.Query) > 0 && strings.HasPrefix(res.Email, filter.Query)) || res.Email == filter.Email) {
			// user email matched but it will be masked later on, so exclude it to prevent data probing
			return false, nil
		}

		if svc.maskName(svc.ctx, res) && (len(filter.Query) > 0 && strings.HasPrefix(res.Name, filter.Query)) {
			// user mail matched but it will be masked later on, so exclude it to prevent data probing
			return false, nil
		}

		return true, nil
	}

	err = func() error {
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

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				svc.ctx,
				svc.store,
				types.User{}.LabelResourceKind(),
				filter.Labels,
			)

			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		uu, f, err = store.SearchUsers(svc.ctx, svc.store, filter)
		if err != nil {
			return err
		}

		if err = label.Load(svc.ctx, svc.store, toLabeledUsers(uu)...); err != nil {
			return err
		}

		return uu.Walk(func(u *types.User) error {
			svc.handlePrivateData(svc.ctx, u)
			return nil
		})
	}()

	return uu, f, svc.recordAction(svc.ctx, uaProps, UserActionSearch, err)
}

func (svc user) Create(new *types.User) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{new: new}
	)

	err = func() (err error) {
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
			var c uint
			if c, err = store.CountUsers(svc.ctx, svc.store, types.UserFilter{}); err != nil {
				return err
			}

			// When we have an active subscription, we need to check
			// if users can be create or did this deployment hit
			// it's user-limit
			err = svc.subscription.CanCreateUser(c)
			if err != nil {
				return err
			}
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.UserBeforeCreate(new, u)); err != nil {
			return
		}

		if new.Handle == "" {
			createUserHandle(svc.ctx, DefaultStore, new)
		}

		if err = uniqueUserCheck(svc.ctx, svc.store, new); err != nil {
			return
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		// consider email confirmed
		// when creating user like this
		new.EmailConfirmed = true

		if err = store.CreateUser(svc.ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(svc.ctx, svc.store, new); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.UserAfterCreate(new, u))
		return
	}()

	return new, svc.recordAction(svc.ctx, uaProps, UserActionCreate, err)
}

func (svc user) CreateWithAvatar(input *types.User, avatar io.Reader) (out *types.User, err error) {
	// @todo: avatar
	return svc.Create(input)
}

func (svc user) Update(upd *types.User) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{update: upd}
	)

	err = func() (err error) {
		if upd.ID == 0 {
			return UserErrInvalidID()
		}

		if !handle.IsValid(upd.Handle) {
			return UserErrInvalidHandle()
		}

		if _, err := mail.ParseAddress(upd.Email); err != nil {
			return UserErrInvalidEmail()
		}

		if u, err = store.LookupUserByID(svc.ctx, svc.store, upd.ID); err != nil {
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
		u.UpdatedAt = now()

		if err = svc.eventbus.WaitFor(svc.ctx, event.UserBeforeUpdate(upd, u)); err != nil {
			return
		}

		if err = uniqueUserCheck(svc.ctx, svc.store, u); err != nil {
			return
		}

		if err = store.UpdateUser(svc.ctx, svc.store, u); err != nil {
			return
		}

		if label.Changed(u.Labels, upd.Labels) {
			if err = label.Update(svc.ctx, svc.store, upd); err != nil {
				return
			}

			u.Labels = upd.Labels
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.UserAfterUpdate(upd, u))
		return
	}()

	return u, svc.recordAction(svc.ctx, uaProps, UserActionUpdate, err)
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

	err = func() (err error) {
		if userID == 0 {
			return UserErrInvalidID()
		}

		if u, err = store.LookupUserByID(svc.ctx, svc.store, userID); err != nil {
			return
		}

		if !svc.ac.CanDeleteUser(svc.ctx, u) {
			return UserErrNotAllowedToDelete()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.UserBeforeUpdate(nil, u)); err != nil {
			return
		}

		u.DeletedAt = now()
		if err = store.UpdateUser(svc.ctx, svc.store, u); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.UserAfterDelete(nil, u))
		return nil
	}()

	return svc.recordAction(svc.ctx, uaProps, UserActionDelete, err)
}

func (svc user) Undelete(userID uint64) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if userID == 0 {
			return UserErrInvalidID()
		}

		if u, err = store.LookupUserByID(svc.ctx, svc.store, userID); err != nil {
			return
		}

		uaProps.setUser(u)

		if err = uniqueUserCheck(svc.ctx, svc.store, u); err != nil {
			return err
		}

		if !svc.ac.CanDeleteUser(svc.ctx, u) {
			return UserErrNotAllowedToDelete()
		}

		u.DeletedAt = nil
		if err = store.UpdateUser(svc.ctx, svc.store, u); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(svc.ctx, uaProps, UserActionUndelete, err)

}

func (svc user) Suspend(userID uint64) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if userID == 0 {
			return UserErrInvalidID()
		}

		if u, err = store.LookupUserByID(svc.ctx, svc.store, userID); err != nil {
			return
		}

		uaProps.setUser(u)

		if !svc.ac.CanSuspendUser(svc.ctx, u) {
			return UserErrNotAllowedToSuspend()
		}

		u.SuspendedAt = now()
		if err = store.UpdateUser(svc.ctx, svc.store, u); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(svc.ctx, uaProps, UserActionSuspend, err)

}

func (svc user) Unsuspend(userID uint64) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if userID == 0 {
			return UserErrInvalidID()
		}

		if u, err = store.LookupUserByID(svc.ctx, svc.store, userID); err != nil {
			return
		}

		uaProps.setUser(u)

		if !svc.ac.CanUnsuspendUser(svc.ctx, u) {
			return UserErrNotAllowedToUnsuspend()
		}

		u.SuspendedAt = nil
		if err = store.UpdateUser(svc.ctx, svc.store, u); err != nil {
			return
		}
		return nil
	}()

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

	err = func() error {
		if u, err = store.LookupUserByID(svc.ctx, svc.store, userID); err != nil {
			return err
		}

		uaProps.setUser(u)

		if !svc.ac.CanUpdateUser(svc.ctx, u) {
			return UserErrNotAllowedToUpdate()
		}

		if !svc.auth.CheckPasswordStrength(newPassword) {
			return UserErrPasswordNotSecure()
		}

		if err := svc.auth.SetPasswordCredentials(svc.ctx, userID, newPassword); err != nil {
			return err
		}

		return nil
	}()

	return svc.recordAction(svc.ctx, uaProps, UserActionSetPassword, err)

}

// Masks (or leaves as-is) private data on user
func (svc user) handlePrivateData(ctx context.Context, u *types.User) {
	if svc.maskEmail(ctx, u) {
		u.Email = maskPrivateDataEmail
	}

	if svc.maskName(ctx, u) {
		u.Name = maskPrivateDataName
	}
}

func (svc user) maskEmail(ctx context.Context, u *types.User) bool {
	return svc.settings.Privacy.Mask.Email && !svc.ac.CanUnmaskEmail(ctx, u)
}

func (svc user) maskName(ctx context.Context, u *types.User) bool {
	return svc.settings.Privacy.Mask.Name && !svc.ac.CanUnmaskName(ctx, u)
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

// UniqueCheck verifies user's email, username and handle
func uniqueUserCheck(ctx context.Context, s store.Storer, u *types.User) (err error) {
	isUnique := func(field string) bool {
		f := types.UserFilter{
			// If user exists and is deleted -- not a dup
			Deleted: filter.StateExcluded,

			// If user exists and is suspended -- duplicate
			Suspended: filter.StateInclusive,
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

		set, _, err := store.SearchUsers(ctx, s, f)
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

func createUserHandle(ctx context.Context, s store.Users, u *types.User) {
	if u.Handle == "" {
		u.Handle, _ = handle.Cast(
			// Must not exist before
			func(lookup string) bool {
				e, err := s.LookupUserByHandle(ctx, lookup)
				return err == store.ErrNotFound && (e == nil || e.ID == u.ID)
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

// toLabeledUsers converts to []label.LabeledResource
//
// This function is auto-generated.
func toLabeledUsers(set []*types.User) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
