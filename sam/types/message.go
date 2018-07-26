package types

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `message.go`, `message.util.go` or `message_test.go` to
	implement your API calls, helper functions and tests. The file `message.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"encoding/json"
	"time"
)

type (
	// Messages -
	Message struct {
		ID        uint64     `db:"id"`
		Type      string     `db:"type"`
		Message   string     `db:"message"`
		UserID    uint64     `db:"rel_user"`
		ChannelID uint64     `db:"rel_channel"`
		ReplyTo   uint64     `db:"reply_to"`
		CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

		changed []string
	}

	// Messages -
	Reaction struct {
		ID        uint64    `db:"id"`
		UserID    uint64    `db:"rel_user"`
		MessageID uint64    `db:"rel_message"`
		ChannelID uint64    `db:"rel_channel"`
		Reaction  string    `db:"reaction"`
		CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`

		changed []string
	}

	// Messages -
	Attachment struct {
		ID         uint64          `db:"id"`
		UserID     uint64          `db:"rel_user"`
		MessageID  uint64          `db:"rel_message"`
		ChannelID  uint64          `db:"rel_channel"`
		Attachment json.RawMessage `db:"attachment"`
		Url        string          `db:"url"`
		PreviewUrl string          `db:"preview_url"`
		Size       uint            `db:"size"`
		Mimetype   string          `db:"mimetype"`
		Name       string          `db:"name"`
		CreatedAt  time.Time       `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt  *time.Time      `json:"updated_at,omitempty" db:"updated_at"`
		DeletedAt  *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`

		changed []string
	}
)

// New constructs a new instance of Message
func (Message) New() *Message {
	return &Message{}
}

// New constructs a new instance of Reaction
func (Reaction) New() *Reaction {
	return &Reaction{}
}

// New constructs a new instance of Attachment
func (Attachment) New() *Attachment {
	return &Attachment{}
}

// Get the value of ID
func (m *Message) GetID() uint64 {
	return m.ID
}

// Set the value of ID
func (m *Message) SetID(value uint64) *Message {
	if m.ID != value {
		m.changed = append(m.changed, "ID")
		m.ID = value
	}
	return m
}

// Get the value of Type
func (m *Message) GetType() string {
	return m.Type
}

// Set the value of Type
func (m *Message) SetType(value string) *Message {
	if m.Type != value {
		m.changed = append(m.changed, "Type")
		m.Type = value
	}
	return m
}

// Get the value of Message
func (m *Message) GetMessage() string {
	return m.Message
}

// Set the value of Message
func (m *Message) SetMessage(value string) *Message {
	if m.Message != value {
		m.changed = append(m.changed, "Message")
		m.Message = value
	}
	return m
}

// Get the value of UserID
func (m *Message) GetUserID() uint64 {
	return m.UserID
}

// Set the value of UserID
func (m *Message) SetUserID(value uint64) *Message {
	if m.UserID != value {
		m.changed = append(m.changed, "UserID")
		m.UserID = value
	}
	return m
}

// Get the value of ChannelID
func (m *Message) GetChannelID() uint64 {
	return m.ChannelID
}

// Set the value of ChannelID
func (m *Message) SetChannelID(value uint64) *Message {
	if m.ChannelID != value {
		m.changed = append(m.changed, "ChannelID")
		m.ChannelID = value
	}
	return m
}

// Get the value of ReplyTo
func (m *Message) GetReplyTo() uint64 {
	return m.ReplyTo
}

// Set the value of ReplyTo
func (m *Message) SetReplyTo(value uint64) *Message {
	if m.ReplyTo != value {
		m.changed = append(m.changed, "ReplyTo")
		m.ReplyTo = value
	}
	return m
}

// Get the value of CreatedAt
func (m *Message) GetCreatedAt() time.Time {
	return m.CreatedAt
}

// Set the value of CreatedAt
func (m *Message) SetCreatedAt(value time.Time) *Message {
	m.changed = append(m.changed, "CreatedAt")
	m.CreatedAt = value
	return m
}

// Get the value of UpdatedAt
func (m *Message) GetUpdatedAt() *time.Time {
	return m.UpdatedAt
}

// Set the value of UpdatedAt
func (m *Message) SetUpdatedAt(value *time.Time) *Message {
	m.changed = append(m.changed, "UpdatedAt")
	m.UpdatedAt = value
	return m
}

// Get the value of DeletedAt
func (m *Message) GetDeletedAt() *time.Time {
	return m.DeletedAt
}

// Set the value of DeletedAt
func (m *Message) SetDeletedAt(value *time.Time) *Message {
	m.changed = append(m.changed, "DeletedAt")
	m.DeletedAt = value
	return m
}

// Changes returns the names of changed fields
func (m *Message) Changes() []string {
	return m.changed
}

// Get the value of ID
func (m *Reaction) GetID() uint64 {
	return m.ID
}

// Set the value of ID
func (m *Reaction) SetID(value uint64) *Reaction {
	if m.ID != value {
		m.changed = append(m.changed, "ID")
		m.ID = value
	}
	return m
}

// Get the value of UserID
func (m *Reaction) GetUserID() uint64 {
	return m.UserID
}

// Set the value of UserID
func (m *Reaction) SetUserID(value uint64) *Reaction {
	if m.UserID != value {
		m.changed = append(m.changed, "UserID")
		m.UserID = value
	}
	return m
}

// Get the value of MessageID
func (m *Reaction) GetMessageID() uint64 {
	return m.MessageID
}

// Set the value of MessageID
func (m *Reaction) SetMessageID(value uint64) *Reaction {
	if m.MessageID != value {
		m.changed = append(m.changed, "MessageID")
		m.MessageID = value
	}
	return m
}

// Get the value of ChannelID
func (m *Reaction) GetChannelID() uint64 {
	return m.ChannelID
}

// Set the value of ChannelID
func (m *Reaction) SetChannelID(value uint64) *Reaction {
	if m.ChannelID != value {
		m.changed = append(m.changed, "ChannelID")
		m.ChannelID = value
	}
	return m
}

// Get the value of Reaction
func (m *Reaction) GetReaction() string {
	return m.Reaction
}

// Set the value of Reaction
func (m *Reaction) SetReaction(value string) *Reaction {
	if m.Reaction != value {
		m.changed = append(m.changed, "Reaction")
		m.Reaction = value
	}
	return m
}

// Get the value of CreatedAt
func (m *Reaction) GetCreatedAt() time.Time {
	return m.CreatedAt
}

// Set the value of CreatedAt
func (m *Reaction) SetCreatedAt(value time.Time) *Reaction {
	m.changed = append(m.changed, "CreatedAt")
	m.CreatedAt = value
	return m
}

// Changes returns the names of changed fields
func (m *Reaction) Changes() []string {
	return m.changed
}

// Get the value of ID
func (m *Attachment) GetID() uint64 {
	return m.ID
}

// Set the value of ID
func (m *Attachment) SetID(value uint64) *Attachment {
	if m.ID != value {
		m.changed = append(m.changed, "ID")
		m.ID = value
	}
	return m
}

// Get the value of UserID
func (m *Attachment) GetUserID() uint64 {
	return m.UserID
}

// Set the value of UserID
func (m *Attachment) SetUserID(value uint64) *Attachment {
	if m.UserID != value {
		m.changed = append(m.changed, "UserID")
		m.UserID = value
	}
	return m
}

// Get the value of MessageID
func (m *Attachment) GetMessageID() uint64 {
	return m.MessageID
}

// Set the value of MessageID
func (m *Attachment) SetMessageID(value uint64) *Attachment {
	if m.MessageID != value {
		m.changed = append(m.changed, "MessageID")
		m.MessageID = value
	}
	return m
}

// Get the value of ChannelID
func (m *Attachment) GetChannelID() uint64 {
	return m.ChannelID
}

// Set the value of ChannelID
func (m *Attachment) SetChannelID(value uint64) *Attachment {
	if m.ChannelID != value {
		m.changed = append(m.changed, "ChannelID")
		m.ChannelID = value
	}
	return m
}

// Get the value of Attachment
func (m *Attachment) GetAttachment() json.RawMessage {
	return m.Attachment
}

// Set the value of Attachment
func (m *Attachment) SetAttachment(value json.RawMessage) *Attachment {
	m.changed = append(m.changed, "Attachment")
	m.Attachment = value
	return m
}

// Get the value of Url
func (m *Attachment) GetUrl() string {
	return m.Url
}

// Set the value of Url
func (m *Attachment) SetUrl(value string) *Attachment {
	if m.Url != value {
		m.changed = append(m.changed, "Url")
		m.Url = value
	}
	return m
}

// Get the value of PreviewUrl
func (m *Attachment) GetPreviewUrl() string {
	return m.PreviewUrl
}

// Set the value of PreviewUrl
func (m *Attachment) SetPreviewUrl(value string) *Attachment {
	if m.PreviewUrl != value {
		m.changed = append(m.changed, "PreviewUrl")
		m.PreviewUrl = value
	}
	return m
}

// Get the value of Size
func (m *Attachment) GetSize() uint {
	return m.Size
}

// Set the value of Size
func (m *Attachment) SetSize(value uint) *Attachment {
	if m.Size != value {
		m.changed = append(m.changed, "Size")
		m.Size = value
	}
	return m
}

// Get the value of Mimetype
func (m *Attachment) GetMimetype() string {
	return m.Mimetype
}

// Set the value of Mimetype
func (m *Attachment) SetMimetype(value string) *Attachment {
	if m.Mimetype != value {
		m.changed = append(m.changed, "Mimetype")
		m.Mimetype = value
	}
	return m
}

// Get the value of Name
func (m *Attachment) GetName() string {
	return m.Name
}

// Set the value of Name
func (m *Attachment) SetName(value string) *Attachment {
	if m.Name != value {
		m.changed = append(m.changed, "Name")
		m.Name = value
	}
	return m
}

// Get the value of CreatedAt
func (m *Attachment) GetCreatedAt() time.Time {
	return m.CreatedAt
}

// Set the value of CreatedAt
func (m *Attachment) SetCreatedAt(value time.Time) *Attachment {
	m.changed = append(m.changed, "CreatedAt")
	m.CreatedAt = value
	return m
}

// Get the value of UpdatedAt
func (m *Attachment) GetUpdatedAt() *time.Time {
	return m.UpdatedAt
}

// Set the value of UpdatedAt
func (m *Attachment) SetUpdatedAt(value *time.Time) *Attachment {
	m.changed = append(m.changed, "UpdatedAt")
	m.UpdatedAt = value
	return m
}

// Get the value of DeletedAt
func (m *Attachment) GetDeletedAt() *time.Time {
	return m.DeletedAt
}

// Set the value of DeletedAt
func (m *Attachment) SetDeletedAt(value *time.Time) *Attachment {
	m.changed = append(m.changed, "DeletedAt")
	m.DeletedAt = value
	return m
}

// Changes returns the names of changed fields
func (m *Attachment) Changes() []string {
	return m.changed
}
