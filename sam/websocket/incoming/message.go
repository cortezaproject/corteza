package incoming

type Payload struct {
	// Channel actions
	*ChannelJoin `json:"chjoin"`
	*ChannelPart `json:"chpart"`

	// Get channel message history
	*ChannelOpen `json:"chopen"`

	// Message actions
	*MessageCreate `json:"msgcre"`
	*MessageUpdate `json:"msgupd"`
	*MessageDelete `json:"msgdel"`
}
