package types

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `channel.go`, `channel.util.go` or `channel_test.go` to
	implement your API calls, helper functions and tests. The file `channel.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"encoding/json"
	"time"
)

type (
	// Channels - A channel is a representation of a sequence of messages. It has meta data like channel subject. Channels may be public, private or direct (between two users).
	Channel struct {
		ID            uint64          `db:"id"`
		Name          string          `db:"name"`
		Topic         string          `db:"-"`
		Meta          json.RawMessage `db:"meta"`
		LastMessageID uint64          `json:",omitempty" db:"rel_last_message"`
		CreatedAt     time.Time       `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt     *time.Time      `json:"updated_at,omitempty" db:"updated_at"`
		ArchivedAt    *time.Time      `json:"archived_at,omitempty" db:"archived_at"`
		DeletedAt     *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`

		changed []string
	}
)

// New constructs a new instance of Channel
func (Channel) New() *Channel {
	return &Channel{}
}

// Get the value of ID
func (c *Channel) GetID() uint64 {
	return c.ID
}

// Set the value of ID
func (c *Channel) SetID(value uint64) *Channel {
	if c.ID != value {
		c.changed = append(c.changed, "ID")
		c.ID = value
	}
	return c
}

// Get the value of Name
func (c *Channel) GetName() string {
	return c.Name
}

// Set the value of Name
func (c *Channel) SetName(value string) *Channel {
	if c.Name != value {
		c.changed = append(c.changed, "Name")
		c.Name = value
	}
	return c
}

// Get the value of Topic
func (c *Channel) GetTopic() string {
	return c.Topic
}

// Set the value of Topic
func (c *Channel) SetTopic(value string) *Channel {
	if c.Topic != value {
		c.changed = append(c.changed, "Topic")
		c.Topic = value
	}
	return c
}

// Get the value of Meta
func (c *Channel) GetMeta() json.RawMessage {
	return c.Meta
}

// Set the value of Meta
func (c *Channel) SetMeta(value json.RawMessage) *Channel {
	c.changed = append(c.changed, "Meta")
	c.Meta = value
	return c
}

// Get the value of LastMessageID
func (c *Channel) GetLastMessageID() uint64 {
	return c.LastMessageID
}

// Set the value of LastMessageID
func (c *Channel) SetLastMessageID(value uint64) *Channel {
	if c.LastMessageID != value {
		c.changed = append(c.changed, "LastMessageID")
		c.LastMessageID = value
	}
	return c
}

// Get the value of CreatedAt
func (c *Channel) GetCreatedAt() time.Time {
	return c.CreatedAt
}

// Set the value of CreatedAt
func (c *Channel) SetCreatedAt(value time.Time) *Channel {
	c.changed = append(c.changed, "CreatedAt")
	c.CreatedAt = value
	return c
}

// Get the value of UpdatedAt
func (c *Channel) GetUpdatedAt() *time.Time {
	return c.UpdatedAt
}

// Set the value of UpdatedAt
func (c *Channel) SetUpdatedAt(value *time.Time) *Channel {
	c.changed = append(c.changed, "UpdatedAt")
	c.UpdatedAt = value
	return c
}

// Get the value of ArchivedAt
func (c *Channel) GetArchivedAt() *time.Time {
	return c.ArchivedAt
}

// Set the value of ArchivedAt
func (c *Channel) SetArchivedAt(value *time.Time) *Channel {
	c.changed = append(c.changed, "ArchivedAt")
	c.ArchivedAt = value
	return c
}

// Get the value of DeletedAt
func (c *Channel) GetDeletedAt() *time.Time {
	return c.DeletedAt
}

// Set the value of DeletedAt
func (c *Channel) SetDeletedAt(value *time.Time) *Channel {
	c.changed = append(c.changed, "DeletedAt")
	c.DeletedAt = value
	return c
}

// Changes returns the names of changed fields
func (c *Channel) Changes() []string {
	return c.changed
}
