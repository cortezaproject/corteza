package types

//go:generate go run ../../codegen/v2/type-set.go --no-pk-types Unread --output unread.gen.go

type (
	Unread struct {
		ChannelID     uint64 `db:"rel_channel"`
		ReplyTo       uint64 `db:"rel_reply_to"`
		UserID        uint64 `db:"rel_user"`
		LastMessageID uint64 `db:"rel_last_message"`

		Count uint32 `db:"count"`
	}

	UnreadFilter struct {
		UserID uint64
	}
)
