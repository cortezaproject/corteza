package service

import (
	"log"
	"sync"

	"github.com/crusttech/crust/internal/store"
)

var (
	o                   sync.Once
	DefaultAttachment   AttachmentService
	DefaultChannel      ChannelService
	DefaultMessage      MessageService
	DefaultOrganisation OrganisationService
	DefaultPubSub       *pubSub
	DefaultTeam         TeamService
	DefaultEvent        EventService
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
		DefaultChannel = Channel()
		DefaultOrganisation = Organisation()
		DefaultPubSub = PubSub()
		DefaultTeam = Team()
	})
}
