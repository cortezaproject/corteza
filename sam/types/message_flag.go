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
	}

	MessageFlagSet []*MessageFlag
)

const (
	MessageFlagPinnedToChannel   string = "pin"
	MessageFlagBookmarkedMessage string = "bookmark"
)

func (ff MessageFlagSet) Walk(w func(*MessageFlag) error) (err error) {
	for i := range ff {
		if err = w(ff[i]); err != nil {
			return
		}
	}

	return
}

func (ff MessageFlagSet) FindById(ID uint64) *MessageFlag {
	for i := range ff {
		if ff[i].ID == ID {
			return ff[i]
		}
	}

	return nil
}

func (ff MessageFlagSet) IsBookmarked(UserID uint64) bool {
	for i := range ff {
		if ff[i].UserID == UserID && ff[i].IsBookmark() {
			return true
		}
	}

	return false
}

func (ff MessageFlagSet) IsPinned() bool {
	for i := range ff {
		if ff[i].IsPin() {
			return true
		}
	}

	return false
}

func (f *MessageFlag) IsReaction() bool {
	return f.Flag != MessageFlagPinnedToChannel && f.Flag != MessageFlagBookmarkedMessage
}

func (f *MessageFlag) IsPin() bool {
	return f.Flag == MessageFlagPinnedToChannel
}

func (f *MessageFlag) IsBookmark() bool {
	return f.Flag == MessageFlagBookmarkedMessage
}
