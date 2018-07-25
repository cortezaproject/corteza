package incoming

type ChannelJoin struct {
	ChannelId string `json:"cid"`
}

type ChannelPart struct {
	ChannelId string `json:"cid"`
}

type ChannelOpen struct {
	ChannelId string `json:"cid"`
	Since     string `json:"since,omitempty"`
	Until     string `json:"until,omitempty"`
}

type MessageCreate struct {
	ChannelId string `json:"cid"`
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
