package repository

import (
	"context"

	"github.com/titpetric/factory"

	_ "github.com/crusttech/crust/crm/types"
)

type (
	WorkflowRepository interface {
		With(ctx context.Context, db *factory.DB) WorkflowRepository
	}

	workflow struct {
		*repository
	}
)

func Workflow(ctx context.Context, db *factory.DB) WorkflowRepository {
	return (&workflow{}).With(ctx, db)
}

func (r *workflow) With(ctx context.Context, db *factory.DB) WorkflowRepository {
	return &workflow{
		repository: r.repository.With(ctx, db),
	}
}
