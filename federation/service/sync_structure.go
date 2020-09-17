package service

import (
	"context"
	"strconv"

	composeService "github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/davecgh/go-spew/spew"
)

type (
	module struct {
		ctx       context.Context
		compose   composeService.ModuleService
		store     store.Storer
		actionlog actionlog.Recorder
	}

	ExposedModuleService interface {
		Find(ctx context.Context, filter types.ExposedModuleFilter) (types.ExposedModuleSet, error)
		FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (*types.ExposedModule, error)
		FindByAny(ctx context.Context, nodeID uint64, identifier interface{}) (*types.ExposedModule, error)
		DeleteByID(ctx context.Context, nodeID, moduleID uint64) error
		// Remove(ctx context.Context, filter types.ExposedModuleFilter) (err error)
	}

	moduleUpdateHandler func(ctx context.Context, ns *types.Node, c *types.ExposedModule) (bool, bool, error)
)

func ExposedModule() ExposedModuleService {
	return &module{
		ctx:       context.Background(),
		compose:   composeService.Module(),
		store:     DefaultStore,
		actionlog: DefaultActionlog,
	}
}

// FindByAny tries to find module in a particular namespace by id, handle or name
func (svc module) FindByAny(ctx context.Context, nodeID uint64, identifier interface{}) (m *types.ExposedModule, err error) {
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

func (svc module) FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (module *types.ExposedModule, err error) {
	err = func() error {
		if module, err = store.LookupFederationExposedModuleByID(svc.ctx, svc.store, moduleID); err != nil {
			return err
		}

		return nil
	}()

	return module, err
}

func (svc module) DeleteByID(ctx context.Context, nodeID, moduleID uint64) error {
	return trim1st(svc.updater(ctx, nodeID, moduleID, ModuleActionDelete, svc.handleDelete))
}

func (svc module) updater(ctx context.Context, nodeID, moduleID uint64, action func(...*moduleActionProps) *moduleAction, fn moduleUpdateHandler) (*types.ExposedModule, error) {
	var (
		moduleChanged, fieldsChanged bool

		n *types.Node
		m *types.ExposedModule
		// m, old *types.ExposedModule
		aProps = &moduleActionProps{module: &types.ExposedModule{ID: moduleID, NodeID: nodeID}}
		err    error
	)

	spew.Dump("before handle delete", fn, n, m)

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

func (svc module) handleDelete(ctx context.Context, n *types.Node, m *types.ExposedModule) (bool, bool, error) {
	if err := store.DeleteFederationExposedModuleByID(ctx, svc.store, m.ID); err != nil {
		return false, false, err
	}

	return false, false, nil
}

func (svc module) Find(ctx context.Context, filter types.ExposedModuleFilter) (set types.ExposedModuleSet, err error) {
	// get all modules per-node
	// feed the id's into the compose moduleservice
	// get the data
	// transform (but not here)

	// go to store and fetch the id's, first as module id in filter
	// then without it

	err = func() error {
		if set, _, err = store.SearchFederationExposedModules(svc.ctx, svc.store, filter); err != nil {
			return err
		}

		return nil

		// return loadModuleFields(svc.ctx, svc.store, set...)
	}()

	spew.Dump("ERR", err)

	return set, err
}

// trim1st removes 1st param and returns only error
func trim1st(_ interface{}, err error) error {
	return err
}
