package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/compose/internal/repository"
	"github.com/crusttech/crust/compose/types"
)

type (
	page struct {
		db  *factory.DB
		ctx context.Context

		prmSvc PermissionsService

		pageRepo   repository.PageRepository
		moduleRepo repository.ModuleRepository
	}

	PageService interface {
		With(ctx context.Context) PageService

		FindByID(namespaceID, pageID uint64) (*types.Page, error)
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
		prmSvc: DefaultPermissions,
	}).With(context.Background())
}

func (svc *page) With(ctx context.Context) PageService {
	db := repository.DB(ctx)
	return &page{
		db:  db,
		ctx: ctx,

		prmSvc: svc.prmSvc.With(ctx),

		pageRepo:   repository.Page(ctx, db),
		moduleRepo: repository.Module(ctx, db),
	}
}

func (svc *page) FindByID(namespaceID, pageID uint64) (p *types.Page, err error) {
	return svc.checkPermissions(svc.pageRepo.FindByID(namespaceID, pageID))
}

func (svc *page) FindByModuleID(namespaceID, moduleID uint64) (p *types.Page, err error) {
	return svc.checkPermissions(svc.pageRepo.FindByModuleID(namespaceID, moduleID))
}

func (svc *page) checkPermissions(p *types.Page, err error) (*types.Page, error) {
	if err != nil {
		return nil, err
	} else if !svc.prmSvc.CanReadPage(p) {
		return nil, errors.New("not allowed to access this page")
	}

	return p, err
}

func (svc *page) FindBySelfID(namespaceID, parentID uint64) (pp types.PageSet, f types.PageFilter, err error) {
	return svc.filterPageSetByPermission(svc.pageRepo.Find(types.PageFilter{
		NamespaceID: namespaceID,
		ParentID:    parentID,
	}))
}

func (svc *page) Find(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error) {
	return svc.filterPageSetByPermission(svc.pageRepo.Find(filter))
}

func (svc *page) Tree(namespaceID uint64) (pages types.PageSet, err error) {
	var (
		tree   types.PageSet
		filter = types.PageFilter{
			NamespaceID: namespaceID,
		}
	)

	return tree, svc.db.Transaction(func() (err error) {
		if pages, _, err = svc.filterPageSetByPermission(svc.pageRepo.Find(filter)); err != nil {
			return
		}

		// No preloading - we do not need (or should have) any modules
		// associated with us
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

		return nil
	})
}

func (svc *page) filterPageSetByPermission(pp types.PageSet, f types.PageFilter, err error) (types.PageSet, types.PageFilter, error) {
	if err != nil {
		return nil, f, err
	}

	// @todo Filter-by-permission can/will mess up filter's count & paging...
	pp, err = pp.Filter(func(m *types.Page) (bool, error) {
		return svc.prmSvc.CanReadPage(m), nil
	})

	return pp, f, err
}

func (svc *page) Reorder(namespaceID, selfID uint64, pageIDs []uint64) error {
	return svc.pageRepo.Reorder(namespaceID, selfID, pageIDs)
}

func (svc *page) Create(page *types.Page) (p *types.Page, err error) {
	validate := func() error {
		if !svc.prmSvc.CanCreatePage(crmNamespace()) {
			return errors.New("not allowed to create this page")
		}

		if page.ModuleID > 0 {
			if p, err = svc.pageRepo.FindByModuleID(page.NamespaceID, page.ModuleID); err != nil {
				return err
			} else if p.ID > 0 {
				return errors.New("Page for module already exists")
			}
		}
		return nil
	}
	if err := validate(); err != nil {
		return nil, err
	}
	return p, svc.db.Transaction(func() (err error) {
		p, err = svc.pageRepo.Create(page)
		return
	})
}

func (svc *page) Update(page *types.Page) (p *types.Page, err error) {
	validate := func() error {
		if page.ID == 0 {
			return errors.New("Error when saving page, invalid ID")
		} else if p, err = svc.pageRepo.FindByID(page.NamespaceID, page.ID); err != nil {
			return errors.Wrap(err, "Error while loading page for update")
		} else {
			if !svc.prmSvc.CanUpdatePage(p) {
				return errors.New("not allowed to update this page")
			}
		}

		if page.ModuleID > 0 {
			if p, err = svc.pageRepo.FindByModuleID(page.NamespaceID, page.ModuleID); err != nil {
				return err
			} else if p.ID > 0 && page.ID != p.ID {
				return errors.New("Page for module already exists")
			}
		}
		return nil
	}
	if err := validate(); err != nil {
		return nil, err
	}
	return p, svc.db.Transaction(func() (err error) {
		p, err = svc.pageRepo.Update(page)
		return
	})
}

func (svc *page) DeleteByID(namespaceID, pageID uint64) error {
	if p, err := svc.pageRepo.FindByID(namespaceID, pageID); err != nil {
		return errors.Wrap(err, "could not delete page")
	} else if !svc.prmSvc.CanDeletePage(p) {
		return errors.New("not allowed to delete this page")
	}

	return svc.pageRepo.DeleteByID(namespaceID, pageID)
}
