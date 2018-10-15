package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	page struct {
		db         *factory.DB
		ctx        context.Context
		repository repository.PageRepository
	}

	PageService interface {
		With(ctx context.Context) PageService

		FindByID(pageID uint64) (*types.Page, error)
		Find() ([]*types.Page, error)

		Create(page *types.Page) (*types.Page, error)
		Update(page *types.Page) (*types.Page, error)
		DeleteByID(pageID uint64) error
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
	}
}

func (s *page) FindByID(id uint64) (*types.Page, error) {
	return s.repository.FindByID(id)
}

func (s *page) Find() ([]*types.Page, error) {
	return s.repository.Find()
}

func (s *page) Create(mod *types.Page) (*types.Page, error) {
	return s.repository.Create(mod)
}

func (s *page) Update(mod *types.Page) (*types.Page, error) {
	return s.repository.Update(mod)
}

func (s *page) DeleteByID(id uint64) error {
	return s.repository.DeleteByID(id)
}
