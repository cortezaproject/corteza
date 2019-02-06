package service

import (
	"log"
	"sync"
	"time"

	"github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/internal/store"

	"github.com/crusttech/crust/messaging/types"
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

		scopes := rules.NewScope()
		scopes.Add(&types.Organisation{})
		scopes.Add(&types.Team{})
		scopes.Add(&types.Channel{})

		DefaultEvent = Event()
		DefaultAttachment = Attachment(fs)
		DefaultMessage = Message()
		DefaultPermissions = Permissions(scopes)
		DefaultChannel = Channel()
		DefaultPubSub = PubSub()
	})
}

func timeNowPtr() *time.Time {
	now := time.Now()
	return &now
}
