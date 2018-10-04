package incoming

type Payload struct {
	// Channel actions
	*Channels    `json:"channels"`
	*ChannelJoin `json:"joinChannel"`
	*ChannelPart `json:"partChannel"`

	*ChannelCreate `json:"createChannel"`
	*ChannelUpdate `json:"updateChannel"`
	*ChannelDelete `json:"deleteChannel"`

	// Get channel message history
	*Messages `json:"messages"`

	// Message actions
	*MessageCreate `json:"createMessage"`
	*MessageUpdate `json:"updateMessage"`
	*MessageDelete `json:"deleteMessage"`

	*Users `json:"getUsers"`

	*ExecCommand `json:"exec"`
}
