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
)

var (
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

func Init() (err error) {
	ctx := context.Background()

	intSet := internalSettings.NewService(internalSettings.NewRepository(repository.DB(ctx), "sys_settings"))

	DefaultLogger = logger.Default().Named("system.service")

	pv := permissions.Service(permissions.Repository(repository.DB(ctx), "compose_permission_rules"))
	DefaultAccessControl = AccessControl(pv)

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
