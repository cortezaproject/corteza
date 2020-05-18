package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	accessControl struct {
		permissions accessControlPermissionServicer
		actionlog   actionlog.Recorder
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
		actionlog:   DefaultActionlog,
	}
}

// Effective returns a list of effective service-level permissions
func (svc accessControl) Effective(ctx context.Context) (ee permissions.EffectiveSet) {
	ee = permissions.EffectiveSet{}

	ee.Push(types.SystemPermissionResource, "access", svc.CanAccess(ctx))
	ee.Push(types.SystemPermissionResource, "grant", svc.CanGrant(ctx))
	ee.Push(types.SystemPermissionResource, "settings.read", svc.CanReadSettings(ctx))
	ee.Push(types.SystemPermissionResource, "settings.manage", svc.CanManageSettings(ctx))
	ee.Push(types.SystemPermissionResource, "application.create", svc.CanCreateApplication(ctx))
	ee.Push(types.SystemPermissionResource, "role.create", svc.CanCreateRole(ctx))
	ee.Push(types.SystemPermissionResource, "organisation.create", svc.CanCreateOrganisation(ctx))

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
	return svc.can(ctx, types.SystemPermissionResource, "user.create")
}

func (svc accessControl) CanCreateRole(ctx context.Context) bool {
	return svc.can(ctx, types.SystemPermissionResource, "role.create")
}

func (svc accessControl) CanCreateApplication(ctx context.Context) bool {
	return svc.can(ctx, types.SystemPermissionResource, "application.create")
}

func (svc accessControl) CanAssignReminder(ctx context.Context) bool {
	return svc.can(ctx, types.SystemPermissionResource, "reminder.assign")
}

func (svc accessControl) CanReadRole(ctx context.Context, rl *types.Role) bool {
	return svc.can(ctx, rl, "read", permissions.Allowed)
}

func (svc accessControl) FilterReadableRoles(ctx context.Context) *permissions.ResourceFilter {
	return svc.permissions.ResourceFilter(ctx, types.RolePermissionResource, "read", permissions.Allow)
}

func (svc accessControl) CanUpdateRole(ctx context.Context, rl *types.Role) bool {
	if rl.ID == permissions.EveryoneRoleID {
		return false
	}

	return svc.can(ctx, rl, "update")
}

func (svc accessControl) CanDeleteRole(ctx context.Context, rl *types.Role) bool {
	if rl.ID == permissions.EveryoneRoleID {
		return false
	}

	return svc.can(ctx, rl, "delete")
}

func (svc accessControl) CanManageRoleMembers(ctx context.Context, rl *types.Role) bool {
	if rl.ID == permissions.EveryoneRoleID {
		return false
	}
	return svc.can(ctx, rl, "members.manage")
}

func (svc accessControl) CanReadApplication(ctx context.Context, app *types.Application) bool {
	return svc.can(ctx, app, "read", permissions.Allowed)
}

func (svc accessControl) FilterReadableApplications(ctx context.Context) *permissions.ResourceFilter {
	return svc.permissions.ResourceFilter(ctx, types.ApplicationPermissionResource, "read", permissions.Deny)
}

func (svc accessControl) CanUpdateApplication(ctx context.Context, app *types.Application) bool {
	return svc.can(ctx, app, "update")
}

func (svc accessControl) CanDeleteApplication(ctx context.Context, app *types.Application) bool {
	return svc.can(ctx, app, "delete")
}

func (svc accessControl) FilterReadableUsers(ctx context.Context) *permissions.ResourceFilter {
	return svc.permissions.ResourceFilter(ctx, types.UserPermissionResource, "read", permissions.Deny)
}

func (svc accessControl) FilterUsersWithUnmaskableEmail(ctx context.Context) *permissions.ResourceFilter {
	return svc.permissions.ResourceFilter(ctx, types.UserPermissionResource, "unmask.email", permissions.Deny)
}

func (svc accessControl) FilterUsersWithUnmaskableName(ctx context.Context) *permissions.ResourceFilter {
	return svc.permissions.ResourceFilter(ctx, types.UserPermissionResource, "unmask.name", permissions.Deny)
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

func (svc accessControl) CanUnmaskEmail(ctx context.Context, u *types.User) bool {
	if internalAuth.GetIdentityFromContext(ctx).Identity() == u.ID {
		// Make an exception when users are reading their own info
		return true
	}

	return svc.can(ctx, u, "unmask.email")
}

func (svc accessControl) CanUnmaskName(ctx context.Context, u *types.User) bool {
	if internalAuth.GetIdentityFromContext(ctx).Identity() == u.ID {
		// Make an exception when users are reading their own info
		return true
	}

	return svc.can(ctx, u, "unmask.name")
}

func (svc accessControl) can(ctx context.Context, res permissionResource, op permissions.Operation, ff ...permissions.CheckAccessFunc) bool {
	return svc.permissions.Can(ctx, res.PermissionResource(), op, ff...)
}

func (svc accessControl) Grant(ctx context.Context, rr ...*permissions.Rule) error {
	if !svc.CanGrant(ctx) {
		return AccessControlErrNotAllowedToSetPermissions()
	}

	if err := svc.permissions.Grant(ctx, svc.Whitelist(), rr...); err != nil {
		return AccessControlErrGeneric().Wrap(err)
	}

	svc.logGrants(ctx, rr)

	return nil
}

func (svc accessControl) logGrants(ctx context.Context, rr []*permissions.Rule) {
	if svc.actionlog == nil {
		return
	}

	for _, r := range rr {
		g := AccessControlActionGrant(&accessControlActionProps{r})
		g.log = r.String()
		g.resource = r.Resource.String()

		svc.actionlog.Record(ctx, g)
	}
}

func (svc accessControl) FindRulesByRoleID(ctx context.Context, roleID uint64) (permissions.RuleSet, error) {
	if !svc.CanGrant(ctx) {
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	return svc.permissions.FindRulesByRoleID(roleID), nil
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
		"reminder.assign",
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
		"unmask.email",
		"unmask.name",
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
