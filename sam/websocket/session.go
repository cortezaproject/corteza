package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/crusttech/crust/sam/websocket/outgoing"
)

type (
	// Session
	Session struct {
		id   uint64
		conn *websocket.Conn

		subs *Subscriptions

		send   chan interface{}
		stop   chan interface{}
		detach chan string

		remoteAddr string

		config configuration
	}
)

func (Session) New(conn *websocket.Conn) *Session {
	return &Session{
		conn:   conn,
		config: config,
		subs:   Subscriptions{}.New(),
		send:   make(chan interface{}, 512),
		stop:   make(chan interface{}, 1),
		detach: make(chan string, 64),
	}
}

func (sess *Session) Handle() error {
	go sess.readLoop()
	return sess.writeLoop()
}

func (sess *Session) Close() {
	sess.conn.Close()
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
			// @todo: log error?
			sess.send <- outgoing.NewError(err)
		}
	}
}

func (sess *Session) writeLoop() error {
	ticker := time.NewTicker(sess.config.pingPeriod)

	defer func() {
		ticker.Stop()
		sess.Close() // break readLoop
	}()

	write := func(mt int, msg interface{}) error {
		sess.conn.SetWriteDeadline(time.Now().Add(sess.config.writeTimeout))

		switch msg := msg.(type) {
		case *outgoing.WsMessage:
			return sess.conn.WriteJSON(msg)
		default:
			return sess.conn.WriteMessage(mt, nil)
		}

	}

	for {
		select {
		case msg, ok := <-sess.send:
			if !ok {
				// channel closed
				return nil
			}

			if err := write(websocket.TextMessage, msg); err != nil {
				return errors.Wrap(err, "writeLoop send")
			}
		case msg := <-sess.stop:
			// Shutdown requested, don't care if the message is delivered
			if msg != nil {
				write(websocket.TextMessage, msg)
			}
			return nil

		case topic := <-sess.detach:
			sess.subs.Delete(topic)

		case <-ticker.C:
			if err := write(websocket.PingMessage, nil); err != nil {
				return errors.Wrap(err, "writeLoop ping")
			}
		}
	}
}
