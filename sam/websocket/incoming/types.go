package incoming

type Login struct {
	Username string `json:"username,omitempty"`
	Password []byte `json:"password"`
}

type Join struct {
	Topic string `json:"topic"`
}

type Leave struct {
	Topic string `json:"topic"`
}

type History struct {
	Topic string `json:"topic"`

	// if 0 = last 50 messages, else where message.id < Since
	Since uint64 `json:"since,omitempty"`

	// @todo: extend API (search,...)
}

type Create struct {
	Topic   string      `json:"topic"`
	Content interface{} `json:"content"`
}

type Edit struct {
	ID      string      `json:"id"`
	Topic   string      `json:"topic"`
	Content interface{} `json:"content"`
}

type Delete struct {
	ID    string `json:"id"`
	Topic string `json:"topic"`
}

type Note struct {
	Topic string `json:"topic"`
	Event string `json:"what"`
}
