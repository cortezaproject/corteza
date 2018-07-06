package sam

type (
	// Websocket
	WebSocket struct {
		UserID uint64
		User   User

		changed []string
	}
)

/* Constructors */
func (WebSocket) New() *WebSocket {
	return &WebSocket{}
}

/* Getters/setters */
func (w *WebSocket) GetUserID() uint64 {
	return w.UserID
}

func (w *WebSocket) SetUserID(value uint64) *WebSocket {
	if w.UserID != value {
		w.changed = append(w.changed, "UserID")
		w.UserID = value
	}
	return w
}
func (w *WebSocket) GetUser() User {
	return w.User
}

func (w *WebSocket) SetUser(value User) *WebSocket {
	w.User = value
	return w
}
