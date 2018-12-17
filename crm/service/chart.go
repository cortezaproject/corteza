package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	chart struct {
		db  *factory.DB
		ctx context.Context

		chartRepo repository.ChartRepository
	}

	ChartService interface {
		With(ctx context.Context) ChartService

		FindByID(chartID uint64) (*types.Chart, error)
		Find() (types.ChartSet, error)

		Create(chart *types.Chart) (*types.Chart, error)
		Update(chart *types.Chart) (*types.Chart, error)
		DeleteByID(chartID uint64) error
	}
)

func Chart() ChartService {
	return (&chart{}).With(context.Background())
}

func (svc *chart) With(ctx context.Context) ChartService {
	db := repository.DB(ctx)
	return &chart{
		db:        db,
		ctx:       ctx,
		chartRepo: repository.Chart(ctx, db),
	}
}

func (svc *chart) FindByID(chartID uint64) (*types.Chart, error) {
	return svc.chartRepo.FindByID(chartID)
}

func (svc *chart) Find() (types.ChartSet, error) {
	return svc.chartRepo.Find()
}

func (svc *chart) Create(mod *types.Chart) (c *types.Chart, err error) {
	return c, svc.db.Transaction(func() error {
		c, err = svc.chartRepo.Create(mod)
		return err
	})
}

func (svc *chart) Update(mod *types.Chart) (c *types.Chart, err error) {
	validate := func() error {
		if mod.ID == 0 {
			return errors.New("Error updating chart: invalid ID")
		} else if c, err = svc.chartRepo.FindByID(mod.ID); err != nil {
			return errors.Wrap(err, "Error while loading chart for update")
		} else {
			mod.CreatedAt = c.CreatedAt
		}

		return nil
	}

	if err = validate(); err != nil {
		return nil, err
	}

	c.Config = mod.Config
	c.Name = mod.Name

	return c, svc.db.Transaction(func() error {
		c, err = svc.chartRepo.Update(c)
		return err
	})
}

func (svc *chart) DeleteByID(chartID uint64) error {
	return svc.chartRepo.DeleteByID(chartID)
}
