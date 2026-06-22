package service

import (
	"context"
	"strconv"

	cs "github.com/crusttech/human/server/compose/service"
	ct "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/store"
	ss "github.com/crusttech/human/server/system/service"
	st "github.com/crusttech/human/server/system/types"
)

type (
	exposedModule struct {
		node      node
		ac        exposedModuleAccessController
		module    cs.ModuleService
		namespace cs.NamespaceService
		role      ss.RoleService
		store     store.Storer
		actionlog actionlog.Recorder
	}

	exposedModuleAccessController interface {
		CanCreateModuleOnNode(ctx context.Context, r *types.Node) bool
		CanManageExposedModule(ctx context.Context, r *types.ExposedModule) bool
	}

	ExposedModuleService interface {
		Create(ctx context.Context, new *types.ExposedModule) (*types.ExposedModule, error)
		Update(ctx context.Context, updated *types.ExposedModule) (*types.ExposedModule, error)
		Search(ctx context.Context, filter types.ExposedModuleFilter) (types.ExposedModuleSet, types.ExposedModuleFilter, error)
		FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (*types.ExposedModule, error)
		DeleteByID(ctx context.Context, nodeID, moduleID uint64) error
	}

)

func ExposedModule() *exposedModule {
	return &exposedModule{
		ac:        DefaultAccessControl,
		node:      *DefaultNode,
		module:    cs.DefaultModule,
		role:      ss.DefaultRole,
		namespace: cs.DefaultNamespace,
		store:     DefaultStore,
		actionlog: DefaultActionlog,
	}
}

// FindByAny tries to find module in a particular namespace by id, handle or name
func (svc *exposedModule) FindByAny(ctx context.Context, nodeID uint64, identifier interface{}) (m *types.ExposedModule, err error) {
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

// onLookup is the generated FindByID body handler (node-scoped compound id).
//
// The recordAction wrapper and aProps (module) are owned by the generated
// exposed_module.gen.go.
func (svc *exposedModule) onLookup(ctx context.Context, nodeID, moduleID uint64, aProps *exposedModuleActionProps) (module *types.ExposedModule, err error) {
	if module, err = loadExposedModuleScoped(ctx, svc.store, nodeID, moduleID); err != nil {
		return nil, err
	}

	if !svc.ac.CanManageExposedModule(ctx, module) {
		return nil, ExposedModuleErrNotAllowedToManage()
	}

	return module, nil
}

// onUpdate is the generated Update body handler.
//
// The recordAction wrapper, Tx, load, and old clone are owned by exposed_module.gen.go.
func (svc *exposedModule) onUpdate(ctx context.Context, s store.Storer, upd, res *types.ExposedModule, aProps *exposedModuleActionProps, _ func() error, _ func() error) error {
	var (
		m    *ct.Module
		node *types.Node
		err  error
	)

	if node, err = svc.node.FindByID(ctx, upd.NodeID); err != nil {
		return ExposedModuleErrNodeNotFound()
	}

	if !svc.ac.CanManageExposedModule(ctx, upd) {
		return ExposedModuleErrNotAllowedToManage()
	}

	if _, err = svc.namespace.FindByID(ctx, upd.ComposeNamespaceID); err != nil {
		return ExposedModuleErrComposeNamespaceNotFound()
	}

	if m, err = svc.module.FindByID(ctx, upd.ComposeNamespaceID, upd.ComposeModuleID); err != nil {
		return ExposedModuleErrComposeModuleNotFound()
	}

	res.ComposeNamespaceID = upd.ComposeNamespaceID
	res.ComposeModuleID = upd.ComposeModuleID
	res.Fields = upd.Fields
	res.UpdatedAt = now()
	res.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()

	AddFederationLabel(m, "federation", node.BaseURL)

	if _, err = svc.module.Update(ctx, m); err != nil {
		return err
	}

	aProps.setModule(res)

	return store.UpdateFederationExposedModule(ctx, s, res)
}

// onDelete is the generated DeleteByID body handler (node-scoped compound id).
//
// The recordAction wrapper, Tx, and load are owned by exposed_module.gen.go.
func (svc *exposedModule) onDelete(ctx context.Context, s store.Storer, nodeID uint64, res *types.ExposedModule, aProps *exposedModuleActionProps) error {
	if _, err := svc.node.FindByID(ctx, nodeID); err != nil {
		return ExposedModuleErrNodeNotFound()
	}

	if !svc.ac.CanManageExposedModule(ctx, res) {
		return ExposedModuleErrNotAllowedToManage()
	}

	res.DeletedAt = now()
	res.DeletedBy = auth.GetIdentityFromContext(ctx).Identity()

	if err := store.UpdateFederationExposedModule(ctx, s, res); err != nil {
		return err
	}

	aProps.delete = res
	return nil
}

// onSearch is the generated Search body handler.
//
// The recordAction wrapper and aProps (filter) are owned by the generated
// exposed_module.gen.go.
func (svc *exposedModule) onSearch(ctx context.Context, filter types.ExposedModuleFilter, aProps *exposedModuleActionProps) (set types.ExposedModuleSet, f types.ExposedModuleFilter, err error) {
	filter.Check = func(res *types.ExposedModule) (bool, error) {
		if !svc.ac.CanManageExposedModule(ctx, res) {
			return false, ExposedModuleErrNotAllowedToManage()
		}

		return true, nil
	}

	if set, f, err = store.SearchFederationExposedModules(ctx, svc.store, filter); err != nil {
		return
	}

	return set, f, nil
}

// onCreate is the generated Create body handler.
//
// The recordAction wrapper and aProps (module/create) are owned by the
// generated exposed_module.gen.go.
func (svc *exposedModule) onCreate(ctx context.Context, new *types.ExposedModule) error {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var (
			m       *ct.Module
			node    *types.Node
			fedRole *st.Role
		)

		if node, err = svc.node.FindByID(ctx, new.NodeID); err != nil {
			return ExposedModuleErrNodeNotFound()
		}

		if !svc.ac.CanCreateModuleOnNode(ctx, node) {
			return ExposedModuleErrNotAllowedToCreate()
		}

		if _, err := svc.namespace.FindByID(ctx, new.ComposeNamespaceID); err != nil {
			return ExposedModuleErrComposeNamespaceNotFound()
		}

		if m, err = svc.module.FindByID(ctx, new.ComposeNamespaceID, new.ComposeModuleID); err != nil {
			return ExposedModuleErrComposeModuleNotFound()
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		ctxIdentity := auth.GetIdentityFromContext(ctx)

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = ctxIdentity.Identity()

		// find the federation role
		if fedRole, err = svc.role.FindByHandle(ctx, "federation"); err != nil {
			return nil
		}

		if fedRole != nil {
			// get first ID from role and add it as allow rule
			err = cs.DefaultAccessControl.Grant(ctx, rbac.AllowRule(fedRole.ID, ct.RecordRbacResource(m.NamespaceID, m.ID, 0), "read"))
		}

		// set labels
		AddFederationLabel(m, "federation", node.BaseURL)

		if _, err := svc.module.Update(ctx, m); err != nil {
			return err
		}

		if err = store.CreateFederationExposedModule(ctx, s, new); err != nil {
			return err
		}

		return nil
	})
}

func (svc *exposedModule) uniqueCheck(ctx context.Context, m *types.ExposedModule) (err error) {
	f := types.ExposedModuleFilter{
		NodeID:             m.NodeID,
		ComposeModuleID:    m.ComposeModuleID,
		ComposeNamespaceID: m.ComposeNamespaceID,
	}

	set, _, err := store.SearchFederationExposedModules(ctx, svc.store, f)

	if len(set) > 0 && err == nil {
		return ExposedModuleErrNotUnique()
	} else if err != nil {
		return err
	}

	return nil
}

func loadExposedModuleScoped(ctx context.Context, s store.Storer, nodeID, moduleID uint64) (res *types.ExposedModule, err error) {
	if res, err = loadExposedModule(ctx, s, moduleID); err == nil && res.NodeID != nodeID {
		return nil, ExposedModuleErrNotFound()
	}
	return
}

