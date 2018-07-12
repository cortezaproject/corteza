package sam

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

type (
	// Messages
	Message struct {
		Service    string      `db:"service"`
		Channel    string      `db:"channel"`
		UserName   string      `db:"user_name"`
		UserID     uint64      `db:"user_id"`
		User       *User       `db:"user"`
		UserAvatar string      `db:"user_avatar"`
		Message    string      `db:"message"`
		MessageID  string      `db:"message_id"`
		Type       MessageType `db:"type"`

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
