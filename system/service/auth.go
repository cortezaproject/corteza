package service

import (
	"context"
	"regexp"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/markbates/goth"
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

		// how fresh can a password be before we consider it
		// reused?
		// @todo make this configurable someday
		PasswordReuseTimeWindow time.Duration
	}

	authAccessController interface {
		CanImpersonateUser(context.Context, *types.User) bool
		CanUpdateUser(context.Context, *types.User) bool
	}
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
func (svc *auth) External(ctx context.Context, profile types.ExternalAuthUser) (u *types.User, err error) {
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
			cc types.CredentialSet
			f  = types.CredentialFilter{Kind: profile.Provider, Credentials: profile.UserID}
		)

		if cc, _, err = store.SearchCredentials(ctx, svc.store, f); err == nil {
			// Credential found, load user
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
						if err = store.DeleteCredentialByID(ctx, svc.store, c.ID); err != nil {
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
		c := &types.Credential{
			ID:          nextID(),
			CreatedAt:   *now(),
			Kind:        profile.Provider,
			OwnerID:     u.ID,
			Credentials: profile.UserID,
			LastUsedAt:  now(),
		}

		if err = store.CreateCredential(ctx, svc.store, c); err != nil {
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
func (svc *auth) InternalSignUp(ctx context.Context, input *types.User, password string) (u *types.User, err error) {
	var (
		authProvider = &types.AuthProvider{Provider: credentialsTypePassword}

		aam = &authActionProps{
			email:       input.Email,
			credentials: &types.Credential{Kind: credentialsTypePassword},
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
			// making sure password is not empty
			// proper strength check is done a bit lower, after user existance check
			return AuthErrPasswordNotSecure(aam)
		}

		var eUser *types.User
		eUser, err = store.LookupUserByEmail(ctx, svc.store, input.Email)

		if err == nil && eUser != nil {
			var (
				c  *types.Credential
				cc types.CredentialSet
				f  = types.CredentialFilter{OwnerID: eUser.ID, Kind: credentialsTypePassword}
			)
			if cc, _, err = store.SearchCredentials(ctx, svc.store, f); err != nil {
				return err
			}

			// does password match any of the valid credentials?
			c = findValidPassword(cc, password)
			if c == nil {
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

		// The check must be after the login fallback so that we still allow
		// logins with old passwords in case the policy has changed since then.
		if !svc.CheckPasswordStrength(password) {
			return AuthErrPasswordNotSecure()
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
func (svc *auth) InternalLogin(ctx context.Context, email string, password string) (u *types.User, err error) {
	var (
		authProvider = &types.AuthProvider{Provider: credentialsTypePassword}

		aam = &authActionProps{
			email:       email,
			credentials: &types.Credential{Kind: credentialsTypePassword},
			user:        u,
		}

		c *types.Credential
	)

	err = func() error {
		if !svc.settings.Auth.Internal.Enabled {
			return AuthErrInternalLoginDisabledByConfig()
		}

		if !reEmail.MatchString(email) {
			return AuthErrInvalidEmailFormat()
		}

		if len(password) == 0 {
			// making sure password is not empty
			// we're not checking for strength here, users might
			// use weak passwords from before the policy was introduced
			return AuthErrInvalidCredentials()
		}

		var (
			cc types.CredentialSet
		)

		u, err = store.LookupUserByEmail(ctx, svc.store, email)
		if errors.IsNotFound(err) {
			return AuthErrInvalidCredentials(aam)
		} else if err != nil {
			return err
		}

		// Update audit meta with found user
		ctx = internalAuth.SetIdentityToContext(ctx, u)
		cc, _, err = store.SearchCredentials(ctx, svc.store, types.CredentialFilter{OwnerID: u.ID, Kind: credentialsTypePassword})
		if err != nil {
			return err
		}

		// find 1st valid credentials that match the hashed password
		c = findValidPassword(cc, password)
		if c == nil {
			return AuthErrInvalidCredentials(aam)
		}

		aam.setCredentials(c)
		ctx = internalAuth.SetIdentityToContext(ctx, u)

		return svc.procLogin(ctx, svc.store, u, c, authProvider)
	}()

	return u, svc.recordAction(ctx, aam, AuthActionAuthenticate, err)
}

// Impersonate verifies if user can impersonate another user and returns that user
//
// For now, it's the caller's responsibility to generate the auth token
func (svc *auth) Impersonate(ctx context.Context, userID uint64) (u *types.User, err error) {
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

// procLogin fn performs standard validation, credentials-update tasks and triggers events
func (svc *auth) procLogin(ctx context.Context, s store.Storer, u *types.User, c *types.Credential, p *types.AuthProvider) (err error) {
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
		if err = store.UpdateCredential(ctx, s, c); err != nil {
			return err
		}
	}

	_ = svc.eventbus.WaitFor(ctx, event.AuthAfterLogin(u, p))
	return nil
}

// Automatically promotes user to super-administrator if it is the first non-system user in the database
func (svc *auth) autoPromote(ctx context.Context, u *types.User) (err error) {
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

// LoadRoleMemberships loads membership info
//
// @todo move this to role service
func (svc *auth) LoadRoleMemberships(ctx context.Context, u *types.User) error {
	rr, _, err := store.SearchRoles(ctx, svc.store, types.RoleFilter{MemberID: u.ID})
	if err != nil {
		return err
	}

	u.SetRoles(rr.IDs()...)
	return nil
}

func (svc *auth) GetProviders() types.ExternalAuthProviderSet {
	return CurrentSettings.Auth.External.Providers
}

func (svc *auth) checkLimits(ctx context.Context) error {
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
