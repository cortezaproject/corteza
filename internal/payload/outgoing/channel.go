package outgoing

import (
	"encoding/json"
	"time"
)

type (
	ChannelJoin struct {
		// ID of the channel user is joining
		ID string `json:"channelID"`

		// ID of the user that is joining
		UserID string `json:"userID"`
	}

	ChannelPart struct {
		// Channel to part (nil) for ALL channels
		ID string `json:"channelID"`

		// Who is parting
		UserID string `json:"userID"`
	}

	Channel struct {
		// Channel to part (nil) for ALL channels
		ID            string       `json:"ID"`
		Name          string       `json:"name"`
		Topic         string       `json:"topic"`
		Type          string       `json:"type"`
		LastMessageID string       `json:"lastMessageID"`
		Members       []string     `json:"members,omitempty"`
		View          *ChannelView `json:"view,omitempty"`

		CreatedAt  time.Time  `json:"createdAt"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
		ArchivedAt *time.Time `json:"archivedAt,omitempty"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty"`
	}

	ChannelSet []*Channel
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

func (p *ChannelSet) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{ChannelSet: p})
}
