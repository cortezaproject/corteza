package service

import (
	"context"
	"github.com/crusttech/human/server/pkg/errors"
	"strconv"

	cs "github.com/crusttech/human/server/compose/service"
	ct "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/store"
	ss "github.com/crusttech/human/server/system/service"
	st "github.com/crusttech/human/server/system/types"
)

type (
	exposedModuleServices struct {
		node      node
		module    cs.ModuleService
		namespace cs.NamespaceService
		role      ss.RoleService
	}

	exposedModuleAccessController interface {
		CanCreateModuleOnNode(ctx context.Context, r *types.Node) bool
		CanManageExposedModule(ctx context.Context, r *types.ExposedModule) bool
	}

	ExposedModuleService interface {
		Create(ctx context.Context, new *types.ExposedModule) (*types.ExposedModule, error)
		Update(ctx context.Context, updated *types.ExposedModule) (*types.ExposedModule, error)
		Find(ctx context.Context, filter types.ExposedModuleFilter) (types.ExposedModuleSet, types.ExposedModuleFilter, error)
		Search(ctx context.Context, filter types.ExposedModuleFilter) (types.ExposedModuleSet, types.ExposedModuleFilter, error)
		FindByID(ctx context.Context, nodeID uint64, moduleID uint64) (*types.ExposedModule, error)
		DeleteByID(ctx context.Context, nodeID uint64, ID uint64) error
	}

	moduleUpdateHandler func(ctx context.Context, ns *types.Node, c *types.ExposedModule) (bool, bool, error)
)

func ExposedModule() *exposedModule {
	return &exposedModule{
		ac:        DefaultAccessControl,
		store:     DefaultStore,
		actionlog: DefaultActionlog,
		services: &exposedModuleServices{
			node:      *DefaultNode,
			module:    cs.DefaultModule,
			role:      ss.DefaultRole,
			namespace: cs.DefaultNamespace,
		},
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
		_, err = svc.FindByID(ctx, nodeID, 0)
	}

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (svc *exposedModule) onLookup(ctx context.Context, nodeID uint64, ID uint64, aProps *exposedModuleActionProps) (*types.ExposedModule, error) {
	res, err := loadExposedModuleByNodeID(ctx, svc.store, nodeID, ID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanManageExposedModule(ctx, res) {
		return nil, ExposedModuleErrNotAllowedToManage()
	}

	aProps.setModule(res)
	return res, nil
}

func (svc *exposedModule) onSearch(ctx context.Context, filter types.ExposedModuleFilter, aProps *exposedModuleActionProps) (types.ExposedModuleSet, types.ExposedModuleFilter, error) {
	filter.Check = func(res *types.ExposedModule) (bool, error) {
		if !svc.ac.CanManageExposedModule(ctx, res) {
			return false, ExposedModuleErrNotAllowedToManage()
		}
		return true, nil
	}

	set, f, err := store.SearchFederationExposedModules(ctx, svc.store, filter)
	return set, f, err
}

func (svc *exposedModule) onCreate(ctx context.Context, new *types.ExposedModule) error {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		var (
			m       *ct.Module
			node    *types.Node
			fedRole *st.Role
			err     error
		)

		if node, err = svc.services.node.FindByID(ctx, new.NodeID); err != nil {
			return ExposedModuleErrNodeNotFound()
		}

		if !svc.ac.CanCreateModuleOnNode(ctx, node) {
			return ExposedModuleErrNotAllowedToCreate()
		}

		if _, err = svc.services.namespace.FindByID(ctx, new.ComposeNamespaceID); err != nil {
			return ExposedModuleErrComposeNamespaceNotFound()
		}

		if m, err = svc.services.module.FindByID(ctx, new.ComposeNamespaceID, new.ComposeModuleID); err != nil {
			return ExposedModuleErrComposeModuleNotFound()
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		ctxIdentity := auth.GetIdentityFromContext(ctx)
		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = ctxIdentity.Identity()

		if fedRole, err = svc.services.role.FindByHandle(ctx, "federation"); err != nil {
			return nil
		}

		if fedRole != nil {
			err = cs.DefaultAccessControl.Grant(ctx, rbac.AllowRule(fedRole.ID, ct.RecordRbacResource(m.NamespaceID, m.ID, 0), "read"))
		}

		AddFederationLabel(m, "federation", node.BaseURL)

		if _, err = svc.services.module.Update(ctx, m); err != nil {
			return err
		}

		return store.CreateFederationExposedModule(ctx, s, new)
	})
}

func (svc *exposedModule) onUpdate(ctx context.Context, s store.Storer, upd *types.ExposedModule, res *types.ExposedModule, aProps *exposedModuleActionProps, before func() error, after func() error) error {
	var (
		m    *ct.Module
		node *types.Node
		err  error
	)

	if node, err = svc.services.node.FindByID(ctx, upd.NodeID); err != nil {
		return ExposedModuleErrNodeNotFound()
	}

	if !svc.ac.CanManageExposedModule(ctx, upd) {
		return ExposedModuleErrNotAllowedToManage()
	}

	if _, err = svc.services.namespace.FindByID(ctx, upd.ComposeNamespaceID); err != nil {
		return ExposedModuleErrComposeNamespaceNotFound()
	}

	if m, err = svc.services.module.FindByID(ctx, upd.ComposeNamespaceID, upd.ComposeModuleID); err != nil {
		return ExposedModuleErrComposeModuleNotFound()
	}

	upd.UpdatedAt = now()
	upd.CreatedAt = res.CreatedAt
	upd.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()

	AddFederationLabel(m, "federation", node.BaseURL)

	if _, err = svc.services.module.Update(ctx, m); err != nil {
		return err
	}

	return nil
}

func (svc *exposedModule) onDelete(ctx context.Context, s store.Storer, nodeID uint64, res *types.ExposedModule, aProps *exposedModuleActionProps) error {
	if _, err := svc.services.node.FindByID(ctx, nodeID); err != nil {
		return ExposedModuleErrNodeNotFound()
	}

	if !svc.ac.CanManageExposedModule(ctx, res) {
		return ExposedModuleErrNotAllowedToManage()
	}

	res.DeletedAt = now()
	res.DeletedBy = auth.GetIdentityFromContext(ctx).Identity()

	return store.UpdateFederationExposedModule(ctx, s, res)
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
		if m, err = loadExposedModuleByNodeID(ctx, svc.store, nodeID, moduleID); err != nil {
			return err
		}

		if moduleChanged, fieldsChanged, err = fn(ctx, n, m); err != nil {
			return err
		}

		_ = moduleChanged
		_ = fieldsChanged

		return err
	})

	return m, svc.recordAction(ctx, aProps, action, err)
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

func loadExposedModuleByNodeID(ctx context.Context, s store.FederationExposedModules, nodeID, ID uint64) (res *types.ExposedModule, err error) {
	if ID == 0 || nodeID == 0 {
		return nil, SharedModuleErrInvalidID()
	}

	if res, err = store.LookupFederationExposedModuleByID(ctx, s, ID); errors.IsNotFound(err) {
		err = SharedModuleErrNotFound()
	}

	if err == nil && nodeID != res.NodeID {
		return nil, SharedModuleErrNotFound()
	}

	return
}

func loadExposedModuleScoped(ctx context.Context, s store.FederationExposedModules, nodeID, ID uint64) (res *types.ExposedModule, err error) {
	return loadExposedModuleByNodeID(ctx, s, nodeID, ID)
}
