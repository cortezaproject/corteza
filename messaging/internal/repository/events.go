package repository

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
)

/*

The event queue table is used for a multi-server broadcast scenario.

If we have two servers, one channel, which have users [1,2,3] and [4,5,6],
the event queue table holds the broadcast message which should be sent
to all these users in the channel;

The reading of the event queue table is triggered by pubsub.

- mostly, as the servers send out all the data, the contents of the
  event queue table can be discarded,
- the events queue table might eventually be not needed if we can
  solve everything on the level of pubsub, (@todo)
- when a client reloads the browser, the events queue table isn't
  read, everything should be in messages
- the event queue has a server id for messages originating from the
  websocket; the rest api should broadcast to all websocket connected
  clients, while the websocket API (currently), performs a local
  broadcast, triggering the event poll only on other servers

*/

type (
	EventsRepository interface {
		Pull(ctx context.Context) (*types.EventQueueItem, error)
		Push(ctx context.Context, item *types.EventQueueItem) error
	}

	events struct {
		pipe chan *types.EventQueueItem
	}
)

var eventsPipe chan *types.EventQueueItem

func Events() EventsRepository {
	if eventsPipe == nil {
		eventsPipe = make(chan *types.EventQueueItem, 512)
	}
	return &events{eventsPipe}
}

func (r *events) Pull(ctx context.Context) (*types.EventQueueItem, error) {
	select {
	case res, ok := <-r.pipe:
		if !ok {
			return res, ErrEventsPullClosed.New()
		}
		return res, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *events) Push(ctx context.Context, item *types.EventQueueItem) error {
	item.ID = factory.Sonyflake.NextID()
	select {
	case r.pipe <- item:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
