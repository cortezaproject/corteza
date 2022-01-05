package websocket

import (
	"io"
	"net/http"
	"sync"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type (
	server struct {
		config options.WebsocketOpt
		logger *zap.Logger

		// user id => session id => session
		sessions map[uint64]map[uint64]io.Writer

		// keep lock on session map changes
		l sync.RWMutex
	}
)

var (
	// upgrader handles websocket requests from peers
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,

		// Allow connections from any Origin
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func Server(logger *zap.Logger, config options.WebsocketOpt) *server {
	if !config.LogEnabled {
		logger = zap.NewNop()
	}

	return &server{
		config:   config,
		logger:   logger.Named("websocket"),
		sessions: make(map[uint64]map[uint64]io.Writer),
	}
}

func (ws *server) Open(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn, err := upgrader.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		errors.ProperlyServeHTTP(w, r, errors.Internal("need a websocket handshake"), false)
		return
	} else if err != nil {
		errors.ProperlyServeHTTP(w, r, errors.Internal("failed to upgrade connection").Wrap(err), false)
		return
	}

	// init new session
	//
	// session will add itself back to server's session map when
	// ready (if user authenticates itself)
	ses := Session(ctx, ws, conn)

	if err = ses.Handle(); err != nil {
		ws.logger.Warn("websocket session handler error", zap.Error(err))
	}
}

// Send delivers payload to one, more or all users
//
// Omit userIDs to deliver to ALL users
func (ws *server) Send(t string, payload interface{}, userIDs ...uint64) error {
	pb, err := MarshalPayload(t, payload)
	if err != nil {
		return err
	}

	var (
		sendToAll = len(userIDs) == 0
		uMap      = slice.ToUint64BoolMap(userIDs)
	)

	ws.l.RLock()
	defer ws.l.RUnlock()

	for uid := range ws.sessions {
		if sendToAll || (!sendToAll && uMap[uid]) {
			for _, s := range ws.sessions[uid] {
				_, err = s.Write(pb)
			}
		}
	}

	return nil
}

func (ws *server) StoreSession(s *session) {
	ws.l.Lock()
	defer ws.l.Unlock()
	if s.identity != nil {
		ws.storeSession(s, s.identity.Identity(), s.id)
	}
}

func (ws *server) storeSession(w io.Writer, uid, sid uint64) {
	if ws.sessions[uid] == nil {
		ws.sessions[uid] = make(map[uint64]io.Writer)

	}

	ws.sessions[uid][sid] = w
}

func (ws *server) RemoveSession(s *session) {
	ws.l.Lock()
	defer ws.l.Unlock()
	if s.identity != nil {
		uid := s.identity.Identity()
		delete(ws.sessions[uid], s.id)

		if len(ws.sessions[uid]) == 0 {
			delete(ws.sessions, uid)
		}
	}
}
