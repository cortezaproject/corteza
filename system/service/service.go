package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/app/options"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
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

	Config struct {
		Storage          options.StorageOpt
		GRPCClientSystem options.GRPCServerOpt
	}

	permitChecker interface {
		Validate(string, bool) error
		CanCreateUser(uint) error
		CanRegister(uint) error
	}

	eventDispatcher interface {
		WaitFor(ctx context.Context, ev eventbus.Event) (err error)
		Dispatch(ctx context.Context, ev eventbus.Event)
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
	// Reloading permissions on change
	DefaultPermissions.Watch(ctx)
}
