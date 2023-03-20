package service

import (
	"context"
	"reflect"

	"github.com/cortezaproject/corteza/server/compose/service/event"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/handle"
	"github.com/cortezaproject/corteza/server/pkg/label"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/store"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
)

type (
	pageLayout struct {
		actionlog actionlog.Recorder
		ac        pageLayoutAccessController
		eventbus  eventDispatcher
		store     store.Storer
		locale    ResourceTranslationsManagerService

		pageLayoutSettings *pageLayoutSettings
	}

	pageLayoutSettings struct {
		hideNew    bool
		hideEdit   bool
		hideSubmit bool
		hideDelete bool
		hideClone  bool
		hideBack   bool
	}

	pageLayoutAccessController interface {
		// @todo...
	}

	pageLayoutUpdateHandler func(ctx context.Context, ns *types.Namespace, c *types.PageLayout) (pageLayoutChanges, error)
	pageLayoutChanges       uint8
)

const (
	pageLayoutUnchanged     pageLayoutChanges = 0
	pageLayoutChanged       pageLayoutChanges = 1
	pageLayoutLabelsChanged pageLayoutChanges = 2
)

func PageLayout() *pageLayout {
	return &pageLayout{
		actionlog:          DefaultActionlog,
		ac:                 DefaultAccessControl,
		eventbus:           eventbus.Service(),
		store:              DefaultStore,
		locale:             DefaultResourceTranslation,
		pageLayoutSettings: &pageLayoutSettings{},
	}
}

func (svc pageLayout) FindByID(ctx context.Context, namespaceID, pageLayoutID uint64) (p *types.PageLayout, err error) {
	return svc.lookup(ctx, namespaceID, func(aProps *pageLayoutActionProps) (*types.PageLayout, error) {
		if pageLayoutID == 0 {
			return nil, PageLayoutErrInvalidID()
		}

		aProps.pageLayout.ID = pageLayoutID
		return store.LookupComposePageLayoutByID(ctx, svc.store, pageLayoutID)
	})
}

func (svc pageLayout) FindByHandle(ctx context.Context, namespaceID uint64, h string) (c *types.PageLayout, err error) {
	return svc.lookup(ctx, namespaceID, func(aProps *pageLayoutActionProps) (*types.PageLayout, error) {
		if !handle.IsValid(h) {
			return nil, PageLayoutErrInvalidHandle()
		}

		aProps.pageLayout.Handle = h
		return store.LookupComposePageLayoutByNamespaceIDHandle(ctx, svc.store, namespaceID, h)
	})
}

func (svc pageLayout) FindByPageLayoutID(ctx context.Context, namespaceID, pageLayoutID uint64) (p *types.PageLayout, err error) {
	return svc.lookup(ctx, namespaceID, func(aProps *pageLayoutActionProps) (*types.PageLayout, error) {
		if pageLayoutID == 0 {
			return nil, PageLayoutErrInvalidID()
		}

		aProps.pageLayout.ID = pageLayoutID
		return store.LookupComposePageLayoutByID(ctx, svc.store, pageLayoutID)
	})
}

func checkPageLayout(ctx context.Context, ac pageLayoutAccessController) func(res *types.PageLayout) (bool, error) {
	return func(res *types.PageLayout) (bool, error) {
		// if !ac.CanReadPageLayout(ctx, res) {
		// 	return false, nil
		// }

		return true, nil
	}
}

// search fn() orchestrates pageLayouts search, namespace preload and check
func (svc pageLayout) search(ctx context.Context, filter types.PageLayoutFilter) (set types.PageLayoutSet, f types.PageLayoutFilter, err error) {
	var (
		aProps = &pageLayoutActionProps{filter: &filter}
		ns     *types.Namespace
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = checkPageLayout(ctx, svc.ac)

	err = func() error {
		ns, err = loadNamespace(ctx, svc.store, filter.NamespaceID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)
		// if !svc.ac.CanSearchPageLayoutsOnNamespace(ctx, ns) {
		// 	return PageLayoutErrNotAllowedToSearch()
		// }

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.PageLayout{}.LabelResourceKind(),
				filter.Labels,
			)

			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if set, f, err = store.SearchComposePageLayouts(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledPageLayouts(set)...); err != nil {
			return err
		}

		// i18n
		tag := locale.GetAcceptLanguageFromContext(ctx)
		set.Walk(func(p *types.PageLayout) error {
			p.DecodeTranslations(svc.locale.Locale().ResourceTranslations(tag, p.ResourceTranslation()))
			return nil
		})

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, PageLayoutActionSearch, err)
}

func (svc pageLayout) Find(ctx context.Context, filter types.PageLayoutFilter) (set types.PageLayoutSet, f types.PageLayoutFilter, err error) {
	return svc.search(ctx, filter)
}

func (svc pageLayout) Create(ctx context.Context, new *types.PageLayout) (*types.PageLayout, error) {
	var (
		ns     *types.Namespace
		aProps = &pageLayoutActionProps{pageLayout: new}
	)

	new.ID = 0

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !handle.IsValid(new.Handle) {
			return PageLayoutErrInvalidID()
		}

		if ns, err = loadNamespace(ctx, s, new.NamespaceID); err != nil {
			return err
		}

		// if !svc.ac.CanCreatePageLayoutOnNamespace(ctx, ns) {
		// 	return PageLayoutErrNotAllowedToCreate()
		// }

		aProps.setNamespace(ns)

		if err = svc.eventbus.WaitFor(ctx, event.PageLayoutBeforeCreate(new, nil, ns, nil)); err != nil {
			return err
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.UpdatedAt = nil
		new.DeletedAt = nil

		// Ensure pageLayout-block IDs
		for i := range new.Blocks {
			new.Blocks[i].BlockID = uint64(i) + 1
		}

		aProps.setChanged(new)

		if err = store.CreateComposePageLayout(ctx, s, new); err != nil {
			return err
		}

		// if err = updateTranslations(ctx, svc.ac, svc.locale, new.EncodeTranslations()...); err != nil {
		// 	return
		// }

		if err = label.Create(ctx, s, new); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.PageLayoutAfterCreate(new, nil, ns, nil))
		return err
	})

	return new, svc.recordAction(ctx, aProps, PageLayoutActionCreate, err)
}

func (svc pageLayout) Update(ctx context.Context, upd *types.PageLayout) (c *types.PageLayout, err error) {
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, res, err := loadPageLayoutCombo(ctx, s, upd.NamespaceID, upd.PageID, upd.ID)
		if err != nil {
			return
		}

		c, err = svc.updater(ctx, svc.store, ns, res, PageLayoutActionUpdate, svc.handleUpdate(ctx, upd))
		return
	})

	return
}

func (svc pageLayout) DeleteByID(ctx context.Context, namespaceID, pageID, pageLayoutID uint64) error {
	var (
		ns  *types.Namespace
		res *types.PageLayout
	)

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		// simply delete the pageLayout and ignore the subpageLayouts
		ns, res, err = loadPageLayoutCombo(ctx, s, namespaceID, pageID, pageLayoutID)
		if err != nil {
			return
		}

		_, err = svc.updater(ctx, svc.store, ns, res, PageLayoutActionDelete, svc.handleDelete)
		return
	})
}

func (svc pageLayout) UndeleteByID(ctx context.Context, namespaceID, pageID, pageLayoutID uint64) error {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, res, err := loadPageLayoutCombo(ctx, s, namespaceID, pageID, pageLayoutID)
		if err != nil {
			return
		}

		_, err = svc.updater(ctx, svc.store, ns, res, PageLayoutActionUpdate, svc.handleUndelete)
		return
	})
}

func (svc pageLayout) updater(ctx context.Context, s store.Storer, ns *types.Namespace, res *types.PageLayout, action func(...*pageLayoutActionProps) *pageLayoutAction, fn pageLayoutUpdateHandler) (*types.PageLayout, error) {
	var (
		changes pageLayoutChanges
		old     *types.PageLayout
		aProps  = &pageLayoutActionProps{pageLayout: res}
		err     error
	)

	err = store.Tx(ctx, s, func(ctx context.Context, s store.Storer) (err error) {
		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		// Get max blockID for later use
		blockID := uint64(0)
		for _, b := range res.Blocks {
			if b.BlockID > blockID {
				blockID = b.BlockID
			}
		}

		old = res.Clone()

		aProps.setNamespace(ns)
		aProps.setChanged(res)

		if res.DeletedAt == nil {
			err = svc.eventbus.WaitFor(ctx, event.PageLayoutBeforeUpdate(old, res, ns, nil))
		} else {
			err = svc.eventbus.WaitFor(ctx, event.PageLayoutBeforeDelete(old, res, ns, nil))
		}

		if err != nil {
			return
		}

		if changes, err = fn(ctx, ns, res); err != nil {
			return err
		}

		if changes&pageLayoutChanged > 0 {
			if err = store.UpdateComposePageLayout(ctx, s, res); err != nil {
				return err
			}
		}

		// if err = updateTranslations(ctx, svc.ac, svc.locale, res.EncodeTranslations()...); err != nil {
		// 	return
		// }

		if changes&pageLayoutLabelsChanged > 0 {
			if err = label.Update(ctx, s, res); err != nil {
				return
			}
		}

		if res.DeletedAt == nil {
			err = svc.eventbus.WaitFor(ctx, event.PageLayoutAfterUpdate(res, res, ns, nil))
		} else {
			err = svc.eventbus.WaitFor(ctx, event.PageLayoutAfterDelete(nil, res, ns, nil))
		}

		return err
	})

	return res, svc.recordAction(ctx, aProps, action, err)
}

// lookup fn() orchestrates pageLayout lookup, namespace preload and check
func (svc pageLayout) lookup(ctx context.Context, namespaceID uint64, lookup func(*pageLayoutActionProps) (*types.PageLayout, error)) (p *types.PageLayout, err error) {
	var aProps = &pageLayoutActionProps{pageLayout: &types.PageLayout{NamespaceID: namespaceID}}

	err = func() error {
		if ns, err := loadNamespace(ctx, svc.store, namespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if p, err = lookup(aProps); errors.IsNotFound(err) {
			return PageLayoutErrNotFound()
		} else if err != nil {
			return err
		}

		p.DecodeTranslations(svc.locale.Locale().ResourceTranslations(locale.GetAcceptLanguageFromContext(ctx), p.ResourceTranslation()))

		aProps.setPageLayout(p)

		// if !svc.ac.CanReadPageLayout(ctx, p) {
		// 	return PageLayoutErrNotAllowedToRead()
		// }

		if err = label.Load(ctx, svc.store, p); err != nil {
			return err
		}

		return nil
	}()

	return p, svc.recordAction(ctx, aProps, PageLayoutActionLookup, err)
}

func (svc pageLayout) uniqueCheck(ctx context.Context, p *types.PageLayout) (err error) {
	if p.Handle != "" {
		if e, _ := store.LookupComposePageLayoutByNamespaceIDHandle(ctx, svc.store, p.NamespaceID, p.Handle); e != nil && e.ID != p.ID {
			return PageLayoutErrHandleNotUnique()
		}
	}

	if p.ModuleID > 0 {
		if e, _ := store.LookupComposePageLayoutByNamespaceIDModuleID(ctx, svc.store, p.NamespaceID, p.ModuleID); e != nil && e.ID != p.ID {
			return PageLayoutErrModuleNotFound()
		}
	}

	return nil
}

func (svc pageLayout) handleUpdate(ctx context.Context, upd *types.PageLayout) pageLayoutUpdateHandler {
	return func(ctx context.Context, ns *types.Namespace, res *types.PageLayout) (changes pageLayoutChanges, err error) {
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return pageLayoutUnchanged, PageLayoutErrStaleData()
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return pageLayoutUnchanged, PageLayoutErrInvalidHandle()
		}

		if err := svc.uniqueCheck(ctx, upd); err != nil {
			return pageLayoutUnchanged, err
		}

		// if !svc.ac.CanUpdatePageLayout(ctx, res) {
		// 	return pageLayoutUnchanged, PageLayoutErrNotAllowedToUpdate()
		// }

		// Get max blockID for later use
		blockID := uint64(0)
		for _, b := range res.Blocks {
			if b.BlockID > blockID {
				blockID = b.BlockID
			}
		}

		if res.ID != upd.ID {
			res.ID = upd.ID
			changes |= pageLayoutChanged
		}

		if res.PageID != upd.PageID {
			res.PageID = upd.PageID
			changes |= pageLayoutChanged
		}

		if res.ParentID != upd.ParentID {
			res.ParentID = upd.ParentID
			changes |= pageLayoutChanged
		}

		if res.NamespaceID != upd.NamespaceID {
			res.NamespaceID = upd.NamespaceID
			changes |= pageLayoutChanged
		}

		if res.ModuleID != upd.ModuleID {
			res.ModuleID = upd.ModuleID
			changes |= pageLayoutChanged
		}

		if !reflect.DeepEqual(res.Config, upd.Config) {
			res.Config = upd.Config
			changes |= pageLayoutChanged
		}

		if !reflect.DeepEqual(res.Blocks, upd.Blocks) {
			res.Blocks = upd.Blocks
			changes |= pageLayoutChanged
		}

		// Assure blockIDs
		for i, b := range res.Blocks {
			if b.BlockID == 0 {
				blockID++
				b.BlockID = blockID
				res.Blocks[i] = b

				changes |= pageLayoutChanged
			}
		}

		if res.Meta != nil && res.Meta.Name != upd.Meta.Name {
			res.Meta.Name = upd.Meta.Name
			changes |= pageLayoutChanged
		}

		if res.Handle != upd.Handle {
			res.Handle = upd.Handle
			changes |= pageLayoutChanged
		}

		if res.Meta != nil && res.Meta.Description != upd.Meta.Description {
			res.Meta.Description = upd.Meta.Description
			changes |= pageLayoutChanged
		}

		if res.Primary != upd.Primary {
			res.Primary = upd.Primary
			changes |= pageLayoutChanged
		}

		if upd.Labels != nil {
			if label.Changed(res.Labels, upd.Labels) {
				changes |= pageLayoutLabelsChanged
				res.Labels = upd.Labels
			}
		}

		if changes&pageLayoutChanged > 0 {
			res.UpdatedAt = now()
		}

		return
	}
}

func (svc pageLayout) handleDelete(ctx context.Context, ns *types.Namespace, m *types.PageLayout) (pageLayoutChanges, error) {
	// if !svc.ac.CanDeletePageLayout(ctx, m) {
	// 	return pageLayoutUnchanged, PageLayoutErrNotAllowedToDelete()
	// }

	if m.DeletedAt != nil {
		// pageLayout already deleted
		return pageLayoutUnchanged, nil
	}

	m.DeletedAt = now()
	return pageLayoutChanged, nil
}

func (svc pageLayout) handleUndelete(ctx context.Context, ns *types.Namespace, m *types.PageLayout) (pageLayoutChanges, error) {
	// if !svc.ac.CanDeletePageLayout(ctx, m) {
	// 	return pageLayoutUnchanged, PageLayoutErrNotAllowedToUndelete()
	// }

	if m.DeletedAt == nil {
		// pageLayout not deleted
		return pageLayoutUnchanged, nil
	}

	m.DeletedAt = nil
	return pageLayoutChanged, nil
}

func (svc *pageLayout) UpdateConfig(ss *systemTypes.AppSettings) {
	a := ss.Compose.UI.RecordToolbar

	svc.pageLayoutSettings = &pageLayoutSettings{
		hideNew:    a.HideNew,
		hideEdit:   a.HideEdit,
		hideSubmit: a.HideSubmit,
		hideDelete: a.HideDelete,
		hideClone:  a.HideClone,
		hideBack:   a.HideBack,
	}
}

func loadPageLayoutCombo(ctx context.Context, s interface {
	store.ComposePageLayouts
	store.ComposeNamespaces
}, namespaceID, pageID, pageLayoutID uint64) (ns *types.Namespace, c *types.PageLayout, err error) {
	ns, err = loadNamespace(ctx, s, namespaceID)
	if err != nil {
		return
	}

	c, err = loadPageLayout(ctx, s, namespaceID, pageID, pageLayoutID)
	return
}

func loadPageLayout(ctx context.Context, s store.ComposePageLayouts, namespaceID, pageID, pageLayoutID uint64) (res *types.PageLayout, err error) {
	if pageLayoutID == 0 || namespaceID == 0 {
		return nil, PageLayoutErrInvalidID()
	}

	if res, err = store.LookupComposePageLayoutByID(ctx, s, pageLayoutID); errors.IsNotFound(err) {
		err = PageLayoutErrNotFound()
	}

	if err == nil && namespaceID != res.NamespaceID {
		// Make sure chart belongs to the right namespace
		return nil, PageLayoutErrNotFound()
	}

	if err == nil && namespaceID != res.NamespaceID {
		// Make sure pageLayout belongs to the right namespace
		return nil, PageLayoutErrNotFound()
	}

	return
}

// toLabeledPageLayouts converts to []label.LabeledResource
//
// This function is auto-generated.
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
