package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	types "github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/scope"
)

type nodeSyncServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *nodeSync) Search(ctx context.Context, filter types.NodeSyncFilter) (set types.NodeSyncSet, f types.NodeSyncFilter, err error) {
	var (
		aProps = &nodeSyncActionProps{nodeSyncFilter: &filter}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, NodeSyncActionSearch, err)
}

func (svc *nodeSync) Create(ctx context.Context, new *types.NodeSync) (res *types.NodeSync, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &nodeSyncActionProps{nodeSync: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, NodeSyncActionCreate, err)
}

func (svc *nodeSync) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *nodeSync) scopeServices(ctx context.Context) *nodeSyncServices {
	return &nodeSyncServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}
