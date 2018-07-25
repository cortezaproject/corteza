package incoming

type (
	ChannelJoin struct {
		ChannelID string `json:"cid"`
	}

	ChannelPart struct {
		ChannelID string `json:"cid"`
	}

	ChannelPartAll struct {
		Leave bool `json:"leave"`
	}

	ChannelOpen struct {
		ChannelID string `json:"cid"`
		Since     string `json:"since,omitempty"`
		Until     string `json:"until,omitempty"`
	}
)
