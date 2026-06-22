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

func (svc *dalSensitivityLevel) FindByID(ctx context.Context, ID uint64) (res *types.DalSensitivityLevel, err error) {
	var (
		aProps = &dalSensitivityLevelActionProps{sensitivityLevel: &types.DalSensitivityLevel{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, DalSensitivityLevelActionLookup, err)
}

func (svc *dalSensitivityLevel) Search(ctx context.Context, filter types.DalSensitivityLevelFilter) (set types.DalSensitivityLevelSet, f types.DalSensitivityLevelFilter, err error) {
	var (
		aProps = &dalSensitivityLevelActionProps{search: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, DalSensitivityLevelActionSearch, err)
}

func (svc *dalSensitivityLevel) Create(ctx context.Context, new *types.DalSensitivityLevel) (res *types.DalSensitivityLevel, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &dalSensitivityLevelActionProps{sensitivityLevel: new, new: new}
	)

	err = func() (err error) {
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, DalSensitivityLevelActionCreate, err)
}

func (svc *dalSensitivityLevel) Update(ctx context.Context, upd *types.DalSensitivityLevel) (res *types.DalSensitivityLevel, err error) {
	var (
		aProps = &dalSensitivityLevelActionProps{update: upd}
		old    *types.DalSensitivityLevel
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadDalSensitivityLevel(ctx, s, upd.ID); err != nil {
			return
		}

		aProps.setSensitivityLevel(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return DalSensitivityLevelErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return DalSensitivityLevelErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		return svc.onUpdate(ctx, s, upd, res, aProps, before, after)
	})

	return res, svc.recordAction(ctx, aProps, DalSensitivityLevelActionUpdate, err, old, res)
}

func (svc *dalSensitivityLevel) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &dalSensitivityLevelActionProps{}
		res    *types.DalSensitivityLevel
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadDalSensitivityLevel(ctx, s, ID); err != nil {
			return
		}

		aProps.setSensitivityLevel(res)

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, DalSensitivityLevelActionDelete, err)
}

func (svc *dalSensitivityLevel) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &dalSensitivityLevelActionProps{}
		res    *types.DalSensitivityLevel
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadDalSensitivityLevel(ctx, s, ID); err != nil {
			return
		}

		aProps.setSensitivityLevel(res)

		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, DalSensitivityLevelActionUndelete, err)
}

func loadDalSensitivityLevel(ctx context.Context, s store.DalSensitivityLevels, ID uint64) (res *types.DalSensitivityLevel, err error) {
	if ID == 0 {
		return nil, DalSensitivityLevelErrInvalidID()
	}

	if res, err = store.LookupDalSensitivityLevelByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, DalSensitivityLevelErrNotFound()
	}

	return
}
