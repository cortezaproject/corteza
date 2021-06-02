package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	queue struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        queueAccessController
	}

	queueAccessController interface {
		CanCreateQueue(ctx context.Context) bool
		CanReadQueue(ctx context.Context, c *messagebus.QueueSettings) bool
		CanUpdateQueue(ctx context.Context, c *messagebus.QueueSettings) bool
		CanDeleteQueue(ctx context.Context, c *messagebus.QueueSettings) bool
	}
)

func Queue() *queue {
	return &queue{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *queue) FindByID(ctx context.Context, ID uint64) (q *messagebus.QueueSettings, err error) {
	var (
		qProps = &queueActionProps{}
	)

	err = func() error {
		if ID == 0 {
			return QueueErrInvalidID()
		}

		if q, err = store.LookupMessagebusQueueSettingByID(ctx, svc.store, ID); err != nil {
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

func (svc *queue) Create(ctx context.Context, new *messagebus.QueueSettings) (q *messagebus.QueueSettings, err error) {
	var (
		qProps = &queueActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateQueue(ctx) {
			return QueueErrNotAllowedToCreate(qProps)
		}

		if !svc.isValidHandler(messagebus.ConsumerType(new.Consumer)) {
			return QueueErrInvalidConsumer(qProps)
		}

		// Set new values after beforeCreate events are emitted
		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateMessagebusQueueSetting(ctx, svc.store, new); err != nil {
			return
		}

		q = new

		// send the signal to reload all queues
		messagebus.Service().ReloadQueues()

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, QueueActionCreate, err)
}

func (svc *queue) Update(ctx context.Context, upd *messagebus.QueueSettings) (q *messagebus.QueueSettings, err error) {
	var (
		qProps = &queueActionProps{update: upd}
		qq     *messagebus.QueueSettings
		e      error
	)

	err = func() (err error) {
		if !svc.ac.CanUpdateQueue(ctx, upd) {
			return QueueErrNotAllowedToUpdate(qProps)
		}

		if qq, e = store.LookupMessagebusQueueSettingByID(ctx, svc.store, upd.ID); e != nil {
			return QueueErrNotFound(qProps)
		}

		if qq, e := store.LookupMessagebusQueueSettingByQueue(ctx, svc.store, upd.Queue); e == nil && qq != nil && qq.ID != upd.ID {
			return QueueErrAlreadyExists(qProps)
		}

		if !svc.isValidHandler(messagebus.ConsumerType(upd.Consumer)) {
			return QueueErrInvalidConsumer(qProps)
		}

		// Set new values after beforeCreate events are emitted
		upd.UpdatedAt = now()
		upd.CreatedAt = qq.CreatedAt

		if err = store.UpdateMessagebusQueueSetting(ctx, svc.store, upd); err != nil {
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
		q      *messagebus.QueueSettings
	)

	err = func() (err error) {
		if ID == 0 {
			return QueueErrInvalidID()
		}

		if q, err = store.LookupMessagebusQueueSettingByID(ctx, svc.store, ID); err != nil {
			return
		}

		qProps.setQueue(q)

		if !svc.ac.CanDeleteQueue(ctx, q) {
			return QueueErrNotAllowedToDelete(qProps)
		}

		q.DeletedAt = now()
		if err = store.UpdateMessagebusQueueSetting(ctx, svc.store, q); err != nil {
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
		q      *messagebus.QueueSettings
	)

	err = func() (err error) {
		if ID == 0 {
			return QueueErrInvalidID()
		}

		if q, err = store.LookupMessagebusQueueSettingByID(ctx, svc.store, ID); err != nil {
			return
		}

		qProps.setQueue(q)

		if !svc.ac.CanDeleteQueue(ctx, q) {
			return QueueErrNotAllowedToDelete(qProps)
		}

		q.DeletedAt = nil
		if err = store.UpdateMessagebusQueueSetting(ctx, svc.store, q); err != nil {
			return
		}

		// send the signal to reload all queues
		messagebus.Service().ReloadQueues()

		return nil
	}()

	return svc.recordAction(ctx, qProps, QueueActionDelete, err)
}

func (svc *queue) Search(ctx context.Context, filter messagebus.QueueSettingsFilter) (q messagebus.QueueSettingsSet, f messagebus.QueueSettingsFilter, err error) {
	var (
		aProps = &queueActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *messagebus.QueueSettings) (bool, error) {
		if !svc.ac.CanReadQueue(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if q, f, err = store.SearchMessagebusQueueSettings(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return q, f, svc.recordAction(ctx, aProps, QueueActionSearch, err)
}

func (svc *queue) isValidHandler(h messagebus.ConsumerType) bool {
	for _, hh := range messagebus.ConsumerTypes() {
		if h == hh {
			return true
		}
	}
	return false
}
