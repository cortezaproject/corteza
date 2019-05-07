package service

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/crusttech/crust/compose/internal/repository"
	"github.com/crusttech/crust/compose/types"
	"github.com/crusttech/crust/internal/logger"
	internalRules "github.com/crusttech/crust/internal/rules"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	permissions struct {
		db     db
		ctx    context.Context
		logger *zap.Logger

		rules systemService.RulesService
	}

	permissionResource interface {
		PermissionResource() internalRules.Resource
	}

	PermissionsService interface {
		With(context.Context) PermissionsService

		Effective() (ee []effectivePermission, err error)

		CanAccess() bool
		CanGrant() bool
		CanCreateNamespace() bool
		CanReadNamespace(r permissionResource) bool
		CanUpdateNamespace(r permissionResource) bool
		CanDeleteNamespace(r permissionResource) bool
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
		logger: DefaultLogger.Named("permissions"),
		rules:  systemService.DefaultRules,
	}).With(context.Background())
}

func (svc permissions) With(ctx context.Context) PermissionsService {
	db := repository.DB(ctx)
	return &permissions{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		rules: systemService.Rules(ctx),
	}
}

// log() returns zap's logger with requestID from current context and fields.
func (svc permissions) log(fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
}

// Return permissions
func (svc permissions) Effective() (ee []effectivePermission, err error) {
	ep := func(res, op string, allow bool) effectivePermission {
		return effectivePermission{
			Resource:  res,
			Operation: op,
			Allow:     allow,
		}
	}

	ee = append(ee, ep("compose", "access", svc.CanAccess()))
	ee = append(ee, ep("compose", "grant", svc.CanGrant()))
	ee = append(ee, ep("compose", "namespace.create", svc.CanCreateNamespace()))

	ee = append(ee, ep("compose:namespace:crm", "module.create", svc.CanCreateModule(crmNamespace())))
	ee = append(ee, ep("compose:namespace:crm", "chart.create", svc.CanCreateChart(crmNamespace())))
	ee = append(ee, ep("compose:namespace:crm", "trigger.create", svc.CanCreateTrigger(crmNamespace())))
	ee = append(ee, ep("compose:namespace:crm", "page.create", svc.CanCreatePage(crmNamespace())))

	return
}

func (svc permissions) CanAccess() bool {
	return svc.checkAccess(types.PermissionResource, "access")
}

func (svc permissions) CanGrant() bool {
	return svc.checkAccess(types.PermissionResource, "grant")
}

func (svc permissions) CanCreateNamespace() bool {
	return svc.checkAccess(types.PermissionResource, "namespace.create")
}

func (svc permissions) CanReadNamespace(r permissionResource) bool {
	return svc.checkAccess(r, "read", svc.allow())
}

func (svc permissions) CanUpdateNamespace(r permissionResource) bool {
	return svc.checkAccess(r, "update")
}

func (svc permissions) CanDeleteNamespace(r permissionResource) bool {
	return svc.checkAccess(r, "delete")
}

func (svc permissions) CanCreateModule(r permissionResource) bool {
	return svc.checkAccess(r, "module.create")
}

func (svc permissions) CanReadModule(r permissionResource) bool {
	return svc.checkAccess(r, "read")
}

func (svc permissions) CanUpdateModule(r permissionResource) bool {
	return svc.checkAccess(r, "update")
}

func (svc permissions) CanDeleteModule(r permissionResource) bool {
	return svc.checkAccess(r, "delete")
}

func (svc permissions) CanCreateRecord(r permissionResource) bool {
	return svc.checkAccess(r, "record.create")
}

func (svc permissions) CanReadRecord(r permissionResource) bool {
	return svc.checkAccess(r, "record.read")
}

func (svc permissions) CanUpdateRecord(r permissionResource) bool {
	return svc.checkAccess(r, "record.update")
}

func (svc permissions) CanDeleteRecord(r permissionResource) bool {
	return svc.checkAccess(r, "record.delete")
}

func (svc permissions) CanCreateChart(r permissionResource) bool {
	return svc.checkAccess(r, "chart.create")
}

func (svc permissions) CanReadChart(r permissionResource) bool {
	return svc.checkAccess(r, "read")
}

func (svc permissions) CanUpdateChart(r permissionResource) bool {
	return svc.checkAccess(r, "update")
}

func (svc permissions) CanDeleteChart(r permissionResource) bool {
	return svc.checkAccess(r, "delete")
}

func (svc permissions) CanCreateTrigger(r permissionResource) bool {
	return svc.checkAccess(r, "trigger.create")
}

func (svc permissions) CanReadTrigger(r permissionResource) bool {
	return svc.checkAccess(r, "read")
}

func (svc permissions) CanUpdateTrigger(r permissionResource) bool {
	return svc.checkAccess(r, "update")
}

func (svc permissions) CanDeleteTrigger(r permissionResource) bool {
	return svc.checkAccess(r, "delete")
}

func (svc permissions) CanCreatePage(r permissionResource) bool {
	// @todo move to func args when namespaces are implemented
	return svc.checkAccess(r, "page.create")
}

func (svc permissions) CanReadPage(r permissionResource) bool {
	return svc.checkAccess(r, "read")
}

func (svc permissions) CanUpdatePage(r permissionResource) bool {
	return svc.checkAccess(r, "update")
}

func (svc permissions) CanDeletePage(r permissionResource) bool {
	return svc.checkAccess(r, "delete")
}

func (svc permissions) checkAccess(resource permissionResource, operation string, fallbacks ...internalRules.CheckAccessFunc) bool {
	access := svc.rules.Check(resource.PermissionResource(), operation, fallbacks...)
	if access == internalRules.Allow {
		return true
	}
	return false
}

func (svc permissions) allow() func() internalRules.Access {
	return func() internalRules.Access {
		return internalRules.Allow
	}
}
