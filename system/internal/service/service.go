package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/permissions"
	internalSettings "github.com/cortezaproject/corteza-server/internal/settings"
	"github.com/cortezaproject/corteza-server/system/internal/repository"
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
	DefaultPermissions permissionServicer
	DefaultIntSettings internalSettings.Service

	DefaultLogger *zap.Logger

	DefaultAccessControl *accessControl

	DefaultSettings         SettingsService
	DefaultAuthNotification AuthNotificationService
	DefaultAuthSettings     AuthSettings

	DefaultAuth         AuthService
	DefaultUser         UserService
	DefaultRole         RoleService
	DefaultOrganisation OrganisationService
	DefaultApplication  ApplicationService
)

func Init(ctx context.Context, log *zap.Logger) (err error) {
	DefaultLogger = log.Named("service")

	DefaultIntSettings = internalSettings.NewService(internalSettings.NewRepository(repository.DB(ctx), "sys_settings"))

	DefaultPermissions = permissions.Service(
		ctx,
		DefaultLogger,
		permissions.Repository(repository.DB(ctx), "sys_permission_rules"))
	DefaultAccessControl = AccessControl(DefaultPermissions)

	DefaultSettings = Settings(ctx, DefaultIntSettings)

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
	DefaultPermissions.Watch(ctx)
}
