package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
	"github.com/cortezaproject/corteza-server/pkg/http"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/store"
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

	DefaultAccessControl *accessControl

	DefaultAttachment AttachmentService
	DefaultChannel    ChannelService
	DefaultMessage    MessageService
	DefaultEvent      EventService
	DefaultCommand    CommandService
	DefaultWebhook    WebhookService
)

func Init(ctx context.Context, log *zap.Logger, c Config) (err error) {
	DefaultLogger = log.Named("service")

	if DefaultStore == nil {
		DefaultStore, err = store.New(c.Storage.Path)
		log.Info("initializing store", zap.String("path", c.Storage.Path), zap.Error(err))
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
		pRepo := permissions.Repository(repository.DB(ctx), "messaging_permission_rules")
		DefaultPermissions = permissions.Service(ctx, DefaultLogger, pRepo)
	}
	DefaultAccessControl = AccessControl(DefaultPermissions)

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
