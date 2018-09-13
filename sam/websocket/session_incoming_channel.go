package websocket

import (
	"context"
	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/incoming"
	"github.com/crusttech/crust/sam/websocket/outgoing"
)

func (s *Session) channelJoin(ctx context.Context, p *incoming.ChannelJoin) error {
	// @todo: check access / can we join this channel (should be done by service layer)

	s.subs.Add(p.ChannelID)

	// Telling all subscribers of the channel we're joining that we are joining.
	var chJoin = &outgoing.ChannelJoin{
		ID:     p.ChannelID,
		UserID: uint64toa(auth.GetIdentityFromContext(ctx).Identity()),
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
		UserID: uint64toa(auth.GetIdentityFromContext(ctx).Identity()),
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

func (s *Session) channelCreate(ctx context.Context, p *incoming.ChannelCreate) (err error) {
	ch := &types.Channel{
		Type:  types.ChannelTypePublic,
		Name:  p.Name,
		Topic: p.Topic,
	}

	ch, err = service.Channel().Create(ctx, ch)
	if err != nil {
		return err
	}

	// Explicitly subscribe to newly created channel
	s.subs.Add(uint64toa(ch.ID))

	// @todo this should go over all user's sessons and subscribe there as well

	pl := payloadFromChannel(ch)

	if ch.Type == types.ChannelTypePublic {
		return s.sendToAll(pl)
	}

	// By default, just send reply to user
	return s.sendReply(pl)
}

func (s *Session) channelDelete(ctx context.Context, p *incoming.ChannelDelete) (err error) {
	err = service.Channel().Delete(ctx, parseUInt64(p.ChannelID))
	if err != nil {
		return err
	}

	return s.sendToAllSubscribers(&outgoing.ChannelDeleted{
		ID:     p.ChannelID,
		UserID: uint64toa(auth.GetIdentityFromContext(ctx).Identity()),
	}, p.ChannelID)
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
