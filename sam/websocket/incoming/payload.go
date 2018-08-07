package incoming

type Payload struct {
	// Channel actions
	*ChannelList `json:"chlist"`
	*ChannelJoin `json:"chjoin"`
	*ChannelPart `json:"chpart"`

	*ChannelChangeTopic `json:"chct"`
	*ChannelRename      `json:"chrn"`
	*ChannelCreate      `json:"chcr"`
	*ChannelDelete      `json:"chdel"`

	// Get channel message history
	*MessageHistory `json:"chopen"`

	// Message actions
	*MessageCreate `json:"msgcre"`
	*MessageUpdate `json:"msgupd"`
	*MessageDelete `json:"msgdel"`

	*UserList `json:"users"`
}
