package outgoing

import (
	"encoding/json"
	"time"
)

type (
	ChannelMember struct {
		// Channel to part (nil) for ALL channels
		UserID    uint64     `json:"userID,string"`
		Type      string     `json:"type"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	}

	ChannelMemberSet []*ChannelMember
)

func (p *ChannelMember) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{ChannelMember: p})
}

func (p *ChannelMemberSet) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{ChannelMemberSet: p})
}
