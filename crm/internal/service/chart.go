package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/internal/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	chart struct {
		db  *factory.DB
		ctx context.Context

		prmSvc PermissionsService

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
	return (&chart{
		prmSvc: DefaultPermissions,
	}).With(context.Background())
}

func (svc *chart) With(ctx context.Context) ChartService {
	db := repository.DB(ctx)
	return &chart{
		db:  db,
		ctx: ctx,

		prmSvc: svc.prmSvc.With(ctx),

		chartRepo: repository.Chart(ctx, db),
	}
}

func (svc *chart) FindByID(chartID uint64) (c *types.Chart, err error) {
	if c, err = svc.chartRepo.FindByID(chartID); err != nil {
		return
	} else if !svc.prmSvc.CanReadChart(c) {
		return nil, errors.New("not allowed to access this chart")
	}

	return
}

func (svc *chart) Find() (cc types.ChartSet, err error) {
	if cc, err = svc.chartRepo.Find(); err != nil {
		return nil, err
	} else {
		return cc.Filter(func(m *types.Chart) (bool, error) {
			return svc.prmSvc.CanReadChart(m), nil
		})
	}
}

func (svc *chart) Create(mod *types.Chart) (c *types.Chart, err error) {
	if !svc.prmSvc.CanCreateChart() {
		return nil, errors.New("not allowed to create this chart")
	}

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
			if !svc.prmSvc.CanUpdateChart(c) {
				return errors.New("not allowed to update this chart")
			}

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

func (svc *chart) DeleteByID(ID uint64) error {
	if !svc.prmSvc.CanDeleteChartByID(ID) {
		return errors.New("not allowed to delete this chart")
	}

	return svc.chartRepo.DeleteByID(ID)
}
