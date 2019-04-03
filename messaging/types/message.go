package types

import (
	"time"

	"encoding/json"

	"database/sql/driver"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/rules"
	systemTypes "github.com/crusttech/crust/system/types"
)

type (
	Message struct {
		ID        uint64       `json:"id" db:"id"`
		Type      MessageType  `json:"type" db:"type"`
		Message   string       `json:"message" db:"message"`
		Meta      *MessageMeta `json:"meta" db:"meta"`
		UserID    uint64       `json:"userId" db:"rel_user"`
		ChannelID uint64       `json:"channelId" db:"rel_channel"`
		ReplyTo   uint64       `json:"replyTo" db:"reply_to"`
		Replies   uint         `json:"replies" db:"replies"`
		CreatedAt time.Time    `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time   `json:"updatedAt,omitempty" db:"updated_at"`
		DeletedAt *time.Time   `json:"deletedAt,omitempty" db:"deleted_at"`

		Attachment *Attachment       `json:"attachment,omitempty"`
		User       *systemTypes.User `json:"user,omitempty"`
		Flags      MessageFlagSet    `json:"flags,omitempty"`

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
		ChannelID uint64

		// Only messages that belong to a user
		UserID uint64

		// Return all replies to a single message
		RepliesTo uint64

		// (FirstID...LastID), for paging
		//
		// Include all messages which IDs range from "first" to "last" (exclusive!)
		FirstID uint64
		LastID  uint64

		// [FromID...ToID, for paging
		//
		// Include all messages which IDs range from "from" to "to" (inclusive!)
		FromID uint64
		ToID   uint64

		PinnedOnly     bool
		BookmarkedOnly bool

		Limit uint
	}

	MessageType string
)

const (
	MessageTypeSimpleMessage MessageType = ""
	MessageTypeChannelEvent              = "channelEvent"
	MessageTypeInlineImage               = "inlineImage"
	MessageTypeAttachment                = "attachment"
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

//func (mtype *MessageType) Scan(value interface{}) error {
//	switch value.(type) {
//	case nil:
//		*mtype = MessageTypeSimpleMessage
//	case []uint8:
//		*mtype = MessageType(string(value.([]uint8)))
//		if !mtype.IsValid() {
//			return errors.Errorf("Can not scan %v into MessageType", value)
//		}
//	}
//
//	return nil
//}

func (mtype MessageType) Value() (driver.Value, error) {
	if mtype == MessageTypeSimpleMessage {
		return nil, nil
	}

	return mtype.String(), nil
}

func (m *Message) IsValid() bool {
	return m.DeletedAt == nil
}

func (m Message) PermissionResource() rules.Resource {
	return ChannelPermissionResource.AppendID(m.ChannelID)
}

func (meta *MessageMeta) Scan(value interface{}) error {
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
