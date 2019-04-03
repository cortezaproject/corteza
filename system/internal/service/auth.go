package service

import (
	"context"
	"log"
	"time"

	"github.com/markbates/goth"
	"github.com/pkg/errors"

	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/types"
)

type (
	auth struct {
		db  db
		ctx context.Context

		credentials repository.CredentialsRepository
		users       repository.UserRepository
	}

	AuthService interface {
		With(ctx context.Context) AuthService

		External(profile goth.User) (*types.User, error)

		CheckPassword(username, password string) (*types.User, error)
		ChangePassword(user *types.User, password string) error
		CheckCredentials(credentialsID uint64, secret string) (*types.User, error)
		RevokeCredentialsByID(user *types.User, credentialsID uint64) error
	}
)

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
	}
}

// Social user verifies existance by using email value from social profile and creates user if needed
//
// It does not update user's info
func (svc *auth) External(profile goth.User) (u *types.User, err error) {
	var lastUsedAt = time.Now()

	if _, err := goth.GetProvider(profile.Provider); err != nil {
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
				return errors.Wrap(err, "could not create user after successful social auth")
			}

			log.Printf("Created new user after successful social auth (%v, %v)", u.ID, u.Email)

			// Owner created
			return nil
		} else if err != nil {
			return err
		} else if !u.Valid() {
			return errors.Errorf(
				"Social login to an invalid/suspended user (user ID: %v)",
				u.ID,
			)
		} else {
			log.Printf(
				"Autheticated user (%v, %v) via %s, existing user",
				u.ID,
				u.Email,
				profile.Provider,
			)
		}

		if c == nil {
			c = &types.Credentials{
				Kind:        profile.Provider,
				OwnerID:     u.ID,
				Credentials: profile.UserID,
				LastUsedAt:  &lastUsedAt,
			}

			if !profile.ExpiresAt.IsZero() {
				// Copy expiration date when provided
				c.ExpiresAt = &profile.ExpiresAt
			}

			if c, err = svc.credentials.Create(c); err != nil {
				return err
			}

			log.Printf(
				"Creating new credential entry (%v, %v) for exisintg user (%v, %v)",
				c.ID,
				profile.Provider,
				u.ID,
				u.Email,
			)
		} else {
			if !profile.ExpiresAt.IsZero() {
				// Copy expiration date when provided
				c.ExpiresAt = &profile.ExpiresAt
			}

			c.LastUsedAt = &lastUsedAt
			if c, err = svc.credentials.Update(c); err != nil {
				return err
			}

			log.Printf(
				"Updating credential entry (%v, %v) for exisintg user (%v, %v)",
				c.ID,
				profile.Provider,
				u.ID,
				u.Email,
			)
		}

		// Owner loaded, carry on.
		return nil
	})
}

// CheckPassword verifies username/password combination
//
// Expects plain text password as an input
func (svc *auth) CheckPassword(username, password string) (*types.User, error) {
	panic("svc.auth.CheckPassword, not implemented")
}

// ChangePassword (soft) deletes old password entry and creates a new one
//
// Expects plain text password as an input
func (svc *auth) ChangePassword(user *types.User, password string) error {
	panic("svc.auth.ChangePassword, not implemented")
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
