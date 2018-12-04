package service

import (
	"log"
	"sync"
	"time"

	"github.com/crusttech/crust/internal/store"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	o                  sync.Once
	DefaultAttachment  AttachmentService
	DefaultChannel     ChannelService
	DefaultMessage     MessageService
	DefaultPermissions PermissionsService
	DefaultPubSub      *pubSub
	DefaultEvent       EventService
)

func Init() {
	o.Do(func() {
		fs, err := store.New("var/store")
		if err != nil {
			log.Fatalf("Failed to initialize stor: %v", err)
		}

		DefaultEvent = Event()
		DefaultAttachment = Attachment(fs)
		DefaultMessage = Message()
		DefaultPermissions = Permissions()
		DefaultChannel = Channel()
		DefaultPubSub = PubSub()
	})
}

func timeNowPtr() *time.Time {
	now := time.Now()
	return &now
}
