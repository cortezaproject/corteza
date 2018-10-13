package outgoing

type (
	ChannelView struct {
		// Channel to part (nil) for ALL channels
		LastMessageID    string `json:"lastMessageID"`
		NewMessagesCount uint32 `json:"newMessagesCount"`
	}
)
