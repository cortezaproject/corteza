package websocket

import (
	"github.com/cortezaproject/corteza-server/pkg/api"
	gWebsocket "github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

var (
	// upgrader handles websocket requests from peers
	upgrader = gWebsocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,

		// Allow connections from any Origin
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

type (
	websocket struct {
		config *Config
		logger *zap.Logger
	}
)

func Websocket(config *Config, logger *zap.Logger) *websocket {
	if !config.LogEnabled {
		logger = zap.NewNop()
	}

	return &websocket{
		config: config,
		logger: logger,
	}
}

func (ws *websocket) Open(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn, err := upgrader.Upgrade(w, r, nil)
	if _, ok := err.(gWebsocket.HandshakeError); ok {
		ws.logger.Error("ws: need a websocket handshake")
		api.Send(w, r, errors.Wrap(err, "ws: need a websocket handshake"))
		return
	} else if err != nil {
		ws.logger.Error("ws: failed to upgrade connection")
		api.Send(w, r, errors.Wrap(err, "ws: failed to upgrade connection"))
		return
	}

	session := Session(ctx, ws.config, conn)

	if err := session.Handle(); err != nil {
		ws.logger.
			WithOptions(zap.AddStacktrace(zap.PanicLevel)).
			Warn("websocket session handler error", zap.Error(err))
	}
}
