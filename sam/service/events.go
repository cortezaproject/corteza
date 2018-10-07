package service

import (
	"context"

	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/outgoing"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	event struct {
		ctx context.Context

		events repository.EventsRepository
	}

	EventService interface {
		With(ctx context.Context) EventService
		Message(m *types.Message) error
		Channel(m *types.Channel) error
	}
)

// Event sends sends events back to all (or specific) subscribers
func Event() EventService {
	return (&event{events: repository.Events()}).With(context.Background())
}

func (svc *event) With(ctx context.Context) EventService {
	return &event{
		ctx:    ctx,
		events: svc.events,
	}
}

// Message sends message events to subscribers
func (svc *event) Message(m *types.Message) error {
	return svc.push(payload.Message(m), types.EventQueueItemSubTypeChannel, m.ChannelID)
}

// Channel notifies subscribers about channel change
//
// If this is a public channel we notify everyone
func (svc *event) Channel(ch *types.Channel) error {
	sub := ch.ID
	if ch.Type == types.ChannelTypePublic {
		sub = 0
	}

	return svc.push(payload.Channel(ch), types.EventQueueItemSubTypeChannel, sub)
}

// Join sends force-channel-join event to all matching subscriptions
//
// Subscription will match when session's user ID is the same as sub
// We pack channel ID (the id to subscribe to) as payload
func (svc *event) Join(userID, channelID uint64) error {
	return svc.push(payload.Channel(&types.Channel{ID: channelID}), types.EventQueueItemSubTypeUser, userID)
}

func (svc *event) push(m outgoing.MessageEncoder, subType types.EventQueueItemSubType, sub uint64) error {
	var enc, err = m.EncodeMessage()
	if err != nil {
		return err
	}

	item := &types.EventQueueItem{Payload: enc}

	if sub > 0 {
		item.Subscriber = payload.Uint64toa(sub)
		item.SubType = subType
	}

	return svc.events.Push(svc.ctx, item)
}
