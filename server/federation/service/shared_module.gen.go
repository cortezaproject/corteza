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
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, SharedModuleActionUpdate, err)
}
