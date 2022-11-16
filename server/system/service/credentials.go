package service

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
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
		caProps.setUser(u)

		if !svc.ac.CanManageCredentialsOnUser(ctx, u) {
			return CredentialsErrNotAllowedToManage()
		}

		cc, _, err = store.SearchCredentials(ctx, svc.store, types.CredentialFilter{OwnerID: u.ID})
		return err
	}()

	return cc, svc.recordAction(ctx, caProps, CredentialsActionSearch, err)
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

		if !svc.ac.CanManageCredentialsOnUser(ctx, u) {
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
