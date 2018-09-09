package types

import (
	"time"
)

type (
	Message struct {
		ID         uint64      `json:"id" db:"id"`
		Type       string      `json:"type" db:"type"`
		Message    string      `json:"message" db:"message"`
		UserID     uint64      `json:"userId" db:"rel_user"`
		ChannelID  uint64      `json:"channelId" db:"rel_channel"`
		ReplyTo    uint64      `json:"replyTo" db:"reply_to"`
		CreatedAt  time.Time   `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt  *time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
		DeletedAt  *time.Time  `json:"deletedAt,omitempty" db:"deleted_at"`
		Attachment *Attachment `json:"attachment,omitempty"`
	}

	MessageSet []*Message

	MessageFilter struct {
		Query          string
		ChannelID      uint64
		FromMessageID  uint64
		UntilMessageID uint64
		Limit          uint
	}
)

func (mm MessageSet) Walk(w func(*Message) error) (err error) {
	for i := range mm {
		if err = w(mm[i]); err != nil {
			return
		}
	}

	return
}

func (mm MessageSet) FindById(ID uint64) *Message {
	for i := range mm {
		if mm[i].ID == ID {
			return mm[i]
		}
	}

	return nil
}
