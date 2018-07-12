package sam

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `websocket.go`, `websocket.util.go` or `websocket_test.go` to
	implement your API calls, helper functions and tests. The file `websocket.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

type (
	// Websocket
	Websocket struct {
		UserID uint64 `db:"user_id"`
		User   User   `db:"user"`

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
