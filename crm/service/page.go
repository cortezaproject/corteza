package service

import (
	"context"
	"errors"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	page struct {
		db  *factory.DB
		ctx context.Context

		pageRepo repository.PageRepository
		moduleRepo repository.ModuleRepository
	}

	PageService interface {
		With(ctx context.Context) PageService

		FindByID(pageID uint64) (*types.Page, error)
		FindByModuleID(moduleID uint64) (*types.Page, error)
		FindBySelfID(selfID uint64) (pages types.PageSet, err error)
		Tree() (pages types.PageSet, err error)

		Create(page *types.Page) (*types.Page, error)
		Update(page *types.Page) (*types.Page, error)
		DeleteByID(pageID uint64) error

		Reorder(selfID uint64, pageIDs []uint64) error
	}
)

func Page() PageService {
	return (&page{}).With(context.Background())
}

func (s *page) With(ctx context.Context) PageService {
	db := repository.DB(ctx)
	return &page{
		db:        db,
		ctx:       ctx,
		pageRepo:  repository.Page(ctx, db),
		moduleRepo: repository.Module(ctx, db),
	}
}

func (s *page) FindByID(id uint64) (*types.Page, error) {
	page, err := s.pageRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := s.preload(page); err != nil {
		return nil, err
	}
	return page, err
}

func (s *page) FindByModuleID(moduleID uint64) (*types.Page, error) {
	page, err := s.pageRepo.FindByModuleID(moduleID)
	if err != nil {
		return nil, err
	}
	if err := s.preload(page); err != nil {
		return nil, err
	}
	return page, err
}

func (s *page) FindBySelfID(selfID uint64) (pages types.PageSet, err error) {
	return pages, s.db.Transaction(func() (err error) {
		if pages, err = s.pageRepo.FindBySelfID(selfID); err != nil {
			return
		}

		if err = s.preloadAll(pages); err != nil {
			return
		}

		return nil
	})
}

func (s *page) Tree() (pages types.PageSet, err error) {
	var tree types.PageSet

	return tree, s.db.Transaction(func() (err error) {
		if pages, err = s.pageRepo.FindAll(); err != nil {
			return
		}

		if err = s.preloadAll(pages); err != nil {
			return
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
				// Ignore orphans :(
			}

			return nil
		})

		return nil
	})
}

func (s *page) Reorder(selfID uint64, pageIDs []uint64) error {
	return s.pageRepo.Reorder(selfID, pageIDs)
}

func (s *page) Create(page *types.Page) (p *types.Page, err error) {
	validate := func() error {
		if page.ModuleID > 0 {
			// @todo check if module exists!
			if p, err = s.pageRepo.FindByModuleID(page.ModuleID); err != nil {
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
	return p, s.db.Transaction(func() (err error) {
		p, err = s.pageRepo.Create(page)
		return
	})
}

func (s *page) Update(page *types.Page) (p *types.Page, err error) {
	validate := func() error {
		if page.ID == 0 {
			return errors.New("Error when savig page, invalid ID")
		}
		if page.ModuleID > 0 {
			// @todo check if module exists!
			if p, err = s.pageRepo.FindByModuleID(page.ModuleID); err != nil {
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
	return p, s.db.Transaction(func() (err error) {
		p, err = s.pageRepo.Update(page)
		return
	})
}

func (s *page) DeleteByID(id uint64) error {
	return s.pageRepo.DeleteByID(id)
}
