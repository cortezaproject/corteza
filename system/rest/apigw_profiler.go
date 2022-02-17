package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	profilerService interface {
		Hits(context.Context, types.ApigwProfilerFilter) (types.ApigwProfilerHitSet, types.ApigwProfilerFilter, error)
		HitsAggregated(context.Context, types.ApigwProfilerFilter) (types.ApigwProfilerAggregationSet, types.ApigwProfilerFilter, error)
	}

	ApigwProfiler struct {
		svc profilerService
		ac  templateAccessController
	}

	profilerRoutePayload struct {
		*types.ApigwProfilerHit
	}

	profilerHitPayload struct {
		*types.ApigwProfilerAggregation
	}

	profilerRouteSetPayload struct {
		Filter types.ApigwProfilerFilter `json:"filter"`
		Set    []*profilerRoutePayload   `json:"set"`
	}

	profilerHitSetPayload struct {
		Filter types.ApigwProfilerFilter `json:"filter"`
		Set    []*profilerHitPayload     `json:"set"`
	}
)

func (ApigwProfiler) New() *ApigwProfiler {
	return &ApigwProfiler{
		svc: service.DefaultApigwRoute,
		ac:  service.DefaultAccessControl,
	}
}

// List displays the the aggregated list of routes
func (ctrl *ApigwProfiler) Aggregation(ctx context.Context, r *request.ApigwProfilerAggregation) (interface{}, error) {
	var (
		err error
		f   = types.ApigwProfilerFilter{
			Path:   r.Path,
			Before: r.Before,
		}
	)

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	if f.Paging, err = filter.NewPaging(r.Limit, ""); err != nil {
		return nil, err
	}

	set, f, err := ctrl.svc.HitsAggregated(ctx, f)

	return ctrl.makeFilterPayload(ctx, set, f, err)
}

// Route displays the list of hits per-route
func (ctrl *ApigwProfiler) Route(ctx context.Context, r *request.ApigwProfilerRoute) (interface{}, error) {
	var (
		err error
		f   = types.ApigwProfilerFilter{
			Path:   r.RouteID,
			Before: r.Before,
		}
	)

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	if f.Paging, err = filter.NewPaging(r.Limit, ""); err != nil {
		return nil, err
	}

	set, f, err := ctrl.svc.Hits(ctx, f)

	return ctrl.makeRouteFilterPayload(ctx, set, f, err)
}

// Hit displays the details of a certain hit on a route
func (ctrl *ApigwProfiler) Hit(ctx context.Context, r *request.ApigwProfilerHit) (interface{}, error) {
	var (
		f = types.ApigwProfilerFilter{
			Hit: r.HitID,
		}
	)

	set, f, err := ctrl.svc.Hits(ctx, f)

	if len(set) != 1 {
		return nil, nil
	}

	return ctrl.makeRoutePayload(ctx, set[0], err)
}

func (ctrl *ApigwProfiler) makePayload(ctx context.Context, q *types.ApigwProfilerAggregation, err error) (*profilerHitPayload, error) {
	if err != nil || q == nil {
		return nil, err
	}

	qq := &profilerHitPayload{
		ApigwProfilerAggregation: q,
	}

	return qq, nil
}

func (ctrl *ApigwProfiler) makeFilterPayload(ctx context.Context, nn types.ApigwProfilerAggregationSet, f types.ApigwProfilerFilter, err error) (*profilerHitSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &profilerHitSetPayload{Filter: f, Set: make([]*profilerHitPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}

func (ctrl *ApigwProfiler) makeRoutePayload(ctx context.Context, q *types.ApigwProfilerHit, err error) (*profilerRoutePayload, error) {
	if err != nil || q == nil {
		return nil, err
	}

	return &profilerRoutePayload{q}, nil
}

func (ctrl *ApigwProfiler) makeRouteFilterPayload(ctx context.Context, nn types.ApigwProfilerHitSet, f types.ApigwProfilerFilter, err error) (*profilerRouteSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &profilerRouteSetPayload{Filter: f, Set: make([]*profilerRoutePayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makeRoutePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
