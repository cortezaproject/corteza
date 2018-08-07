package outgoing

import (
	"encoding/json"
)

type (
	User struct {
		// Channel to part (nil) for ALL channels
		ID          string `json:"id"`
		Name        string `json:"name"`
		Username    string `json:"username"`
		Connections uint   `json:"connections"`
	}

	Users []*User
)

func (p *User) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{User: p})
}

func (p *Users) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Users: p})
}
