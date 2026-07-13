package service

import (
	"context"
	"github.com/crusttech/human/server/pkg/errors"

	composeService "github.com/crusttech/human/server/compose/service"
	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/store"
)

type (
	sharedModuleServices struct {
		node    node
		compose composeService.ModuleService
	}

	sharedModuleAccessController interface {
		CanCreateModuleOnNode(ctx context.Context, r *types.Node) bool
	}

	SharedModuleService interface {
		Create(ctx context.Context, new *types.SharedModule) (*types.SharedModule, error)
		Update(ctx context.Context, updated *types.SharedModule) (*types.SharedModule, error)
		Find(ctx context.Context, filter types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error)
		FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (*types.SharedModule, error)
		Search(ctx context.Context, filter types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error)
	}
)

func SharedModule() *sharedModule {
	return &sharedModule{
		ac:        DefaultAccessControl,
		store:     DefaultStore,
		actionlog: DefaultActionlog,
		services: &sharedModuleServices{
			node:    *DefaultNode,
			compose: composeService.DefaultModule,
		},
	}
}

func (svc *sharedModule) onLookup(ctx context.Context, nodeID uint64, ID uint64, aProps *sharedModuleActionProps) (*types.SharedModule, error) {
	res, err := loadSharedModuleByNodeID(ctx, svc.store, nodeID, ID)
	if err != nil {
		return nil, err
	}
	aProps.setModule(res)
	return res, nil
}

func (svc *sharedModule) onSearch(ctx context.Context, filter types.SharedModuleFilter, aProps *sharedModuleActionProps) (types.SharedModuleSet, types.SharedModuleFilter, error) {
	return store.SearchFederationSharedModules(ctx, svc.store, filter)
}

func (svc *sharedModule) onCreate(ctx context.Context, new *types.SharedModule) error {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		var (
			node *types.Node
			err  error
		)

		if node, err = svc.services.node.FindByID(ctx, new.NodeID); err != nil {
			return SharedModuleErrNodeNotFound()
		}

		if !svc.ac.CanCreateModuleOnNode(ctx, node) {
			return SharedModuleErrNotAllowedToCreate()
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = auth.GetIdentityFromContext(ctx).Identity()

		return store.CreateFederationSharedModule(ctx, s, new)
	})
}

func (svc *sharedModule) onUpdate(ctx context.Context, s store.Storer, upd *types.SharedModule, res *types.SharedModule, aProps *sharedModuleActionProps, before func() error, after func() error) error {
	if _, err := svc.services.node.FindByID(ctx, upd.NodeID); err != nil {
		return SharedModuleErrNodeNotFound()
	}

	upd.UpdatedAt = now()
	upd.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()

	return store.UpdateFederationSharedModule(ctx, s, upd)
}

func (svc sharedModule) uniqueCheck(ctx context.Context, m *types.SharedModule) (err error) {
	f := types.SharedModuleFilter{
		NodeID: m.NodeID,
		Handle: m.Handle,
		Name:   m.Name,
	}

	if set, _, err := store.SearchFederationSharedModules(ctx, svc.store, f); len(set) > 0 && err == nil {
		return SharedModuleErrNotUnique()
	} else if err != nil {
		return err
	}

	return nil
}

func (svc sharedModule) Find(ctx context.Context, filter types.SharedModuleFilter) (set types.SharedModuleSet, f types.SharedModuleFilter, err error) {
	var (
		aProps = &sharedModuleActionProps{filter: &filter}
	)

	err = func() error {
		if set, f, err = store.SearchFederationSharedModules(ctx, svc.store, filter); err != nil {
			return err
		}
		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, SharedModuleActionSearch, err)
}

func loadSharedModuleByNodeID(ctx context.Context, s store.FederationSharedModules, nodeID, ID uint64) (res *types.SharedModule, err error) {
	if ID == 0 || nodeID == 0 {
		return nil, SharedModuleErrInvalidID()
	}

	if res, err = store.LookupFederationSharedModuleByID(ctx, s, ID); errors.IsNotFound(err) {
		err = SharedModuleErrNotFound()
	}

	if err == nil && nodeID != res.NodeID {
		return nil, SharedModuleErrNotFound()
	}

	return
}

func loadSharedModuleScoped(ctx context.Context, s store.FederationSharedModules, nodeID, ID uint64) (res *types.SharedModule, err error) {
	return loadSharedModuleByNodeID(ctx, s, nodeID, ID)
}
