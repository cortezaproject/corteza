package incoming

type (
	MessageCreate struct {
		ChannelID string `json:"channelId"`
		Message   string `json:"message"`
	}

	MessageUpdate struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	}

	MessageDelete struct {
		ChannelID string `json:"channelId"`
		ID        string `json:"id"`
	}

	Messages struct {
		ChannelID string `json:"channelId"`
		FromID    string `json:"fromId,omitempty"`
		UntilID   string `json:"untilId,omitempty"`
	}
)
