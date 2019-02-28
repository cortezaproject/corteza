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

		prm systemService.PermissionsService
	}

	PermissionsService interface {
		With(context.Context) PermissionsService

		CanAccessCompose() bool
	}
)

func Permissions() PermissionsService {
	return (&permissions{
		prm: systemService.DefaultPermissions,
	}).With(context.Background())
}

func (p *permissions) With(ctx context.Context) PermissionsService {
	db := repository.DB(ctx)
	return &permissions{
		db:  db,
		ctx: ctx,

		prm: p.prm.With(ctx),
	}
}

func (p *permissions) CanAccessCompose() bool {
	return p.checkAccess("compose", "access")
}

func (p *permissions) checkAccess(resource string, operation string, fallbacks ...internalRules.CheckAccessFunc) bool {
	access := p.prm.Check(resource, operation, fallbacks...)
	if access == internalRules.Allow {
		return true
	}
	return false
}
