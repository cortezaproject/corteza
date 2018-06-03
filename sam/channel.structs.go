package sam

// Channels
type Channel struct {
	ID    uint64
	Name  string
	Topic string
}

func (Channel) new() *Channel {
	return &Channel{}
}

func (c *Channel) GetID() uint64 {
	return c.ID
}

func (c *Channel) SetID(value uint64) *Channel {
	c.ID = value
	return c
}
func (c *Channel) GetName() string {
	return c.Name
}

func (c *Channel) SetName(value string) *Channel {
	c.Name = value
	return c
}
func (c *Channel) GetTopic() string {
	return c.Topic
}

func (c *Channel) SetTopic(value string) *Channel {
	c.Topic = value
	return c
}
