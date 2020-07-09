package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	actionlogRepository "github.com/cortezaproject/corteza-server/pkg/actionlog/repository"
	"github.com/cortezaproject/corteza-server/pkg/app/options"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/pkg/store"
	"github.com/cortezaproject/corteza-server/pkg/store/minio"
	"github.com/cortezaproject/corteza-server/pkg/store/plain"
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
		ActionLog options.ActionLogOpt
		Storage   options.StorageOpt
	}
)

var (
	DefaultStore       store.Store
	DefaultPermissions permissionServicer

	DefaultLogger *zap.Logger

	DefaultActionlog actionlog.Recorder

	DefaultSettings      settings.Service
	DefaultAccessControl *accessControl

	// CurrentSettings represents current messaging settings
	CurrentSettings = &types.Settings{}

	DefaultAttachment AttachmentService
	DefaultChannel    ChannelService
	DefaultMessage    MessageService
	DefaultEvent      EventService
	DefaultCommand    CommandService
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
			// will log directly to system schema for now
			actionlogRepository.Mysql(repository.DB(ctx).Quiet(), "sys_actionlog"),
			log,
			tee,
			policy,
		)
	}

	if DefaultPermissions == nil {
		// Do not override permissions service stored under DefaultPermissions
		// to allow integration tests to inject own permission service
		DefaultPermissions = permissions.Service(ctx, DefaultLogger, repository.DB(ctx), "messaging_permission_rules")
	}

	DefaultAccessControl = AccessControl(DefaultPermissions)

	DefaultSettings = settings.NewService(
		settings.NewRepository(repository.DB(ctx), "messaging_settings"),
		DefaultLogger,
		DefaultAccessControl,
		CurrentSettings,
	)

	if DefaultStore == nil {
		const svcPath = "messaging"
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

	DefaultEvent = Event(ctx)
	DefaultChannel = Channel(ctx)
	DefaultAttachment = Attachment(ctx, DefaultStore)
	DefaultMessage = Message(ctx)
	DefaultCommand = Command(ctx)

	return nil
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
	DefaultPermissions.Watch(ctx)
}

func timeNowPtr() *time.Time {
	now := time.Now()
	return &now
}
