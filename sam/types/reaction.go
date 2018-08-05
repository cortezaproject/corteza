package types

import (
	"time"
)

type (
	Reaction struct {
		ID        uint64    `json:"id" db:"id"`
		UserID    uint64    `json:"userId" db:"rel_user"`
		MessageID uint64    `json:"messageId" db:"rel_message"`
		ChannelID uint64    `json:"channelId" db:"rel_channel"`
		Reaction  string    `json:"reaction" db:"reaction"`
		CreatedAt time.Time `json:"createdAt,omitempty" db:"created_at"`
	}
)
