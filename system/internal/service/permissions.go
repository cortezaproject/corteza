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

	PermissionsService interface {
		With(context.Context) PermissionsService

		Effective() (ee []effectivePermission, err error)

		CanCreateOrganisation() bool
		CanCreateRole() bool
		CanCreateApplication() bool

		CanReadRole(rl *types.Role) bool
		CanUpdateRole(rl *types.Role) bool
		CanDeleteRole(rl *types.Role) bool
		CanManageRoleMembers(rl *types.Role) bool

		CanReadApplication(app *types.Application) bool
		CanUpdateApplication(app *types.Application) bool
		CanDeleteApplication(app *types.Application) bool
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
	return p.checkAccess("system", "access")
}

func (p *permissions) CanCreateOrganisation() bool {
	return p.checkAccess("system", "organisation.create")
}

func (p *permissions) CanCreateRole() bool {
	return p.checkAccess("system", "role.create")
}

func (p *permissions) CanCreateApplication() bool {
	return p.checkAccess("system", "application.create")
}

func (p *permissions) CanReadRole(rl *types.Role) bool {
	return p.checkAccess(rl.Resource().String(), "read", p.allow())
}

func (p *permissions) CanUpdateRole(rl *types.Role) bool {
	return p.checkAccess(rl.Resource().String(), "update")
}

func (p *permissions) CanDeleteRole(rl *types.Role) bool {
	return p.checkAccess(rl.Resource().String(), "delete")
}

func (p *permissions) CanManageRoleMembers(rl *types.Role) bool {
	return p.checkAccess(rl.Resource().String(), "members.manage")
}

func (p *permissions) CanReadApplication(app *types.Application) bool {
	return p.checkAccess(app.Resource().String(), "read", p.allow())
}

func (p *permissions) CanUpdateApplication(app *types.Application) bool {
	return p.checkAccess(app.Resource().String(), "update")
}

func (p *permissions) CanDeleteApplication(app *types.Application) bool {
	return p.checkAccess(app.Resource().String(), "delete")
}

func (p *permissions) checkAccess(resource string, operation string, fallbacks ...internalRules.CheckAccessFunc) bool {
	access := p.rules.Check(resource, operation, fallbacks...)
	if access == internalRules.Allow {
		return true
	}
	return false
}

func (p permissions) allow() func() internalRules.Access {
	return func() internalRules.Access {
		return internalRules.Allow
	}
}
