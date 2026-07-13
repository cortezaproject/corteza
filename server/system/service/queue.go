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

type (
	queueAccessController interface {
		CanCreateQueue(ctx context.Context) bool
		CanSearchQueues(ctx context.Context) bool
		CanReadQueue(ctx context.Context, c *types.Queue) bool
		CanUpdateQueue(ctx context.Context, c *types.Queue) bool
		CanDeleteQueue(ctx context.Context, c *types.Queue) bool
	}
)

func Queue() *queue {
	return &queue{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *queue) beforeCreate(_ context.Context, new *types.Queue) error {
	if !svc.isValidHandler(mt.ConsumerType(new.Consumer)) {
		return QueueErrInvalidConsumer()
	}
	return nil
}

func (svc *queue) afterCreate(_ context.Context, res *types.Queue) error {
	messagebus.Service().ReloadQueues()
	return nil
}

func (svc *queue) beforeUpdate(_ context.Context, upd *types.Queue, _ *types.Queue) error {
	if !svc.isValidHandler(mt.ConsumerType(upd.Consumer)) {
		return QueueErrInvalidConsumer()
	}
	return nil
}

func (svc *queue) afterUpdate(_ context.Context, _ *types.Queue) error {
	messagebus.Service().ReloadQueues()
	return nil
}

func (svc *queue) afterDelete(_ context.Context, _ *types.Queue) error {
	messagebus.Service().ReloadQueues()
	return nil
}

func (svc *queue) afterUndelete(_ context.Context, _ *types.Queue) error {
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
