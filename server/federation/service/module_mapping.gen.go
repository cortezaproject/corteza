package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	types "github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
)

type moduleMapping struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        moduleMappingAccessController
	services  *moduleMappingServices
}

func (svc *moduleMapping) FindByID(ctx context.Context, ID uint64) (res *types.ModuleMapping, err error) {
	var (
		aProps = &moduleMappingActionProps{mapping: &types.ModuleMapping{}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ModuleMappingActionLookup, err)
}

func (svc *moduleMapping) Search(ctx context.Context, filter types.ModuleMappingFilter) (set types.ModuleMappingSet, f types.ModuleMappingFilter, err error) {
	var (
		aProps = &moduleMappingActionProps{filter: &filter}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, ModuleMappingActionSearch, err)
}

func (svc *moduleMapping) Create(ctx context.Context, new *types.ModuleMapping) (res *types.ModuleMapping, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &moduleMappingActionProps{mapping: new, created: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, ModuleMappingActionCreate, err)
}

func (svc *moduleMapping) Update(ctx context.Context, upd *types.ModuleMapping) (res *types.ModuleMapping, err error) {
	var (
		aProps = &moduleMappingActionProps{changed: upd}
		old    *types.ModuleMapping
	)
	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ModuleMappingActionUpdate, err, old, res)
}
func (svc *moduleMapping) guard(_ context.Context, _ *types.ModuleMapping) error { return nil }

func (svc *moduleMapping) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}
