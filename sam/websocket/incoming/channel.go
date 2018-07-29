package incoming

type (
	ChannelList struct{}

	ChannelJoin struct {
		ChannelID string `json:"cid"`
	}

	ChannelPart struct {
		ChannelID string `json:"cid"`
	}

	ChannelRename struct {
		ChannelID string `json:"cid"`
		Name      string `json:"name"`
	}

	ChannelChangeTopic struct {
		ChannelID string `json:"cid"`
		Topic     string `json:"topic"`
	}
)
