package service

import (
	"context"

	composeService "github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	sharedModule struct {
		node      node
		compose   composeService.ModuleService
		store     store.Storer
		actionlog actionlog.Recorder
	}

	SharedModuleService interface {
		Create(ctx context.Context, new *types.SharedModule) (*types.SharedModule, error)
		Update(ctx context.Context, updated *types.SharedModule) (*types.SharedModule, error)
		Find(ctx context.Context, filter types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error)
		FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (*types.SharedModule, error)
	}
)

func SharedModule() SharedModuleService {
	return &sharedModule{
		node:      *DefaultNode,
		compose:   composeService.DefaultModule,
		store:     DefaultStore,
		actionlog: DefaultActionlog,
	}
}

func (svc sharedModule) FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (module *types.SharedModule, err error) {
	err = func() error {
		if module, err = store.LookupFederationSharedModuleByID(ctx, svc.store, moduleID); err != nil {
			return err
		}

		return nil
	}()

	return module, err
}

func (svc sharedModule) Create(ctx context.Context, new *types.SharedModule) (*types.SharedModule, error) {
	var (
		aProps = &sharedModuleActionProps{changed: new}
	)

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		// TODO
		// if !svc.ac.CanCreateFederationExposedModule(ctx, ns) {
		// 	return ExposedModuleErrNotAllowedToCreate()
		// }

		if _, err := svc.node.FindByID(ctx, new.NodeID); err != nil {
			return SharedModuleErrNodeNotFound()
		}

		// Check for node - compose.Module combo
		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = auth.GetIdentityFromContext(ctx).Identity()

		// check if Fields can be unmarshaled to the fields structure
		if new.Fields != nil {
		}

		aProps.setModule(new)

		if err = store.CreateFederationSharedModule(ctx, s, new); err != nil {
			return err
		}

		return nil
	})

	return new, svc.recordAction(ctx, aProps, SharedModuleActionCreate, err)
}

func (svc sharedModule) Update(ctx context.Context, updated *types.SharedModule) (*types.SharedModule, error) {
	var (
		aProps = &sharedModuleActionProps{module: updated}
	)

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		updated.UpdatedAt = now()
		updated.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()

		aProps.setModule(updated)

		if _, err := svc.node.FindByID(ctx, updated.NodeID); err != nil {
			return SharedModuleErrNodeNotFound()
		}

		if err = store.UpdateFederationSharedModule(ctx, s, updated); err != nil {
			return err
		}

		return nil
	})

	return updated, svc.recordAction(ctx, aProps, SharedModuleActionUpdate, err)
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
