package websocket

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/outgoing"
)

type (
	// Session
	Session struct {
		id   uint64
		conn *websocket.Conn
		ctx  context.Context

		subs *Subscriptions

		send chan []byte
		stop chan []byte

		remoteAddr string

		config configuration

		user *types.User
	}
)

func (Session) New(ctx context.Context, conn *websocket.Conn) *Session {
	return &Session{
		conn:   conn,
		ctx:    ctx,
		config: config,
		subs:   NewSubscriptions(),
		send:   make(chan []byte, 512),
		stop:   make(chan []byte, 1),
	}
}

func (sess *Session) Context() context.Context {
	return sess.ctx
}

func (sess *Session) connected() {
	// Tell everyone that user has connected
	sess.sendToAll(&outgoing.Connected{UserID: uint64toa(sess.user.ID)})

	// Subscribe this user to all channels
	if chs, err := service.Channel().FindByMembership(sess.ctx); err != nil {
		log.Printf("Error: %v", err)
	} else {
		for _, ch := range chs {
			sess.subs.Add(uint64toa(ch.ID))
		}
		log.Printf("Subscribing %d to %d channels", sess.user.ID, len(chs))
	}
}

func (sess *Session) disconnected() {
	// Tell everyone that user has disconnected
	sess.sendToAll(&outgoing.Disconnected{UserID: uint64toa(sess.user.ID)})
}

func (sess *Session) Handle() error {
	sess.connected()
	go sess.readLoop()
	return sess.writeLoop()
}

func (sess *Session) Close() {
	// @todo fix "writeLoop send: websocket: close sent" caused by sending disconnect here
	sess.disconnected()
	sess.conn.Close()
	store.Delete(sess.id)
}

func (sess *Session) readLoop() error {
	defer func() {
		log.Println("serveWebsocket - stop")
		sess.Close()
	}()

	sess.conn.SetReadDeadline(time.Now().Add(sess.config.pingTimeout))
	sess.conn.SetPongHandler(func(string) error {
		sess.conn.SetReadDeadline(time.Now().Add(sess.config.pingTimeout))
		return nil
	})
	sess.remoteAddr = sess.conn.RemoteAddr().String()

	for {
		_, raw, err := sess.conn.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "sess.readLoop")
		}

		if err := sess.dispatch(raw); err != nil {
			log.Printf("Error: %v", err)
			sess.sendReply(outgoing.NewError(err))
		}
	}
}

func (sess *Session) writeLoop() error {
	ticker := time.NewTicker(sess.config.pingPeriod)

	defer func() {
		ticker.Stop()
		sess.Close() // break readLoop
	}()

	write := func(msg []byte) error {
		sess.conn.SetWriteDeadline(time.Now().Add(sess.config.writeTimeout))
		if msg != nil {
			return sess.conn.WriteMessage(websocket.TextMessage, msg)
		}
		return nil
	}

	ping := func() error {
		sess.conn.SetWriteDeadline(time.Now().Add(sess.config.writeTimeout))
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
			write(msg)
			return nil

		case <-ticker.C:
			if err := ping(); err != nil {
				return errors.Wrap(err, "writeLoop ping")
			}
		}
	}
}
