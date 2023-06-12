package service

import (
	"context"
	"io"
	"mime/multipart"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/handle"
	"github.com/cortezaproject/corteza/server/pkg/label"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/service/event"
	"github.com/cortezaproject/corteza/server/system/types"
)

const (
	maskPrivateDataEmail = "####.#######@######.###"
	maskPrivateDataName  = "##### ##########"
)

type (
	user struct {
		actionlog actionlog.Recorder

		settings *types.AppSettings

		auth userAuth

		ac       userAccessController
		eventbus eventDispatcher

		store store.Storer

		opt UserOptions

		// List (cache) of preloaded users, accessible by handle
		//
		// It also does negative caching by assigning empty User structs
		preloaded map[string]*types.User

		att AttachmentService
	}

	synteticUserDataGen interface {
		Name() string
		Username() string
		Number(int, int) int
	}

	UserOptions struct {
		LimitUsers int
	}

	userAuth interface {
		CheckPasswordStrength(string) bool
		SetPasswordCredentials(context.Context, uint64, string) error
		RemovePasswordCredentials(context.Context, uint64) error
		RemoveAccessTokens(context.Context, *types.User) error
	}

	userAccessController interface {
		CanSearchUsers(context.Context) bool
		CanCreateUser(context.Context) bool
		CanReadUser(context.Context, *types.User) bool
		CanUpdateUser(context.Context, *types.User) bool
		CanDeleteUser(context.Context, *types.User) bool
		CanSuspendUser(context.Context, *types.User) bool
		CanUnsuspendUser(context.Context, *types.User) bool
		CanUnmaskEmailOnUser(context.Context, *types.User) bool
		CanUnmaskNameOnUser(context.Context, *types.User) bool
	}

	UserService interface {
		FindByEmail(ctx context.Context, email string) (*types.User, error)
		FindByHandle(ctx context.Context, handle string) (*types.User, error)
		FindByID(ctx context.Context, id uint64) (*types.User, error)
		FindByAny(ctx context.Context, identifier interface{}) (*types.User, error)
		Find(context.Context, types.UserFilter) (types.UserSet, types.UserFilter, error)

		Create(ctx context.Context, input *types.User) (*types.User, error)
		Update(ctx context.Context, mod *types.User) (*types.User, error)
		ToggleEmailConfirmation(ctx context.Context, userID uint64, confirm bool) error

		CreateWithAvatar(ctx context.Context, input *types.User, avatar io.Reader) (*types.User, error)
		UpdateWithAvatar(ctx context.Context, mod *types.User, avatar io.Reader) (*types.User, error)

		Delete(ctx context.Context, id uint64) error
		Suspend(ctx context.Context, id uint64) error
		Unsuspend(ctx context.Context, id uint64) error
		Undelete(ctx context.Context, id uint64) error

		SetPassword(ctx context.Context, userID uint64, password string) error

		DeleteAuthTokensByUserID(ctx context.Context, userID uint64) (err error)
		DeleteAuthSessionsByUserID(ctx context.Context, userID uint64) (err error)

		UploadAvatar(ctx context.Context, userID uint64, Upload *multipart.FileHeader) (err error)
		GenerateAvatar(ctx context.Context, userID uint64, bgColor string, initialColor string) (err error)
		DeleteAvatar(ctx context.Context, id uint64) error
	}
)

func User(opt UserOptions) *user {
	return &user{
		eventbus: eventbus.Service(),
		ac:       DefaultAccessControl,
		settings: CurrentSettings,
		auth:     DefaultAuth,

		store: DefaultStore,

		actionlog: DefaultActionlog,

		opt: opt,

		preloaded: make(map[string]*types.User),
		att:       DefaultAttachment,
	}
}

func (svc user) FindByID(ctx context.Context, userID uint64) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() error {
		u, err = loadUser(ctx, svc.store, userID)
		if u, err = svc.proc(ctx, u, err); err != nil {
			return err
		}

		uaProps.setUser(u)

		// If profile avatar settings is enabled and a user doesn't have an avatar image,
		// generate one automatically when fetching their user information.
		if svc.settings.Auth.Internal.ProfileAvatar.Enabled && u.Meta.AvatarID == 0 && u.Meta.AvatarColor == "" {
			if err = svc.generateUserAvatarInitial(ctx, u); err != nil {
				return err
			}
		}

		if !svc.ac.CanReadUser(ctx, u) {
			return UserErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, u); err != nil {
			return err
		}

		return nil
	}()

	return u, svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc user) FindByEmail(ctx context.Context, email string) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{user: &types.User{Email: email}}
	)

	err = func() error {

		u, err = store.LookupUserByEmail(ctx, svc.store, email)
		if u, err = svc.proc(ctx, u, err); err != nil {
			return err
		}

		uaProps.setUser(u)

		if !svc.ac.CanReadUser(ctx, u) {
			return UserErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, u); err != nil {
			return err
		}

		return nil
	}()

	return u, svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

func (svc user) FindByHandle(ctx context.Context, handle string) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{user: &types.User{Handle: handle}}
	)

	err = func() error {
		u, err = store.LookupUserByHandle(ctx, svc.store, handle)
		if u, err = svc.proc(ctx, u, err); err != nil {
			return err
		}

		uaProps.setUser(u)

		if !svc.ac.CanReadUser(ctx, u) {
			return UserErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, u); err != nil {
			return err
		}

		return nil
	}()

	return u, svc.recordAction(ctx, uaProps, UserActionLookup, err)
}

// FindByAny finds user by given identifier (context, id, handle, email)
func (svc user) FindByAny(ctx context.Context, identifier interface{}) (u *types.User, err error) {
	if ctx, ok := identifier.(context.Context); ok {
		identifier = internalAuth.GetIdentityFromContext(ctx).Identity()
	}

	if ID, ok := identifier.(uint64); ok {
		u, err = svc.FindByID(ctx, ID)
	} else if identity, ok := identifier.(internalAuth.Identifiable); ok {
		u, err = svc.FindByID(ctx, identity.Identity())
	} else if strIdentifier, ok := identifier.(string); ok {
		if ID, _ := strconv.ParseUint(strIdentifier, 10, 64); ID > 0 {
			u, err = svc.FindByID(ctx, ID)
		} else if strings.Contains(strIdentifier, "@") {
			u, err = svc.FindByEmail(ctx, strIdentifier)
		} else {
			u, err = svc.FindByHandle(ctx, strIdentifier)
		}
	} else {
		err = UserErrInvalidID()
	}

	if err != nil {
		return
	}

	rr, _, err := store.SearchRoles(ctx, svc.store, types.RoleFilter{MemberID: u.ID})
	if err != nil {
		return nil, err
	}

	u.SetRoles(rr.IDs()...)
	return
}

func (svc user) proc(ctx context.Context, u *types.User, err error) (*types.User, error) {
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, UserErrNotFound()
		}

		return nil, err
	}

	svc.handlePrivateData(ctx, u)

	return u, nil
}

// Find interacts with backend storage and
//
// @todo rename to Search() for consistency
func (svc user) Find(ctx context.Context, filter types.UserFilter) (uu types.UserSet, f types.UserFilter, err error) {
	var (
		uaProps = &userActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.MaskedEmailsEnabled = svc.settings.Privacy.Mask.Email
	filter.MaskedNamesEnabled = svc.settings.Privacy.Mask.Name
	filter.Check = func(res *types.User) (bool, error) {
		if !svc.ac.CanReadUser(ctx, res) {
			return false, nil
		}

		if svc.maskEmail(ctx, res) && ((len(filter.Query) > 0 && strings.HasPrefix(res.Email, filter.Query)) || res.Email == filter.Email) {
			// user email matched but it will be masked later on, so exclude it to prevent data probing
			return false, nil
		}

		if svc.maskName(ctx, res) && (len(filter.Query) > 0 && strings.HasPrefix(res.Name, filter.Query)) {
			// user mail matched but it will be masked later on, so exclude it to prevent data probing
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchUsers(ctx) {
			return UserErrNotAllowedToSearch()
		}

		if filter.Deleted > 0 {
			// If list with deleted users is requested
			// user must have access permissions to system (ie: is admin)
			//
			// not the best solution but ATM it allows us to have at least
			// some kind of control over who can see deleted users
			//if !svc.ac.CanAccess(ctx) {
			//	return UserErrNotAllowedToListUsers()
			//}
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
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

		uu, f, err = store.SearchUsers(ctx, svc.store, filter)
		if err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledUsers(uu)...); err != nil {
			return err
		}

		return uu.Walk(func(u *types.User) error {
			svc.handlePrivateData(ctx, u)
			return nil
		})
	}()

	return uu, f, svc.recordAction(ctx, uaProps, UserActionSearch, err)
}

func (svc user) Create(ctx context.Context, new *types.User) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{user: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateUser(ctx) {
			return UserErrNotAllowedToCreate()
		}

		if new.Kind == types.SystemUser {
			return UserErrNotAllowedToCreateSystem()
		}

		if !handle.IsValid(new.Handle) {
			return UserErrInvalidHandle()
		}

		if _, err := mail.ParseAddress(new.Email); err != nil {
			return UserErrInvalidEmail()
		}

		if err = svc.checkLimits(ctx); err != nil {
			return err
		}

		if err = svc.eventbus.WaitFor(ctx, event.UserBeforeCreate(new, u)); err != nil {
			return
		}

		if new.Handle == "" {
			createUserHandle(ctx, DefaultStore, new)
		}

		if err = uniqueUserCheck(ctx, svc.store, new); err != nil {
			return
		}

		// Process avatar initials Image
		if err = svc.generateUserAvatarInitial(ctx, new); err != nil {
			return
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		// consider email confirmed
		// when creating user like this
		new.EmailConfirmed = true

		if err = store.CreateUser(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.UserAfterCreate(new, u))
		return
	}()

	return new, svc.recordAction(ctx, uaProps, UserActionCreate, err)
}

func (svc user) CreateWithAvatar(ctx context.Context, input *types.User, avatar io.Reader) (out *types.User, err error) {
	// @todo: avatar
	return svc.Create(ctx, input)
}

func (svc user) Update(ctx context.Context, upd *types.User) (u *types.User, err error) {
	var (
		uaProps = &userActionProps{update: upd}
	)

	err = func() (err error) {
		if !handle.IsValid(upd.Handle) {
			return UserErrInvalidHandle()
		}

		if _, err := mail.ParseAddress(upd.Email); err != nil {
			return UserErrInvalidEmail()
		}

		if u, err = loadUser(ctx, svc.store, upd.ID); err != nil {
			return
		}

		uaProps.setUser(u)

		if upd.Kind == types.SystemUser || u.Kind == types.SystemUser {
			return UserErrNotAllowedToUpdateSystem()
		}

		if upd.ID != internalAuth.GetIdentityFromContext(ctx).Identity() {
			if !svc.ac.CanUpdateUser(ctx, u) {
				return UserErrNotAllowedToUpdate()
			}
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, u.UpdatedAt, u.CreatedAt) {
			return UserErrStaleData()
		}

		// Assign changed values
		u.Email = upd.Email
		u.Username = upd.Username
		u.Name = upd.Name
		u.Handle = upd.Handle
		u.Kind = upd.Kind
		u.UpdatedAt = now()

		if upd.Meta != nil {
			// Only update meta when set
			u.Meta = upd.Meta
		}

		if err = svc.generateUserAvatarInitial(ctx, u); err != nil {
			return
		}

		if err = svc.eventbus.WaitFor(ctx, event.UserBeforeUpdate(upd, u)); err != nil {
			return
		}

		if err = uniqueUserCheck(ctx, svc.store, u); err != nil {
			return
		}

		if err = store.UpdateUser(ctx, svc.store, u); err != nil {
			return
		}

		if label.Changed(u.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}

			u.Labels = upd.Labels
		}

		_ = svc.eventbus.WaitFor(ctx, event.UserAfterUpdate(upd, u))
		return
	}()

	return u, svc.recordAction(ctx, uaProps, UserActionUpdate, err)
}

func (svc user) ToggleEmailConfirmation(ctx context.Context, userID uint64, confirmed bool) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{}
	)

	err = func() (err error) {
		if u, err = loadUser(ctx, svc.store, userID); err != nil {
			return
		}

		uaProps.setUser(u)

		if userID != internalAuth.GetIdentityFromContext(ctx).Identity() {
			if !svc.ac.CanUpdateUser(ctx, u) {
				return UserErrNotAllowedToUpdate()
			}
		}

		u.EmailConfirmed = confirmed

		if err = store.UpdateUser(ctx, svc.store, u); err != nil {
			return
		}
		return
	}()

	return svc.recordAction(ctx, uaProps, UserActionUpdate, err)
}

func (svc user) UpdateWithAvatar(ctx context.Context, mod *types.User, avatar io.Reader) (out *types.User, err error) {
	// @todo: avatar
	return svc.Create(ctx, mod)
}

func (svc user) Delete(ctx context.Context, userID uint64) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if u, err = loadUser(ctx, svc.store, userID); err != nil {
			return
		}

		if u.Kind == types.SystemUser {
			return UserErrNotAllowedToDelete()
		}

		if !svc.ac.CanDeleteUser(ctx, u) {
			return UserErrNotAllowedToDelete()
		}

		if err = svc.eventbus.WaitFor(ctx, event.UserBeforeDelete(nil, u)); err != nil {
			return
		}

		u.DeletedAt = now()
		if err = store.UpdateUser(ctx, svc.store, u); err != nil {
			return
		}

		if err = svc.auth.RemoveAccessTokens(ctx, u); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.UserAfterDelete(nil, u))
		return nil
	}()

	return svc.recordAction(ctx, uaProps, UserActionDelete, err)
}

func (svc user) Undelete(ctx context.Context, userID uint64) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if u, err = loadUser(ctx, svc.store, userID); err != nil {
			return
		}

		uaProps.setUser(u)

		if err = uniqueUserCheck(ctx, svc.store, u); err != nil {
			return err
		}

		if u.Kind == types.SystemUser {
			return UserErrNotAllowedToDelete()
		}

		if err = svc.checkLimits(ctx); err != nil {
			return err
		}

		if !svc.ac.CanDeleteUser(ctx, u) {
			return UserErrNotAllowedToDelete()
		}

		u.DeletedAt = nil
		if err = store.UpdateUser(ctx, svc.store, u); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, uaProps, UserActionUndelete, err)

}

func (svc user) Suspend(ctx context.Context, userID uint64) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if u, err = loadUser(ctx, svc.store, userID); err != nil {
			return
		}

		uaProps.setUser(u)

		if u.Kind == types.SystemUser {
			return UserErrNotAllowedToSuspend()
		}

		if !svc.ac.CanSuspendUser(ctx, u) {
			return UserErrNotAllowedToSuspend()
		}

		// Clone u to oldUser
		oldUser := *u
		u.SuspendedAt = now()

		if err = svc.eventbus.WaitFor(ctx, event.UserBeforeSuspend(u, &oldUser)); err != nil {
			return
		}

		if err = store.UpdateUser(ctx, svc.store, u); err != nil {
			return
		}

		if err = svc.auth.RemoveAccessTokens(ctx, u); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.UserAfterSuspend(u, &oldUser))
		return nil
	}()

	return svc.recordAction(ctx, uaProps, UserActionSuspend, err)

}

func (svc user) Unsuspend(ctx context.Context, userID uint64) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if u, err = loadUser(ctx, svc.store, userID); err != nil {
			return
		}

		uaProps.setUser(u)

		if u.Kind == types.SystemUser {
			return UserErrNotAllowedToUnsuspend()
		}

		if !svc.ac.CanUnsuspendUser(ctx, u) {
			return UserErrNotAllowedToUnsuspend()
		}

		if err = svc.checkLimits(ctx); err != nil {
			return err
		}

		u.SuspendedAt = nil
		if err = store.UpdateUser(ctx, svc.store, u); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, uaProps, UserActionUnsuspend, err)
}

// SetPassword sets new password for a user
//
// Expecting setter to have permissions to update users
func (svc user) SetPassword(ctx context.Context, userID uint64, newPassword string) (err error) {
	var (
		u *types.User

		uaProps = &userActionProps{user: &types.User{ID: userID}}
		a       = UserActionSetPassword

		self = internalAuth.GetIdentityFromContext(ctx).Identity() == userID
	)

	err = func() (err error) {
		if u, err = loadUser(ctx, svc.store, userID); err != nil {
			return err
		}

		uaProps.setUser(u)

		if !svc.ac.CanUpdateUser(ctx, u) {
			return UserErrNotAllowedToUpdate()
		}

		if u.Kind == types.SystemUser {
			return UserErrNotAllowedToUpdateSystem()
		}

		if !self {
			// when user is changing password for herself
			// we should not remove the tokens!
			//
			// without this, user needs to log-in again
			// and we do not want that if he is using general
			// user management API/UI
			if err = svc.auth.RemoveAccessTokens(ctx, u); err != nil {
				return
			}
		}

		if newPassword == "" {
			a = UserActionRemovePassword
			return svc.auth.RemovePasswordCredentials(ctx, userID)
		}

		// note on password reuse:
		//
		// we do not really care if user is setting same password
		// to someone else (or to self for that matter)
		//
		// he has rights to update the user and is doing so
		// through general user management API

		if !svc.auth.CheckPasswordStrength(newPassword) {
			return UserErrPasswordNotSecure()
		}

		if err = svc.auth.SetPasswordCredentials(ctx, userID, newPassword); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, uaProps, a, err)

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
	return svc.settings.Privacy.Mask.Email && !svc.ac.CanUnmaskEmailOnUser(ctx, u)
}

func (svc user) maskName(ctx context.Context, u *types.User) bool {
	return svc.settings.Privacy.Mask.Name && !svc.ac.CanUnmaskNameOnUser(ctx, u)
}

func (svc *user) Get(ctx context.Context, h string) (u *types.User, err error) {
	if svc.preloaded[h] == nil {
		svc.preloaded[h], err = svc.FindByHandle(ctx, h)
		if err != nil {
			svc.preloaded[h] = &types.User{}
			return
		}
	}

	if svc.preloaded[h] == nil || svc.preloaded[h].ID == 0 {
		return nil, UserErrNotFound()
	}

	return svc.preloaded[h], nil
}

// DeleteAuthTokensByUserID will delete all auth tokens of user which will un-authorize all auth clients of user
func (svc user) DeleteAuthTokensByUserID(ctx context.Context, userID uint64) (err error) {
	var (
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if userID == 0 {
			return UserErrInvalidID()
		}

		if err = store.DeleteAuthOA2TokenByUserID(ctx, svc.store, userID); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, uaProps, UserActionDeleteAuthTokens, err)
}

// DeleteAuthSessionsByUserID will delete all auth session of user
func (svc user) DeleteAuthSessionsByUserID(ctx context.Context, userID uint64) (err error) {
	var (
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if userID == 0 {
			return UserErrInvalidID()
		}

		if err = store.DeleteAuthSessionsByUserID(ctx, svc.store, userID); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, uaProps, UserActionDeleteAuthSessions, err)
}

func (svc user) checkLimits(ctx context.Context) error {
	if svc.opt.LimitUsers == 0 {
		return nil
	}

	if c, err := countValidUsers(ctx, svc.store); err != nil {
		return err
	} else if c >= uint(svc.opt.LimitUsers) {
		return UserErrMaxUserLimitReached()
	}

	return nil
}

// CreateSynthetic generates, saves and returns new user
//
// Generated users will have their handles prefixed with "synthetic_" and email domain "synthetic.tld"
//
// Function checks if user can create users but avoids all other checks (besides unique value)
func (svc user) CreateSynthetic(ctx context.Context, src synteticUserDataGen, total uint) error {
	if !svc.ac.CanCreateUser(ctx) {
		return UserErrNotAllowedToCreate()
	}

	const maxRetries = 10

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var retry uint
		for total > 0 || maxRetries < retry {
			// even with pre-check for unique users this one
			// still returns not unique error from time to time ?!
			err = store.CreateUser(ctx, s, syntheticUser(src))
			if errors.IsDuplicateData(err) {
				retry++
				continue
			}

			if err != nil {
				return
			}

			retry = 0
			total--
		}

		return
	})
}

func syntheticUser(src synteticUserDataGen) (r *types.User) {
	r = &types.User{
		ID:             nextID(),
		Kind:           types.NormalUser,
		Name:           src.Name(),
		Handle:         "synthetic_" + src.Username(),
		EmailConfirmed: src.Number(0, 1) > 0,

		// Make sure all users are created in the past
		CreatedAt: time.Now().Add(time.Hour * time.Duration(src.Number(100000, 1000000)*-1)),
	}

	r.Email = strings.ToLower(strings.ReplaceAll(r.Name, " ", ".")) + "@synthetic.tld"

	if src.Number(0, 1) > 0 {
		aux := time.Now().Add(time.Hour * time.Duration(src.Number(100, 100000)*-1))
		r.UpdatedAt = &aux
	}

	return
}

// CreateSynthetic generates, saves and returns new user
//
// Generated users will have their handles prefixed with "synthetic_" and email domain "synthetic.tld"
func (svc user) RemoveSynthetic(ctx context.Context) error {
	// not a mistake, we do not need or want to check if user can be deleted
	if !svc.ac.CanCreateUser(ctx) {
		return UserErrNotAllowedToCreate()
	}

	var (
		f  = types.UserFilter{Query: "@synthetic.tld"}
		uu types.UserSet
	)

	f.Limit = 1000

	// @todo this should be optimized by using store.DeleteUserByFilter
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		for {
			uu, _, err = store.SearchUsers(ctx, s, f)
			if len(uu) == 0 || err != nil {
				// when nothing is fetch or error returned
				// break out of the loop
				return
			}

			for _, u := range uu {
				// check if handle starts with synthetic_
				if !strings.HasPrefix(u.Handle, "synthetic_") {
					continue
				}

				// check if email ends with synthetic.tld
				if !strings.HasSuffix(u.Email, "@synthetic.tld") {
					continue
				}

				if err = store.DeleteUser(ctx, s, u); err != nil {
					return
				}
			}
		}

		return
	})

}

func loadUser(ctx context.Context, s store.Users, ID uint64) (res *types.User, err error) {
	if ID == 0 {
		return nil, UserErrInvalidID()
	}

	if res, err = store.LookupUserByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, UserErrNotFound()
	}

	return
}

func countValidUsers(ctx context.Context, s store.Users) (c uint, err error) {
	return store.CountUsers(ctx, s, types.UserFilter{Kind: types.NormalUser})
}

// uniqueUserCheck verifies user's email, username and handle
func uniqueUserCheck(ctx context.Context, s store.Storer, u *types.User) (err error) {
	isUnique := func(field string) bool {
		f := types.UserFilter{
			// If user exists and is deleted -- not a dup
			Deleted: filter.StateExcluded,

			// If user exists and is suspended -- duplicate
			Suspended: filter.StateInclusive,
		}

		f.Limit = 1

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

func (svc user) UploadAvatar(ctx context.Context, userID uint64, upload *multipart.FileHeader) (err error) {
	var (
		u       *types.User
		att     *types.Attachment
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if u, err = loadUser(ctx, svc.store, userID); err != nil {
			return
		}

		if userID != internalAuth.GetIdentityFromContext(ctx).Identity() {
			if !svc.ac.CanUpdateUser(ctx, u) {
				return UserErrNotAllowedToUpdate()
			}
		}

		if u.Meta.AvatarID != 0 {
			if err = svc.att.DeleteByID(ctx, u.Meta.AvatarID); err != nil {
				return
			}
		}

		file, err := upload.Open()
		if err != nil {
			return err
		}

		defer file.Close()

		att, err = svc.att.CreateAuthAttachment(
			ctx,
			upload.Filename,
			upload.Size,
			file,
			map[string]string{"key": types.AttachmentKindAvatar},
		)

		if err != nil {
			return err
		}

		u.Meta.AvatarID = att.ID
		u.Meta.AvatarKind = types.AttachmentKindAvatar

		if err = store.UpdateUser(ctx, svc.store, u); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, uaProps, UserActionUploadAvatar, err)
}

// DeleteAvatar will delete user's avatar
func (svc user) DeleteAvatar(ctx context.Context, userID uint64) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if u, err = svc.FindByID(ctx, userID); err != nil {
			return
		}

		if u.Kind == types.SystemUser {
			return UserErrNotAllowedToDeleteAvatar()
		}

		att, err := svc.att.FindByID(ctx, u.Meta.AvatarID)
		if err != nil {
			return err
		}

		if att.Meta.Labels["key"] != types.AttachmentKindAvatar {
			return nil
		}

		if !svc.ac.CanUpdateUser(ctx, u) {
			return UserErrNotAllowedToDeleteAvatar()
		}

		if err = svc.att.DeleteByID(ctx, u.Meta.AvatarID); err != nil {
			return err
		}

		u.Meta.AvatarID = 0

		// When an uploaded avatar is deleted, generate avatar initial
		if err = svc.generateUserAvatarInitial(ctx, u); err != nil {
			return err
		}

		if err = store.UpdateUser(ctx, svc.store, u); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, uaProps, UserActionDeleteAvatar, err)
}

func processAvatarInitials(u *types.User) (initial string) {
	var (
		chars string
		parts []string
	)

	if u.Name != "" {
		parts = strings.Fields(u.Name)
		if len(parts) > 2 {
			chars = string(parts[0][0]) + string(parts[1][0]) + string(parts[2][0])
		} else if len(parts) > 1 {
			chars = string(parts[0][0]) + string(parts[1][0])
		} else {
			if len(parts[0]) > 1 {
				chars = string(parts[0][0]) + string(parts[0][1])
			} else {
				chars = string(parts[0][0])
			}
		}
	} else if u.Handle != "" {
		if strings.ContainsAny(u.Handle, "._-") {
			for _, del := range "._-" {
				if strings.Contains(u.Handle, string(del)) {
					parts = strings.Split(u.Handle, string(del))
					break
				}
			}

			chars = string(parts[0][0]) + string(parts[1][0])
		} else {
			chars = string(u.Handle[0])
		}
	} else {
		email := strings.Split(u.Email, "@")
		if strings.ContainsAny(email[0], "._-") {
			for _, del := range "._-" {
				if strings.Contains(email[0], string(del)) {
					parts = strings.Split(email[0], string(del))
					break
				}
			}

			chars = string(parts[0][0]) + string(parts[1][0])
		} else {
			chars = string(email[0][0])
		}
	}

	// Validate initials: if initials are letters if not assign a default "CU"
	for _, c := range chars {
		if unicode.IsLetter(c) {
			initial += string(c)
		}
		continue
	}
	if initial == "" {
		initial = "CU"
	}

	initial = strings.ToUpper(initial)

	return
}

func (svc user) GenerateAvatar(ctx context.Context, userID uint64, bgColor string, initialColor string) (err error) {
	var (
		u       *types.User
		uaProps = &userActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if u, err = loadUser(ctx, svc.store, userID); err != nil {
			return
		}

		if userID != internalAuth.GetIdentityFromContext(ctx).Identity() {
			if !svc.ac.CanUpdateUser(ctx, u) {
				return UserErrNotAllowedToUpdate()
			}
		}

		u.Meta.AvatarColor = initialColor
		u.Meta.AvatarBgColor = bgColor
		if err = svc.generateUserAvatarInitial(ctx, u); err != nil {
			return err
		}

		if err = store.UpdateUser(ctx, svc.store, u); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, uaProps, UserActionGenerateAvatar, err)
}

func (svc user) generateUserAvatarInitial(ctx context.Context, u *types.User) (err error) {
	var (
		att *types.Attachment
	)

	initial := processAvatarInitials(u)

	if u.Meta == nil {
		u.Meta = &types.UserMeta{}
	}

	if u.Meta.AvatarID != 0 {
		if att, err = svc.att.FindByID(ctx, u.Meta.AvatarID); err != nil {
			return err
		}

		if att.Meta.Labels["key"] == types.AttachmentKindAvatar {
			return nil
		}

		colorLogic := att.Meta.Original.Image.BackgroundColor == u.Meta.AvatarBgColor && att.Meta.Original.Image.InitialColor == u.Meta.AvatarColor
		if att.Meta.Original.Image.Initial == initial && colorLogic {
			return nil
		}

		if err = svc.att.DeleteByID(ctx, att.ID); err != nil {
			return err
		}
	}

	if att, err = svc.att.CreateAvatarInitialsAttachment(ctx, initial, u.Meta.AvatarBgColor, u.Meta.AvatarColor); err != nil {
		return err
	}

	if u.Meta == nil {
		u.Meta = &types.UserMeta{}
	}

	u.Meta.AvatarID = att.ID
	u.Meta.AvatarKind = types.AttachmentKindAvatarInitials

	if u.Meta.AvatarBgColor == "" {
		u.Meta.AvatarBgColor = att.Meta.Original.Image.BackgroundColor
	}

	if u.Meta.AvatarColor == "" {
		u.Meta.AvatarColor = att.Meta.Original.Image.InitialColor
	}

	return nil
}
