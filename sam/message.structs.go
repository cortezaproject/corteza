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
}

func (Message) new() *Message {
	return &Message{}
}

func (m *Message) GetService() string {
	return m.Service
}

func (m *Message) SetService(value string) *Message {
	m.Service = value
	return m
}
func (m *Message) GetChannel() string {
	return m.Channel
}

func (m *Message) SetChannel(value string) *Message {
	m.Channel = value
	return m
}
func (m *Message) GetUserName() string {
	return m.UserName
}

func (m *Message) SetUserName(value string) *Message {
	m.UserName = value
	return m
}
func (m *Message) GetUserID() uint64 {
	return m.UserID
}

func (m *Message) SetUserID(value uint64) *Message {
	m.UserID = value
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
	m.UserAvatar = value
	return m
}
func (m *Message) GetMessage() string {
	return m.Message
}

func (m *Message) SetMessage(value string) *Message {
	m.Message = value
	return m
}
func (m *Message) GetMessageID() string {
	return m.MessageID
}

func (m *Message) SetMessageID(value string) *Message {
	m.MessageID = value
	return m
}
func (m *Message) GetType() MessageType {
	return m.Type
}

func (m *Message) SetType(value MessageType) *Message {
	m.Type = value
	return m
}
