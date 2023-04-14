package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	credentials struct {
		actionlog actionlog.Recorder
		ac        credentialsAccessController
		store     store.Storer
	}

	credentialsAccessController interface {
		CanManageCredentialsOnUser(context.Context, *types.User) bool
	}
)

func Credentials() *credentials {
	return &credentials{
		ac:        DefaultAccessControl,
		store:     DefaultStore,
		actionlog: DefaultActionlog,
	}
}

func (svc *credentials) List(ctx context.Context, userID uint64) (cc types.CredentialSet, err error) {
	var (
		u *types.User

		caProps = &credentialsActionProps{user: &types.User{ID: userID}}
	)

	err = func() error {
		u, err = loadUser(ctx, svc.store, userID)
		if err != nil {
			return err
		}
		caProps.setUser(u)

		// Allow users to manage their own credentials
		ci := internalAuth.GetIdentityFromContext(ctx)
		if ci.Identity() != u.ID && !svc.ac.CanManageCredentialsOnUser(ctx, u) {
			return CredentialsErrNotAllowedToManage()
		}

		cc, _, err = store.SearchCredentials(ctx, svc.store, types.CredentialFilter{OwnerID: u.ID})
		return err
	}()

	return cc, svc.recordAction(ctx, caProps, CredentialsActionSearch, err)
}

func (svc *credentials) Create(ctx context.Context, c *types.Credential) (out *types.Credential, err error) {
	var (
		u       *types.User
		caProps = &credentialsActionProps{user: &types.User{ID: c.OwnerID}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		u, err = loadUser(ctx, svc.store, c.OwnerID)
		if err != nil {
			return
		}
		caProps.setUser(u)

		// Allow users to manage their own credentials
		ci := internalAuth.GetIdentityFromContext(ctx)
		if ci.Identity() != u.ID && !svc.ac.CanManageCredentialsOnUser(ctx, u) {
			return CredentialsErrNotAllowedToManage()
		}

		c.ID = nextID()
		c.CreatedAt = *now()

		err = store.CreateCredential(ctx, svc.store, c)
		return
	})

	return c, svc.recordAction(ctx, caProps, CredentialsActionCreate, err)
}

func (svc *credentials) Update(ctx context.Context, c *types.Credential) (out *types.Credential, err error) {
	var (
		u       *types.User
		caProps = &credentialsActionProps{user: &types.User{ID: c.OwnerID}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		u, err = loadUser(ctx, svc.store, c.OwnerID)
		if err != nil {
			return
		}
		caProps.setUser(u)

		// Allow users to manage their own credentials
		ci := internalAuth.GetIdentityFromContext(ctx)
		if ci.Identity() != u.ID && !svc.ac.CanManageCredentialsOnUser(ctx, u) {
			return CredentialsErrNotAllowedToManage()
		}

		old, err := store.LookupCredentialByID(ctx, s, c.ID)
		if err != nil {
			return
		}

		old.OwnerID = c.OwnerID
		old.Label = c.Label
		old.Kind = c.Kind
		old.Credentials = c.Credentials
		old.Meta = c.Meta
		old.LastUsedAt = c.LastUsedAt
		old.ExpiresAt = c.ExpiresAt
		old.CreatedAt = c.CreatedAt
		old.UpdatedAt = now()
		old.DeletedAt = c.DeletedAt

		err = store.UpdateCredential(ctx, s, old)

		return
	})

	return c, svc.recordAction(ctx, caProps, CredentialsActionUpdate, err)
}

func (svc *credentials) Delete(ctx context.Context, userID, credentialsID uint64) (err error) {
	var (
		c       *types.Credential
		u       *types.User
		caProps = &credentialsActionProps{user: &types.User{ID: userID}}
	)

	err = func() (err error) {
		if u, err = loadUser(ctx, svc.store, userID); err != nil {
			return
		}

		if u.Kind == types.SystemUser {
			return CredentialsErrNotAllowedToManage()
		}

		// Allow users to manage their own credentials
		ci := internalAuth.GetIdentityFromContext(ctx)
		if ci.Identity() != u.ID && !svc.ac.CanManageCredentialsOnUser(ctx, u) {
			return CredentialsErrNotAllowedToManage()
		}

		if c, err = store.LookupCredentialByID(ctx, svc.store, credentialsID); err != nil {
			return
		}

		caProps.setCredentials(c)
		c.DeletedAt = now()
		if err = store.UpdateCredential(ctx, svc.store, c); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, caProps, CredentialsActionDelete, err)
}
