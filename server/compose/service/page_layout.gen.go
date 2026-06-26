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

func (svc *pageLayout) FindByID(ctx context.Context, namespaceID uint64, ID uint64) (res *types.PageLayout, err error) {
	var (
		aProps = &pageLayoutActionProps{pageLayout: &types.PageLayout{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, namespaceID, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, PageLayoutActionLookup, err)
}

func (svc *pageLayout) Search(ctx context.Context, filter types.PageLayoutFilter) (set types.PageLayoutSet, f types.PageLayoutFilter, err error) {
	var (
		aProps = &pageLayoutActionProps{filter: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, PageLayoutActionSearch, err)
}

func (svc *pageLayout) Create(ctx context.Context, new *types.PageLayout) (res *types.PageLayout, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &pageLayoutActionProps{pageLayout: new}
	)

	err = func() (err error) {
		before := func() error { return nil }
		after := func() error { return nil }
		res = new
		return svc.onCreate(ctx, new, before, after)
	}()

	return res, svc.recordAction(ctx, aProps, PageLayoutActionCreate, err)
}

func (svc *pageLayout) Update(ctx context.Context, upd *types.PageLayout) (res *types.PageLayout, err error) {
	var (
		aProps = &pageLayoutActionProps{changed: upd}
		old    *types.PageLayout
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadPageLayout(ctx, s, upd.ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setPageLayout(res)
		aProps.setChanged(res)
		old = res.Clone()

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return PageLayoutErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return PageLayoutErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Handle = upd.Handle
		res.PageID = upd.PageID
		res.ParentID = upd.ParentID
		res.NamespaceID = upd.NamespaceID
		res.Weight = upd.Weight
		res.OwnedBy = upd.OwnedBy
		res.CreatedByAgent = upd.CreatedByAgent
		res.UpdatedAt = now()

		if err = store.UpdateComposePageLayout(ctx, s, res); err != nil {
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

	return res, svc.recordAction(ctx, aProps, PageLayoutActionUpdate, err, old, res)
}

func (svc *pageLayout) DeleteByID(ctx context.Context, namespaceID uint64, pageID uint64, ID uint64) (err error) {
	var (
		aProps = &pageLayoutActionProps{}
		res    *types.PageLayout
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadPageLayout(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setPageLayout(res)

		return svc.onDelete(ctx, s, namespaceID, pageID, res, aProps)
	})

	return svc.recordAction(ctx, aProps, PageLayoutActionDelete, err)
}

func (svc *pageLayout) UndeleteByID(ctx context.Context, namespaceID uint64, pageID uint64, ID uint64) (err error) {
	var (
		aProps = &pageLayoutActionProps{}
		res    *types.PageLayout
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadPageLayout(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setPageLayout(res)

		return svc.onUndelete(ctx, s, namespaceID, pageID, res, aProps)
	})

	return svc.recordAction(ctx, aProps, PageLayoutActionUndelete, err)
}

func loadPageLayout(ctx context.Context, s store.ComposePageLayouts, ID uint64) (res *types.PageLayout, err error) {
	if ID == 0 {
		return nil, PageLayoutErrInvalidID()
	}

	if res, err = store.LookupComposePageLayoutByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, PageLayoutErrNotFound()
	}

	return
}

// toLabeledPageLayouts converts to []label.LabeledResource
func toLabeledPageLayouts(set []*types.PageLayout) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
