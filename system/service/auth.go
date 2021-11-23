package service

import (
	"context"
	"fmt"
	"math"
	rand2 "math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/dgryski/dgoogauth"
	"github.com/markbates/goth"
	"golang.org/x/crypto/bcrypt"
)

type (
	auth struct {
		actionlog actionlog.Recorder
		ac        authAccessController
		eventbus  eventDispatcher

		store         store.Storer
		settings      *types.AppSettings
		notifications AuthNotificationService

		opt AuthOptions

		providerValidator func(string) error
	}

	AuthOptions struct {
		LimitUsers int
	}

	authAccessController interface {
		CanImpersonateUser(context.Context, *types.User) bool
		CanUpdateUser(context.Context, *types.User) bool
	}
)

const (
	credentialsTypePassword                    = "password"
	credentialsTypePersistentSession           = "persistent-session"
	credentialsTypeEmailAuthToken              = "email-authentication-token"
	credentialsTypeResetPasswordToken          = "password-reset-token"
	credentialsTypeResetPasswordTokenExchanged = "password-reset-token-exchanged"
	credentialsTypeCreatePasswordToken         = "password-create-token"
	credentialsTypeMfaTotpSecret               = "mfa-totp-secret"
	credentialsTypeMFAEmailOTP                 = "mfa-email-otp"

	credentialsTokenLength = 32

	passwordMinLength = 8
	passwordMaxLength = 256
)

var (
	reEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func defaultProviderValidator(provider string) error {
	switch provider {
	case "saml":
		return nil
	default:
		_, err := goth.GetProvider(provider)
		return err
	}
}

func Auth(opt AuthOptions) *auth {
	return &auth{
		eventbus:      eventbus.Service(),
		ac:            DefaultAccessControl,
		settings:      CurrentSettings,
		notifications: DefaultAuthNotification,

		actionlog: DefaultActionlog,
		store:     DefaultStore,

		providerValidator: defaultProviderValidator,

		opt: opt,
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
func (svc auth) External(ctx context.Context, profile types.ExternalAuthUser) (u *types.User, err error) {
	var (
		authProvider = &types.AuthProvider{Provider: profile.Provider}

		aam = &authActionProps{
			email:    profile.Email,
			provider: profile.Provider,
			user:     u,
		}
	)

	err = func() (err error) {
		if !svc.settings.Auth.External.Enabled {
			return AuthErrExternalDisabledByConfig(aam)
		}

		if err = svc.providerValidator(profile.Provider); err != nil {
			return err
		}

		if !reEmail.MatchString(profile.Email) {
			return AuthErrProfileWithoutValidEmail(aam)
		}

		var (
			cc types.CredentialsSet
			f  = types.CredentialsFilter{Kind: profile.Provider, Credentials: profile.UserID}
		)

		if cc, _, err = store.SearchCredentials(ctx, svc.store, f); err == nil {
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
							return
						} else {
							goto findByEmail
						}
					}

					return
				}

				aam.setUser(u)
				ctx = internalAuth.SetIdentityToContext(ctx, u)

				if err = svc.procLogin(ctx, svc.store, u, c, authProvider); err != nil {
					// Scenario: linked to an invalid user
					if len(cc) > 1 {
						// try with next credentials
						u = nil
						continue
					}

					return AuthErrCredentialsLinkedToInvalidUser(aam)
				}

				// Assuming we can trust that email has been verified by the provider
				if !u.EmailConfirmed {
					u.EmailConfirmed = true
					if err = store.UpdateUser(ctx, svc.store, u); err != nil {
						return
					}
				}

				return svc.recordAction(ctx, aam, AuthActionUpdateCredentials, nil)
			}

			// If we could not find anything useful,
			// we can search for user via email
			// (using goto for consistency)
			goto findByEmail
		} else {
			// A serious error occurred, bail out...
			return
		}

	findByEmail:
		// Reset audit meta data that might got set during credentials check
		aam.setEmail(profile.Email).
			setCredentials(nil).
			setUser(nil)

		// Find user via his email
		if u, err = store.LookupUserByEmail(ctx, svc.store, profile.Email); errors.IsNotFound(err) {
			// @todo check if it is ok to auto-create a user here
			// In case we do not have this email, create a new user
			u = &types.User{
				Email:          profile.Email,
				Name:           profile.Name,
				Username:       profile.NickName,
				EmailConfirmed: true,
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

			if err = svc.checkLimits(ctx); err != nil {
				return err
			}

			u.ID = nextID()
			u.CreatedAt = *now()
			if err = store.CreateUser(ctx, svc.store, u); err != nil {
				return err
			}

			aam.setUser(u)
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
			ctx = internalAuth.SetIdentityToContext(ctx, u)

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

		// Assuming we can trust that email has been verified by the provider
		if !u.EmailConfirmed {
			u.EmailConfirmed = true
			if err = store.UpdateUser(ctx, svc.store, u); err != nil {
				return
			}
		}

		return nil
	}()

	return u, svc.recordAction(ctx, aam, AuthActionAuthenticate, err)
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
			ctx = internalAuth.SetIdentityToContext(ctx, u)

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

		if err = svc.checkLimits(ctx); err != nil {
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

		if err = svc.LoadRoleMemberships(ctx, nUser); err != nil {
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
			err = svc.SendEmailAddressConfirmationToken(ctx, nUser)
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
			return AuthErrInternalLoginDisabledByConfig()
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
		ctx = internalAuth.SetIdentityToContext(ctx, u)

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
			return AuthErrInternalLoginDisabledByConfig(aam)
		}

		if !svc.CheckPasswordStrength(password) {
			return AuthErrPasswordNotSecure(aam)
		}

		u, err = store.LookupUserByID(ctx, svc.store, userID)
		if errors.IsNotFound(err) {
			return AuthErrPasswordChangeFailedForUnknownUser(aam)
		}

		aam.setUser(u)
		ctx = internalAuth.SetIdentityToContext(ctx, u)

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
func (svc auth) ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) (err error) {
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
			return AuthErrInternalLoginDisabledByConfig(aam)
		}

		if len(oldPassword) == 0 {
			return AuthErrPasswordNotSecure(aam)
		}

		if !svc.CheckPasswordStrength(newPassword) {
			return AuthErrPasswordNotSecure(aam)
		}

		u, err = store.LookupUserByID(ctx, svc.store, userID)
		if errors.IsNotFound(err) {
			return AuthErrPasswordChangeFailedForUnknownUser(aam)
		}

		aam.setUser(u)
		ctx = internalAuth.SetIdentityToContext(ctx, u)

		cc, _, err = store.SearchCredentials(ctx, svc.store, types.CredentialsFilter{Kind: credentialsTypePassword, OwnerID: userID})
		if err != nil {
			return err
		}

		if !svc.checkPassword(oldPassword, cc) {
			return AuthErrPasswodResetFailedOldPasswordCheckFailed(aam)
		}

		if err != svc.SetPasswordCredentials(ctx, userID, newPassword) {
			return err
		}

		if err = svc.RemoveAccessTokens(ctx, u); err != nil {
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
	pwdL := len(password)

	// Ignore defined password constraints
	if !svc.settings.Auth.Internal.PasswordConstraints.PasswordSecurity {
		return true
	}

	// Check the password length
	minL := math.Max(float64(passwordMinLength), float64(svc.settings.Auth.Internal.PasswordConstraints.MinLength))
	if pwdL < int(minL) || pwdL > passwordMaxLength {
		return false
	}

	// Check special constraints
	// - numeric characters
	count := svc.settings.Auth.Internal.PasswordConstraints.MinNumCount
	if count > 0 {
		rr := regexp.MustCompile("[0-9]")
		if uint(len(rr.FindAllStringIndex(password, -1))) < count {
			return false
		}
	}

	// - special characters
	count = svc.settings.Auth.Internal.PasswordConstraints.MinSpecialCount
	if count > 0 {
		rr := regexp.MustCompile("[^0-9a-zA-Z]")
		if uint(len(rr.FindAllStringIndex(password, -1))) < count {
			return false
		}
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
	)

	if hash, err = svc.hashPassword(password); err != nil {
		return
	}

	if err = svc.removePasswordCredentials(ctx, userID); err != nil {
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

// RemovePasswordCredentials (soft) deletes old password entry
func (svc auth) RemovePasswordCredentials(ctx context.Context, userID uint64) (err error) {
	// Do a partial update and soft-delete all
	return svc.removePasswordCredentials(ctx, userID)
}

// RemovePasswordCredentials (soft) deletes old password entry
func (svc auth) removePasswordCredentials(ctx context.Context, userID uint64) (err error) {
	var (
		cc types.CredentialsSet
		f  = types.CredentialsFilter{Kind: credentialsTypePassword, OwnerID: userID}
	)

	if cc, _, err = store.SearchCredentials(ctx, svc.store, f); err != nil {
		return nil
	}

	// Mark all credentials as deleted
	_ = cc.Walk(func(c *types.Credentials) error {
		c.DeletedAt = now()
		return nil
	})

	// Do a partial update and soft-delete all
	return store.UpdateCredentials(ctx, svc.store, cc...)
}

// ValidateEmailConfirmationToken issues a validation token that can be used for
func (svc auth) ValidateEmailConfirmationToken(ctx context.Context, token string) (user *types.User, err error) {
	return svc.loadFromTokenAndConfirmEmail(ctx, token, credentialsTypeEmailAuthToken)
}

// ValidatePasswordResetToken validates password reset token
func (svc auth) ValidatePasswordResetToken(ctx context.Context, token string) (user *types.User, err error) {
	return svc.loadFromTokenAndConfirmEmail(ctx, token, credentialsTypeResetPasswordToken)
}

// ValidatePasswordCreateToken validates password create token
func (svc auth) ValidatePasswordCreateToken(ctx context.Context, token string) (user *types.User, err error) {
	return svc.loadFromTokenAndConfirmEmail(ctx, token, credentialsTypeCreatePasswordToken)
}

// PasswordSet checks and returns true if user's password is set
// False is also returned in case user does not exist.
func (svc *auth) PasswordSet(ctx context.Context, email string) (is bool) {

	//svc.settings.Auth.External.Enabled
	u, err := store.LookupUserByEmail(ctx, svc.store, email)
	if err != nil {
		return
	}

	cc, _, err := store.SearchCredentials(ctx, svc.store, types.CredentialsFilter{
		OwnerID: u.ID,
		Kind:    credentialsTypePassword,
	})
	if err != nil {
		return
	}

	if len(cc) > 0 && svc.settings.Auth.Internal.Enabled {
		return true
	}

	return
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

		if !u.EmailConfirmed {
			// User's email is not confirmed but going through password reset flow
			// we can confirm it
			u.EmailConfirmed = true
			u.UpdatedAt = now()
			if err = store.UpdateUser(ctx, svc.store, u); err != nil {
				return err
			}
		}

		if err = svc.LoadRoleMemberships(ctx, u); err != nil {
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

func (svc auth) SendEmailAddressConfirmationToken(ctx context.Context, u *types.User) (err error) {
	var (
		token string

		aam = &authActionProps{
			user:        u,
			credentials: &types.Credentials{Kind: credentialsTypeEmailAuthToken},
		}
	)

	if token, err = svc.createUserToken(ctx, u, credentialsTypeEmailAuthToken); err != nil {
		return
	}

	if err = svc.notifications.EmailConfirmation(ctx, u.Email, token); err != nil {
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

func (svc auth) sendPasswordResetToken(ctx context.Context, u *types.User) (err error) {
	token, err := svc.createUserToken(ctx, u, credentialsTypeResetPasswordToken)
	if err != nil {
		return err
	}

	return svc.notifications.PasswordReset(ctx, u.Email, token)
}

// GeneratePasswordCreateToken generates password create token
func (svc auth) GeneratePasswordCreateToken(ctx context.Context, email string) (url string, err error) {
	var (
		u *types.User

		aam = &authActionProps{
			user:  u,
			email: email,
		}
	)

	err = func() error {
		if !svc.settings.Auth.Internal.Enabled || !svc.settings.Auth.Internal.PasswordCreate.Enabled {
			return AuthErrPasswordCreateDisabledByConfig(aam)
		}

		if u, err = store.LookupUserByEmail(ctx, svc.store, email); err != nil {
			return err
		}

		ctx = internalAuth.SetIdentityToContext(ctx, u)

		if url, err = svc.sendPasswordCreateToken(ctx, u); err != nil {
			return err
		}

		return nil
	}()

	return url, svc.recordAction(ctx, aam, AuthActionGeneratePasswordCreateToken, err)
}

func (svc auth) sendPasswordCreateToken(ctx context.Context, u *types.User) (url string, err error) {
	token, err := svc.createUserToken(ctx, u, credentialsTypeCreatePasswordToken)
	if err != nil {
		return
	}

	return svc.notifications.PasswordCreate(token)
}

// procLogin fn performs standard validation, credentials-update tasks and triggers events
func (svc auth) procLogin(ctx context.Context, s store.Storer, u *types.User, c *types.Credentials, p *types.AuthProvider) (err error) {
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
	case u.Kind == types.SystemUser:
		return AuthErrFailedForSystemUser()

	}

	if err = svc.LoadRoleMemberships(ctx, u); err != nil {
		return
	}

	{
		var eapSec *types.ExternalAuthProviderSecurity

		switch p.Provider {
		case credentialsTypePassword:
			// nothing to do with password provider

		case "saml":
			// we need to fetch SAML security settings from different part of settings
			eapSec = &CurrentSettings.Auth.External.Saml.Security

		default:
			eap := CurrentSettings.Auth.External.Providers.FindByHandle(p.Provider)

			if eap != nil {
				eapSec = &eap.Security
			}

		}

		if eapSec != nil {
			// if authenticated with external auth provider
			// there might be additional roles that need to be
			// set to this security session
			u.SetRoles(internalAuth.ApplyRoleSecurity(
				payload.ParseUint64s(eapSec.PermittedRoles),
				payload.ParseUint64s(eapSec.ProhibitedRoles),
				payload.ParseUint64s(eapSec.ForcedRoles),
				u.Roles()...,
			)...)
		}
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
		case credentialsTypeMFAEmailOTP:
			expSec := svc.settings.Auth.MultiFactor.EmailOTP.Expires
			if expSec == 0 {
				expSec = 60
			}

			expiresAt = now().Add(time.Second * time.Duration(expSec))

			// random number, 6 chars
			token = fmt.Sprintf("%06d", rand2.Int())[0:6]
		case credentialsTypeCreatePasswordToken:
			expSec := svc.settings.Auth.Internal.PasswordCreate.Expires
			if expSec == 0 {
				expSec = 24
			}

			expiresAt = now().Add(time.Hour * time.Duration(expSec))

			// random password string, "3i[g0|)z"
			token = fmt.Sprintf("%s", rand.Password(credentialsTokenLength))
		default:
			// 1h expiration for all tokens send via email
			expiresAt = now().Add(time.Minute * 60)
			token = string(rand.Bytes(credentialsTokenLength))
		}

		c := &types.Credentials{
			ID:          nextID(),
			CreatedAt:   *now(),
			OwnerID:     u.ID,
			Kind:        kind,
			Credentials: token,
			ExpiresAt:   &expiresAt,
		}

		err = store.CreateCredentials(ctx, svc.store, c)

		if err != nil {
			return err
		}

		switch kind {
		case credentialsTypeMFAEmailOTP:
			// do not alter the final token
		default:
			// suffixing tokens with credentials ID
			// this will help us with token lookups
			token = fmt.Sprintf("%s%d", token, c.ID)
		}

		return nil
	}()

	return token, svc.recordAction(ctx, aam, AuthActionIssueToken, err)
}

// Automatically promotes user to super-administrator if it is the first non-system user in the database
func (svc auth) autoPromote(ctx context.Context, u *types.User) (err error) {
	var (
		c   uint
		aam = &authActionProps{user: u, role: &types.Role{}}
	)

	if c, err = countValidUsers(ctx, svc.store); err != nil {
		return err
	}

	if c > 1 || u.ID == 0 {
		return nil
	}

	for _, r := range internalAuth.BypassRoles() {
		m := &types.RoleMember{UserID: u.ID, RoleID: r.ID}
		if err = store.CreateRoleMember(ctx, svc.store, m); err != nil {
			return err
		}

		aam.role = r
		_ = svc.recordAction(ctx, aam, AuthActionAutoPromote, nil)
	}

	return nil
}

// ValidateTOTP checks given code with the current secret
// Fn fails if no secret is set
func (svc auth) ValidateTOTP(ctx context.Context, code string) (err error) {
	var (
		c    *types.Credentials
		u    *types.User
		kind = credentialsTypeMfaTotpSecret
		aam  = &authActionProps{credentials: &types.Credentials{Kind: kind}}
		i    = internalAuth.GetIdentityFromContext(ctx)
	)

	err = svc.store.Tx(ctx, func(ctx context.Context, s store.Storer) error {
		if !svc.settings.Auth.MultiFactor.TOTP.Enabled {
			return AuthErrDisabledMFAWithTOTP()
		}

		u, err = store.LookupUserByID(ctx, svc.store, i.Identity())
		if errors.IsNotFound(err) {
			return AuthErrFailedForUnknownUser(aam)
		}

		aam.setUser(u)

		if !u.Meta.SecurityPolicy.MFA.EnforcedTOTP {
			return AuthErrUnconfiguredTOTP()
		}

		if c, err = svc.getTOTPSecret(ctx, s, u.ID); err != nil {
			return err
		} else if err = svc.validateTOTP(c.Credentials, code); err != nil {
			return err
		} else {
			c.LastUsedAt = now()
			return store.UpdateCredentials(ctx, s, c)
		}
	})

	return svc.recordAction(ctx, aam, AuthActionTotpValidate, err)
}

// ConfigureTOTP stores totp secret in user's credentials
//
// It returns the user with security policy changes
func (svc auth) ConfigureTOTP(ctx context.Context, secret string, code string) (u *types.User, err error) {
	var (
		kind = credentialsTypeMfaTotpSecret
		aam  = &authActionProps{credentials: &types.Credentials{Kind: kind}}
		i    = internalAuth.GetIdentityFromContext(ctx)
	)

	err = svc.store.Tx(ctx, func(ctx context.Context, s store.Storer) error {
		if !svc.settings.Auth.MultiFactor.TOTP.Enabled {
			return AuthErrDisabledMFAWithTOTP()
		}

		if err = svc.validateTOTP(secret, code); err != nil {
			return err
		}

		u, err = store.LookupUserByID(ctx, svc.store, i.Identity())
		if errors.IsNotFound(err) {
			return AuthErrFailedForUnknownUser(aam)
		}

		aam.setUser(u)

		if i == nil || u.Meta.SecurityPolicy.MFA.EnforcedTOTP {
			// TOTP is already enforced on the user,
			// this means that we cannot just allow the change
			return AuthErrNotAllowedToConfigureTOTP()
		}

		// revoke (soft-delete) all existing secrets
		if err = svc.revokeAllTOTP(ctx, s, u.ID); err != nil {
			return err
		}

		cred := &types.Credentials{
			ID:          nextID(),
			CreatedAt:   *now(),
			OwnerID:     u.ID,
			Kind:        kind,
			Credentials: secret,
		}

		if err = store.CreateCredentials(ctx, s, cred); err != nil {
			return err
		}

		u.Meta.SecurityPolicy.MFA.EnforcedTOTP = true
		return store.UpdateUser(ctx, s, u)
	})

	return u, svc.recordAction(ctx, aam, AuthActionTotpConfigure, err)
}

// RemoveTOTP removes TOTP secret from user's credentials
//
// If user is removing own TOTP code is required
// When removing TOTP for another user, remover shou
//
// It returns the user with security policy changes
func (svc auth) RemoveTOTP(ctx context.Context, userID uint64, code string) (u *types.User, err error) {
	var (
		c    *types.Credentials
		kind = credentialsTypeMfaTotpSecret
		aam  = &authActionProps{credentials: &types.Credentials{Kind: kind}}
		i    = internalAuth.GetIdentityFromContext(ctx)
		self = i != nil && i.Identity() == userID
	)

	err = svc.store.Tx(ctx, func(ctx context.Context, s store.Storer) error {
		if !svc.settings.Auth.MultiFactor.TOTP.Enabled {
			return AuthErrDisabledMFAWithTOTP()
		}
		if svc.settings.Auth.MultiFactor.TOTP.Enforced {
			return AuthErrEnforcedMFAWithTOTP()
		}

		u, err = store.LookupUserByID(ctx, svc.store, userID)
		if errors.IsNotFound(err) {
			return AuthErrFailedForUnknownUser(aam)
		}

		aam.setUser(u)

		if i != nil && u != nil && self {
			if c, err = svc.getTOTPSecret(ctx, s, u.ID); err != nil {
				return err
			}

			if err = svc.validateTOTP(c.Credentials, code); err != nil {
				return err
			}
		} else if !svc.ac.CanUpdateUser(ctx, u) {
			return AuthErrNotAllowedToRemoveTOTP()
		}

		if err = svc.revokeAllTOTP(ctx, s, u.ID); err != nil {
			return err
		}

		u.Meta.SecurityPolicy.MFA.EnforcedTOTP = false
		return store.UpdateUser(ctx, s, u)

	})

	return u, svc.recordAction(ctx, aam, AuthActionTotpConfigure, err)
}

// Searches for all valid TOTP secret credentials
func (svc auth) getTOTPSecret(ctx context.Context, s store.Credentials, userID uint64) (*types.Credentials, error) {
	cc, _, err := store.SearchCredentials(ctx, s, types.CredentialsFilter{
		OwnerID: userID,
		Kind:    credentialsTypeMfaTotpSecret,
		Deleted: filter.StateExcluded,
	})

	if err != nil {
		return nil, err
	}

	if len(cc) != 1 {
		return nil, AuthErrInvalidTOTP()
	}

	return cc[0], nil
}

// Verifies TOTP code and secret
func (auth) validateTOTP(secret string, code string) error {
	// removes all non-numeric characters
	code = regexp.MustCompile(`[^0-9]`).ReplaceAllString(code, "")
	if len(code) != 6 {
		return AuthErrInvalidTOTP()
	}

	otpc := &dgoogauth.OTPConfig{
		Secret:     secret,
		WindowSize: 5,
	}

	if ok, err := otpc.Authenticate(code); err != nil {
		return AuthErrInvalidTOTP().Wrap(err)
	} else if !ok {
		return AuthErrInvalidTOTP()
	}

	return nil
}

// Revokes all existing user's TOTPs
func (auth) revokeAllTOTP(ctx context.Context, s store.Credentials, userID uint64) error {
	// revoke (soft-delete) all existing secrets
	cc, _, err := store.SearchCredentials(ctx, s, types.CredentialsFilter{
		OwnerID: userID,
		Kind:    credentialsTypeMfaTotpSecret,
		Deleted: filter.StateExcluded,
	})

	if err != nil {
		return err
	}

	return cc.Walk(func(c *types.Credentials) error {
		c.DeletedAt = now()
		return store.UpdateCredentials(ctx, s, c)
	})
}

func (svc auth) SendEmailOTP(ctx context.Context) (err error) {
	var (
		otp  string
		u    *types.User
		kind = credentialsTypeMFAEmailOTP
		aam  = &authActionProps{credentials: &types.Credentials{Kind: kind}}
		i    = internalAuth.GetIdentityFromContext(ctx)
	)

	err = svc.store.Tx(ctx, func(ctx context.Context, s store.Storer) (err error) {
		if !svc.settings.Auth.MultiFactor.EmailOTP.Enabled {
			return AuthErrDisabledMFAWithEmailOTP()
		}

		u, err = store.LookupUserByID(ctx, svc.store, i.Identity())
		if errors.IsNotFound(err) {
			return AuthErrFailedForUnknownUser(aam)
		}

		aam.setUser(u)

		if otp, err = svc.createUserToken(ctx, u, kind); err != nil {
			return
		}

		if err = svc.notifications.EmailOTP(ctx, u.Email, otp); err != nil {
			return
		}

		return
	})

	return svc.recordAction(ctx, aam, AuthActionSendEmailConfirmationToken, err)
}

func (svc auth) ConfigureEmailOTP(ctx context.Context, userID uint64, enable bool) (u *types.User, err error) {
	var (
		kind = credentialsTypeMFAEmailOTP
		aam  = &authActionProps{credentials: &types.Credentials{Kind: kind}}
	)

	err = svc.store.Tx(ctx, func(ctx context.Context, s store.Storer) (err error) {
		if !svc.settings.Auth.MultiFactor.EmailOTP.Enabled {
			return AuthErrDisabledMFAWithEmailOTP()
		}

		if svc.settings.Auth.MultiFactor.EmailOTP.Enforced && !enable {
			return AuthErrEnforcedMFAWithEmailOTP()
		}

		u, err = store.LookupUserByID(ctx, svc.store, userID)
		if errors.IsNotFound(err) {
			return AuthErrFailedForUnknownUser(aam)
		}

		aam.setUser(u)
		u.Meta.SecurityPolicy.MFA.EnforcedEmailOTP = enable

		return store.UpdateUser(ctx, s, u)
	})

	return u, svc.recordAction(ctx, aam, AuthActionSendEmailConfirmationToken, err)
}

// ValidateEmailOTP issues a validation OTP
func (svc auth) ValidateEmailOTP(ctx context.Context, code string) (err error) {
	var (
		cc   types.CredentialsSet
		u    *types.User
		kind = credentialsTypeMFAEmailOTP
		aam  = &authActionProps{credentials: &types.Credentials{Kind: kind}}
		i    = internalAuth.GetIdentityFromContext(ctx)
	)

	err = svc.store.Tx(ctx, func(ctx context.Context, s store.Storer) error {
		if !svc.settings.Auth.MultiFactor.EmailOTP.Enabled {
			return AuthErrDisabledMFAWithEmailOTP()
		}

		u, err = store.LookupUserByID(ctx, svc.store, i.Identity())
		if errors.IsNotFound(err) {
			return AuthErrFailedForUnknownUser(aam)
		}

		aam.setUser(u)

		// removes all non-numeric characters
		code = regexp.MustCompile(`[^0-9]`).ReplaceAllString(code, "")
		if len(code) != 6 {
			return AuthErrInvalidEmailOTP()
		}

		cc, _, err = store.SearchCredentials(ctx, s, types.CredentialsFilter{
			OwnerID: u.ID,
			Kind:    kind,
			Deleted: filter.StateExcluded,
		})

		if err != nil {
			return err
		}

		for _, c := range cc {
			if c.ExpiresAt.Before(*now()) {
				continue
			}

			if c.Credentials != code {
				continue
			}

			// Credentials found, remove it
			return store.DeleteCredentials(ctx, s, c)
		}

		return AuthErrInvalidEmailOTP()
	})

	return svc.recordAction(ctx, aam, AuthActionEmailOtpVerify, err)
}

// LoadRoleMemberships loads membership info
//
// @todo move this to role service
func (svc auth) LoadRoleMemberships(ctx context.Context, u *types.User) error {
	rr, _, err := store.SearchRoles(ctx, svc.store, types.RoleFilter{MemberID: u.ID})
	if err != nil {
		return err
	}

	u.SetRoles(rr.IDs()...)
	return nil
}

func (svc auth) GetProviders() types.ExternalAuthProviderSet {
	return CurrentSettings.Auth.External.Providers
}

func (svc auth) checkLimits(ctx context.Context) error {
	if svc.opt.LimitUsers == 0 {
		return nil
	}

	if c, err := countValidUsers(ctx, svc.store); err != nil {
		return err
	} else if c >= uint(svc.opt.LimitUsers) {
		return AuthErrMaxUserLimitReached()
	}

	return nil
}

// RemoveAccessTokens removes all user's access tokens when suspended,
// deleted or security context changes
func (svc auth) RemoveAccessTokens(ctx context.Context, user *types.User) error {
	return svc.recordAction(
		ctx,
		&authActionProps{user: user},
		AuthActionAccessTokensRemoved,
		svc.store.DeleteAuthOA2TokenByUserID(ctx, user.ID),
	)
}
