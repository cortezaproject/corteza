package websocket

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/outgoing"
	"github.com/crusttech/crust/messaging/internal/repository"
	"github.com/crusttech/crust/messaging/internal/service"
	"github.com/crusttech/crust/messaging/types"
)

type (
	// Session
	Session struct {
		id   uint64
		once sync.Once
		conn *websocket.Conn
		ctx  context.Context

		subs *Subscriptions

		send chan []byte
		stop chan []byte

		remoteAddr string

		config *repository.Flags

		user auth.Identifiable

		svc struct {
			ch  service.ChannelService
			msg service.MessageService
		}
	}
)

func (Session) New(ctx context.Context, config *repository.Flags, conn *websocket.Conn) *Session {
	s := &Session{
		conn:   conn,
		ctx:    ctx,
		config: config,
		subs:   NewSubscriptions(),
		send:   make(chan []byte, 512),
		stop:   make(chan []byte, 1),
	}

	s.svc.ch = service.DefaultChannel
	s.svc.msg = service.DefaultMessage

	return s
}

func (sess *Session) Context() context.Context {
	return sess.ctx
}

func (sess *Session) connected() (err error) {
	var (
		cc types.ChannelSet
	)

	// Push user info about all channels he has access to...
	if cc, err = sess.svc.ch.With(sess.ctx).Find(&types.ChannelFilter{}); err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Subscribing %d to %d channels", sess.user.Identity(), len(cc))

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
	if err = sess.sendToAll(&outgoing.Connected{UserID: payload.Uint64toa(sess.user.Identity())}); err != nil {
		return
	}

	return nil
}

func (sess *Session) disconnected() {
	// Tell everyone that user has disconnected
	_ = sess.sendToAll(&outgoing.Disconnected{UserID: payload.Uint64toa(sess.user.Identity())})
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
		sess.conn.Close()
		sess.conn = nil
		store.Delete(sess.id)
	})
}

func (sess *Session) readLoop() (err error) {
	defer func() {
		sess.Close()
	}()

	if err = sess.conn.SetReadDeadline(time.Now().Add(sess.config.Websocket.PingTimeout)); err != nil {
		return
	}

	sess.conn.SetPongHandler(func(string) error {
		return sess.conn.SetReadDeadline(time.Now().Add(sess.config.Websocket.PingTimeout))
	})

	sess.remoteAddr = sess.conn.RemoteAddr().String()

	for {
		_, raw, err := sess.conn.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "sess.readLoop")
		}

		if err = sess.dispatch(raw); err != nil {
			log.Printf("Error: %v", err)

			sess.sendReply(outgoing.NewError(err))
		}
	}
}

func (sess *Session) writeLoop() error {
	ticker := time.NewTicker(sess.config.Websocket.PingPeriod)

	defer func() {
		ticker.Stop()
		sess.Close() // break readLoop
	}()

	write := func(msg []byte) (err error) {
		if sess.conn == nil {
			// Connection closed, nowhere to write
			return
		}

		if err = sess.conn.SetWriteDeadline(time.Now().Add(sess.config.Websocket.Timeout)); err != nil {
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

		if err = sess.conn.SetWriteDeadline(time.Now().Add(sess.config.Websocket.Timeout)); err != nil {
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
