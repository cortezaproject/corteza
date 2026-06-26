package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	types "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
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
		before := func() error { return nil }
		after := func() error { return nil }
		res = new
		return svc.onCreate(ctx, new, before, after)
	}()

	return res, svc.recordAction(ctx, aProps, ModuleActionCreate, err)
}

func (svc *module) Update(ctx context.Context, upd *types.Module) (res *types.Module, err error) {
	var (
		aProps = &moduleActionProps{changed: upd}
		old    *types.Module
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadModule(ctx, s, upd.ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setModule(res)
		aProps.setChanged(res)
		old = res.Clone()

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return ModuleErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ModuleErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.NamespaceID = upd.NamespaceID
		res.Handle = upd.Handle
		res.Name = upd.Name
		res.CreatedByAgent = upd.CreatedByAgent
		res.UpdatedAt = now()

		if err = store.UpdateComposeModule(ctx, s, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, s, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, ModuleActionUpdate, err, old, res)
}

func (svc *module) DeleteByID(ctx context.Context, namespaceID uint64, ID uint64) (err error) {
	var (
		aProps = &moduleActionProps{}
		res    *types.Module
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadModule(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setModule(res)

		return svc.onDelete(ctx, s, namespaceID, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ModuleActionDelete, err)
}

func (svc *module) UndeleteByID(ctx context.Context, namespaceID uint64, ID uint64) (err error) {
	var (
		aProps = &moduleActionProps{}
		res    *types.Module
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadModule(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setModule(res)

		return svc.onUndelete(ctx, s, namespaceID, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ModuleActionUndelete, err)
}

func loadModule(ctx context.Context, s store.ComposeModules, ID uint64) (res *types.Module, err error) {
	if ID == 0 {
		return nil, ModuleErrInvalidID()
	}

	if res, err = store.LookupComposeModuleByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ModuleErrNotFound()
	}

	return
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
