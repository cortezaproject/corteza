package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	types "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/label"
)

func (svc *module) FindByID(ctx context.Context, namespaceID uint64, ID uint64) (res *types.Module, err error) {
	var (
		aProps = &moduleActionProps{module: &types.Module{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, namespaceID, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ModuleActionLookup, err)
}

func (svc *module) Search(ctx context.Context, filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error) {
	var (
		aProps = &moduleActionProps{filter: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, ModuleActionSearch, err)
}

func (svc *module) Create(ctx context.Context, new *types.Module) (res *types.Module, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &moduleActionProps{module: new}
	)

	err = func() (err error) {
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, ModuleActionCreate, err)
}

func (svc *module) Update(ctx context.Context, upd *types.Module) (res *types.Module, err error) {
	var (
		aProps = &moduleActionProps{changed: upd}
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ModuleActionUpdate, err)
}

func (svc *module) DeleteByID(ctx context.Context, namespaceID uint64, ID uint64) (err error) {
	var (
		aProps = &moduleActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, namespaceID, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, ModuleActionDelete, err)
}

func (svc *module) UndeleteByID(ctx context.Context, namespaceID uint64, ID uint64) (err error) {
	var (
		aProps = &moduleActionProps{}
	)

	err = func() (err error) {
		return svc.onUndelete(ctx, namespaceID, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, ModuleActionUndelete, err)
}

// toLabeledModules converts to []label.LabeledResource
func toLabeledModules(set []*types.Module) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
