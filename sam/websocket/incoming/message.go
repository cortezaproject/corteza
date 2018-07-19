package incoming

import (
	"time"
)

type Message struct {
	// User login
	Login *Login `json:"login"`

	// Channel actions
	Join  *Join  `json:"join"`
	Leave *Leave `json:"leave"`

	// Get channel message history
	History *History `json:"history"`

	// Message actions
	Create *Create `json:"create"`
	Edit   *Edit   `json:"edit"`
	Delete *Delete `json:"delete"`

	// Client notifications (message received, message read, typing indicator)
	Note *Note `json:"note"`

	timestamp time.Time
}

func (Message) New() *Message {
	return &Message{
		timestamp: time.Now().UTC(),
	}
}