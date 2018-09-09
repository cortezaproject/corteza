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

		Attachment *Attachment `json:"att,omitempty"`

		CreatedAt time.Time  `json:"cat,omitempty"`
		UpdatedAt *time.Time `json:"uat,omitempty"`
	}

	Messages []*Message

	MessageUpdate struct {
		ID      string `json:"id"`
		Message string `json:"m"`

		UpdatedAt time.Time `json:"uat,omitempty"`
	}

	MessageDelete struct {
		ID string `json:"id"`
	}

	Attachment struct {
		ID         string     `json:"id"`
		UserID     string     `json:"uid"`
		Url        string     `json:"url"`
		PreviewUrl string     `json:"prw"`
		Size       int64      `json:"sze"`
		Mimetype   string     `json:"typ"`
		Name       string     `json:"nme"`
		CreatedAt  time.Time  `json:"cat,omitempty"`
		UpdatedAt  *time.Time `json:"uat,omitempty"`
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
