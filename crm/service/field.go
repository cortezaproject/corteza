package service

import (
	"context"
	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	field struct {
		repository repository.Field
	}

	FieldService interface {
		With(ctx context.Context) FieldService
		FindByType(t string) (*types.Field, error)
		Find() ([]*types.Field, error)
	}
)

func Field() FieldService {
	return &field{
		repository: repository.NewField(context.Background()),
	}
}

func (s *field) With(ctx context.Context) FieldService {
	return &field{
		repository: s.repository.With(ctx),
	}
}

func (s *field) FindByType(t string) (*types.Field, error) {
	return s.repository.FindByType(t)
}

func (s *field) Find() ([]*types.Field, error) {
	return s.repository.Find()
}
