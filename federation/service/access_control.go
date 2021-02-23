package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"

	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

type (
	accessControl struct {
		permissions accessControlRBACServicer
		actionlog   actionlog.Recorder
	}

	accessControlRBACServicer interface {
		Can([]uint64, rbac.Resource, rbac.Operation, ...rbac.CheckAccessFunc) bool
		Grant(context.Context, rbac.Whitelist, ...*rbac.Rule) error
		FindRulesByRoleID(roleID uint64) (rr rbac.RuleSet)
	}

	RBACResource interface {
		RBACResource() rbac.Resource
	}
)

func AccessControl(perm accessControlRBACServicer) *accessControl {
	return &accessControl{
		permissions: perm,
		actionlog:   DefaultActionlog,
	}
}

// Effective returns a list of effective service-level permissions
func (svc accessControl) Effective(ctx context.Context) (ee rbac.EffectiveSet) {
	ee = rbac.EffectiveSet{}

	ee.Push(types.FederationRBACResource, "grant", svc.CanGrant(ctx))
	ee.Push(types.FederationRBACResource, "pair", svc.CanPair(ctx))
	ee.Push(types.FederationRBACResource, "node.create", svc.CanCreateNode(ctx))
	ee.Push(types.FederationRBACResource, "settings.read", svc.CanReadSettings(ctx))
	ee.Push(types.FederationRBACResource, "settings.manage", svc.CanManageSettings(ctx))

	return
}

func (svc accessControl) CanGrant(ctx context.Context) bool {
	return svc.can(ctx, types.FederationRBACResource, "grant")
}

func (svc accessControl) CanPair(ctx context.Context) bool {
	return svc.can(ctx, types.FederationRBACResource, "pair")
}

func (svc accessControl) CanReadSettings(ctx context.Context) bool {
	return svc.can(ctx, types.FederationRBACResource, "settings.read")
}

func (svc accessControl) CanManageSettings(ctx context.Context) bool {
	return svc.can(ctx, types.FederationRBACResource, "settings.manage")
}

func (svc accessControl) CanCreateNode(ctx context.Context) bool {
	return svc.can(ctx, types.FederationRBACResource, "node.create")
}

func (svc accessControl) CanManageNode(ctx context.Context, r *types.Node) bool {
	return svc.can(ctx, r.RBACResource(), "manage")
}

func (svc accessControl) CanCreateModule(ctx context.Context, r *types.Node) bool {
	return svc.can(ctx, r.RBACResource(), "module.create")
}

func (svc accessControl) CanManageModule(ctx context.Context, r *types.ExposedModule) bool {
	return svc.can(ctx, r.RBACResource(), "manage")
}

func (svc accessControl) CanMapModule(ctx context.Context, r *types.SharedModule) bool {
	return svc.can(ctx, r.RBACResource(), "map")
}

func (svc accessControl) can(ctx context.Context, res rbac.Resource, op rbac.Operation, ff ...rbac.CheckAccessFunc) bool {
	var (
		u = internalAuth.GetIdentityFromContext(ctx)
	)

	if internalAuth.IsSuperUser(u) {
		// Temp solution to allow migration from passing context to ResourceFilter
		// and checking "superuser" privileges there to more sustainable solution
		// (eg: creating super-role with allow-all)
		return true
	}

	return svc.permissions.Can(
		append(u.Roles(), res.DynamicRoles(u.Identity())...),
		res.RBACResource(),
		op,
		ff...,
	)
}

func (svc accessControl) Grant(ctx context.Context, rr ...*rbac.Rule) error {
	if !svc.CanGrant(ctx) {
		return AccessControlErrNotAllowedToSetPermissions()
	}

	if err := svc.permissions.Grant(ctx, svc.Whitelist(), rr...); err != nil {
		return AccessControlErrGeneric().Wrap(err)
	}

	svc.logGrants(ctx, rr)

	return nil
}

func (svc accessControl) logGrants(ctx context.Context, rr []*rbac.Rule) {
	if svc.actionlog == nil {
		return
	}

	for _, r := range rr {
		g := AccessControlActionGrant(&accessControlActionProps{r})
		g.log = r.String()
		g.resource = r.Resource.String()

		svc.actionlog.Record(ctx, g.ToAction())
	}
}

func (svc accessControl) FindRulesByRoleID(ctx context.Context, roleID uint64) (rbac.RuleSet, error) {
	if !svc.CanGrant(ctx) {
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	return svc.permissions.FindRulesByRoleID(roleID), nil
}

func (svc accessControl) Whitelist() rbac.Whitelist {
	var wl = rbac.Whitelist{}

	wl.Set(
		types.FederationRBACResource,
		"grant",
		"pair",
		"node.create",
		"settings.read",
		"settings.manage",
	)

	wl.Set(
		types.NodeRBACResource,
		"manage",
		"module.create",
	)

	wl.Set(
		types.ModuleRBACResource,
		"manage",
		"map",
	)

	return wl
}
