package types

import (
	"time"
)

type (
	Message struct {
		ID        uint64     `db:"id"`
		Type      string     `db:"type"`
		Message   string     `db:"message"`
		UserID    uint64     `db:"rel_user"`
		ChannelID uint64     `db:"rel_channel"`
		ReplyTo   uint64     `db:"reply_to"`
		CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	}

	MessageFilter struct {
		Query          string
		ChannelID      uint64
		FromMessageID  uint64
		UntilMessageID uint64
		Limit          uint
	}
)
