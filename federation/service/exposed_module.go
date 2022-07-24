package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"strconv"

	cs "github.com/cortezaproject/corteza-server/compose/service"
	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	ss "github.com/cortezaproject/corteza-server/system/service"
	st "github.com/cortezaproject/corteza-server/system/types"
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
		Find(ctx context.Context, filter types.ExposedModuleFilter) (types.ExposedModuleSet, types.ExposedModuleFilter, error)
		FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (*types.ExposedModule, error)
		DeleteByID(ctx context.Context, nodeID, moduleID uint64) (*types.ExposedModule, error)
	}

	moduleUpdateHandler func(ctx context.Context, ns *types.Node, c *types.ExposedModule) (bool, bool, error)
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
		if module, err = loadExposedModule(ctx, svc.store, nodeID, moduleID); err != nil {
			return err
		}

		if !svc.ac.CanManageExposedModule(ctx, module) {
			return ExposedModuleErrNotAllowedToManage()
		}

		return nil
	}()

	return module, err
}

func (svc exposedModule) Update(ctx context.Context, updated *types.ExposedModule) (*types.ExposedModule, error) {
	var (
		aProps = &exposedModuleActionProps{update: updated}
	)

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var (
			m    *ct.Module
			node *types.Node
			old  *types.ExposedModule
		)

		if node, err = svc.node.FindByID(ctx, updated.NodeID); err != nil {
			return ExposedModuleErrNodeNotFound()
		}

		if !svc.ac.CanManageExposedModule(ctx, updated) {
			return ExposedModuleErrNotAllowedToManage()
		}

		if _, err := svc.namespace.FindByID(ctx, updated.ComposeNamespaceID); err != nil {
			return ExposedModuleErrComposeNamespaceNotFound()
		}

		if m, err = svc.module.FindByID(ctx, updated.ComposeNamespaceID, updated.ComposeModuleID); err != nil {
			return ExposedModuleErrComposeModuleNotFound()
		}

		if old, err = svc.FindByID(ctx, updated.NodeID, updated.ID); err != nil {
			return ExposedModuleErrNotFound()
		}

		updated.UpdatedAt = now()
		updated.CreatedAt = old.CreatedAt
		updated.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()

		// set labels
		AddFederationLabel(m, "federation", node.BaseURL)

		if _, err := svc.module.Update(ctx, m); err != nil {
			return err
		}

		aProps.setModule(updated)

		if err = store.UpdateFederationExposedModule(ctx, s, updated); err != nil {
			return err
		}

		return nil
	})

	return updated, svc.recordAction(ctx, aProps, ExposedModuleActionUpdate, err)
}

func (svc exposedModule) updater(ctx context.Context, nodeID, moduleID uint64, action func(...*exposedModuleActionProps) *exposedModuleAction, fn moduleUpdateHandler) (*types.ExposedModule, error) {
	var (
		moduleChanged, fieldsChanged bool

		n      *types.Node
		m      *types.ExposedModule
		aProps = &exposedModuleActionProps{module: &types.ExposedModule{ID: moduleID, NodeID: nodeID}}
		err    error
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if m, err = loadExposedModule(ctx, svc.store, nodeID, moduleID); err != nil {
			return err
		}

		// TODO - handle node id also
		if moduleChanged, fieldsChanged, err = fn(ctx, n, m); err != nil {
			return err
		}

		_ = moduleChanged
		_ = fieldsChanged

		return err
	})

	return m, svc.recordAction(ctx, aProps, action, err)
}

func (svc exposedModule) DeleteByID(ctx context.Context, nodeID, moduleID uint64) (m *types.ExposedModule, err error) {
	var (
		aProps = &exposedModuleActionProps{}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var (
			m *types.ExposedModule
		)

		if _, err = svc.node.FindByID(ctx, nodeID); err != nil {
			return ExposedModuleErrNodeNotFound()
		}

		if m, err = svc.FindByID(ctx, nodeID, moduleID); err != nil {
			return err
		}

		if !svc.ac.CanManageExposedModule(ctx, m) {
			return ExposedModuleErrNotAllowedToManage()
		}

		m.DeletedAt = now()
		m.DeletedBy = auth.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateFederationExposedModule(ctx, s, m); err != nil {
			return err
		}

		aProps.delete = m

		return nil
	})

	return m, svc.recordAction(ctx, aProps, ExposedModuleActionDelete, err)
}

func (svc exposedModule) Find(ctx context.Context, filter types.ExposedModuleFilter) (set types.ExposedModuleSet, f types.ExposedModuleFilter, err error) {
	filter.Check = func(res *types.ExposedModule) (bool, error) {
		if !svc.ac.CanManageExposedModule(ctx, res) {
			return false, ExposedModuleErrNotAllowedToManage()
		}

		return true, nil
	}

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
		aProps = &exposedModuleActionProps{create: new}
	)

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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

func loadExposedModule(ctx context.Context, s store.FederationExposedModules, nodeID, ID uint64) (res *types.ExposedModule, err error) {
	if ID == 0 || nodeID == 0 {
		return nil, SharedModuleErrInvalidID()
	}

	if res, err = store.LookupFederationExposedModuleByID(ctx, s, ID); errors.IsNotFound(err) {
		err = SharedModuleErrNotFound()
	}

	if err == nil && nodeID != res.NodeID {
		// Make sure chart belongs to the right namespace
		return nil, SharedModuleErrNotFound()
	}

	return
}
