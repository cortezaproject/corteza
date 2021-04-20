package websocket

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

var (
	// Handles websocket requests from peers
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,

		// Allow connections from any Origin
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

type (
	Websocket struct {
		config *Config
	}
)

func New(config *Config) *Websocket {
	ws := &Websocket{
		config: config,
	}

	return ws
}

func (ws Websocket) Open(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn, err := upgrader.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		fmt.Println("ws: need a websocket handshake")
		api.Send(w, r, errors.Wrap(err, "ws: need a websocket handshake"))
		return
	} else if err != nil {
		fmt.Println("ws: failed to upgrade connection")
		api.Send(w, r, errors.Wrap(err, "ws: failed to upgrade connection"))
		return
	}

	session := (&Session{}).New(ctx, ws.config, conn)

	if err := session.Handle(); err != nil {
		logger.Default().
			WithOptions(zap.AddStacktrace(zap.PanicLevel)).
			Warn("websocket session handler error", zap.Error(err))
	}
}
