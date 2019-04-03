package service

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/markbates/goth"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/types"
)

type (
	auth struct {
		db  db
		ctx context.Context

		credentials repository.CredentialsRepository
		users       repository.UserRepository

		providerValidator func(string) error
		now               func() *time.Time
	}

	AuthService interface {
		With(ctx context.Context) AuthService

		External(profile goth.User) (*types.User, error)

		CheckPassword(email string, password []byte) (*types.User, error)
		ChangePassword(user *types.User, password []byte) error
		CheckCredentials(credentialsID uint64, secret string) (*types.User, error)
		RevokeCredentialsByID(user *types.User, credentialsID uint64) error
	}
)

const (
	CredentialsTypePassword = "password"
)

var (
	reEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func defaultProviderValidator(provider string) error {
	_, err := goth.GetProvider(provider)
	return err
}

func Auth(ctx context.Context) AuthService {
	return (&auth{}).With(ctx)
}

func (svc *auth) With(ctx context.Context) AuthService {
	db := repository.DB(ctx)
	return &auth{
		db:  db,
		ctx: ctx,

		credentials: repository.Credentials(ctx, db),
		users:       repository.User(ctx, db),

		providerValidator: defaultProviderValidator,
		now: func() *time.Time {
			var now = time.Now()
			return &now
		},
	}
}

// Social user verifies existance by using email value from social profile and creates user if needed
//
// It does not update user's info
func (svc *auth) External(profile goth.User) (u *types.User, err error) {
	if err = svc.providerValidator(profile.Provider); err != nil {
		return nil, err
	}

	if profile.Email == "" {
		return nil, errors.New("Can not use profile data without an email")
	}

	return u, svc.db.Transaction(func() error {
		var c *types.Credentials
		if cc, err := svc.credentials.FindByCredentials(profile.Provider, profile.UserID); err == nil {
			// Credentials found, load user
			for _, c := range cc {
				if !c.Valid() {
					continue
				}

				if u, err = svc.users.FindByID(c.OwnerID); err != nil {
					if err == repository.ErrUserNotFound {
						// Orphaned credentials (no owner)
						// try to auto-fix this by removing credentials and recreating user
						if err := svc.credentials.DeleteByID(c.ID); err != nil {
							return errors.Wrap(err, "could not cleanup orphaned credentials")
						} else {
							goto findByEmail
						}
					}
					return nil
				} else if u.Valid() && u.Email != profile.Email {
					return errors.Errorf(
						"Refusing to authenticate with non matching emails (profile: %v, db: %v) on credentials (ID: %v)",
						profile.Email,
						u.Email,
						c.ID)
				} else if u.Valid() {
					// Valid user, matching emails. Bingo!
					c.LastUsedAt = svc.now()
					if c, err = svc.credentials.Update(c); err != nil {
						return err
					}

					log.Printf(
						"updating credential entry (%v, %v) for exisintg user (%v, %v)",
						c.ID,
						profile.Provider,
						u.ID,
						u.Email,
					)

					return nil
				} else {
					// Scenario: linked to an invalid user
					u = nil
					continue
				}
			}

			// If we could not find anything useful,
			// we can search for user via email
		} else {
			// A serious error occured, bail out...
			return err
		}

	findByEmail:
		// Find user via his email
		if u, err = svc.users.FindByEmail(profile.Email); err == repository.ErrUserNotFound {
			// @todo check if it is ok to auto-create a user here

			// In case we do not have this email, create a new user
			u = &types.User{
				Email:    profile.Email,
				Name:     profile.Name,
				Username: profile.NickName,
				Handle:   profile.NickName,
			}

			if u, err = svc.users.Create(u); err != nil {
				return errors.Wrap(err, "could not create user after successful external authentication")
			}

			log.Printf("created new user after successful social auth (%v, %v)", u.ID, u.Email)
		} else if err != nil {
			return err
		} else if !u.Valid() {
			return errors.Errorf(
				"can not use invalid user (user ID: %v)",
				u.ID,
			)
		} else {
			log.Printf(
				"autheticated user (%v, %v) via %s, existing user",
				u.ID,
				u.Email,
				profile.Provider,
			)
		}

		c = &types.Credentials{
			Kind:        profile.Provider,
			OwnerID:     u.ID,
			Credentials: profile.UserID,
			LastUsedAt:  svc.now(),
		}

		if !profile.ExpiresAt.IsZero() {
			// Copy expiration date when provided
			c.ExpiresAt = &profile.ExpiresAt
		}

		if c, err = svc.credentials.Create(c); err != nil {
			return err
		}

		log.Printf(
			"creating new credential entry (%v, %v) for exisintg user (%v, %v)",
			c.ID,
			profile.Provider,
			u.ID,
			u.Email,
		)

		// Owner loaded, carry on.
		return nil
	})
}

// CheckPassword verifies username/password combination
//
// Expects plain text password as an input
func (svc *auth) CheckPassword(email string, password []byte) (u *types.User, err error) {
	if err = svc.validateCredentials(email, password); err != nil {
		return
	}

	return u, svc.db.Transaction(func() error {
		var (
			cc types.CredentialsSet
		)

		u, err = svc.users.FindByEmail(email)
		if err != repository.ErrUserNotFound {
			return errors.New("invalid username/password combination")
		}

		if err != nil {
			return errors.Wrap(err, "could not check password")
		}

		cc, err := svc.credentials.FindByKind(u.ID, CredentialsTypePassword)
		if err != nil {
			return errors.Wrap(err, "could not find credentials")
		}

		return svc.checkPassword(password, cc)
	})
}

// validateCredentials does basic format & length check
func (svc auth) validateCredentials(email string, password []byte) error {
	if !reEmail.MatchString(email) {
		return errors.New("invalid email format")
	}

	if len(password) == 0 {
		return errors.New("empty password")
	}

	return nil
}

func (svc auth) checkPassword(password []byte, cc types.CredentialsSet) (err error) {
	// We need only valid credentials (skip deleted, expired)
	cc, _ = cc.Filter(func(c *types.Credentials) (b bool, e error) {
		return c.Valid(), nil
	})

	for _, c := range cc {
		if len(c.Credentials) == 0 {
			continue
		}

		err = bcrypt.CompareHashAndPassword([]byte(c.Credentials), password)
		if err == bcrypt.ErrMismatchedHashAndPassword {
			// Mismatch, continue with the checking
			continue
		} else if err != nil {
			// Some other error
			return errors.Wrap(err, "could not compare passwords")
		} else {
			// Password matched one of credentials
			return nil
		}
	}

	return errors.New("invalid username/password combination")
}

// ChangePassword (soft) deletes old password entry and creates a new one
//
// Expects plain text password as an input
func (svc *auth) ChangePassword(user *types.User, password []byte) (err error) {
	var hash []byte

	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "could not hash password")
	}

	return svc.db.Transaction(func() error {
		if err = svc.credentials.DeleteByKind(user.ID, CredentialsTypePassword); err != nil {
			return errors.Wrap(err, "could not remove old passswords")
		}

		_, err = svc.credentials.Create(&types.Credentials{
			OwnerID:     user.ID,
			Kind:        CredentialsTypePassword,
			Credentials: string(hash),
		})

		return errors.Wrap(err, "could not create new password")
	})
}

// CheckCredentials searches for credentials/secret combination and returns loaded user if successful
func (svc *auth) CheckCredentials(credentialsID uint64, secret string) (*types.User, error) {
	panic("svc.auth.CheckCredentials, not implemented")
}

// RevokeCredentialsByID (soft) deletes credentials by id
func (svc *auth) RevokeCredentialsByID(user *types.User, credentialsID uint64) error {
	panic("svc.auth.RevokeCredentialsByID, not implemented")
}

var _ AuthService = &auth{}
