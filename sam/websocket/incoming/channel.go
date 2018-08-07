package incoming

type (
	ChannelList struct{}

	ChannelJoin struct {
		ChannelID string `json:"id"`
	}

	ChannelPart struct {
		ChannelID string `json:"id"`
	}

	ChannelCreate struct {
		Name  string `json:"name"`
		Topic string `json:"topic"`
	}

	ChannelRename struct {
		ChannelID string `json:"id"`
		Name      string `json:"name"`
	}

	ChannelChangeTopic struct {
		ChannelID string `json:"id"`
		Topic     string `json:"topic"`
	}

	ChannelDelete struct {
		ChannelID string `json:"id"`
	}
)
