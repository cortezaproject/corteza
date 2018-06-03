package sam

// Websocket
type Websocket struct {
	UserID uint64
	User   User
}

func (Websocket) new() *Websocket {
	return &Websocket{}
}

func (w *Websocket) GetUserID() uint64 {
	return w.UserID
}

func (w *Websocket) SetUserID(value uint64) *Websocket {
	w.UserID = value
	return w
}
func (w *Websocket) GetUser() User {
	return w.User
}

func (w *Websocket) SetUser(value User) *Websocket {
	w.User = value
	return w
}
