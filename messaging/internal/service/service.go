package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/http"
	"github.com/crusttech/crust/internal/logger"
	"github.com/crusttech/crust/internal/store"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	DefaultLogger *zap.Logger

	DefaultAttachment  AttachmentService
	DefaultChannel     ChannelService
	DefaultMessage     MessageService
	DefaultPubSub      *pubSub
	DefaultEvent       EventService
	DefaultPermissions PermissionsService
	DefaultCommand     CommandService
	DefaultWebhook     WebhookService
)

func Init() error {
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

	ctx := context.Background()

	DefaultPermissions = Permissions(ctx)
	DefaultEvent = Event(ctx)
	DefaultAttachment = Attachment(ctx, fs)
	DefaultMessage = Message(ctx)
	DefaultChannel = Channel(ctx)
	DefaultPubSub = PubSub()
	DefaultCommand = Command(ctx)
	DefaultWebhook = Webhook(ctx, client)

	return nil
}

func timeNowPtr() *time.Time {
	now := time.Now()
	return &now
}
