package websocket

import (
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/options"
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
		config options.WebsocketOpt
		logger *zap.Logger
	}
)

func Websocket(logger *zap.Logger, config options.WebsocketOpt) *websocket {
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

	session := Session(ctx, ws.logger, ws.config, conn)

	if err := session.Handle(); err != nil {
		ws.logger.
			WithOptions(zap.AddStacktrace(zap.PanicLevel)).
			Warn("websocket session handler error", zap.Error(err))
	}
}

// Send delivers message to user to ones we want to
// if len(userIDs) == 0 -- it delivers to everyone
func (ws *websocket) Send(kind string, payload interface{}, userIDs ...uint64) error {
	pb, err := Response(kind, payload).Marshal()
	if err != nil {
		return err
	}

	sendsToAll := len(userIDs) == 0
	userIDMap := make(map[uint64]bool)
	for _, userID := range userIDs {
		userIDMap[userID] = true
	}

	for uid, uSessions := range sessions {
		if sendsToAll || (!sendsToAll && userIDMap[uid]) {
			for _, sess := range uSessions {
				_ = sess.sendBytes(pb)
			}
		}
	}

	return nil
}
