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
	page struct {
		actionlog actionlog.Recorder
		ac        pageAccessController
		eventbus  eventDispatcher
		store     store.Storer
		locale    ResourceTranslationsManagerService

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
		actionlog:    DefaultActionlog,
		ac:           DefaultAccessControl,
		eventbus:     eventbus.Service(),
		store:        DefaultStore,
		locale:       DefaultResourceTranslation,
		pageSettings: &pageSettings{},
	}
}

func (svc page) FindByID(ctx context.Context, namespaceID, pageID uint64) (p *types.Page, err error) {
	return svc.lookup(ctx, namespaceID, func(aProps *pageActionProps) (*types.Page, error) {
		if pageID == 0 {
			return nil, PageErrInvalidID()
		}

		aProps.page.ID = pageID
		return store.LookupComposePageByID(ctx, svc.store, pageID)
	})
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
			p.DecodeTranslations(svc.locale.Locale().ResourceTranslations(tag, p.ResourceTranslation()))
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

func (svc page) Tree(ctx context.Context, namespaceID uint64) (tree types.PageSet, err error) {
	var (
		pages  types.PageSet
		filter = types.PageFilter{
			NamespaceID: namespaceID,
			Check:       checkPage(ctx, svc.ac),
		}
	)

	if err = filter.Sort.Set("weight ASC"); err != nil {
		return
	}

	if pages, _, err = svc.search(ctx, filter); err != nil {
		return
	}

	// safe to ignore errors
	_ = pages.Walk(func(p *types.Page) error {
		if p.SelfID == 0 {
			tree = append(tree, p)
		} else if c := pages.FindByID(p.SelfID); c != nil {
			if c.Children == nil {
				c.Children = types.PageSet{}
			}

			c.Children = append(c.Children, p)
		} else {
			// Move orphans to root
			p.SelfID = 0
			tree = append(tree, p)
		}

		return nil
	})

	return tree, nil
}

// Reorder pages
func (svc page) Reorder(ctx context.Context, namespaceID, parentID uint64, pageIDs []uint64) (err error) {
	var (
		aProps = &pageActionProps{page: &types.Page{ID: parentID}}
		ns     *types.Namespace
		p      *types.Page
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if ns, err = loadNamespace(ctx, s, namespaceID); err != nil {
			return err
		}

		if parentID == 0 {
			// Reordering on root mode -- check if user can create pages.
			if !svc.ac.CanCreatePageOnNamespace(ctx, ns) {
				return PageErrNotAllowedToUpdate()
			}
		} else {
			// Validate permissions on parent page
			if p, err = store.LookupComposePageByID(ctx, s, parentID); errors.IsNotFound(err) {
				return PageErrNotFound()
			} else if err != nil {
				return err
			}

			aProps.setPage(p)

			if !svc.ac.CanUpdatePage(ctx, p) {
				return PageErrNotAllowedToUpdate()
			}
		}

		return store.ReorderComposePages(ctx, s, namespaceID, parentID, pageIDs)
	})

	return svc.recordAction(ctx, aProps, PageActionReorder, err)

}

func (svc page) Create(ctx context.Context, new *types.Page) (*types.Page, error) {
	var (
		ns     *types.Namespace
		aProps = &pageActionProps{page: new}
	)

	new.ID = 0

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !handle.IsValid(new.Handle) {
			return PageErrInvalidID()
		}

		if ns, err = loadNamespace(ctx, s, new.NamespaceID); err != nil {
			return err
		}

		if !svc.ac.CanCreatePageOnNamespace(ctx, ns) {
			return PageErrNotAllowedToCreate()
		}

		aProps.setNamespace(ns)

		if err = svc.eventbus.WaitFor(ctx, event.PageBeforeCreate(new, nil, ns, nil)); err != nil {
			return err
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.UpdatedAt = nil
		new.DeletedAt = nil

		// Ensure page-block IDs
		for i := range new.Blocks {
			new.Blocks[i].BlockID = uint64(i) + 1
		}

		aProps.setChanged(new)

		if err = store.CreateComposePage(ctx, s, new); err != nil {
			return err
		}

		if err = updateTranslations(ctx, svc.ac, svc.locale, new.EncodeTranslations()...); err != nil {
			return
		}

		if err = label.Create(ctx, s, new); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.PageAfterCreate(new, nil, ns, nil))
		return err
	})

	return new, svc.recordAction(ctx, aProps, PageActionCreate, err)
}

func (svc page) Update(ctx context.Context, upd *types.Page) (c *types.Page, err error) {
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, res, err := loadPageCombo(ctx, s, upd.NamespaceID, upd.ID)
		if err != nil {
			return
		}

		c, err = svc.updater(ctx, svc.store, ns, res, PageActionUpdate, svc.handleUpdate(ctx, upd))
		return
	})

	return
}

func (svc page) DeleteByID(ctx context.Context, namespaceID, pageID uint64, strategy types.PageChildrenDeleteStrategy) error {
	var (
		validChildren, pp types.PageSet

		ns  *types.Namespace
		res *types.Page

		skipUndeleted = func(p *types.Page) (bool, error) {
			return p.DeletedAt == nil, nil
		}
	)

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if strategy == types.PageChildrenOnDeleteForce {
			// simply delete the page and ignore the subpages
			ns, res, err = loadPageCombo(ctx, s, namespaceID, pageID)
			if err != nil {
				return
			}
		} else {
			// Load all pages in the namespace and
			// try to figure out the family tree
			pp, _, err = store.SearchComposePages(ctx, s, types.PageFilter{
				NamespaceID: namespaceID,
			})

			if res = pp.FindByID(pageID); res == nil {
				return PageErrNotFound()
			}

			validChildren, _ = pp.FindByParent(res.ID).Filter(skipUndeleted)

			switch strategy {
			case types.PageChildrenOnDeleteAbort:
				// Abort if there are any valid (undeleted) children
				if len(validChildren) > 0 {
					return PageErrDeleteAbortedForPageWithSubpages()
				}

			case types.PageChildrenOnDeleteRebase:
				// update all our children to point to our parent
				err = validChildren.Walk(func(child *types.Page) (err error) {
					updChild := child.Clone()
					updChild.SelfID = res.SelfID
					_, err = svc.updater(ctx, s, ns, child, PageActionUpdate, svc.handleUpdate(ctx, updChild))
					return err
				})

				if err != nil {
					return
				}

			case types.PageChildrenOnDeleteCascade:
				// update all our children to point to our parent
				err = pp.RecursiveWalk(res, func(child *types.Page, _ *types.Page) (err error) {
					if child.DeletedAt != nil {
						// skip the ones that are already deleted
						return nil
					}

					_, err = svc.updater(ctx, s, ns, child, PageActionDelete, svc.handleDelete)
					return err
				})

				if err != nil {
					return
				}
			default:
				return PageErrUnknownDeleteStrategy()
			}

			if ns, err = loadNamespace(ctx, s, namespaceID); err != nil {
				return
			}
		}

		_, err = svc.updater(ctx, svc.store, ns, res, PageActionDelete, svc.handleDelete)
		return
	})
}

func (svc page) UndeleteByID(ctx context.Context, namespaceID, pageID uint64) error {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, res, err := loadPageCombo(ctx, s, namespaceID, pageID)
		if err != nil {
			return
		}

		_, err = svc.updater(ctx, svc.store, ns, res, PageActionUpdate, svc.handleUndelete)
		return
	})
}

func (svc page) UpdateIcon(ctx context.Context, namespaceID, pageID uint64, icon *types.PageConfigIcon) (out *types.PageConfigIcon, err error) {
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, p, err := loadPageCombo(ctx, s, namespaceID, pageID)
		if err != nil {
			return
		}

		p.Config.NavItem.Icon = icon
		p, err = svc.updater(ctx, svc.store, ns, p, PageActionUpdate, svc.handleUpdate(ctx, p))
		out = p.Config.NavItem.Icon

		return
	})

	return
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
			err = svc.eventbus.WaitFor(ctx, event.PageBeforeUpdate(old, res, ns, nil))
		} else {
			err = svc.eventbus.WaitFor(ctx, event.PageBeforeDelete(old, res, ns, nil))
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

		if err = updateTranslations(ctx, svc.ac, svc.locale, res.EncodeTranslations()...); err != nil {
			return
		}

		if changes&pageLabelsChanged > 0 {
			if err = label.Update(ctx, s, res); err != nil {
				return
			}
		}

		if res.DeletedAt == nil {
			err = svc.eventbus.WaitFor(ctx, event.PageAfterUpdate(res, res, ns, nil))
		} else {
			err = svc.eventbus.WaitFor(ctx, event.PageAfterDelete(nil, res, ns, nil))
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

		p.DecodeTranslations(svc.locale.Locale().ResourceTranslations(locale.GetAcceptLanguageFromContext(ctx), p.ResourceTranslation()))

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

		if !reflect.DeepEqual(res.Meta.Notifications, upd.Meta.Notifications) {
			res.Meta.Notifications = upd.Meta.Notifications
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

	svc.pageSettings = &pageSettings{
		hideNew:    a.HideNew,
		hideEdit:   a.HideEdit,
		hideSubmit: a.HideSubmit,
		hideDelete: a.HideDelete,
		hideClone:  a.HideClone,
		hideBack:   a.HideBack,
	}
}

func loadPageCombo(ctx context.Context, s interface {
	store.ComposePages
	store.ComposeNamespaces
}, namespaceID, pageID uint64) (ns *types.Namespace, c *types.Page, err error) {
	ns, err = loadNamespace(ctx, s, namespaceID)
	if err != nil {
		return
	}

	c, err = loadPage(ctx, s, namespaceID, pageID)
	return
}

func loadPage(ctx context.Context, s store.ComposePages, namespaceID, pageID uint64) (res *types.Page, err error) {
	if pageID == 0 || namespaceID == 0 {
		return nil, PageErrInvalidID()
	}

	if res, err = store.LookupComposePageByID(ctx, s, pageID); errors.IsNotFound(err) {
		err = PageErrNotFound()
	}

	if err == nil && namespaceID != res.NamespaceID {
		// Make sure chart belongs to the right namespace
		return nil, PageErrNotFound()
	}

	if err == nil && namespaceID != res.NamespaceID {
		// Make sure page belongs to the right namespace
		return nil, PageErrNotFound()
	}

	return
}

// toLabeledPages converts to []label.LabeledResource
//
// This function is auto-generated.
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
