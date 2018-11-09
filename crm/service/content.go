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

		Find(moduleID uint64, query string, page int, perPage int) (*repository.FindResponse, error)

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
		repository: repository.Content(ctx, db),
	}
}

func (s *content) FindByID(id uint64) (*types.Content, error) {
	return s.repository.FindByID(id)
}

func (s *content) Find(moduleID uint64, query string, page int, perPage int) (*repository.FindResponse, error) {
	return s.repository.Find(moduleID, query, page, perPage)
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
