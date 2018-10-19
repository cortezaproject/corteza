package outgoing

import (
	"encoding/json"
	"time"
)

type (
	Message struct {
		ID        uint64 `json:"ID,string"`
		Type      string `json:"type"`
		Message   string `json:"message"`
		ChannelID string `json:"channelID"`
		ReplyTo   uint64 `json:"replyTo,omitempty,string"`
		Replies   uint   `json:"replies,omitempty"`

		User       *User       `json:"user"`
		Attachment *Attachment `json:"att,omitempty"`

		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	MessageSet []*Message
)

func (p *Message) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Message: p})
}

func (p *MessageSet) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{MessageSet: p})
}
