package incoming

import (
	"time"
)

type Message struct {
	// Channel actions
	*ChannelJoin `json:"chjoin"`
	*ChannelPart `json:"chpart"`

	// Get channel message history
	*ChannelOpen `json:"chopen"`

	// Message actions
	*MessageCreate `json:"msgcre"`
	*MessageUpdate `json:"msgupd"`
	*MessageDelete `json:"msgdel"`

	timestamp time.Time
}
