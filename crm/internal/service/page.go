package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
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

		FindByID(pageID uint64) (*types.Page, error)
		FindByModuleID(moduleID uint64) (*types.Page, error)
		FindBySelfID(selfID uint64) (pages types.PageSet, err error)
		Find() (pages types.PageSet, err error)
		Tree() (pages types.PageSet, err error)
		FindRecordPages() (pages types.PageSet, err error)

		Create(page *types.Page) (*types.Page, error)
		Update(page *types.Page) (*types.Page, error)
		DeleteByID(pageID uint64) error

		Reorder(selfID uint64, pageIDs []uint64) error
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

func (svc *page) FindByID(id uint64) (p *types.Page, err error) {
	return svc.checkPermissions(svc.pageRepo.FindByID(id))
}

func (svc *page) FindByModuleID(moduleID uint64) (p *types.Page, err error) {
	return svc.checkPermissions(svc.pageRepo.FindByModuleID(moduleID))
}

func (svc *page) checkPermissions(p *types.Page, err error) (*types.Page, error) {
	if err != nil {
		return nil, err
	} else if !svc.prmSvc.CanReadPage(p) {
		return nil, errors.New("not allowed to access this page")
	}

	return p, err
}

func (svc *page) FindBySelfID(selfID uint64) (pp types.PageSet, err error) {
	return svc.filterPageSet(svc.pageRepo.FindBySelfID(selfID))
}

func (svc *page) Find() (pages types.PageSet, err error) {
	return svc.filterPageSet(svc.pageRepo.Find())
}

func (svc *page) Tree() (pages types.PageSet, err error) {
	var tree types.PageSet

	return tree, svc.db.Transaction(func() (err error) {
		if pages, err = svc.filterPageSet(svc.pageRepo.Find()); err != nil {
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

func (svc *page) FindRecordPages() (pages types.PageSet, err error) {
	return svc.pageRepo.FindRecordPages()
}

func (svc *page) filterPageSet(pp types.PageSet, err error) (types.PageSet, error) {
	if err != nil {
		return nil, err
	}

	return pp.Filter(func(m *types.Page) (bool, error) {
		return svc.prmSvc.CanReadPage(m), nil
	})
}

func (svc *page) Reorder(selfID uint64, pageIDs []uint64) error {
	return svc.pageRepo.Reorder(selfID, pageIDs)
}

func (svc *page) Create(page *types.Page) (p *types.Page, err error) {
	validate := func() error {
		if !svc.prmSvc.CanCreatePage() {
			return errors.New("not allowed to create this module")
		}

		if page.ModuleID > 0 {
			if p, err = svc.pageRepo.FindByModuleID(page.ModuleID); err != nil {
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
			return errors.New("Error when savig page, invalid ID")
		} else if p, err = svc.pageRepo.FindByID(page.ID); err != nil {
			return errors.Wrap(err, "Error while loading module for update")
		} else {
			if !svc.prmSvc.CanUpdatePage(p) {
				return errors.New("not allowed to update this pahe")
			}
		}

		if page.ModuleID > 0 {
			if p, err = svc.pageRepo.FindByModuleID(page.ModuleID); err != nil {
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

func (svc *page) DeleteByID(ID uint64) error {
	if !svc.prmSvc.CanDeletePageByID(ID) {
		return errors.New("not allowed to delete this page")
	}

	return svc.pageRepo.DeleteByID(ID)
}
