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

func (svc *exposedModule) FindByID(ctx context.Context, nodeID uint64, ID uint64) (res *types.ExposedModule, err error) {
	var (
		aProps = &exposedModuleActionProps{module: &types.ExposedModule{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, nodeID, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ExposedModuleActionLookup, err)
}

func (svc *exposedModule) Search(ctx context.Context, filter types.ExposedModuleFilter) (set types.ExposedModuleSet, f types.ExposedModuleFilter, err error) {
	var (
		aProps = &exposedModuleActionProps{filter: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, ExposedModuleActionSearch, err)
}

func (svc *exposedModule) Create(ctx context.Context, new *types.ExposedModule) (res *types.ExposedModule, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &exposedModuleActionProps{module: new, create: new}
	)

	err = func() (err error) {
		before := func() error { return nil }
		after := func() error { return nil }
		res = new
		return svc.onCreate(ctx, new, before, after)
	}()

	return res, svc.recordAction(ctx, aProps, ExposedModuleActionCreate, err)
}

func (svc *exposedModule) Update(ctx context.Context, upd *types.ExposedModule) (res *types.ExposedModule, err error) {
	var (
		aProps = &exposedModuleActionProps{update: upd}
		old    *types.ExposedModule
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadExposedModule(ctx, s, upd.ID); err != nil {
			return
		}

		aProps.setModule(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return ExposedModuleErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ExposedModuleErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Handle = upd.Handle
		res.Name = upd.Name
		res.NodeID = upd.NodeID
		res.CreatedBy = upd.CreatedBy
		res.UpdatedBy = upd.UpdatedBy
		res.DeletedBy = upd.DeletedBy
		res.UpdatedAt = now()

		if err = store.UpdateFederationExposedModule(ctx, s, res); err != nil {
			return err
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, ExposedModuleActionUpdate, err, old, res)
}

func (svc *exposedModule) DeleteByID(ctx context.Context, nodeID uint64, ID uint64) (err error) {
	var (
		aProps = &exposedModuleActionProps{}
		res    *types.ExposedModule
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadExposedModule(ctx, s, ID); err != nil {
			return
		}

		aProps.setModule(res)

		return svc.onDelete(ctx, s, nodeID, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ExposedModuleActionDelete, err)
}

func loadExposedModule(ctx context.Context, s store.FederationExposedModules, ID uint64) (res *types.ExposedModule, err error) {
	if ID == 0 {
		return nil, ExposedModuleErrInvalidID()
	}

	if res, err = store.LookupFederationExposedModuleByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ExposedModuleErrNotFound()
	}

	return
}
