package service

import (
	"context"

	internalRules "github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/types"
)

type (
	permissions struct {
		db  db
		ctx context.Context

		rules RulesService
	}

	resource interface {
		PermissionResource() internalRules.Resource
	}

	PermissionsService interface {
		With(context.Context) PermissionsService

		Effective() (ee []effectivePermission, err error)

		CanCreateOrganisation() bool
		CanCreateUser() bool
		CanCreateRole() bool
		CanCreateApplication() bool

		CanReadRole(rl *types.Role) bool
		CanUpdateRole(rl *types.Role) bool
		CanDeleteRole(rl *types.Role) bool
		CanManageRoleMembers(rl *types.Role) bool

		CanReadApplication(app *types.Application) bool
		CanUpdateApplication(app *types.Application) bool
		CanDeleteApplication(app *types.Application) bool

		CanUpdateUser(u *types.User) bool
		CanSuspendUser(u *types.User) bool
		CanUnsuspendUser(u *types.User) bool
		CanDeleteUser(u *types.User) bool
	}

	effectivePermission struct {
		Resource  string `json:"resource"`
		Operation string `json:"operation"`
		Allow     bool   `json:"allow"`
	}
)

func Permissions() PermissionsService {
	return (&permissions{
		rules: DefaultRules,
	}).With(context.Background())
}

func (p *permissions) With(ctx context.Context) PermissionsService {
	db := repository.DB(ctx)
	return &permissions{
		db:  db,
		ctx: ctx,

		rules: p.rules.With(ctx),
	}
}

func (p *permissions) Effective() (ee []effectivePermission, err error) {
	ep := func(res, op string, allow bool) effectivePermission {
		return effectivePermission{
			Resource:  res,
			Operation: op,
			Allow:     allow,
		}
	}

	ee = append(ee, ep("system", "access", p.CanAccess()))
	ee = append(ee, ep("system", "application.create", p.CanCreateApplication()))
	ee = append(ee, ep("system", "role.create", p.CanCreateRole()))
	ee = append(ee, ep("system", "organisation.create", p.CanCreateOrganisation()))
	ee = append(ee, ep("system", "grant", p.CanCreateRole()))

	return
}

func (p *permissions) CanAccess() bool {
	return p.checkAccess(types.PermissionResource, "access")
}

func (p *permissions) CanCreateOrganisation() bool {
	return p.checkAccess(types.PermissionResource, "organisation.create")
}

func (p *permissions) CanCreateUser() bool {
	return p.checkAccess(types.PermissionResource, "user.create", p.allow())
}

func (p *permissions) CanCreateRole() bool {
	return p.checkAccess(types.PermissionResource, "role.create")
}

func (p *permissions) CanCreateApplication() bool {
	return p.checkAccess(types.PermissionResource, "application.create")
}

func (p *permissions) CanReadRole(rl *types.Role) bool {
	return p.checkAccess(rl, "read", p.allow())
}

func (p *permissions) CanUpdateRole(rl *types.Role) bool {
	return p.checkAccess(rl, "update")
}

func (p *permissions) CanDeleteRole(rl *types.Role) bool {
	return p.checkAccess(rl, "delete")
}

func (p *permissions) CanManageRoleMembers(rl *types.Role) bool {
	return p.checkAccess(rl, "members.manage")
}

func (p *permissions) CanReadApplication(app *types.Application) bool {
	return p.checkAccess(app, "read", p.allow())
}

func (p *permissions) CanUpdateApplication(app *types.Application) bool {
	return p.checkAccess(app, "update")
}

func (p *permissions) CanDeleteApplication(app *types.Application) bool {
	return p.checkAccess(app, "delete")
}

func (p *permissions) CanUpdateUser(u *types.User) bool {
	return p.checkAccess(u, "update")
}

func (p *permissions) CanSuspendUser(u *types.User) bool {
	return p.checkAccess(u, "suspend")
}

func (p *permissions) CanUnsuspendUser(u *types.User) bool {
	return p.checkAccess(u, "unsuspend")
}

func (p *permissions) CanDeleteUser(u *types.User) bool {
	return p.checkAccess(u, "delete")
}

func (p *permissions) checkAccess(r resource, operation string, fallbacks ...internalRules.CheckAccessFunc) bool {
	return p.rules.Check(r.PermissionResource(), operation, fallbacks...) == internalRules.Allow
}

func (p permissions) allow() func() internalRules.Access {
	return func() internalRules.Access {
		return internalRules.Allow
	}
}
