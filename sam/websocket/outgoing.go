package websocket

import (
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/websocket/outgoing"
	"log"
	"time"
)

// Sends message to subscribers
func (s *Session) sendToAllSubscribers(p outgoing.MessageEncoder, channelID string) error {
	pb, err := p.EncodeMessage()
	if err != nil {
		return err
	}

	eq.push(s.ctx, &repository.EventQueueItem{Payload: pb, Subscriber: channelID})

	store.Walk(func(sess *Session) {
		// send message only to users with subscribed channels
		if sess.subs.Get(channelID) != nil {
			sess.sendBytes(pb)
		}
	})

	return nil
}

// Sends message to all connected clients
func (s *Session) sendToAll(p outgoing.MessageEncoder) error {
	pb, err := p.EncodeMessage()
	if err != nil {
		return err
	}

	eq.push(s.ctx, &repository.EventQueueItem{Payload: pb})

	store.Walk(func(sess *Session) {
		// send message only to users with subscribed channels
		sess.sendBytes(pb)
	})

	return nil
}

func (s *Session) sendReply(p outgoing.MessageEncoder) error {
	pb, err := p.EncodeMessage()
	if err != nil {
		return err
	}

	return s.sendBytes(pb)
}

func (s *Session) sendBytes(p []byte) error {
	select {
	case s.send <- p:
	case <-time.After(2 * time.Millisecond):
		log.Println("websocket.sendBytes send timeout")
	}
	return nil
}
