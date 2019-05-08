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
		ID             string   `json:"channelID"`
		Name           string   `json:"name"`
		Topic          string   `json:"topic"`
		Type           string   `json:"type"`
		LastMessageID  string   `json:"lastMessageID"`
		Unread         *Unread  `json:"unread,omitempty"`
		Members        []string `json:"members,omitempty"`
		MembershipFlag string   `json:"membershipFlag"`

		CanJoin           bool `json:"canJoin"`
		CanPart           bool `json:"canPart"`
		CanObserve        bool `json:"canObserve"`
		CanSendMessages   bool `json:"canSendMessages"`
		CanDeleteMessages bool `json:"canDeleteMessages"`
		CanChangeMembers  bool `json:"canChangeMembers"`
		CanUpdate         bool `json:"canUpdate"`
		CanArchive        bool `json:"canArchive"`
		CanDelete         bool `json:"canDelete"`

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
