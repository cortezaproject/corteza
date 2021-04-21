package websocket

import (
	"encoding/json"
)

const (
	StatusOK    = "ok"
	StatusError = "error"

	WorkflowApplication = "Workflow"
)

type (
	message struct {
		Status      string      `json:"status"`
		Application string      `json:"application"`
		Data        interface{} `json:"data"`
	}

	MessageEncoder interface {
		EncodeMessage() ([]byte, error)
	}
)

func Message(status, application string, data interface{}) *message {
	return &message{
		Status:      status,
		Application: application,
		Data:        data,
	}
}

func (m message) EncodeMessage() ([]byte, error) {
	return json.Marshal(m)
}
