package service

import (
	"context"

	"github.com/crusttech/crust/crm/repository"
	internalRules "github.com/crusttech/crust/internal/rules"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	permissions struct {
		db  db
		ctx context.Context

		rules systemService.RulesService
	}

	PermissionsService interface {
		With(context.Context) PermissionsService

		CanAccessCompose() bool
	}
)

func Permissions() PermissionsService {
	return (&permissions{
		rules: systemService.DefaultRules,
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

func (p *permissions) CanAccessCompose() bool {
	return p.checkAccess("compose", "access")
}

func (p *permissions) CanCreateNamspace() bool {
	return p.checkAccess("compose", "namespace.create")
}

func (p *permissions) checkAccess(resource string, operation string, fallbacks ...internalRules.CheckAccessFunc) bool {
	access := p.rules.Check(resource, operation, fallbacks...)
	if access == internalRules.Allow {
		return true
	}
	return false
}
