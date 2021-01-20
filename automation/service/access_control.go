package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"
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

	ee.Push(types.AutomationRBACResource, "access", svc.CanAccess(ctx))
	ee.Push(types.AutomationRBACResource, "grant", svc.CanGrant(ctx))
	ee.Push(types.AutomationRBACResource, "workflow.create", svc.CanCreateWorkflow(ctx))
	ee.Push(types.AutomationRBACResource, "sessions.search", svc.CanSearchSessions(ctx))
	ee.Push(types.AutomationRBACResource, "triggers.search", svc.CanSearchTriggers(ctx))

	return
}

func (svc accessControl) CanAccess(ctx context.Context) bool {
	return svc.can(ctx, types.AutomationRBACResource, "access")
}

func (svc accessControl) CanGrant(ctx context.Context) bool {
	return svc.can(ctx, types.AutomationRBACResource, "grant")
}

func (svc accessControl) CanCreateWorkflow(ctx context.Context) bool {
	return svc.can(ctx, types.AutomationRBACResource, "workflow.create")
}

func (svc accessControl) CanSearchTriggers(ctx context.Context) bool {
	return svc.can(ctx, types.AutomationRBACResource, "triggers.search")
}

func (svc accessControl) CanSearchSessions(ctx context.Context) bool {
	return svc.can(ctx, types.AutomationRBACResource, "sessions.search")
}

func (svc accessControl) CanReadWorkflow(ctx context.Context, u *types.Workflow) bool {
	return svc.can(ctx, u.RBACResource(), "read")
}

func (svc accessControl) CanUpdateWorkflow(ctx context.Context, u *types.Workflow) bool {
	return svc.can(ctx, u.RBACResource(), "update")
}

func (svc accessControl) CanDeleteWorkflow(ctx context.Context, u *types.Workflow) bool {
	return svc.can(ctx, u.RBACResource(), "delete")
}

func (svc accessControl) CanUndeleteWorkflow(ctx context.Context, u *types.Workflow) bool {
	return svc.can(ctx, u.RBACResource(), "undelete")
}

func (svc accessControl) CanExecuteWorkflow(ctx context.Context, u *types.Workflow) bool {
	return svc.can(ctx, u.RBACResource(), "execute")
}

func (svc accessControl) CanManageWorkflowTriggers(ctx context.Context, u *types.Workflow) bool {
	return svc.can(ctx, u.RBACResource(), "triggers.manage")
}

func (svc accessControl) CanManageWorkflowSessions(ctx context.Context, u *types.Workflow) bool {
	return svc.can(ctx, u.RBACResource(), "sessions.manage")
}

func (svc accessControl) can(ctx context.Context, res rbac.Resource, op rbac.Operation, ff ...rbac.CheckAccessFunc) bool {
	var (
		u     = internalAuth.GetIdentityFromContext(ctx)
		roles = u.Roles()
	)

	if internalAuth.IsSuperUser(u) {
		// Temp solution to allow migration from passing context to ResourceFilter
		// and checking "superuser" privileges there to more sustainable solution
		// (eg: creating super-role with allow-all)
		return true
	}

	return svc.permissions.Can(roles, res.RBACResource(), op, ff...)
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
		types.AutomationRBACResource,
		"access",
		"grant",
		"workflow.create",
		"triggers.search",
		"sessions.search",
	)

	wl.Set(
		types.WorkflowRBACResource,
		"read",
		"update",
		"delete",
		"undelete",
		"execute",
		"triggers.manage",
		"sessions.manage",
	)

	return wl
}
