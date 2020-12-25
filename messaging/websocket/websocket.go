package websocket

import (
	"github.com/cortezaproject/corteza-server/pkg/api"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/logger"
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

// Handles websocket requests from peers
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// Allow connections from any Origin
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (ws Websocket) Open(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Disallow all unauthorized!
	identity := auth.GetIdentityFromContext(ctx)
	if !identity.Valid() {
		api.Send(w, r, errors.New("Unauthorized"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		api.Send(w, r, errors.Wrap(err, "ws: need a websocket handshake"))
		return
	} else if err != nil {
		api.Send(w, r, errors.Wrap(err, "ws: failed to upgrade connection"))
		return
	}

	session := store.Save((&Session{}).New(ctx, ws.config, conn))
	session.user = identity

	if err := session.Handle(); err != nil {
		logger.Default().
			WithOptions(zap.AddStacktrace(zap.PanicLevel)).
			Warn("websocket session handler error", zap.Error(err))
	}

}
