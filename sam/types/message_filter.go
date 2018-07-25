package types

type (
	MessageFilter struct {
		Query          string
		ChannelID      uint64
		FromMessageID  uint64
		UntilMessageID uint64
		Limit          uint
	}
)
