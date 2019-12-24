package service

import (
	"context"

	"go.uber.org/zap"

	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/automation/corredor"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	db interface {
		Transaction(callback func() error) error
	}

	permissionServicer interface {
		accessControlPermissionServicer
		Watch(ctx context.Context)
	}

	automationManager interface {
		automationScriptManager
		automationTriggerManager
		automationScriptsFinder
	}

	Config struct {
		Storage          options.StorageOpt
		Corredor         options.CorredorOpt
		GRPCClientSystem options.GRPCServerOpt
	}

	permitChecker interface {
		Validate(string, bool) error
		CanCreateUser(uint) error
		CanRegister(uint) error
	}
)

var (
	DefaultLogger *zap.Logger

	// CurrentSubscription holds current subscription info,
	// and functions for domain validation, user limit checks and
	// warning texts
	//
	// By default, Corteza (community edition) has this set to nil
	// and with that all checks & validations are skipped
	//
	// Other flavours or distributions can set this to
	// something that suits their needs.
	CurrentSubscription permitChecker

	// DefaultPermissions Retrieves & stores permissions
	DefaultPermissions permissionServicer

	// DefaultSettings controls system's settings
	DefaultSettings settings.Service

	// DefaultAccessControl Access control checking
	DefaultAccessControl *accessControl

	// DefaultInternalAutomationManager manages automation scripts, triggers, runnable scripts
	DefaultInternalAutomationManager automationManager

	// DefaultAutomationScriptManager manages scripts
	DefaultAutomationScriptManager automationScript

	// DefaultAutomationTriggerManager manages triggerManager
	DefaultAutomationTriggerManager automationTrigger

	// DefaultAutomationRunner runs automation scripts by listening to triggerManager and invoking Corredor service
	DefaultAutomationRunner automationRunner

	DefaultAuthNotification AuthNotificationService

	// CurrentSettings represents current system settings
	CurrentSettings = &types.Settings{}

	DefaultSink *sink

	DefaultAuth         AuthService
	DefaultUser         UserService
	DefaultRole         RoleService
	DefaultOrganisation OrganisationService
	DefaultApplication  ApplicationService
	DefaultReminder     ReminderService

	DefaultStatistics *statistics
)

func Initialize(ctx context.Context, log *zap.Logger, c Config) (err error) {
	DefaultLogger = log.Named("service")

	if DefaultPermissions == nil {
		// Do not override permissions service stored under DefaultPermissions
		// to allow integration tests to inject own permission service
		DefaultPermissions = permissions.Service(ctx, DefaultLogger, repository.DB(ctx), "sys_permission_rules")
	}

	DefaultAccessControl = AccessControl(DefaultPermissions)

	DefaultSettings = settings.NewService(
		settings.NewRepository(repository.DB(ctx), "sys_settings"),
		DefaultLogger,
		DefaultAccessControl,
		CurrentSettings,
	)

	DefaultUser = User(ctx)
	DefaultRole = Role(ctx)
	DefaultOrganisation = Organisation(ctx)
	DefaultApplication = Application(ctx)
	DefaultReminder = Reminder(ctx)

	DefaultAuthNotification = AuthNotification(ctx)
	DefaultAuth = Auth(ctx)

	{
		if DefaultInternalAutomationManager == nil {
			// handles script & trigger management & keeping runnables cripts in internal cache
			DefaultInternalAutomationManager = automation.Service(automation.AutomationServiceConfig{
				Logger:        DefaultLogger,
				DbTablePrefix: "sys",
				DB:            repository.DB(ctx),
				TokenMaker: func(ctx context.Context, userID uint64) (jwt string, err error) {
					var u *types.User

					ctx = intAuth.SetSuperUserContext(ctx)
					if u, err = DefaultUser.FindByID(userID); err != nil {
						return
					} else if err = DefaultAuth.LoadRoleMemberships(u); err != nil {
						return
					}

					return intAuth.DefaultJwtHandler.Encode(u), nil
				},
			})
		}

		// Pass automation manager to
		DefaultAutomationTriggerManager = AutomationTrigger(DefaultInternalAutomationManager)
		DefaultAutomationScriptManager = AutomationScript(DefaultInternalAutomationManager)

		var scriptRunnerClient corredor.ScriptRunnerClient

		if c.Corredor.Enabled {
			conn, err := corredor.NewConnection(ctx, c.Corredor, DefaultLogger)

			log.Info("initializing corredor connection", zap.String("addr", c.Corredor.Addr), zap.Error(err))
			if err != nil {
				return err
			}

			scriptRunnerClient = corredor.NewScriptRunnerClient(conn)
		}

		DefaultAutomationRunner = AutomationRunner(
			AutomationRunnerOpt{
				ApiBaseURLSystem:    c.Corredor.ApiBaseURLSystem,
				ApiBaseURLMessaging: c.Corredor.ApiBaseURLMessaging,
				ApiBaseURLCompose:   c.Corredor.ApiBaseURLCompose,
			},
			DefaultInternalAutomationManager,
			scriptRunnerClient,
		)
	}

	DefaultSink = Sink()
	DefaultStatistics = Statistics(ctx)

	return
}

func Activate(ctx context.Context) (err error) {
	// Run initial update of current settings with super-user credentials
	err = DefaultSettings.UpdateCurrent(intAuth.SetSuperUserContext(ctx))
	if err != nil {
		return
	}

	return
}

func Watchers(ctx context.Context) {
	// Reloading automation scripts on change
	DefaultAutomationRunner.Watch(ctx)

	// Reloading permissions on change
	DefaultPermissions.Watch(ctx)
}
