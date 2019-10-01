package outgoing

import (
	"encoding/json"
)

type (
	Command struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	CommandSet []*Command
)

func (p *CommandSet) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{CommandSet: p})
}
