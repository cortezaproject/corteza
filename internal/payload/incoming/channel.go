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
)
