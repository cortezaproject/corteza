package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/markbates/goth"
	"golang.org/x/crypto/bcrypt"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	auth struct {
		db  db
		ctx context.Context

		actionlog actionlog.Recorder
		eventbus  eventDispatcher

		subscription  authSubscriptionChecker
		credentials   repository.CredentialsRepository
		users         repository.UserRepository
		roles         repository.RoleRepository
		settings      *types.Settings
		notifications AuthNotificationService

		providerValidator func(string) error
		now               func() *time.Time
	}

	AuthService interface {
		With(ctx context.Context) AuthService

		External(profile goth.User) (*types.User, error)
		FrontendRedirectURL() string

		InternalSignUp(input *types.User, password string) (*types.User, error)
		InternalLogin(email string, password string) (*types.User, error)
		SetPassword(userID uint64, AuthActionPassword string) error
		ChangePassword(userID uint64, oldPassword, AuthActionPassword string) error

		IssueAuthRequestToken(user *types.User) (token string, err error)
		ValidateAuthRequestToken(token string) (user *types.User, err error)
		ValidateEmailConfirmationToken(token string) (user *types.User, err error)
		ExchangePasswordResetToken(token string) (user *types.User, exchangedToken string, err error)
		ValidatePasswordResetToken(token string) (user *types.User, err error)
		SendEmailAddressConfirmationToken(email string) (err error)
		SendPasswordResetToken(email string) (err error)

		CanRegister() error

		LoadRoleMemberships(*types.User) error

		checkPasswordStrength(string) bool
		changePassword(uint64, string) error
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

func Auth(ctx context.Context) AuthService {
	return (&auth{
		eventbus:      eventbus.Service(),
		subscription:  CurrentSubscription,
		settings:      CurrentSettings,
		notifications: DefaultAuthNotification,

		actionlog: DefaultActionlog,

		providerValidator: defaultProviderValidator,

		now: func() *time.Time {
			var now = time.Now()
			return &now
		},
	}).With(ctx)
}

// With returns copy of service with new context
// obsolete approach, will be removed ASAP
func (svc auth) With(ctx context.Context) AuthService {
	db := repository.DB(ctx)
	return &auth{
		db:  db,
		ctx: ctx,

		credentials: repository.Credentials(ctx, db),
		users:       repository.User(ctx, db),
		roles:       repository.Role(ctx, db),

		subscription:      svc.subscription,
		settings:          svc.settings,
		notifications:     svc.notifications,
		eventbus:          svc.eventbus,
		providerValidator: svc.providerValidator,

		actionlog: svc.actionlog,

		now: svc.now,
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
func (svc auth) External(profile goth.User) (u *types.User, err error) {
	var (
		authProvider = &types.AuthProvider{Provider: profile.Provider}

		aam = &authActionProps{
			email:    profile.Email,
			provider: profile.Provider,
			user:     u,
		}
	)

	err = svc.db.Transaction(func() error {
		if !svc.settings.Auth.External.Enabled {
			return AuthErrExternalDisabledByConfig(aam)
		}

		if err = svc.providerValidator(profile.Provider); err != nil {
			return err
		}

		if !reEmail.MatchString(profile.Email) {
			return AuthErrProfileWithoutValidEmail(aam)
		}

		if cc, err := svc.credentials.FindByCredentials(profile.Provider, profile.UserID); err == nil {
			// Credentials found, load user
			for _, c := range cc {
				if !c.Valid() {
					continue
				}

				// Add credentials ID for audit log
				aam.setCredentials(c)

				if u, err = svc.users.FindByID(c.OwnerID); err != nil {
					if repository.ErrUserNotFound.Eq(err) {
						// Orphaned credentials (no owner)
						// try to auto-fix this by removing credentials and recreating user
						if err = svc.credentials.DeleteByID(c.ID); err != nil {
							return err
						} else {
							goto findByEmail
						}
					}
					return err
				}

				// Add user ID for audit log
				aam.setUser(u)
				svc.ctx = internalAuth.SetIdentityToContext(svc.ctx, u)

				if err = svc.eventbus.WaitFor(svc.ctx, event.AuthBeforeLogin(u, authProvider)); err != nil {
					return err
				}

				if u.Valid() {
					// Valid user, Bingo!
					c.LastUsedAt = svc.now()
					if c, err = svc.credentials.Update(c); err != nil {
						return err
					}

					_ = svc.eventbus.WaitFor(svc.ctx, event.AuthAfterLogin(u, authProvider))
					return svc.recordAction(svc.ctx, aam, AuthActionUpdateCredentials, nil)
				} else {
					// Scenario: linked to an invalid user
					if len(cc) > 1 {
						// try with next credentials
						u = nil
						continue
					}

					return AuthErrCredentialsLinkedToInvalidUser(aam)
				}
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
		// Reset audit meta data that might got set during credentials check
		aam.setEmail(profile.Email).
			setCredentials(nil).
			setUser(nil)

		// Find user via his email
		if u, err = svc.users.FindByEmail(profile.Email); repository.ErrUserNotFound.Eq(err) {
			// @todo check if it is ok to auto-create a user here

			// In case we do not have this email, create a new user
			u = &types.User{
				Email:    profile.Email,
				Name:     profile.Name,
				Username: profile.NickName,
			}

			if !handle.IsValid(profile.NickName) {
				u.Handle = profile.NickName
			}

			if err = svc.CanRegister(); err != nil {
				return AuthErrSubscription(aam).Wrap(err)
			}

			if err = svc.eventbus.WaitFor(svc.ctx, event.AuthBeforeSignup(u, authProvider)); err != nil {
				return err
			}

			if u.Handle == "" {
				createHandle(svc.users, u)
			}

			if u, err = svc.users.Create(u); err != nil {
				return err
			}

			aam.setUser(nil)
			svc.ctx = internalAuth.SetIdentityToContext(svc.ctx, u)

			_ = svc.eventbus.WaitFor(svc.ctx, event.AuthAfterSignup(u, authProvider))

			svc.recordAction(svc.ctx, aam, AuthActionExternalSignup, nil)

			// Auto-promote first user
			if err = svc.autoPromote(u); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			// User found
			aam.setUser(u)
			svc.ctx = internalAuth.SetIdentityToContext(svc.ctx, u)

			if err = svc.eventbus.WaitFor(svc.ctx, event.AuthBeforeLogin(u, authProvider)); err != nil {
				return err
			}

			_ = svc.eventbus.WaitFor(svc.ctx, event.AuthAfterLogin(u, authProvider))

			// If user
			if !u.Valid() {
				return AuthErrFailedForDisabledUser(aam).Wrap(err)
			}
		}

		// If we got to this point, assume that user is authenticated
		// but credentials need to be stored
		c := &types.Credentials{
			Kind:        profile.Provider,
			OwnerID:     u.ID,
			Credentials: profile.UserID,
			LastUsedAt:  svc.now(),
		}

		if c, err = svc.credentials.Create(c); err != nil {
			return err
		}

		aam.setCredentials(c)
		svc.recordAction(svc.ctx, aam, AuthActionCreateCredentials, nil)

		// Owner loaded, carry on.
		return nil
	})

	return u, svc.recordAction(svc.ctx, aam, AuthActionAuthenticate, err)
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
func (svc auth) InternalSignUp(input *types.User, password string) (u *types.User, err error) {
	var (
		authProvider = &types.AuthProvider{Provider: credentialsTypePassword}

		aam = &authActionProps{
			email:       input.Email,
			credentials: &types.Credentials{Kind: credentialsTypePassword},
			user:        u,
		}
	)

	err = func() error {
		if !svc.settings.Auth.Internal.Enabled || !svc.settings.Auth.Internal.Signup.Enabled {
			return AuthErrInternalSignupDisabledByConfig(aam)
		}

		if input == nil || !reEmail.MatchString(input.Email) {
			return AuthErrInvalidEmailFormat(aam)
		}

		if len(password) == 0 {
			return AuthErrPasswordNotSecure(aam)
		}

		if !handle.IsValid(input.Handle) {
			return AuthErrInvalidHandle(aam)
		}

		var eUser *types.User
		eUser, err = svc.users.FindByEmail(input.Email)
		if err == nil && eUser.Valid() {
			var cc types.CredentialsSet
			cc, err = svc.credentials.FindByKind(eUser.ID, credentialsTypePassword)
			if err != nil {
				return err
			}

			if c := cc.CompareHashAndPassword(password); c == nil {
				return AuthErrInvalidCredentials(aam)
			} else {
				// Update last-used-by timestamp on matching credentials
				c.LastUsedAt = svc.now()
				aam.setCredentials(c)

				if _, err = svc.credentials.Update(c); err != nil {
					return err
				}
			}

			// We're not actually doing sign-up here - user exists,
			// password is a match, so lets trigger before/after user login events
			if err = svc.eventbus.WaitFor(svc.ctx, event.AuthBeforeLogin(eUser, authProvider)); err != nil {
				return err
			}

			if !eUser.EmailConfirmed {
				err = svc.sendEmailAddressConfirmationToken(eUser)
				if err != nil {
					return err
				}
			}

			_ = svc.eventbus.WaitFor(svc.ctx, event.AuthAfterLogin(eUser, authProvider))
			u = eUser
			return nil

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
		} else if !repository.ErrUserNotFound.Eq(err) {
			return err
		}

		if err = svc.CanRegister(); err != nil {
			return err
		}

		var nUser = &types.User{
			Email:    input.Email,
			Name:     input.Name,
			Username: input.Username,
			Handle:   input.Handle,

			// Do we need confirmed email?
			EmailConfirmed: !svc.settings.Auth.Internal.Signup.EmailConfirmationRequired,
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.AuthBeforeSignup(nUser, authProvider)); err != nil {
			return err
		}

		if input.Handle == "" {
			createHandle(svc.users, input)
		}

		// Whitelisted user data to copy
		u, err = svc.users.Create(nUser)

		if err != nil {
			return err
		}

		aam.setUser(u)
		_ = svc.eventbus.WaitFor(svc.ctx, event.AuthAfterSignup(u, authProvider))

		if err = svc.autoPromote(u); err != nil {
			return err
		}

		if len(password) > 0 {
			err = svc.changePassword(u.ID, password)
			if err != nil {
				return err
			}
		}

		if !u.EmailConfirmed {
			err = svc.sendEmailAddressConfirmationToken(u)
			if err != nil {
				return err
			}

			return svc.recordAction(svc.ctx, aam, AuthActionSendEmailConfirmationToken, nil)
		}

		return nil
	}()

	return u, svc.recordAction(svc.ctx, aam, AuthActionInternalSignup, err)
}

// InternalLogin verifies username/password combination in the internal credentials table
//
// Expects plain text password as an input
func (svc auth) InternalLogin(email string, password string) (u *types.User, err error) {
	var (
		authProvider = &types.AuthProvider{Provider: credentialsTypePassword}

		aam = &authActionProps{
			email:       email,
			credentials: &types.Credentials{Kind: credentialsTypePassword},
			user:        u,
		}
	)

	err = svc.db.Transaction(func() error {
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

		u, err = svc.users.FindByEmail(email)
		if repository.ErrUserNotFound.Eq(err) {
			return AuthErrFailedForUnknownUser()
		}

		if err != nil {
			return err
		}

		// Update audit meta with found user
		svc.ctx = internalAuth.SetIdentityToContext(svc.ctx, u)

		cc, err = svc.credentials.FindByKind(u.ID, credentialsTypePassword)
		if err != nil {
			return err
		}

		if c := cc.CompareHashAndPassword(password); c == nil {
			return AuthErrInvalidCredentials(aam)
		} else {
			// Update last-used-by timestamp on matching credentials
			c.LastUsedAt = svc.now()
			aam.setCredentials(c)

			if _, err = svc.credentials.Update(c); err != nil {
				return err
			}
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.AuthBeforeLogin(u, authProvider)); err != nil {
			return err
		}

		if !u.Valid() {
			return AuthErrFailedForDisabledUser()
		}

		if !u.EmailConfirmed {
			if err = svc.sendEmailAddressConfirmationToken(u); err != nil {
				return err
			}

			return AuthErrFailedUnconfirmedEmail()
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.AuthAfterLogin(u, authProvider))
		return nil
	})

	return u, svc.recordAction(svc.ctx, aam, AuthActionAuthenticate, err)
}

// checkPassword returns true if given (encrypted) password matches any of the credentials
func (svc auth) checkPassword(password string, cc types.CredentialsSet) bool {
	return cc.CompareHashAndPassword(password) != nil
}

// SetPassword sets new password for a user
func (svc auth) SetPassword(userID uint64, password string) (err error) {
	var (
		u *types.User

		aam = &authActionProps{
			user:        u,
			credentials: &types.Credentials{Kind: credentialsTypePassword},
		}
	)

	err = svc.db.Transaction(func() error {
		if !svc.settings.Auth.Internal.Enabled {
			return AuthErrInteralLoginDisabledByConfig(aam)
		}

		if !svc.checkPasswordStrength(password) {
			return AuthErrPasswordNotSecure(aam)
		}

		u, err = svc.users.FindByID(userID)
		if repository.ErrUserNotFound.Eq(err) {
			return AuthErrPasswordChangeFailedForUnknownUser(aam)
		}

		if err != svc.changePassword(userID, password) {
			return err
		}

		return nil
	})

	return svc.recordAction(svc.ctx, aam, AuthActionChangePassword, err)
}

// ChangePassword validates old password and changes it with new
func (svc auth) ChangePassword(userID uint64, oldPassword, AuthActionPassword string) (err error) {
	var (
		u  *types.User
		cc types.CredentialsSet

		aam = &authActionProps{
			user:        u,
			credentials: &types.Credentials{Kind: credentialsTypePassword},
		}
	)

	err = svc.db.Transaction(func() error {
		if !svc.settings.Auth.Internal.Enabled {
			return AuthErrInteralLoginDisabledByConfig(aam)
		}

		if len(oldPassword) == 0 {
			return AuthErrPasswordNotSecure(aam)
		}

		if !svc.checkPasswordStrength(AuthActionPassword) {
			return AuthErrPasswordNotSecure(aam)
		}

		u, err = svc.users.FindByID(userID)
		if repository.ErrUserNotFound.Eq(err) {
			return AuthErrPasswordChangeFailedForUnknownUser(aam)
		}

		cc, err = svc.credentials.FindByKind(userID, credentialsTypePassword)
		if err != nil {
			return err
		}

		if !svc.checkPassword(oldPassword, cc) {
			return AuthErrPasswodResetFailedOldPasswordCheckFailed(aam)
		}

		if err != svc.changePassword(userID, AuthActionPassword) {
			return err
		}

		return nil
	})

	return svc.recordAction(svc.ctx, aam, AuthActionChangePassword, err)
}

func (svc auth) hashPassword(password string) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (svc auth) checkPasswordStrength(password string) bool {
	if len(password) <= 4 {
		return false
	}

	return true
}

// ChangePassword (soft) deletes old password entry and creates a new one
//
// Expects hashed password as an input
func (svc auth) changePassword(userID uint64, password string) (err error) {
	var hash []byte
	if hash, err = svc.hashPassword(password); err != nil {
		return
	}

	if err = svc.credentials.DeleteByKind(userID, credentialsTypePassword); err != nil {
		return
	}

	_, err = svc.credentials.Create(&types.Credentials{
		OwnerID:     userID,
		Kind:        credentialsTypePassword,
		Credentials: string(hash),
	})

	return err
}

// IssueAuthRequestToken returns token that can be used for authentication
func (svc auth) IssueAuthRequestToken(user *types.User) (token string, err error) {
	return svc.createUserToken(user, credentialsTypeAuthToken)
}

// ValidateAuthRequestToken returns user that requested auth token
func (svc auth) ValidateAuthRequestToken(token string) (u *types.User, err error) {
	var (
		aam = &authActionProps{
			credentials: &types.Credentials{Kind: credentialsTypeAuthToken},
		}
	)

	err = svc.db.Transaction(func() error {
		u, err = svc.loadUserFromToken(token, credentialsTypeAuthToken)
		if err != nil && u != nil {
			aam.setUser(u)
			svc.ctx = internalAuth.SetIdentityToContext(svc.ctx, u)
		}
		return err
	})

	return u, svc.recordAction(svc.ctx, aam, AuthActionValidateToken, err)
}

// ValidateEmailConfirmationToken issues a validation token that can be used for
func (svc auth) ValidateEmailConfirmationToken(token string) (user *types.User, err error) {
	return svc.loadFromTokenAndConfirmEmail(token, credentialsTypeEmailAuthToken)
}

// ValidatePasswordResetToken validates password reset token
func (svc auth) ValidatePasswordResetToken(token string) (user *types.User, err error) {
	return svc.loadFromTokenAndConfirmEmail(token, credentialsTypeEmailAuthToken)
}

// loadFromTokenAndConfirmEmail loads token, confirms user's
func (svc auth) loadFromTokenAndConfirmEmail(token, tokenType string) (u *types.User, err error) {
	var (
		aam = &authActionProps{
			user:        u,
			credentials: &types.Credentials{Kind: tokenType},
		}
	)

	err = svc.db.Transaction(func() error {
		if !svc.settings.Auth.Internal.Enabled {
			return AuthErrInternalSignupDisabledByConfig(aam)
		}

		u, err = svc.loadUserFromToken(token, tokenType)
		if err != nil {
			return err
		}

		aam.setUser(u)
		svc.ctx = internalAuth.SetIdentityToContext(svc.ctx, u)

		if u.EmailConfirmed {
			return nil
		}

		u.EmailConfirmed = true
		if u, err = svc.users.Update(u); err != nil {
			return err
		}

		return nil
	})

	return u, svc.recordAction(svc.ctx, aam, AuthActionConfirmEmail, err)
}

// ExchangePasswordResetToken exchanges reset password token for a new one and returns it with user info
func (svc auth) ExchangePasswordResetToken(token string) (u *types.User, t string, err error) {
	var (
		aam = &authActionProps{
			user:        u,
			credentials: &types.Credentials{Kind: credentialsTypeResetPasswordToken},
		}
	)

	err = svc.db.Transaction(func() error {
		if !svc.settings.Auth.Internal.Enabled || !svc.settings.Auth.Internal.PasswordReset.Enabled {
			return AuthErrPasswordResetDisabledByConfig(aam)
		}

		u, err = svc.loadUserFromToken(token, credentialsTypeResetPasswordToken)
		if err != nil {
			return AuthErrInvalidToken(aam).Wrap(err)
		}

		aam.setUser(u)
		svc.ctx = internalAuth.SetIdentityToContext(svc.ctx, u)

		t, err = svc.createUserToken(u, credentialsTypeResetPasswordTokenExchanged)
		if err != nil {
			u = nil
			t = ""
			return AuthErrInvalidToken(aam).Wrap(err)
		}

		return nil
	})

	return u, t, svc.recordAction(svc.ctx, aam, AuthActionExchangePasswordResetToken, err)
}

// SendEmailAddressConfirmationToken sends email with email address confirmation token
func (svc auth) SendEmailAddressConfirmationToken(email string) (err error) {
	var (
		aam = &authActionProps{
			email: email,
		}
	)

	err = svc.db.Transaction(func() error {
		if !svc.settings.Auth.Internal.Enabled || !svc.settings.Auth.Internal.PasswordReset.Enabled {
			return AuthErrPasswordResetDisabledByConfig(aam)
		}

		u, err := svc.users.FindByEmail(email)
		if err != nil {
			return AuthErrInvalidToken(aam)
		}

		return svc.sendEmailAddressConfirmationToken(u)
	})

	return svc.recordAction(svc.ctx, aam, AuthActionSendEmailConfirmationToken, err)
}

func (svc auth) sendEmailAddressConfirmationToken(u *types.User) (err error) {
	var (
		notificationLang = "en"
		token            string

		aam = &authActionProps{
			user:        u,
			credentials: &types.Credentials{Kind: credentialsTypeEmailAuthToken},
		}
	)

	if token, err = svc.createUserToken(u, credentialsTypeEmailAuthToken); err != nil {
		return
	}

	if err = svc.notifications.EmailConfirmation(notificationLang, u.Email, token); err != nil {
		return
	}

	return svc.recordAction(svc.ctx, aam, AuthActionSendEmailConfirmationToken, err)
}

// SendPasswordResetToken sends password reset token to email
func (svc auth) SendPasswordResetToken(email string) (err error) {
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

		if u, err = svc.users.FindByEmail(email); err != nil {
			return err
		}

		svc.ctx = internalAuth.SetIdentityToContext(svc.ctx, u)

		if err = svc.sendPasswordResetToken(u); err != nil {
			return err
		}

		return nil
	}()

	return svc.recordAction(svc.ctx, aam, AuthActionSendPasswordResetToken, err)
}

// CanRegister verifies if user can register
func (svc auth) CanRegister() error {
	if svc.subscription != nil {
		// When we have an active subscription, we need to check
		// if users can register or did this deployment hit
		// it's user-limit
		return svc.subscription.CanRegister(svc.users.Total())
	}

	return nil
}

func (svc auth) sendPasswordResetToken(u *types.User) (err error) {
	var (
		notificationLang = "en"
	)

	token, err := svc.createUserToken(u, credentialsTypeResetPasswordToken)
	if err != nil {
		return err
	}

	return svc.notifications.PasswordReset(notificationLang, u.Email, token)
}

func (svc auth) loadUserFromToken(token, kind string) (u *types.User, err error) {
	var (
		aam = &authActionProps{
			credentials: &types.Credentials{Kind: kind},
		}
	)

	credentialsID, credentials := svc.validateToken(token)
	if credentialsID == 0 {
		return nil, AuthErrInvalidToken(aam)
	}

	c, err := svc.credentials.FindByID(credentialsID)
	if err == repository.ErrCredentialsNotFound {
		return nil, AuthErrInvalidToken(aam)
	}

	aam.setCredentials(c)

	if err != nil {
		return
	}

	if err = svc.credentials.DeleteByID(c.ID); err != nil {
		return
	}

	if !c.Valid() || c.Credentials != credentials {
		return nil, AuthErrInvalidToken(aam)
	}

	u, err = svc.users.FindByID(c.OwnerID)
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
func (svc auth) createUserToken(u *types.User, kind string) (token string, err error) {
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
			expiresAt = svc.now().Add(time.Second * 15)
		default:
			// 1h expiration for all tokens send via email
			expiresAt = svc.now().Add(time.Minute * 60)
		}

		c, err := svc.credentials.Create(&types.Credentials{
			OwnerID:     u.ID,
			Kind:        kind,
			Credentials: string(rand.Bytes(credentialsTokenLength)),
			ExpiresAt:   &expiresAt,
		})

		if err != nil {
			return err
		}

		token = fmt.Sprintf("%s%d", c.Credentials, c.ID)
		return nil
	}()

	return token, svc.recordAction(svc.ctx, aam, AuthActionIssueToken, err)
}

// Automatically promotes user to administrator if it is the first user in the database
func (svc auth) autoPromote(u *types.User) (err error) {
	if svc.users.Total() > 1 || u.ID == 0 {
		return nil
	}

	if svc.roles == nil {
		// no role repository; auto-promotion disabled
		return nil
	}

	var (
		roleID = permissions.AdminsRoleID
		aam    = &authActionProps{user: u, role: &types.Role{ID: roleID}}
	)

	err = svc.roles.MemberAddByID(roleID, u.ID)
	return svc.recordAction(svc.ctx, aam, AuthActionAutoPromote, err)
}

// LoadRoleMemberships loads membership info
func (svc auth) LoadRoleMemberships(u *types.User) error {
	rr, _, err := svc.roles.Find(types.RoleFilter{MemberID: u.ID})
	if err != nil {
		return err
	}

	u.SetRoles(rr.IDs())
	return nil
}
