package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/compose/internal/repository"
	"github.com/crusttech/crust/compose/types"
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

		FindByID(namespaceID, chartID uint64) (*types.Chart, error)
		Find(filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error)

		Create(chart *types.Chart) (*types.Chart, error)
		Update(chart *types.Chart) (*types.Chart, error)
		DeleteByID(namespaceID, chartID uint64) error
	}
)

func Chart() ChartService {
	return (&chart{
		prmSvc: DefaultPermissions,
	}).With(context.Background())
}

func (svc chart) With(ctx context.Context) ChartService {
	db := repository.DB(ctx)
	return &chart{
		db:  db,
		ctx: ctx,

		prmSvc: svc.prmSvc.With(ctx),

		chartRepo: repository.Chart(ctx, db),
	}
}

func (svc chart) FindByID(namespaceID, chartID uint64) (c *types.Chart, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired
	}

	if c, err = svc.chartRepo.FindByID(namespaceID, chartID); err != nil {
		return
	} else if !svc.prmSvc.CanReadChart(c) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return
}

func (svc chart) Find(filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error) {
	set, f, err = svc.chartRepo.Find(filter)
	if err != nil {
		return
	}

	set, _ = set.Filter(func(m *types.Chart) (bool, error) {
		return svc.prmSvc.CanReadChart(m), nil
	})

	return
}

func (svc chart) Create(mod *types.Chart) (c *types.Chart, err error) {
	if !svc.prmSvc.CanCreateChart(crmNamespace()) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	return svc.chartRepo.Create(mod)
}

func (svc chart) Update(mod *types.Chart) (c *types.Chart, err error) {
	if mod.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if c, err = svc.chartRepo.FindByID(mod.NamespaceID, mod.ID); err != nil {
		return
	}

	if isStale(mod.UpdatedAt, c.UpdatedAt, c.CreatedAt) {
		return nil, ErrStaleData.withStack()
	}

	if !svc.prmSvc.CanUpdateChart(c) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	c.Config = mod.Config
	c.Name = mod.Name

	return svc.chartRepo.Update(c)
}

func (svc chart) DeleteByID(namespaceID, chartID uint64) error {
	if namespaceID == 0 {
		return ErrNamespaceRequired.withStack()
	}

	if c, err := svc.chartRepo.FindByID(namespaceID, chartID); err != nil {
		return err
	} else if !svc.prmSvc.CanDeleteChart(c) {
		return ErrNoDeletePermissions.withStack()
	}

	return svc.chartRepo.DeleteByID(namespaceID, chartID)
}
