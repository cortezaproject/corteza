package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type authClientServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *authClient) Create(ctx context.Context, new *types.AuthClient) (res *types.AuthClient, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &authClientActionProps{authClient: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if err = svc.validate(ctx, new); err != nil {
			return err
		}
		if !svc.ac.CanCreateAuthClient(ctx) {
			return AuthClientErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateAuthClient(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, AuthClientActionCreate, err)
}

func (svc *authClient) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &authClientActionProps{}
		res    *types.AuthClient
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadAuthClient(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setAuthClient(res)

		if !svc.ac.CanDeleteAuthClient(ctx, res) {
			return AuthClientErrNotAllowedToUndelete()
		}

		res.DeletedAt = nil
		if err = store.UpdateAuthClient(ctx, svc.store, res); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, aProps, AuthClientActionUndelete, err)
}

func loadAuthClient(ctx context.Context, s store.AuthClients, ID uint64) (res *types.AuthClient, err error) {
	if ID == 0 {
		return nil, AuthClientErrInvalidID()
	}

	if res, err = store.LookupAuthClientByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, AuthClientErrNotFound()
	}

	return
}

// toLabeledAuthClients converts to []label.LabeledResource
func toLabeledAuthClients(set []*types.AuthClient) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}

func (svc *authClient) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *authClient) scopeServices(ctx context.Context) *authClientServices {
	return &authClientServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}

func (svc *authClient) RegenerateSecret(ctx context.Context, ID uint64) (secret string, err error) {
	var (
		aProps = &authClientActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		secret, err = svc.onRegenerateSecret(ctx, aProps, ID)
		return err
	}()

	return secret, svc.recordAction(ctx, aProps, AuthClientActionRegenerateSecret, err)
}
