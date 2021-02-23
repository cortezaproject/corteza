package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	AuthClient struct {
		authClient authClientService
		ac         authClientAccessController
	}

	authClientService interface {
		LookupByID(ctx context.Context, ID uint64) (app *types.AuthClient, err error)
		Search(ctx context.Context, filter types.AuthClientFilter) (aa types.AuthClientSet, f types.AuthClientFilter, err error)
		Create(ctx context.Context, new *types.AuthClient) (app *types.AuthClient, err error)
		Update(ctx context.Context, upd *types.AuthClient) (app *types.AuthClient, err error)
		Delete(ctx context.Context, ID uint64) (err error)
		Undelete(ctx context.Context, ID uint64) (err error)
		ExposeSecret(ctx context.Context, ID uint64) (secret string, err error)
		RegenerateSecret(ctx context.Context, ID uint64) (secret string, err error)
	}

	authClientAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateAuthClient(context.Context, *types.AuthClient) bool
		CanDeleteAuthClient(context.Context, *types.AuthClient) bool
	}

	authClientPayload struct {
		*types.AuthClient

		CanGrant            bool `json:"canGrant"`
		CanUpdateAuthClient bool `json:"canUpdateAuthClient"`
		CanDeleteAuthClient bool `json:"canDeleteAuthClient"`
	}

	authClientSetPayload struct {
		Filter types.AuthClientFilter `json:"filter"`
		Set    []*authClientPayload   `json:"set"`
	}
)

func (AuthClient) New() *AuthClient {
	return &AuthClient{
		authClient: service.DefaultAuthClient,
		ac:         service.DefaultAccessControl,
	}
}

func (ctrl *AuthClient) List(ctx context.Context, r *request.AuthClientList) (interface{}, error) {
	var (
		err error
		f   = types.AuthClientFilter{
			Handle:  r.Handle,
			Labels:  r.Labels,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.authClient.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *AuthClient) Create(ctx context.Context, r *request.AuthClientCreate) (interface{}, error) {
	var (
		err error
		app = &types.AuthClient{
			Handle:      r.Handle,
			Meta:        r.Meta,
			ValidGrant:  r.ValidGrant,
			RedirectURI: r.RedirectURI,
			Scope:       r.Scope,
			Trusted:     r.Trusted,
			Enabled:     r.Enabled,
			ValidFrom:   r.ValidFrom,
			ExpiresAt:   r.ExpiresAt,
			Security:    r.Security,
			Labels:      r.Labels,
		}
	)

	app, err = ctrl.authClient.Create(ctx, app)
	return ctrl.makePayload(ctx, app, err)
}

func (ctrl *AuthClient) Update(ctx context.Context, r *request.AuthClientUpdate) (interface{}, error) {
	var (
		err error
		app = &types.AuthClient{
			ID:          r.ClientID,
			Handle:      r.Handle,
			Meta:        r.Meta,
			ValidGrant:  r.ValidGrant,
			RedirectURI: r.RedirectURI,
			Scope:       r.Scope,
			Trusted:     r.Trusted,
			Enabled:     r.Enabled,
			ValidFrom:   r.ValidFrom,
			ExpiresAt:   r.ExpiresAt,
			Security:    r.Security,
			Labels:      r.Labels,
		}
	)

	app, err = ctrl.authClient.Update(ctx, app)
	return ctrl.makePayload(ctx, app, err)
}

func (ctrl *AuthClient) Read(ctx context.Context, r *request.AuthClientRead) (interface{}, error) {
	app, err := ctrl.authClient.LookupByID(ctx, r.ClientID)
	return ctrl.makePayload(ctx, app, err)
}

func (ctrl *AuthClient) ExposeSecret(ctx context.Context, r *request.AuthClientExposeSecret) (interface{}, error) {
	return ctrl.authClient.ExposeSecret(ctx, r.ClientID)
}

func (ctrl *AuthClient) RegenerateSecret(ctx context.Context, r *request.AuthClientRegenerateSecret) (interface{}, error) {
	return ctrl.authClient.RegenerateSecret(ctx, r.ClientID)
}

func (ctrl *AuthClient) Delete(ctx context.Context, r *request.AuthClientDelete) (interface{}, error) {
	return api.OK(), ctrl.authClient.Delete(ctx, r.ClientID)
}

func (ctrl *AuthClient) Undelete(ctx context.Context, r *request.AuthClientUndelete) (interface{}, error) {
	return api.OK(), ctrl.authClient.Undelete(ctx, r.ClientID)
}

func (ctrl AuthClient) makePayload(ctx context.Context, m *types.AuthClient, err error) (*authClientPayload, error) {
	if err != nil || m == nil {
		return nil, err
	}

	return &authClientPayload{
		AuthClient: m,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateAuthClient: ctrl.ac.CanUpdateAuthClient(ctx, m),
		CanDeleteAuthClient: ctrl.ac.CanDeleteAuthClient(ctx, m),
	}, nil
}

func (ctrl AuthClient) makeFilterPayload(ctx context.Context, nn types.AuthClientSet, f types.AuthClientFilter, err error) (*authClientSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &authClientSetPayload{Filter: f, Set: make([]*authClientPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
