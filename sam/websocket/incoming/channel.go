package incoming

type (
	ChannelList struct{}

	ChannelJoin struct {
		ChannelID string `json:"cid"`
	}

	ChannelPart struct {
		ChannelID string `json:"cid"`
	}

	ChannelOpen struct {
		ChannelID string `json:"cid"`
		Since     string `json:"since,omitempty"`
		Until     string `json:"until,omitempty"`
	}
)
