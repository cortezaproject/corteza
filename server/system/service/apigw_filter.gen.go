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

func (svc *apigwFilter) FindByID(ctx context.Context, ID uint64) (res *types.ApigwFilter, err error) {
	var (
		aProps = &apigwFilterActionProps{filter: &types.ApigwFilter{ID: ID}}
	)

	err = func() error {
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

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ApigwFilterErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		return svc.onUpdate(ctx, s, upd, res, aProps, before, after)
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
