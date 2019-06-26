package types

type (
	Unread struct {
		ChannelID     uint64 `db:"rel_channel"`
		ReplyTo       uint64 `db:"rel_reply_to"`
		UserID        uint64 `db:"rel_user"`
		LastMessageID uint64 `db:"rel_last_message"`

		Count         uint32 `db:"count"`
		InThreadCount uint32 `db:"-"`
	}

	UnreadFilter struct {
		UserID    uint64
		ChannelID uint64
		ThreadIDs []uint64
	}
)
