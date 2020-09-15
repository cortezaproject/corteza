package service

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	page struct {
		ctx       context.Context
		actionlog actionlog.Recorder
		ac        pageAccessController
		eventbus  eventDispatcher
		store     store.Storer
	}

	pageAccessController interface {
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanCreatePage(context.Context, *types.Namespace) bool
		CanReadPage(context.Context, *types.Page) bool
		CanUpdatePage(context.Context, *types.Page) bool
		CanDeletePage(context.Context, *types.Page) bool
	}

	PageService interface {
		With(ctx context.Context) PageService

		FindByID(namespaceID, pageID uint64) (*types.Page, error)
		FindByHandle(namespaceID uint64, handle string) (*types.Page, error)
		FindByPageID(namespaceID, pageID uint64) (*types.Page, error)
		FindBySelfID(namespaceID, selfID uint64) (pages types.PageSet, f types.PageFilter, err error)
		Find(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error)
		Tree(namespaceID uint64) (pages types.PageSet, err error)

		Create(page *types.Page) (*types.Page, error)
		Update(page *types.Page) (*types.Page, error)
		DeleteByID(namespaceID, pageID uint64) error

		Reorder(namespaceID, selfID uint64, pageIDs []uint64) error
	}

	pageUpdateHandler func(ctx context.Context, ns *types.Namespace, c *types.Page) (bool, error)
)

func Page() PageService {
	return (&page{
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),
		store:    DefaultStore,
	}).With(context.Background())
}

func (svc page) With(ctx context.Context) PageService {
	return &page{
		ctx:       ctx,
		actionlog: DefaultActionlog,
		ac:        svc.ac,
		eventbus:  svc.eventbus,
		store:     svc.store,
	}
}

func (svc page) FindByID(namespaceID, pageID uint64) (p *types.Page, err error) {
	return svc.lookup(namespaceID, func(aProps *pageActionProps) (*types.Page, error) {
		if pageID == 0 {
			return nil, PageErrInvalidID()
		}

		aProps.page.ID = pageID
		return store.LookupComposePageByID(svc.ctx, svc.store, pageID)
	})
}

func (svc page) FindByHandle(namespaceID uint64, h string) (c *types.Page, err error) {
	return svc.lookup(namespaceID, func(aProps *pageActionProps) (*types.Page, error) {
		if !handle.IsValid(h) {
			return nil, PageErrInvalidHandle()
		}

		aProps.page.Handle = h
		return store.LookupComposePageByNamespaceIDHandle(svc.ctx, svc.store, namespaceID, h)
	})
}

func (svc page) FindByPageID(namespaceID, pageID uint64) (p *types.Page, err error) {
	return svc.lookup(namespaceID, func(aProps *pageActionProps) (*types.Page, error) {
		if pageID == 0 {
			return nil, PageErrInvalidID()
		}

		aProps.page.ID = pageID
		return store.LookupComposePageByID(svc.ctx, svc.store, pageID)
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
func (svc page) search(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error) {
	var (
		aProps = &pageActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = checkPage(svc.ctx, svc.ac)

	err = func() error {
		if ns, err := loadNamespace(svc.ctx, svc.store, filter.NamespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if set, f, err = store.SearchComposePages(svc.ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(svc.ctx, aProps, PageActionSearch, err)
}

func (svc page) FindBySelfID(namespaceID, parentID uint64) (pp types.PageSet, f types.PageFilter, err error) {
	return svc.search(types.PageFilter{
		NamespaceID: namespaceID,
		ParentID:    parentID,

		// This will enable parentID=0 query
		Root: true,

		Check: checkPage(svc.ctx, svc.ac),
	})
}

func (svc page) Find(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error) {
	return svc.search(filter)
}

func (svc page) Tree(namespaceID uint64) (tree types.PageSet, err error) {
	var (
		pages  types.PageSet
		filter = types.PageFilter{
			NamespaceID: namespaceID,
			Check:       checkPage(svc.ctx, svc.ac),
		}
	)

	if err = filter.Sort.Set("weight ASC"); err != nil {
		return
	}

	if pages, _, err = svc.search(filter); err != nil {
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
//
func (svc page) Reorder(namespaceID, parentID uint64, pageIDs []uint64) (err error) {
	var (
		aProps = &pageActionProps{page: &types.Page{ID: parentID}}
		ns     *types.Namespace
		p      *types.Page
	)

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if ns, err = loadNamespace(ctx, s, namespaceID); err != nil {
			return err
		}

		if parentID == 0 {
			// Reordering on root mode -- check if user can create pages.
			if !svc.ac.CanCreatePage(svc.ctx, ns) {
				return PageErrNotAllowedToUpdate()
			}
		} else {
			// Validate permissions on parent page
			if p, err = store.LookupComposePageByID(ctx, s, parentID); errors.Is(err, store.ErrNotFound) {
				return PageErrNotFound()
			} else if err != nil {
				return err
			}

			aProps.setPage(p)

			if !svc.ac.CanUpdatePage(svc.ctx, p) {
				return PageErrNotAllowedToUpdate()
			}
		}

		return store.ReorderComposePages(ctx, s, namespaceID, parentID, pageIDs)
	})

	return svc.recordAction(svc.ctx, aProps, PageActionReorder, err)

}

func (svc page) Create(new *types.Page) (p *types.Page, err error) {
	var (
		ns     *types.Namespace
		aProps = &pageActionProps{changed: new}
	)

	new.ID = 0

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if !handle.IsValid(new.Handle) {
			return PageErrInvalidID()
		}

		if ns, err = loadNamespace(ctx, s, new.NamespaceID); err != nil {
			return err
		}

		if !svc.ac.CanCreatePage(svc.ctx, ns) {
			return PageErrNotAllowedToCreate()
		}

		aProps.setNamespace(ns)

		if err = svc.eventbus.WaitFor(svc.ctx, event.PageBeforeCreate(new, nil, ns)); err != nil {
			return err
		}

		if err = svc.uniqueCheck(new); err != nil {
			return err
		}

		if err = store.CreateComposePage(ctx, s, new); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.PageAfterCreate(new, nil, ns))
		return err
	})

	return p, svc.recordAction(svc.ctx, aProps, PageActionCreate, err)
}

func (svc page) Update(upd *types.Page) (c *types.Page, err error) {
	return svc.updater(upd.NamespaceID, upd.ID, PageActionUpdate, svc.handleUpdate(upd))
}

func (svc page) DeleteByID(namespaceID, pageID uint64) error {
	return trim1st(svc.updater(namespaceID, pageID, PageActionDelete, svc.handleDelete))
}

func (svc page) UndeleteByID(namespaceID, pageID uint64) error {
	return trim1st(svc.updater(namespaceID, pageID, PageActionUndelete, svc.handleUndelete))
}

func (svc page) updater(namespaceID, pageID uint64, action func(...*pageActionProps) *pageAction, fn pageUpdateHandler) (*types.Page, error) {
	var (
		changed bool

		ns     *types.Namespace
		p, old *types.Page
		aProps = &pageActionProps{page: &types.Page{ID: pageID, NamespaceID: namespaceID}}
		err    error
	)

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, p, err = loadPage(svc.ctx, s, namespaceID, pageID)
		if err != nil {
			return
		}

		old = p.Clone()

		aProps.setNamespace(ns)
		aProps.setChanged(p)

		if p.DeletedAt == nil {
			err = svc.eventbus.WaitFor(svc.ctx, event.PageBeforeUpdate(p, old, ns))
		} else {
			err = svc.eventbus.WaitFor(svc.ctx, event.PageBeforeDelete(p, old, ns))
		}

		if err != nil {
			return
		}

		if changed, err = fn(svc.ctx, ns, p); err != nil {
			return err
		}

		if changed {
			if err = store.UpdateComposePage(svc.ctx, svc.store, p); err != nil {
				return err
			}
		}

		if p.DeletedAt == nil {
			err = svc.eventbus.WaitFor(svc.ctx, event.PageAfterUpdate(p, old, ns))
		} else {
			err = svc.eventbus.WaitFor(svc.ctx, event.PageAfterDelete(nil, old, ns))
		}

		return err
	})

	return p, svc.recordAction(svc.ctx, aProps, action, err)
}

// lookup fn() orchestrates page lookup, namespace preload and check
func (svc page) lookup(namespaceID uint64, lookup func(*pageActionProps) (*types.Page, error)) (p *types.Page, err error) {
	var aProps = &pageActionProps{page: &types.Page{NamespaceID: namespaceID}}

	err = func() error {
		if ns, err := loadNamespace(svc.ctx, svc.store, namespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if p, err = lookup(aProps); errors.Is(err, store.ErrNotFound) {
			return PageErrNotFound()
		} else if err != nil {
			return err
		}

		aProps.setPage(p)

		if !svc.ac.CanReadPage(svc.ctx, p) {
			return PageErrNotAllowedToRead()
		}

		return nil
	}()

	return p, svc.recordAction(svc.ctx, aProps, PageActionLookup, err)
}

func (svc page) uniqueCheck(p *types.Page) (err error) {
	if p.Handle != "" {
		if e, _ := store.LookupComposePageByNamespaceIDHandle(svc.ctx, svc.store, p.NamespaceID, p.Handle); e != nil && e.ID != p.ID {
			return PageErrHandleNotUnique()
		}
	}

	if p.ModuleID > 0 {
		if e, _ := store.LookupComposePageByNamespaceIDModuleID(svc.ctx, svc.store, p.NamespaceID, p.ModuleID); e != nil && e.ID != p.ID {
			return PageErrModuleNotFound()
		}
	}

	return nil
}

func (svc page) handleUpdate(upd *types.Page) pageUpdateHandler {
	return func(ctx context.Context, ns *types.Namespace, p *types.Page) (bool, error) {
		if isStale(upd.UpdatedAt, p.UpdatedAt, p.CreatedAt) {
			return false, PageErrStaleData()
		}

		if upd.Handle != p.Handle && !handle.IsValid(upd.Handle) {
			return false, PageErrInvalidHandle()
		}

		if err := svc.uniqueCheck(upd); err != nil {
			return false, err
		}

		if !svc.ac.CanUpdatePage(svc.ctx, p) {
			return false, PageErrNotAllowedToUpdate()
		}

		p.ID = upd.ID
		p.SelfID = upd.SelfID
		p.Blocks = upd.Blocks
		p.Title = upd.Title
		p.Handle = upd.Handle
		p.Description = upd.Description
		p.Visible = upd.Visible
		p.Weight = upd.Weight
		p.UpdatedAt = nowPtr()

		return true, nil
	}
}

func (svc page) handleDelete(ctx context.Context, ns *types.Namespace, m *types.Page) (bool, error) {
	if !svc.ac.CanDeletePage(ctx, m) {
		return false, PageErrNotAllowedToDelete()
	}

	if m.DeletedAt != nil {
		// page already deleted
		return false, nil
	}

	m.DeletedAt = nowPtr()
	return true, nil
}

func (svc page) handleUndelete(ctx context.Context, ns *types.Namespace, m *types.Page) (bool, error) {
	if !svc.ac.CanDeletePage(ctx, m) {
		return false, PageErrNotAllowedToUndelete()
	}

	if m.DeletedAt == nil {
		// page not deleted
		return false, nil
	}

	m.DeletedAt = nil
	return true, nil
}

func loadPage(ctx context.Context, s store.Storer, namespaceID, pageID uint64) (ns *types.Namespace, m *types.Page, err error) {
	if pageID == 0 {
		return nil, nil, PageErrInvalidID()
	}

	if ns, err = loadNamespace(ctx, s, namespaceID); err == nil {
		if m, err = store.LookupComposePageByID(ctx, s, pageID); errors.Is(err, store.ErrNotFound) {
			err = PageErrNotFound()
		}
	}

	if err != nil {
		return nil, nil, err
	}

	if namespaceID != m.NamespaceID {
		// Make sure chart belongs to the right namespace
		return nil, nil, PageErrNotFound()
	}

	return
}
