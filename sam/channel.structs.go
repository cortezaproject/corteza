package sam

import (
	"time"
)

// Channels
type Channel struct {
	ID    uint64
	Name  string
	Topic string

	ArchivedAt *time.Time `json:",omitempty"`
	DeletedAt  *time.Time `json:",omitempty"`

	changed []string
}

func (Channel) new() *Channel {
	return &Channel{}
}

func (c *Channel) GetID() uint64 {
	return c.ID
}

func (c *Channel) SetID(value uint64) *Channel {
	if c.ID != value {
		c.changed = append(c.changed, "id")
		c.ID = value
	}
	return c
}
func (c *Channel) GetName() string {
	return c.Name
}

func (c *Channel) SetName(value string) *Channel {
	if c.Name != value {
		c.changed = append(c.changed, "name")
		c.Name = value
	}
	return c
}
func (c *Channel) GetTopic() string {
	return c.Topic
}

func (c *Channel) SetTopic(value string) *Channel {
	if c.Topic != value {
		c.changed = append(c.changed, "topic")
		c.Topic = value
	}
	return c
}
