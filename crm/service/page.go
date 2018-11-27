package service

import (
	"context"
	"errors"

	"github.com/davecgh/go-spew/spew"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	page struct {
		db         *factory.DB
		ctx        context.Context
		repository repository.PageRepository
		moduleRepo repository.ModuleRepository
	}

	PageService interface {
		With(ctx context.Context) PageService

		FindByID(pageID uint64) (*types.Page, error)
		Find(selfID uint64) (pages types.PageSet, err error)
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
		db:         db,
		ctx:        ctx,
		repository: repository.Page(ctx, db),
		moduleRepo: repository.Module(ctx, db),
	}
}

func (s *page) FindByID(id uint64) (*types.Page, error) {
	return s.repository.FindByID(id)
}

func (s *page) Find(selfID uint64) (pages types.PageSet, err error) {
	return pages, s.db.Transaction(func() (err error) {
		if pages, err = s.repository.FindBySelfID(selfID); err != nil {
			return
		}

		if err = s.preload(pages); err != nil {
			return
		}

		return nil
	})
}

func (s *page) Tree() (pages types.PageSet, err error) {
	var tree types.PageSet

	return tree, s.db.Transaction(func() (err error) {
		if pages, err = s.repository.FindAll(); err != nil {
			return
		}

		if err = s.preload(pages); err != nil {
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
	return s.repository.Reorder(selfID, pageIDs)
}

func (s *page) Create(mod *types.Page) (p *types.Page, err error) {
	return p, s.db.Transaction(func() (err error) {
		if mod.ModuleID > 0 {
			// @todo check if module exists!
			if p, err = s.repository.FindByModuleID(mod.ModuleID); err != nil {
				return err
			} else if p.ID > 0 {
				return errors.New("Page for module already exists")
			}
		}

		p, err = s.repository.Create(mod)
		return
	})
}

func (s *page) Update(mod *types.Page) (p *types.Page, err error) {
	return p, s.db.Transaction(func() (err error) {
		if mod.ModuleID > 0 {
			// @todo check if module exists!
			if p, err = s.repository.FindByModuleID(mod.ModuleID); err != nil {
				return err
			} else if p.ID > 0 && mod.ID != p.ID {
				spew.Dump(mod, p)
				return errors.New("Page for module already exists")
			}
		}

		p, err = s.repository.Update(mod)
		return
	})
}

func (s *page) DeleteByID(id uint64) error {
	return s.repository.DeleteByID(id)
}

// Preloads modules for all pages
func (s *page) preload(pages types.PageSet) error {
	if modules, err := s.moduleRepo.Find(); err != nil {
		return err
	} else {
		_ = pages.Walk(func(page *types.Page) error {
			page.Module = modules.FindByID(page.ModuleID)
			return nil
		})
	}

	return nil
}
