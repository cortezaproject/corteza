package websocket

import (
	"context"
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/websocket/incoming"
	"github.com/crusttech/crust/sam/websocket/outgoing"
)

func (s *Session) channelJoin(ctx context.Context, p *incoming.ChannelJoin) error {
	// @todo: check access / can we join this channel (should be done by service layer)

	s.subs.Add(p.ChannelID)

	// Telling all subscribers of the channel we're joining that we are joining.
	var chJoin = &outgoing.ChannelJoin{
		ID:     p.ChannelID,
		UserID: uint64toa(auth.GetIdentityFromContext(ctx).GetID()),
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
		UserID: uint64toa(auth.GetIdentityFromContext(ctx).GetID()),
	}

	s.sendToAllSubscribers(chPart, p.ChannelID)

	return nil
}

func (s *Session) channelList(ctx context.Context, p *incoming.ChannelList) error {
	channels, err := service.Channel().Find(ctx, nil)
	if err != nil {
		return err
	}

	return s.sendReply(payloadFromChannels(channels))
}

func (s *Session) channelRename(ctx context.Context, p *incoming.ChannelRename) error {
	ch, err := service.Channel().FindByID(ctx, parseUInt64(p.ChannelID))
	if err != nil {
		return err
	}

	if ch.Name == p.Name {
		// No changes, ignore
		return nil
	}

	ch.Name = p.Name

	ch, err = service.Channel().Update(ctx, ch)
	if err != nil {
		return err
	}

	return s.sendToAllSubscribers(payloadFromChannel(ch), p.ChannelID)
}

func (s *Session) channelChangeTopic(ctx context.Context, p *incoming.ChannelChangeTopic) error {
	ch, err := service.Channel().FindByID(ctx, parseUInt64(p.ChannelID))
	if err != nil {
		return err
	}

	if ch.Topic == p.Topic {
		// No changes, ignore
		return nil
	}

	ch.Topic = p.Topic

	ch, err = service.Channel().Update(ctx, ch)
	if err != nil {
		return err
	}

	return s.sendToAllSubscribers(payloadFromChannel(ch), p.ChannelID)
}
