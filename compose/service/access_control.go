package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	accessControl struct {
		permissions accessControlPermissionServicer
	}

	accessControlPermissionServicer interface {
		Can(context.Context, permissions.Resource, permissions.Operation, ...permissions.CheckAccessFunc) bool
		Grant(context.Context, permissions.Whitelist, ...*permissions.Rule) error
		FindRulesByRoleID(roleID uint64) (rr permissions.RuleSet)
		ResourceFilter(context.Context, permissions.Resource, permissions.Operation, permissions.Access) *permissions.ResourceFilter
	}

	permissionResource interface {
		PermissionResource() permissions.Resource
	}
)

func AccessControl(perm accessControlPermissionServicer) *accessControl {
	return &accessControl{
		permissions: perm,
	}
}

// Effective returns a list of effective service-level permissions
func (svc accessControl) Effective(ctx context.Context) (ee permissions.EffectiveSet) {
	ee = permissions.EffectiveSet{}

	ee.Push(types.ComposePermissionResource, "access", svc.CanAccess(ctx))
	ee.Push(types.ComposePermissionResource, "grant", svc.CanGrant(ctx))
	ee.Push(types.ComposePermissionResource, "namespace.create", svc.CanCreateNamespace(ctx))
	ee.Push(types.ComposePermissionResource, "settings.read", svc.CanReadSettings(ctx))
	ee.Push(types.ComposePermissionResource, "settings.manage", svc.CanManageSettings(ctx))

	return
}

func (svc accessControl) CanAccess(ctx context.Context) bool {
	return svc.can(ctx, types.ComposePermissionResource, "access")
}

func (svc accessControl) CanGrant(ctx context.Context) bool {
	return svc.can(ctx, types.ComposePermissionResource, "grant")
}

func (svc accessControl) CanReadSettings(ctx context.Context) bool {
	return svc.can(ctx, types.ComposePermissionResource, "settings.read")
}

func (svc accessControl) CanManageSettings(ctx context.Context) bool {
	return svc.can(ctx, types.ComposePermissionResource, "settings.manage")
}

func (svc accessControl) CanCreateNamespace(ctx context.Context) bool {
	return svc.can(ctx, types.ComposePermissionResource, "namespace.create")
}

func (svc accessControl) CanReadNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "read", permissions.Allowed)
}

func (svc accessControl) FilterReadableNamespaces(ctx context.Context) *permissions.ResourceFilter {
	return svc.permissions.ResourceFilter(ctx, types.NamespacePermissionResource, "read", permissions.Deny)
}

func (svc accessControl) CanUpdateNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "update")
}

func (svc accessControl) CanDeleteNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "delete")
}

func (svc accessControl) CanManageNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "manage")
}

func (svc accessControl) CanCreateModule(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "module.create")
}

func (svc accessControl) CanCreateAutomationScript(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "automation-script.create")
}

func (svc accessControl) CanReadModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, r, "read")
}

func (svc accessControl) FilterReadableModules(ctx context.Context) *permissions.ResourceFilter {
	return svc.permissions.ResourceFilter(ctx, types.ModulePermissionResource, "read", permissions.Deny)
}

func (svc accessControl) CanUpdateModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, r, "update")
}

func (svc accessControl) CanDeleteModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, r, "delete")
}

func (svc accessControl) CanReadRecordValue(ctx context.Context, r *types.ModuleField) bool {
	return svc.can(ctx, r, "record.value.read", permissions.Allowed)
}

func (svc accessControl) CanUpdateRecordValue(ctx context.Context, r *types.ModuleField) bool {
	return svc.can(ctx, r, "record.value.update", permissions.Allowed)
}

func (svc accessControl) CanCreateRecord(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, r, "record.create")
}

func (svc accessControl) CanReadRecord(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, r, "record.read")
}

func (svc accessControl) CanUpdateRecord(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, r, "record.update")
}

func (svc accessControl) CanDeleteRecord(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, r, "record.delete")
}

func (svc accessControl) CanManageAutomationTriggersOnModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, r, "automation-trigger.manage")
}

func (svc accessControl) CanCreateChart(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "chart.create")
}

func (svc accessControl) CanReadChart(ctx context.Context, r *types.Chart) bool {
	return svc.can(ctx, r, "read")
}

func (svc accessControl) FilterReadableCharts(ctx context.Context) *permissions.ResourceFilter {
	return svc.permissions.ResourceFilter(ctx, types.ChartPermissionResource, "read", permissions.Deny)
}

func (svc accessControl) CanUpdateChart(ctx context.Context, r *types.Chart) bool {
	return svc.can(ctx, r, "update")
}

func (svc accessControl) CanDeleteChart(ctx context.Context, r *types.Chart) bool {
	return svc.can(ctx, r, "delete")
}

func (svc accessControl) CanCreatePage(ctx context.Context, r *types.Namespace) bool {
	// @todo move to func args when namespaces are implemented
	return svc.can(ctx, r, "page.create")
}

func (svc accessControl) CanReadPage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, r, "read")
}

func (svc accessControl) FilterReadablePages(ctx context.Context) *permissions.ResourceFilter {
	return svc.permissions.ResourceFilter(ctx, types.PagePermissionResource, "read", permissions.Deny)
}

func (svc accessControl) CanUpdatePage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, r, "update")
}

func (svc accessControl) CanDeletePage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, r, "delete")
}

func (svc accessControl) CanReadAutomationScript(ctx context.Context, r *automation.Script) bool {
	return svc.can(ctx, types.AutomationScriptPermissionResource.AppendID(r.ID), "read")
}

func (svc accessControl) FilterReadableScripts(ctx context.Context) *permissions.ResourceFilter {
	return svc.permissions.ResourceFilter(ctx, types.AutomationScriptPermissionResource, "read", permissions.Deny)
}

func (svc accessControl) CanUpdateAutomationScript(ctx context.Context, r *automation.Script) bool {
	return svc.can(ctx, types.AutomationScriptPermissionResource.AppendID(r.ID), "update")
}

func (svc accessControl) CanDeleteAutomationScript(ctx context.Context, r *automation.Script) bool {
	return svc.can(ctx, types.AutomationScriptPermissionResource.AppendID(r.ID), "delete")
}

func (svc accessControl) CanRunAutomationTrigger(ctx context.Context, r *automation.Trigger) bool {
	return svc.can(ctx, types.AutomationScriptPermissionResource.AppendID(r.ID), "run", permissions.Allowed)
}

func (svc accessControl) can(ctx context.Context, res permissionResource, op permissions.Operation, ff ...permissions.CheckAccessFunc) bool {
	return svc.permissions.Can(ctx, res.PermissionResource(), op, ff...)
}

func (svc accessControl) Grant(ctx context.Context, rr ...*permissions.Rule) error {
	if !svc.CanGrant(ctx) {
		return ErrNoGrantPermissions
	}

	return svc.permissions.Grant(ctx, svc.Whitelist(), rr...)
}

func (svc accessControl) FindRulesByRoleID(ctx context.Context, roleID uint64) (permissions.RuleSet, error) {
	if !svc.CanGrant(ctx) {
		return nil, ErrNoPermissions
	}

	return svc.permissions.FindRulesByRoleID(roleID), nil
}

func (svc accessControl) Whitelist() permissions.Whitelist {
	var wl = permissions.Whitelist{}

	wl.Set(
		types.ComposePermissionResource,
		"access",
		"grant",
		"namespace.create",
		"settings.read",
		"settings.manage",
	)

	wl.Set(
		types.NamespacePermissionResource,
		"read",
		"update",
		"delete",
		"manage",
		"module.create",
		"chart.create",
		"page.create",
		"automation-script.create",
	)

	wl.Set(
		types.ModulePermissionResource,
		"read",
		"update",
		"delete",
		"record.create",
		"record.read",
		"record.update",
		"record.delete",
		"automation-trigger.manage",
	)

	wl.Set(
		types.ModuleFieldPermissionResource,
		"record.value.read",
		"record.value.update",
	)

	wl.Set(
		types.ChartPermissionResource,
		"read",
		"update",
		"delete",
	)

	wl.Set(
		types.PagePermissionResource,
		"read",
		"update",
		"delete",
	)

	wl.Set(
		types.AutomationScriptPermissionResource,
		"read",
		"update",
		"delete",
	)

	wl.Set(
		types.AutomationTriggerPermissionResource,
		"run",
	)

	return wl
}
