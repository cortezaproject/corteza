package sam

import (
	"net/http"
)

// RequestFiller is an interface for typed request parameters
type RequestFiller interface {
	Fill(r *http.Request) error
}

// MessageType is a type used to determine the CRUD state of a message.
type MessageType string

const (
	// MessageTypeCreate is the message type for message creation.
	MessageTypeCreate MessageType = "create"
	// MessageTypeUpdate is the message type for message updates.
	MessageTypeUpdate = "update"
	// MessageTypeDelete is the message type for message deletion.
	MessageTypeDelete = "delete"
)
