package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
	"github.com/cortezaproject/corteza-server/pkg/http"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	internalSettings "github.com/cortezaproject/corteza-server/pkg/settings"
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
		Storage options.StorageOpt
	}
)

var (
	DefaultStore       store.Store
	DefaultPermissions permissionServicer

	DefaultLogger *zap.Logger

	DefaultInternalSettings internalSettings.Service
	DefaultSettings         SettingsService
	DefaultAccessControl    *accessControl

	DefaultAttachment AttachmentService
	DefaultChannel    ChannelService
	DefaultMessage    MessageService
	DefaultEvent      EventService
	DefaultCommand    CommandService
	DefaultWebhook    WebhookService
)

func Init(ctx context.Context, log *zap.Logger, c Config) (err error) {
	DefaultLogger = log.Named("service")

	DefaultInternalSettings = internalSettings.NewService(internalSettings.NewRepository(repository.DB(ctx), "messaging_settings"))

	if DefaultStore == nil {
		if c.Storage.MinioEndpoint != "" {
			if c.Storage.MinioBucket == "" {
				c.Storage.MinioBucket = "messaging"
			}

			DefaultStore, err = minio.New(c.Storage.MinioBucket, minio.Options{
				Endpoint:        c.Storage.MinioEndpoint,
				Secure:          c.Storage.MinioSecure,
				Strict:          c.Storage.MinioStrict,
				AccessKeyID:     c.Storage.MinioAccessKey,
				SecretAccessKey: c.Storage.MinioSecretKey,

				ServerSideEncryptKey: []byte(c.Storage.MinioSSECKey),
			})

			log.Info("initializing minio",
				zap.String("bucket", c.Storage.MinioBucket),
				zap.String("endpoint", c.Storage.MinioEndpoint),
				zap.Error(err))
		} else {
			DefaultStore, err = plain.New(c.Storage.Path)

			log.Info("initializing store",
				zap.String("path", c.Storage.Path),
				zap.Error(err))
		}

		if err != nil {
			return err
		}
	}

	client, err := http.New(&http.Config{
		Timeout: 10,
	})
	if err != nil {
		return err
	}

	if DefaultPermissions == nil {
		// Do not override permissions service stored under DefaultPermissions
		// to allow integration tests to inject own permission service
		DefaultPermissions = permissions.Service(ctx, DefaultLogger, repository.DB(ctx), "messaging_permission_rules")
	}
	DefaultAccessControl = AccessControl(DefaultPermissions)

	DefaultSettings = Settings(ctx, DefaultInternalSettings)

	DefaultEvent = Event(ctx)
	DefaultChannel = Channel(ctx)
	DefaultAttachment = Attachment(ctx, DefaultStore)
	DefaultMessage = Message(ctx)
	DefaultCommand = Command(ctx)
	DefaultWebhook = Webhook(ctx, client)

	return nil
}

func Watchers(ctx context.Context) {
	DefaultPermissions.Watch(ctx)
}

func timeNowPtr() *time.Time {
	now := time.Now()
	return &now
}
