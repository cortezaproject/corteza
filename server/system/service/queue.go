package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/messagebus"
	mt "github.com/crusttech/human/server/pkg/messagebus/types"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service/event"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD skeleton (FindByID, Create, Update, DeleteByID, UndeleteByID and the
// loadQueue helper) is generated in queue.gen.go from system/queue.cue.
//
// The struct, queueAccessController interface and Queue() constructor are also
// generated (genConstructor + genAccessController), as is the standard Search
// (its action-log filter prop is named `search`, set via filterProp).
//
// This file owns the before/after hooks the generated CRUD calls into (consumer
// validation, unique-name check and the messagebus ReloadQueues() signal) and the
// custom methods (CreateQueueEvent, ProcessQueueMessage, CreateQueueMessage, the
// messagebus SearchQueues + makeFilter, isValidHandler).

// beforeCreate runs after the create access check, before id/timestamps are set.
func (svc *queue) beforeCreate(ctx context.Context, new *types.Queue) error {
	if !svc.isValidHandler(mt.ConsumerType(new.Consumer)) {
		return QueueErrInvalidConsumer(&queueActionProps{new: new})
	}

	return nil
}

// afterCreate runs after the store create.
func (svc *queue) afterCreate(ctx context.Context, res *types.Queue) error {
	// send the signal to reload all queues
	messagebus.Service().ReloadQueues()

	return nil
}

// beforeUpdate runs after the stale check, before the field copy.
func (svc *queue) beforeUpdate(ctx context.Context, upd, existing *types.Queue) error {
	if qq, e := store.LookupQueueByQueue(ctx, svc.store, upd.Queue); e == nil && qq != nil && qq.ID != upd.ID {
		return QueueErrAlreadyExists(&queueActionProps{update: upd})
	}

	if !svc.isValidHandler(mt.ConsumerType(upd.Consumer)) {
		return QueueErrInvalidConsumer(&queueActionProps{update: upd})
	}

	return nil
}

// afterUpdate runs after the store update.
func (svc *queue) afterUpdate(ctx context.Context, res *types.Queue) error {
	// send the signal to reload all queues
	messagebus.Service().ReloadQueues()

	return nil
}

// afterDelete runs after the store soft-delete.
func (svc *queue) afterDelete(ctx context.Context, res *types.Queue) error {
	// send the signal to reload all queues
	messagebus.Service().ReloadQueues()

	return nil
}

// afterUndelete runs after the store undelete.
func (svc *queue) afterUndelete(ctx context.Context, res *types.Queue) error {
	// send the signal to reload all queues
	messagebus.Service().ReloadQueues()

	return nil
}

func (svc *queue) CreateQueueEvent(q string, p []byte) eventbus.Event {
	return event.QueueOnMessage(&types.QueueMessage{
		Queue:   q,
		Payload: p,
	})
}

func (svc *queue) ProcessQueueMessage(ctx context.Context, ID uint64, m mt.QueueMessage) error {
	store.UpdateQueueMessage(ctx, svc.store, &types.QueueMessage{
		ID:        ID,
		Processed: now(),
		Queue:     m.Queue,
		Payload:   m.Payload,
	})

	return nil
}

func (svc *queue) CreateQueueMessage(ctx context.Context, m mt.QueueMessage) error {
	store.CreateQueueMessage(ctx, svc.store, &types.QueueMessage{
		ID:      nextID(),
		Created: now(),
		Queue:   m.Queue,
		Payload: m.Payload,
	})

	return nil
}

func (svc *queue) SearchQueues(ctx context.Context, ff mt.QueueFilter) (l []mt.QueueDb, f mt.QueueFilter, err error) {
	list, _, err := store.SearchQueues(ctx, svc.store, *(makeFilter(&ff)))

	if err != nil {
		return
	}

	l = make([]mt.QueueDb, len(list))

	for i, q := range list {
		l[i] = mt.QueueDb{
			Queue:    q.Queue,
			Consumer: q.Consumer,
			Meta:     mt.QueueMeta(q.Meta),
		}
	}

	return
}

func makeFilter(ff *mt.QueueFilter) (f *types.QueueFilter) {
	return &types.QueueFilter{
		Query:   ff.Query,
		Deleted: ff.Deleted,
		Sorting: ff.Sorting,
		Paging:  ff.Paging,
	}
}

func (svc *queue) isValidHandler(h mt.ConsumerType) bool {
	for _, hh := range mt.ConsumerTypes() {
		if h == hh {
			return true
		}
	}
	return false
}
