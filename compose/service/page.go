package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	page struct {
		db  *factory.DB
		ctx context.Context

		actionlog actionlog.Recorder

		ac       pageAccessController
		eventbus eventDispatcher

		pageRepo   repository.PageRepository
		moduleRepo repository.ModuleRepository
		nsRepo     repository.NamespaceRepository
	}

	pageAccessController interface {
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanCreatePage(context.Context, *types.Namespace) bool
		CanReadPage(context.Context, *types.Page) bool
		CanUpdatePage(context.Context, *types.Page) bool
		CanDeletePage(context.Context, *types.Page) bool

		FilterReadablePages(ctx context.Context) *permissions.ResourceFilter
	}

	PageService interface {
		With(ctx context.Context) PageService

		FindByID(namespaceID, pageID uint64) (*types.Page, error)
		FindByHandle(namespaceID uint64, handle string) (*types.Page, error)
		FindByModuleID(namespaceID, moduleID uint64) (*types.Page, error)
		FindBySelfID(namespaceID, selfID uint64) (pages types.PageSet, f types.PageFilter, err error)
		Find(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error)
		Tree(namespaceID uint64) (pages types.PageSet, err error)

		Create(page *types.Page) (*types.Page, error)
		Update(page *types.Page) (*types.Page, error)
		DeleteByID(namespaceID, pageID uint64) error

		Reorder(namespaceID, selfID uint64, pageIDs []uint64) error
	}
)

func Page() PageService {
	return (&page{
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),
	}).With(context.Background())
}

func (svc page) With(ctx context.Context) PageService {
	db := repository.DB(ctx)
	return &page{
		db:  db,
		ctx: ctx,

		actionlog: DefaultActionlog,

		ac:       svc.ac,
		eventbus: svc.eventbus,

		pageRepo:   repository.Page(ctx, db),
		moduleRepo: repository.Module(ctx, db),
		nsRepo:     repository.Namespace(ctx, db),
	}
}

// lookup fn() orchestrates page lookup, namespace preload and check
func (svc page) lookup(namespaceID uint64, lookup func(*pageActionProps) (*types.Page, error)) (p *types.Page, err error) {
	var aProps = &pageActionProps{page: &types.Page{NamespaceID: namespaceID}}

	err = svc.db.Transaction(func() error {
		if ns, err := svc.loadNamespace(namespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if p, err = lookup(aProps); err != nil {
			if repository.ErrPageNotFound.Eq(err) {
				return PageErrNotFound()
			}

			return err
		}

		aProps.setPage(p)

		if !svc.ac.CanReadPage(svc.ctx, p) {
			return PageErrNotAllowedToRead()
		}

		return nil
	})

	return p, svc.recordAction(svc.ctx, aProps, PageActionLookup, err)
}

func (svc page) FindByID(namespaceID, pageID uint64) (p *types.Page, err error) {
	return svc.lookup(namespaceID, func(aProps *pageActionProps) (*types.Page, error) {
		if pageID == 0 {
			return nil, PageErrInvalidID()
		}

		aProps.page.ID = pageID
		return svc.pageRepo.FindByID(namespaceID, pageID)
	})
}

func (svc page) FindByHandle(namespaceID uint64, h string) (c *types.Page, err error) {
	return svc.lookup(namespaceID, func(aProps *pageActionProps) (*types.Page, error) {
		if !handle.IsValid(h) {
			return nil, PageErrInvalidHandle()
		}

		aProps.page.Handle = h
		return svc.pageRepo.FindByHandle(namespaceID, h)
	})
}

func (svc page) FindByModuleID(namespaceID, moduleID uint64) (p *types.Page, err error) {
	return svc.lookup(namespaceID, func(aProps *pageActionProps) (*types.Page, error) {
		if moduleID == 0 {
			return nil, PageErrInvalidID()
		}

		aProps.page.ModuleID = moduleID
		return svc.pageRepo.FindByModuleID(namespaceID, moduleID)
	})
}

// search fn() orchestrates pages search, namespace preload and check
func (svc page) search(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error) {
	var (
		aProps = &pageActionProps{filter: &filter}
	)

	f = filter
	f.IsReadable = svc.ac.FilterReadablePages(svc.ctx)

	err = svc.db.Transaction(func() error {
		if ns, err := svc.loadNamespace(f.NamespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if set, f, err = svc.pageRepo.Find(f); err != nil {
			return err
		}

		return nil
	})

	return set, f, svc.recordAction(svc.ctx, aProps, PageActionSearch, err)
}

func (svc page) FindBySelfID(namespaceID, parentID uint64) (pp types.PageSet, f types.PageFilter, err error) {
	return svc.search(types.PageFilter{
		NamespaceID: namespaceID,
		ParentID:    parentID,

		// This will enable parentID=0 query
		Root: true,

		IsReadable: svc.ac.FilterReadablePages(svc.ctx),
	})
}

func (svc page) Find(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error) {
	filter.IsReadable = svc.ac.FilterReadablePages(svc.ctx)
	return svc.search(filter)
}

func (svc page) Tree(namespaceID uint64) (tree types.PageSet, err error) {
	var pages types.PageSet
	pages, _, err = svc.search(types.PageFilter{
		NamespaceID: namespaceID,
		IsReadable:  svc.ac.FilterReadablePages(svc.ctx),
		Sort:        "weight ASC",
	})

	if err != nil {
		return nil, err
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

	err = svc.db.Transaction(func() (err error) {
		if ns, err = svc.loadNamespace(namespaceID); err != nil {
			return err
		}

		if parentID == 0 {
			// Reordering on root mode -- check if user can create pages.
			if !svc.ac.CanCreatePage(svc.ctx, ns) {
				return PageErrNotAllowedToUpdate()
			}
		} else {
			// Validate permissions on parent page
			if p, err = svc.pageRepo.FindByID(ns.ID, parentID); err != nil {
				if repository.ErrPageNotFound.Eq(err) {
					return PageErrNotFound()
				}

				return
			}

			aProps.setPage(p)

			if !svc.ac.CanUpdatePage(svc.ctx, p) {
				return PageErrNotAllowedToUpdate()
			}
		}

		return svc.pageRepo.Reorder(namespaceID, parentID, pageIDs)
	})

	return svc.recordAction(svc.ctx, aProps, PageActionReorder, err)

}

func (svc page) Create(new *types.Page) (p *types.Page, err error) {
	var (
		ns     *types.Namespace
		aProps = &pageActionProps{changed: new}
	)

	new.ID = 0

	err = svc.db.Transaction(func() error {
		if !handle.IsValid(new.Handle) {
			return PageErrInvalidID()
		}

		if ns, err = svc.loadNamespace(new.NamespaceID); err != nil {
			return err
		} else if !svc.ac.CanCreatePage(svc.ctx, ns) {
			return PageErrNotAllowedToCreate()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.PageBeforeCreate(new, nil, ns)); err != nil {
			return err
		}

		if err = svc.UniqueCheck(new); err != nil {
			return err
		}

		if p, err = svc.pageRepo.Create(new); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.PageAfterCreate(new, nil, ns))
		return err
	})

	return p, svc.recordAction(svc.ctx, aProps, PageActionCreate, err)
}

func (svc page) Update(upd *types.Page) (p *types.Page, err error) {
	var (
		ns     *types.Namespace
		aProps = &pageActionProps{changed: upd}
	)

	err = svc.db.Transaction(func() error {
		if upd.ID == 0 {
			return PageErrInvalidID()
		}

		if !handle.IsValid(upd.Handle) {
			return PageErrInvalidHandle()
		}

		if ns, err = svc.loadNamespace(upd.NamespaceID); err != nil {
			return err
		}

		if p, err = svc.pageRepo.FindByID(upd.NamespaceID, upd.ID); err != nil {
			if repository.ErrPageNotFound.Eq(err) {
				return PageErrNotFound()
			}

			return err
		}

		if isStale(upd.UpdatedAt, p.UpdatedAt, p.CreatedAt) {
			return PageErrStaleData()
		}

		if !svc.ac.CanUpdatePage(svc.ctx, p) {
			return PageErrNotAllowedToUpdate()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.PageBeforeUpdate(upd, p, ns)); err != nil {
			return err
		}

		if err = svc.UniqueCheck(upd); err != nil {
			return err
		}

		p.ModuleID = upd.ModuleID
		p.SelfID = upd.SelfID
		p.Blocks = upd.Blocks
		p.Title = upd.Title
		p.Handle = upd.Handle
		p.Description = upd.Description
		p.Visible = upd.Visible
		p.Weight = upd.Weight

		if p, err = svc.pageRepo.Update(p); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.PageAfterUpdate(upd, p, ns))
		return err
	})

	return p, svc.recordAction(svc.ctx, aProps, PageActionUpdate, err)
}

func (svc page) DeleteByID(namespaceID, pageID uint64) (err error) {
	var (
		del    *types.Page
		ns     *types.Namespace
		aProps = &pageActionProps{page: &types.Page{ID: pageID, NamespaceID: namespaceID}}
	)

	err = svc.db.Transaction(func() error {

		if pageID == 0 {
			return PageErrInvalidID()
		}

		if ns, err = svc.loadNamespace(namespaceID); err != nil {
			return err
		}

		if del, err = svc.pageRepo.FindByID(namespaceID, pageID); err != nil {
			if repository.ErrPageNotFound.Eq(err) {
				return PageErrNotFound()
			}

			return err
		}

		aProps.setChanged(del)

		if !svc.ac.CanDeletePage(svc.ctx, del) {
			return PageErrNotAllowedToDelete()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.PageBeforeDelete(nil, del, ns)); err != nil {
			return err
		}

		if err = svc.pageRepo.DeleteByID(namespaceID, pageID); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.PageAfterDelete(nil, del, ns))
		return err
	})

	return svc.recordAction(svc.ctx, aProps, PageActionDelete, err)
}

func (svc page) UniqueCheck(p *types.Page) (err error) {
	if p.Handle != "" {
		if e, _ := svc.pageRepo.FindByHandle(p.NamespaceID, p.Handle); e != nil && e.ID != p.ID {
			return PageErrHandleNotUnique()
		}
	}

	if p.ModuleID > 0 {
		if e, _ := svc.pageRepo.FindByModuleID(p.NamespaceID, p.ModuleID); e != nil && e.ID != p.ID {
			return PageErrModuleNotFound()
		}
	}

	return nil
}

func (svc page) loadNamespace(namespaceID uint64) (ns *types.Namespace, err error) {
	if namespaceID == 0 {
		return nil, PageErrInvalidNamespaceID()
	}

	if ns, err = svc.nsRepo.FindByID(namespaceID); err != nil {
		return
	}

	if !svc.ac.CanReadNamespace(svc.ctx, ns) {
		return nil, PageErrNotAllowedToReadNamespace()
	}

	return
}
