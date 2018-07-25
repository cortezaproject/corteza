package websocket

import (
	"encoding/json"
	"github.com/crusttech/crust/sam/websocket/incoming"
	"github.com/pkg/errors"
)

func (s *Session) dispatch(raw []byte) (err error) {
	var p = &incoming.Payload{}
	if err = json.Unmarshal(raw, p); err != nil {
		return errors.Wrap(err, "Session.incoming: payload malformed")
	}

	ctx := s.Context()
	switch {

	// message actions
	case p.MessageCreate != nil:
		return s.messageCreate(ctx, p)
	case p.MessageUpdate != nil:
		return s.messageUpdate(ctx, p)
	case p.MessageDelete != nil:
		return s.messageDelete(ctx, p)

	// channel actions
	case p.ChannelJoin != nil:
		return s.channelJoin(ctx, p)
	case p.ChannelPart != nil:
		return s.channelPart(ctx, p)
	case p.ChannelOpen != nil:
		return s.channelOpen(ctx, p)

	}

	return nil
}
