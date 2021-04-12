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
		ac        templateAccessController
	}

	QueueService interface {
		FindByID(ctx context.Context, ID uint64) (*messagebus.QueueSettings, error)
		Search(context.Context, messagebus.QueueSettingsFilter) (messagebus.QueueSettingsSet, messagebus.QueueSettingsFilter, error)

		Create(ctx context.Context, q *messagebus.QueueSettings) (*messagebus.QueueSettings, error)
		Update(ctx context.Context, q *messagebus.QueueSettings) (*messagebus.QueueSettings, error)

		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
	}
)

func Queue() QueueService {
	return (&queue{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	})
}

func (svc *queue) FindByID(ctx context.Context, ID uint64) (q *messagebus.QueueSettings, err error) {
	var (
		qProps = &queueActionProps{}
	)

	err = func() error {
		if ID == 0 {
			return QueueErrInvalidID()
		}

		if q, err = store.LookupMessagebusQueuesettingByID(ctx, svc.store, ID); err != nil {
			return TemplateErrInvalidID().Wrap(err)
		}

		qProps.setQueue(q)

		// if !svc.ac.CanReadTemplate(ctx, tpl) {
		// 	return TemplateErrNotAllowedToRead()
		// }

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, QueueActionLookup, err)
}

func (svc *queue) Create(ctx context.Context, new *messagebus.QueueSettings) (q *messagebus.QueueSettings, err error) {
	var (
		qProps = &queueActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.isValidHandler(messagebus.HandlerType(new.Handler)) {
			return QueueErrInvalidHandler(qProps)
		}

		// Set new values after beforeCreate events are emitted
		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateMessagebusQueuesetting(ctx, svc.store, new); err != nil {
			return
		}

		q = new

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, QueueActionCreate, err)
}

func (svc *queue) Update(ctx context.Context, upd *messagebus.QueueSettings) (q *messagebus.QueueSettings, err error) {
	var (
		qProps = &queueActionProps{update: upd}
	)

	err = func() (err error) {
		if _, e := store.LookupMessagebusQueuesettingByID(ctx, svc.store, upd.ID); e != nil {
			return QueueErrNotFound(qProps)
		}

		if qq, e := store.LookupMessagebusQueuesettingByQueue(ctx, svc.store, upd.Queue); e == nil && qq != nil && qq.ID != upd.ID {
			return QueueErrAlreadyExists(qProps)
		}

		if !svc.isValidHandler(messagebus.HandlerType(upd.Handler)) {
			return QueueErrInvalidHandler(qProps)
		}

		// Set new values after beforeCreate events are emitted
		upd.UpdatedAt = now()

		if err = store.UpdateMessagebusQueuesetting(ctx, svc.store, upd); err != nil {
			return
		}

		q = upd

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

		if q, err = store.LookupMessagebusQueuesettingByID(ctx, svc.store, ID); err != nil {
			return
		}

		qProps.setQueue(q)

		// if !svc.ac.CanDeleteTemplate(ctx, tpl) {
		// 	return TemplateErrNotAllowedToDelete()
		// }

		q.DeletedAt = now()
		if err = store.UpdateMessagebusQueuesetting(ctx, svc.store, q); err != nil {
			return
		}

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

		if q, err = store.LookupMessagebusQueuesettingByID(ctx, svc.store, ID); err != nil {
			return
		}

		qProps.setQueue(q)

		// if !svc.ac.CanDeleteTemplate(ctx, tpl) {
		// 	return TemplateErrNotAllowedToDelete()
		// }

		q.DeletedAt = nil
		if err = store.UpdateMessagebusQueuesetting(ctx, svc.store, q); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, qProps, QueueActionDelete, err)
}

func (svc *queue) Search(ctx context.Context, filter messagebus.QueueSettingsFilter) (q messagebus.QueueSettingsSet, f messagebus.QueueSettingsFilter, err error) {
	var (
		aProps = &queueActionProps{search: &filter}
	)

	err = func() error {
		if q, f, err = store.SearchMessagebusQueuesettings(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return q, f, svc.recordAction(ctx, aProps, QueueActionSearch, err)
}

func (svc *queue) isValidHandler(h messagebus.HandlerType) bool {
	for _, hh := range messagebus.HandlerTypes() {
		if h == hh {
			return true
		}
	}
	return false
}
