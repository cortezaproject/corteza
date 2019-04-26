package incoming

type Payload struct {
	// Channel actions
	*Channels    `json:"channels"`
	*ChannelJoin `json:"joinChannel"`
	*ChannelPart `json:"partChannel"`

	*ChannelCreate `json:"createChannel"`
	*ChannelUpdate `json:"updateChannel"`

	*ChannelViewRecord `json:"recordChannelView"`

	*ChannelActivity `json:"channelActivity"`

	// Get channel message history
	*Messages       `json:"messages"`
	*MessageThreads `json:"messageThreads"`

	*MessageActivity `json:"messageActivity"`

	// Message actions
	*MessageCreate `json:"createMessage"`
	*MessageUpdate `json:"updateMessage"`
	*MessageDelete `json:"deleteMessage"`

	*Users `json:"getUsers"`
}
