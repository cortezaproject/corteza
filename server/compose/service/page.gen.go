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
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, PageActionUpdate, err)
}

func (svc *page) DeleteByID(ctx context.Context, namespaceID uint64, ID uint64, strategy types.PageChildrenDeleteStrategy) (err error) {
	var (
		aProps = &pageActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, namespaceID, ID, strategy, aProps)
	}()

	return svc.recordAction(ctx, aProps, PageActionDelete, err)
}

func (svc *page) UndeleteByID(ctx context.Context, namespaceID uint64, ID uint64) (err error) {
	var (
		aProps = &pageActionProps{}
	)

	err = func() (err error) {
		return svc.onUndelete(ctx, namespaceID, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, PageActionUndelete, err)
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
