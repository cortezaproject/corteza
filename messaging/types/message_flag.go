package types

import (
	"time"
)

type (
	MessageFlag struct {
		ID        uint64    `json:"id" db:"id"`
		UserID    uint64    `json:"userId" db:"rel_user"`
		MessageID uint64    `json:"messageId" db:"rel_message"`
		ChannelID uint64    `json:"channelId" db:"rel_channel"`
		Flag      string    `json:"flag" db:"flag"`
		CreatedAt time.Time `json:"createdAt,omitempty" db:"created_at"`

		// Internal only
		DeletedAt *time.Time `json:"-" db:"-"`
	}
)

const (
	MessageFlagPinnedToChannel   string = "pin"
	MessageFlagBookmarkedMessage string = "bookmark"
)

func (f *MessageFlag) IsReaction() bool {
	return f.Flag != MessageFlagPinnedToChannel && f.Flag != MessageFlagBookmarkedMessage
}

func (f *MessageFlag) IsPin() bool {
	return f.Flag == MessageFlagPinnedToChannel
}

func (f *MessageFlag) IsBookmark() bool {
	return f.Flag == MessageFlagBookmarkedMessage
}
