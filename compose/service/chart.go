package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	chart struct {
		db  *factory.DB
		ctx context.Context

		actionlog actionlog.Recorder

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

		FilterReadableCharts(ctx context.Context) *permissions.ResourceFilter
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
		ac: DefaultAccessControl,
	}).With(context.Background())
}

func (svc chart) With(ctx context.Context) ChartService {
	db := repository.DB(ctx)
	return &chart{
		db:  db,
		ctx: ctx,

		actionlog: DefaultActionlog,

		ac: svc.ac,

		chartRepo: repository.Chart(ctx, db),
		nsRepo:    repository.Namespace(ctx, db),
	}
}

// lookup fn() orchestrates chart lookup, namespace preload and check
func (svc chart) lookup(namespaceID uint64, lookup func(*chartActionProps) (*types.Chart, error)) (m *types.Chart, err error) {
	var aProps = &chartActionProps{chart: &types.Chart{NamespaceID: namespaceID}}

	err = svc.db.Transaction(func() error {
		if ns, err := svc.loadNamespace(namespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if m, err = lookup(aProps); err != nil {
			if repository.ErrChartNotFound.Eq(err) {
				return ChartErrNotFound()
			}

			return err
		}

		aProps.setChart(m)

		if !svc.ac.CanReadChart(svc.ctx, m) {
			return ChartErrNotAllowedToRead()
		}

		return nil
	})

	return m, svc.recordAction(svc.ctx, aProps, ChartActionLookup, err)
}

func (svc chart) FindByID(namespaceID, chartID uint64) (p *types.Chart, err error) {
	return svc.lookup(namespaceID, func(aProps *chartActionProps) (*types.Chart, error) {
		if chartID == 0 {
			return nil, ChartErrInvalidID()
		}

		aProps.chart.ID = chartID
		return svc.chartRepo.FindByID(namespaceID, chartID)
	})
}

func (svc chart) FindByHandle(namespaceID uint64, h string) (c *types.Chart, err error) {
	return svc.lookup(namespaceID, func(aProps *chartActionProps) (*types.Chart, error) {
		if !handle.IsValid(h) {
			return nil, ChartErrInvalidHandle()
		}

		aProps.chart.Handle = h
		return svc.chartRepo.FindByHandle(namespaceID, h)
	})
}

// search fn() orchestrates charts search, namespace preload and check
func (svc chart) search(filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error) {
	var (
		aProps = &chartActionProps{filter: &filter}
	)

	f = filter
	f.IsReadable = svc.ac.FilterReadableCharts(svc.ctx)

	err = svc.db.Transaction(func() error {
		if ns, err := svc.loadNamespace(f.NamespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if set, f, err = svc.chartRepo.Find(f); err != nil {
			return err
		}

		return nil
	})

	return set, f, svc.recordAction(svc.ctx, aProps, ChartActionSearch, err)
}

func (svc chart) Find(filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error) {
	return svc.search(filter)
}

func (svc chart) Create(new *types.Chart) (c *types.Chart, err error) {
	var (
		ns     *types.Namespace
		aProps = &chartActionProps{changed: new}
	)

	err = svc.db.Transaction(func() error {

		if !handle.IsValid(new.Handle) {
			return ChartErrInvalidHandle()
		}

		if ns, err = svc.loadNamespace(new.NamespaceID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if !svc.ac.CanCreateChart(svc.ctx, ns) {
			return ChartErrNotAllowedToCreate()
		}

		if err = svc.UniqueCheck(new); err != nil {
			return err
		}
		c, err = svc.chartRepo.Create(new)
		return err
	})

	return c, svc.recordAction(svc.ctx, aProps, ChartActionCreate, err)
}

func (svc chart) Update(upd *types.Chart) (c *types.Chart, err error) {
	var (
		ns     *types.Namespace
		aProps = &chartActionProps{changed: upd}
	)

	err = svc.db.Transaction(func() error {
		if !handle.IsValid(upd.Handle) {
			return ChartErrInvalidHandle()
		}

		if upd.ID == 0 {
			return ChartErrInvalidID()
		}

		if ns, err = svc.loadNamespace(upd.NamespaceID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if c, err = svc.chartRepo.FindByID(upd.NamespaceID, upd.ID); err != nil {
			if repository.ErrModuleNotFound.Eq(err) {
				return ModuleErrNotFound()
			}

			return err
		}

		if isStale(upd.UpdatedAt, c.UpdatedAt, c.CreatedAt) {
			return ChartErrStaleData()
		}

		if err = svc.UniqueCheck(upd); err != nil {
			return err
		}

		if !svc.ac.CanUpdateChart(svc.ctx, c) {
			return ChartErrNotAllowedToUpdate()
		}

		c.Config = upd.Config
		c.Name = upd.Name
		c.Handle = upd.Handle

		c, err = svc.chartRepo.Update(c)
		return err
	})

	return c, svc.recordAction(svc.ctx, aProps, ChartActionUpdate, err)
}

func (svc chart) DeleteByID(namespaceID, chartID uint64) (err error) {
	var (
		ns     *types.Namespace
		c      *types.Chart
		aProps = &chartActionProps{chart: &types.Chart{ID: chartID, NamespaceID: namespaceID}}
	)

	err = svc.db.Transaction(func() (err error) {
		if chartID == 0 {
			return ChartErrInvalidID()
		}

		if ns, err = svc.loadNamespace(namespaceID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if c, err = svc.chartRepo.FindByID(namespaceID, chartID); err != nil {
			if repository.ErrChartNotFound.Eq(err) {
				return ChartErrNotFound()
			}

			return err
		}

		aProps.setChanged(c)

		if !svc.ac.CanDeleteChart(svc.ctx, c) {
			return ChartErrNotAllowedToDelete()
		}

		return svc.chartRepo.DeleteByID(namespaceID, chartID)
	})

	return svc.recordAction(svc.ctx, aProps, ChartActionDelete, err)

}

func (svc chart) UniqueCheck(c *types.Chart) (err error) {
	if c.Handle != "" {
		if e, _ := svc.chartRepo.FindByHandle(c.NamespaceID, c.Handle); e != nil && e.ID != c.ID {
			return ChartErrHandleNotUnique()
		}
	}

	return nil
}

func (svc chart) loadNamespace(namespaceID uint64) (ns *types.Namespace, err error) {
	if namespaceID == 0 {
		return nil, ChartErrInvalidNamespaceID()
	}

	if ns, err = svc.nsRepo.FindByID(namespaceID); err != nil {
		return
	}

	if !svc.ac.CanReadNamespace(svc.ctx, ns) {
		return nil, ChartErrNotAllowedToReadNamespace()
	}

	return
}
