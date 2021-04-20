package websocket

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/payload/outgoing"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	Session struct {
		id   uint64
		once sync.Once
		Conn *websocket.Conn

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

func (*Session) New(ctx context.Context, config *Config, conn *websocket.Conn) *Session {
	s := &Session{
		Conn:   conn,
		config: config,
		send:   make(chan []byte, 512),
		stop:   make(chan []byte, 1),
	}

	s.ctx, s.ctxCancel = context.WithCancel(ctx)

	s.logger = logger.AddRequestID(s.ctx, logger.Default().Named("websocket"))

	return s
}

func (s *Session) log(fields ...zapcore.Field) *zap.Logger {
	return s.logger.With(fields...)
}

func (s *Session) Context() context.Context {
	return s.ctx
}

func (s *Session) connected() (err error) {
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
				s.sendPresence("")
			}
		}
	}()

	return nil
}

func (s *Session) disconnected() {
	// Tell everyone that user has disconnected
	_ = s.sendPresence("disconnected")

	// Cancel context
	s.ctxCancel()

	// Close connection
	s.Conn.Close()
	s.Conn = nil
}

// Sends user presence information to all subscribers
//
// It sends "connected", "disconnected" and "" activity kinds
func (s *Session) sendPresence(kind string) error {
	connections := store.CountConnections(s.user.Identity())
	if kind == "disconnected" {
		connections--
	}

	return nil
}

func (s *Session) Handle() (err error) {
	if err = s.connected(); err != nil {
		s.Close()
		return
	}

	go s.readLoop()
	return s.writeLoop()
}

func (s *Session) Close() {
	s.once.Do(func() {
		s.disconnected()
		store.Delete(s.id)
	})
}

func (s *Session) readLoop() (err error) {
	defer func() {
		s.Close()
	}()

	if err = s.Conn.SetReadDeadline(time.Now().Add(s.config.PingTimeout)); err != nil {
		return
	}

	s.Conn.SetPongHandler(func(string) error {
		return s.Conn.SetReadDeadline(time.Now().Add(s.config.PingTimeout))
	})

	s.remoteAddr = s.Conn.RemoteAddr().String()

	for {
		_, raw, err := s.Conn.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "s.readLoop")
		}

		if err = s.dispatch(raw); err != nil {
			s.log(zap.Error(err)).Error("could not dispatch")
			_ = s.sendReply(outgoing.NewError(err))
		}
	}
}

func (s *Session) writeLoop() error {
	ticker := time.NewTicker(s.config.PingPeriod)

	defer func() {
		ticker.Stop()
		s.Close() // break readLoop
	}()

	write := func(msg []byte) (err error) {
		if s.Conn == nil {
			// Connection closed, nowhere to write
			return
		}

		if err = s.Conn.SetWriteDeadline(time.Now().Add(s.config.Timeout)); err != nil {
			return
		}

		if msg != nil && s.Conn != nil {
			return s.Conn.WriteMessage(websocket.TextMessage, msg)
		}

		return
	}

	ping := func() (err error) {
		if s.Conn == nil {
			// Connection closed, nothing to ping
			return
		}

		if err = s.Conn.SetWriteDeadline(time.Now().Add(s.config.Timeout)); err != nil {
			return
		}

		if s.Conn != nil {
			return s.Conn.WriteMessage(websocket.PingMessage, nil)
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
