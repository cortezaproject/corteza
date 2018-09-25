package repository

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/types"
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
	Events interface {
		With(ctx context.Context, db *factory.DB) Events

		Pull(origin uint64) ([]*types.EventQueueItem, error)
		Push(eqi *types.EventQueueItem) error
		Sync(origin, id uint64) error
	}

	events struct {
		*repository
	}
)

func NewEvents(ctx context.Context, db *factory.DB) Events {
	return (&events{}).With(ctx, db)
}

func (r *events) With(ctx context.Context, db *factory.DB) Events {
	return &events{
		repository: r.repository.With(ctx, db),
	}
}

func (r *events) Pull(origin uint64) ([]*types.EventQueueItem, error) {
	var ee = make([]*types.EventQueueItem, 0)

	return ee, r.db().Quiet().Select(&ee, `
		SELECT * 
		  FROM event_queue 
         WHERE origin <> ? 
           AND id > GREATEST(COALESCE((SELECT rel_last FROM event_queue_synced WHERE origin = ?), 0), ?)
         LIMIT 50`, origin, origin, origin)
}

func (r *events) Push(eqi *types.EventQueueItem) error {
	eqi.ID = factory.Sonyflake.NextID()
	return r.db().Quiet().Insert("event_queue", eqi)
}

func (r *events) Sync(origin, id uint64) error {
	type evqs struct {
		Origin    uint64 `db:"origin"`
		LastEvent uint64 `db:"rel_last"`
	}

	// @todo do we even need this?
	return r.db().Quiet().Replace("event_queue_synced", evqs{
		Origin:    origin,
		LastEvent: id,
	})
}

func (r *events) Cleanup() error {
	return exec(r.db().Exec("DELETE FROM event_queue WHERE id < (SELECT MIN(rel_last) FROM event_queue_synced)"))
}

/*

do we need event_queue_synced??
do we need stable server id or can it be regenerad on each run?


*/
