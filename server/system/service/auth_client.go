package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/pkg/rand"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service/event"
	"github.com/crusttech/human/server/system/types"
	oauth2def "github.com/go-oauth2/oauth2/v4"
)

type (
	authClientServices struct {
		eventbus eventDispatcher
		opt      options.AuthOpt
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
		services: &authClientServices{
			eventbus: eb,
			opt:      opt,
		},
	}
}

func (svc *authClient) validate(_ context.Context, new *types.AuthClient) error {
	if new.Meta == nil || new.Meta.Name == "" {
		return AuthClientErrMissingName()
	}

	if new.ValidGrant == oauth2def.ClientCredentials.String() {
		if new.Security == nil || new.Security.ImpersonateUser == 0 {
			return errors.Internal("auth client security configuration invalid")
		}
	}

	return nil
}

func (svc *authClient) beforeCreate(ctx context.Context, new *types.AuthClient) error {
	if err := svc.services.eventbus.WaitFor(ctx, event.AuthClientBeforeCreate(new, nil)); err != nil {
		return err
	}

	new.Secret = string(rand.Bytes(64))

	if new.Security == nil {
		new.Security = &types.AuthClientSecurity{}
	}

	if new.Meta == nil {
		new.Meta = &types.AuthClientMeta{}
	}

	return nil
}

func (svc *authClient) onRegenerateSecret(ctx context.Context, aProps *authClientActionProps, ID uint64) (secret string, err error) {
	var client *types.AuthClient

	if client, err = svc.lookupByID(ctx, ID); err != nil {
		return "", err
	}

	aProps.setAuthClient(client)

	secret = string(rand.Bytes(64))
	client.Secret = secret

	if err = store.UpdateAuthClient(ctx, svc.store, client); err != nil {
		return "", err
	}

	return secret, nil
}

func (svc *authClient) FindByID(ctx context.Context, ID uint64) (*types.AuthClient, error) {
	return svc.LookupByID(ctx, ID)
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

func (svc *authClient) IsDefaultClient(c *types.AuthClient) bool {
	if c == nil {
		return false
	}

	return c.Handle == svc.services.opt.DefaultClient
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
			// deleted auth clients visible only to admins; placeholder for future RBAC check
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
			a.Secret = ""
			return nil
		})

		return nil
	}()

	return aa, f, svc.recordAction(ctx, aaProps, AuthClientActionSearch, err)
}

func (svc *authClient) Update(ctx context.Context, upd *types.AuthClient) (res *types.AuthClient, err error) {
	var (
		aaProps                = &authClientActionProps{update: upd}
		defaultClientValidator = func(old, upd *types.AuthClient) error {
			if old.Handle != svc.services.opt.DefaultClient {
				return nil
			}

			if old.Handle != upd.Handle {
				return AuthClientErrUnableToChangeDefaultClientHandle()
			}

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
		if upd.Meta == nil || upd.Meta.Name == "" {
			return AuthClientErrMissingName()
		}

		if res, err = loadAuthClient(ctx, svc.store, upd.ID); err != nil {
			return
		}

		aaProps.setAuthClient(res)

		if !svc.ac.CanUpdateAuthClient(ctx, res) {
			return AuthClientErrNotAllowedToUpdate()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return AuthClientErrStaleData()
		}

		if err = defaultClientValidator(res, upd); err != nil {
			return err
		}

		if upd.ValidGrant == oauth2def.ClientCredentials.String() {
			if upd.Security == nil || upd.Security.ImpersonateUser == 0 {
				return errors.Internal("auth client security configuration invalid")
			}
		}

		if err = svc.services.eventbus.WaitFor(ctx, event.AuthClientBeforeUpdate(upd, res)); err != nil {
			return
		}

		if err = defaultClientValidator(res, upd); err != nil {
			return err
		}

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

		_ = svc.services.eventbus.WaitFor(ctx, event.AuthClientAfterUpdate(upd, res))
		return nil
	}()

	return res, svc.recordAction(ctx, aaProps, AuthClientActionUpdate, err)
}

func (svc *authClient) DeleteByID(ctx context.Context, ID uint64) (err error) {
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

		if res.Handle == svc.services.opt.DefaultClient {
			return AuthClientErrUnableToDeleteDefaultClient()
		}

		if err = svc.services.eventbus.WaitFor(ctx, event.AuthClientBeforeDelete(nil, res)); err != nil {
			return
		}

		res.DeletedAt = now()
		if err = store.UpdateAuthClient(ctx, svc.store, res); err != nil {
			return
		}

		_ = svc.services.eventbus.WaitFor(ctx, event.AuthClientAfterDelete(nil, res))
		return nil
	}()

	return svc.recordAction(ctx, aaProps, AuthClientActionDelete, err)
}
