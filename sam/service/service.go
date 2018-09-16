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
)

func Init() {
	o.Do(func() {
		fs, err := store.New("var/store")
		if err != nil {
			log.Fatalf("Failed to initialize stor: %v", err)
		}

		DefaultAttachment = Attachment(fs)
		DefaultChannel = Channel()
		DefaultMessage = Message(DefaultAttachment)
		DefaultOrganisation = Organisation()
		DefaultPubSub = PubSub()
		DefaultTeam = Team()
	})
}
