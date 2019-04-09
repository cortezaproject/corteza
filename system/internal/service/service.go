package service

import (
	"context"

	internalSettings "github.com/crusttech/crust/internal/settings"
	"github.com/crusttech/crust/system/internal/repository"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	DefaultSettings         SettingsService
	DefaultAuthNotification AuthNotificationService
	DefaultAuthSettings     authSettings

	DefaultAuth         AuthService
	DefaultUser         UserService
	DefaultRole         RoleService
	DefaultRules        RulesService
	DefaultOrganisation OrganisationService
	DefaultApplication  ApplicationService
	DefaultPermissions  PermissionsService
)

func Init() (err error) {
	ctx := context.Background()

	intSet := internalSettings.NewService(internalSettings.NewRepository(repository.DB(ctx), "sys_settings"))

	DefaultSettings = Settings(ctx, intSet)
	DefaultRules = Rules(ctx)
	DefaultPermissions = Permissions(ctx)

	DefaultUser = User(ctx)
	DefaultRole = Role(ctx)
	DefaultOrganisation = Organisation(ctx)
	DefaultApplication = Application(ctx)

	// Authentication helpers & services
	DefaultAuthSettings, err = DefaultSettings.LoadAuthSettings()
	DefaultAuthSettings.internalEnabled = true
	DefaultAuthSettings.internalSignUpEmailConfirmationRequired = true
	DefaultAuthSettings.internalPasswordResetEnabled = true
	DefaultAuthSettings.mailFromAddress = "denis.arh@gmail.com"
	if err != nil {
		return
	}
	DefaultAuthNotification = AuthNotification(ctx)
	DefaultAuth = Auth(ctx)

	return
}
