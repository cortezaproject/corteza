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
		FromID    uint64 `json:"fromID,string"`
		ToID      uint64 `json:"toID,string"`
		FirstID   uint64 `json:"firstID,string"`
		LastID    uint64 `json:"lastID,string"`
		RepliesTo uint64 `json:"repliesTo,string"`

		PinnedOnly     bool `json:pinned`
		BookmarkedOnly bool `json:bookmarked`
	}

	MessageThreads struct {
		ChannelID uint64 `json:"channelId,string"`
		FirstID   uint64 `json:"firstID,string"`
		LastID    uint64 `json:"lastID,string"`
	}

	// MessageActivity is sent from the client when there is an activity on the message...
	MessageActivity struct {
		MessageID uint64 `json:"messageID,string"`
		ChannelID uint64 `json:"channelID,string"`
		Kind      string `json:"kind,omitempty"`
	}
)
