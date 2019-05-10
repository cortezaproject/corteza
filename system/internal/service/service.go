package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/crusttech/crust/internal/logger"
	"github.com/crusttech/crust/internal/permissions"
	internalSettings "github.com/crusttech/crust/internal/settings"
	"github.com/crusttech/crust/system/internal/repository"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
	permissionServicer interface {
		accessControlPermissionServicer
		Watch(ctx context.Context)
	}
)

var (
	permSvc       permissionServicer
	DefaultLogger *zap.Logger

	DefaultAccessControl *accessControl

	DefaultSettings         SettingsService
	DefaultAuthNotification AuthNotificationService
	DefaultAuthSettings     authSettings

	DefaultAuth         AuthService
	DefaultUser         UserService
	DefaultRole         RoleService
	DefaultOrganisation OrganisationService
	DefaultApplication  ApplicationService
)

func Init(ctx context.Context) (err error) {
	intSet := internalSettings.NewService(internalSettings.NewRepository(repository.DB(ctx), "sys_settings"))

	DefaultLogger = logger.Default().Named("system.service")

	permSvc = permissions.Service(
		ctx,
		DefaultLogger,
		permissions.Repository(repository.DB(ctx), "compose_permission_rules"))
	DefaultAccessControl = AccessControl(permSvc)

	DefaultSettings = Settings(ctx, intSet)

	DefaultUser = User(ctx)
	DefaultRole = Role(ctx)
	DefaultOrganisation = Organisation(ctx)
	DefaultApplication = Application(ctx)

	// Authentication helpers & services
	DefaultAuthSettings, err = DefaultSettings.LoadAuthSettings()
	if err != nil {
		return
	}
	DefaultAuthNotification = AuthNotification(ctx)
	DefaultAuth = Auth(ctx)

	return
}

func Watchers(ctx context.Context) {
	permSvc.Watch(ctx)
}
