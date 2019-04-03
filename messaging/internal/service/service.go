package service

import (
	"context"
	"time"

	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/http"
	"github.com/crusttech/crust/internal/store"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	DefaultAttachment  AttachmentService
	DefaultChannel     ChannelService
	DefaultMessage     MessageService
	DefaultPubSub      *pubSub
	DefaultEvent       EventService
	DefaultPermissions PermissionsService
)

func Init() error {
	fs, err := store.New("var/store")
	if err != nil {
		return err
	}

	_, err = http.New(&config.HTTPClient{
		Timeout: 10,
	})
	if err != nil {
		return err
	}

	ctx := context.Background()

	DefaultPermissions = Permissions(ctx)
	DefaultEvent = Event(ctx)
	DefaultAttachment = Attachment(ctx, fs)
	DefaultMessage = Message(ctx)
	DefaultChannel = Channel(ctx)
	DefaultPubSub = PubSub()

	return nil
}

func timeNowPtr() *time.Time {
	now := time.Now()
	return &now
}
