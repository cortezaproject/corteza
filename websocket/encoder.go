package websocket

import (
	"encoding/json"
)

const (
	StatusOK    = "ok"
	StatusError = "error"
)

type (
	message struct {
		Status      string      `json:"status"`
		Application uint64      `json:"application"`
		Data        interface{} `json:"data"`
	}

	MessageEncoder interface {
		EncodeMessage() ([]byte, error)
	}
)

func Message(status string, application uint64, data interface{}) *message {
	return &message{
		Status:      status,
		Application: application,
		Data:        data,
	}
}

func (m message) EncodeMessage() ([]byte, error) {
	return json.Marshal(m)
}
