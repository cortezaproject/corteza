package service

import (
	"context"

	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
	internalRules "github.com/crusttech/crust/internal/rules"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	permissions struct {
		db  db
		ctx context.Context

		rules systemService.RulesService
	}

	resource interface {
		Resource() internalRules.Resource
	}

	PermissionsService interface {
		With(context.Context) PermissionsService

		Effective() (ee []effectivePermission, err error)

		CanAccess() bool
		CanCreateNamspace() bool
		CanCreateModule() bool
		CanReadModule(r resource) bool
		CanUpdateModule(r resource) bool
		CanDeleteModule(r resource) bool
		CanDeleteModuleByID(ID uint64) bool
		CanCreateRecord(r resource) bool
		CanReadRecord(r resource) bool
		CanUpdateRecord(r resource) bool
		CanDeleteRecord(r resource) bool
		CanDeleteRecordByModuleID(ID uint64) bool
		CanCreateChart() bool
		CanReadChart(r resource) bool
		CanUpdateChart(r resource) bool
		CanDeleteChart(r resource) bool
		CanDeleteChartByID(ID uint64) bool
		CanCreateTrigger() bool
		CanReadTrigger(r resource) bool
		CanUpdateTrigger(r resource) bool
		CanDeleteTrigger(r resource) bool
		CanDeleteTriggerByID(ID uint64) bool
		CanCreatePage() bool
		CanReadPage(r resource) bool
		CanUpdatePage(r resource) bool
		CanDeletePage(r resource) bool
		CanDeletePageByID(ID uint64) bool
	}

	effectivePermission struct {
		Resource  string `json:"resource"`
		Operation string `json:"operation"`
		Allow     bool   `json:"allow"`
	}

	Compose struct{}
)

func (Compose) Resource() internalRules.Resource {
	return internalRules.Resource{Service: "compose"}
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

		rules: p.rules.With(ctx),
	}
}

func (p *permissions) baseResource() resource {
	return &Compose{}
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

	ee = append(ee, ep("compose:namespace:crm", "module.create", p.CanCreateModule()))
	ee = append(ee, ep("compose:namespace:crm", "chart.create", p.CanCreateChart()))
	ee = append(ee, ep("compose:namespace:crm", "trigger.create", p.CanCreateTrigger()))
	ee = append(ee, ep("compose:namespace:crm", "page.create", p.CanCreatePage()))

	return
}

func (p *permissions) CanAccess() bool {
	return p.checkAccess(p.baseResource(), "access")
}

func (p *permissions) CanGrant() bool {
	return p.checkAccess(p.baseResource(), "grant")
}

func (p *permissions) CanCreateNamspace() bool {
	return p.checkAccess(p.baseResource(), "namespace.create")
}

func (p *permissions) CanCreateModule() bool {
	// @todo move to func args when namespaces are implemented
	ns := &types.Namespace{ID: types.NamespaceCRM}
	return p.checkAccess(ns, "module.create")
}

func (p *permissions) CanReadModule(r resource) bool {
	return p.checkAccess(r, "read")
}

func (p *permissions) CanUpdateModule(r resource) bool {
	return p.checkAccess(r, "update")
}

func (p *permissions) CanDeleteModule(r resource) bool {
	return p.checkAccess(r, "delete")
}

func (p *permissions) CanDeleteModuleByID(ID uint64) bool {
	return p.CanDeleteModule(&types.Module{ID: ID})
}

func (p *permissions) CanCreateRecord(r resource) bool {
	return p.checkAccess(r, "record.create")
}

func (p *permissions) CanReadRecord(r resource) bool {
	return p.checkAccess(r, "record.read")
}

func (p *permissions) CanUpdateRecord(r resource) bool {
	return p.checkAccess(r, "record.update")
}

func (p *permissions) CanDeleteRecord(r resource) bool {
	return p.checkAccess(r, "record.delete")
}

func (p *permissions) CanDeleteRecordByModuleID(moduleID uint64) bool {
	return p.CanDeleteRecord(&types.Record{ModuleID: moduleID})
}

func (p *permissions) CanCreateChart() bool {
	// @todo move to func args when namespaces are implemented
	ns := &types.Namespace{ID: types.NamespaceCRM}
	return p.checkAccess(ns, "chart.create")
}

func (p *permissions) CanReadChart(r resource) bool {
	return p.checkAccess(r, "read")
}

func (p *permissions) CanUpdateChart(r resource) bool {
	return p.checkAccess(r, "update")
}

func (p *permissions) CanDeleteChart(r resource) bool {
	return p.checkAccess(r, "delete")
}

func (p *permissions) CanDeleteChartByID(ID uint64) bool {
	return p.CanDeleteChart(&types.Chart{ID: ID})
}

func (p *permissions) CanCreateTrigger() bool {
	// @todo move to func args when namespaces are implemented
	ns := &types.Namespace{ID: types.NamespaceCRM}
	return p.checkAccess(ns, "trigger.create")
}

func (p *permissions) CanReadTrigger(r resource) bool {
	return p.checkAccess(r, "read")
}

func (p *permissions) CanUpdateTrigger(r resource) bool {
	return p.checkAccess(r, "update")
}

func (p *permissions) CanDeleteTrigger(r resource) bool {
	return p.checkAccess(r, "delete")
}

func (p *permissions) CanDeleteTriggerByID(ID uint64) bool {
	return p.CanDeleteTrigger(&types.Trigger{ID: ID})
}

func (p *permissions) CanCreatePage() bool {
	// @todo move to func args when namespaces are implemented
	ns := &types.Namespace{ID: types.NamespaceCRM}
	return p.checkAccess(ns, "page.create")
}

func (p *permissions) CanReadPage(r resource) bool {
	return p.checkAccess(r, "read")
}

func (p *permissions) CanUpdatePage(r resource) bool {
	return p.checkAccess(r, "update")
}

func (p *permissions) CanDeletePage(r resource) bool {
	return p.checkAccess(r, "delete")
}

func (p *permissions) CanDeletePageByID(ID uint64) bool {
	return p.CanDeletePage(&types.Page{ID: ID})
}

func (p *permissions) checkAccess(resource resource, operation string, fallbacks ...internalRules.CheckAccessFunc) bool {
	access := p.rules.Check(resource.Resource().String(), operation, fallbacks...)
	if access == internalRules.Allow {
		return true
	}
	return false
}
