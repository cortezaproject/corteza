package repository

import (
	"context"

	"github.com/titpetric/factory"

	_ "github.com/crusttech/crust/crm/types"
)

type (
	JobRepository interface {
		With(ctx context.Context, db *factory.DB) JobRepository
	}

	job struct {
		*repository
	}
)

func Job(ctx context.Context, db *factory.DB) JobRepository {
	return (&job{}).With(ctx, db)
}

func (r *job) With(ctx context.Context, db *factory.DB) JobRepository {
	return &job{
		repository: r.repository.With(ctx, db),
	}
}
