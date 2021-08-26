package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/apigw"
	agtypes "github.com/cortezaproject/corteza-server/pkg/apigw/types"
	a "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	apigwFilter struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        functionAccessController
		route     *apigwRoute
	}

	functionAccessController interface {
		CanSearchApigwFilters(context.Context) bool

		CanCreateApigwFilter(context.Context) bool
		CanReadApigwFilter(context.Context, *types.ApigwFilter) bool
		CanUpdateApigwFilter(context.Context, *types.ApigwFilter) bool
		CanDeleteApigwFilter(context.Context, *types.ApigwFilter) bool
	}
)

func Filter() *apigwFilter {
	return (&apigwFilter{
		route:     DefaultApigwRoute,
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	})
}

func (svc *apigwFilter) FindByID(ctx context.Context, ID uint64) (q *types.ApigwFilter, err error) {
	var (
		rProps = &apigwFilterActionProps{}
	)

	err = func() error {
		if ID == 0 {
			return ApigwFilterErrInvalidID()
		}

		if !svc.ac.CanSearchApigwFilters(ctx) {
			return ApigwFilterErrNotAllowedToRead(rProps)
		}

		if q, err = store.LookupApigwFilterByID(ctx, svc.store, ID); err != nil {
			return TemplateErrInvalidID().Wrap(err)
		}

		rProps.setFilter(q)

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
		if !svc.ac.CanCreateApigwFilter(ctx) {
			return ApigwFilterErrNotAllowedToCreate(qProps)
		}

		// Set new values after beforeCreate events are emitted
		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		if r, err = svc.route.FindByID(ctx, new.Route); err != nil {
			return ApigwFilterErrNotFound(qProps)
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

		// send the signal to reload all routes
		if r.Enabled {
			apigw.Service().Reload(ctx)
		}

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, ApigwFilterActionCreate, err)
}

func (svc *apigwFilter) validateAsyncRoute(ctx context.Context, r *types.ApigwRoute, f *types.ApigwFilter, props *apigwFilterActionProps) (err error) {
	filters, _, err := svc.Search(ctx, types.ApigwFilterFilter{
		RouteID: r.ID,
		Enabled: true,
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

		if !svc.ac.CanUpdateApigwFilter(ctx, qq) {
			return ApigwFilterErrNotAllowedToUpdate(qProps)
		}

		if r, err = svc.route.FindByID(ctx, upd.Route); err != nil {
			return err
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

		// send the signal to reload all routes
		if r.Enabled {
			apigw.Service().Reload(ctx)
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

		if !svc.ac.CanDeleteApigwFilter(ctx, q) {
			return ApigwFilterErrNotAllowedToDelete(qProps)
		}

		if r, err = svc.route.FindByID(ctx, q.Route); err != nil {
			return err
		}

		qProps.setFilter(q)

		q.DeletedAt = now()
		q.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateApigwFilter(ctx, svc.store, q); err != nil {
			return
		}

		// send the signal to reload all routes
		if r.Enabled {
			apigw.Service().Reload(ctx)
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

		if !svc.ac.CanDeleteApigwFilter(ctx, q) {
			return ApigwFilterErrNotAllowedToDelete(qProps)
		}

		if r, err = svc.route.FindByID(ctx, q.Route); err != nil {
			return err
		}

		qProps.setFilter(q)

		q.DeletedAt = nil
		q.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateApigwFilter(ctx, svc.store, q); err != nil {
			return
		}

		// send the signal to reload all routes
		if r.Enabled {
			apigw.Service().Reload(ctx)
		}

		return nil
	}()

	return svc.recordAction(ctx, qProps, ApigwFilterActionDelete, err)
}

func (svc *apigwFilter) Search(ctx context.Context, filter types.ApigwFilterFilter) (r types.ApigwFilterSet, f types.ApigwFilterFilter, err error) {
	var (
		aProps = &apigwFilterActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ApigwFilter) (bool, error) {
		if !svc.ac.CanReadApigwFilter(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
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
		if !svc.ac.CanSearchApigwFilters(ctx) {
			return ApigwFilterErrNotAllowedToRead(qProps)
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
		if !svc.ac.CanSearchApigwFilters(ctx) {
			return ApigwFilterErrNotAllowedToRead(qProps)
		}

		// get the definitions from registry
		l = apigw.Service().ProxyAuthDef()

		return nil
	}()

	return l, svc.recordAction(ctx, qProps, ApigwFilterActionSearch, err)

}
