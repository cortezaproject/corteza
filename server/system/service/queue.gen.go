package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type queueAccessController interface {
	CanCreateQueue(context.Context) bool
	CanSearchQueues(context.Context) bool
	CanReadQueue(context.Context, *types.Queue) bool
	CanUpdateQueue(context.Context, *types.Queue) bool
	CanDeleteQueue(context.Context, *types.Queue) bool
}

type queue struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        queueAccessController
}

func Queue() *queue {
	return &queue{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

func (svc *queue) FindByID(ctx context.Context, ID uint64) (res *types.Queue, err error) {
	var (
		aProps = &queueActionProps{queue: &types.Queue{ID: ID}}
	)

	err = func() error {
		if res, err = loadQueue(ctx, svc.store, ID); err != nil {
			return QueueErrInvalidID().Wrap(err)
		}

		aProps.setQueue(res)

		if !svc.ac.CanReadQueue(ctx, res) {
			return QueueErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, QueueActionLookup, err)
}

func (svc *queue) Search(ctx context.Context, filter types.QueueFilter) (set types.QueueSet, f types.QueueFilter, err error) {
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

		if set, f, err = store.SearchQueues(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, QueueActionSearch, err)
}

func (svc *queue) Create(ctx context.Context, new *types.Queue) (res *types.Queue, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &queueActionProps{queue: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateQueue(ctx) {
			return QueueErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateQueue(ctx, svc.store, new); err != nil {
			return
		}

		res = new

		if err = svc.afterCreate(ctx, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, QueueActionCreate, err)
}

func (svc *queue) Update(ctx context.Context, upd *types.Queue) (res *types.Queue, err error) {
	var (
		aProps = &queueActionProps{update: upd}
		old    *types.Queue
	)
	err = func() (err error) {
		if res, err = loadQueue(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.setQueue(res)

		if !svc.ac.CanUpdateQueue(ctx, res) {
			return QueueErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return QueueErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.Consumer = upd.Consumer
		res.Queue = upd.Queue
		res.Meta = upd.Meta
		res.UpdatedAt = now()

		if err = store.UpdateQueue(ctx, svc.store, res); err != nil {
			return err
		}

		if err = svc.afterUpdate(ctx, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, QueueActionUpdate, err, old, res)
}

func (svc *queue) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &queueActionProps{}
		res    *types.Queue
	)
	err = func() (err error) {
		if res, err = loadQueue(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setQueue(res)

		if !svc.ac.CanDeleteQueue(ctx, res) {
			return QueueErrNotAllowedToDelete()
		}

		res.DeletedAt = now()
		if err = store.UpdateQueue(ctx, svc.store, res); err != nil {
			return
		}

		if err = svc.afterDelete(ctx, res); err != nil {
			return err
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, QueueActionDelete, err)
}

func (svc *queue) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &queueActionProps{}
		res    *types.Queue
	)
	err = func() (err error) {
		if res, err = loadQueue(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setQueue(res)

		if !svc.ac.CanDeleteQueue(ctx, res) {
			return QueueErrNotAllowedToUndelete()
		}

		res.DeletedAt = nil
		if err = store.UpdateQueue(ctx, svc.store, res); err != nil {
			return
		}

		if err = svc.afterUndelete(ctx, res); err != nil {
			return err
		}

		return nil
	}()

	return svc.recordAction(ctx, aProps, QueueActionUndelete, err)
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
