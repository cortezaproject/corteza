package websocket

import (
	"encoding/json"
	"github.com/crusttech/crust/sam/websocket/outgoing"
	"log"
	"time"
)

// Sends message to all connected subscribers/users
//
// If channelID is nil, it broadcasts payload to anyone connected,
// if it is set, it broadcasts only to clients subscribed to that channel
func (s *Session) broadcast(message outgoing.PayloadType, channelID *string) error {
	p := outgoing.Payload{}.New().Load(message)
	// encode message once and send bytes
	pb, err := json.Marshal(p)
	if err != nil {
		return err
	}

	store.Walk(func(sess *Session) {
		// send message only to users with subscribed channels
		if channelID == nil || sess.subs.Get(*channelID) != nil {
			sess.sendBytes(pb)
		}
	})

	return nil
}

func (s *Session) respond(message outgoing.PayloadType) error {
	p := outgoing.Payload{}.New().Load(message)
	select {
	case s.send <- p:
	case <-time.After(2 * time.Millisecond):
		log.Println("websocket.respond send timeout")
	}
	return nil
}

func (s *Session) sendBytes(p []byte) error {
	select {
	case s.send <- p:
	case <-time.After(2 * time.Millisecond):
		log.Println("websocket.sendBytes send timeout")
	}
	return nil
}
