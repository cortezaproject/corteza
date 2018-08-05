package websocket

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"log"
	"time"
)

type (
	eventQueuePuller interface {
		EventQueuePull(origin uint64) ([]*types.EventQueueItem, error)
		EventQueueSync(origin uint64, ID uint64) error
	}
	eventQueuePusher interface {
		EventQueuePush(*types.EventQueueItem) error
	}

	eventQueueWalker interface {
		Walk(func(session *Session))
	}

	eventQueue struct {
		origin uint64
		queue  chan *types.EventQueueItem
	}
)

var eq *eventQueue

func init() {
	eq = EventQueue(factory.Sonyflake.NextID())
}

func EventQueue(origin uint64) *eventQueue {
	return &eventQueue{
		origin: origin,
		queue:  make(chan *types.EventQueueItem, 512),
	}
}

func (eq *eventQueue) store(ctx context.Context, qp eventQueuePusher) {
	go func() {
		for {
			select {
			case <-ctx.Done():
			case eqi := <-eq.queue:
				qp.EventQueuePush(eqi)
			}
		}
	}()
}

func (eq *eventQueue) feedSessions(ctx context.Context, qp eventQueuePuller, store eventQueueWalker) {
	var items []*types.EventQueueItem

	go func() {
		var err error
	mainLoop:
		for {
			select {
			case <-ctx.Done():
				log.Printf("Error: %v", ctx.Err())
			case <-time.After(time.Second * 1):
				// How often do we check the database for new events?
				// @todo make this interval configurable
			}

			for {
				items, err = qp.EventQueuePull(eq.origin)
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}

				if len(items) == 0 {
					// No more items to sync, continue the mainLoop loop
					continue mainLoop
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
					qp.EventQueueSync(eq.origin, lastSyncedId)
				}

			}
		}
	}()

	return
}

// Adds origin to the event and puts it into queue.
func (eq *eventQueue) push(ctx context.Context, eqi *types.EventQueueItem) {
	eqi.Origin = eq.origin

	select {
	case <-ctx.Done():
	case eq.queue <- eqi:
	}
}
