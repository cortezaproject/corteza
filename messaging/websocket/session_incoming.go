package websocket

import (
	"github.com/cortezaproject/corteza-server/internal/payload"
	"github.com/pkg/errors"
)

func (s *Session) dispatch(raw []byte) error {
	var p, err = payload.Unmarshal(raw)
	if err != nil {
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

	// channel actions
	case p.ChannelJoin != nil:
		return s.channelJoin(ctx, p.ChannelJoin)
	case p.ChannelPart != nil:
		return s.channelPart(ctx, p.ChannelPart)
	case p.Channels != nil:
		return s.channelList(ctx, p.Channels)
	case p.ChannelCreate != nil:
		return s.channelCreate(ctx, p.ChannelCreate)
	case p.ChannelUpdate != nil:
		return s.channelUpdate(ctx, p.ChannelUpdate)

	// @deprecated
	case p.ChannelViewRecord != nil:
		return s.channelViewRecord(ctx, p.ChannelViewRecord)
	}

	return nil
}
