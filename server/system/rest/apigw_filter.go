package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
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

func (ctrl *ApigwFilter) List(ctx context.Context, r *request.ApigwFilterList) (interface{}, error) {
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
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, f, err := ctrl.svc.Search(ctx, f)

	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl *ApigwFilter) Create(ctx context.Context, r *request.ApigwFilterCreate) (interface{}, error) {
	var (
		err error
		q   = &types.ApigwFilter{
			Route:   r.RouteID,
			Weight:  r.Weight,
			Kind:    r.Kind,
			Ref:     r.Ref,
			Enabled: r.Enabled,
			Params:  r.Params,
		}
	)

	q, err = ctrl.svc.Create(ctx, q)

	return ctrl.makePayload(ctx, q, err)
}

func (ctrl *ApigwFilter) Read(ctx context.Context, r *request.ApigwFilterRead) (interface{}, error) {
	return ctrl.svc.FindByID(ctx, r.FilterID)
}

func (ctrl *ApigwFilter) Update(ctx context.Context, r *request.ApigwFilterUpdate) (interface{}, error) {
	var (
		err error
		q   = &types.ApigwFilter{
			ID:      r.FilterID,
			Route:   r.RouteID,
			Weight:  r.Weight,
			Kind:    r.Kind,
			Ref:     r.Ref,
			Enabled: r.Enabled,
			Params:  r.Params,
		}
	)

	q, err = ctrl.svc.Update(ctx, q)

	return ctrl.makePayload(ctx, q, err)
}

func (ctrl *ApigwFilter) Delete(ctx context.Context, r *request.ApigwFilterDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.FilterID)
}

func (ctrl *ApigwFilter) DefFilter(ctx context.Context, r *request.ApigwFilterDefFilter) (interface{}, error) {
	return ctrl.svc.DefFilter(ctx, r.Kind)
}

func (ctrl *ApigwFilter) DefProxyAuth(ctx context.Context, r *request.ApigwFilterDefProxyAuth) (interface{}, error) {
	return ctrl.svc.DefProxyAuth(ctx)
}

func (ctrl *ApigwFilter) Undelete(ctx context.Context, r *request.ApigwFilterUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.FilterID)
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
