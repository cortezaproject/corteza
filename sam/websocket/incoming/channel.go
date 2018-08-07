package incoming

type (
	ChannelList struct{}

	ChannelJoin struct {
		ChannelID string `json:"cid"`
	}

	ChannelPart struct {
		ChannelID string `json:"cid"`
	}

	ChannelCreate struct {
		Name  string `json:"name"`
		Topic string `json:"topic"`
	}

	ChannelRename struct {
		ChannelID string `json:"cid"`
		Name      string `json:"name"`
	}

	ChannelChangeTopic struct {
		ChannelID string `json:"cid"`
		Topic     string `json:"topic"`
	}

	ChannelDelete struct {
		ChannelID string `json:"cid"`
	}
)
