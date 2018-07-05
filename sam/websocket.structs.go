package sam

// Websocket
type Websocket struct {
	UserID uint64
	User   User

	changed []string
}

func (Websocket) new() *Websocket {
	return &Websocket{}
}

func (w *Websocket) GetUserID() uint64 {
	return w.UserID
}

func (w *Websocket) SetUserID(value uint64) *Websocket {
	if w.UserID != value {
		w.changed = append(w.changed, "userid")
		w.UserID = value
	}
	return w
}
func (w *Websocket) GetUser() User {
	return w.User
}

func (w *Websocket) SetUser(value User) *Websocket {
	if w.User != value {
		w.changed = append(w.changed, "user")
		w.User = value
	}
	return w
}
