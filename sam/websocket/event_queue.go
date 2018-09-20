package websocket

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
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

func (eq *eventQueue) store(ctx context.Context, qp repository.Events) {
	go func() {
		for {
			select {
			case <-ctx.Done():
			case eqi := <-eq.queue:
				qp.Push(eqi)
			}
		}
	}()
}

func (eq *eventQueue) feedSessions(ctx context.Context, config *repository.Flags, qp repository.Events, store eventQueueWalker) error {
	newMessageEvent := make(chan struct{}, eventQueueBacklog)
	done := make(chan error, 1)

	pubsub := service.PubSub()
	go func() {
		onConnect := func() error {
			return nil
		}
		onMessage := func(message string, payload []byte) error {
			newMessageEvent <- struct{}{}
			return nil
		}
		done <- pubsub.Subscribe(ctx, "events", onConnect, onMessage)
	}()

	poll := func() error {
		for {
			items, err := qp.Pull(eq.origin)
			if err != nil {
				return err
			}
			if len(items) == 0 {
				return nil
			}

			var lastSyncedId uint64

			for _, item := range items {
				if item.Subscriber == "" {
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

				lastSyncedId = item.ID

			}

			if lastSyncedId > 0 {
				qp.Sync(eq.origin, lastSyncedId)
			}
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-newMessageEvent:
			if err := poll(); err != nil {
				return err
			}
		case err := <-done:
			return err
		}
	}
}

// Adds origin to the event and puts it into queue.
func (eq *eventQueue) push(ctx context.Context, eqi *types.EventQueueItem) {
	eqi.Origin = eq.origin

	select {
	case <-ctx.Done():
	case eq.queue <- eqi:
	}
}
