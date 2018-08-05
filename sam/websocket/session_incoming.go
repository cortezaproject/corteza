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
		return s.messageCreate(ctx, p.MessageCreate)
	case p.MessageUpdate != nil:
		return s.messageUpdate(ctx, p.MessageUpdate)
	case p.MessageDelete != nil:
		return s.messageDelete(ctx, p.MessageDelete)
	case p.MessageHistory != nil:
		return s.messageHistory(ctx, p.MessageHistory)

	// channel actions
	case p.ChannelJoin != nil:
		return s.channelJoin(ctx, p.ChannelJoin)
	case p.ChannelPart != nil:
		return s.channelPart(ctx, p.ChannelPart)
	case p.ChannelList != nil:
		return s.channelList(ctx, p.ChannelList)
	case p.ChannelRename != nil:
		return s.channelRename(ctx, p.ChannelRename)
	case p.ChannelChangeTopic != nil:
		return s.channelChangeTopic(ctx, p.ChannelChangeTopic)

	}

	return nil
}
