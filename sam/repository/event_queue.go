package repository

import (
	"encoding/json"
	"github.com/titpetric/factory"
)

type (
	EventQueue interface {
		EventQueuePull(origin uint64) ([]*EventQueueItem, error)
		EventQueuePush(eqi *EventQueueItem) error
		EventQueueSync(origin uint64, ID uint64) error
	}

	EventQueueItem struct {
		ID         uint64          `db:"id"`
		Origin     uint64          `db:"origin"`
		Subscriber string          `db:"subscriber"`
		Payload    json.RawMessage `db:"payload"`
	}

	evqs struct {
		Origin    uint64 `db:"origin"`
		LastEvent uint64 `db:"rel_last"`
	}
)

func (r *repository) EventQueuePull(origin uint64) ([]*EventQueueItem, error) {
	var ee = make([]*EventQueueItem, 0)

	return ee, r.db().Quiet().Select(&ee, `
		SELECT * 
		  FROM event_queue 
         WHERE origin <> ? 
           AND id > GREATEST(COALESCE((SELECT rel_last FROM event_queue_synced WHERE origin = ?), 0), ?)
         LIMIT 50`, origin, origin, origin)
}

func (r *repository) EventQueuePush(eqi *EventQueueItem) error {
	eqi.ID = factory.Sonyflake.NextID()
	return r.db().Quiet().Insert("event_queue", eqi)
}

func (r *repository) EventQueueSync(origin uint64, ID uint64) error {
	// @todo do we even need this?
	return r.db().Quiet().Replace("event_queue_synced", evqs{
		Origin:    origin,
		LastEvent: ID,
	})
}

func (r *repository) EventQueueCleanup() error {
	return exec(r.db().Exec("DELETE FROM event_queue WHERE id < (SELECT MIN(rel_last) FROM event_queue_synced)"))
}

/*

do we need event_queue_synced??
do we need stable server id or can it be regenerad on each run?


*/
