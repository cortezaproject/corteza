package incoming

type (
	MessageCreate struct {
		ChannelID string `json:"cid"`
		Message   string `json:"msg"`
	}

	MessageUpdate struct {
		ID      string `json:"id"`
		Message string `json:"msg"`
	}

	MessageDelete struct {
		ChannelID string `json:"cid"`
		ID        string `json:"id"`
	}

	MessageHistory struct {
		ChannelID string `json:"cid"`
		FromID    string `json:"fid,omitempty"`
		UntilID   string `json:"uid,omitempty"`
	}
)
