package service

// part of auth service
// collection of functions that handle password checking, setting, resetting, changing
//
// general credential handling functions should still be part of auth.go

import (
	"context"
	"fmt"
	internalAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/dgryski/dgoogauth"
	"golang.org/x/crypto/bcrypt"
	rand2 "math/rand"
	"regexp"
	"sort"
	"strconv"
	"time"
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

	tokenReqMaxCount  = 5
	tokenReqMaxWindow = time.Minute * 15

	passwordMinLength = 8
	passwordMaxLength = 256
)

var (
	oneTokenPerUser = map[string]bool{
		credentialsTypeResetPasswordToken: true,
	}
)

// ValidateEmailConfirmationToken issues a validation token that can be used for
func (svc *auth) ValidateEmailConfirmationToken(ctx context.Context, token string) (user *types.User, err error) {
	return svc.loadFromTokenAndConfirmEmail(ctx, token, credentialsTypeEmailAuthToken)
}

// loadFromTokenAndConfirmEmail loads token, confirms user's
func (svc *auth) loadFromTokenAndConfirmEmail(ctx context.Context, token, tokenType string) (u *types.User, err error) {
	var (
		aam = &authActionProps{
			user:        u,
			credentials: &types.Credential{Kind: tokenType},
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

func (svc *auth) SendEmailAddressConfirmationToken(ctx context.Context, u *types.User) (err error) {
	var (
		token string

		aam = &authActionProps{
			user:        u,
			credentials: &types.Credential{Kind: credentialsTypeEmailAuthToken},
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

// Loads user from token and removes that token right after
func (svc *auth) loadUserFromToken(ctx context.Context, token, kind string) (u *types.User, _ error) {
	var (
		aam = &authActionProps{
			credentials: &types.Credential{Kind: kind},
		}
	)

	return u, store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		credentialsID, credentials := validateToken(token)
		if credentialsID == 0 {
			return AuthErrInvalidToken(aam)
		}

		c, err := store.LookupCredentialByID(ctx, s, credentialsID)
		if errors.IsNotFound(err) {
			return AuthErrInvalidToken(aam)
		}

		aam.setCredentials(c)

		if err != nil {
			return
		}

		if err = store.DeleteCredentialByID(ctx, s, c.ID); err != nil {
			return
		}

		if !c.Valid() || c.Credentials != credentials {
			return AuthErrInvalidToken(aam)
		}

		u, err = store.LookupUserByID(ctx, s, c.OwnerID)
		if err != nil {
			return err
		}

		aam.setUser(u)

		// context will be updated with new identity
		// in the caller fn

		if !u.Valid() {
			return AuthErrInvalidCredentials(aam)
		}

		return nil
	})
}

// Generates & stores user token
// it returns combined value of token + token ID to help with the lookups
func (svc *auth) createUserToken(ctx context.Context, u *types.User, kind string) (token string, err error) {
	var (
		expiresAt time.Time
		aam       = &authActionProps{
			user:        u,
			credentials: &types.Credential{Kind: kind},
		}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if u == nil || u.ID == 0 {
			return AuthErrGeneric()
		}

		// Rate limit requests
		cc, _, err := store.SearchCredentials(ctx, s, types.CredentialFilter{
			OwnerID: u.ID,
			Kind:    kind,

			// we want to count deleted tokens as well
			Deleted: filter.StateInclusive,
		})

		if err != nil {
			return err
		}

		// gt/eq since this current request is not yet stored
		if err = svc.checkTokenRate(cc, tokenReqMaxWindow, tokenReqMaxCount); err != nil {
			return
		}

		// removes expired and soft-deleted tokens
		// and enforces one-token-per-user rule
		if err = svc.cleanupCredentials(ctx, s, cc); err != nil {
			return
		}

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

		c := &types.Credential{
			ID:          nextID(),
			CreatedAt:   *now(),
			OwnerID:     u.ID,
			Kind:        kind,
			Credentials: token,
			ExpiresAt:   &expiresAt,
		}

		err = store.CreateCredential(ctx, s, c)

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
	})

	return token, svc.recordAction(ctx, aam, AuthActionIssueToken, err)
}

// checks existing tokens and ensure that the creation rate is within limits
func (svc *auth) checkTokenRate(cc types.CredentialSet, window time.Duration, max int) error {
	if len(cc) == 0 || window == 0 || max == 0 {
		return nil
	}

	var (
		cutoff = now().Add(window * -1)
		count  = 0
	)

	for _, c := range cc {
		if c.CreatedAt.Before(cutoff) {
			// skip tokens created before cutoff
			continue
		}

		count++

		if count > max {
			break
		}
	}

	if count > max {
		return AuthErrRateLimitExceeded()
	}

	return nil
}

func (svc *auth) cleanupCredentials(ctx context.Context, s store.Credentials, cc types.CredentialSet) (err error) {
	var (
		update types.CredentialSet
		remove types.CredentialSet
	)

	for _, c := range cc {
		switch {
		case oneTokenPerUser[c.Kind]:
			// if token type is shortlisted in one-token-per-user
			// mark all existing tokens as deleted if to
			//
			// only want to mark them as deleted ad
			c.DeletedAt = now()
			update = append(update, c)

		case false, // just a placeholder
			(c.DeletedAt != nil && c.DeletedAt.Add(tokenReqMaxWindow).Before(*now())),
			(c.ExpiresAt != nil && c.ExpiresAt.Before(*now())):
			// schedule all soft-deleted and expired token
			// for removal
			remove = append(remove, c)
		}
	}

	if err = store.UpdateCredential(ctx, s, update...); err != nil {
		return
	}

	if err = store.DeleteCredential(ctx, s, remove...); err != nil {
		return
	}

	return
}

// SendPasswordResetToken sends password reset token to email
func (svc *auth) SendPasswordResetToken(ctx context.Context, email string) (err error) {
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

func (svc *auth) sendPasswordResetToken(ctx context.Context, u *types.User) (err error) {
	token, err := svc.createUserToken(ctx, u, credentialsTypeResetPasswordToken)
	if err != nil {
		return err
	}

	return svc.notifications.PasswordReset(ctx, u.Email, token)
}

// GeneratePasswordCreateToken generates password create token
func (svc *auth) GeneratePasswordCreateToken(ctx context.Context, email string) (url string, err error) {
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

func (svc *auth) sendPasswordCreateToken(ctx context.Context, u *types.User) (url string, err error) {
	token, err := svc.createUserToken(ctx, u, credentialsTypeCreatePasswordToken)
	if err != nil {
		return
	}

	return svc.notifications.PasswordCreate(token)
}

// ValidatePasswordResetToken validates password reset token
func (svc *auth) ValidatePasswordResetToken(ctx context.Context, token string) (user *types.User, err error) {
	return svc.loadFromTokenAndConfirmEmail(ctx, token, credentialsTypeResetPasswordToken)
}

// ValidatePasswordCreateToken validates password create token
func (svc *auth) ValidatePasswordCreateToken(ctx context.Context, token string) (user *types.User, err error) {
	return svc.loadFromTokenAndConfirmEmail(ctx, token, credentialsTypeCreatePasswordToken)
}

// ExchangePasswordResetToken exchanges reset password token for a new one and returns it with user info
func (svc *auth) ExchangePasswordResetToken(ctx context.Context, token string) (u *types.User, t string, err error) {
	var (
		aam = &authActionProps{
			user:        u,
			credentials: &types.Credential{Kind: credentialsTypeResetPasswordToken},
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

// ChangePassword validates old password and changes it with new
func (svc *auth) ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) (err error) {
	var (
		u  *types.User
		cc types.CredentialSet

		aam = &authActionProps{
			user:        u,
			credentials: &types.Credential{Kind: credentialsTypePassword},
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

		cc, _, err = store.SearchCredentials(ctx, svc.store, types.CredentialFilter{
			Kind:    credentialsTypePassword,
			OwnerID: userID,
			Deleted: filter.StateInclusive})

		if err != nil {
			return err
		}

		if c := findValidPassword(cc, oldPassword); c == nil {
			return AuthErrPasswordResetFailedOldPasswordCheckFailed(aam)
		}

		if isPasswordReused(cc, newPassword, svc.opt.PasswordReuseTimeWindow) {
			return AuthErrPasswordSetFailedReusedPasswordCheckFailed(aam)
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

func hashPassword(password string) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (svc *auth) CheckPasswordStrength(password string) bool {
	return checkPasswordStrength(password, svc.settings.Auth.Internal.PasswordConstraints)
}

// SetPasswordCredentials (soft) deletes old password entry and creates a new entry with new password on every change
//
// # This way we can implement more strict password-change policies in the future
//
// This method is used by auth and user procedures to unify password hashing and updating
// credentials
func (svc *auth) SetPasswordCredentials(ctx context.Context, userID uint64, password string) (err error) {
	if err = svc.removePasswordCredentials(ctx, userID); err != nil {
		return
	}

	return SetPasswordCredentials(ctx, svc.store, userID, password)
}

// SetPasswordCredentials creates a new password entry
func SetPasswordCredentials(ctx context.Context, s store.Storer, userID uint64, password string) (err error) {
	var (
		hash []byte
	)

	if hash, err = hashPassword(password); err != nil {
		return
	}

	// Add new credentials with new password
	c := &types.Credential{
		ID:          nextID(),
		CreatedAt:   *now(),
		OwnerID:     userID,
		Kind:        credentialsTypePassword,
		Credentials: string(hash),
	}

	return store.CreateCredential(ctx, s, c)
}

//CheckPassword verifies if password matches any of the valid credentials
//func (svc *auth) CheckPassword(cc types.CredentialSet, password string) bool {
//	return findMatchingCredentials(cc, password, true) != nil
//}

// SetPassword sets new password for a user
//
// # This function also records an action
//
// this method is used in 2 scenarios:
//
//	SELF:
//	  user forgot the password and needs to reset it
//	  there should be protocols prior to this point that
//	  authenticate and validate users
//
//	USER MANAGEMENT:
//	  administrator is resetting password for another user
func (svc *auth) SetPassword(ctx context.Context, userID uint64, password string) (err error) {
	var (
		u  *types.User
		cc types.CredentialSet

		aam = &authActionProps{
			user:        u,
			credentials: &types.Credential{Kind: credentialsTypePassword},
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

		cc, _, err = store.SearchCredentials(ctx, svc.store, types.CredentialFilter{
			Kind:    credentialsTypePassword,
			OwnerID: userID,
			Deleted: filter.StateInclusive})

		if err != nil {
			return err
		}

		if isPasswordReused(cc, password, svc.opt.PasswordReuseTimeWindow) {
			return AuthErrPasswordSetFailedReusedPasswordCheckFailed(aam)
		}

		if err != svc.SetPasswordCredentials(ctx, userID, password) {
			return err
		}

		return nil
	}()

	return svc.recordAction(ctx, aam, AuthActionChangePassword, err)
}

// PasswordSet checks and returns true if user's password is set
// False is also returned in case user does not exist.
func (svc *auth) PasswordSet(ctx context.Context, email string) (is bool) {
	//svc.settings.Auth.External.Enabled
	u, err := store.LookupUserByEmail(ctx, svc.store, email)
	if err != nil {
		return
	}

	cc, _, err := store.SearchCredentials(ctx, svc.store, types.CredentialFilter{
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

// RemovePasswordCredentials (soft) deletes old password entry
func (svc *auth) RemovePasswordCredentials(ctx context.Context, userID uint64) (err error) {
	// Do a partial update and soft-delete all
	return svc.removePasswordCredentials(ctx, userID)
}

// RemovePasswordCredentials (soft) deletes old password entry
func (svc *auth) removePasswordCredentials(ctx context.Context, userID uint64) (err error) {
	var (
		cc types.CredentialSet
		f  = types.CredentialFilter{Kind: credentialsTypePassword, OwnerID: userID}
	)

	if cc, _, err = store.SearchCredentials(ctx, svc.store, f); err != nil {
		return nil
	}

	// Mark all credentials as deleted
	_ = cc.Walk(func(c *types.Credential) error {
		c.DeletedAt = now()
		return nil
	})

	// Do a partial update and soft-delete all
	return store.UpdateCredential(ctx, svc.store, cc...)
}

// RemoveAccessTokens removes all user's access tokens when suspended,
// deleted or security context changes
func (svc *auth) RemoveAccessTokens(ctx context.Context, user *types.User) error {
	return svc.recordAction(
		ctx,
		&authActionProps{user: user},
		AuthActionAccessTokensRemoved,
		store.DeleteAuthOA2TokenByUserID(ctx, svc.store, user.ID),
	)
}

// ValidateTOTP checks given code with the current secret
// Fn fails if no secret is set
func (svc *auth) ValidateTOTP(ctx context.Context, code string) (err error) {
	var (
		c    *types.Credential
		u    *types.User
		kind = credentialsTypeMfaTotpSecret
		aam  = &authActionProps{credentials: &types.Credential{Kind: kind}}
		i    = internalAuth.GetIdentityFromContext(ctx)
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
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
			return store.UpdateCredential(ctx, s, c)
		}
	})

	return svc.recordAction(ctx, aam, AuthActionTotpValidate, err)
}

// ConfigureTOTP stores totp secret in user's credentials
//
// It returns the user with security policy changes
func (svc *auth) ConfigureTOTP(ctx context.Context, secret string, code string) (u *types.User, err error) {
	var (
		kind = credentialsTypeMfaTotpSecret
		aam  = &authActionProps{credentials: &types.Credential{Kind: kind}}
		i    = internalAuth.GetIdentityFromContext(ctx)
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
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

		cred := &types.Credential{
			ID:          nextID(),
			CreatedAt:   *now(),
			OwnerID:     u.ID,
			Kind:        kind,
			Credentials: secret,
		}

		if err = store.CreateCredential(ctx, s, cred); err != nil {
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
func (svc *auth) RemoveTOTP(ctx context.Context, userID uint64, code string) (u *types.User, err error) {
	var (
		c    *types.Credential
		kind = credentialsTypeMfaTotpSecret
		aam  = &authActionProps{credentials: &types.Credential{Kind: kind}}
		i    = internalAuth.GetIdentityFromContext(ctx)
		self = i != nil && i.Identity() == userID
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
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
func (svc *auth) getTOTPSecret(ctx context.Context, s store.Credentials, userID uint64) (*types.Credential, error) {
	cc, _, err := store.SearchCredentials(ctx, s, types.CredentialFilter{
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
	cc, _, err := store.SearchCredentials(ctx, s, types.CredentialFilter{
		OwnerID: userID,
		Kind:    credentialsTypeMfaTotpSecret,
		Deleted: filter.StateExcluded,
	})

	if err != nil {
		return err
	}

	return cc.Walk(func(c *types.Credential) error {
		c.DeletedAt = now()
		return store.UpdateCredential(ctx, s, c)
	})
}

func (svc *auth) SendEmailOTP(ctx context.Context) (err error) {
	var (
		otp  string
		u    *types.User
		kind = credentialsTypeMFAEmailOTP
		aam  = &authActionProps{credentials: &types.Credential{Kind: kind}}
		i    = internalAuth.GetIdentityFromContext(ctx)
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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

func (svc *auth) ConfigureEmailOTP(ctx context.Context, userID uint64, enable bool) (u *types.User, err error) {
	var (
		kind = credentialsTypeMFAEmailOTP
		aam  = &authActionProps{credentials: &types.Credential{Kind: kind}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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
func (svc *auth) ValidateEmailOTP(ctx context.Context, code string) (err error) {
	var (
		cc   types.CredentialSet
		u    *types.User
		kind = credentialsTypeMFAEmailOTP
		aam  = &authActionProps{credentials: &types.Credential{Kind: kind}}
		i    = internalAuth.GetIdentityFromContext(ctx)
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
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

		cc, _, err = store.SearchCredentials(ctx, s, types.CredentialFilter{
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
			return store.DeleteCredential(ctx, s, c)
		}

		return AuthErrInvalidEmailOTP()
	})

	return svc.recordAction(ctx, aam, AuthActionEmailOtpVerify, err)
}

func validateToken(token string) (ID uint64, credentials string) {
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

// returns matching credential if (hashed version of a) password is found in the
// list of (valid) credentials
//
// should be used as a parameter for credentialsFilter fn
func findValidPassword(cc []*types.Credential, password string) (c *types.Credential) {
	cList := credentialsFilter(cc, 1, skipInvalid, compareHashedCredentials(password))
	if len(cList) > 0 {
		c = cList[0]
	}
	return
}

// returns true if (hashed version of a) password is found in the
// list of given credentials
//
// should be used as a parameter for credentialsFilter fn
func isPasswordReused(cc []*types.Credential, password string, reuseWindow time.Duration) bool {
	return len(credentialsFilter(cc, -1, compareHashedCredentials(password), skipNewerCredentials(reuseWindow))) > 0
}

// skips all invalid credentials
//
// should be used as a parameter for credentialsFilter fn
func skipInvalid(c *types.Credential) bool {
	return c.Valid()
}

func compareHashedCredentials(password string) func(c *types.Credential) bool {
	var (
		p = []byte(password)
	)

	return func(c *types.Credential) bool {
		return bcrypt.CompareHashAndPassword([]byte(c.Credentials), p) == nil
	}
}

func skipNewerCredentials(cutoff time.Duration) func(c *types.Credential) bool {
	var (
		t = now().Add(cutoff * -1)
	)

	return func(c *types.Credential) bool {
		return c.CreatedAt.Before(t)
	}
}

// CompareHashAndPassword returns first valid credentials with matching hash
func credentialsFilter(cc []*types.Credential, limit int, mm ...func(*types.Credential) bool) (out []*types.Credential) {
	// sort credentials by ID (and effectively from newest to oldest)
	sort.Slice(cc, func(i, j int) bool {
		return cc[i].ID > cc[j].ID
	})

	for _, c := range cc {
		if len(c.Credentials) == 0 {
			continue
		}

		next := false
		for _, m := range mm {
			if !m(c) {
				next = true
				break
			}
		}

		if next {
			continue
		}

		out = append(out, c)
		if limit == len(out) {
			break
		}
	}

	return
}

func checkPasswordStrength(password string, pc types.PasswordConstraints) bool {
	var (
		length = len(password)
		re     *regexp.Regexp
		mt     [][]int
	)

	// Always check system constraints
	if length < passwordMinLength || length > passwordMaxLength {
		return false
	}

	// Ignore defined password constraints
	if !pc.PasswordSecurity {
		return true
	}

	// Check the password length
	if length < int(pc.MinLength) {
		return false
	}

	// Check special constraints
	// - numeric characters
	if count := int(pc.MinNumCount); count > 0 {
		re = regexp.MustCompile("[0-9]")
		mt = re.FindAllStringIndex(password, -1)
		if len(mt) < count {
			return false
		}
	}

	// Check for lowercase characters
	if count := int(pc.MinLowerCase); count > 0 {
		re = regexp.MustCompile("[a-z]")
		mt = re.FindAllStringIndex(password, -1)
		if len(mt) < count {
			return false
		}
	}

	// Check for upper-case characters
	if count := int(pc.MinUpperCase); count > 0 {
		re = regexp.MustCompile("[A-Z]")
		mt = re.FindAllStringIndex(password, -1)
		if len(mt) < count {
			return false
		}
	}

	// - special characters
	if count := int(pc.MinSpecialCount); count > 0 {
		re = regexp.MustCompile("[^0-9a-zA-Z]")
		mt = re.FindAllStringIndex(password, -1)
		if len(mt) < count {
			return false
		}
	}

	return true
}
