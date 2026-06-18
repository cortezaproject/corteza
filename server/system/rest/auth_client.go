package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	AuthClient struct {
		authClient authClientService
		ac         authClientAccessController
		opt        options.AuthOpt
	}

	authClientService interface {
		FindByID(ctx context.Context, ID uint64) (app *types.AuthClient, err error)
		Search(ctx context.Context, filter types.AuthClientFilter) (aa types.AuthClientSet, f types.AuthClientFilter, err error)
		Create(ctx context.Context, new *types.AuthClient) (app *types.AuthClient, err error)
		Update(ctx context.Context, upd *types.AuthClient) (app *types.AuthClient, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		UndeleteByID(ctx context.Context, ID uint64) (err error)
		ExposeSecret(ctx context.Context, ID uint64) (secret string, err error)
		RegenerateSecret(ctx context.Context, ID uint64) (secret string, err error)
		IsDefaultClient(c *types.AuthClient) bool
	}

	authClientAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateAuthClient(context.Context, *types.AuthClient) bool
		CanDeleteAuthClient(context.Context, *types.AuthClient) bool
	}

	authClientPayload struct {
		*types.AuthClient

		IsDefault bool `json:"isDefault"`

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

// makeFilter builds the search filter for the generated List controller.
func (ctrl AuthClient) makeFilter(ctx context.Context, r *request.AuthClientList) (types.AuthClientFilter, error) {
	var (
		err error
		f   = types.AuthClientFilter{
			Handle:  r.Handle,
			Labels:  r.Labels,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params.
func (ctrl AuthClient) beforeCreate(ctx context.Context, res *types.AuthClient, r *request.AuthClientCreate) error {
	res.Meta = r.Meta
	res.Security = r.Security
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl AuthClient) beforeUpdate(ctx context.Context, res *types.AuthClient, r *request.AuthClientUpdate) error {
	res.Meta = r.Meta
	res.Security = r.Security
	return nil
}

func (ctrl *AuthClient) ExposeSecret(ctx context.Context, r *request.AuthClientExposeSecret) (interface{}, error) {
	return ctrl.authClient.ExposeSecret(ctx, r.ClientID)
}

func (ctrl *AuthClient) RegenerateSecret(ctx context.Context, r *request.AuthClientRegenerateSecret) (interface{}, error) {
	return ctrl.authClient.RegenerateSecret(ctx, r.ClientID)
}

func (ctrl AuthClient) makePayload(ctx context.Context, m *types.AuthClient, err error) (*authClientPayload, error) {
	if err != nil || m == nil {
		return nil, err
	}

	return &authClientPayload{
		AuthClient: m,

		IsDefault: ctrl.authClient.IsDefaultClient(m),

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
