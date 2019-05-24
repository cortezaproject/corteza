package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/http"
	"github.com/cortezaproject/corteza-server/internal/logger"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/internal/store"
	"github.com/cortezaproject/corteza-server/messaging/internal/repository"
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

func Init(ctx context.Context) error {
	fs, err := store.New("var/store")
	if err != nil {
		return err
	}

	client, err := http.New(&http.Config{
		Timeout: 10,
	})
	if err != nil {
		return err
	}

	DefaultLogger = logger.Default().Named("messaging.service")

	DefaultPermissions = permissions.Service(
		ctx,
		DefaultLogger,
		permissions.Repository(repository.DB(ctx), "messaging_permission_rules"))
	DefaultAccessControl = AccessControl(DefaultPermissions)

	DefaultEvent = Event(ctx)
	DefaultChannel = Channel(ctx)
	DefaultAttachment = Attachment(ctx, fs)
	DefaultMessage = Message(ctx)
	DefaultCommand = Command(ctx)
	DefaultWebhook = Webhook(ctx, client)

	return nil
}

func timeNowPtr() *time.Time {
	now := time.Now()
	return &now
}
