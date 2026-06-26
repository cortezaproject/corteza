package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

func (svc *dalConnection) FindByID(ctx context.Context, ID uint64) (res *types.DalConnection, err error) {
	var (
		aProps = &dalConnectionActionProps{connection: &types.DalConnection{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, DalConnectionActionLookup, err)
}

func (svc *dalConnection) Search(ctx context.Context, filter types.DalConnectionFilter) (set types.DalConnectionSet, f types.DalConnectionFilter, err error) {
	var (
		aProps = &dalConnectionActionProps{search: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, DalConnectionActionSearch, err)
}

func (svc *dalConnection) Create(ctx context.Context, new *types.DalConnection) (res *types.DalConnection, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &dalConnectionActionProps{connection: new, new: new}
	)

	err = func() (err error) {
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, DalConnectionActionCreate, err)
}

func (svc *dalConnection) Update(ctx context.Context, upd *types.DalConnection) (res *types.DalConnection, err error) {
	var (
		aProps = &dalConnectionActionProps{update: upd}
		old    *types.DalConnection
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadDalConnection(ctx, s, upd.ID); err != nil {
			return
		}

		aProps.setConnection(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return DalConnectionErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return DalConnectionErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Handle = upd.Handle
		res.Type = upd.Type
		res.CreatedBy = upd.CreatedBy
		res.UpdatedBy = upd.UpdatedBy
		res.DeletedBy = upd.DeletedBy
		res.UpdatedAt = now()

		if err = store.UpdateDalConnection(ctx, s, res); err != nil {
			return err
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, DalConnectionActionUpdate, err, old, res)
}

func (svc *dalConnection) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &dalConnectionActionProps{}
		res    *types.DalConnection
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadDalConnection(ctx, s, ID); err != nil {
			return
		}

		aProps.setConnection(res)

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, DalConnectionActionDelete, err)
}

func loadDalConnection(ctx context.Context, s store.DalConnections, ID uint64) (res *types.DalConnection, err error) {
	if ID == 0 {
		return nil, DalConnectionErrInvalidID()
	}

	if res, err = store.LookupDalConnectionByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, DalConnectionErrNotFound()
	}

	return
}
