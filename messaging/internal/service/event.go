package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/payload"
	"github.com/cortezaproject/corteza-server/internal/payload/outgoing"
	"github.com/cortezaproject/corteza-server/messaging/internal/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	event struct {
		ctx    context.Context
		logger *zap.Logger

		events repository.EventsRepository
	}

	EventService interface {
		With(ctx context.Context) EventService
		Activity(a *types.Activity) error
		Message(m *types.Message) error
		MessageFlag(m *types.MessageFlag) error
		Channel(m *types.Channel) error
		Join(userID, channelID uint64) error
		Part(userID, channelID uint64) error
	}
)

// Event sends sends events back to all (or specific) subscribers
func Event(ctx context.Context) EventService {
	return (&event{
		logger: DefaultLogger.Named("event"),
	}).With(ctx)
}

func (svc event) With(ctx context.Context) EventService {
	return &event{
		ctx:    ctx,
		logger: svc.logger,

		events: repository.Events(),
	}
}

// log() returns zap's logger with requestID from current context and fields.
// func (svc event) log(fields ...zapcore.Field) *zap.Logger {
// 	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
// }

// Message sends message events to subscribers
func (svc event) Message(m *types.Message) error {
	return svc.push(payload.Message(svc.ctx, m), types.EventQueueItemSubTypeChannel, m.ChannelID)
}

// Activity sends activity event to subscribers
func (svc event) Activity(a *types.Activity) error {
	return svc.push(payload.Activity(a), types.EventQueueItemSubTypeChannel, a.ChannelID)
}

// MessageFlag sends message flag events to subscribers
func (svc event) MessageFlag(f *types.MessageFlag) error {
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

	return svc.push(p, types.EventQueueItemSubTypeChannel, f.ChannelID)
	return nil
}

// Channel notifies subscribers about channel change
//
// If this is a public channel we notify everyone
func (svc event) Channel(ch *types.Channel) error {
	var sub uint64 = 0

	if ch.Type != types.ChannelTypePublic {
		sub = ch.ID
	}

	return svc.push(payload.Channel(ch), types.EventQueueItemSubTypeChannel, sub)
}

// Join sends force-channel-join event to all matching subscriptions
//
// Subscription will match when session's user ID is the same as sub
// We pack channel ID (the id to subscribe to) as payload
func (svc event) Join(userID, channelID uint64) (err error) {
	join := payload.ChannelJoin(channelID, userID)

	// Subscribe user to the channel
	if err = svc.push(join, types.EventQueueItemSubTypeUser, 0); err != nil {
		return
	}

	// Let subscribers know that this user has joined the channel
	if err = svc.push(join, types.EventQueueItemSubTypeChannel, channelID); err != nil {
		return
	}

	return
}

// Part sends force-channel-part event to all matching subscriptions
//
// Subscription will match when session's user ID is the same as sub
// We pack channel ID (the id to subscribe to) as payload
func (svc event) Part(userID, channelID uint64) (err error) {
	part := payload.ChannelPart(channelID, userID)

	// Let subscribers know that this user has parted the channel
	if err = svc.push(part, types.EventQueueItemSubTypeChannel, channelID); err != nil {
		return
	}

	// Subscribe user to the channel
	if err = svc.push(part, types.EventQueueItemSubTypeUser, 0); err != nil {
		return
	}

	return
}

func (svc event) push(m outgoing.MessageEncoder, subType types.EventQueueItemSubType, sub uint64) error {
	var enc, err = m.EncodeMessage()
	if err != nil {
		return err
	}

	item := &types.EventQueueItem{Payload: enc, SubType: subType}

	if sub > 0 {
		item.Subscriber = payload.Uint64toa(sub)
	}

	return svc.events.Push(svc.ctx, item)
}
