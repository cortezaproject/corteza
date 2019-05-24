package service

import (
	"context"

	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/internal/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	chart struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac chartAccessController

		chartRepo repository.ChartRepository
		nsRepo    repository.NamespaceRepository
	}

	chartAccessController interface {
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanCreateChart(context.Context, *types.Namespace) bool
		CanReadChart(context.Context, *types.Chart) bool
		CanUpdateChart(context.Context, *types.Chart) bool
		CanDeleteChart(context.Context, *types.Chart) bool
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
		logger: DefaultLogger.Named("chart"),
		ac:     DefaultAccessControl,
	}).With(context.Background())
}

func (svc chart) With(ctx context.Context) ChartService {
	db := repository.DB(ctx)
	return &chart{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		ac: svc.ac,

		chartRepo: repository.Chart(ctx, db),
		nsRepo:    repository.Namespace(ctx, db),
	}
}

// log() returns zap's logger with requestID from current context and fields.
// func (svc chart) log(fields ...zapcore.Field) *zap.Logger {
// 	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
// }

func (svc chart) FindByID(namespaceID, chartID uint64) (c *types.Chart, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired
	}

	if c, err = svc.chartRepo.FindByID(namespaceID, chartID); err != nil {
		return
	} else if !svc.ac.CanReadChart(svc.ctx, c) {
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
		return svc.ac.CanReadChart(svc.ctx, m), nil
	})

	return
}

func (svc chart) Create(mod *types.Chart) (c *types.Chart, err error) {
	if ns, err := svc.loadNamespace(mod.NamespaceID); err != nil {
		return nil, err
	} else if !svc.ac.CanCreateChart(svc.ctx, ns) {
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

	if !svc.ac.CanUpdateChart(svc.ctx, c) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	c.Config = mod.Config
	c.Name = mod.Name

	return svc.chartRepo.Update(c)
}

func (svc chart) DeleteByID(namespaceID, chartID uint64) error {
	if chartID == 0 {
		return ErrInvalidID.withStack()
	}

	if namespaceID == 0 {
		return ErrNamespaceRequired.withStack()
	}

	if c, err := svc.chartRepo.FindByID(namespaceID, chartID); err != nil {
		return err
	} else if !svc.ac.CanDeleteChart(svc.ctx, c) {
		return ErrNoDeletePermissions.withStack()
	}

	return svc.chartRepo.DeleteByID(namespaceID, chartID)
}

func (svc chart) loadNamespace(namespaceID uint64) (ns *types.Namespace, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired.withStack()
	}

	if ns, err = svc.nsRepo.FindByID(namespaceID); err != nil {
		return
	}

	if !svc.ac.CanReadNamespace(svc.ctx, ns) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return
}
