package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type apigwFilter struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        apigwFilterAccessController
	services  *apigwFilterServices
}

func (svc *apigwFilter) FindByID(ctx context.Context, ID uint64) (res *types.ApigwFilter, err error) {
	var (
		aProps = &apigwFilterActionProps{filter: &types.ApigwFilter{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ApigwFilterActionLookup, err)
}

func (svc *apigwFilter) Search(ctx context.Context, filter types.ApigwFilterFilter) (set types.ApigwFilterSet, f types.ApigwFilterFilter, err error) {
	var (
		aProps = &apigwFilterActionProps{search: &filter}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, ApigwFilterActionSearch, err)
}

func (svc *apigwFilter) Create(ctx context.Context, new *types.ApigwFilter) (res *types.ApigwFilter, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &apigwFilterActionProps{filter: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, ApigwFilterActionCreate, err)
}

func (svc *apigwFilter) Update(ctx context.Context, upd *types.ApigwFilter) (res *types.ApigwFilter, err error) {
	var (
		aProps = &apigwFilterActionProps{filter: upd}
		old    *types.ApigwFilter
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadApigwFilter(ctx, s, upd.ID); err != nil {
			return
		}

		aProps.setFilter(res)
		aProps.setFilter(res)
		old = res.Clone()

		if err = svc.guard(ctx, res); err != nil {
			return
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ApigwFilterErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Route = upd.Route
		res.Weight = upd.Weight
		res.Kind = upd.Kind
		res.Ref = upd.Ref
		res.Enabled = upd.Enabled
		res.CreatedBy = upd.CreatedBy
		res.UpdatedBy = upd.UpdatedBy
		res.DeletedBy = upd.DeletedBy
		res.UpdatedAt = now()

		if err = store.UpdateApigwFilter(ctx, s, res); err != nil {
			return err
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, ApigwFilterActionUpdate, err, old, res)
}

func (svc *apigwFilter) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &apigwFilterActionProps{}
		res    *types.ApigwFilter
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadApigwFilter(ctx, s, ID); err != nil {
			return
		}

		aProps.setFilter(res)

		if err = svc.guard(ctx, res); err != nil {
			return
		}

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ApigwFilterActionDelete, err)
}

func loadApigwFilter(ctx context.Context, s store.ApigwFilters, ID uint64) (res *types.ApigwFilter, err error) {
	if ID == 0 {
		return nil, ApigwFilterErrInvalidID()
	}

	if res, err = store.LookupApigwFilterByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ApigwFilterErrNotFound()
	}

	return
}
func (svc *apigwFilter) guard(_ context.Context, _ *types.ApigwFilter) error { return nil }

func (svc *apigwFilter) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *apigwFilter) DefFilter(ctx context.Context, kind string) (l interface{}, err error) {
	var (
		aProps = &apigwFilterActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		l, err = svc.onDefFilter(ctx, aProps, kind)
		return err
	}()

	return l, svc.recordAction(ctx, aProps, ApigwFilterActionSearch, err)
}

func (svc *apigwFilter) DefProxyAuth(ctx context.Context) (l interface{}, err error) {
	var (
		aProps = &apigwFilterActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		l, err = svc.onDefProxyAuth(ctx, aProps)
		return err
	}()

	return l, svc.recordAction(ctx, aProps, ApigwFilterActionSearch, err)
}
