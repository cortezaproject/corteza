package incoming

type (
	Channels struct{}

	ChannelJoin struct {
		ChannelID string `json:"id"`
	}

	ChannelPart struct {
		ChannelID string `json:"id"`
	}

	ChannelViewRecord struct {
		ChannelID     string `json:"channelID"`
		LastMessageID string `json:"lastMessageID"`
	}

	ChannelCreate struct {
		Name  *string `json:"name"`
		Topic *string `json:"topic"`
		Type  *string `json:"type"`
	}

	ChannelUpdate struct {
		ID    string  `json:"id"`
		Name  *string `json:"name"`
		Topic *string `json:"topic"`
		Type  *string `json:"type"`
	}

	ChannelDelete struct {
		ChannelID string `json:"id"`
	}

	// ChannelActivity is sent from the client when there is an activity on the channel...
	ChannelActivity struct {
		ChannelID uint64 `json:"ID,string"`
		Kind      string `json:"kind,omitempty"`
	}
)
