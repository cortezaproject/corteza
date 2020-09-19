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
	sharedModule struct {
		ctx       context.Context
		compose   composeService.ModuleService
		store     store.Storer
		actionlog actionlog.Recorder
	}

	SharedModuleService interface {
		Find(ctx context.Context, filter types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error)
		FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (*types.SharedModule, error)
		FindByAny(ctx context.Context, nodeID uint64, identifier interface{}) (*types.SharedModule, error)
		// DeleteByID(ctx context.Context, nodeID, moduleID uint64) error
	}

	// moduleUpdateHandler func(ctx context.Context, ns *types.Node, c *types.SharedModule) (bool, bool, error)
)

func SharedModule() SharedModuleService {
	return &sharedModule{
		ctx:       context.Background(),
		compose:   composeService.Module(),
		store:     DefaultStore,
		actionlog: DefaultActionlog,
	}
}

// FindByAny tries to find module in a particular namespace by id, handle or name
func (svc sharedModule) FindByAny(ctx context.Context, nodeID uint64, identifier interface{}) (m *types.SharedModule, err error) {
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

func (svc sharedModule) FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (module *types.SharedModule, err error) {
	err = func() error {
		if module, err = store.LookupFederationSharedModuleByID(ctx, svc.store, moduleID); err != nil {
			return err
		}

		return nil
	}()

	return module, err
}

// func (svc sharedModule) DeleteByID(ctx context.Context, nodeID, moduleID uint64) error {
// 	return trim1st(svc.updater(ctx, nodeID, moduleID, ModuleActionDelete, svc.handleDelete))
// }

// func (svc sharedModule) updater(ctx context.Context, nodeID, moduleID uint64, action func(...*moduleActionProps) *moduleAction, fn moduleUpdateHandler) (*types.SharedModule, error) {
// 	var (
// 		moduleChanged, fieldsChanged bool

// 		n *types.Node
// 		m *types.SharedModule
// 		// m, old *types.SharedModule
// 		aProps = &moduleActionProps{module: &types.SharedModule{ID: moduleID, NodeID: nodeID}}
// 		err    error
// 	)

// 	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
// 		if m, err = svc.store.LookupFederationSharedModuleByID(ctx, moduleID); err != nil {
// 			return err
// 		}

// 		// TODO - handle node id also
// 		if moduleChanged, fieldsChanged, err = fn(ctx, n, m); err != nil {
// 			return err
// 		}

// 		return err
// 	})

// 	return m, svc.recordAction(ctx, aProps, action, err)
// }

// func (svc sharedModule) handleDelete(ctx context.Context, n *types.Node, m *types.SharedModule) (bool, bool, error) {
// 	if err := store.DeleteFederationSharedModuleByID(ctx, svc.store, m.ID); err != nil {
// 		return false, false, err
// 	}

// 	return false, false, nil
// }

func (svc sharedModule) Find(ctx context.Context, filter types.SharedModuleFilter) (set types.SharedModuleSet, f types.SharedModuleFilter, err error) {
	var (
		aProps = &sharedModuleActionProps{filter: &filter}
	)

	err = func() error {
		// handle node for actionlog here?
		if f.NodeID > 0 {
		}

		if set, f, err = store.SearchFederationSharedModules(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, SharedModuleActionSearch, err)
}

// // trim1st removes 1st param and returns only error
// func trim1st(_ interface{}, err error) error {
// 	return err
// }
