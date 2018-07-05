package sam

// Messages
type Message struct {
	Service    string
	Channel    string
	UserName   string
	UserID     uint64
	User       *User
	UserAvatar string
	Message    string
	MessageID  string
	Type       MessageType

	changed []string
}

func (Message) new() *Message {
	return &Message{}
}

func (m *Message) GetService() string {
	return m.Service
}

func (m *Message) SetService(value string) *Message {
	if m.Service != value {
		m.changed = append(m.changed, "service")
		m.Service = value
	}
	return m
}
func (m *Message) GetChannel() string {
	return m.Channel
}

func (m *Message) SetChannel(value string) *Message {
	if m.Channel != value {
		m.changed = append(m.changed, "channel")
		m.Channel = value
	}
	return m
}
func (m *Message) GetUserName() string {
	return m.UserName
}

func (m *Message) SetUserName(value string) *Message {
	if m.UserName != value {
		m.changed = append(m.changed, "username")
		m.UserName = value
	}
	return m
}
func (m *Message) GetUserID() uint64 {
	return m.UserID
}

func (m *Message) SetUserID(value uint64) *Message {
	if m.UserID != value {
		m.changed = append(m.changed, "userid")
		m.UserID = value
	}
	return m
}
func (m *Message) GetUser() *User {
	return m.User
}

func (m *Message) SetUser(value *User) *Message {
	if m.User != value {
		m.changed = append(m.changed, "user")
		m.User = value
	}
	return m
}
func (m *Message) GetUserAvatar() string {
	return m.UserAvatar
}

func (m *Message) SetUserAvatar(value string) *Message {
	if m.UserAvatar != value {
		m.changed = append(m.changed, "useravatar")
		m.UserAvatar = value
	}
	return m
}
func (m *Message) GetMessage() string {
	return m.Message
}

func (m *Message) SetMessage(value string) *Message {
	if m.Message != value {
		m.changed = append(m.changed, "message")
		m.Message = value
	}
	return m
}
func (m *Message) GetMessageID() string {
	return m.MessageID
}

func (m *Message) SetMessageID(value string) *Message {
	if m.MessageID != value {
		m.changed = append(m.changed, "messageid")
		m.MessageID = value
	}
	return m
}
func (m *Message) GetType() MessageType {
	return m.Type
}

func (m *Message) SetType(value MessageType) *Message {
	if m.Type != value {
		m.changed = append(m.changed, "type")
		m.Type = value
	}
	return m
}
