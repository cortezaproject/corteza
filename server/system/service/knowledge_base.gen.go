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
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type knowledgeBaseAccessController interface {
	CanCreateKnowledgeBase(context.Context) bool
	CanSearchKnowledgeBases(context.Context) bool
	CanReadKnowledgeBase(context.Context, *types.KnowledgeBase) bool
	CanUpdateKnowledgeBase(context.Context, *types.KnowledgeBase) bool
	CanDeleteKnowledgeBase(context.Context, *types.KnowledgeBase) bool
}

type knowledgeBase struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        knowledgeBaseAccessController
}

func KnowledgeBase() *knowledgeBase {
	return &knowledgeBase{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

func (svc *knowledgeBase) FindByID(ctx context.Context, ID uint64) (res *types.KnowledgeBase, err error) {
	var (
		aProps = &knowledgeBaseActionProps{knowledgeBase: &types.KnowledgeBase{ID: ID}}
	)

	err = func() error {
		if res, err = loadKnowledgeBase(ctx, svc.store, ID); err != nil {
			return KnowledgeBaseErrInvalidID().Wrap(err)
		}

		aProps.setKnowledgeBase(res)

		if !svc.ac.CanReadKnowledgeBase(ctx, res) {
			return KnowledgeBaseErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, KnowledgeBaseActionLookup, err)
}

func (svc *knowledgeBase) Search(ctx context.Context, filter types.KnowledgeBaseFilter) (set types.KnowledgeBaseSet, f types.KnowledgeBaseFilter, err error) {
	var (
		aProps = &knowledgeBaseActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.KnowledgeBase) (bool, error) {
		if !svc.ac.CanReadKnowledgeBase(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchKnowledgeBases(ctx) {
			return KnowledgeBaseErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchKnowledgeBases(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, KnowledgeBaseActionSearch, err)
}

func (svc *knowledgeBase) Create(ctx context.Context, new *types.KnowledgeBase) (res *types.KnowledgeBase, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &knowledgeBaseActionProps{knowledgeBase: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateKnowledgeBase(ctx) {
			return KnowledgeBaseErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateKnowledgeBase(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, KnowledgeBaseActionCreate, err)
}

func (svc *knowledgeBase) Update(ctx context.Context, upd *types.KnowledgeBase) (res *types.KnowledgeBase, err error) {
	var (
		aProps = &knowledgeBaseActionProps{update: upd}
		old    *types.KnowledgeBase
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadKnowledgeBase(ctx, s, upd.ID); err != nil {
			return
		}

		aProps.setKnowledgeBase(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return KnowledgeBaseErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return KnowledgeBaseErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		return svc.onUpdate(ctx, s, upd, res, aProps, before, after)
	})

	return res, svc.recordAction(ctx, aProps, KnowledgeBaseActionUpdate, err, old, res)
}

func (svc *knowledgeBase) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &knowledgeBaseActionProps{}
		res    *types.KnowledgeBase
	)
	err = func() (err error) {
		if res, err = loadKnowledgeBase(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setKnowledgeBase(res)

		if !svc.ac.CanDeleteKnowledgeBase(ctx, res) {
			return KnowledgeBaseErrNotAllowedToDelete()
		}

		if err = svc.beforeDelete(ctx, res); err != nil {
			return err
		}

		res.DeletedAt = now()
		if err = store.UpdateKnowledgeBase(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, KnowledgeBaseActionDelete, err)
}

func (svc *knowledgeBase) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &knowledgeBaseActionProps{}
		res    *types.KnowledgeBase
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadKnowledgeBase(ctx, s, ID); err != nil {
			return
		}

		aProps.setKnowledgeBase(res)

		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, KnowledgeBaseActionUndelete, err)
}

func loadKnowledgeBase(ctx context.Context, s store.KnowledgeBases, ID uint64) (res *types.KnowledgeBase, err error) {
	if ID == 0 {
		return nil, KnowledgeBaseErrInvalidID()
	}

	if res, err = store.LookupKnowledgeBaseByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, KnowledgeBaseErrNotFound()
	}

	return
}
