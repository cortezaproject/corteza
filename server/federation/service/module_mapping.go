package service

import (
	"context"

	cs "github.com/crusttech/human/server/compose/service"
	ct "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/store"
)

// The CRUD skeleton (FindByID, Search, Create, Update + their recordAction
// scaffolds) is generated in module_mapping.gen.go from federation/module_mapping.cue.
//
// This file owns the struct, access-controller interface, constructor, the
// public ModuleMappingService interface, the on<Op> handlers the generated
// methods delegate to (each op is bespoke: federation-specific store funcs,
// shared-module access checks, compose validation and label sync), and the
// uniqueCheck helper.

type (
	moduleMapping struct {
		store     store.Storer
		node      node
		ac        moduleMappingAccessController
		module    cs.ModuleService
		smodule   SharedModuleService
		namespace cs.NamespaceService
		actionlog actionlog.Recorder
	}

	moduleMappingAccessController interface {
		CanMapSharedModule(ctx context.Context, r *types.SharedModule) bool
	}

	ModuleMappingService interface {
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
		node:      *DefaultNode,
		store:     DefaultStore,
		actionlog: DefaultActionlog,
		smodule:   DefaultSharedModule,
		module:    cs.DefaultModule,
		namespace: cs.DefaultNamespace,
	}
}

// onLookup is the custom body for the generated FindByID. The generated method
// owns the action-log scaffold + recordAction; the federation-specific lookup
// (by federation module id) and the shared-module access check live here.
func (svc *moduleMapping) onLookup(ctx context.Context, federationModuleID uint64, aProps *moduleMappingActionProps) (mm *types.ModuleMapping, err error) {
	var (
		sm *types.SharedModule
	)

	if mm, err = store.LookupFederationModuleMappingByFederationModuleID(ctx, svc.store, federationModuleID); err != nil {
		return nil, err
	}

	// fetch shared module for access check
	if sm, err = svc.smodule.FindByID(ctx, mm.NodeID, federationModuleID); err != nil {
		return nil, err
	}

	if !svc.ac.CanMapSharedModule(ctx, sm) {
		return nil, ModuleMappingErrNotAllowedToMap()
	}

	return mm, nil
}

// onSearch is the custom body for the generated Search. The generated method
// owns the action-log scaffold + recordAction; the per-item access check
// (via shared module) and the node-id walk live here.
func (svc *moduleMapping) onSearch(ctx context.Context, filter types.ModuleMappingFilter, aProps *moduleMappingActionProps) (set types.ModuleMappingSet, f types.ModuleMappingFilter, err error) {
	// @todo - optimise this access check
	filter.Check = func(res *types.ModuleMapping) (bool, error) {
		// fetch shared module for this
		sm, err := svc.smodule.FindByID(ctx, res.NodeID, res.FederationModuleID)

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

// onCreate is the custom body for the generated Create. The generated method
// owns the action-log scaffold + recordAction; the compose namespace/module/node
// validation, unique check, shared-module access check, store create and the
// federation-label sync on the compose module live here.
func (svc *moduleMapping) onCreate(ctx context.Context, new *types.ModuleMapping) error {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var (
			m  *ct.Module
			sm *types.SharedModule
		)

		if _, err := svc.namespace.FindByID(ctx, new.ComposeNamespaceID); err != nil {
			return ModuleMappingErrComposeNamespaceNotFound()
		}

		if _, err = svc.node.FindByID(ctx, new.NodeID); err != nil {
			return ModuleMappingErrNodeNotFound()
		}

		// Check for federation module - compose.Module combo
		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		if m, err = svc.module.FindByID(ctx, new.ComposeNamespaceID, new.ComposeModuleID); err != nil {
			return ModuleMappingErrComposeModuleNotFound()
		}

		if sm, err = svc.smodule.FindByID(ctx, new.NodeID, new.FederationModuleID); err != nil {
			return err
		}

		if !svc.ac.CanMapSharedModule(ctx, sm) {
			return ModuleMappingErrNotAllowedToMap()
		}

		if err = store.CreateFederationModuleMapping(ctx, s, new); err != nil {
			return err
		}

		// set labels
		AddFederationLabel(m, "federation", "")

		if _, err := svc.module.Update(ctx, m); err != nil {
			return err
		}

		return nil
	})
}

// onUpdate is the custom body for the generated Update. The generated method
// owns the action-log scaffold + recordAction; the compose namespace/module
// validation, shared-module access check, store update and the federation-label
// sync on the compose module live here.
func (svc *moduleMapping) onUpdate(ctx context.Context, updated *types.ModuleMapping, aProps *moduleMappingActionProps) (*types.ModuleMapping, error) {
	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var (
			m  *ct.Module
			sm *types.SharedModule
		)

		if _, err := svc.namespace.FindByID(ctx, updated.ComposeNamespaceID); err != nil {
			return ModuleMappingErrComposeNamespaceNotFound()
		}

		if m, err = svc.module.FindByID(ctx, updated.ComposeNamespaceID, updated.ComposeModuleID); err != nil {
			return ModuleMappingErrComposeModuleNotFound()
		}

		if sm, err = svc.smodule.FindByID(ctx, updated.NodeID, updated.FederationModuleID); err != nil {
			return err
		}

		if !svc.ac.CanMapSharedModule(ctx, sm) {
			return ModuleMappingErrNotAllowedToMap()
		}

		if err = store.UpdateFederationModuleMapping(ctx, s, updated); err != nil {
			return err
		}

		// set labels
		AddFederationLabel(m, "federation", "")

		if _, err := svc.module.Update(ctx, m); err != nil {
			return err
		}

		return nil
	})

	return updated, err
}

func (svc *moduleMapping) uniqueCheck(ctx context.Context, m *types.ModuleMapping) (err error) {
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
