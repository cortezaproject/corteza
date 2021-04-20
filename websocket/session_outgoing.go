package websocket

import (
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"time"
)

// SendReply sends message only on this session, no need to enqueue item
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
		logger.Default().Warn("websocket.sendBytes send timeout")
	}
	return nil
}
