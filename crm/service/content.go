package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	content struct {
		db         *factory.DB
		ctx        context.Context
		repository repository.ContentRepository
	}

	ContentService interface {
		With(ctx context.Context) ContentService

		FindByID(contentID uint64) (*types.Content, error)
		Find() ([]*types.Content, error)

		Create(content *types.Content) (*types.Content, error)
		Update(content *types.Content) (*types.Content, error)
		DeleteByID(contentID uint64) error
	}
)

func Content() ContentService {
	return (&content{}).With(context.Background())
}

func (s *content) With(ctx context.Context) ContentService {
	db := repository.DB(ctx)
	return &content{
		db:         db,
		ctx:        ctx,
		repository: s.repository.With(ctx, db),
	}
}

func (s *content) FindByID(id uint64) (*types.Content, error) {
	return s.repository.FindByID(id)
}

func (s *content) Find() ([]*types.Content, error) {
	return s.repository.Find()
}

func (s *content) Create(mod *types.Content) (*types.Content, error) {
	return s.repository.Create(mod)
}

func (s *content) Update(mod *types.Content) (*types.Content, error) {
	return s.repository.Update(mod)
}

func (s *content) DeleteByID(id uint64) error {
	return s.repository.DeleteByID(id)
}
