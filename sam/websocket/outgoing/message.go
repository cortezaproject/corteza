package outgoing

import (
	"fmt"
	"time"
)

type Payload struct {
	*Error   `json:"error,omitempty"`
	*Message `json:"m"`

	// @todo: implement outgoing message types
	timestamp time.Time
}

func (p *Payload) Load(payload PayloadType) *Payload {
	switch val := payload.(type) {
	case *Error:
		p.Error = val
	case *Message:
		p.Message = val
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
