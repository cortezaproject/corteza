package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

type (
	Message struct {
		ID        uint64       `json:"id"`
		Type      MessageType  `json:"type"`
		Message   string       `json:"message"`
		Meta      *MessageMeta `json:"meta"`
		UserID    uint64       `json:"userId"`
		ChannelID uint64       `json:"channelId"`
		ReplyTo   uint64       `json:"replyTo"`
		Replies   uint         `json:"replies"`
		CreatedAt time.Time    `json:"createdAt,omitempty"`
		UpdatedAt *time.Time   `json:"updatedAt,omitempty"`
		DeletedAt *time.Time   `json:"deletedAt,omitempty"`

		Attachment *Attachment    `json:"attachment,omitempty"`
		Flags      MessageFlagSet `json:"flags,omitempty"`

		Unread *Unread `json:"-"`

		Mentions    MentionSet
		RepliesFrom []uint64
	}

	MessageMeta struct {
		// Bot users can override username/avatar
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}

	MessageFilter struct {
		Query string

		// Required param to filter accessible messages
		CurrentUserID uint64

		// All messages that belong to a channel
		ChannelID []uint64

		// Only messages that belong to a user
		UserID []uint64

		// Replies to a message
		ThreadID []uint64

		// Filter by type
		Type []string

		// (AfterID...BeforeID), for paging
		//
		// Include all messages which IDs range from "first" to "last" (exclusive!)
		AfterID  uint64
		BeforeID uint64

		// [FromID...ToID, for paging
		//
		// Include all messages which IDs range from "from" to "to" (inclusive!)
		FromID uint64
		ToID   uint64

		PinnedOnly      bool
		BookmarkedOnly  bool
		AttachmentsOnly bool

		Limit uint
	}

	MessageType string
)

const (
	MessageTypeSimpleMessage MessageType = ""
	MessageTypeChannelEvent  MessageType = "channelEvent"
	MessageTypeInlineImage   MessageType = "inlineImage"
	MessageTypeAttachment    MessageType = "attachment"
	MessageTypeIlleism       MessageType = "illeism"
)

func (mtype MessageType) String() string {
	return string(mtype)
}

func (mtype MessageType) IsValid() bool {
	switch mtype {
	case MessageTypeSimpleMessage,
		MessageTypeChannelEvent,
		MessageTypeInlineImage,
		MessageTypeAttachment:
		return true
	}

	return false
}

func (mtype MessageType) IsRepliable() bool {
	return mtype.IsEditable()
}

func (mtype MessageType) IsEditable() bool {
	switch mtype {
	case MessageTypeSimpleMessage,
		MessageTypeInlineImage,
		MessageTypeAttachment:
		return true
	}

	return false
}

func (mtype *MessageType) Scan(value interface{}) error {
	switch value.(type) {
	case nil:
		*mtype = MessageTypeSimpleMessage
	case []uint8:
		*mtype = MessageType(string(value.([]uint8)))
		if !mtype.IsValid() {
			return errors.Errorf("Can not scan %v into MessageType", value)
		}
	}

	return nil
}

func (mtype MessageType) Value() (driver.Value, error) {
	if mtype == MessageTypeSimpleMessage {
		return nil, nil
	}

	return mtype.String(), nil
}

func (m *Message) IsValid() bool {
	return m.DeletedAt == nil
}

func (m Message) RBACResource() rbac.Resource {
	return ChannelRBACResource.AppendID(m.ChannelID)
}

func (c Message) DynamicRoles(userID uint64) []uint64 {
	return nil
}

func (meta *MessageMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*meta = MessageMeta{}
		return nil
	case []uint8:
		if err := json.Unmarshal(value.([]byte), meta); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into Message.Meta", value)
		}
		return nil
	}
	return errors.Errorf("Message.Meta: unknown type %T, expected []uint8", value)
}

func (meta *MessageMeta) Value() (driver.Value, error) {
	return json.Marshal(meta)
}
