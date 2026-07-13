package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/apigw"
	agtypes "github.com/crusttech/human/server/pkg/apigw/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	apigwFilterServices struct {
		route *apigwRoute
	}

	apigwFilterAccessController interface {
		CanSearchApigwRoutes(ctx context.Context) bool
		CanReadApigwRoute(context.Context, *types.ApigwRoute) bool
		CanUpdateApigwRoute(context.Context, *types.ApigwRoute) bool
		CanDeleteApigwRoute(context.Context, *types.ApigwRoute) bool
	}
)

func Filter() *apigwFilter {
	return &apigwFilter{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		services:  &apigwFilterServices{route: DefaultApigwRoute},
	}
}

func (svc *apigwFilter) onLookup(ctx context.Context, filterID uint64, aProps *apigwFilterActionProps) (*types.ApigwFilter, error) {
	var (
		q *types.ApigwFilter
		r *types.ApigwRoute
	)

	if filterID == 0 {
		return nil, ApigwFilterErrInvalidID()
	}

	var err error
	if q, err = store.LookupApigwFilterByID(ctx, svc.store, filterID); err != nil {
		return nil, TemplateErrInvalidID().Wrap(err)
	}

	aProps.setFilter(q)

	if r, err = svc.services.route.FindByID(ctx, q.Route); err != nil {
		return nil, err
	}

	if !svc.ac.CanReadApigwRoute(ctx, r) {
		return nil, ApigwRouteErrNotAllowedToRead()
	}

	return q, nil
}

func (svc *apigwFilter) onSearch(ctx context.Context, f types.ApigwFilterFilter, aProps *apigwFilterActionProps) (types.ApigwFilterSet, types.ApigwFilterFilter, error) {
	var (
		r     types.ApigwFilterSet
		route *types.ApigwRoute
	)

	if f.RouteID == 0 {
		return nil, f, ApigwRouteErrInvalidID()
	}

	var err error
	if route, err = store.LookupApigwRouteByID(ctx, svc.store, f.RouteID); err != nil {
		return nil, f, ApigwRouteErrNotFound()
	}

	if !svc.ac.CanReadApigwRoute(ctx, route) {
		return nil, f, ApigwRouteErrNotAllowedToRead()
	}

	f.Check = func(res *types.ApigwFilter) (bool, error) {
		if !svc.ac.CanReadApigwRoute(ctx, route) {
			return false, nil
		}
		return true, nil
	}

	if r, f, err = store.SearchApigwFilters(ctx, svc.store, f); err != nil {
		return nil, f, err
	}

	return r, f, nil
}

func (svc *apigwFilter) onCreate(ctx context.Context, new *types.ApigwFilter) error {
	new.ID = nextID()
	new.CreatedAt = *now()
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

	r, err := store.LookupApigwRouteByID(ctx, svc.store, new.Route)
	if err != nil {
		return err
	}

	if !svc.ac.CanUpdateApigwRoute(ctx, r) {
		return ApigwRouteErrNotAllowedToUpdate()
	}

	if r.Meta.Async {
		if err = svc.validateAsyncRoute(ctx, r, new, &apigwFilterActionProps{filter: new}); err != nil {
			return err
		}
	}

	if err = store.CreateApigwFilter(ctx, svc.store, new); err != nil {
		return err
	}

	if r.Enabled {
		if err = apigw.Service().ReloadEndpoint(ctx, r.Method, r.Endpoint); err != nil {
			return err
		}
	}

	return nil
}

func (svc *apigwFilter) onUpdate(ctx context.Context, s store.Storer, upd *types.ApigwFilter, res *types.ApigwFilter, aProps *apigwFilterActionProps, before func() error, after func() error) error {
	r, err := svc.services.route.FindByID(ctx, upd.Route)
	if err != nil {
		return err
	}

	if !svc.ac.CanUpdateApigwRoute(ctx, r) {
		return ApigwRouteErrNotAllowedToUpdate()
	}

	if err = before(); err != nil {
		return err
	}

	if err = after(); err != nil {
		return err
	}

	if r.Enabled {
		if err = apigw.Service().ReloadEndpoint(ctx, r.Method, r.Endpoint); err != nil {
			return err
		}
	}

	return nil
}

func (svc *apigwFilter) onDelete(ctx context.Context, s store.Storer, res *types.ApigwFilter, aProps *apigwFilterActionProps) error {
	r, err := store.LookupApigwRouteByID(ctx, s, res.Route)
	if err == store.ErrNotFound {
		return ApigwRouteErrNotFound()
	} else if err != nil {
		return err
	}

	if !svc.ac.CanDeleteApigwRoute(ctx, r) {
		return ApigwRouteErrNotAllowedToDelete()
	}

	res.DeletedAt = now()
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

	if err = store.UpdateApigwFilter(ctx, s, res); err != nil {
		return err
	}

	if r.Enabled {
		if err = apigw.Service().ReloadEndpoint(ctx, r.Method, r.Endpoint); err != nil {
			return err
		}
	}

	return nil
}

func (svc *apigwFilter) onDefFilter(ctx context.Context, aProps *apigwFilterActionProps, kind string) (interface{}, error) {
	if !svc.ac.CanSearchApigwRoutes(ctx) {
		return nil, ApigwRouteErrNotAllowedToRead()
	}
	return apigw.Service().Funcs(kind), nil
}

func (svc *apigwFilter) onDefProxyAuth(ctx context.Context, aProps *apigwFilterActionProps) (interface{}, error) {
	if !svc.ac.CanSearchApigwRoutes(ctx) {
		return nil, ApigwRouteErrNotAllowedToRead()
	}
	return apigw.Service().ProxyAuthDef(), nil
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

		if r.Enabled {
			if err = apigw.Service().ReloadEndpoint(ctx, r.Method, r.Endpoint); err != nil {
				return err
			}
		}

		return nil
	}()

	return svc.recordAction(ctx, qProps, ApigwFilterActionDelete, err)
}
