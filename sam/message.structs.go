package sam

type (
	// Messages
	Message struct {
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
)

/* Constructors */
func (Message) New() *Message {
	return &Message{}
}

/* Getters/setters */
func (m *Message) GetService() string {
	return m.Service
}

func (m *Message) SetService(value string) *Message {
	if m.Service != value {
		m.changed = append(m.changed, "Service")
		m.Service = value
	}
	return m
}
func (m *Message) GetChannel() string {
	return m.Channel
}

func (m *Message) SetChannel(value string) *Message {
	if m.Channel != value {
		m.changed = append(m.changed, "Channel")
		m.Channel = value
	}
	return m
}
func (m *Message) GetUserName() string {
	return m.UserName
}

func (m *Message) SetUserName(value string) *Message {
	if m.UserName != value {
		m.changed = append(m.changed, "UserName")
		m.UserName = value
	}
	return m
}
func (m *Message) GetUserID() uint64 {
	return m.UserID
}

func (m *Message) SetUserID(value uint64) *Message {
	if m.UserID != value {
		m.changed = append(m.changed, "UserID")
		m.UserID = value
	}
	return m
}
func (m *Message) GetUser() *User {
	return m.User
}

func (m *Message) SetUser(value *User) *Message {
	m.User = value
	return m
}
func (m *Message) GetUserAvatar() string {
	return m.UserAvatar
}

func (m *Message) SetUserAvatar(value string) *Message {
	if m.UserAvatar != value {
		m.changed = append(m.changed, "UserAvatar")
		m.UserAvatar = value
	}
	return m
}
func (m *Message) GetMessage() string {
	return m.Message
}

func (m *Message) SetMessage(value string) *Message {
	if m.Message != value {
		m.changed = append(m.changed, "Message")
		m.Message = value
	}
	return m
}
func (m *Message) GetMessageID() string {
	return m.MessageID
}

func (m *Message) SetMessageID(value string) *Message {
	if m.MessageID != value {
		m.changed = append(m.changed, "MessageID")
		m.MessageID = value
	}
	return m
}
func (m *Message) GetType() MessageType {
	return m.Type
}

func (m *Message) SetType(value MessageType) *Message {
	if m.Type != value {
		m.changed = append(m.changed, "Type")
		m.Type = value
	}
	return m
}
