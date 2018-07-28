package outgoing

import (
	"encoding/json"
	"time"
)

type (
	Message struct {
		ID        string `json:"id"`
		Type      string `json:"t"`
		Message   string `json:"m"`
		UserID    string `json:"uid"`
		ChannelID string `json:"cid"`
		ReplyTo   string `json:"rid"`

		CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	}

	Messages []*Message

	MessageUpdate struct {
		ID      string `json:"id"`
		Message string `json:"m"`

		UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
	}

	MessageDelete struct {
		ID string `json:"id"`
	}
)

func (p *Message) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Message: p})
}

func (p *Messages) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Messages: p})
}

func (p *MessageUpdate) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{MessageUpdate: p})
}

func (p *MessageDelete) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{MessageDelete: p})
}
