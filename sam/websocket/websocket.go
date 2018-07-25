package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"
)

type (
	Websocket struct{}
)

func (Websocket) New() *Websocket {
	return &Websocket{}
}

// Handles websocket requests from peers
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow connections from any Origin
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (Websocket) Open(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		resputil.JSON(w, errors.Wrap(err, "ws: need a websocket handshake"))
		return
	} else if err != nil {
		resputil.JSON(w, errors.Wrap(err, "ws: failed to upgrade connection"))
		return
	}

	session := store.Save(Session{}.New(r.Context(), ws))
	if err := session.Handle(); err != nil {
		// @todo: log error, because at this point we can't really write it to w
	}
}
