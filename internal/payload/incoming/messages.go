package incoming

type (
	MessageCreate struct {
		ChannelID string `json:"channelID"`
		ReplyTo   uint64 `json:"replyTo,omitempty,string"`
		Message   string `json:"message"`
	}

	MessageUpdate struct {
		ID      string `json:"messageID"`
		Message string `json:"message"`
	}

	MessageDelete struct {
		ID string `json:"messageID"`
	}

	// MessageActivity is sent from the client when there is an activity on the message...
	MessageActivity struct {
		MessageID uint64 `json:"messageID,string"`
		ChannelID uint64 `json:"channelID,string"`
		Kind      string `json:"kind,omitempty"`
	}
)
