package service

import (
	"context"
	"strconv"

	composeService "github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	exposedModule struct {
		ctx       context.Context
		compose   composeService.ModuleService
		store     store.Storer
		actionlog actionlog.Recorder
	}

	ExposedModuleService interface {
		Find(ctx context.Context, filter types.ExposedModuleFilter) (types.ExposedModuleSet, types.ExposedModuleFilter, error)
		FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (*types.ExposedModule, error)
		FindByAny(ctx context.Context, nodeID uint64, identifier interface{}) (*types.ExposedModule, error)
		DeleteByID(ctx context.Context, nodeID, moduleID uint64) error
		Create(ctx context.Context, new *types.ExposedModule) (*types.ExposedModule, error)
	}

	moduleUpdateHandler func(ctx context.Context, ns *types.Node, c *types.ExposedModule) (bool, bool, error)
)

func ExposedModule() ExposedModuleService {
	return &exposedModule{
		ctx:       context.Background(),
		compose:   composeService.Module(),
		store:     DefaultStore,
		actionlog: DefaultActionlog,
	}
}

// FindByAny tries to find module in a particular namespace by id, handle or name
func (svc exposedModule) FindByAny(ctx context.Context, nodeID uint64, identifier interface{}) (m *types.ExposedModule, err error) {
	if ID, ok := identifier.(uint64); ok {
		m, err = svc.FindByID(ctx, nodeID, ID)
	} else if strIdentifier, ok := identifier.(string); ok {
		if ID, _ := strconv.ParseUint(strIdentifier, 10, 64); ID > 0 {
			m, err = svc.FindByID(ctx, nodeID, ID)
		}
	} else {
		// force invalid ID error
		// we do that to wrap error with lookup action context
		_, err = svc.FindByID(ctx, nodeID, 0)
	}

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (svc exposedModule) FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (module *types.ExposedModule, err error) {
	err = func() error {
		if module, err = store.LookupFederationExposedModuleByID(ctx, svc.store, moduleID); err != nil {
			return err
		}

		return nil
	}()

	return module, err
}

func (svc exposedModule) DeleteByID(ctx context.Context, nodeID, moduleID uint64) error {
	return trim1st(svc.updater(ctx, nodeID, moduleID, ExposedModuleActionDelete, svc.handleDelete))
}

func (svc exposedModule) updater(ctx context.Context, nodeID, moduleID uint64, action func(...*exposedModuleActionProps) *exposedModuleAction, fn moduleUpdateHandler) (*types.ExposedModule, error) {
	var (
		moduleChanged, fieldsChanged bool

		n *types.Node
		m *types.ExposedModule
		// m, old *types.ExposedModule
		aProps = &exposedModuleActionProps{module: &types.ExposedModule{ID: moduleID, NodeID: nodeID}}
		err    error
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if m, err = svc.store.LookupFederationExposedModuleByID(ctx, moduleID); err != nil {
			return err
		}

		// TODO - handle node id also
		if moduleChanged, fieldsChanged, err = fn(ctx, n, m); err != nil {
			return err
		}

		return err
	})

	return m, svc.recordAction(ctx, aProps, action, err)
}

func (svc exposedModule) handleDelete(ctx context.Context, n *types.Node, m *types.ExposedModule) (bool, bool, error) {
	if err := store.DeleteFederationExposedModuleByID(ctx, svc.store, m.ID); err != nil {
		return false, false, err
	}

	return false, false, nil
}

func (svc exposedModule) Find(ctx context.Context, filter types.ExposedModuleFilter) (set types.ExposedModuleSet, f types.ExposedModuleFilter, err error) {
	err = func() error {
		if set, f, err = store.SearchFederationExposedModules(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, err
}

func (svc exposedModule) Create(ctx context.Context, new *types.ExposedModule) (*types.ExposedModule, error) {
	var (
		aProps = &exposedModuleActionProps{changed: new}
	)

	// check if compose module actually exists
	// TODO - how do we handle namespace?
	// if _, err := svc.compose.With(ctx).FindByID(r.NamespaceID, new.ComposeModuleID); err == nil {
	// 	return nil, ExposedModuleErrComposeModuleNotFound()
	// }

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		// TODO
		// if !svc.ac.CanCreateFederationExposedModule(ctx, ns) {
		// 	return ExposedModuleErrNotAllowedToCreate()
		// }

		// TODO - fetch Node
		aProps.setNode(nil)

		// Check for node - compose.Module combo
		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.UpdatedAt = nil
		new.DeletedAt = nil

		// check if Fields can be unmarshaled to the fields structure
		if new.Fields != nil {
		}

		aProps.setModule(new)

		if err = store.CreateFederationExposedModule(ctx, s, new); err != nil {
			return err
		}

		return nil
	})

	return new, svc.recordAction(ctx, aProps, ExposedModuleActionCreate, err)
}

func (svc exposedModule) uniqueCheck(ctx context.Context, m *types.ExposedModule) (err error) {
	f := types.ExposedModuleFilter{
		NodeID:          m.NodeID,
		ComposeModuleID: m.ComposeModuleID,
	}

	if set, _, err := store.SearchFederationExposedModules(ctx, svc.store, f); len(set) > 0 && err != nil {
		return ExposedModuleErrNotUnique()
	}

	return nil
}

// trim1st removes 1st param and returns only error
func trim1st(_ interface{}, err error) error {
	return err
}
