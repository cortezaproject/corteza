package outgoing

import (
	"fmt"
	"time"
)

type (
	Payload struct {
		*Error         `json:"error,omitempty"`
		*Message       `json:"m,omitempty"`
		*MessageDelete `json:"md,omitempty"`
		*MessageUpdate `json:"mu,omitempty"`
		*Messages      `json:"ms,omitempty"`

		// @todo: implement outgoing message types
		timestamp time.Time
	}
	PayloadType interface {
		valid() bool
	}
)

func (p *Payload) Load(payload PayloadType) *Payload {
	switch val := payload.(type) {
	case *Error:
		p.Error = val
	case *Message:
		p.Message = val
	case *Messages:
		p.Messages = val
	case *MessageDelete:
		p.MessageDelete = val
	case *MessageUpdate:
		p.MessageUpdate = val
	default:
		panic(fmt.Sprintf("Unknown/unsupported Payload type: %T", val))
	}
	return p
}

func (Payload) New() *Payload {
	return &Payload{
		timestamp: time.Now().UTC(),
	}
}

func NewError(err error) *Payload {
	return Payload{}.New().Load(&Error{Message: err.Error()})
}
