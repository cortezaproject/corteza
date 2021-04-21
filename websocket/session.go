package websocket

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"

	gWebsocket "github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

var sessions = make(map[uint64]*session)

type (
	session struct {
		id   uint64
		once sync.Once
		conn *gWebsocket.Conn

		ctx       context.Context
		ctxCancel context.CancelFunc

		logger *zap.Logger

		send chan []byte
		stop chan []byte

		remoteAddr string

		config *Config

		user auth.Identifiable
	}
)

func Session(ctx context.Context, config *Config, conn *gWebsocket.Conn) *session {
	s := &session{
		conn:   conn,
		config: config,
		send:   make(chan []byte, 512),
		stop:   make(chan []byte, 1),
	}

	s.ctx, s.ctxCancel = context.WithCancel(ctx)

	s.logger = logger.AddRequestID(s.ctx, logger.Default().Named("websocket"))

	return s
}

func (s *session) log(fields ...zapcore.Field) *zap.Logger {
	return s.logger.With(fields...)
}

func (s *session) Context() context.Context {
	return s.ctx
}

func (s *session) User() auth.Identifiable {
	return s.user
}

func (s *session) connected() (err error) {
	// Tell everyone that user has connected
	if err = s.sendPresence("connected"); err != nil {
		return
	}

	// Create a heartbeat every minute for this user
	go func() {
		defer sentry.Recover()

		t := time.NewTicker(time.Second * 60)
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-t.C:
				_ = s.sendPresence("")
			}
		}
	}()

	return nil
}

func (s *session) disconnected() {
	// Tell everyone that user has disconnected
	_ = s.sendPresence("disconnected")

	// Cancel context
	s.ctxCancel()

	// Close connection
	s.conn.Close()
	s.conn = nil
}

// sendPresence sends user presence: "connected", "disconnected" and "" activity kinds
func (s *session) sendPresence(kind string) error {
	//connections := store.CountConnections(s.user.Identity())
	//if kind == "disconnected" {
	//	connections--
	//}

	return nil
}

func (s *session) Handle() (err error) {
	if err = s.connected(); err != nil {
		s.Close()
		return
	}

	go func() {
		_ = s.readLoop()
	}()
	return s.writeLoop()
}

func (s *session) Close() {
	s.once.Do(func() {
		s.disconnected()
		s.Delete()
	})
}

func (s *session) readLoop() (err error) {
	defer func() {
		s.Close()
	}()

	if err = s.conn.SetReadDeadline(time.Now().Add(s.config.PingTimeout)); err != nil {
		return
	}

	s.conn.SetPongHandler(func(string) error {
		return s.conn.SetReadDeadline(time.Now().Add(s.config.PingTimeout))
	})

	s.remoteAddr = s.conn.RemoteAddr().String()

	for {
		_, raw, err := s.conn.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "s.readLoop")
		}

		if err = s.dispatch(raw); err != nil {
			s.log(zap.Error(err)).Error("could not dispatch")
			//_ = s.send(outgoing.NewError(err))
		}
	}
}

func (s *session) writeLoop() error {
	ticker := time.NewTicker(s.config.PingPeriod)

	defer func() {
		ticker.Stop()
		s.Close() // break readLoop
	}()

	write := func(msg []byte) (err error) {
		if s.conn == nil {
			// Connection closed, nowhere to write
			return
		}

		if err = s.conn.SetWriteDeadline(time.Now().Add(s.config.Timeout)); err != nil {
			return
		}

		if msg != nil && s.conn != nil {
			return s.conn.WriteMessage(gWebsocket.TextMessage, msg)
		}

		return
	}

	ping := func() (err error) {
		if s.conn == nil {
			// Connection closed, nothing to ping
			return
		}

		if err = s.conn.SetWriteDeadline(time.Now().Add(s.config.Timeout)); err != nil {
			return
		}

		if s.conn != nil {
			return s.conn.WriteMessage(gWebsocket.PingMessage, nil)
		}

		return
	}

	for {
		select {
		case msg, ok := <-s.send:
			if !ok {
				// channel closed
				return nil
			}

			if err := write(msg); err != nil {
				return errors.Wrap(err, "writeLoop send")
			}

		case msg := <-s.stop:
			// Shutdown requested, don't care if the message is delivered
			_ = write(msg)
			return nil

		case <-ticker.C:
			if err := ping(); err != nil {
				return errors.Wrap(err, "writeLoop ping")
			}
		}
	}
}

func (s *session) dispatch(raw []byte) error {
	var p, err = Unmarshal(raw)
	if err != nil {
		return errors.Wrap(err, "Session.incoming: payload malformed")
	}

	ctx := s.Context()

	if p.Auth != nil {
		return s.authenticate(ctx, p.Auth)
	}

	return nil
}

func (s *session) Save() *session {
	if s.id == 0 {
		s.id = id.Next()
	}

	if s.user != nil {
		if _, ok := sessions[s.user.Identity()]; !ok {
			sessions[s.user.Identity()] = s
		}
	}

	return s
}

func (s *session) Get(userID uint64) *session {
	if sess, ok := sessions[userID]; ok {
		return sess
	}

	return nil
}

func (s *session) Delete() error {
	if s.id == 0 {
		return nil
	}

	if s.user != nil {
		delete(sessions, s.user.Identity())
	}

	return nil
}

// Send sends message to user to ones we want to
// if len(userIDs) == 0 -- it sends to everyone
func (s *session) Send(m *message, userIDs ...uint64) error {
	pb, err := m.EncodeMessage()
	if err != nil {
		return err
	}

	sendsToAll := len(userIDs) == 0
	userIDMap := make(map[uint64]uint64)
	for _, userID := range userIDs {
		userIDMap[userID] = userID
	}

	for _, sess := range sessions {
		_, validUser := userIDMap[sess.user.Identity()]
		if sendsToAll || (!sendsToAll && validUser) {
			_ = sess.sendBytes(pb)
		}
	}

	return nil
}

// sendBytes sends byte to channel or timout
func (s *session) sendBytes(p []byte) error {
	select {
	case s.send <- p:
	case <-time.After(2 * time.Millisecond):
		s.logger.Warn("websocket.sendBytes send timeout")
	}
	return nil
}

func GetActiveSessions() map[uint64]*session {
	return sessions
}
