package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/payload/outgoing"
)

type (
	event struct {
		events repository.EventsRepository
	}

	EventService interface {
		Activity(ctx context.Context, a *types.Activity) error
		Message(ctx context.Context, m *types.Message) error
		MessageFlag(ctx context.Context, m *types.MessageFlag) error
		UnreadCounters(ctx context.Context, uu types.UnreadSet) error
		Channel(ctx context.Context, m *types.Channel) error
		Join(ctx context.Context, userID, channelID uint64) error
		Part(ctx context.Context, userID, channelID uint64) error
	}
)

// Event sends sends events back to all (or specific) subscribers
func Event() EventService {
	return &event{events: repository.Events()}
}

// Message sends message events to subscribers
func (svc event) Message(ctx context.Context, m *types.Message) error {
	return svc.push(ctx, payload.Message(ctx, m), types.EventQueueItemSubTypeChannel, m.ChannelID)
}

// Activity sends activity event to subscribers
func (svc event) Activity(ctx context.Context, a *types.Activity) error {
	return svc.push(ctx, payload.Activity(a), types.EventQueueItemSubTypeChannel, a.ChannelID)
}

// MessageFlag sends message flag events to subscribers
func (svc event) MessageFlag(ctx context.Context, f *types.MessageFlag) error {
	var p outgoing.MessageEncoder

	switch {
	case f.IsBookmark():
		// Leaving this here so it is obvious.
		return nil
	case f.IsPin() && f.DeletedAt != nil:
		p = payload.MessagePinRemoved(f)
	case f.IsReaction() && f.DeletedAt != nil:
		p = payload.MessageReactionRemoved(f)
	case f.IsPin():
		p = payload.MessagePin(f)
	case f.IsReaction():
		p = payload.MessageReaction(f)
	default:
		return nil
	}

	return svc.push(ctx, p, types.EventQueueItemSubTypeChannel, f.ChannelID)
}

func (svc event) UnreadCounters(ctx context.Context, uu types.UnreadSet) error {
	return uu.Walk(func(u *types.Unread) error {
		return svc.push(ctx, payload.Unread(u), types.EventQueueItemSubTypeUser, u.UserID)
	})
}

// Channel notifies subscribers about channel change
//
// If this is a public channel we notify everyone
func (svc event) Channel(ctx context.Context, ch *types.Channel) error {
	var sub uint64 = 0

	if ch.Type != types.ChannelTypePublic {
		sub = ch.ID
	}

	return svc.push(ctx, payload.Channel(ch), types.EventQueueItemSubTypeChannel, sub)
}

// Join sends force-channel-join event to all matching subscriptions
//
// Subscription will match when session's user ID is the same as sub
// We pack channel ID (the id to subscribe to) as payload
func (svc event) Join(ctx context.Context, userID, channelID uint64) (err error) {
	join := payload.ChannelJoin(channelID, userID)

	// Subscribe user to the channel
	if err = svc.push(ctx, join, types.EventQueueItemSubTypeUser, userID); err != nil {
		return
	}

	// Let subscribers know that this user has joined the channel
	if err = svc.push(ctx, join, types.EventQueueItemSubTypeChannel, channelID); err != nil {
		return
	}

	return
}

// Part sends force-channel-part event to all matching subscriptions
//
// Subscription will match when session's user ID is the same as sub
// We pack channel ID (the id to subscribe to) as payload
func (svc event) Part(ctx context.Context, userID, channelID uint64) (err error) {
	part := payload.ChannelPart(channelID, userID)

	// Let subscribers know that this user has parted the channel
	if err = svc.push(ctx, part, types.EventQueueItemSubTypeChannel, channelID); err != nil {
		return
	}

	// Subscribe user to the channel
	if err = svc.push(ctx, part, types.EventQueueItemSubTypeUser, userID); err != nil {
		return
	}

	return
}

func (svc event) push(ctx context.Context, m outgoing.MessageEncoder, subType types.EventQueueItemSubType, sub uint64) error {
	var enc, err = m.EncodeMessage()
	if err != nil {
		return err
	}

	item := &types.EventQueueItem{Payload: enc, SubType: subType}

	if sub > 0 {
		item.Subscriber = payload.Uint64toa(sub)
	}

	return svc.events.Push(ctx, item)
}
