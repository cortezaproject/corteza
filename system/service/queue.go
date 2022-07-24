package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	mt "github.com/cortezaproject/corteza-server/pkg/messagebus/types"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	queue struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        queueAccessController
	}

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

func (svc *queue) CreateQueueEvent(q string, p []byte) eventbus.Event {
	return event.QueueOnMessage(&types.QueueMessage{
		Queue:   q,
		Payload: p,
	})
}

func (svc *queue) ProcessQueueMessage(ctx context.Context, ID uint64, m mt.QueueMessage) error {
	svc.store.UpdateQueueMessage(ctx, &types.QueueMessage{
		ID:        ID,
		Processed: now(),
		Queue:     m.Queue,
		Payload:   m.Payload,
	})

	return nil
}

func (svc *queue) CreateQueueMessage(ctx context.Context, m mt.QueueMessage) error {
	svc.store.CreateQueueMessage(ctx, &types.QueueMessage{
		ID:      nextID(),
		Created: now(),
		Queue:   m.Queue,
		Payload: m.Payload,
	})

	return nil
}

func (svc *queue) SearchQueues(ctx context.Context, ff mt.QueueFilter) (l []mt.QueueDb, f mt.QueueFilter, err error) {
	list, _, err := svc.store.SearchQueues(ctx, *(makeFilter(&ff)))

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

func (svc *queue) FindByID(ctx context.Context, ID uint64) (q *types.Queue, err error) {
	var (
		qProps = &queueActionProps{}
	)

	err = func() error {
		if q, err = loadQueue(ctx, svc.store, ID); err != nil {
			return TemplateErrInvalidID().Wrap(err)
		}

		qProps.setQueue(q)

		if !svc.ac.CanReadQueue(ctx, q) {
			return QueueErrNotAllowedToRead(qProps)
		}

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, QueueActionLookup, err)
}

func (svc *queue) Create(ctx context.Context, new *types.Queue) (q *types.Queue, err error) {
	var (
		qProps = &queueActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateQueue(ctx) {
			return QueueErrNotAllowedToCreate(qProps)
		}

		if !svc.isValidHandler(mt.ConsumerType(new.Consumer)) {
			return QueueErrInvalidConsumer(qProps)
		}

		// Set new values after beforeCreate events are emitted
		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateQueue(ctx, svc.store, new); err != nil {
			return
		}

		q = new

		// send the signal to reload all queues
		messagebus.Service().ReloadQueues()

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, QueueActionCreate, err)
}

func (svc *queue) Update(ctx context.Context, upd *types.Queue) (q *types.Queue, err error) {
	var (
		qProps = &queueActionProps{update: upd}
		qq     *types.Queue
		e      error
	)

	err = func() (err error) {
		if !svc.ac.CanUpdateQueue(ctx, upd) {
			return QueueErrNotAllowedToUpdate(qProps)
		}

		if qq, e = store.LookupQueueByID(ctx, svc.store, upd.ID); e != nil {
			return QueueErrNotFound(qProps)
		}

		if qq, e := store.LookupQueueByQueue(ctx, svc.store, upd.Queue); e == nil && qq != nil && qq.ID != upd.ID {
			return QueueErrAlreadyExists(qProps)
		}

		if !svc.isValidHandler(mt.ConsumerType(upd.Consumer)) {
			return QueueErrInvalidConsumer(qProps)
		}

		// Set new values after beforeCreate events are emitted
		upd.UpdatedAt = now()
		upd.CreatedAt = qq.CreatedAt

		if err = store.UpdateQueue(ctx, svc.store, upd); err != nil {
			return
		}

		q = upd

		// send the signal to reload all queues
		messagebus.Service().ReloadQueues()

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, QueueActionUpdate, err)
}

func (svc *queue) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &queueActionProps{}
		q      *types.Queue
	)

	err = func() (err error) {
		if q, err = loadQueue(ctx, svc.store, ID); err != nil {
			return
		}

		qProps.setQueue(q)

		if !svc.ac.CanDeleteQueue(ctx, q) {
			return QueueErrNotAllowedToDelete(qProps)
		}

		q.DeletedAt = now()
		if err = store.UpdateQueue(ctx, svc.store, q); err != nil {
			return
		}

		// send the signal to reload all queues
		messagebus.Service().ReloadQueues()

		return nil
	}()

	return svc.recordAction(ctx, qProps, QueueActionDelete, err)
}

func (svc *queue) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &queueActionProps{}
		q      *types.Queue
	)

	err = func() (err error) {
		if q, err = loadQueue(ctx, svc.store, ID); err != nil {
			return
		}

		qProps.setQueue(q)

		if !svc.ac.CanDeleteQueue(ctx, q) {
			return QueueErrNotAllowedToDelete(qProps)
		}

		q.DeletedAt = nil
		if err = store.UpdateQueue(ctx, svc.store, q); err != nil {
			return
		}

		// send the signal to reload all queues
		messagebus.Service().ReloadQueues()

		return nil
	}()

	return svc.recordAction(ctx, qProps, QueueActionDelete, err)
}

func (svc *queue) Search(ctx context.Context, filter types.QueueFilter) (q types.QueueSet, f types.QueueFilter, err error) {
	var (
		aProps = &queueActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Queue) (bool, error) {
		if !svc.ac.CanReadQueue(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchQueues(ctx) {
			return QueueErrNotAllowedToSearch()
		}

		if q, f, err = store.SearchQueues(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return q, f, svc.recordAction(ctx, aProps, QueueActionSearch, err)
}

func loadQueue(ctx context.Context, s store.Queues, ID uint64) (res *types.Queue, err error) {
	if ID == 0 {
		return nil, QueueErrInvalidID()
	}

	if res, err = store.LookupQueueByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, QueueErrNotFound()
	}

	return
}

func (svc *queue) isValidHandler(h mt.ConsumerType) bool {
	for _, hh := range mt.ConsumerTypes() {
		if h == hh {
			return true
		}
	}
	return false
}
