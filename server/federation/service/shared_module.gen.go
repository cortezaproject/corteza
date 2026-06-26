package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	types "github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/store"
)

func (svc *sharedModule) FindByID(ctx context.Context, nodeID uint64, ID uint64) (res *types.SharedModule, err error) {
	var (
		aProps = &sharedModuleActionProps{module: &types.SharedModule{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, nodeID, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, SharedModuleActionLookup, err)
}

func (svc *sharedModule) Search(ctx context.Context, filter types.SharedModuleFilter) (set types.SharedModuleSet, f types.SharedModuleFilter, err error) {
	var (
		aProps = &sharedModuleActionProps{filter: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, SharedModuleActionSearch, err)
}

func (svc *sharedModule) Create(ctx context.Context, new *types.SharedModule) (res *types.SharedModule, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &sharedModuleActionProps{module: new, changed: new}
	)

	err = func() (err error) {
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, SharedModuleActionCreate, err)
}

func (svc *sharedModule) Update(ctx context.Context, upd *types.SharedModule) (res *types.SharedModule, err error) {
	var (
		aProps = &sharedModuleActionProps{module: upd}
		old    *types.SharedModule
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadSharedModule(ctx, s, upd.ID); err != nil {
			return
		}

		aProps.setModule(res)
		aProps.setModule(res)
		old = res.Clone()

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return SharedModuleErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return SharedModuleErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Handle = upd.Handle
		res.Name = upd.Name
		res.CreatedBy = upd.CreatedBy
		res.UpdatedBy = upd.UpdatedBy
		res.DeletedBy = upd.DeletedBy
		res.UpdatedAt = now()

		if err = store.UpdateFederationSharedModule(ctx, s, res); err != nil {
			return err
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, SharedModuleActionUpdate, err, old, res)
}

func loadSharedModule(ctx context.Context, s store.FederationSharedModules, ID uint64) (res *types.SharedModule, err error) {
	if ID == 0 {
		return nil, SharedModuleErrInvalidID()
	}

	if res, err = store.LookupFederationSharedModuleByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, SharedModuleErrNotFound()
	}

	return
}
