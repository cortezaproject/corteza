package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/repository"
	_ "github.com/crusttech/crust/crm/types"
)

type (
	job struct {
		db  *factory.DB
		ctx context.Context

		repository repository.JobRepository
	}

	JobService interface {
		With(ctx context.Context) JobService
	}
)

func Job() JobService {
	return (&job{}).With(context.Background())
}

func (s *job) With(ctx context.Context) JobService {
	db := repository.DB(ctx)
	return &job{
		db:         db,
		ctx:        ctx,
		repository: repository.Job(ctx, db),
	}
}
