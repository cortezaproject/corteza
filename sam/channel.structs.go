package sam

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
	// Channels
	Channel struct {
		ID            uint64
		Name          string
		Topic         string
		Meta          json.RawMessage
		LastMessageId uint64     `json:",omitempty" db:"rel_last_message"`
		ArchivedAt    *time.Time `json:",omitempty" db:"archived_at"`
		DeletedAt     *time.Time `json:",omitempty" db:"deleted_at"`

		changed []string
	}
)

/* Constructors */
func (Channel) New() *Channel {
	return &Channel{}
}

/* Getters/setters */
func (c *Channel) GetID() uint64 {
	return c.ID
}

func (c *Channel) SetID(value uint64) *Channel {
	if c.ID != value {
		c.changed = append(c.changed, "ID")
		c.ID = value
	}
	return c
}
func (c *Channel) GetName() string {
	return c.Name
}

func (c *Channel) SetName(value string) *Channel {
	if c.Name != value {
		c.changed = append(c.changed, "Name")
		c.Name = value
	}
	return c
}
func (c *Channel) GetTopic() string {
	return c.Topic
}

func (c *Channel) SetTopic(value string) *Channel {
	if c.Topic != value {
		c.changed = append(c.changed, "Topic")
		c.Topic = value
	}
	return c
}
func (c *Channel) GetMeta() json.RawMessage {
	return c.Meta
}

func (c *Channel) SetMeta(value json.RawMessage) *Channel {
	c.Meta = value
	return c
}
func (c *Channel) GetLastMessageId() uint64 {
	return c.LastMessageId
}

func (c *Channel) SetLastMessageId(value uint64) *Channel {
	if c.LastMessageId != value {
		c.changed = append(c.changed, "LastMessageId")
		c.LastMessageId = value
	}
	return c
}
func (c *Channel) GetArchivedAt() *time.Time {
	return c.ArchivedAt
}

func (c *Channel) SetArchivedAt(value *time.Time) *Channel {
	c.ArchivedAt = value
	return c
}
func (c *Channel) GetDeletedAt() *time.Time {
	return c.DeletedAt
}

func (c *Channel) SetDeletedAt(value *time.Time) *Channel {
	c.DeletedAt = value
	return c
}
