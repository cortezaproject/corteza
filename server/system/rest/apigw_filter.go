package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ApigwFilter struct {
		svc filterService
		ac  templateAccessController
	}

	functionPayload struct {
		*types.ApigwFilter
	}

	functionSetPayload struct {
		Filter types.ApigwFilterFilter `json:"filter"`
		Set    []*functionPayload      `json:"set"`
	}

	filterService interface {
		FindByID(ctx context.Context, ID uint64) (*types.ApigwFilter, error)
		Search(ctx context.Context, filter types.ApigwFilterFilter) (types.ApigwFilterSet, types.ApigwFilterFilter, error)
		Create(ctx context.Context, new *types.ApigwFilter) (*types.ApigwFilter, error)
		Update(ctx context.Context, upd *types.ApigwFilter) (*types.ApigwFilter, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error

		DefFilter(context.Context, string) (interface{}, error)
		DefProxyAuth(context.Context) (interface{}, error)
	}
)

func (ApigwFilter) New() *ApigwFilter {
	return &ApigwFilter{
		svc: service.DefaultApigwFilter,
		ac:  service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl *ApigwFilter) makeFilter(ctx context.Context, r *request.ApigwFilterList) (types.ApigwFilterFilter, error) {
	var (
		err error
		f   = types.ApigwFilterFilter{
			RouteID: r.RouteID,
			Deleted: filter.State(r.Deleted),

			// todo: this should dynamic as Delete
			//		but making it default to `1`, until UI is aligned with this
			Disabled: filter.StateInclusive,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params.
func (ctrl *ApigwFilter) beforeCreate(ctx context.Context, res *types.ApigwFilter, r *request.ApigwFilterCreate) error {
	res.Route = r.RouteID
	res.Params = r.Params
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl *ApigwFilter) beforeUpdate(ctx context.Context, res *types.ApigwFilter, r *request.ApigwFilterUpdate) error {
	res.Route = r.RouteID
	res.Params = r.Params
	return nil
}

func (ctrl *ApigwFilter) DefFilter(ctx context.Context, r *request.ApigwFilterDefFilter) (interface{}, error) {
	return ctrl.svc.DefFilter(ctx, r.Kind)
}

func (ctrl *ApigwFilter) DefProxyAuth(ctx context.Context, r *request.ApigwFilterDefProxyAuth) (interface{}, error) {
	return ctrl.svc.DefProxyAuth(ctx)
}

func (ctrl *ApigwFilter) makePayload(ctx context.Context, q *types.ApigwFilter, err error) (*functionPayload, error) {
	if err != nil || q == nil {
		return nil, err
	}

	qq := &functionPayload{
		ApigwFilter: q,
	}

	return qq, nil
}

func (ctrl *ApigwFilter) makeFilterPayload(ctx context.Context, nn types.ApigwFilterSet, f types.ApigwFilterFilter, err error) (*functionSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &functionSetPayload{Filter: f, Set: make([]*functionPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
