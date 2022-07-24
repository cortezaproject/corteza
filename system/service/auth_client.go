package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	oauth2def "github.com/go-oauth2/oauth2/v4"
)

type (
	authClient struct {
		ac        authClientAccessController
		eventbus  eventDispatcher
		actionlog actionlog.Recorder
		store     store.Storer
		opt       options.AuthOpt
	}

	authClientAccessController interface {
		CanSearchAuthClients(context.Context) bool
		CanCreateAuthClient(context.Context) bool
		CanReadAuthClient(context.Context, *types.AuthClient) bool
		CanUpdateAuthClient(context.Context, *types.AuthClient) bool
		CanDeleteAuthClient(context.Context, *types.AuthClient) bool
	}
)

// AuthClient is a default authClient service initializer
func AuthClient(s store.Storer, ac authClientAccessController, al actionlog.Recorder, eb eventDispatcher, opt options.AuthOpt) *authClient {
	return &authClient{
		store:     s,
		ac:        ac,
		actionlog: al,
		eventbus:  eb,
		opt:       opt,
	}
}

func (svc *authClient) LookupByID(ctx context.Context, ID uint64) (client *types.AuthClient, err error) {
	var (
		aaProps = &authClientActionProps{authClient: &types.AuthClient{ID: ID}}
	)

	client, err = svc.lookupByID(ctx, ID)

	if client != nil {
		client.Secret = ""
	}

	return client, svc.recordAction(ctx, aaProps, AuthClientActionLookup, err)
}

func (svc *authClient) ExposeSecret(ctx context.Context, ID uint64) (secret string, err error) {
	var (
		client  *types.AuthClient
		aaProps = &authClientActionProps{authClient: &types.AuthClient{ID: ID}}
	)

	client, err = svc.lookupByID(ctx, ID)
	if client != nil {
		secret = client.Secret
	}

	return secret, svc.recordAction(ctx, aaProps, AuthClientActionExposeSecret, err)
}

func (svc *authClient) RegenerateSecret(ctx context.Context, ID uint64) (secret string, err error) {
	var (
		client  *types.AuthClient
		aaProps = &authClientActionProps{authClient: &types.AuthClient{ID: ID}}
	)

	client, err = svc.lookupByID(ctx, ID)
	if client != nil {
		secret = string(rand.Bytes(64))
		client.Secret = secret
		err = store.UpdateAuthClient(ctx, svc.store, client)
	}

	return secret, svc.recordAction(ctx, aaProps, AuthClientActionRegenerateSecret, err)
}

func (svc *authClient) IsDefaultClient(c *types.AuthClient) bool {
	if c == nil {
		return false
	}

	return c.Handle == svc.opt.DefaultClient
}

func (svc *authClient) lookupByID(ctx context.Context, ID uint64) (client *types.AuthClient, err error) {
	err = func() error {
		if client, err = loadAuthClient(ctx, svc.store, ID); err != nil {
			return AuthClientErrInvalidID().Wrap(err)
		}

		if !svc.ac.CanReadAuthClient(ctx, client) {
			return AuthClientErrNotAllowedToRead()
		}

		return nil
	}()

	return client, err
}

func (svc *authClient) Search(ctx context.Context, af types.AuthClientFilter) (aa types.AuthClientSet, f types.AuthClientFilter, err error) {
	var (
		aaProps = &authClientActionProps{filter: &af}
	)

	// For each fetched item, store backend will check if it is valid or not
	af.Check = func(res *types.AuthClient) (bool, error) {
		if !svc.ac.CanReadAuthClient(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchAuthClients(ctx) {
			return AuthClientErrNotAllowedToSearch()
		}

		if af.Deleted > filter.StateExcluded {
			// If list with deleted authClients is requested
			// user must have access permissions to system (ie: is admin)
			//
			// not the best solution but ATM it allows us to have at least
			// some kind of control over who can see deleted authClients
			//if !svc.ac.CanAccess(ctx) {
			//	return AuthClientErrNotAllowedToListAuthClients()
			//}
		}

		if len(af.Labels) > 0 {
			af.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.AuthClient{}.LabelResourceKind(),
				af.Labels,
			)

			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(af.LabeledIDs) == 0 {
				return nil
			}
		}

		if aa, f, err = store.SearchAuthClients(ctx, svc.store, af); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledAuthClients(aa)...); err != nil {
			return err
		}

		_ = aa.Walk(func(a *types.AuthClient) error {
			// make sure we do not leak client's secret without explicit request
			a.Secret = ""
			return nil
		})

		return nil

	}()

	return aa, f, svc.recordAction(ctx, aaProps, AuthClientActionSearch, err)
}

func (svc *authClient) Create(ctx context.Context, new *types.AuthClient) (res *types.AuthClient, err error) {
	var (
		aaProps = &authClientActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateAuthClient(ctx) {
			return AuthClientErrNotAllowedToCreate()
		}

		if err = svc.eventbus.WaitFor(ctx, event.AuthClientBeforeCreate(new, nil)); err != nil {
			return
		}

		// Set new values after beforeCreate events are emitted
		new.ID = nextID()
		new.CreatedAt = *now()
		new.Secret = string(rand.Bytes(64))

		if new.Security == nil {
			new.Security = &types.AuthClientSecurity{}
		}

		if new.Meta == nil {
			new.Meta = &types.AuthClientMeta{}
		}

		// Validate impersonated user
		if new.ValidGrant == oauth2def.ClientCredentials.String() {
			if new.Security == nil || new.Security.ImpersonateUser == 0 {
				return errors.Internal("auth client security configuration invalid")
			}
		}

		if err = store.CreateAuthClient(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		res = new

		_ = svc.eventbus.WaitFor(ctx, event.AuthClientAfterCreate(new, nil))
		return nil
	}()

	return res, svc.recordAction(ctx, aaProps, AuthClientActionCreate, err)
}

func (svc *authClient) Update(ctx context.Context, upd *types.AuthClient) (res *types.AuthClient, err error) {
	var (
		aaProps                = &authClientActionProps{update: upd}
		defaultClientValidator = func(old, upd *types.AuthClient) error {
			if old.Handle != svc.opt.DefaultClient {
				return nil
			}

			// The handle may not change
			if old.Handle != upd.Handle {
				return AuthClientErrUnableToChangeDefaultClientHandle()
			}

			// The client may not get disabled
			if !upd.Enabled {
				return AuthClientErrUnableToDisableDefaultClient()
			}

			return nil
		}
	)

	err = func() (err error) {
		if upd.ID == 0 {
			return AuthClientErrInvalidID()
		}

		if res, err = loadAuthClient(ctx, svc.store, upd.ID); err != nil {
			return
		}

		aaProps.setAuthClient(res)

		if !svc.ac.CanUpdateAuthClient(ctx, res) {
			return AuthClientErrNotAllowedToUpdate()
		}

		// Firstly validate default clients before the automation occurs
		if err = defaultClientValidator(res, upd); err != nil {
			return err
		}

		// Validate impersonated user
		if upd.ValidGrant == oauth2def.ClientCredentials.String() {
			if upd.Security == nil || upd.Security.ImpersonateUser == 0 {
				return errors.Internal("auth client security configuration invalid")
			}
		}

		if err = svc.eventbus.WaitFor(ctx, event.AuthClientBeforeUpdate(upd, res)); err != nil {
			return
		}

		// Next validate default clients after the automation occurs
		if err = defaultClientValidator(res, upd); err != nil {
			return err
		}

		// Assign changed values after afterUpdate events are emitted
		res.Handle = upd.Handle
		res.ValidGrant = upd.ValidGrant
		res.RedirectURI = upd.RedirectURI
		res.Scope = upd.Scope
		res.Enabled = upd.Enabled
		res.Trusted = upd.Trusted
		res.ValidFrom = upd.ValidFrom
		res.ExpiresAt = upd.ExpiresAt
		res.UpdatedAt = now()

		if upd.Meta != nil {
			res.Meta = upd.Meta
		}

		if upd.Security != nil {
			res.Security = upd.Security
		}

		if err = store.UpdateAuthClient(ctx, svc.store, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}

		_ = svc.eventbus.WaitFor(ctx, event.AuthClientAfterUpdate(upd, res))
		return nil
	}()

	return res, svc.recordAction(ctx, aaProps, AuthClientActionUpdate, err)
}

func (svc *authClient) Delete(ctx context.Context, ID uint64) (err error) {
	var (
		aaProps = &authClientActionProps{}
		res     *types.AuthClient
	)

	err = func() (err error) {
		if res, err = loadAuthClient(ctx, svc.store, ID); err != nil {
			return
		}

		aaProps.setAuthClient(res)

		if !svc.ac.CanDeleteAuthClient(ctx, res) {
			return AuthClientErrNotAllowedToDelete()
		}

		if res.Handle == svc.opt.DefaultClient {
			return AuthClientErrUnableToDeleteDefaultClient()
		}

		if err = svc.eventbus.WaitFor(ctx, event.AuthClientBeforeDelete(nil, res)); err != nil {
			return
		}

		res.DeletedAt = now()
		if err = store.UpdateAuthClient(ctx, svc.store, res); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.AuthClientAfterDelete(nil, res))
		return nil
	}()

	return svc.recordAction(ctx, aaProps, AuthClientActionDelete, err)
}

func (svc *authClient) Undelete(ctx context.Context, ID uint64) (err error) {
	var (
		aaProps = &authClientActionProps{}
		res     *types.AuthClient
	)

	err = func() (err error) {
		if res, err = loadAuthClient(ctx, svc.store, ID); err != nil {
			return
		}

		aaProps.setAuthClient(res)

		if !svc.ac.CanDeleteAuthClient(ctx, res) {
			return AuthClientErrNotAllowedToUndelete()
		}

		// @todo add event
		//       if err = svc.eventbus.WaitFor(ctx, event.AuthClientBeforeUndelete(nil, res)); err != nil {
		//       	return
		//       }

		res.DeletedAt = nil
		if err = store.UpdateAuthClient(ctx, svc.store, res); err != nil {
			return
		}

		// @todo add event
		//       _ = svc.eventbus.WaitFor(ctx, event.AuthClientAfterUndelete(nil, res))
		return nil
	}()

	return svc.recordAction(ctx, aaProps, AuthClientActionUndelete, err)
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
//
// This function is auto-generated.
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
