package types

import (
	"time"
)

type (
	Message struct {
		ID        uint64     `json:"id" db:"id"`
		Type      string     `json:"type" db:"type"`
		Message   string     `json:"message" db:"message"`
		UserID    uint64     `json:"userId" db:"rel_user"`
		ChannelID uint64     `json:"channelId" db:"rel_channel"`
		ReplyTo   uint64     `json:"replyTo" db:"reply_to"`
		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	MessageFilter struct {
		Query          string
		ChannelID      uint64
		FromMessageID  uint64
		UntilMessageID uint64
		Limit          uint
	}
)
