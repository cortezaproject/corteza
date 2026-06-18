package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ApigwRoute struct {
		svc routeService
		ac  apiGwRouteAccessController
	}

	routePayload struct {
		*types.ApigwRoute

		CanGrant            bool `json:"canGrant"`
		CanUpdateApigwRoute bool `json:"canUpdateApigwRoute"`
		CanDeleteApigwRoute bool `json:"canDeleteApigwRoute"`
	}

	routeSetPayload struct {
		Filter types.ApigwRouteFilter `json:"filter"`
		Set    []*routePayload        `json:"set"`
	}

	routeService interface {
		FindByID(ctx context.Context, ID uint64) (*types.ApigwRoute, error)
		Create(ctx context.Context, new *types.ApigwRoute) (*types.ApigwRoute, error)
		Update(ctx context.Context, upd *types.ApigwRoute) (*types.ApigwRoute, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.ApigwRouteFilter) (types.ApigwRouteSet, types.ApigwRouteFilter, error)
	}

	apiGwRouteAccessController interface {
		CanGrant(context.Context) bool

		CanCreateApigwRoute(context.Context) bool
		CanUpdateApigwRoute(context.Context, *types.ApigwRoute) bool
		CanDeleteApigwRoute(context.Context, *types.ApigwRoute) bool
	}
)

func (ApigwRoute) New() *ApigwRoute {
	return &ApigwRoute{
		svc: service.DefaultApigwRoute,
		ac:  service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl *ApigwRoute) makeFilter(ctx context.Context, r *request.ApigwRouteList) (types.ApigwRouteFilter, error) {
	var (
		err error
		f   = types.ApigwRouteFilter{
			// todo: this should renamed to r.Endpoint after UI is aligned with this
			Endpoint: r.Query,
			Deleted:  filter.State(r.Deleted),

			// todo: this should dynamic as Delete
			//		but making it default to `1`, until UI is aligned with this
			Disabled: filter.StateInclusive,
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
func (ctrl *ApigwRoute) beforeCreate(ctx context.Context, res *types.ApigwRoute, r *request.ApigwRouteCreate) error {
	res.Meta = r.Meta
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl *ApigwRoute) beforeUpdate(ctx context.Context, res *types.ApigwRoute, r *request.ApigwRouteUpdate) error {
	res.Meta = r.Meta
	return nil
}

func (ctrl *ApigwRoute) makePayload(ctx context.Context, q *types.ApigwRoute, err error) (*routePayload, error) {
	if err != nil || q == nil {
		return nil, err
	}

	qq := &routePayload{
		ApigwRoute: q,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateApigwRoute: ctrl.ac.CanUpdateApigwRoute(ctx, q),
		CanDeleteApigwRoute: ctrl.ac.CanDeleteApigwRoute(ctx, q),
	}

	return qq, nil
}

func (ctrl *ApigwRoute) makeFilterPayload(ctx context.Context, nn types.ApigwRouteSet, f types.ApigwRouteFilter, err error) (*routeSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &routeSetPayload{Filter: f, Set: make([]*routePayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
