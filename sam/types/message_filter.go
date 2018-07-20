package types

type (
	MessageFilter struct {
		Query         string
		ChannelId     uint64
		LastMessageId uint64
	}
)
