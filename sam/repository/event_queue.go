package repository

import (
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
)

type (
	EventQueue interface {
		EventQueuePull(origin uint64) ([]*types.EventQueueItem, error)
		EventQueuePush(eqi *types.EventQueueItem) error
		EventQueueSync(origin, id uint64) error
	}
)

func (r *repository) EventQueuePull(origin uint64) ([]*types.EventQueueItem, error) {
	var ee = make([]*types.EventQueueItem, 0)

	return ee, r.db().Quiet().Select(&ee, `
		SELECT * 
		  FROM event_queue 
         WHERE origin <> ? 
           AND id > GREATEST(COALESCE((SELECT rel_last FROM event_queue_synced WHERE origin = ?), 0), ?)
         LIMIT 50`, origin, origin, origin)
}

func (r *repository) EventQueuePush(eqi *types.EventQueueItem) error {
	eqi.ID = factory.Sonyflake.NextID()
	return r.db().Quiet().Insert("event_queue", eqi)
}

func (r *repository) EventQueueSync(origin, id uint64) error {
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

func (r *repository) EventQueueCleanup() error {
	return exec(r.db().Exec("DELETE FROM event_queue WHERE id < (SELECT MIN(rel_last) FROM event_queue_synced)"))
}

/*

do we need event_queue_synced??
do we need stable server id or can it be regenerad on each run?


*/
