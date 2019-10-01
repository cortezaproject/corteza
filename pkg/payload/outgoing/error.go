package outgoing

import (
	"encoding/json"
)

type (
	Error struct {
		Message string `json:"m"`
	}
)

func (p *Error) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Error: p})
}

func NewError(err error) *Error {
	return &Error{Message: err.Error()}
}
