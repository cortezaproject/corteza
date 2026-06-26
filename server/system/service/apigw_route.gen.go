package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

func (svc *apigwRoute) FindByID(ctx context.Context, ID uint64) (res *types.ApigwRoute, err error) {
	var (
		aProps = &apigwRouteActionProps{route: &types.ApigwRoute{ID: ID}}
	)

	err = func() error {
		if res, err = loadApigwRoute(ctx, svc.store, ID); err != nil {
			return ApigwRouteErrInvalidID().Wrap(err)
		}

		aProps.setRoute(res)

		if !svc.ac.CanReadApigwRoute(ctx, res) {
			return ApigwRouteErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ApigwRouteActionLookup, err)
}

func (svc *apigwRoute) Search(ctx context.Context, filter types.ApigwRouteFilter) (set types.ApigwRouteSet, f types.ApigwRouteFilter, err error) {
	var (
		aProps = &apigwRouteActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ApigwRoute) (bool, error) {
		if !svc.ac.CanReadApigwRoute(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchApigwRoutes(ctx) {
			return ApigwRouteErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchApigwRoutes(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ApigwRouteActionSearch, err)
}

func (svc *apigwRoute) Create(ctx context.Context, new *types.ApigwRoute) (res *types.ApigwRoute, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &apigwRouteActionProps{route: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateApigwRoute(ctx) {
			return ApigwRouteErrNotAllowedToCreate()
		}
		before := func() error { return nil }
		after := func() error { return nil }
		res = new
		return svc.onCreate(ctx, new, before, after)
	}()

	return res, svc.recordAction(ctx, aProps, ApigwRouteActionCreate, err)
}

func (svc *apigwRoute) Update(ctx context.Context, upd *types.ApigwRoute) (res *types.ApigwRoute, err error) {
	var (
		aProps = &apigwRouteActionProps{update: upd}
		old    *types.ApigwRoute
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadApigwRoute(ctx, s, upd.ID); err != nil {
			return
		}

		aProps.setRoute(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ApigwRouteErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Endpoint = upd.Endpoint
		res.Method = upd.Method
		res.Enabled = upd.Enabled
		res.Group = upd.Group
		res.CreatedBy = upd.CreatedBy
		res.UpdatedBy = upd.UpdatedBy
		res.DeletedBy = upd.DeletedBy
		res.UpdatedAt = now()

		if err = store.UpdateApigwRoute(ctx, s, res); err != nil {
			return err
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, ApigwRouteActionUpdate, err, old, res)
}

func (svc *apigwRoute) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &apigwRouteActionProps{}
		res    *types.ApigwRoute
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadApigwRoute(ctx, s, ID); err != nil {
			return
		}

		aProps.setRoute(res)

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ApigwRouteActionDelete, err)
}

func (svc *apigwRoute) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &apigwRouteActionProps{}
		res    *types.ApigwRoute
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadApigwRoute(ctx, s, ID); err != nil {
			return
		}

		aProps.setRoute(res)

		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ApigwRouteActionUndelete, err)
}

func loadApigwRoute(ctx context.Context, s store.ApigwRoutes, ID uint64) (res *types.ApigwRoute, err error) {
	if ID == 0 {
		return nil, ApigwRouteErrInvalidID()
	}

	if res, err = store.LookupApigwRouteByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ApigwRouteErrNotFound()
	}

	return
}
