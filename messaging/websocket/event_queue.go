package websocket

import (
	"context"
	"encoding/json"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/internal/payload"
	"github.com/cortezaproject/corteza-server/internal/payload/outgoing"
	"github.com/cortezaproject/corteza-server/messaging/internal/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	eventQueueWalker interface {
		Walk(func(session *Session))
	}

	eventQueue struct {
		origin uint64
		queue  chan *types.EventQueueItem
	}
)

const (
	eventQueueBacklog = 512
)

var eq *eventQueue

func init() {
	eq = EventQueue(factory.Sonyflake.NextID())
}

func EventQueue(origin uint64) *eventQueue {
	return &eventQueue{
		origin: origin,
		queue:  make(chan *types.EventQueueItem, eventQueueBacklog),
	}
}

// @todo: retire this function, use Events().Push(ctx, item) directly.
func (eq *eventQueue) store(ctx context.Context, qp repository.EventsRepository) {
	go func() {
		for {
			select {
			case <-ctx.Done():
			case eqi := <-eq.queue:
				qp.Push(ctx, eqi)
			}
		}
	}()
}

func (eq *eventQueue) feedSessions(ctx context.Context, qp repository.EventsRepository, store eventQueueWalker) error {
	for {
		item, err := qp.Pull(ctx)
		if err != nil {
			return err
		}

		if item.SubType == types.EventQueueItemSubTypeUser {
			p := &outgoing.Payload{}

			if err := json.Unmarshal(item.Payload, p); err != nil {
				return err
			}

			if p.ChannelJoin != nil {
				// Handle subscribing

				// This store.Walk handler does not send to subscribed sessions but
				// subscribes all sessions that belong to the same user
				store.Walk(func(s *Session) {
					if payload.Uint64toa(s.user.Identity()) == p.ChannelJoin.UserID {
						s.subs.Add(p.ChannelJoin.ID)
					}
				})
			} else if p.ChannelPart != nil {
				// Handle un-subscribing

				// This store.Walk handler does not send to subscribed sessions but
				// subscribes all sessions that belong to the same user
				store.Walk(func(s *Session) {
					if payload.Uint64toa(s.user.Identity()) == p.ChannelPart.UserID {
						s.subs.Delete(p.ChannelPart.ID)
					}
				})
			} else {
				// No other payload types we can handle at the moment.
				return nil
			}

		} else if item.Subscriber == "" {
			// Distribute payload to all connected sessions
			store.Walk(func(s *Session) {
				s.sendBytes(item.Payload)
			})
		} else {
			// Distribute payload to specific subscribers
			store.Walk(func(s *Session) {
				if s.subs.Get(item.Subscriber) != nil {
					s.sendBytes(item.Payload)
				}
			})
		}

	}
}

// Adds origin to the event and puts it into queue.
// func (eq *eventQueue) push(ctx context.Context, eqi *types.EventQueueItem) {
// 	eqi.Origin = eq.origin
//
// 	select {
// 	case <-ctx.Done():
// 	case eq.queue <- eqi:
// 	}
// }
