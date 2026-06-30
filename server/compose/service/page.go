package service

import (
	"context"
	"reflect"

	"github.com/crusttech/human/server/compose/service/event"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/store"
	systemTypes "github.com/crusttech/human/server/system/types"
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

// onLookup is the generated FindByID body handler (namespace-scoped compound id).
func (svc *page) onLookup(ctx context.Context, namespaceID, pageID uint64, aProps *pageActionProps) (p *types.Page, err error) {
	return svc.lookup(ctx, namespaceID, aProps, func(aProps *pageActionProps) (*types.Page, error) {
		if pageID == 0 {
			return nil, PageErrInvalidID()
		}

		aProps.page.ID = pageID
		return store.LookupComposePageByID(ctx, svc.store, pageID)
	})
}

func (svc *page) FindByHandle(ctx context.Context, namespaceID uint64, h string) (c *types.Page, err error) {
	var aProps = &pageActionProps{page: &types.Page{NamespaceID: namespaceID}}

	c, err = svc.lookup(ctx, namespaceID, aProps, func(aProps *pageActionProps) (*types.Page, error) {
		if !handle.IsValid(h) {
			return nil, PageErrInvalidHandle()
		}

		aProps.page.Handle = h
		return store.LookupComposePageByNamespaceIDHandle(ctx, svc.store, namespaceID, h)
	})

	return c, svc.recordAction(ctx, aProps, PageActionLookup, err)
}

func (svc *page) FindByPageID(ctx context.Context, namespaceID, pageID uint64) (p *types.Page, err error) {
	var aProps = &pageActionProps{page: &types.Page{NamespaceID: namespaceID}}

	p, err = svc.lookup(ctx, namespaceID, aProps, func(aProps *pageActionProps) (*types.Page, error) {
		if pageID == 0 {
			return nil, PageErrInvalidID()
		}

		aProps.page.ID = pageID
		return store.LookupComposePageByID(ctx, svc.store, pageID)
	})

	return p, svc.recordAction(ctx, aProps, PageActionLookup, err)
}

func checkPage(ctx context.Context, ac pageAccessController) func(res *types.Page) (bool, error) {
	return func(res *types.Page) (bool, error) {
		if !ac.CanReadPage(ctx, res) {
			return false, nil
		}

		return true, nil
	}
}

// onSearch is the generated Search body handler.
//
// The recordAction wrapper and aProps (filter) are owned by the generated
// page.gen.go; this handler runs the namespace preload, access check and store
// search.
func (svc *page) onSearch(ctx context.Context, filter types.PageFilter, aProps *pageActionProps) (set types.PageSet, f types.PageFilter, err error) {
	var ns *types.Namespace

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = checkPage(ctx, svc.ac)

	ns, err = loadNamespace(ctx, svc.store, filter.NamespaceID)
	if err != nil {
		return
	}

	aProps.setNamespace(ns)
	if !svc.ac.CanSearchPagesOnNamespace(ctx, ns) {
		return nil, f, PageErrNotAllowedToSearch()
	}

	if len(filter.Labels) > 0 {
		filter.LabeledIDs, err = label.Search(
			ctx,
			svc.store,
			types.Page{}.LabelResourceKind(),
			filter.Labels,
		)

		if err != nil {
			return
		}

		// labels specified but no labeled resources found
		if len(filter.LabeledIDs) == 0 {
			return
		}
	}

	if set, f, err = store.SearchComposePages(ctx, svc.store, filter); err != nil {
		return
	}

	if err = label.Load(ctx, svc.store, toLabeledPages(set)...); err != nil {
		return
	}

	// i18n
	tag := locale.GetAcceptLanguageFromContext(ctx)
	set.Walk(func(p *types.Page) error {
		p.DecodeTranslations(svc.locale.Locale().ResourceTranslations(tag, p.ResourceTranslation()))
		return nil
	})

	return
}

// search fn() orchestrates pages search; it is retained for the custom methods
// (FindBySelfID/Tree) that need the recorded-action wrapper.
func (svc *page) search(ctx context.Context, filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error) {
	var aProps = &pageActionProps{filter: &filter}

	set, f, err = svc.onSearch(ctx, filter, aProps)

	return set, f, svc.recordAction(ctx, aProps, PageActionSearch, err)
}

func (svc *page) FindBySelfID(ctx context.Context, namespaceID, parentID uint64) (pp types.PageSet, f types.PageFilter, err error) {
	return svc.search(ctx, types.PageFilter{
		NamespaceID: namespaceID,
		ParentID:    parentID,

		// This will enable parentID=0 query
		Root: true,

		Check: checkPage(ctx, svc.ac),
	})
}

func (svc *page) Tree(ctx context.Context, namespaceID uint64) (tree types.PageSet, err error) {
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
func (svc *page) Reorder(ctx context.Context, namespaceID, parentID uint64, pageIDs []uint64) (err error) {
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

// onCreate is the generated Create body handler.
//
// The recordAction wrapper, aProps and res=new assignment are owned by the
// generated page.gen.go.
func (svc *page) onCreate(ctx context.Context, new *types.Page) error {
	var (
		ns     *types.Namespace
		aProps = &pageActionProps{page: new}
	)

	new.ID = 0

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

		// ProjectID comes from the caller (the REST create param), not derived
		// from the namespace — the project flow (record pages, standalone pages)
		// stamps it explicitly. Callers that don't send it create tenant-level
		// pages (project_id 0), the same as modules.

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
}

// onUpdate is the generated Update body handler.
func (svc *page) onUpdate(ctx context.Context, s store.Storer, upd, res *types.Page, aProps *pageActionProps, _ func() error, _ func() error) error {
	// ModuleID is the page's record binding, set at create and immutable on
	// update; preserve the existing value so a caller that doesn't resend it
	// (e.g. the page builder save) can't wipe the binding. res is the loaded
	// record; the generated Update copies upd onto it after this hook.
	upd.ModuleID = res.ModuleID

	ns, err := loadNamespace(ctx, s, res.NamespaceID)
	if err != nil {
		return err
	}

	aProps.setNamespace(ns)

	old := res.Clone()
	if err = svc.eventbus.WaitFor(ctx, event.PageBeforeUpdate(old, res, ns, nil)); err != nil {
		return err
	}

	changes, err := svc.handleUpdate(ctx, upd)(ctx, ns, res)
	if err != nil {
		return err
	}

	if changes&pageChanged > 0 {
		if err = store.UpdateComposePage(ctx, s, res); err != nil {
			return err
		}
	}

	if err = updateTranslations(ctx, svc.ac, svc.locale, res.EncodeTranslations()...); err != nil {
		return err
	}

	if changes&pageLabelsChanged > 0 {
		if err = label.Update(ctx, s, res); err != nil {
			return err
		}
	}

	return svc.eventbus.WaitFor(ctx, event.PageAfterUpdate(res, res, ns, nil))
}

// onDelete is the generated DeleteByID body handler (namespace-scoped compound
// id + child-delete strategy). The recordAction wrapper and aProps are owned by
// the generated page.gen.go.
func (svc *page) onDelete(ctx context.Context, s store.Storer, namespaceID uint64, res *types.Page, strategy types.PageChildrenDeleteStrategy, aProps *pageActionProps) error {
	var (
		validChildren, pp types.PageSet
		ns                *types.Namespace

		skipUndeleted = func(p *types.Page) (bool, error) {
			return p.DeletedAt == nil, nil
		}
	)

	if strategy == types.PageChildrenOnDeleteForce {
		ns, err := loadNamespace(ctx, s, namespaceID)
		if err != nil {
			return err
		}
		return svc.deleteOne(ctx, s, ns, res)
	}

	// Load all pages in the namespace to figure out the family tree
	pp, _, err := store.SearchComposePages(ctx, s, types.PageFilter{
		NamespaceID: namespaceID,
	})
	if err != nil {
		return err
	}

	// res was already loaded by gen.go; find it in the full set so we can use FindByParent
	if ppRes := pp.FindByID(res.ID); ppRes != nil {
		res = ppRes
	}

	validChildren, _ = pp.FindByParent(res.ID).Filter(skipUndeleted)

	switch strategy {
	case types.PageChildrenOnDeleteAbort:
		if len(validChildren) > 0 {
			return PageErrDeleteAbortedForPageWithSubpages()
		}

	case types.PageChildrenOnDeleteRebase:
		if ns, err = loadNamespace(ctx, s, namespaceID); err != nil {
			return err
		}
		err = validChildren.Walk(func(child *types.Page) (err error) {
			updChild := child.Clone()
			updChild.SelfID = res.SelfID
			_, err = svc.updater(ctx, s, ns, child, PageActionUpdate, svc.handleUpdate(ctx, updChild))
			return err
		})
		if err != nil {
			return err
		}

	case types.PageChildrenOnDeleteCascade:
		if ns, err = loadNamespace(ctx, s, namespaceID); err != nil {
			return err
		}
		err = pp.RecursiveWalk(res, func(child *types.Page, _ *types.Page) (err error) {
			if child.DeletedAt != nil {
				return nil
			}
			_, err = svc.updater(ctx, s, ns, child, PageActionDelete, svc.handleDelete)
			return err
		})
		if err != nil {
			return err
		}

	default:
		return PageErrUnknownDeleteStrategy()
	}

	if ns == nil {
		if ns, err = loadNamespace(ctx, s, namespaceID); err != nil {
			return err
		}
	}

	return svc.deleteOne(ctx, s, ns, res)
}

// onUndelete is the generated UndeleteByID body handler (namespace-scoped
// compound id).
func (svc *page) onUndelete(ctx context.Context, s store.Storer, namespaceID uint64, res *types.Page, aProps *pageActionProps) error {
	ns, err := loadNamespace(ctx, s, res.NamespaceID)
	if err != nil {
		return err
	}

	if !svc.ac.CanDeletePage(ctx, res) {
		return PageErrNotAllowedToUndelete()
	}
	if res.DeletedAt == nil {
		return nil
	}

	old := res.Clone()
	// res.DeletedAt != nil before undelete → fires delete-before event
	if err = svc.eventbus.WaitFor(ctx, event.PageBeforeDelete(old, res, ns, nil)); err != nil {
		return err
	}

	res.DeletedAt = nil
	if err = store.UpdateComposePage(ctx, s, res); err != nil {
		return err
	}

	// res.DeletedAt == nil after undelete → fires update-after event
	return svc.eventbus.WaitFor(ctx, event.PageAfterUpdate(res, res, ns, nil))
}

// deleteOne applies handleDelete logic for a single page without calling
// recordAction (the generated DeleteByID scaffold owns the action log entry).
func (svc *page) deleteOne(ctx context.Context, s store.Storer, ns *types.Namespace, res *types.Page) error {
	if !svc.ac.CanDeletePage(ctx, res) {
		return PageErrNotAllowedToDelete()
	}
	if res.DeletedAt != nil {
		return nil
	}

	old := res.Clone()
	// res.DeletedAt == nil before delete → fires update-before event (matches updater behaviour)
	if err := svc.eventbus.WaitFor(ctx, event.PageBeforeUpdate(old, res, ns, nil)); err != nil {
		return err
	}

	res.DeletedAt = now()
	if err := store.UpdateComposePage(ctx, s, res); err != nil {
		return err
	}

	// res.DeletedAt != nil after delete → fires delete-after event
	return svc.eventbus.WaitFor(ctx, event.PageAfterDelete(nil, res, ns, nil))
}

func (svc *page) UpdateIcon(ctx context.Context, namespaceID, pageID uint64, icon *types.PageConfigIcon) (out *types.PageConfigIcon, err error) {
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

func (svc *page) updater(ctx context.Context, s store.Storer, ns *types.Namespace, res *types.Page, action func(...*pageActionProps) *pageAction, fn pageUpdateHandler) (*types.Page, error) {
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

	return res, svc.recordAction(ctx, aProps, action, err, old, res)
}

// lookup fn() orchestrates page lookup, namespace preload and check.
//
// The recordAction wrapper is owned by the caller (the generated onLookup path
// via page.gen.go, or the custom FindBy* methods).
func (svc *page) lookup(ctx context.Context, namespaceID uint64, aProps *pageActionProps, lookup func(*pageActionProps) (*types.Page, error)) (p *types.Page, err error) {
	if aProps.page == nil {
		aProps.page = &types.Page{NamespaceID: namespaceID}
	}

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

	return p, err
}

func (svc *page) uniqueCheck(ctx context.Context, p *types.Page) (err error) {
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

func (svc *page) handleUpdate(ctx context.Context, upd *types.Page) pageUpdateHandler {
	return func(ctx context.Context, ns *types.Namespace, res *types.Page) (changes pageChanges, err error) {
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

func (svc *page) handleDelete(ctx context.Context, ns *types.Namespace, m *types.Page) (pageChanges, error) {
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

func (svc *page) handleUndelete(ctx context.Context, ns *types.Namespace, m *types.Page) (pageChanges, error) {
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

	c, err = loadPageScoped(ctx, s, namespaceID, pageID)
	return
}

// loadPageScoped loads a page by ID and checks it belongs to the given namespace.
func loadPageScoped(ctx context.Context, s store.ComposePages, namespaceID, pageID uint64) (res *types.Page, err error) {
	if pageID == 0 || namespaceID == 0 {
		return nil, PageErrInvalidID()
	}

	if res, err = store.LookupComposePageByID(ctx, s, pageID); errors.IsNotFound(err) {
		err = PageErrNotFound()
	}

	if err == nil && namespaceID != res.NamespaceID {
		return nil, PageErrNotFound()
	}

	return
}
