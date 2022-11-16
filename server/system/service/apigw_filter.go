package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/apigw"
	agtypes "github.com/cortezaproject/corteza/server/pkg/apigw/types"
	a "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

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

func (svc *apigwFilter) FindByID(ctx context.Context, filterID uint64) (q *types.ApigwFilter, err error) {
	var (
		rProps = &apigwFilterActionProps{}
		r      *types.ApigwRoute
	)

	err = func() error {
		if filterID == 0 {
			return ApigwFilterErrInvalidID()
		}

		if q, err = store.LookupApigwFilterByID(ctx, svc.store, filterID); err != nil {
			return TemplateErrInvalidID().Wrap(err)
		}

		rProps.setFilter(q)

		// Get route
		if r, err = svc.route.FindByID(ctx, q.Route); err != nil {
			return err
		}

		if !svc.ac.CanReadApigwRoute(ctx, r) {
			return ApigwRouteErrNotAllowedToRead()
		}

		return nil
	}()

	return q, svc.recordAction(ctx, rProps, ApigwFilterActionLookup, err)
}

func (svc *apigwFilter) Create(ctx context.Context, new *types.ApigwFilter) (q *types.ApigwFilter, err error) {
	var (
		qProps = &apigwFilterActionProps{filter: new}
		r      *types.ApigwRoute
	)

	err = func() (err error) {
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

		q = new

		// send the signal to reload current route
		if r.Enabled {
			if err = apigw.Service().ReloadEndpoint(ctx, r.Method, r.Endpoint); err != nil {
				return err
			}
		}

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, ApigwFilterActionCreate, err)
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

func (svc *apigwFilter) Update(ctx context.Context, upd *types.ApigwFilter) (q *types.ApigwFilter, err error) {
	var (
		qProps = &apigwFilterActionProps{filter: upd}
		qq     *types.ApigwFilter
		r      *types.ApigwRoute
		e      error
	)

	err = func() (err error) {
		if qq, e = store.LookupApigwFilterByID(ctx, svc.store, upd.ID); e != nil {
			return ApigwFilterErrNotFound(qProps)
		}

		if r, err = svc.route.FindByID(ctx, upd.Route); err != nil {
			return err
		}

		if !svc.ac.CanUpdateApigwRoute(ctx, r) {
			return ApigwRouteErrNotAllowedToUpdate()
		}

		if qq, e = store.LookupApigwFilterByID(ctx, svc.store, upd.ID); e == nil && qq == nil {
			return ApigwFilterErrNotFound(qProps)
		}

		upd.UpdatedAt = now()
		upd.CreatedAt = qq.CreatedAt
		upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateApigwFilter(ctx, svc.store, upd); err != nil {
			return
		}

		q = upd

		// send the signal to reload current route
		if r.Enabled {
			if err = apigw.Service().ReloadEndpoint(ctx, r.Method, r.Endpoint); err != nil {
				return err
			}
		}

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, ApigwFilterActionUpdate, err)
}

func (svc *apigwFilter) DeleteByID(ctx context.Context, ID uint64) (err error) {
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
			return ApigwRouteErrNotAllowedToDelete()
		}

		qProps.setFilter(q)

		q.DeletedAt = now()
		q.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

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

func (svc *apigwFilter) Search(ctx context.Context, filter types.ApigwFilterFilter) (r types.ApigwFilterSet, f types.ApigwFilterFilter, err error) {
	var (
		aProps = &apigwFilterActionProps{search: &filter}
		route  *types.ApigwRoute
	)

	err = func() error {
		// Preload the corresponding API GW route for access control
		if filter.RouteID == 0 {
			return ApigwRouteErrInvalidID()
		}

		if route, err = store.LookupApigwRouteByID(ctx, svc.store, filter.RouteID); err != nil {
			return ApigwRouteErrNotFound()
		}

		if !svc.ac.CanReadApigwRoute(ctx, route) {
			return ApigwRouteErrNotAllowedToRead()
		}

		// Prepare the filter checker so we can evaluate access to specific filters
		filter.Check = func(res *types.ApigwFilter) (bool, error) {
			if !svc.ac.CanReadApigwRoute(ctx, route) {
				return false, nil
			}
			return true, nil
		}

		// Go!
		if r, f, err = store.SearchApigwFilters(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return r, f, svc.recordAction(ctx, aProps, ApigwFilterActionSearch, err)
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
