package service

import (
	"context"

	"github.com/crusttech/crust/internal/permissions"
	"github.com/crusttech/crust/system/types"
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

	ee.Push(types.SystemPermissionResource, "access", svc.CanAccess(ctx))
	ee.Push(types.SystemPermissionResource, "settings.read", svc.CanReadSettings(ctx))
	ee.Push(types.SystemPermissionResource, "settings.manage", svc.CanManageSettings(ctx))
	ee.Push(types.SystemPermissionResource, "application.create", svc.CanCreateApplication(ctx))
	ee.Push(types.SystemPermissionResource, "role.create", svc.CanCreateRole(ctx))
	ee.Push(types.SystemPermissionResource, "organisation.create", svc.CanCreateOrganisation(ctx))
	ee.Push(types.SystemPermissionResource, "grant", svc.CanCreateRole(ctx))

	return
}

func (svc accessControl) CanAccess(ctx context.Context) bool {
	return svc.can(ctx, types.SystemPermissionResource, "access")
}

func (svc accessControl) CanReadSettings(ctx context.Context) bool {
	return svc.can(ctx, types.SystemPermissionResource, "settings.read")
}

func (svc accessControl) CanManageSettings(ctx context.Context) bool {
	return svc.can(ctx, types.SystemPermissionResource, "settings.manage")
}

func (svc accessControl) CanCreateOrganisation(ctx context.Context) bool {
	return svc.can(ctx, types.SystemPermissionResource, "organisation.create")
}

func (svc accessControl) CanCreateUser(ctx context.Context) bool {
	return svc.can(ctx, types.SystemPermissionResource, "user.create", permissions.Allowed)
}

func (svc accessControl) CanCreateRole(ctx context.Context) bool {
	return svc.can(ctx, types.SystemPermissionResource, "role.create")
}

func (svc accessControl) CanCreateApplication(ctx context.Context) bool {
	return svc.can(ctx, types.SystemPermissionResource, "application.create")
}

func (svc accessControl) CanReadRole(ctx context.Context, rl *types.Role) bool {
	return svc.can(ctx, rl, "read", permissions.Allowed)
}

func (svc accessControl) CanUpdateRole(ctx context.Context, rl *types.Role) bool {
	return svc.can(ctx, rl, "update")
}

func (svc accessControl) CanDeleteRole(ctx context.Context, rl *types.Role) bool {
	return svc.can(ctx, rl, "delete")
}

func (svc accessControl) CanManageRoleMembers(ctx context.Context, rl *types.Role) bool {
	return svc.can(ctx, rl, "members.manage")
}

func (svc accessControl) CanReadApplication(ctx context.Context, app *types.Application) bool {
	return svc.can(ctx, app, "read", permissions.Allowed)
}

func (svc accessControl) CanUpdateApplication(ctx context.Context, app *types.Application) bool {
	return svc.can(ctx, app, "update")
}

func (svc accessControl) CanDeleteApplication(ctx context.Context, app *types.Application) bool {
	return svc.can(ctx, app, "delete")
}

func (svc accessControl) CanUpdateUser(ctx context.Context, u *types.User) bool {
	return svc.can(ctx, u, "update")
}

func (svc accessControl) CanSuspendUser(ctx context.Context, u *types.User) bool {
	return svc.can(ctx, u, "suspend")
}

func (svc accessControl) CanUnsuspendUser(ctx context.Context, u *types.User) bool {
	return svc.can(ctx, u, "unsuspend")
}

func (svc accessControl) CanDeleteUser(ctx context.Context, u *types.User) bool {
	return svc.can(ctx, u, "delete")
}

func (svc accessControl) can(ctx context.Context, res permissionResource, op permissions.Operation, ff ...permissions.CheckAccessFunc) bool {
	return svc.permissions.Can(ctx, res.PermissionResource(), op, ff...)
}
