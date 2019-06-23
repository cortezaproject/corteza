package outgoing

import (
	"encoding/json"
	"time"
)

type (
	Message struct {
		ID        uint64 `json:"messageID,string"`
		Type      string `json:"type"`
		Message   string `json:"message"`
		ChannelID uint64 `json:"channelID,string"`
		UserID    uint64 `json:"userID,string"`

		ReplyTo     uint64   `json:"replyTo,omitempty,string"`
		Replies     uint     `json:"replies,omitempty"`
		RepliesFrom []string `json:"repliesFrom,omitempty"`
		Unread      *Unread  `json:"unread,omitempty"`

		Attachment   *Attachment           `json:"att,omitempty"`
		Mentions     MessageMentionSet     `json:"mentions,omitempty"`
		Reactions    MessageReactionSumSet `json:"reactions,omitempty"`
		IsBookmarked bool                  `json:"isBookmarked"`
		IsPinned     bool                  `json:"isPinned"`

		CanReply  bool `json:"canReply"`
		CanEdit   bool `json:"canEdit"`
		CanDelete bool `json:"canDelete"`

		CreatedAt time.Time  `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	MessageSet []*Message

	MessageMentionSet []string

	// Used for single reaction event notification
	MessageReactionSum struct {
		UserIDs  []string `json:"userIDs"`
		Reaction string   `json:"reaction"`
		Count    uint     `json:"count"`
	}

	MessageReactionSumSet []*MessageReactionSum

	// Used for single reaction event notification
	MessageReaction struct {
		MessageID uint64 `json:"messageID,string"`
		UserID    uint64 `json:"userID,string"`
		Reaction  string `json:"reaction"`
	}

	MessageReactionRemoved MessageReaction

	// Used for single pin/unpin event notification
	MessagePin struct {
		MessageID uint64 `json:"messageID,string"`
		UserID    uint64 `json:"userID,string"`
	}

	MessagePinRemoved MessagePin
)

func (p *Message) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{Message: p})
}

func (p *MessageSet) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{MessageSet: p})
}

func (p *MessageReaction) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{MessageReaction: p})
}

func (p *MessageReactionRemoved) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{MessageReactionRemoved: p})
}

func (p *MessagePin) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{MessagePin: p})
}

func (p *MessagePinRemoved) EncodeMessage() ([]byte, error) {
	return json.Marshal(Payload{MessagePinRemoved: p})
}
