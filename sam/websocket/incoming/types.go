package incoming

type ChannelJoin struct {
	ChannelID string `json:"cid"`
}

type ChannelPart struct {
	ChannelID string `json:"cid"`
}

type ChannelOpen struct {
	ChannelID string `json:"cid"`
	Since     string `json:"since,omitempty"`
	Until     string `json:"until,omitempty"`
}

type MessageCreate struct {
	ChannelID string `json:"cid"`
	Message   string `json:"msg"`
}

type MessageUpdate struct {
	ID      string `json:"id"`
	Message string `json:"msg"`
}

type MessageDelete struct {
	ID    string `json:"id"`
	Topic string `json:"topic"`
}
