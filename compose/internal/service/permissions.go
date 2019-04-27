package service

import (
	"context"

	"github.com/crusttech/crust/compose/internal/repository"
	"github.com/crusttech/crust/compose/types"
	internalRules "github.com/crusttech/crust/internal/rules"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	permissions struct {
		db  db
		ctx context.Context

		rules systemService.RulesService
	}

	permissionResource interface {
		PermissionResource() internalRules.Resource
	}

	PermissionsService interface {
		With(context.Context) PermissionsService

		Effective() (ee []effectivePermission, err error)

		CanAccess() bool
		CanCreateNamspace() bool
		CanCreateModule(r permissionResource) bool
		CanReadModule(r permissionResource) bool
		CanUpdateModule(r permissionResource) bool
		CanDeleteModule(r permissionResource) bool
		CanCreateRecord(r permissionResource) bool
		CanReadRecord(r permissionResource) bool
		CanUpdateRecord(r permissionResource) bool
		CanDeleteRecord(r permissionResource) bool
		CanCreateChart(r permissionResource) bool
		CanReadChart(r permissionResource) bool
		CanUpdateChart(r permissionResource) bool
		CanDeleteChart(r permissionResource) bool
		CanCreateTrigger(r permissionResource) bool
		CanReadTrigger(r permissionResource) bool
		CanUpdateTrigger(r permissionResource) bool
		CanDeleteTrigger(r permissionResource) bool
		CanCreatePage(r permissionResource) bool
		CanReadPage(r permissionResource) bool
		CanUpdatePage(r permissionResource) bool
		CanDeletePage(r permissionResource) bool
	}

	effectivePermission struct {
		Resource  string `json:"resource"`
		Operation string `json:"operation"`
		Allow     bool   `json:"allow"`
	}
)

// Creates a virtual namespace for CRM
//
// We need this until we're through with complete migration
// to Crust Compose
func crmNamespace() *types.Namespace {
	return &types.Namespace{
		ID: types.NamespaceCRM,
	}
}

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

		rules: systemService.Rules(ctx),
	}
}

// Return permissions
func (p *permissions) Effective() (ee []effectivePermission, err error) {
	ep := func(res, op string, allow bool) effectivePermission {
		return effectivePermission{
			Resource:  res,
			Operation: op,
			Allow:     allow,
		}
	}

	ee = append(ee, ep("compose", "access", p.CanAccess()))
	ee = append(ee, ep("compose", "grant", p.CanGrant()))
	ee = append(ee, ep("compose", "namespace.create", p.CanCreateNamspace()))

	ee = append(ee, ep("compose:namespace:crm", "module.create", p.CanCreateModule(crmNamespace())))
	ee = append(ee, ep("compose:namespace:crm", "chart.create", p.CanCreateChart(crmNamespace())))
	ee = append(ee, ep("compose:namespace:crm", "trigger.create", p.CanCreateTrigger(crmNamespace())))
	ee = append(ee, ep("compose:namespace:crm", "page.create", p.CanCreatePage(crmNamespace())))

	return
}

func (p *permissions) CanAccess() bool {
	return p.checkAccess(types.PermissionResource, "access")
}

func (p *permissions) CanGrant() bool {
	return p.checkAccess(types.PermissionResource, "grant")
}

func (p *permissions) CanCreateNamspace() bool {
	return p.checkAccess(types.PermissionResource, "namespace.create")
}

func (p *permissions) CanCreateModule(ns permissionResource) bool {
	// @todo move to func args when namespaces are implemented
	return p.checkAccess(ns, "module.create")
}

func (p *permissions) CanReadModule(r permissionResource) bool {
	return p.checkAccess(r, "read")
}

func (p *permissions) CanUpdateModule(r permissionResource) bool {
	return p.checkAccess(r, "update")
}

func (p *permissions) CanDeleteModule(r permissionResource) bool {
	return p.checkAccess(r, "delete")
}

func (p *permissions) CanCreateRecord(r permissionResource) bool {
	return p.checkAccess(r, "record.create")
}

func (p *permissions) CanReadRecord(r permissionResource) bool {
	return p.checkAccess(r, "record.read")
}

func (p *permissions) CanUpdateRecord(r permissionResource) bool {
	return p.checkAccess(r, "record.update")
}

func (p *permissions) CanDeleteRecord(r permissionResource) bool {
	return p.checkAccess(r, "record.delete")
}

func (p *permissions) CanCreateChart(r permissionResource) bool {
	return p.checkAccess(r, "chart.create")
}

func (p *permissions) CanReadChart(r permissionResource) bool {
	return p.checkAccess(r, "read")
}

func (p *permissions) CanUpdateChart(r permissionResource) bool {
	return p.checkAccess(r, "update")
}

func (p *permissions) CanDeleteChart(r permissionResource) bool {
	return p.checkAccess(r, "delete")
}

func (p *permissions) CanCreateTrigger(r permissionResource) bool {
	return p.checkAccess(r, "trigger.create")
}

func (p *permissions) CanReadTrigger(r permissionResource) bool {
	return p.checkAccess(r, "read")
}

func (p *permissions) CanUpdateTrigger(r permissionResource) bool {
	return p.checkAccess(r, "update")
}

func (p *permissions) CanDeleteTrigger(r permissionResource) bool {
	return p.checkAccess(r, "delete")
}

func (p *permissions) CanCreatePage(r permissionResource) bool {
	// @todo move to func args when namespaces are implemented
	return p.checkAccess(r, "page.create")
}

func (p *permissions) CanReadPage(r permissionResource) bool {
	return p.checkAccess(r, "read")
}

func (p *permissions) CanUpdatePage(r permissionResource) bool {
	return p.checkAccess(r, "update")
}

func (p *permissions) CanDeletePage(r permissionResource) bool {
	return p.checkAccess(r, "delete")
}

func (p *permissions) checkAccess(resource permissionResource, operation string, fallbacks ...internalRules.CheckAccessFunc) bool {
	access := p.rules.Check(resource.PermissionResource(), operation, fallbacks...)
	if access == internalRules.Allow {
		return true
	}
	return false
}
