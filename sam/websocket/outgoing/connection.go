package outgoing

import (
	"encoding/json"
)

type (
	Connected struct {
		// Who did connect?
		UserID string `json:"uid"`
	}

	Disconnected struct {
		// Who did disconnect?
		UserID string `json:"uid"`
	}
)

func (p *Connected) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Connected: p})
}

func (p *Disconnected) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Disconnected: p})
}
