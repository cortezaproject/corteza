package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	accessControl struct {
		permissions accessControlPermissionServicer
	}

	accessControlPermissionServicer interface {
		Can(context.Context, permissions.Resource, permissions.Operation, ...permissions.CheckAccessFunc) bool
		Grant(context.Context, permissions.Whitelist, ...*permissions.Rule) error
		FindRulesByRoleID(roleID uint64) (rr permissions.RuleSet)
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

func (svc accessControl) CanGrant(ctx context.Context) bool {
	return svc.can(ctx, types.SystemPermissionResource, "grant")
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

func (svc accessControl) Grant(ctx context.Context, rr ...*permissions.Rule) error {
	if !svc.CanGrant(ctx) {
		return ErrNoPermissions
	}

	return svc.permissions.Grant(ctx, svc.Whitelist(), rr...)
}

func (svc accessControl) FindRulesByRoleID(ctx context.Context, roleID uint64) (permissions.RuleSet, error) {
	if !svc.CanGrant(ctx) {
		return nil, ErrNoPermissions
	}

	return svc.permissions.FindRulesByRoleID(roleID), nil
}

// DefaultRules returns list of default rules for this compose service
func (svc accessControl) DefaultRules() permissions.RuleSet {
	var (
		sys           = types.SystemPermissionResource
		applications  = types.ApplicationPermissionResource.AppendWildcard()
		organisations = types.OrganisationPermissionResource.AppendWildcard()
		roles         = types.RolePermissionResource.AppendWildcard()
		users         = types.UserPermissionResource.AppendWildcard()

		allowAdm = func(res permissions.Resource, op permissions.Operation) *permissions.Rule {
			return permissions.AllowRule(permissions.AdminRoleID, res, op)
		}
	)

	return permissions.RuleSet{
		permissions.AllowRule(permissions.EveryoneRoleID, sys, "user.create"),

		allowAdm(sys, "access"),
		allowAdm(sys, "grant"),
		allowAdm(sys, "settings.read"),
		allowAdm(sys, "settings.manage"),
		allowAdm(sys, "organisation.create"),
		allowAdm(sys, "application.create"),
		allowAdm(sys, "user.create"),
		allowAdm(sys, "role.create"),

		allowAdm(organisations, "access"),
		allowAdm(applications, "read"),
		allowAdm(applications, "update"),
		allowAdm(applications, "delete"),

		allowAdm(users, "read"),
		allowAdm(users, "update"),
		allowAdm(users, "suspend"),
		allowAdm(users, "unsuspend"),
		allowAdm(users, "delete"),

		allowAdm(roles, "read"),
		allowAdm(roles, "update"),
		allowAdm(roles, "delete"),
		allowAdm(roles, "members.manage"),
	}
}

func (svc accessControl) Whitelist() permissions.Whitelist {
	var wl = permissions.Whitelist{}

	wl.Set(
		types.SystemPermissionResource,
		"access",
		"grant",
		"settings.read",
		"settings.manage",
		"organisation.create",
		"role.create",
		"user.create",
		"application.create",
	)

	wl.Set(
		types.OrganisationPermissionResource,
		"access",
	)

	wl.Set(
		types.ApplicationPermissionResource,
		"read",
		"update",
		"delete",
	)

	wl.Set(
		types.UserPermissionResource,
		"read",
		"update",
		"delete",
		"suspend",
		"unsuspend",
	)

	wl.Set(
		types.RolePermissionResource,
		"read",
		"update",
		"delete",
		"members.manage",
	)

	return wl
}
