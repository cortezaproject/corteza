package service

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/crusttech/crust/internal/logger"
	internalRules "github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/types"
)

type (
	permissions struct {
		db     db
		ctx    context.Context
		logger *zap.Logger

		rules RulesService
	}

	resource interface {
		PermissionResource() internalRules.Resource
	}

	PermissionsService interface {
		With(context.Context) PermissionsService

		Effective() (ee []effectivePermission, err error)

		CanReadSettings() bool
		CanManageSettings() bool

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

func Permissions(ctx context.Context) PermissionsService {
	return (&permissions{
		logger: DefaultLogger.Named("permissions"),
	}).With(ctx)
}

func (svc permissions) With(ctx context.Context) PermissionsService {
	db := repository.DB(ctx)
	return &permissions{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		rules: Rules(ctx),
	}
}

// log() returns zap's logger with requestID from current context and fields.
func (svc permissions) log(fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
}

func (svc permissions) Effective() (ee []effectivePermission, err error) {
	ep := func(res, op string, allow bool) effectivePermission {
		return effectivePermission{
			Resource:  res,
			Operation: op,
			Allow:     allow,
		}
	}

	ee = append(ee, ep("system", "access", svc.CanAccess()))
	ee = append(ee, ep("system", "settings.read", svc.CanReadSettings()))
	ee = append(ee, ep("system", "settings.manage", svc.CanManageSettings()))
	ee = append(ee, ep("system", "application.create", svc.CanCreateApplication()))
	ee = append(ee, ep("system", "role.create", svc.CanCreateRole()))
	ee = append(ee, ep("system", "organisation.create", svc.CanCreateOrganisation()))
	ee = append(ee, ep("system", "grant", svc.CanCreateRole()))

	return
}

func (svc permissions) CanAccess() bool {
	return svc.checkAccess(types.PermissionResource, "access")
}

func (svc permissions) CanReadSettings() bool {
	return svc.checkAccess(types.PermissionResource, "settings.read")
}

func (svc permissions) CanManageSettings() bool {
	return svc.checkAccess(types.PermissionResource, "settings.manage")
}

func (svc permissions) CanCreateOrganisation() bool {
	return svc.checkAccess(types.PermissionResource, "organisation.create")
}

func (svc permissions) CanCreateUser() bool {
	return svc.checkAccess(types.PermissionResource, "user.create", svc.allow())
}

func (svc permissions) CanCreateRole() bool {
	return svc.checkAccess(types.PermissionResource, "role.create")
}

func (svc permissions) CanCreateApplication() bool {
	return svc.checkAccess(types.PermissionResource, "application.create")
}

func (svc permissions) CanReadRole(rl *types.Role) bool {
	return svc.checkAccess(rl, "read", svc.allow())
}

func (svc permissions) CanUpdateRole(rl *types.Role) bool {
	return svc.checkAccess(rl, "update")
}

func (svc permissions) CanDeleteRole(rl *types.Role) bool {
	return svc.checkAccess(rl, "delete")
}

func (svc permissions) CanManageRoleMembers(rl *types.Role) bool {
	return svc.checkAccess(rl, "members.manage")
}

func (svc permissions) CanReadApplication(app *types.Application) bool {
	return svc.checkAccess(app, "read", svc.allow())
}

func (svc permissions) CanUpdateApplication(app *types.Application) bool {
	return svc.checkAccess(app, "update")
}

func (svc permissions) CanDeleteApplication(app *types.Application) bool {
	return svc.checkAccess(app, "delete")
}

func (svc permissions) CanUpdateUser(u *types.User) bool {
	return svc.checkAccess(u, "update")
}

func (svc permissions) CanSuspendUser(u *types.User) bool {
	return svc.checkAccess(u, "suspend")
}

func (svc permissions) CanUnsuspendUser(u *types.User) bool {
	return svc.checkAccess(u, "unsuspend")
}

func (svc permissions) CanDeleteUser(u *types.User) bool {
	return svc.checkAccess(u, "delete")
}

func (svc permissions) checkAccess(r resource, operation string, fallbacks ...internalRules.CheckAccessFunc) bool {
	return svc.rules.Check(r.PermissionResource(), operation, fallbacks...) == internalRules.Allow
}

func (svc permissions) allow() func() internalRules.Access {
	return func() internalRules.Access {
		return internalRules.Allow
	}
}
