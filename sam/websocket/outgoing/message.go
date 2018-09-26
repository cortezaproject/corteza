package outgoing

import (
	"encoding/json"
	"time"
)

type (
	Message struct {
		ID        string `json:"ID"`
		Type      string `json:"type"`
		Message   string `json:"message"`
		ChannelID string `json:"channelID"`
		ReplyTo   string `json:"replyID"`

		User       *User       `json:"user"`
		Attachment *Attachment `json:"att,omitempty"`

		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	}

	MessageSet []*Message

	MessageUpdate struct {
		ID      string `json:"ID"`
		Message string `json:"message"`

		UpdatedAt time.Time `json:"updatedAt,omitempty"`
	}

	MessageDelete struct {
		ID string `json:"ID"`
	}
)

func (p *Message) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Message: p})
}

func (p *MessageSet) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{MessageSet: p})
}

func (p *MessageUpdate) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{MessageUpdate: p})
}

func (p *MessageDelete) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{MessageDelete: p})
}
