package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	types "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/label"
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
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, PageLayoutActionCreate, err)
}

func (svc *pageLayout) Update(ctx context.Context, upd *types.PageLayout) (res *types.PageLayout, err error) {
	var (
		aProps = &pageLayoutActionProps{changed: upd}
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, PageLayoutActionUpdate, err)
}

func (svc *pageLayout) DeleteByID(ctx context.Context, namespaceID uint64, pageID uint64, ID uint64) (err error) {
	var (
		aProps = &pageLayoutActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, namespaceID, pageID, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, PageLayoutActionDelete, err)
}

func (svc *pageLayout) UndeleteByID(ctx context.Context, namespaceID uint64, pageID uint64, ID uint64) (err error) {
	var (
		aProps = &pageLayoutActionProps{}
	)

	err = func() (err error) {
		return svc.onUndelete(ctx, namespaceID, pageID, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, PageLayoutActionUndelete, err)
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
