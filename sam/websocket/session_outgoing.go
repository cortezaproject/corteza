package websocket

import (
	"log"
	"time"

	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

// Sends message to subscribers
func (s *Session) sendToAllSubscribers(p MessageEncoder, channelID string) error {
	pb, err := p.EncodeMessage()
	if err != nil {
		return err
	}

	return repository.Events().Push(s.ctx, &types.EventQueueItem{Payload: pb, Subscriber: channelID})
}

// Sends message to all connected clients
func (s *Session) sendToAll(p MessageEncoder) error {
	pb, err := p.EncodeMessage()
	if err != nil {
		return err
	}

	return repository.Events().Push(s.ctx, &types.EventQueueItem{Payload: pb})
}

// @todo: this isn't going to be correct - a user may have open multiple clients,
//        that will connect to different edge SAM servers. It should also go
//        through a repository.Events().Push (EventQueueItem) path.
func (s *Session) sendReply(p MessageEncoder) error {
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
