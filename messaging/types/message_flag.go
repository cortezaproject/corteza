package types

import (
	"time"
)

type (
	MessageFlag struct {
		ID        uint64    `json:"id"`
		UserID    uint64    `json:"userId"`
		MessageID uint64    `json:"messageId"`
		ChannelID uint64    `json:"channelId"`
		Flag      string    `json:"flag"`
		CreatedAt time.Time `json:"createdAt,omitempty"`

		// Internal only
		DeletedAt *time.Time `json:"-"`
	}

	MessageFlagFilter struct {
		Flag      string
		MessageID []uint64
	}
)

const (
	MessageFlagPinnedToChannel   string = "pin"
	MessageFlagBookmarkedMessage string = "bookmark"
)

func (f MessageFlag) IsReaction() bool {
	return f.Flag != MessageFlagPinnedToChannel && f.Flag != MessageFlagBookmarkedMessage
}

func (f MessageFlag) IsPin() bool {
	return f.Flag == MessageFlagPinnedToChannel
}

func (f MessageFlag) IsBookmark() bool {
	return f.Flag == MessageFlagBookmarkedMessage
}
