package service

import (
	"context"

	composeService "github.com/crusttech/human/server/compose/service"
	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/store"
)

type (
	sharedModule struct {
		node      node
		ac        sharedModuleAccessController
		compose   composeService.ModuleService
		store     store.Storer
		actionlog actionlog.Recorder
	}

	sharedModuleAccessController interface {
		CanCreateModuleOnNode(ctx context.Context, r *types.Node) bool
	}

	SharedModuleService interface {
		Create(ctx context.Context, new *types.SharedModule) (*types.SharedModule, error)
		Update(ctx context.Context, updated *types.SharedModule) (*types.SharedModule, error)
		Search(ctx context.Context, filter types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error)
		FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (*types.SharedModule, error)
	}
)

func SharedModule() *sharedModule {
	return &sharedModule{
		ac:        DefaultAccessControl,
		node:      *DefaultNode,
		compose:   composeService.DefaultModule,
		store:     DefaultStore,
		actionlog: DefaultActionlog,
	}
}

// onLookup is the generated FindByID body handler (node-scoped compound id).
//
// The recordAction wrapper and aProps are owned by the generated
// shared_module.gen.go.
func (svc *sharedModule) onLookup(ctx context.Context, nodeID, ID uint64, aProps *sharedModuleActionProps) (module *types.SharedModule, err error) {
	if module, err = loadSharedModuleScoped(ctx, svc.store, nodeID, ID); err != nil {
		return nil, err
	}

	return module, nil
}

// onCreate is the generated Create body handler.
//
// The recordAction wrapper, aProps (module/changed) and res=new assignment are
// owned by the generated shared_module.gen.go.
func (svc *sharedModule) onCreate(ctx context.Context, new *types.SharedModule) error {
	var (
		aProps = &sharedModuleActionProps{changed: new}
	)

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var (
			node *types.Node
		)

		if node, err = svc.node.FindByID(ctx, new.NodeID); err != nil {
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

		// check if Fields can be unmarshaled to the fields structure
		if new.Fields != nil {
		}

		aProps.setModule(new)

		if err = store.CreateFederationSharedModule(ctx, s, new); err != nil {
			return err
		}

		return nil
	})
}

// onUpdate is the generated Update body handler.
//
// The recordAction wrapper, Tx, load, and old clone are owned by shared_module.gen.go.
func (svc *sharedModule) onUpdate(ctx context.Context, s store.Storer, upd, res *types.SharedModule, aProps *sharedModuleActionProps, _ func() error, _ func() error) error {
	if _, err := svc.node.FindByID(ctx, upd.NodeID); err != nil {
		return SharedModuleErrNodeNotFound()
	}

	res.Fields = upd.Fields
	res.Handle = upd.Handle
	res.Name = upd.Name
	res.UpdatedAt = now()
	res.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()

	aProps.setModule(res)

	return store.UpdateFederationSharedModule(ctx, s, res)
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

// onSearch is the generated Search body handler.
//
// The recordAction wrapper and aProps (filter) are owned by the generated
// shared_module.gen.go.
func (svc *sharedModule) onSearch(ctx context.Context, filter types.SharedModuleFilter, aProps *sharedModuleActionProps) (set types.SharedModuleSet, f types.SharedModuleFilter, err error) {
	if set, f, err = store.SearchFederationSharedModules(ctx, svc.store, filter); err != nil {
		return nil, f, err
	}

	return set, f, nil
}

func loadSharedModuleScoped(ctx context.Context, s store.Storer, nodeID, moduleID uint64) (res *types.SharedModule, err error) {
	if res, err = loadSharedModule(ctx, s, moduleID); err == nil && res.NodeID != nodeID {
		return nil, SharedModuleErrNotFound()
	}
	return
}

