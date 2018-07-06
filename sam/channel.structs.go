package sam

import (
	"time"
)

type (
	// Channels
	Channel struct {
		ID         uint64
		Name       string
		Topic      string
		ArchivedAt *time.Time `json:",omitempty"`
		DeletedAt  *time.Time `json:",omitempty"`

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
