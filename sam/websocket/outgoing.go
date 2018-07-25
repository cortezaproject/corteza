package websocket

import (
	"log"
	"encoding/json"
	"time"
	"github.com/crusttech/crust/sam/websocket/outgoing"
)

func (s *Session) sendMessageChannel(channelID string, message outgoing.PayloadType) error {
	p := outgoing.Payload{}.New().Load(message)
	// encode message once and send bytes
	pb, err := json.Marshal(p)
	if err != nil {
		return err
	}
	store.Walk(func(sess *Session) {
		// send message only to users with subscribed channels
		if sess.subs.Get(channelID) != nil {
			select {
			case sess.send <- pb:
			case <-time.After(2 * time.Millisecond):
				log.Println("websocket.messageChannel send timeout")
			}
		}
	})
	return nil
}

func (s *Session) sendMessage(message outgoing.PayloadType) error {
	p := outgoing.Payload{}.New().Load(message)
	select {
	case s.send <- p:
	case <-time.After(2 * time.Millisecond):
		log.Println("websocket.messageChannel send timeout")
	}
	return nil
}
