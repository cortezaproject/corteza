package types

import (
	"time"
)

type (
	Mention struct {
		ID            uint64    `db:"id"`
		MessageID     uint64    `db:"rel_message"`
		ChannelID     uint64    `db:"rel_channel"`
		UserID        uint64    `db:"rel_user"`
		MentionedByID uint64    `db:"rel_mentioned_by"`
		CreatedAt     time.Time `db:"created_at"`
	}

	MentionFilter struct {
		// All mentions by this user
		MentionedByID uint64

		// All mentions of this user
		UserID uint64

		// How many entries
		Limit uint
	}
)
