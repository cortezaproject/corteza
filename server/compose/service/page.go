package service

import (
	"context"
	"reflect"

	"github.com/crusttech/human/server/compose/service/event"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/store"
	systemTypes "github.com/crusttech/human/server/system/types"
)

type (
	pageServices struct {
		eventbus     eventDispatcher
		locale       ResourceTranslationsManagerService
		pageSettings *pageSettings
	}

	pageSettings struct {
		hideNew    bool
		hideEdit   bool
		hideSubmit bool
		hideDelete bool
		hideClone  bool
		hideBack   bool
	}

	pageAccessController interface {
		CanManageResourceTranslations(ctx context.Context) bool
		CanSearchPagesOnNamespace(context.Context, *types.Namespace) bool
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanCreatePageOnNamespace(context.Context, *types.Namespace) bool
		CanReadPage(context.Context, *types.Page) bool
		CanUpdatePage(context.Context, *types.Page) bool
		CanDeletePage(context.Context, *types.Page) bool
	}

	pageUpdateHandler func(ctx context.Context, ns *types.Namespace, c *types.Page) (pageChanges, error)
	pageChanges       uint8
)

const (
	pageUnchanged     pageChanges = 0
	pageChanged       pageChanges = 1
	pageLabelsChanged pageChanges = 2
)

func Page() *page {
	return &page{
		actionlog: DefaultActionlog,
		ac:        DefaultAccessControl,
		store:     DefaultStore,
		services: &pageServices{
			eventbus:     eventbus.Service(),
			locale:       DefaultResourceTranslation,
			pageSettings: &pageSettings{},
		},
	}
}

func (svc page) FindByHandle(ctx context.Context, namespaceID uint64, h string) (c *types.Page, err error) {
	return svc.lookup(ctx, namespaceID, func(aProps *pageActionProps) (*types.Page, error) {
		if !handle.IsValid(h) {
			return nil, PageErrInvalidHandle()
		}

		aProps.page.Handle = h
		return store.LookupComposePageByNamespaceIDHandle(ctx, svc.store, namespaceID, h)
	})
}

func (svc page) FindByPageID(ctx context.Context, namespaceID, pageID uint64) (p *types.Page, err error) {
	return svc.lookup(ctx, namespaceID, func(aProps *pageActionProps) (*types.Page, error) {
		if pageID == 0 {
			return nil, PageErrInvalidID()
		}

		aProps.page.ID = pageID
		return store.LookupComposePageByID(ctx, svc.store, pageID)
	})
}

func checkPage(ctx context.Context, ac pageAccessController) func(res *types.Page) (bool, error) {
	return func(res *types.Page) (bool, error) {
		if !ac.CanReadPage(ctx, res) {
			return false, nil
		}

		return true, nil
	}
}

// search fn() orchestrates pages search, namespace preload and check
func (svc page) search(ctx context.Context, filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error) {
	var (
		aProps = &pageActionProps{filter: &filter}
		ns     *types.Namespace
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = checkPage(ctx, svc.ac)

	err = func() error {
		ns, err = loadNamespace(ctx, svc.store, filter.NamespaceID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)
		if !svc.ac.CanSearchPagesOnNamespace(ctx, ns) {
			return PageErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Page{}.LabelResourceKind(),
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

		if set, f, err = store.SearchComposePages(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledPages(set)...); err != nil {
			return err
		}

		// i18n
		tag := locale.GetAcceptLanguageFromContext(ctx)
		set.Walk(func(p *types.Page) error {
			p.DecodeTranslations(svc.services.locale.Locale().ResourceTranslations(tag, p.ResourceTranslation()))
			return nil
		})

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, PageActionSearch, err)
}

func (svc page) FindBySelfID(ctx context.Context, namespaceID, parentID uint64) (pp types.PageSet, f types.PageFilter, err error) {
	return svc.search(ctx, types.PageFilter{
		NamespaceID: namespaceID,
		ParentID:    parentID,

		// This will enable parentID=0 query
		Root: true,

		Check: checkPage(ctx, svc.ac),
	})
}

func (svc page) Find(ctx context.Context, filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error) {
	return svc.search(ctx, filter)
}

func (svc page) updater(ctx context.Context, s store.Storer, ns *types.Namespace, res *types.Page, action func(...*pageActionProps) *pageAction, fn pageUpdateHandler) (*types.Page, error) {
	var (
		changes pageChanges
		old     *types.Page
		aProps  = &pageActionProps{page: res}
		err     error
	)

	err = store.Tx(ctx, s, func(ctx context.Context, s store.Storer) (err error) {
		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		old = res.Clone()

		aProps.setNamespace(ns)
		aProps.setChanged(res)

		if res.DeletedAt == nil {
			err = svc.services.eventbus.WaitFor(ctx, event.PageBeforeUpdate(old, res, ns, nil))
		} else {
			err = svc.services.eventbus.WaitFor(ctx, event.PageBeforeDelete(old, res, ns, nil))
		}

		if err != nil {
			return
		}

		if changes, err = fn(ctx, ns, res); err != nil {
			return err
		}

		if changes&pageChanged > 0 {
			if err = store.UpdateComposePage(ctx, s, res); err != nil {
				return err
			}
		}

		if err = updateTranslations(ctx, svc.ac, svc.services.locale, res.EncodeTranslations()...); err != nil {
			return
		}

		if changes&pageLabelsChanged > 0 {
			if err = label.Update(ctx, s, res); err != nil {
				return
			}
		}

		if res.DeletedAt == nil {
			err = svc.services.eventbus.WaitFor(ctx, event.PageAfterUpdate(res, res, ns, nil))
		} else {
			err = svc.services.eventbus.WaitFor(ctx, event.PageAfterDelete(nil, res, ns, nil))
		}

		return err
	})

	return res, svc.recordAction(ctx, aProps, action, err)
}

// lookup fn() orchestrates page lookup, namespace preload and check
func (svc page) lookup(ctx context.Context, namespaceID uint64, lookup func(*pageActionProps) (*types.Page, error)) (p *types.Page, err error) {
	var aProps = &pageActionProps{page: &types.Page{NamespaceID: namespaceID}}

	err = func() error {
		if ns, err := loadNamespace(ctx, svc.store, namespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if p, err = lookup(aProps); errors.IsNotFound(err) {
			return PageErrNotFound()
		} else if err != nil {
			return err
		}

		p.DecodeTranslations(svc.services.locale.Locale().ResourceTranslations(locale.GetAcceptLanguageFromContext(ctx), p.ResourceTranslation()))

		aProps.setPage(p)

		if !svc.ac.CanReadPage(ctx, p) {
			return PageErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, p); err != nil {
			return err
		}

		return nil
	}()

	return p, svc.recordAction(ctx, aProps, PageActionLookup, err)
}

func (svc page) uniqueCheck(ctx context.Context, p *types.Page) (err error) {
	if p.Handle != "" {
		if e, _ := store.LookupComposePageByNamespaceIDHandle(ctx, svc.store, p.NamespaceID, p.Handle); e != nil && e.ID != p.ID {
			return PageErrHandleNotUnique()
		}
	}

	if p.ModuleID > 0 {
		if e, _ := store.LookupComposePageByNamespaceIDModuleID(ctx, svc.store, p.NamespaceID, p.ModuleID); e != nil && e.ID != p.ID {
			return PageErrModuleNotFound()
		}
	}

	return nil
}

func (svc page) handleUpdate(ctx context.Context, upd *types.Page) pageUpdateHandler {
	return func(ctx context.Context, ns *types.Namespace, res *types.Page) (changes pageChanges, err error) {
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return pageUnchanged, PageErrStaleData()
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return pageUnchanged, PageErrInvalidHandle()
		}

		if err := svc.uniqueCheck(ctx, upd); err != nil {
			return pageUnchanged, err
		}

		if !svc.ac.CanUpdatePage(ctx, res) {
			return pageUnchanged, PageErrNotAllowedToUpdate()
		}

		// Get max blockID for later use
		blockID := uint64(0)
		for _, b := range res.Blocks {
			if b.BlockID > blockID {
				blockID = b.BlockID
			}
		}

		if res.ID != upd.ID {
			res.ID = upd.ID
			changes |= pageChanged
		}

		if res.SelfID != upd.SelfID {
			res.SelfID = upd.SelfID
			changes |= pageChanged
		}

		if !reflect.DeepEqual(res.Config, upd.Config) {
			res.Config = upd.Config
			changes |= pageChanged
		}

		if !reflect.DeepEqual(res.Blocks, upd.Blocks) {
			res.Blocks = upd.Blocks
			changes |= pageChanged
		}

		// Assure blockIDs
		for i, b := range res.Blocks {
			if b.BlockID == 0 {
				blockID++
				b.BlockID = blockID
				res.Blocks[i] = b

				changes |= pageChanged
			}
		}

		if res.Meta.AllowPersonalLayouts != upd.Meta.AllowPersonalLayouts {
			res.Meta.AllowPersonalLayouts = upd.Meta.AllowPersonalLayouts
			changes |= pageChanged
		}

		if res.Title != upd.Title {
			res.Title = upd.Title
			changes |= pageChanged
		}

		if res.Handle != upd.Handle {
			res.Handle = upd.Handle
			changes |= pageChanged
		}

		if res.Description != upd.Description {
			res.Description = upd.Description
			changes |= pageChanged
		}

		if res.Visible != upd.Visible {
			res.Visible = upd.Visible
			changes |= pageChanged
		}

		if res.Weight != upd.Weight {
			res.Weight = upd.Weight
			changes |= pageChanged
		}

		if upd.Labels != nil {
			if label.Changed(res.Labels, upd.Labels) {
				changes |= pageLabelsChanged
				res.Labels = upd.Labels
			}
		}

		if changes&pageChanged > 0 {
			res.UpdatedAt = now()
		}

		return
	}
}

func (svc page) handleDelete(ctx context.Context, ns *types.Namespace, m *types.Page) (pageChanges, error) {
	if !svc.ac.CanDeletePage(ctx, m) {
		return pageUnchanged, PageErrNotAllowedToDelete()
	}

	if m.DeletedAt != nil {
		// page already deleted
		return pageUnchanged, nil
	}

	m.DeletedAt = now()
	return pageChanged, nil
}

func (svc page) handleUndelete(ctx context.Context, ns *types.Namespace, m *types.Page) (pageChanges, error) {
	if !svc.ac.CanDeletePage(ctx, m) {
		return pageUnchanged, PageErrNotAllowedToUndelete()
	}

	if m.DeletedAt == nil {
		// page not deleted
		return pageUnchanged, nil
	}

	m.DeletedAt = nil
	return pageChanged, nil
}

func (svc *page) UpdateConfig(ss *systemTypes.AppSettings) {
	a := ss.Compose.UI.RecordToolbar

	svc.services.pageSettings = &pageSettings{
		hideNew:    a.HideNew,
		hideEdit:   a.HideEdit,
		hideSubmit: a.HideSubmit,
		hideDelete: a.HideDelete,
		hideClone:  a.HideClone,
		hideBack:   a.HideBack,
	}
}

func (svc *page) onLookup(ctx context.Context, namespaceID uint64, ID uint64, aProps *pageActionProps) (*types.Page, error) {
	return svc.lookup(ctx, namespaceID, func(props *pageActionProps) (*types.Page, error) {
		if ID == 0 {
			return nil, PageErrInvalidID()
		}
		props.page.ID = ID
		return store.LookupComposePageByID(ctx, svc.store, ID)
	})
}

func (svc *page) onSearch(ctx context.Context, filter types.PageFilter, _ *pageActionProps) (types.PageSet, types.PageFilter, error) {
	return svc.search(ctx, filter)
}

func (svc *page) onCreate(ctx context.Context, new *types.Page) error {
	var (
		ns *types.Namespace
	)

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !handle.IsValid(new.Handle) {
			return PageErrInvalidID()
		}

		if ns, err = loadNamespace(ctx, s, new.NamespaceID); err != nil {
			return err
		}

		if !svc.ac.CanCreatePageOnNamespace(ctx, ns) {
			return PageErrNotAllowedToCreate()
		}

		if err = svc.services.eventbus.WaitFor(ctx, event.PageBeforeCreate(new, nil, ns, nil)); err != nil {
			return err
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.UpdatedAt = nil
		new.DeletedAt = nil

		for i := range new.Blocks {
			new.Blocks[i].BlockID = uint64(i) + 1
		}

		if err = store.CreateComposePage(ctx, s, new); err != nil {
			return err
		}

		if err = updateTranslations(ctx, svc.ac, svc.services.locale, new.EncodeTranslations()...); err != nil {
			return
		}

		if err = label.Create(ctx, s, new); err != nil {
			return
		}

		_ = svc.services.eventbus.WaitFor(ctx, event.PageAfterCreate(new, nil, ns, nil))
		return err
	})
}

func (svc *page) onUpdate(ctx context.Context, s store.Storer, upd *types.Page, res *types.Page, _ *pageActionProps, before func() error, after func() error) error {
	if err := before(); err != nil {
		return err
	}
	ns, err := loadNamespace(ctx, svc.store, res.NamespaceID)
	if err != nil {
		return err
	}
	_, err = svc.updater(ctx, s, ns, res, PageActionUpdate, svc.handleUpdate(ctx, upd))
	if err != nil {
		return err
	}
	return after()
}

func (svc *page) onDelete(ctx context.Context, s store.Storer, namespaceID uint64, res *types.Page, strategy types.PageChildrenDeleteStrategy, _ *pageActionProps) error {
	var (
		validChildren, pp types.PageSet

		skipUndeleted = func(p *types.Page) (bool, error) {
			return p.DeletedAt == nil, nil
		}
	)

	deleteOne := func(ns *types.Namespace, p *types.Page) error {
		_, err := svc.updater(ctx, s, ns, p, PageActionDelete, svc.handleDelete)
		return err
	}

	// Force skips the family-tree resolution and deletes only this page.
	if strategy == types.PageChildrenOnDeleteForce {
		ns, err := loadNamespace(ctx, svc.store, namespaceID)
		if err != nil {
			return err
		}
		return deleteOne(ns, res)
	}

	// Load all pages in the namespace to resolve the family tree.
	pp, _, err := store.SearchComposePages(ctx, s, types.PageFilter{NamespaceID: namespaceID})
	if err != nil {
		return err
	}

	// res was loaded by the generated wrapper; use the copy from the full set so
	// FindByParent/RecursiveWalk operate on the same instances.
	if ppRes := pp.FindByID(res.ID); ppRes != nil {
		res = ppRes
	}

	validChildren, _ = pp.FindByParent(res.ID).Filter(skipUndeleted)

	ns, err := loadNamespace(ctx, svc.store, namespaceID)
	if err != nil {
		return err
	}

	switch strategy {
	case types.PageChildrenOnDeleteAbort:
		if len(validChildren) > 0 {
			return PageErrDeleteAbortedForPageWithSubpages()
		}

	case types.PageChildrenOnDeleteRebase:
		err = validChildren.Walk(func(child *types.Page) error {
			updChild := child.Clone()
			updChild.SelfID = res.SelfID
			_, err := svc.updater(ctx, s, ns, child, PageActionUpdate, svc.handleUpdate(ctx, updChild))
			return err
		})
		if err != nil {
			return err
		}

	case types.PageChildrenOnDeleteCascade:
		err = pp.RecursiveWalk(res, func(child *types.Page, _ *types.Page) error {
			if child.DeletedAt != nil {
				return nil
			}
			return deleteOne(ns, child)
		})
		if err != nil {
			return err
		}

	default:
		return PageErrUnknownDeleteStrategy()
	}

	return deleteOne(ns, res)
}

func (svc *page) onUndelete(ctx context.Context, s store.Storer, namespaceID uint64, res *types.Page, _ *pageActionProps) error {
	ns, err := loadNamespace(ctx, svc.store, namespaceID)
	if err != nil {
		return err
	}
	_, err = svc.updater(ctx, s, ns, res, PageActionUpdate, svc.handleUndelete)
	return err
}

func (svc *page) onTree(ctx context.Context, _ *pageActionProps, namespaceID uint64) (types.PageSet, error) {
	var (
		pages  types.PageSet
		tree   types.PageSet
		filter = types.PageFilter{
			NamespaceID: namespaceID,
			Check:       checkPage(ctx, svc.ac),
		}
		err error
	)

	if err = filter.Sort.Set("weight ASC"); err != nil {
		return nil, err
	}

	if pages, _, err = svc.search(ctx, filter); err != nil {
		return nil, err
	}

	_ = pages.Walk(func(p *types.Page) error {
		if p.SelfID == 0 {
			tree = append(tree, p)
		} else if c := pages.FindByID(p.SelfID); c != nil {
			if c.Children == nil {
				c.Children = types.PageSet{}
			}
			c.Children = append(c.Children, p)
		} else {
			p.SelfID = 0
			tree = append(tree, p)
		}
		return nil
	})

	return tree, nil
}

func (svc *page) onReorder(ctx context.Context, _ *pageActionProps, namespaceID uint64, parentID uint64, pageIDs []uint64) error {
	var (
		ns  *types.Namespace
		p   *types.Page
		err error
	)

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if ns, err = loadNamespace(ctx, s, namespaceID); err != nil {
			return err
		}
		_ = ns

		if parentID == 0 {
			if !svc.ac.CanCreatePageOnNamespace(ctx, ns) {
				return PageErrNotAllowedToUpdate()
			}
		} else {
			if p, err = store.LookupComposePageByID(ctx, s, parentID); errors.IsNotFound(err) {
				return PageErrNotFound()
			} else if err != nil {
				return err
			}
			if !svc.ac.CanUpdatePage(ctx, p) {
				return PageErrNotAllowedToUpdate()
			}
		}

		return store.ReorderComposePages(ctx, s, namespaceID, parentID, pageIDs)
	})
}

func (svc *page) onUpdateIcon(ctx context.Context, _ *pageActionProps, namespaceID uint64, pageID uint64, icon *types.PageConfigIcon) (out *types.PageConfigIcon, err error) {
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, p, err := loadPageCombo(ctx, s, namespaceID, pageID)
		if err != nil {
			return err
		}

		p.Config.NavItem.Icon = icon
		p, err = svc.updater(ctx, svc.store, ns, p, PageActionUpdate, svc.handleUpdate(ctx, p))
		if err != nil {
			return err
		}
		out = p.Config.NavItem.Icon
		return nil
	})
	return
}

func loadPageCombo(ctx context.Context, s interface {
	store.ComposePages
	store.ComposeNamespaces
}, namespaceID, pageID uint64) (ns *types.Namespace, c *types.Page, err error) {
	ns, err = loadNamespace(ctx, s, namespaceID)
	if err != nil {
		return
	}

	if c, err = loadPage(ctx, s, pageID); err != nil {
		return
	}

	if namespaceID != c.NamespaceID {
		return nil, nil, PageErrNotFound()
	}

	return
}
