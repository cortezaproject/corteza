package service

import (
	"context"

	"github.com/crusttech/crust/compose/types"
	"github.com/crusttech/crust/internal/permissions"
)

type (
	accessControl struct {
		permissions permissions.Verifier
	}

	permissionResource interface {
		PermissionResource() permissions.Resource
	}
)

func AccessControl(pv permissions.Verifier) *accessControl {
	return &accessControl{
		permissions: pv,
	}
}

// Effective returns a list of effective service-level permissions
func (svc accessControl) Effective(ctx context.Context) (ee permissions.EffectiveSet) {
	ee = permissions.EffectiveSet{}

	ee.Push(types.ComposePermissionResource, "access", svc.CanAccess(ctx))
	ee.Push(types.ComposePermissionResource, "grant", svc.CanGrant(ctx))
	ee.Push(types.ComposePermissionResource, "namespace.create", svc.CanCreateNamespace(ctx))

	return
}

func (svc accessControl) CanAccess(ctx context.Context) bool {
	return svc.can(ctx, types.ComposePermissionResource, "access")
}

func (svc accessControl) CanGrant(ctx context.Context) bool {
	return svc.can(ctx, types.ComposePermissionResource, "grant")
}

func (svc accessControl) CanCreateNamespace(ctx context.Context) bool {
	return svc.can(ctx, types.ComposePermissionResource, "namespace.create")
}

func (svc accessControl) CanReadNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "read", permissions.Allowed)
}

func (svc accessControl) CanUpdateNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "update")
}

func (svc accessControl) CanDeleteNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "delete")
}

func (svc accessControl) CanCreateModule(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "module.create")
}

func (svc accessControl) CanReadModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, r, "read")
}

func (svc accessControl) CanUpdateModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, r, "update")
}

func (svc accessControl) CanDeleteModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, r, "delete")
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

func (svc accessControl) CanCreateChart(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "chart.create")
}

func (svc accessControl) CanReadChart(ctx context.Context, r *types.Chart) bool {
	return svc.can(ctx, r, "read")
}

func (svc accessControl) CanUpdateChart(ctx context.Context, r *types.Chart) bool {
	return svc.can(ctx, r, "update")
}

func (svc accessControl) CanDeleteChart(ctx context.Context, r *types.Chart) bool {
	return svc.can(ctx, r, "delete")
}

func (svc accessControl) CanCreateTrigger(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, r, "trigger.create")
}

func (svc accessControl) CanReadTrigger(ctx context.Context, r *types.Trigger) bool {
	return svc.can(ctx, r, "read")
}

func (svc accessControl) CanUpdateTrigger(ctx context.Context, r *types.Trigger) bool {
	return svc.can(ctx, r, "update")
}

func (svc accessControl) CanDeleteTrigger(ctx context.Context, r *types.Trigger) bool {
	return svc.can(ctx, r, "delete")
}

func (svc accessControl) CanCreatePage(ctx context.Context, r *types.Namespace) bool {
	// @todo move to func args when namespaces are implemented
	return svc.can(ctx, r, "page.create")
}

func (svc accessControl) CanReadPage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, r, "read")
}

func (svc accessControl) CanUpdatePage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, r, "update")
}

func (svc accessControl) CanDeletePage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, r, "delete")
}

func (svc accessControl) can(ctx context.Context, res permissionResource, op permissions.Operation, ff ...permissions.CheckAccessFunc) bool {
	return svc.permissions.Can(ctx, res.PermissionResource(), op, ff...)
}
