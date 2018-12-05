package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/repository"
	_ "github.com/crusttech/crust/crm/types"
)

type (
	workflow struct {
		db  *factory.DB
		ctx context.Context

		repository repository.WorkflowRepository
	}

	WorkflowService interface {
		With(ctx context.Context) WorkflowService
	}
)

func Workflow() WorkflowService {
	return (&workflow{}).With(context.Background())
}

func (s *workflow) With(ctx context.Context) WorkflowService {
	db := repository.DB(ctx)
	return &workflow{
		db:         db,
		ctx:        ctx,
		repository: repository.Workflow(ctx, db),
	}
}
