package websocket

import (
	"context"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/incoming"
	"github.com/crusttech/crust/internal/payload/outgoing"
	"github.com/crusttech/crust/sam/types"
)

func (s *Session) channelJoin(ctx context.Context, p *incoming.ChannelJoin) error {
	// @todo: check access / can we join this channel (should be done by service layer)

	s.subs.Add(p.ChannelID)

	// Telling all subscribers of the channel we're joining that we are joining.
	var chJoin = &outgoing.ChannelJoin{
		ID:     p.ChannelID,
		UserID: payload.Uint64toa(auth.GetIdentityFromContext(ctx).Identity()),
	}

	// Send to all channel subscribers
	s.sendToAllSubscribers(chJoin, p.ChannelID)

	return nil
}

func (s *Session) channelPart(ctx context.Context, p *incoming.ChannelPart) error {
	// @todo: check access / can we part this channel? (should be done by service layer)

	// First, let's unsubscribe, so we don't hear echos
	s.subs.Delete(p.ChannelID)

	// This payload will tell everyone that we're parting from ALL channels
	var chPart = &outgoing.ChannelPart{
		ID:     p.ChannelID,
		UserID: payload.Uint64toa(auth.GetIdentityFromContext(ctx).Identity()),
	}

	s.sendToAllSubscribers(chPart, p.ChannelID)

	return nil
}

func (s *Session) channelList(ctx context.Context, p *incoming.Channels) error {
	channels, err := s.svc.ch.With(ctx).Find(&types.ChannelFilter{IncludeMembers: true})
	if err != nil {
		return err
	}

	// @todo count members for all channels

	return s.sendReply(payload.Channels(channels))
}

func (s *Session) channelCreate(ctx context.Context, p *incoming.ChannelCreate) (err error) {
	ch := &types.Channel{
		Type: types.ChannelTypePublic,
	}

	if p.Name != nil {
		ch.Name = *p.Name
	}

	if p.Topic != nil {
		ch.Topic = *p.Topic
	}

	if p.Type != nil {
		ch.Type = types.ChannelType(*p.Type)
	}

	ch, err = s.svc.ch.With(ctx).Create(ch)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) channelDelete(ctx context.Context, p *incoming.ChannelDelete) (err error) {
	return s.svc.ch.With(ctx).Delete(payload.ParseUInt64(p.ChannelID))
}

func (s *Session) channelUpdate(ctx context.Context, p *incoming.ChannelUpdate) error {
	ch, err := s.svc.ch.With(ctx).FindByID(payload.ParseUInt64(p.ID))
	if err != nil {
		return err
	}

	if p.Name != nil {
		ch.Name = *p.Name
	}

	if p.Topic != nil {
		ch.Topic = *p.Topic
	}

	if p.Type != nil {
		ch.Type = types.ChannelType(*p.Type)
	}

	_, err = s.svc.ch.With(ctx).Update(ch)
	return err
}

func (s *Session) channelViewRecord(ctx context.Context, p *incoming.ChannelViewRecord) error {
	var (
		channelID     = payload.ParseUInt64(p.ChannelID)
		lastMessageID = payload.ParseUInt64(p.LastMessageID)
		userID        = auth.GetIdentityFromContext(ctx).Identity()
	)

	if channelID == 0 || lastMessageID == 0 {
		return nil
	}

	return s.svc.ch.With(ctx).RecordView(channelID, userID, lastMessageID)
}
