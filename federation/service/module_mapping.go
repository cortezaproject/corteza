package service

import (
	"context"

	cs "github.com/cortezaproject/corteza-server/compose/service"
	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/store"
)

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
		CanMapModule(ctx context.Context, r *types.SharedModule) bool
	}

	ModuleMappingService interface {
		Find(ctx context.Context, filter types.ModuleMappingFilter) (types.ModuleMappingSet, types.ModuleMappingFilter, error)
		FindByID(ctx context.Context, federationModuleID uint64) (*types.ModuleMapping, error)
		Create(ctx context.Context, new *types.ModuleMapping) (*types.ModuleMapping, error)
		Update(ctx context.Context, updated *types.ModuleMapping) (*types.ModuleMapping, error)
	}

	moduleMappingUpdateHandler func(ctx context.Context, c *types.ModuleMapping) (bool, bool, error)
)

func ModuleMapping() ModuleMappingService {
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

func (svc moduleMapping) FindByID(ctx context.Context, federationModuleID uint64) (mm *types.ModuleMapping, err error) {
	err = func() error {
		var (
			sm *types.SharedModule
		)

		if mm, err = store.LookupFederationModuleMappingByFederationModuleID(ctx, svc.store, federationModuleID); err != nil {
			return err
		}

		// fetch shared module for access check
		if sm, err = svc.smodule.FindByID(ctx, mm.NodeID, federationModuleID); err != nil {
			return err
		}

		if !svc.ac.CanMapModule(ctx, sm) {
			return ModuleMappingErrNotAllowedToMap()
		}

		return nil
	}()

	return
}

func (svc moduleMapping) Find(ctx context.Context, filter types.ModuleMappingFilter) (set types.ModuleMappingSet, f types.ModuleMappingFilter, err error) {
	// @todo - optimise this access check
	filter.Check = func(res *types.ModuleMapping) (bool, error) {
		// fetch shared module for this
		sm, err := svc.smodule.FindByID(ctx, res.NodeID, res.FederationModuleID)

		if err != nil {
			return false, err
		}

		if !svc.ac.CanMapModule(ctx, sm) {
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

func (svc moduleMapping) Create(ctx context.Context, new *types.ModuleMapping) (*types.ModuleMapping, error) {
	var (
		aProps = &moduleMappingActionProps{created: new}
	)

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var (
			m  *ct.Module
			sm *types.SharedModule
		)

		if _, err := svc.namespace.With(ctx).FindByID(new.ComposeNamespaceID); err != nil {
			return ModuleMappingErrComposeNamespaceNotFound()
		}

		if _, err = svc.node.FindByID(ctx, new.NodeID); err != nil {
			return ModuleMappingErrNodeNotFound()
		}

		// Check for federation module - compose.Module combo
		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		if m, err = svc.module.With(ctx).FindByID(new.ComposeNamespaceID, new.ComposeModuleID); err != nil {
			return ModuleMappingErrComposeModuleNotFound()
		}

		if sm, err = svc.smodule.FindByID(ctx, new.NodeID, new.FederationModuleID); err != nil {
			return err
		}

		if !svc.ac.CanMapModule(ctx, sm) {
			return ModuleMappingErrNotAllowedToMap()
		}

		if err = store.CreateFederationModuleMapping(ctx, s, new); err != nil {
			return err
		}

		// set labels
		AddFederationLabel(m, "federation", "")

		if _, err := svc.module.With(ctx).Update(m); err != nil {
			return err
		}

		return nil
	})

	return new, svc.recordAction(ctx, aProps, ModuleMappingActionCreate, err)
}

func (svc moduleMapping) Update(ctx context.Context, updated *types.ModuleMapping) (*types.ModuleMapping, error) {
	var (
		aProps = &moduleMappingActionProps{changed: updated}
	)

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var (
			m  *ct.Module
			sm *types.SharedModule
		)

		if _, err := svc.namespace.With(ctx).FindByID(updated.ComposeNamespaceID); err != nil {
			return ModuleMappingErrComposeNamespaceNotFound()
		}

		if m, err = svc.module.With(ctx).FindByID(updated.ComposeNamespaceID, updated.ComposeModuleID); err != nil {
			return ModuleMappingErrComposeModuleNotFound()
		}

		if sm, err = svc.smodule.FindByID(ctx, updated.NodeID, updated.FederationModuleID); err != nil {
			return err
		}

		if !svc.ac.CanMapModule(ctx, sm) {
			return ModuleMappingErrNotAllowedToMap()
		}

		if err = store.UpdateFederationModuleMapping(ctx, s, updated); err != nil {
			return err
		}

		// set labels
		AddFederationLabel(m, "federation", "")

		if _, err := svc.module.With(ctx).Update(m); err != nil {
			return err
		}

		return nil
	})

	return updated, svc.recordAction(ctx, aProps, ModuleMappingActionUpdate, err)
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
