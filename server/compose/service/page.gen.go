package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	types "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
)

func (svc *page) FindByID(ctx context.Context, namespaceID uint64, ID uint64) (res *types.Page, err error) {
	var (
		aProps = &pageActionProps{page: &types.Page{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, namespaceID, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, PageActionLookup, err)
}

func (svc *page) Search(ctx context.Context, filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error) {
	var (
		aProps = &pageActionProps{filter: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, PageActionSearch, err)
}

func (svc *page) Create(ctx context.Context, new *types.Page) (res *types.Page, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &pageActionProps{page: new}
	)

	err = func() (err error) {
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, PageActionCreate, err)
}

func (svc *page) Update(ctx context.Context, upd *types.Page) (res *types.Page, err error) {
	var (
		aProps = &pageActionProps{changed: upd}
		old    *types.Page
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadPage(ctx, s, upd.ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setPage(res)
		aProps.setChanged(res)
		old = res.Clone()

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return PageErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return PageErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Title = upd.Title
		res.Handle = upd.Handle
		res.SelfID = upd.SelfID
		res.ModuleID = upd.ModuleID
		res.NamespaceID = upd.NamespaceID
		res.Visible = upd.Visible
		res.Weight = upd.Weight
		res.Description = upd.Description
		res.CreatedByAgent = upd.CreatedByAgent
		res.UpdatedAt = now()

		if err = store.UpdateComposePage(ctx, s, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, s, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, PageActionUpdate, err, old, res)
}

func (svc *page) DeleteByID(ctx context.Context, namespaceID uint64, ID uint64, strategy types.PageChildrenDeleteStrategy) (err error) {
	var (
		aProps = &pageActionProps{}
		res    *types.Page
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadPage(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setPage(res)

		return svc.onDelete(ctx, s, namespaceID, res, strategy, aProps)
	})

	return svc.recordAction(ctx, aProps, PageActionDelete, err)
}

func (svc *page) UndeleteByID(ctx context.Context, namespaceID uint64, ID uint64) (err error) {
	var (
		aProps = &pageActionProps{}
		res    *types.Page
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadPage(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setPage(res)

		return svc.onUndelete(ctx, s, namespaceID, res, aProps)
	})

	return svc.recordAction(ctx, aProps, PageActionUndelete, err)
}

func loadPage(ctx context.Context, s store.ComposePages, ID uint64) (res *types.Page, err error) {
	if ID == 0 {
		return nil, PageErrInvalidID()
	}

	if res, err = store.LookupComposePageByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, PageErrNotFound()
	}

	return
}

// toLabeledPages converts to []label.LabeledResource
func toLabeledPages(set []*types.Page) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
