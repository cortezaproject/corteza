package outgoing

import (
	"encoding/json"
)

type (
	ChannelJoin struct {
		// ID of the channel user is joining
		ID string `json:"id"`

		// ID of the user that is joining
		UserID string `json:"uid"`
	}

	ChannelPart struct {
		// Channel to part (nil) for ALL channels
		ID string `json:"id"`

		// Who is parting
		UserID string `json:"uid"`
	}

	Channel struct {
		// Channel to part (nil) for ALL channels
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	Channels []*Channel
)

func (p *ChannelJoin) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{ChannelJoin: p})
}

func (p *ChannelPart) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{ChannelPart: p})
}

func (p *Channel) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Channel: p})
}

func (p *Channels) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Channels: p})
}
