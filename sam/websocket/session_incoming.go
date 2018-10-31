package websocket

import (
	"github.com/crusttech/crust/internal/payload"
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
	case p.Messages != nil:
		return s.messageHistory(ctx, p.Messages)
	case p.MessageThreads != nil:
		return s.messageThreads(ctx, p.MessageThreads)

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
	case p.ChannelViewRecord != nil:
		return s.channelViewRecord(ctx, p.ChannelViewRecord)
	case p.ChannelActivity != nil:
		return s.channelActivity(ctx, p.ChannelActivity)

	case p.Users != nil:
		return s.userList(ctx, p.Users)

	case p.ExecCommand != nil:
		return s.execCommand(ctx, p.ExecCommand)
	}

	return nil
}
