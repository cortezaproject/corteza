package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// active sessions of users
var (
	// wrapper around nextID that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)

type (
	conection interface {
		Close() error
		RemoteAddr() net.Addr
		WriteMessage(messageType int, data []byte) error
		SetWriteDeadline(t time.Time) error
		ReadMessage() (messageType int, p []byte, err error)
		SetReadDeadline(t time.Time) error
		SetPongHandler(h func(appData string) error)
	}

	session struct {
		l sync.RWMutex

		id   uint64
		once sync.Once
		conn conection

		ctx       context.Context
		ctxCancel context.CancelFunc

		logger *zap.Logger

		send chan []byte
		stop chan []byte

		remoteAddr string

		config options.WebsocketOpt

		identity auth.Identifiable

		server *server
	}
)

func Session(ctx context.Context, ws *server, conn *websocket.Conn) *session {
	s := &session{
		id:     nextID(),
		conn:   conn,
		config: ws.config,
		send:   make(chan []byte, 512),
		stop:   make(chan []byte, 1),
		server: ws,
	}

	s.ctx, s.ctxCancel = context.WithCancel(ctx)

	s.logger = ws.logger.
		Named("session").
		With(
			zap.Uint64("id", s.id),
		)

	return s
}

func (s *session) Identity() auth.Identifiable {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.identity
}

func (s *session) disconnect() {
	s.l.Lock()
	defer s.l.Unlock()

	// Cancel context
	s.ctxCancel()

	s.logger.Info("disconnected")

	// Close connection
	_ = s.conn.Close()

	close(s.send)
	close(s.stop)
	s.conn = nil
}

func (s *session) Handle() error {
	go func() {
		// Close unidentified connections in 5sec
		<-time.NewTimer(time.Second * 5).C
		if s.Identity() == nil {
			_, _ = s.Write(closingUnidentifiedConn)
			s.logger.Info("closing unidentified connection")
			s.Close()
		}
	}()

	go func() {
		if err := s.readLoop(); err != nil {
			if errors.Is(err, net.ErrClosed) {
				// read will return net.ErrClosed when
				// recovering from panic
				return
			}

			s.logger.Error("read failure", zap.Error(err))
		}
		s.Close()
	}()

	if err := s.writeLoop(); err != nil {
		if errors.Is(err, net.ErrClosed) {
			// write will return net.ErrClosed when
			// recovering from panic
			return nil
		}

		s.logger.Error("write failure", zap.Error(err))
	}

	return nil
}

func (s *session) Close() {
	s.once.Do(func() {
		s.disconnect()
		s.server.RemoveSession(s)
	})
}

func (s *session) readLoop() (err error) {
	if err = s.conn.SetReadDeadline(time.Now().Add(s.config.PingTimeout)); err != nil {
		return
	}

	s.conn.SetPongHandler(func(string) error {
		return s.conn.SetReadDeadline(time.Now().Add(s.config.PingTimeout))
	})

	s.remoteAddr = s.conn.RemoteAddr().String()

	var (
		raw []byte
	)

	for {
		if raw, err = s.read(); err != nil {
			return
		}

		if raw == nil {
			continue
		}

		if err = s.procRawMessage(raw); err != nil {
			return
		}
	}
}

func (s *session) read() (raw []byte, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			s.logger.Debug("recovering from websocket read panic", zap.Any("recovered-error", recovered))
			err = net.ErrClosed
		}
	}()

	s.l.RLock()
	defer s.l.RUnlock()

	if _, raw, err = s.conn.ReadMessage(); err != nil {
		return nil, errHandler("websocket read failed", err)
	}

	return
}

func (s *session) procRawMessage(raw []byte) (err error) {
	pw := payloadWrap{}
	if err = json.Unmarshal(raw, &pw); err != nil {
		return fmt.Errorf("could not unmarshal session message: %w", err)
	}

	if pw.Type == payloadTypeCredentials {
		authPayload := &payloadAuth{}
		if err = pw.UnmarshalValue(authPayload); err != nil {
			return fmt.Errorf("could not unmarshal session payload: %w", err)
		}

		if err = s.authenticate(authPayload); err != nil {
			return fmt.Errorf("unauthorized: %w", err)
		}

		i := s.Identity()
		s.logger.Debug(
			"authenticated",
			zap.Uint64("userID", i.Identity()),
			zap.Uint64s("roles", i.Roles()),
		)

		s.server.StoreSession(s)

		// not expecting anything else
		return
	}

	if s.Identity() == nil {
		return fmt.Errorf("unauthenticated session")
	}

	// at the moment we do not support any other kinds of message types
	return fmt.Errorf("unknown message type '%s'", pw.Type)
}

// reads send & stop channels and sends received messages to websocket connection via write fn()
func (s *session) writeLoop() error {
	ticker := time.NewTicker(s.config.PingPeriod)

	defer ticker.Stop()
	defer s.Close() // break readLoop

	for {
		select {
		case msg, ok := <-s.send:
			if !ok {
				// channel closed
				return nil
			}

			if err := s.write(websocket.TextMessage, msg); err != nil {
				return err
			}

			// continue with wait & write/ping loop

		case msg, ok := <-s.stop:
			if !ok {
				// channel closed
				return nil
			}

			// Shutdown requested, don't care if the message is delivered
			if err := s.write(websocket.TextMessage, msg); err != nil {
				return err
			}

			// stopping, break the loop.
			return nil

		case <-ticker.C:
			if err := s.write(websocket.PingMessage, nil); err != nil {
				return err
			}

			// continue with wait & write/ping loop
		}
	}
}

// writes messages to websocket connection
func (s *session) write(t int, msg []byte) (err error) {
	s.l.RLock()
	defer s.l.RUnlock()

	defer func() {
		if recovered := recover(); recovered != nil {
			s.logger.Debug("recovering from websocket write panic", zap.Any("recovered-error", recovered))
			err = net.ErrClosed
		}
	}()

	if err = s.conn.SetWriteDeadline(time.Now().Add(s.config.Timeout)); err != nil {
		return fmt.Errorf("deadline error: %w", err)
	}

	return errHandler("websocket write failed", s.conn.WriteMessage(t, msg))
}

func (s *session) authenticate(p *payloadAuth) error {
	identity, err := s.server.tokenValidator(s.ctx, p.AccessToken)
	if err != nil {
		return err
	}

	if i := s.Identity(); i != nil {
		if i.Identity() != identity.Identity() {
			return fmt.Errorf("identity does not match")
		}
	}

	if !identity.Valid() {
		return fmt.Errorf("invalid identity")
	}

	s.l.Lock()
	defer s.l.Unlock()

	s.identity = identity
	_, _ = s.Write(ok)
	return nil
}

// sendBytes sends byte to channel or timeout
func (s *session) Write(p []byte) (int, error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			s.logger.Debug("recovering from websocket write panic", zap.Any("recovered-error", recovered))
		}
	}()

	select {
	case s.send <- p:
		return len(p), nil
	case <-time.After(2 * time.Millisecond):
		return 0, fmt.Errorf("write timedout")
	}
}

func errHandler(prefix string, err error) error {
	if err == nil {
		return nil
	}

	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		// normal closing
		return nil
	}

	if errors.Is(err, net.ErrClosed) {
		// suppress errors when reading/writing from/to a closed connection
		return nil
	}

	return fmt.Errorf(prefix+": %w", err)
}
