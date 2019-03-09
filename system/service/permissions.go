package service

import (
	"context"

	internalRules "github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/repository"
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

func (p *permissions) CanCreateOrganisation() bool {
	return p.checkAccess("system", "application.create")
}

func (p *permissions) CanCreateRole() bool {
	return p.checkAccess("system", "role.create")
}

func (p *permissions) CanCreateApplication() bool {
	return p.checkAccess("system", "application.create")
}

func (p *permissions) CanReadRole(rl *types.Role) bool {
	return p.checkAccess(rl.Resource().String(), "read")
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
	return p.checkAccess(app.Resource().String(), "read")
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
