package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	types "github.com/crusttech/human/server/federation/types"
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
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, ExposedModuleActionCreate, err)
}

func (svc *exposedModule) Update(ctx context.Context, upd *types.ExposedModule) (res *types.ExposedModule, err error) {
	var (
		aProps = &exposedModuleActionProps{update: upd}
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ExposedModuleActionUpdate, err)
}

func (svc *exposedModule) DeleteByID(ctx context.Context, nodeID uint64, ID uint64) (err error) {
	var (
		aProps = &exposedModuleActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, nodeID, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, ExposedModuleActionDelete, err)
}
