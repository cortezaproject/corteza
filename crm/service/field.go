package service

import (
	"context"
	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	field struct {
		db         *factory.DB
		ctx        context.Context
		repository repository.FieldRepository
	}

	FieldService interface {
		With(ctx context.Context) FieldService
		FindByType(t string) (*types.Field, error)
		Find() ([]*types.Field, error)
	}
)

func Field() FieldService {
	return (&field{}).With(context.Background())
}

func (s *field) With(ctx context.Context) FieldService {
	db := repository.DB(ctx)
	return &field{
		db:         db,
		ctx:        ctx,
		repository: s.repository.With(ctx, db),
	}
}

func (s *field) FindByType(t string) (*types.Field, error) {
	return s.repository.FindByType(t)
}

func (s *field) Find() ([]*types.Field, error) {
	return s.repository.Find()
}
