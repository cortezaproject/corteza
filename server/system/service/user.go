package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	internalAuth "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/sass"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service/event"
	"github.com/crusttech/human/server/system/types"
)

const (
	maskPrivateDataEmail = "####.#######@######.###"
	maskPrivateDataName  = "##### ##########"
)

type (
	userServices struct {
		settings  *types.AppSettings
		auth      userAuth
		eventbus  eventDispatcher
		opt       UserOptions
		preloaded map[string]*types.User
		att       AttachmentService
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
		Search(context.Context, types.UserFilter) (types.UserSet, types.UserFilter, error)

		Create(ctx context.Context, input *types.User) (*types.User, error)
		Update(ctx context.Context, mod *types.User) (*types.User, error)
		ToggleEmailConfirmation(ctx context.Context, userID uint64, confirm bool) error

		CreateWithAvatar(ctx context.Context, input *types.User, avatar io.Reader) (*types.User, error)
		UpdateWithAvatar(ctx context.Context, mod *types.User, avatar io.Reader) (*types.User, error)

		Delete(ctx context.Context, id uint64) error
		DeleteByID(ctx context.Context, id uint64) error
		Suspend(ctx context.Context, id uint64) error
		Unsuspend(ctx context.Context, id uint64) error
		Undelete(ctx context.Context, id uint64) error
		UndeleteByID(ctx context.Context, id uint64) error

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
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		services: &userServices{
			eventbus:  eventbus.Service(),
			settings:  CurrentSettings,
			auth:      DefaultAuth,
			opt:       opt,
			preloaded: make(map[string]*types.User),
			att:       DefaultAttachment,
		},
	}
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

// Find interacts with backend storage
func (svc user) Find(ctx context.Context, filter types.UserFilter) (uu types.UserSet, f types.UserFilter, err error) {
	return svc.Search(ctx, filter)
}

func (svc user) CreateWithAvatar(ctx context.Context, input *types.User, avatar io.Reader) (out *types.User, err error) {
	return svc.Create(ctx, input)
}

func (svc user) UpdateWithAvatar(ctx context.Context, mod *types.User, avatar io.Reader) (out *types.User, err error) {
	return svc.Update(ctx, mod)
}

// Delete soft-deletes a user by ID
func (svc user) Delete(ctx context.Context, userID uint64) error {
	return svc.DeleteByID(ctx, userID)
}

// Undelete restores a soft-deleted user by ID
func (svc user) Undelete(ctx context.Context, userID uint64) error {
	return svc.UndeleteByID(ctx, userID)
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

func (svc *user) Get(ctx context.Context, h string) (u *types.User, err error) {
	if svc.services.preloaded[h] == nil {
		svc.services.preloaded[h], err = svc.FindByHandle(ctx, h)
		if err != nil {
			svc.services.preloaded[h] = &types.User{}
			return
		}
	}

	if svc.services.preloaded[h] == nil || svc.services.preloaded[h].ID == 0 {
		return nil, UserErrNotFound()
	}

	return svc.services.preloaded[h], nil
}

func (svc user) checkLimits(ctx context.Context) error {
	if svc.services.opt.LimitUsers == 0 {
		return nil
	}

	if c, err := countValidUsers(ctx, svc.store); err != nil {
		return err
	} else if c >= uint(svc.services.opt.LimitUsers) {
		return UserErrMaxUserLimitReached()
	}

	return nil
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
	return svc.services.settings.Privacy.Mask.Email && !svc.ac.CanUnmaskEmailOnUser(ctx, u)
}

func (svc user) maskName(ctx context.Context, u *types.User) bool {
	return svc.services.settings.Privacy.Mask.Name && !svc.ac.CanUnmaskNameOnUser(ctx, u)
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
		n := []string{
			fmt.Sprintf("%s_%s", u.Name, u.Username),
			regexp.
				MustCompile("(@.*)$").
				ReplaceAllString(u.Email, ""),
		}

		for i := 1; i <= 10; i++ {
			n = append(n, fmt.Sprintf("%s_%s%d", u.Name, u.Username, i))
		}

		u.Handle, _ = handle.Cast(
			func(lookup string) bool {
				e, err := s.LookupUserByHandle(ctx, lookup)
				return err == store.ErrNotFound && (e == nil || e.ID == u.ID)
			},
			n...,
		)
	}
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

func (svc user) generateUserAvatarInitial(ctx context.Context, u *types.User) (err error) {
	var (
		att *types.Attachment
	)

	initial := processAvatarInitials(u)

	if u.Meta == nil {
		u.Meta = &types.UserMeta{}
	}

	if u.Meta.AvatarID != 0 {
		if att, err = svc.services.att.FindByID(ctx, u.Meta.AvatarID); err != nil {
			return err
		}

		if att.Meta.Labels["key"] == types.AttachmentKindAvatar {
			return nil
		}

		colorLogic := att.Meta.Original.Image.BackgroundColor == u.Meta.AvatarBgColor && att.Meta.Original.Image.InitialColor == u.Meta.AvatarColor
		if att.Meta.Original.Image.Initial == initial && colorLogic {
			return nil
		}

		if err = svc.services.att.DeleteByID(ctx, att.ID); err != nil {
			return err
		}
	}

	if att, err = svc.services.att.CreateAvatarInitialsAttachment(ctx, initial, u.Meta.AvatarBgColor, u.Meta.AvatarColor); err != nil {
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

// --- on-hooks called by generated wrappers ---

func (svc *user) onLookup(ctx context.Context, ID uint64, aProps *userActionProps) (*types.User, error) {
	u, err := loadUser(ctx, svc.store, ID)
	if u, err = svc.proc(ctx, u, err); err != nil {
		return nil, err
	}

	aProps.setUser(u)

	// Auto-generate avatar initials when profile avatar is enabled and user has none
	if svc.services.settings.Auth.Internal.ProfileAvatar.Enabled && u.Meta.AvatarID == 0 && u.Meta.AvatarColor == "" {
		if err = svc.generateUserAvatarInitial(ctx, u); err != nil {
			return nil, err
		}
	}

	if !svc.ac.CanReadUser(ctx, u) {
		return nil, UserErrNotAllowedToRead()
	}

	if err = label.Load(ctx, svc.store, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (svc *user) onSearch(ctx context.Context, f types.UserFilter, aProps *userActionProps) (types.UserSet, types.UserFilter, error) {
	f.MaskedEmailsEnabled = svc.services.settings.Privacy.Mask.Email
	f.MaskedNamesEnabled = svc.services.settings.Privacy.Mask.Name
	f.Check = func(res *types.User) (bool, error) {
		if !svc.ac.CanReadUser(ctx, res) {
			return false, nil
		}

		if svc.maskEmail(ctx, res) && ((len(f.Query) > 0 && strings.HasPrefix(res.Email, f.Query)) || res.Email == f.Email) {
			return false, nil
		}

		if svc.maskName(ctx, res) && (len(f.Query) > 0 && strings.HasPrefix(res.Name, f.Query)) {
			return false, nil
		}

		return true, nil
	}

	if f.Deleted > 0 {
		// If list with deleted users is requested
		// user must have access permissions to system (ie: is admin)
		//
		// not the best solution but ATM it allows us to have at least
		// some kind of control over who can see deleted users
	}

	var err error
	if len(f.Labels) > 0 {
		f.LabeledIDs, err = label.Search(
			ctx,
			svc.store,
			types.User{}.LabelResourceKind(),
			f.Labels,
		)

		if err != nil {
			return nil, f, err
		}

		if len(f.LabeledIDs) == 0 {
			return nil, f, nil
		}
	}

	uu, out, err := store.SearchUsers(ctx, svc.store, f)
	if err != nil {
		return nil, out, err
	}

	if err = label.Load(ctx, svc.store, toLabeledUsers(uu)...); err != nil {
		return nil, out, err
	}

	err = uu.Walk(func(u *types.User) error {
		svc.handlePrivateData(ctx, u)
		return nil
	})

	return uu, out, err
}

func (svc *user) onCreate(ctx context.Context, new *types.User) error {
	if new.Kind == types.SystemUser {
		return UserErrNotAllowedToCreateSystem()
	}

	if !handle.IsValid(new.Handle) {
		return UserErrInvalidHandle()
	}

	if _, err := mail.ParseAddress(new.Email); err != nil {
		return UserErrInvalidEmail()
	}

	if err := svc.checkLimits(ctx); err != nil {
		return err
	}

	if err := svc.services.eventbus.WaitFor(ctx, event.UserBeforeCreate(new, nil)); err != nil {
		return err
	}

	if new.Handle == "" {
		createUserHandle(ctx, DefaultStore, new)
	}

	if err := uniqueUserCheck(ctx, svc.store, new); err != nil {
		return err
	}

	if new.Meta == nil {
		new.Meta = &types.UserMeta{}
	}

	if err := svc.generateUserAvatarInitial(ctx, new); err != nil {
		return err
	}

	new.Meta.Theme = sass.LightTheme

	new.ID = nextID()
	new.CreatedAt = *now()
	new.EmailConfirmed = true

	if err := store.CreateUser(ctx, svc.store, new); err != nil {
		return err
	}

	if err := label.Create(ctx, svc.store, new); err != nil {
		return err
	}

	_ = svc.services.eventbus.WaitFor(ctx, event.UserAfterCreate(new, nil))
	return nil
}

func (svc *user) onUpdate(ctx context.Context, s store.Storer, upd, res *types.User, aProps *userActionProps, before, after func() error) error {
	if upd.Kind == types.SystemUser || res.Kind == types.SystemUser {
		return UserErrNotAllowedToUpdateSystem()
	}

	if upd.ID != internalAuth.GetIdentityFromContext(ctx).Identity() {
		if !svc.ac.CanUpdateUser(ctx, res) {
			return UserErrNotAllowedToUpdate()
		}
	}

	if _, err := mail.ParseAddress(upd.Email); err != nil {
		return UserErrInvalidEmail()
	}

	if err := before(); err != nil {
		return err
	}

	res.Kind = upd.Kind
	if upd.Meta != nil {
		res.Meta = upd.Meta
	}

	if err := svc.generateUserAvatarInitial(ctx, res); err != nil {
		return err
	}

	if err := svc.services.eventbus.WaitFor(ctx, event.UserBeforeUpdate(upd, res)); err != nil {
		return err
	}

	if err := uniqueUserCheck(ctx, svc.store, res); err != nil {
		return err
	}

	if err := after(); err != nil {
		return err
	}

	return nil
}

func (svc *user) onDelete(ctx context.Context, s store.Storer, res *types.User, aProps *userActionProps) error {
	if res.Kind == types.SystemUser {
		return UserErrNotAllowedToDelete()
	}

	if !svc.ac.CanDeleteUser(ctx, res) {
		return UserErrNotAllowedToDelete()
	}

	if err := svc.services.eventbus.WaitFor(ctx, event.UserBeforeDelete(nil, res)); err != nil {
		return err
	}

	res.DeletedAt = now()
	if err := store.UpdateUser(ctx, s, res); err != nil {
		return err
	}

	if err := svc.services.auth.RemoveAccessTokens(ctx, res); err != nil {
		return err
	}

	_ = svc.services.eventbus.WaitFor(ctx, event.UserAfterDelete(nil, res))
	return nil
}

func (svc *user) onUndelete(ctx context.Context, s store.Storer, res *types.User, aProps *userActionProps) error {
	if err := uniqueUserCheck(ctx, svc.store, res); err != nil {
		return err
	}

	if res.Kind == types.SystemUser {
		return UserErrNotAllowedToDelete()
	}

	if err := svc.checkLimits(ctx); err != nil {
		return err
	}

	if !svc.ac.CanDeleteUser(ctx, res) {
		return UserErrNotAllowedToDelete()
	}

	res.DeletedAt = nil
	return store.UpdateUser(ctx, s, res)
}

func (svc *user) onFindByEmail(ctx context.Context, aProps *userActionProps, email string) (*types.User, error) {
	u, err := store.LookupUserByEmail(ctx, svc.store, email)
	if u, err = svc.proc(ctx, u, err); err != nil {
		return nil, err
	}

	aProps.setUser(u)

	if !svc.ac.CanReadUser(ctx, u) {
		return nil, UserErrNotAllowedToRead()
	}

	if err = label.Load(ctx, svc.store, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (svc *user) onFindByHandle(ctx context.Context, aProps *userActionProps, h string) (*types.User, error) {
	u, err := store.LookupUserByHandle(ctx, svc.store, h)
	if u, err = svc.proc(ctx, u, err); err != nil {
		return nil, err
	}

	aProps.setUser(u)

	if !svc.ac.CanReadUser(ctx, u) {
		return nil, UserErrNotAllowedToRead()
	}

	if err = label.Load(ctx, svc.store, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (svc *user) onToggleEmailConfirmation(ctx context.Context, aProps *userActionProps, userID uint64, confirmed bool) error {
	u, err := loadUser(ctx, svc.store, userID)
	if err != nil {
		return err
	}

	aProps.setUser(u)

	if userID != internalAuth.GetIdentityFromContext(ctx).Identity() {
		if !svc.ac.CanUpdateUser(ctx, u) {
			return UserErrNotAllowedToUpdate()
		}
	}

	u.EmailConfirmed = confirmed
	return store.UpdateUser(ctx, svc.store, u)
}

func (svc *user) onSuspend(ctx context.Context, aProps *userActionProps, userID uint64) error {
	u, err := loadUser(ctx, svc.store, userID)
	if err != nil {
		return err
	}

	aProps.setUser(u)

	if u.Kind == types.SystemUser {
		return UserErrNotAllowedToSuspend()
	}

	if !svc.ac.CanSuspendUser(ctx, u) {
		return UserErrNotAllowedToSuspend()
	}

	oldUser := *u
	u.SuspendedAt = now()

	if err = svc.services.eventbus.WaitFor(ctx, event.UserBeforeSuspend(u, &oldUser)); err != nil {
		return err
	}

	if err = store.UpdateUser(ctx, svc.store, u); err != nil {
		return err
	}

	if err = svc.services.auth.RemoveAccessTokens(ctx, u); err != nil {
		return err
	}

	_ = svc.services.eventbus.WaitFor(ctx, event.UserAfterSuspend(u, &oldUser))
	return nil
}

func (svc *user) onUnsuspend(ctx context.Context, aProps *userActionProps, userID uint64) error {
	u, err := loadUser(ctx, svc.store, userID)
	if err != nil {
		return err
	}

	aProps.setUser(u)

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
	return store.UpdateUser(ctx, svc.store, u)
}

func (svc *user) onSetPassword(ctx context.Context, aProps *userActionProps, userID uint64, newPassword string) error {
	u, err := loadUser(ctx, svc.store, userID)
	if err != nil {
		return err
	}

	aProps.setUser(u)

	if !svc.ac.CanUpdateUser(ctx, u) {
		return UserErrNotAllowedToUpdate()
	}

	if u.Kind == types.SystemUser {
		return UserErrNotAllowedToUpdateSystem()
	}

	self := internalAuth.GetIdentityFromContext(ctx).Identity() == userID
	if !self {
		if err = svc.services.auth.RemoveAccessTokens(ctx, u); err != nil {
			return err
		}
	}

	if newPassword == "" {
		return svc.services.auth.RemovePasswordCredentials(ctx, userID)
	}

	if !svc.services.auth.CheckPasswordStrength(newPassword) {
		return UserErrPasswordNotSecure()
	}

	return svc.services.auth.SetPasswordCredentials(ctx, userID, newPassword)
}

func (svc *user) onDeleteAuthTokensByUserID(ctx context.Context, aProps *userActionProps, userID uint64) error {
	if userID == 0 {
		return UserErrInvalidID()
	}

	return store.DeleteAuthOA2TokenByUserID(ctx, svc.store, userID)
}

func (svc *user) onDeleteAuthSessionsByUserID(ctx context.Context, aProps *userActionProps, userID uint64) error {
	if userID == 0 {
		return UserErrInvalidID()
	}

	return store.DeleteAuthSessionsByUserID(ctx, svc.store, userID)
}

func (svc *user) onCreateSynthetic(ctx context.Context, aProps *userActionProps, src synteticUserDataGen, total uint) error {
	const maxRetries = 10

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var retry uint
		for total > 0 || maxRetries < retry {
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

func (svc *user) onRemoveSynthetic(ctx context.Context, aProps *userActionProps) error {
	var (
		f  = types.UserFilter{Query: "@synthetic.tld"}
		uu types.UserSet
	)

	f.Limit = 1000

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		for {
			uu, _, err = store.SearchUsers(ctx, s, f)
			if len(uu) == 0 || err != nil {
				return
			}

			for _, u := range uu {
				if !strings.HasPrefix(u.Handle, "synthetic_") {
					continue
				}

				if !strings.HasSuffix(u.Email, "@synthetic.tld") {
					continue
				}

				if err = store.DeleteUser(ctx, s, u); err != nil {
					return
				}
			}
		}
	})
}

func (svc *user) onUploadAvatar(ctx context.Context, aProps *userActionProps, userID uint64, upload *multipart.FileHeader) error {
	u, err := loadUser(ctx, svc.store, userID)
	if err != nil {
		return err
	}

	if userID != internalAuth.GetIdentityFromContext(ctx).Identity() {
		if !svc.ac.CanUpdateUser(ctx, u) {
			return UserErrNotAllowedToUpdate()
		}
	}

	if u.Meta.AvatarID != 0 {
		if err = svc.services.att.DeleteByID(ctx, u.Meta.AvatarID); err != nil {
			return err
		}
	}

	file, err := upload.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	att, err := svc.services.att.CreateAuthAttachment(
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

	return store.UpdateUser(ctx, svc.store, u)
}

func (svc *user) onDeleteAvatar(ctx context.Context, aProps *userActionProps, userID uint64) error {
	u, err := svc.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if u.Kind == types.SystemUser {
		return UserErrNotAllowedToDeleteAvatar()
	}

	att, err := svc.services.att.FindByID(ctx, u.Meta.AvatarID)
	if err != nil {
		return err
	}

	if att.Meta.Labels["key"] != types.AttachmentKindAvatar {
		return nil
	}

	if !svc.ac.CanUpdateUser(ctx, u) {
		return UserErrNotAllowedToDeleteAvatar()
	}

	if err = svc.services.att.DeleteByID(ctx, u.Meta.AvatarID); err != nil {
		return err
	}

	u.Meta.AvatarID = 0

	if err = svc.generateUserAvatarInitial(ctx, u); err != nil {
		return err
	}

	return store.UpdateUser(ctx, svc.store, u)
}

func (svc *user) onGenerateAvatar(ctx context.Context, aProps *userActionProps, userID uint64, bgColor string, initialColor string) error {
	u, err := loadUser(ctx, svc.store, userID)
	if err != nil {
		return err
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

	return store.UpdateUser(ctx, svc.store, u)
}
