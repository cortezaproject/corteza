package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/http"
	"github.com/crusttech/crust/internal/logger"
	"github.com/crusttech/crust/internal/permissions"
	"github.com/crusttech/crust/internal/store"
	"github.com/crusttech/crust/messaging/internal/repository"
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
	permSvc permissionServicer

	DefaultLogger *zap.Logger

	DefaultAccessControl *accessControl

	DefaultAttachment AttachmentService
	DefaultChannel    ChannelService
	DefaultMessage    MessageService
	DefaultPubSub     *pubSub
	DefaultEvent      EventService
	DefaultCommand    CommandService
	DefaultWebhook    WebhookService
)

func Init(ctx context.Context) error {
	fs, err := store.New("var/store")
	if err != nil {
		return err
	}

	client, err := http.New(&config.HTTPClient{
		Timeout: 10,
	})
	if err != nil {
		return err
	}

	DefaultLogger = logger.Default().Named("messaging.service")

	permSvc = permissions.Service(
		ctx,
		DefaultLogger,
		permissions.Repository(repository.DB(ctx), "compose_permission_rules"))
	DefaultAccessControl = AccessControl(permSvc)

	DefaultEvent = Event(ctx)
	DefaultChannel = Channel(ctx)
	DefaultAttachment = Attachment(ctx, fs)
	DefaultMessage = Message(ctx)
	DefaultPubSub = PubSub()
	DefaultCommand = Command(ctx)
	DefaultWebhook = Webhook(ctx, client)

	return nil
}

func Watchers(ctx context.Context) {
	permSvc.Watch(ctx)
}

func timeNowPtr() *time.Time {
	now := time.Now()
	return &now
}
