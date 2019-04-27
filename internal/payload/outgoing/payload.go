package outgoing

type (
	Payload struct {
		*Error `json:"error,omitempty"`

		*Message    `json:"message,omitempty"`
		*MessageSet `json:"messages,omitempty"`

		*Activity `json:"activity,omitempty"`

		*MessageReaction        `json:"messageReaction,omitempty"`
		*MessageReactionRemoved `json:"messageReactionRemoved,omitempty"`
		*MessagePin             `json:"messagePin,omitempty"`
		*MessagePinRemoved      `json:"messagePinRemoved,omitempty"`

		*ChannelJoin `json:"channelJoin,omitempty"`
		*ChannelPart `json:"channelPart,omitempty"`
		*Channel     `json:"channel,omitempty"`
		*ChannelSet  `json:"channels,omitempty"`

		*ChannelMember    `json:"channelMember,omitempty"`
		*ChannelMemberSet `json:"channelMembers,omitempty"`

		*CommandSet `json:"commands,omitempty"`
	}

	// This is same-same but different as using the json.Marshaler
	// (this one does not cause json.Marshal to call itself)
	MessageEncoder interface {
		EncodeMessage() ([]byte, error)
	}
)
