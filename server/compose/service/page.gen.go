package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	types "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
)

type page struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        pageAccessController
	services  *pageServices
}

func (svc *page) FindByID(ctx context.Context, namespaceID uint64, ID uint64) (res *types.Page, err error) {
	var (
		aProps = &pageActionProps{page: &types.Page{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
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
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
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
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
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

		if err = svc.guard(ctx, res); err != nil {
			return
		}

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

		if err = svc.guard(ctx, res); err != nil {
			return
		}

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

		if err = svc.guard(ctx, res); err != nil {
			return
		}

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
func (svc *page) guard(_ context.Context, _ *types.Page) error { return nil }

func (svc *page) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *page) Tree(ctx context.Context, namespaceID uint64) (tree types.PageSet, err error) {
	var (
		aProps = &pageActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		tree, err = svc.onTree(ctx, aProps, namespaceID)
		return err
	}()

	return tree, svc.recordAction(ctx, aProps, PageActionTree, err)
}

func (svc *page) Reorder(ctx context.Context, namespaceID uint64, parentID uint64, pageIDs []uint64) (err error) {
	var (
		aProps = &pageActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onReorder(ctx, aProps, namespaceID, parentID, pageIDs)
		return err
	}()

	return svc.recordAction(ctx, aProps, PageActionReorder, err)
}

func (svc *page) UpdateIcon(ctx context.Context, namespaceID uint64, pageID uint64, icon *types.PageConfigIcon) (out *types.PageConfigIcon, err error) {
	var (
		aProps = &pageActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		out, err = svc.onUpdateIcon(ctx, aProps, namespaceID, pageID, icon)
		return err
	}()

	return out, svc.recordAction(ctx, aProps, PageActionUpdateIcon, err)
}
