package websocket

import (
	"encoding/json"
	"log"

	"github.com/crusttech/crust/sam/websocket/incoming"
	"github.com/pkg/errors"
)

func (s *Session) dispatch(raw []byte) error {
	log.Printf("%s> %s", s.remoteAddr, string(raw))

	msg := incoming.Message{}.New()
	if err := json.Unmarshal(raw, msg); err != nil {
		return errors.Wrap(err, "Session.incoming: malformed json payload")
	}

	// @todo: do stuff with msg

	return nil
}
