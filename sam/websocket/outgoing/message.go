package outgoing

import (
	"time"

	"github.com/titpetric/factory"
)

type Message struct {
	Error *Error `json:"error,omitempty"`

	// @todo: implement outgoing message types

	id        uint64
	timestamp time.Time
}

func (Message) New() *Message {
	return &Message{
		id:        factory.Sonyflake.NextID(),
		timestamp: time.Now().UTC(),
	}
}

func (m *Message) FromError(err error) *Message {
	m.Error = &Error{err.Error()}
	return m
}
