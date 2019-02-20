package service

import (
	"log"
	"sync"
	"time"

	internalRules "github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/internal/store"

	"github.com/crusttech/crust/messaging/types"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	o                 sync.Once
	DefaultAttachment AttachmentService
	DefaultChannel    ChannelService
	DefaultMessage    MessageService
	DefaultPubSub     *pubSub
	DefaultEvent      EventService
)

func Init() {
	o.Do(func() {
		fs, err := store.New("var/store")
		if err != nil {
			log.Fatalf("Failed to initialize stor: %v", err)
		}

		scopes := internalRules.NewScope()
		scopes.Add(&types.Organisation{})
		scopes.Add(&types.Role{})
		scopes.Add(&types.Channel{})

		DefaultEvent = Event()
		DefaultAttachment = Attachment(fs)
		DefaultMessage = Message()
		DefaultChannel = Channel()
		DefaultPubSub = PubSub()
	})
}

func timeNowPtr() *time.Time {
	now := time.Now()
	return &now
}
