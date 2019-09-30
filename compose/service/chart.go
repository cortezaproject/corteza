package service

import (
	"context"

	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/logger"
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
		FindByHandle(namespaceID uint64, handle string) (*types.Chart, error)
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
func (svc chart) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc chart) FindByID(namespaceID, chartID uint64) (c *types.Chart, err error) {
	if _, err = svc.loadNamespace(namespaceID); err != nil {
		return
	} else {
		return svc.checkPermissions(svc.chartRepo.FindByID(namespaceID, chartID))
	}
}

func (svc chart) FindByHandle(namespaceID uint64, handle string) (c *types.Chart, err error) {
	if _, err = svc.loadNamespace(namespaceID); err != nil {
		return
	} else {
		return svc.checkPermissions(svc.chartRepo.FindByHandle(namespaceID, handle))
	}
}

func (svc chart) checkPermissions(c *types.Chart, err error) (*types.Chart, error) {
	if err != nil {
		return nil, err
	} else if !svc.ac.CanReadChart(svc.ctx, c) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return c, err
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
	if !handle.IsValid(mod.Handle) {
		return nil, ErrInvalidHandle
	}

	if ns, err := svc.loadNamespace(mod.NamespaceID); err != nil {
		return nil, err
	} else if !svc.ac.CanCreateChart(svc.ctx, ns) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	if err = svc.UniqueCheck(mod); err != nil {
		return
	}

	return svc.chartRepo.Create(mod)
}

func (svc chart) Update(mod *types.Chart) (c *types.Chart, err error) {
	if !handle.IsValid(mod.Handle) {
		return nil, ErrInvalidHandle
	}

	if mod.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if c, err = svc.chartRepo.FindByID(mod.NamespaceID, mod.ID); err != nil {
		return
	}

	if isStale(mod.UpdatedAt, c.UpdatedAt, c.CreatedAt) {
		return nil, ErrStaleData.withStack()
	}

	if err = svc.UniqueCheck(mod); err != nil {
		return
	}

	if !svc.ac.CanUpdateChart(svc.ctx, c) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	c.Config = mod.Config
	c.Name = mod.Name
	c.Handle = mod.Handle

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

func (svc chart) UniqueCheck(c *types.Chart) (err error) {
	if c.Handle != "" {
		if e, _ := svc.chartRepo.FindByHandle(c.NamespaceID, c.Handle); e != nil && e.ID != c.ID {
			return repository.ErrChartHandleNotUnique
		}
	}

	return nil
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
