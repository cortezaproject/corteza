package service

import (
	"context"

	cs "github.com/crusttech/human/server/compose/service"
	ct "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/store"
)

type (
	moduleMappingServices struct {
		node      node
		module    cs.ModuleService
		smodule   SharedModuleService
		namespace cs.NamespaceService
	}

	moduleMappingAccessController interface {
		CanMapSharedModule(ctx context.Context, r *types.SharedModule) bool
	}

	ModuleMappingService interface {
		Find(ctx context.Context, filter types.ModuleMappingFilter) (types.ModuleMappingSet, types.ModuleMappingFilter, error)
		Search(ctx context.Context, filter types.ModuleMappingFilter) (types.ModuleMappingSet, types.ModuleMappingFilter, error)
		FindByID(ctx context.Context, federationModuleID uint64) (*types.ModuleMapping, error)
		Create(ctx context.Context, new *types.ModuleMapping) (*types.ModuleMapping, error)
		Update(ctx context.Context, updated *types.ModuleMapping) (*types.ModuleMapping, error)
	}

	moduleMappingUpdateHandler func(ctx context.Context, c *types.ModuleMapping) (bool, bool, error)
)

func ModuleMapping() *moduleMapping {
	return &moduleMapping{
		ac:        DefaultAccessControl,
		store:     DefaultStore,
		actionlog: DefaultActionlog,
		services: &moduleMappingServices{
			node:      *DefaultNode,
			smodule:   DefaultSharedModule,
			module:    cs.DefaultModule,
			namespace: cs.DefaultNamespace,
		},
	}
}

func (svc *moduleMapping) onLookup(ctx context.Context, ID uint64, aProps *moduleMappingActionProps) (*types.ModuleMapping, error) {
	mm, err := store.LookupFederationModuleMappingByFederationModuleID(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	sm, err := svc.services.smodule.FindByID(ctx, mm.NodeID, ID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanMapSharedModule(ctx, sm) {
		return nil, ModuleMappingErrNotAllowedToMap()
	}

	aProps.setMapping(mm)
	return mm, nil
}

func (svc *moduleMapping) onSearch(ctx context.Context, filter types.ModuleMappingFilter, aProps *moduleMappingActionProps) (types.ModuleMappingSet, types.ModuleMappingFilter, error) {
	filter.Check = func(res *types.ModuleMapping) (bool, error) {
		sm, err := svc.services.smodule.FindByID(ctx, res.NodeID, res.FederationModuleID)
		if err != nil {
			return false, err
		}

		if !svc.ac.CanMapSharedModule(ctx, sm) {
			return false, ModuleMappingErrNotAllowedToMap()
		}

		return true, nil
	}

	set, f, err := store.SearchFederationModuleMappings(ctx, svc.store, filter)
	if err != nil {
		return nil, f, err
	}

	set.Walk(func(mm *types.ModuleMapping) error {
		mm.NodeID = f.NodeID
		return nil
	})

	return set, f, nil
}

func (svc *moduleMapping) onCreate(ctx context.Context, new *types.ModuleMapping) error {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		var (
			m   *ct.Module
			sm  *types.SharedModule
			err error
		)

		if _, err = svc.services.namespace.FindByID(ctx, new.ComposeNamespaceID); err != nil {
			return ModuleMappingErrComposeNamespaceNotFound()
		}

		if _, err = svc.services.node.FindByID(ctx, new.NodeID); err != nil {
			return ModuleMappingErrNodeNotFound()
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		if m, err = svc.services.module.FindByID(ctx, new.ComposeNamespaceID, new.ComposeModuleID); err != nil {
			return ModuleMappingErrComposeModuleNotFound()
		}

		if sm, err = svc.services.smodule.FindByID(ctx, new.NodeID, new.FederationModuleID); err != nil {
			return err
		}

		if !svc.ac.CanMapSharedModule(ctx, sm) {
			return ModuleMappingErrNotAllowedToMap()
		}

		if err = store.CreateFederationModuleMapping(ctx, s, new); err != nil {
			return err
		}

		AddFederationLabel(m, "federation", "")

		if _, err = svc.services.module.Update(ctx, m); err != nil {
			return err
		}

		return nil
	})
}

func (svc *moduleMapping) onUpdate(ctx context.Context, upd *types.ModuleMapping, aProps *moduleMappingActionProps) (*types.ModuleMapping, error) {
	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		var (
			m   *ct.Module
			sm  *types.SharedModule
			err error
		)

		if _, err = svc.services.namespace.FindByID(ctx, upd.ComposeNamespaceID); err != nil {
			return ModuleMappingErrComposeNamespaceNotFound()
		}

		if m, err = svc.services.module.FindByID(ctx, upd.ComposeNamespaceID, upd.ComposeModuleID); err != nil {
			return ModuleMappingErrComposeModuleNotFound()
		}

		if sm, err = svc.services.smodule.FindByID(ctx, upd.NodeID, upd.FederationModuleID); err != nil {
			return err
		}

		if !svc.ac.CanMapSharedModule(ctx, sm) {
			return ModuleMappingErrNotAllowedToMap()
		}

		if err = store.UpdateFederationModuleMapping(ctx, s, upd); err != nil {
			return err
		}

		AddFederationLabel(m, "federation", "")

		if _, err = svc.services.module.Update(ctx, m); err != nil {
			return err
		}

		return nil
	})

	return upd, err
}

func (svc moduleMapping) Find(ctx context.Context, filter types.ModuleMappingFilter) (set types.ModuleMappingSet, f types.ModuleMappingFilter, err error) {
	filter.Check = func(res *types.ModuleMapping) (bool, error) {
		sm, err := svc.services.smodule.FindByID(ctx, res.NodeID, res.FederationModuleID)
		if err != nil {
			return false, err
		}

		if !svc.ac.CanMapSharedModule(ctx, sm) {
			return false, ModuleMappingErrNotAllowedToMap()
		}

		return true, nil
	}

	err = func() error {
		if set, f, err = store.SearchFederationModuleMappings(ctx, svc.store, filter); err != nil {
			return err
		}
		return nil
	}()

	set.Walk(func(mm *types.ModuleMapping) error {
		mm.NodeID = f.NodeID
		return nil
	})

	return
}

func (svc moduleMapping) uniqueCheck(ctx context.Context, m *types.ModuleMapping) (err error) {
	f := types.ModuleMappingFilter{
		FederationModuleID: m.FederationModuleID,
		ComposeModuleID:    m.ComposeModuleID,
		ComposeNamespaceID: m.ComposeNamespaceID,
	}

	if set, _, err := store.SearchFederationModuleMappings(ctx, svc.store, f); len(set) > 0 && err == nil {
		return ModuleMappingErrModuleMappingExists()
	} else if err != nil {
		return err
	}

	return err
}
