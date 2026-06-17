package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

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
	)

	err = func() (err error) {
		if res, err = loadQueue(ctx, svc.store, upd.ID); err != nil {
			return
		}

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

	return res, svc.recordAction(ctx, aProps, QueueActionUpdate, err)
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
