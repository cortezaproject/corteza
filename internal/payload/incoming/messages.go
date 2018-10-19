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

	Messages struct {
		ChannelID uint64 `json:"channelId,string"`
		FirstID   uint64 `json:"firstID,string"`
		LastID    uint64 `json:"lastID,string"`
		RepliesTo uint64 `json:"repliesTo,string"`
	}
)
