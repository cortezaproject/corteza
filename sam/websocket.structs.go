package sam

type (
	// Websocket
	Websocket struct {
		UserID uint64
		User   User

		changed []string
	}
)

/* Constructors */
func (Websocket) New() *Websocket {
	return &Websocket{}
}

/* Getters/setters */
func (w *Websocket) GetUserID() uint64 {
	return w.UserID
}

func (w *Websocket) SetUserID(value uint64) *Websocket {
	if w.UserID != value {
		w.changed = append(w.changed, "UserID")
		w.UserID = value
	}
	return w
}
func (w *Websocket) GetUser() User {
	return w.User
}

func (w *Websocket) SetUser(value User) *Websocket {
	w.User = value
	return w
}
