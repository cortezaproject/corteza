package service

import (
	"context"
	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	fieldType struct {
		repository fieldTypeRepository
	}

	fieldTypeRepository interface {
		FindByName(ctx context.Context, name string) (*types.Field, error)
		Find(ctx context.Context) ([]*types.Field, error)
	}
)

func Field() fieldType {
	return fieldType{
		repository: repository.Field(),
	}
}

func (svc fieldType) FindByName(ctx context.Context, name string) (*types.Field, error) {
	return svc.repository.FindByName(ctx, name)
}

func (svc fieldType) Find(ctx context.Context) ([]*types.Field, error) {
	return svc.repository.Find(ctx)
}
