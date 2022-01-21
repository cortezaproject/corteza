package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/options"
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
	session struct {
		id   uint64
		once sync.Once
		conn *websocket.Conn

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

func (s *session) connected() (err error) {
	s.logger.Info("connected", zap.String("remoteAddr", s.conn.RemoteAddr().String()))

	//// Tell everyone that user has connected
	//if err = s.sendPresence("connected"); err != nil {
	//	return
	//}
	//
	//
	//// Create a heartbeat every minute for this user
	//go func() {
	//	defer sentry.Recover()
	//
	//	t := time.NewTicker(time.Second * 60)
	//	for {
	//		select {
	//		case <-s.ctx.Done():
	//			return
	//		case <-t.C:
	//			_ = s.sendPresence("active")
	//		}
	//	}
	//}()

	return nil
}

func (s *session) disconnected() {
	// Cancel context
	s.ctxCancel()

	s.logger.Info("disconnected")

	// Close connection
	_ = s.conn.Close()
	s.conn = nil
}

//func (s *session) sendPresence(_ string) error {
//	return nil
//}

func (s *session) Handle() (err error) {
	if err = s.connected(); err != nil {
		s.Close()
		return
	}

	go func() {
		// Close unidentified connections in 5sec
		<-time.NewTimer(time.Second * 5).C
		if s.identity == nil {
			s.Write([]byte(closingUnidentifiedConn))
			s.logger.Info("closing unidentified connection")
			s.Close()
		}
	}()

	go func() {
		if err = s.readLoop(); err != nil {
			s.logger.Error("read failure", zap.Error(err))
		}
		s.Close()
	}()

	if err = s.writeLoop(); err != nil {
		s.logger.Error("write failure", zap.Error(err))
	}

	return
}

func (s *session) Close() {
	s.once.Do(func() {
		s.disconnected()
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
		if s.conn == nil {
			return nil
		}

		if _, raw, err = s.conn.ReadMessage(); err != nil {
			return errHandler("read failed", err)
		}

		if err = s.procRawMessage(raw); err != nil {
			return
		}
	}
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

		s.logger.Debug(
			"authenticated",
			zap.Uint64("userID", s.identity.Identity()),
			zap.Uint64s("roles", s.identity.Roles()),
		)

		s.server.StoreSession(s)

		// not expecting anything else
		return
	}

	if s.identity == nil {
		return fmt.Errorf("unauthenticated session")
	}

	// at the moment we do not support any other kinds of message types
	return fmt.Errorf("unknown message type '%s'", pw.Type)
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
			return fmt.Errorf("deadline error: %w", err)
		}

		if msg != nil && s.conn != nil {
			return s.conn.WriteMessage(websocket.TextMessage, msg)
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
			return s.conn.WriteMessage(websocket.PingMessage, nil)
		}

		return
	}

	for {
		if s.conn == nil {
			return nil
		}

		select {
		case msg, ok := <-s.send:
			if !ok {
				// channel closed
				return nil
			}

			if err := errHandler("send failed", write(msg)); err != nil {
				return err
			}

		case msg := <-s.stop:
			// Shutdown requested, don't care if the message is delivered
			_ = write(msg)
			return nil

		case <-ticker.C:
			if err := ping(); err != nil {
				return errHandler("ping failed", err)
			}
		}
	}
}

func (s *session) authenticate(p *payloadAuth) error {
	identity, err := s.server.tokenValidator(s.ctx, p.AccessToken)
	if err != nil {
		return err
	}

	if s.identity != nil {
		if s.identity.Identity() != identity.Identity() {
			return fmt.Errorf("identity does not match")
		}
	}

	if !identity.Valid() {
		return fmt.Errorf("invalid identity")
	}

	s.identity = identity
	_, _ = s.Write(ok)
	return nil
}

// sendBytes sends byte to channel or timeout
func (s *session) Write(p []byte) (int, error) {
	select {
	case s.send <- p:
		return len(p), nil
	case <-time.After(2 * time.Millisecond):
		return 0, fmt.Errorf("write timedout")
	}
}

func errHandler(wrap string, err error) error {
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
	return fmt.Errorf(wrap+": %w", err)
}
