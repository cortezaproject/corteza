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

type resourceTranslation struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        resourceTranslationAccessController
}

func ResourceTranslation() *resourceTranslation {
	return &resourceTranslation{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

func (svc *resourceTranslation) FindByID(ctx context.Context, ID uint64) (res *types.ResourceTranslation, err error) {
	var (
		aProps = &resourceTranslationActionProps{resourceTranslation: &types.ResourceTranslation{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ResourceTranslationActionLookup, err)
}

func (svc *resourceTranslation) Search(ctx context.Context, filter types.ResourceTranslationFilter) (set types.ResourceTranslationSet, f types.ResourceTranslationFilter, err error) {
	var (
		aProps = &resourceTranslationActionProps{filter: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, ResourceTranslationActionSearch, err)
}

func (svc *resourceTranslation) Create(ctx context.Context, new *types.ResourceTranslation) (res *types.ResourceTranslation, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &resourceTranslationActionProps{resourceTranslation: new, new: new}
	)

	err = func() (err error) {
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, ResourceTranslationActionCreate, err)
}

func (svc *resourceTranslation) Update(ctx context.Context, upd *types.ResourceTranslation) (res *types.ResourceTranslation, err error) {
	var (
		aProps = &resourceTranslationActionProps{update: upd}
		old    *types.ResourceTranslation
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadResourceTranslation(ctx, s, upd.ID); err != nil {
			return
		}

		aProps.setResourceTranslation(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ResourceTranslationErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		return svc.onUpdate(ctx, s, upd, res, aProps, before, after)
	})

	return res, svc.recordAction(ctx, aProps, ResourceTranslationActionUpdate, err, old, res)
}

func (svc *resourceTranslation) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &resourceTranslationActionProps{}
		res    *types.ResourceTranslation
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadResourceTranslation(ctx, s, ID); err != nil {
			return
		}

		aProps.setResourceTranslation(res)

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ResourceTranslationActionDelete, err)
}

func (svc *resourceTranslation) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &resourceTranslationActionProps{}
		res    *types.ResourceTranslation
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadResourceTranslation(ctx, s, ID); err != nil {
			return
		}

		aProps.setResourceTranslation(res)

		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ResourceTranslationActionUndelete, err)
}

func loadResourceTranslation(ctx context.Context, s store.ResourceTranslations, ID uint64) (res *types.ResourceTranslation, err error) {
	if ID == 0 {
		return nil, ResourceTranslationErrInvalidID()
	}

	if res, err = store.LookupResourceTranslationByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ResourceTranslationErrNotFound()
	}

	return
}
