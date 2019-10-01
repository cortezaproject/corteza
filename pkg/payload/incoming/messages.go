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
)
