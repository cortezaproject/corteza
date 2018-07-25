package types

type (
	MessageFilter struct {
		Query          string
		ChannelId      uint64
		FromMessageId  uint64
		UntilMessageId uint64
		Limit          uint
	}
)
