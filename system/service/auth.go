package service

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/markbates/goth"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strconv"
	"time"
)

type (
	auth struct {
		actionlog actionlog.Recorder
		ac        authAccessController
		eventbus  eventDispatcher

		subscription  authSubscriptionChecker
		store         store.Storer
		settings      *types.AppSettings
		notifications AuthNotificationService

		providerValidator func(string) error
	}

	authAccessController interface {
		CanImpersonateUser(context.Context, *types.User) bool
	}

	authSubscriptionChecker interface {
		CanRegister(uint) error
	}
)

const (
	credentialsTypePassword                    = "password"
	credentialsTypeEmailAuthToken              = "email-authentication-token"
	credentialsTypeResetPasswordToken          = "password-reset-token"
	credentialsTypeResetPasswordTokenExchanged = "password-reset-token-exchanged"
	credentialsTypeAuthToken                   = "auth-token"

	credentialsTokenLength = 32
)

var (
	reEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func defaultProviderValidator(provider string) error {
	_, err := goth.GetProvider(provider)
	return err
}

func Auth() *auth {
	return &auth{
		eventbus:      eventbus.Service(),
		ac:            DefaultAccessControl,
		subscription:  CurrentSubscription,
		settings:      CurrentSettings,
		notifications: DefaultAuthNotification,

		actionlog: DefaultActionlog,
		store:     DefaultStore,

		providerValidator: defaultProviderValidator,
	}
}

// External func performs login/signup procedures
//
// We fully trust external auth sources (see system/auth/external) to provide a valid & validates
// profile (goth.User) that we use for user discovery and/or creation
//
// Flow
// 1.   check for existing credentials using profile provider & provider's user ID
// 1.1. find existing local -or- "shadow" user
// 1.2. if user exists and is valid, update credentials (last-used-at) and complete the procedure
//
// 2.   check for existing users using email from the profile
// 2.1. validate existing user -or-
// 2.2. create user on-the-fly if it does not exist
// 2.3. create credentials for that social login
//
// External login/signup does not:
//  - validate provider on profile, only uses it for matching credentials
func (svc auth) External(ctx context.Context, profile goth.User) (u *types.User, err error) {
	var (
		authProvider = &types.AuthProvider{Provider: profile.Provider}

		aam = &authActionProps{
			email:    profile.Email,
			provider: profile.Provider,
			user:     u,
		}
	)

	err = func() error {
		if !svc.settings.Auth.External.Enabled {
			return AuthErrExternalDisabledByConfig(aam)
		}

		if err = svc.providerValidator(profile.Provider); err != nil {
			return err
		}

		if !reEmail.MatchString(profile.Email) {
			return AuthErrProfileWithoutValidEmail(aam)
		}

		f := types.CredentialsFilter{Kind: profile.Provider, Credentials: profile.UserID}
		if cc, _, err := store.SearchCredentials(ctx, svc.store, f); err == nil {
			// Credentials found, load user
			for _, c := range cc {
				if !c.Valid() {
					continue
				}

				// Add credentials ID for audit log
				aam.setCredentials(c)

				if u, err = store.LookupUserByID(ctx, svc.store, c.OwnerID); err != nil {
					if errors.IsNotFound(err) {
						// Orphaned credentials (no owner)
						// try to auto-fix this by removing credentials and recreating user
						if err = store.DeleteCredentialsByID(ctx, svc.store, c.ID); err != nil {
							return err
						} else {
							goto findByEmail
						}
					}

					return err
				}

				aam.setUser(u)

				if err = svc.procLogin(ctx, svc.store, u, c, authProvider); err != nil {
					// Scenario: linked to an invalid user
					if len(cc) > 1 {
						// try with next credentials
						u = nil
						continue
					}

					return AuthErrCredentialsLinkedToInvalidUser(aam)
				}

				return svc.recordAction(ctx, aam, AuthActionUpdateCredentials, nil)
			}

			// If we could not find anything useful,
			// we can search for user via email
			// (using goto for consistency)
			goto findByEmail
		} else {
			// A serious error occurred, bail out...
			return err
		}

	findByEmail:
		// Reset audit meta data that might got workflows during credentials check
		aam.setEmail(profile.Email).
			setCredentials(nil).
			setUser(nil)

		// Find user via his email
		if u, err = store.LookupUserByEmail(ctx, svc.store, profile.Email); errors.IsNotFound(err) {
			// @todo check if it is ok to auto-create a user here
			if err = svc.CanRegister(ctx); err != nil {
				return AuthErrSubscription(aam).Wrap(err)
			}

			// In case we do not have this email, create a new user
			u = &types.User{
				Email:    profile.Email,
				Name:     profile.Name,
				Username: profile.NickName,
			}

			if !handle.IsValid(profile.NickName) {
				u.Handle = profile.NickName
			}

			if err = svc.eventbus.WaitFor(ctx, event.AuthBeforeSignup(u, authProvider)); err != nil {
				return err
			}

			if u.Handle == "" {
				createUserHandle(ctx, svc.store, u)
			}

			if err = uniqueUserCheck(ctx, svc.store, u); err != nil {
				return err
			}

			u.ID = nextID()
			u.CreatedAt = *now()
			if err = store.CreateUser(ctx, svc.store, u); err != nil {
				return err
			}

			aam.setUser(nil)
			ctx = internalAuth.SetIdentityToContext(ctx, u)

			_ = svc.eventbus.WaitFor(ctx, event.AuthAfterSignup(u, authProvider))

			if err = svc.recordAction(ctx, aam, AuthActionExternalSignup, nil); err != nil {
				return err
			}

			// Auto-promote first user
			if err = svc.autoPromote(ctx, u); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			// User found
			aam.setUser(u)

			if err = svc.procLogin(ctx, svc.store, u, nil, authProvider); err != nil {
				return err
			}
		}

		// If we got to this point, assume that user is authenticated
		// but credentials need to be stored
		c := &types.Credentials{
			ID:          nextID(),
			CreatedAt:   *now(),
			Kind:        profile.Provider,
			OwnerID:     u.ID,
			Credentials: profile.UserID,
			LastUsedAt:  now(),
		}

		if err = store.CreateCredentials(ctx, svc.store, c); err != nil {
			return err
		}

		aam.setCredentials(c)
		svc.recordAction(ctx, aam, AuthActionCreateCredentials, nil)

		// Owner loaded, carry on.
		return nil
	}()

	return u, svc.recordAction(ctx, aam, AuthActionAuthenticate, err)
}

// FrontendRedirectURL - a proxy to frontend redirect url setting
func (svc auth) FrontendRedirectURL() string {
	return svc.settings.Auth.Frontend.Url.Redirect
}

// InternalSignUp protocol
//
// Forgiving but strict: valid existing users get notified
//
// We're accepting the whole user object here and copy all we need to the new user
func (svc auth) InternalSignUp(ctx context.Context, input *types.User, password string) (u *types.User, err error) {
	var (
		authProvider = &types.AuthProvider{Provider: credentialsTypePassword}

		aam = &authActionProps{
			email:       input.Email,
			credentials: &types.Credentials{Kind: credentialsTypePassword},
		}
	)

	err = func() error {
		if !svc.settings.Auth.Internal.Enabled || !svc.settings.Auth.Internal.Signup.Enabled {
			return AuthErrInternalSignupDisabledByConfig(aam)
		}

		if input == nil || !reEmail.MatchString(input.Email) {
			return AuthErrInvalidEmailFormat(aam)
		}

		if !handle.IsValid(input.Handle) {
			return AuthErrInvalidHandle(aam)
		}

		if len(password) == 0 {
			return AuthErrPasswordNotSecure(aam)
		}

		var eUser *types.User
		eUser, err = store.LookupUserByEmail(ctx, svc.store, input.Email)

		if err == nil && eUser != nil {
			var (
				c  *types.Credentials
				cc types.CredentialsSet
				f  = types.CredentialsFilter{OwnerID: eUser.ID, Kind: credentialsTypePassword}
			)
			if cc, _, err = store.SearchCredentials(ctx, svc.store, f); err != nil {
				return err
			}

			if c = cc.CompareHashAndPassword(password); c == nil {
				return AuthErrInvalidCredentials(aam)
			}

			aam.setCredentials(c)
			u = eUser
			return svc.procLogin(ctx, svc.store, eUser, c, authProvider)
		} else if !errors.IsNotFound(err) {
			return err
		}

		// if !svc.settings.internalSignUpSendEmailOnExisting {
		// 	return nil,errors.Wrap(err, "user with this email already exists")
		// }

		// User already exists, but we're nice and we'll send this user an
		// email that will help him to login
		// if !u.Valid() {
		// 	return nil,errors.New("could not validate the user")
		// }
		//
		// return nil,nil

		if err = svc.CanRegister(ctx); err != nil {
			return err
		}

		var nUser = &types.User{
			ID:        nextID(),
			CreatedAt: *now(),

			Email:    input.Email,
			Name:     input.Name,
			Username: input.Username,
			Handle:   input.Handle,

			// Do we need confirmed email?
			EmailConfirmed: !svc.settings.Auth.Internal.Signup.EmailConfirmationRequired,
		}

		if err = svc.eventbus.WaitFor(ctx, event.AuthBeforeSignup(nUser, authProvider)); err != nil {
			return err
		}

		if nUser.Handle == "" {
			createUserHandle(ctx, svc.store, nUser)
		}

		if err = uniqueUserCheck(ctx, svc.store, nUser); err != nil {
			return err
		}

		// Whitelisted user data to copy
		err = store.CreateUser(ctx, svc.store, nUser)
		if err != nil {
			return err
		}

		aam.setUser(nUser)
		_ = svc.eventbus.WaitFor(ctx, event.AuthAfterSignup(nUser, authProvider))

		if err = svc.autoPromote(ctx, nUser); err != nil {
			return err
		}

		if len(password) > 0 {
			err = svc.SetPassword(ctx, nUser.ID, password)
			if err != nil {
				return err
			}
		}

		u = nUser
		if !nUser.EmailConfirmed {
			err = svc.sendEmailAddressConfirmationToken(ctx, nUser)
			if err != nil {
				return err
			}

			return svc.recordAction(ctx, aam, AuthActionSendEmailConfirmationToken, nil)
		}

		return nil
	}()

	return u, svc.recordAction(ctx, aam, AuthActionInternalSignup, err)
}

// InternalLogin verifies username/password combination in the internal credentials table
//
// Expects plain text password as an input
func (svc auth) InternalLogin(ctx context.Context, email string, password string) (u *types.User, err error) {
	var (
		authProvider = &types.AuthProvider{Provider: credentialsTypePassword}

		aam = &authActionProps{
			email:       email,
			credentials: &types.Credentials{Kind: credentialsTypePassword},
			user:        u,
		}
	)

	err = func() error {
		if !svc.settings.Auth.Internal.Enabled {
			return AuthErrInteralLoginDisabledByConfig()
		}

		if !reEmail.MatchString(email) {
			return AuthErrInvalidEmailFormat()
		}

		if len(password) == 0 {
			return AuthErrInvalidCredentials()
		}

		var (
			cc types.CredentialsSet
		)

		u, err = store.LookupUserByEmail(ctx, svc.store, email)
		if errors.IsNotFound(err) {
			return AuthErrInvalidCredentials(aam)
		} else if err != nil {
			return err
		}

		// Update audit meta with found user
		ctx = internalAuth.SetIdentityToContext(ctx, u)

		cc, _, err = store.SearchCredentials(ctx, svc.store, types.CredentialsFilter{OwnerID: u.ID, Kind: credentialsTypePassword})
		if err != nil {
			return err
		}

		c := cc.CompareHashAndPassword(password)
		if c == nil {
			return AuthErrInvalidCredentials(aam)
		}

		aam.setCredentials(c)
		return svc.procLogin(ctx, svc.store, u, c, authProvider)
	}()

	return u, svc.recordAction(ctx, aam, AuthActionAuthenticate, err)
}

// checkPassword returns true if given (encrypted) password matches any of the credentials
func (svc auth) checkPassword(password string, cc types.CredentialsSet) bool {
	return cc.CompareHashAndPassword(password) != nil
}

// SetPassword sets new password for a user
//
// This function also records an action
func (svc auth) SetPassword(ctx context.Context, userID uint64, password string) (err error) {
	var (
		u *types.User

		aam = &authActionProps{
			user:        u,
			credentials: &types.Credentials{Kind: credentialsTypePassword},
		}
	)

	err = func() error {
		if !svc.settings.Auth.Internal.Enabled {
			return AuthErrInteralLoginDisabledByConfig(aam)
		}

		if !svc.CheckPasswordStrength(password) {
			return AuthErrPasswordNotSecure(aam)
		}

		u, err = store.LookupUserByID(ctx, svc.store, userID)
		if errors.IsNotFound(err) {
			return AuthErrPasswordChangeFailedForUnknownUser(aam)
		}

		if err != svc.SetPasswordCredentials(ctx, userID, password) {
			return err
		}

		return nil
	}()

	return svc.recordAction(ctx, aam, AuthActionChangePassword, err)
}

// Impersonate verifies if user can impersonate another user and returns that user
//
// For now, it's the caller's responsibility to generate the auth token
func (svc auth) Impersonate(ctx context.Context, userID uint64) (u *types.User, err error) {
	var (
		aam = &authActionProps{user: u}
	)

	err = func() error {
		if u, err = store.LookupUserByID(ctx, svc.store, userID); err != nil {
			return err
		}

		if !svc.ac.CanImpersonateUser(ctx, u) {
			return AuthErrNotAllowedToImpersonate()
		}

		return err
	}()

	return u, svc.recordAction(ctx, aam, AuthActionImpersonate, err)
}

// ChangePassword validates old password and changes it with new
func (svc auth) ChangePassword(ctx context.Context, userID uint64, oldPassword, AuthActionPassword string) (err error) {
	var (
		u  *types.User
		cc types.CredentialsSet

		aam = &authActionProps{
			user:        u,
			credentials: &types.Credentials{Kind: credentialsTypePassword},
		}
	)

	err = func() error {
		if !svc.settings.Auth.Internal.Enabled {
			return AuthErrInteralLoginDisabledByConfig(aam)
		}

		if len(oldPassword) == 0 {
			return AuthErrPasswordNotSecure(aam)
		}

		if !svc.CheckPasswordStrength(AuthActionPassword) {
			return AuthErrPasswordNotSecure(aam)
		}

		u, err = store.LookupUserByID(ctx, svc.store, userID)
		if errors.IsNotFound(err) {
			return AuthErrPasswordChangeFailedForUnknownUser(aam)
		}

		cc, _, err = store.SearchCredentials(ctx, svc.store, types.CredentialsFilter{Kind: credentialsTypePassword, OwnerID: userID})
		if err != nil {
			return err
		}

		if !svc.checkPassword(oldPassword, cc) {
			return AuthErrPasswodResetFailedOldPasswordCheckFailed(aam)
		}

		if err != svc.SetPasswordCredentials(ctx, userID, AuthActionPassword) {
			return err
		}

		return nil
	}()

	return svc.recordAction(ctx, aam, AuthActionChangePassword, err)
}

func (svc auth) hashPassword(password string) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (svc auth) CheckPasswordStrength(password string) bool {
	if len(password) <= 4 {
		return false
	}

	return true
}

// SetPasswordCredentials (soft) deletes old password entry and creates a new entry with new password on every change
//
// This way we can implement more strict password-change policies in the future
//
// This method is used by auth and user procedures to unify password hashing and updating
// credentials
func (svc auth) SetPasswordCredentials(ctx context.Context, userID uint64, password string) (err error) {
	var (
		hash []byte
		cc   types.CredentialsSet
		f    = types.CredentialsFilter{Kind: credentialsTypePassword, OwnerID: userID}
	)

	if hash, err = svc.hashPassword(password); err != nil {
		return
	}

	if cc, _, err = store.SearchCredentials(ctx, svc.store, f); err != nil {
		return nil
	}

	// Mark all credentials as deleted
	_ = cc.Walk(func(c *types.Credentials) error {
		c.DeletedAt = now()
		return nil
	})

	// Do a partial update and soft-delete all
	if err = store.UpdateCredentials(ctx, svc.store, cc...); err != nil {
		return
	}

	// Add new credentials with new password
	c := &types.Credentials{
		ID:          nextID(),
		CreatedAt:   *now(),
		OwnerID:     userID,
		Kind:        credentialsTypePassword,
		Credentials: string(hash),
	}

	return store.CreateCredentials(ctx, svc.store, c)
}

// IssueAuthRequestToken returns token that can be used for authentication
func (svc auth) IssueAuthRequestToken(ctx context.Context, user *types.User) (token string, err error) {
	return svc.createUserToken(ctx, user, credentialsTypeAuthToken)
}

// ValidateAuthRequestToken returns user that requested auth token
func (svc auth) ValidateAuthRequestToken(ctx context.Context, token string) (u *types.User, err error) {
	var (
		aam = &authActionProps{
			credentials: &types.Credentials{Kind: credentialsTypeAuthToken},
		}
	)

	err = func() error {
		u, err = svc.loadUserFromToken(ctx, token, credentialsTypeAuthToken)
		if err != nil && u != nil {
			aam.setUser(u)
			ctx = internalAuth.SetIdentityToContext(ctx, u)
		}
		return err
	}()

	return u, svc.recordAction(ctx, aam, AuthActionValidateToken, err)
}

// ValidateEmailConfirmationToken issues a validation token that can be used for
func (svc auth) ValidateEmailConfirmationToken(ctx context.Context, token string) (user *types.User, err error) {
	return svc.loadFromTokenAndConfirmEmail(ctx, token, credentialsTypeEmailAuthToken)
}

// ValidatePasswordResetToken validates password reset token
func (svc auth) ValidatePasswordResetToken(ctx context.Context, token string) (user *types.User, err error) {
	return svc.loadFromTokenAndConfirmEmail(ctx, token, credentialsTypeEmailAuthToken)
}

// loadFromTokenAndConfirmEmail loads token, confirms user's
func (svc auth) loadFromTokenAndConfirmEmail(ctx context.Context, token, tokenType string) (u *types.User, err error) {
	var (
		aam = &authActionProps{
			user:        u,
			credentials: &types.Credentials{Kind: tokenType},
		}
	)

	err = func() error {
		if !svc.settings.Auth.Internal.Enabled {
			return AuthErrInternalSignupDisabledByConfig(aam)
		}

		u, err = svc.loadUserFromToken(ctx, token, tokenType)
		if err != nil {
			return err
		}

		aam.setUser(u)
		ctx = internalAuth.SetIdentityToContext(ctx, u)

		if u.EmailConfirmed {
			return nil
		}

		u.EmailConfirmed = true
		u.UpdatedAt = now()
		if err = store.UpdateUser(ctx, svc.store, u); err != nil {
			return err
		}

		return nil
	}()

	return u, svc.recordAction(ctx, aam, AuthActionConfirmEmail, err)
}

// ExchangePasswordResetToken exchanges reset password token for a new one and returns it with user info
func (svc auth) ExchangePasswordResetToken(ctx context.Context, token string) (u *types.User, t string, err error) {
	var (
		aam = &authActionProps{
			user:        u,
			credentials: &types.Credentials{Kind: credentialsTypeResetPasswordToken},
		}
	)

	err = func() error {
		if !svc.settings.Auth.Internal.Enabled || !svc.settings.Auth.Internal.PasswordReset.Enabled {
			return AuthErrPasswordResetDisabledByConfig(aam)
		}

		u, err = svc.loadUserFromToken(ctx, token, credentialsTypeResetPasswordToken)
		if err != nil {
			return AuthErrInvalidToken(aam).Wrap(err)
		}

		aam.setUser(u)
		ctx = internalAuth.SetIdentityToContext(ctx, u)

		t, err = svc.createUserToken(ctx, u, credentialsTypeResetPasswordTokenExchanged)
		if err != nil {
			u = nil
			t = ""
			return AuthErrInvalidToken(aam).Wrap(err)
		}

		return nil
	}()

	return u, t, svc.recordAction(ctx, aam, AuthActionExchangePasswordResetToken, err)
}

func (svc auth) sendEmailAddressConfirmationToken(ctx context.Context, u *types.User) (err error) {
	var (
		notificationLang = "en"
		token            string

		aam = &authActionProps{
			user:        u,
			credentials: &types.Credentials{Kind: credentialsTypeEmailAuthToken},
		}
	)

	if token, err = svc.createUserToken(ctx, u, credentialsTypeEmailAuthToken); err != nil {
		return
	}

	if err = svc.notifications.EmailConfirmation(ctx, notificationLang, u.Email, token); err != nil {
		return
	}

	return svc.recordAction(ctx, aam, AuthActionSendEmailConfirmationToken, err)
}

// SendPasswordResetToken sends password reset token to email
func (svc auth) SendPasswordResetToken(ctx context.Context, email string) (err error) {
	var (
		u *types.User

		aam = &authActionProps{
			user:  u,
			email: email,
		}
	)

	err = func() error {
		if !svc.settings.Auth.Internal.Enabled || !svc.settings.Auth.Internal.PasswordReset.Enabled {
			return AuthErrPasswordResetDisabledByConfig(aam)
		}

		if u, err = store.LookupUserByEmail(ctx, svc.store, email); err != nil {
			return err
		}

		ctx = internalAuth.SetIdentityToContext(ctx, u)

		if err = svc.sendPasswordResetToken(ctx, u); err != nil {
			return err
		}

		return nil
	}()

	return svc.recordAction(ctx, aam, AuthActionSendPasswordResetToken, err)
}

// CanRegister verifies if user can register
func (svc auth) CanRegister(ctx context.Context) error {
	if svc.subscription != nil {
		c, err := store.CountUsers(ctx, svc.store, types.UserFilter{})
		if err != nil {
			return fmt.Errorf("can not check if user can register: %w", err)
		}

		// When we have an active subscription, we need to check
		// if users can register or did this deployment hit
		// it's user-limit
		return svc.subscription.CanRegister(c)
	}

	return nil
}

func (svc auth) sendPasswordResetToken(ctx context.Context, u *types.User) (err error) {
	var (
		notificationLang = "en"
	)

	token, err := svc.createUserToken(ctx, u, credentialsTypeResetPasswordToken)
	if err != nil {
		return err
	}

	return svc.notifications.PasswordReset(ctx, notificationLang, u.Email, token)
}

// procLogin fn performs standard validation, credentials-update tasks and triggers events
func (svc auth) procLogin(ctx context.Context, s store.Storer, u *types.User, c *types.Credentials, p *types.AuthProvider) (err error) {
	ctx = internalAuth.SetIdentityToContext(ctx, u)
	if err = svc.eventbus.WaitFor(ctx, event.AuthBeforeLogin(u, p)); err != nil {
		return err
	}

	// all checks (suspension, deleted, confirmed email) are checked AFTER
	// before-login event to enable before-login hooks to alter user and make her
	// valid for login
	switch true {
	case u.SuspendedAt != nil:
		return AuthErrFailedForSuspendedUser()
	case u.DeletedAt != nil:
		return AuthErrFailedForDeletedUser()
	case !u.EmailConfirmed && svc.settings.Auth.Internal.Signup.EmailConfirmationRequired:
		// Re-send email-confirmation when not confirmed and signup email confirmation required
		if err = svc.sendEmailAddressConfirmationToken(ctx, u); err != nil {
			return err
		}

		return AuthErrFailedUnconfirmedEmail()
	}

	if c != nil {
		c.LastUsedAt = now()
		if err = store.UpdateCredentials(ctx, s, c); err != nil {
			return err
		}
	}

	_ = svc.eventbus.WaitFor(ctx, event.AuthAfterLogin(u, p))
	return nil
}

func (svc auth) loadUserFromToken(ctx context.Context, token, kind string) (u *types.User, err error) {
	var (
		aam = &authActionProps{
			credentials: &types.Credentials{Kind: kind},
		}
	)

	credentialsID, credentials := svc.validateToken(token)
	if credentialsID == 0 {
		return nil, AuthErrInvalidToken(aam)
	}

	c, err := store.LookupCredentialsByID(ctx, svc.store, credentialsID)
	if errors.IsNotFound(err) {
		return nil, AuthErrInvalidToken(aam)
	}

	aam.setCredentials(c)

	if err != nil {
		return
	}

	if err = store.DeleteCredentialsByID(ctx, svc.store, c.ID); err != nil {
		return
	}

	if !c.Valid() || c.Credentials != credentials {
		return nil, AuthErrInvalidToken(aam)
	}

	u, err = store.LookupUserByID(ctx, svc.store, c.OwnerID)
	if err != nil {
		return nil, err
	}

	aam.setUser(u)

	// context will be updated with new identity
	// in the caller fn

	if !u.Valid() {
		return nil, AuthErrInvalidCredentials(aam)
	}

	return u, nil
}

func (svc auth) validateToken(token string) (ID uint64, credentials string) {
	// Token = <32 random chars><credentials-id>
	if len(token) <= credentialsTokenLength {
		return
	}

	ID, _ = strconv.ParseUint(token[credentialsTokenLength:], 10, 64)
	if ID == 0 {
		return
	}

	credentials = token[:credentialsTokenLength]
	return
}

// Generates & stores user token
// it returns combined value of token + token ID to help with the lookups
func (svc auth) createUserToken(ctx context.Context, u *types.User, kind string) (token string, err error) {
	var (
		expiresAt time.Time
		aam       = &authActionProps{
			user:        u,
			credentials: &types.Credentials{Kind: kind},
		}
	)

	err = func() error {
		switch kind {
		case credentialsTypeAuthToken:
			// 15 sec expiration for all tokens that are part of redirection
			expiresAt = now().Add(time.Second * 15)
		default:
			// 1h expiration for all tokens send via email
			expiresAt = now().Add(time.Minute * 60)
		}

		c := &types.Credentials{
			ID:          nextID(),
			CreatedAt:   *now(),
			OwnerID:     u.ID,
			Kind:        kind,
			Credentials: string(rand.Bytes(credentialsTokenLength)),
			ExpiresAt:   &expiresAt,
		}

		err := store.CreateCredentials(ctx, svc.store, c)

		if err != nil {
			return err
		}

		token = fmt.Sprintf("%s%d", c.Credentials, c.ID)
		return nil
	}()

	return token, svc.recordAction(ctx, aam, AuthActionIssueToken, err)
}

// Automatically promotes user to administrator if it is the first user in the database
func (svc auth) autoPromote(ctx context.Context, u *types.User) (err error) {
	var (
		c      uint
		roleID = rbac.AdminsRoleID
		aam    = &authActionProps{user: u, role: &types.Role{ID: roleID}}
	)

	err = func() error {
		if c, err = store.CountUsers(ctx, svc.store, types.UserFilter{}); err != nil {
			return err
		}

		if c > 1 || u.ID == 0 {
			return nil
		}

		return store.CreateRoleMember(ctx, svc.store, &types.RoleMember{RoleID: roleID, UserID: u.ID})
	}()

	return svc.recordAction(ctx, aam, AuthActionAutoPromote, err)
}

// LoadRoleMemberships loads membership info
//
// @todo move this to role service
func (svc auth) LoadRoleMemberships(ctx context.Context, u *types.User) error {
	rr, _, err := store.SearchRoles(ctx, svc.store, types.RoleFilter{MemberID: u.ID})
	if err != nil {
		return err
	}

	u.SetRoles(rr.IDs())
	return nil
}
