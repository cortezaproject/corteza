package service

import (
	"context"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	actionlogRepository "github.com/cortezaproject/corteza-server/pkg/actionlog/repository"
	"github.com/cortezaproject/corteza-server/pkg/app/options"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/pkg/store"
	"github.com/cortezaproject/corteza-server/pkg/store/minio"
	"github.com/cortezaproject/corteza-server/pkg/store/plain"
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
		ActionLog        options.ActionLogOpt
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
	DefaultStore store.Store

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

	DefaultActionlog actionlog.Recorder

	DefaultSink *sink

	DefaultAuth         AuthService
	DefaultUser         UserService
	DefaultRole         RoleService
	DefaultOrganisation OrganisationService
	DefaultApplication  ApplicationService
	DefaultReminder     ReminderService
	DefaultAttachment   AttachmentService

	DefaultStatistics *statistics
)

func Initialize(ctx context.Context, log *zap.Logger, c Config) (err error) {
	DefaultLogger = log.Named("service")

	{
		tee := zap.NewNop()
		policy := actionlog.MakeProductionPolicy()

		if !c.ActionLog.Enabled {
			policy = actionlog.MakeDisabledPolicy()
		} else if c.ActionLog.Debug {
			policy = actionlog.MakeDebugPolicy()
			tee = log
		}

		DefaultActionlog = actionlog.NewService(
			actionlogRepository.Mysql(repository.DB(ctx).Quiet(), "sys_actionlog"),
			log,
			tee,
			policy,
		)
	}

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

	if DefaultStore == nil {
		const svcPath = "compose"
		if c.Storage.MinioEndpoint != "" {
			var bucket = svcPath
			if c.Storage.MinioBucket != "" {
				bucket = c.Storage.MinioBucket + "/" + svcPath
			}

			DefaultStore, err = minio.New(bucket, minio.Options{
				Endpoint:        c.Storage.MinioEndpoint,
				Secure:          c.Storage.MinioSecure,
				Strict:          c.Storage.MinioStrict,
				AccessKeyID:     c.Storage.MinioAccessKey,
				SecretAccessKey: c.Storage.MinioSecretKey,

				ServerSideEncryptKey: []byte(c.Storage.MinioSSECKey),
			})

			log.Info("initializing minio",
				zap.String("bucket", bucket),
				zap.String("endpoint", c.Storage.MinioEndpoint),
				zap.Error(err))
		} else {
			path := c.Storage.Path + "/" + svcPath
			DefaultStore, err = plain.New(path)
			log.Info("initializing store",
				zap.String("path", path),
				zap.Error(err))
		}

		if err != nil {
			return err
		}
	}

	DefaultAuthNotification = AuthNotification(ctx)
	DefaultAuth = Auth(ctx)
	DefaultUser = User(ctx)
	DefaultRole = Role(ctx)
	DefaultOrganisation = Organisation(ctx)
	DefaultApplication = Application(ctx)
	DefaultReminder = Reminder(ctx)
	DefaultSink = Sink()
	DefaultStatistics = Statistics()
	DefaultAttachment = Attachment(DefaultStore)

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
