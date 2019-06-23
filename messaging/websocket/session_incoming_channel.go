package websocket

import (
	"context"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/payload"
	"github.com/cortezaproject/corteza-server/internal/payload/incoming"
	"github.com/cortezaproject/corteza-server/internal/payload/outgoing"
	"github.com/cortezaproject/corteza-server/messaging/types"
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
	channels, err := s.svc.ch.With(ctx).Find(&types.ChannelFilter{})
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

	_, err = s.svc.ch.With(ctx).Create(ch)
	if err != nil {
		return err
	}

	return nil
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
