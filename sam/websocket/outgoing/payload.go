package outgoing

type (
	Payload struct {
		*Error `json:"error,omitempty"`

		*Connected    `json:"clientConnected,omitempty"`
		*Disconnected `json:"clientDisconnected,omitempty"`

		*Message       `json:"message,omitempty"`
		*MessageDelete `json:"messageDeleted,omitempty"`
		*MessageUpdate `json:"messageUpdated,omitempty"`
		*Messages      `json:"messages,omitempty"`

		*ChannelJoin `json:"channelJoin,omitempty"`
		*ChannelPart `json:"channelPart,omitempty"`
		*Channel     `json:"channel,omitempty"`
		*Channels    `json:"channels,omitempty"`

		*User  `json:"user,omitempty"`
		*Users `json:"users,omitempty"`
	}

	// This is same-same but different as using the json.Marshaler
	// (this one does not cause json.Marshal to call itself)
	MessageEncoder interface {
		EncodeMessage() ([]byte, error)
	}
)
