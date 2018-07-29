package types

import (
	"time"
)

type (
	Reaction struct {
		ID        uint64    `db:"id"`
		UserID    uint64    `db:"rel_user"`
		MessageID uint64    `db:"rel_message"`
		ChannelID uint64    `db:"rel_channel"`
		Reaction  string    `db:"reaction"`
		CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	}
)
