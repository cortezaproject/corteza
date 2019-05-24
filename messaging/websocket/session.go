package websocket

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/payload"
	"github.com/cortezaproject/corteza-server/internal/payload/outgoing"
	"github.com/cortezaproject/corteza-server/messaging/internal/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

type (
	// Session
	Session struct {
		id   uint64
		once sync.Once
		conn *websocket.Conn

		ctx       context.Context
		ctxCancel context.CancelFunc

		logger *zap.Logger

		subs *Subscriptions

		send chan []byte
		stop chan []byte

		remoteAddr string

		config *Config

		user auth.Identifiable

		svc struct {
			ch  service.ChannelService
			msg service.MessageService
		}
	}
)

func (Session) New(ctx context.Context, config *Config, conn *websocket.Conn) *Session {

	s := &Session{
		conn:   conn,
		config: config,
		subs:   NewSubscriptions(),
		send:   make(chan []byte, 512),
		stop:   make(chan []byte, 1),
	}

	s.ctx, s.ctxCancel = context.WithCancel(ctx)

	s.svc.ch = service.DefaultChannel
	s.svc.msg = service.DefaultMessage

	s.logger = logger.AddRequestID(s.ctx, logger.Default().Named("websocket"))

	return s
}

func (sess Session) log(fields ...zapcore.Field) *zap.Logger {
	return sess.logger.With(fields...)
}

func (sess *Session) Context() context.Context {
	return sess.ctx
}

func (sess *Session) connected() (err error) {
	var (
		cc types.ChannelSet
	)

	// Push user info about all channels he has access to...
	// @todo filter out all muted/non-joined channels
	if cc, err = sess.svc.ch.With(sess.ctx).Find(&types.ChannelFilter{}); err != nil {
		sess.log(zap.Error(err)).Error("Could not load subscribed channels")
	} else {
		sess.log().Debug(
			"websocket session connected",
			zap.Uint64("userID", sess.user.Identity()),
			zap.Int("subscriptions", len(cc)),
		)

		err = cc.Walk(func(c *types.Channel) error {
			// Subscribe this user/session to all channels
			sess.subs.Add(payload.Uint64toa(c.ID))
			return nil
		})

		if err != nil {
			return
		}
	}

	// Tell everyone that user has connected
	if err = sess.sendPresence("connected"); err != nil {
		return
	}

	// Create a heartbeat every minute for this user
	go func() {
		t := time.NewTicker(time.Second * 60)
		for {
			select {
			case <-sess.ctx.Done():
				return
			case <-t.C:
				sess.sendPresence("")
			}
		}
	}()

	return nil
}

func (sess *Session) disconnected() {
	// Tell everyone that user has disconnected
	_ = sess.sendPresence("disconnected")

	// Cancel context
	sess.ctxCancel()

	// Close connection
	sess.conn.Close()
	sess.conn = nil
}

// Sends user presence information to all subscribers
//
// It sends "connected", "disconnected" and "" activity kinds
func (sess *Session) sendPresence(kind string) error {
	connections := store.CountConnections(sess.user.Identity())
	if kind == "disconnected" {
		connections--
	}

	// Tell everyone that user has disconnected
	return sess.sendToAll(&outgoing.Activity{
		UserID:  sess.user.Identity(),
		Kind:    kind,
		Present: connections > 0,
	})
}

func (sess *Session) Handle() (err error) {
	if err = sess.connected(); err != nil {
		sess.Close()
		return
	}

	go sess.readLoop()
	return sess.writeLoop()
}

func (sess *Session) Close() {
	sess.once.Do(func() {
		sess.disconnected()
		store.Delete(sess.id)
	})
}

func (sess *Session) readLoop() (err error) {
	defer func() {
		sess.Close()
	}()

	if err = sess.conn.SetReadDeadline(time.Now().Add(sess.config.PingTimeout)); err != nil {
		return
	}

	sess.conn.SetPongHandler(func(string) error {
		return sess.conn.SetReadDeadline(time.Now().Add(sess.config.PingTimeout))
	})

	sess.remoteAddr = sess.conn.RemoteAddr().String()

	for {
		_, raw, err := sess.conn.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "sess.readLoop")
		}

		if err = sess.dispatch(raw); err != nil {
			sess.log(zap.Error(err)).Error("could not dispatch")
			_ = sess.sendReply(outgoing.NewError(err))
		}
	}
}

func (sess *Session) writeLoop() error {
	ticker := time.NewTicker(sess.config.PingPeriod)

	defer func() {
		ticker.Stop()
		sess.Close() // break readLoop
	}()

	write := func(msg []byte) (err error) {
		if sess.conn == nil {
			// Connection closed, nowhere to write
			return
		}

		if err = sess.conn.SetWriteDeadline(time.Now().Add(sess.config.Timeout)); err != nil {
			return
		}

		if msg != nil {
			return sess.conn.WriteMessage(websocket.TextMessage, msg)
		}

		return
	}

	ping := func() (err error) {
		if sess.conn == nil {
			// Connection closed, nothing to ping
			return
		}

		if err = sess.conn.SetWriteDeadline(time.Now().Add(sess.config.Timeout)); err != nil {
			return
		}

		return sess.conn.WriteMessage(websocket.PingMessage, nil)
	}

	for {
		select {
		case msg, ok := <-sess.send:
			if !ok {
				// channel closed
				return nil
			}

			if err := write(msg); err != nil {
				return errors.Wrap(err, "writeLoop send")
			}

		case msg := <-sess.stop:
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
