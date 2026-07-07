package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type configuredConnectionServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *configuredConnection) Create(ctx context.Context, new *types.ConfiguredConnection) (res *types.ConfiguredConnection, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &configuredConnectionActionProps{connection: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateConfiguredConnection(ctx) {
			return ConfiguredConnectionErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateConfiguredConnection(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ConfiguredConnectionActionCreate, err)
}

// toLabeledConfiguredConnections converts to []label.LabeledResource
func toLabeledConfiguredConnections(set []*types.ConfiguredConnection) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}

func (svc *configuredConnection) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *configuredConnection) scopeServices(ctx context.Context) *configuredConnectionServices {
	return &configuredConnectionServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}

func (svc *configuredConnection) Enable(ctx context.Context, ID uint64) (res *types.ConfiguredConnection, err error) {
	var (
		aProps = &configuredConnectionActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		res, err = svc.onEnable(ctx, aProps, ID)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ConfiguredConnectionActionEnable, err)
}

func (svc *configuredConnection) Check(ctx context.Context, ID uint64) (res *types.ConfiguredConnectionCheckResult, err error) {
	var (
		aProps = &configuredConnectionActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		res, err = svc.onCheck(ctx, aProps, ID)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ConfiguredConnectionActionCheck, err)
}

func (svc *configuredConnection) RefreshDiscovery(ctx context.Context, ID uint64) (res map[string]any, err error) {
	var (
		aProps = &configuredConnectionActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		res, err = svc.onRefreshDiscovery(ctx, aProps, ID)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ConfiguredConnectionActionRefreshDiscovery, err)
}
