package outgoing

type (
	Payload struct {
		*Error `json:"error,omitempty"`

		*Connected    `json:"conn,omitempty"`
		*Disconnected `json:"disc,omitempty"`

		*Message       `json:"m,omitempty"`
		*MessageDelete `json:"md,omitempty"`
		*MessageUpdate `json:"mu,omitempty"`
		*Messages      `json:"ms,omitempty"`

		*ChannelJoin `json:"chj,omitempty"`
		*ChannelPart `json:"chp,omitempty"`
		*Channel     `json:"ch,omitempty"`
		*Channels    `json:"chs,omitempty"`
	}

	// This is same-same but different as using the json.Marshaler
	// (this one does not cause json.Marshal to call itself)
	MessageEncoder interface {
		EncodeMessage() ([]byte, error)
	}
)
