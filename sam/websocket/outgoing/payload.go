package outgoing

type (
	Payload struct {
		*Error `json:"error,omitempty"`

		*Connected    `json:"clientConnected,omitempty"`
		*Disconnected `json:"clientDisconnected,omitempty"`

		*Message       `json:"message,omitempty"`
		*MessageDelete `json:"messageDeleted,omitempty"`
		*MessageUpdate `json:"messageUpdated,omitempty"`
		*MessageSet    `json:"messages,omitempty"`

		*ChannelJoin    `json:"channelJoin,omitempty"`
		*ChannelPart    `json:"channelPart,omitempty"`
		*ChannelDeleted `json:"channelDeleted,omitempty"`
		*Channel        `json:"channel,omitempty"`
		*ChannelSet     `json:"channels,omitempty"`

		*User    `json:"user,omitempty"`
		*UserSet `json:"users,omitempty"`
	}

	// This is same-same but different as using the json.Marshaler
	// (this one does not cause json.Marshal to call itself)
	MessageEncoder interface {
		EncodeMessage() ([]byte, error)
	}
)
