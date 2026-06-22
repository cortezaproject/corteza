package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/apigw"
	agtypes "github.com/crusttech/human/server/pkg/apigw/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// The public CRUD skeleton (FindByID, Search, Create, Update, DeleteByID) is
// generated in apigw_filter.gen.go from system/apigw_filter.cue.
//
// Every op delegates to a hand-written on<Op> handler (customBodyOps) and skips
// the generated access check (customAccessOps): access control is enforced on
// the parent ApigwRoute, not the filter, and each op loads the route, validates
// and fires an endpoint-reload side effect.
//
// This file owns the struct, constructor, the on<Op> bodies (verbatim from the
// original service) and the resource-specific methods: validateAsyncRoute,
// UndeleteByID (kept hand-written -- it records ApigwFilterActionDelete, which
// the generated body cannot reproduce), DefFilter and DefProxyAuth.

type (
	apigwFilter struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        routeAccessController
		route     *apigwRoute
	}
)

func Filter() *apigwFilter {
	return &apigwFilter{
		route:     DefaultApigwRoute,
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

// onLookup is the custom body for the generated FindByID. The generated method
// owns the action-log scaffold + recordAction; access control is enforced on
// the parent route here.
func (svc *apigwFilter) onLookup(ctx context.Context, filterID uint64, rProps *apigwFilterActionProps) (q *types.ApigwFilter, err error) {
	var (
		r *types.ApigwRoute
	)

	if filterID == 0 {
		return nil, ApigwFilterErrInvalidID()
	}

	if q, err = store.LookupApigwFilterByID(ctx, svc.store, filterID); err != nil {
		return nil, TemplateErrInvalidID().Wrap(err)
	}

	rProps.setFilter(q)

	// Get route
	if r, err = svc.route.FindByID(ctx, q.Route); err != nil {
		return nil, err
	}

	if !svc.ac.CanReadApigwRoute(ctx, r) {
		return nil, ApigwRouteErrNotAllowedToRead()
	}

	return q, nil
}

// onCreate is the custom body for the generated Create. The generated method
// owns the action-log scaffold + recordAction; access control is enforced on
// the parent route here.
func (svc *apigwFilter) onCreate(ctx context.Context, new *types.ApigwFilter) (err error) {
	var (
		qProps = &apigwFilterActionProps{filter: new}
		r      *types.ApigwRoute
	)

	// Set new values after beforeCreate events are emitted
	new.ID = nextID()
	new.CreatedAt = *now()
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

	if r, err = store.LookupApigwRouteByID(ctx, svc.store, new.Route); err != nil {
		return
	}

	if !svc.ac.CanUpdateApigwRoute(ctx, r) {
		return ApigwRouteErrNotAllowedToUpdate()
	}

	// check for existing filters if route is async
	if r.Meta.Async {
		if err = svc.validateAsyncRoute(ctx, r, new, qProps); err != nil {
			return
		}
	}

	if err = store.CreateApigwFilter(ctx, svc.store, new); err != nil {
		return err
	}

	// send the signal to reload current route
	if r.Enabled {
		if err = apigw.Service().ReloadEndpoint(ctx, r.Method, r.Endpoint); err != nil {
			return err
		}
	}

	return nil
}

func (svc *apigwFilter) validateAsyncRoute(ctx context.Context, r *types.ApigwRoute, f *types.ApigwFilter, props *apigwFilterActionProps) (err error) {
	filters, _, err := svc.Search(ctx, types.ApigwFilterFilter{
		RouteID:  r.ID,
		Disabled: filter.StateExcluded,
	})

	if err != nil {
		return err
	}

	if f.Kind == string(agtypes.Processer) {
		processers, _ := filters.Filter(func(af *types.ApigwFilter) (bool, error) {
			return af.Kind == string(agtypes.Processer), nil
		})

		if len(processers) == 1 {
			return ApigwFilterErrAsyncRouteTooManyProcessers(props)
		}
	}

	if f.Kind == string(agtypes.PostFilter) {
		return ApigwFilterErrAsyncRouteTooManyAfterFilters(props)
	}

	return
}

// onUpdate is the custom body for the generated Update. The generated method
// owns the action-log scaffold + recordAction + load + tx; access control is
// enforced on the parent route here.
func (svc *apigwFilter) onUpdate(ctx context.Context, s store.Storer, upd, res *types.ApigwFilter, qProps *apigwFilterActionProps, _ func() error, _ func() error) (err error) {
	var r *types.ApigwRoute

	if r, err = svc.route.FindByID(ctx, upd.Route); err != nil {
		return err
	}

	if !svc.ac.CanUpdateApigwRoute(ctx, r) {
		return ApigwRouteErrNotAllowedToUpdate()
	}

	upd.UpdatedAt = now()
	upd.CreatedAt = res.CreatedAt
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	if err = store.UpdateApigwFilter(ctx, s, upd); err != nil {
		return err
	}

	// Reflect updated values back into res so the caller sees the final state.
	*res = *upd

	// send the signal to reload current route
	if r.Enabled {
		if err = apigw.Service().ReloadEndpoint(ctx, r.Method, r.Endpoint); err != nil {
			return err
		}
	}

	return nil
}

// onDelete is the custom body for the generated DeleteByID. The generated
// method owns the action-log scaffold + recordAction + load + tx; access
// control is enforced on the parent route here.
func (svc *apigwFilter) onDelete(ctx context.Context, s store.Storer, res *types.ApigwFilter, qProps *apigwFilterActionProps) (err error) {
	var r *types.ApigwRoute

	if r, err = store.LookupApigwRouteByID(ctx, s, res.Route); err == store.ErrNotFound {
		return ApigwRouteErrNotFound()
	} else if err != nil {
		return
	}

	if !svc.ac.CanDeleteApigwRoute(ctx, r) {
		return ApigwRouteErrNotAllowedToDelete()
	}

	res.DeletedAt = now()
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

	if err = store.UpdateApigwFilter(ctx, s, res); err != nil {
		return
	}

	// send the signal to reload current route
	if r.Enabled {
		if err = apigw.Service().ReloadEndpoint(ctx, r.Method, r.Endpoint); err != nil {
			return err
		}
	}

	return nil
}

func (svc *apigwFilter) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &apigwFilterActionProps{}
		q      *types.ApigwFilter
		r      *types.ApigwRoute
	)

	err = func() (err error) {
		if q, err = store.LookupApigwFilterByID(ctx, svc.store, ID); err != nil {
			return ApigwFilterErrNotFound(qProps)
		}

		if r, err = store.LookupApigwRouteByID(ctx, svc.store, q.Route); err == store.ErrNotFound {
			return ApigwRouteErrNotFound()
		} else if err != nil {
			return
		}

		if !svc.ac.CanDeleteApigwRoute(ctx, r) {
			return ApigwRouteErrNotAllowedToUndelete()
		}

		qProps.setFilter(q)

		q.DeletedAt = nil
		q.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateApigwFilter(ctx, svc.store, q); err != nil {
			return
		}

		// send the signal to reload current route
		if r.Enabled {
			if err = apigw.Service().ReloadEndpoint(ctx, r.Method, r.Endpoint); err != nil {
				return err
			}
		}

		return nil
	}()

	return svc.recordAction(ctx, qProps, ApigwFilterActionDelete, err)
}

// onSearch is the custom body for the generated Search. The generated method
// owns the action-log scaffold + recordAction; access control is enforced on
// the parent route here.
func (svc *apigwFilter) onSearch(ctx context.Context, f types.ApigwFilterFilter, aProps *apigwFilterActionProps) (r types.ApigwFilterSet, outFilter types.ApigwFilterFilter, err error) {
	var (
		route *types.ApigwRoute
	)

	// Preload the corresponding API GW route for access control
	if f.RouteID == 0 {
		return nil, f, ApigwRouteErrInvalidID()
	}

	if route, err = store.LookupApigwRouteByID(ctx, svc.store, f.RouteID); err != nil {
		return nil, f, ApigwRouteErrNotFound()
	}

	if !svc.ac.CanReadApigwRoute(ctx, route) {
		return nil, f, ApigwRouteErrNotAllowedToRead()
	}

	// Prepare the filter checker so we can evaluate access to specific filters
	f.Check = func(res *types.ApigwFilter) (bool, error) {
		if !svc.ac.CanReadApigwRoute(ctx, route) {
			return false, nil
		}
		return true, nil
	}

	// Go!
	if r, outFilter, err = store.SearchApigwFilters(ctx, svc.store, f); err != nil {
		return nil, f, err
	}

	return r, outFilter, nil
}

func (svc *apigwFilter) DefFilter(ctx context.Context, kind string) (l interface{}, err error) {
	var (
		qProps = &apigwFilterActionProps{}
	)

	err = func() error {
		if !svc.ac.CanSearchApigwRoutes(ctx) {
			return ApigwRouteErrNotAllowedToRead()
		}
		// get the definitions from registry
		l = apigw.Service().Funcs(kind)

		return nil
	}()

	return l, svc.recordAction(ctx, qProps, ApigwFilterActionSearch, err)

}

func (svc *apigwFilter) DefProxyAuth(ctx context.Context) (l interface{}, err error) {
	var (
		qProps = &apigwFilterActionProps{}
	)

	err = func() error {
		if !svc.ac.CanSearchApigwRoutes(ctx) {
			return ApigwRouteErrNotAllowedToRead()
		}
		// get the definitions from registry
		l = apigw.Service().ProxyAuthDef()

		return nil
	}()

	return l, svc.recordAction(ctx, qProps, ApigwFilterActionSearch, err)

}
