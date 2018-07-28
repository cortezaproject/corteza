package outgoing

import (
	"encoding/json"
)

type (
	Message struct {
		ID        string `json:"id"`
		ChannelID string `json:"cid"`
		Message   string `json:"m"`
		Type      string `json:"t"`
		ReplyTo   string `json:"rid"`
		UserID    string `json:"uid"`
	}

	Messages []*Message

	MessageUpdate struct {
		ID      string `json:"id"`
		Message string `json:"m"`
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
