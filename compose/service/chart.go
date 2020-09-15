package service

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	chart struct {
		ctx       context.Context
		actionlog actionlog.Recorder
		ac        chartAccessController
		store     store.Storer
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

	chartUpdateHandler func(ctx context.Context, ns *types.Namespace, c *types.Chart) (bool, error)
)

func Chart() ChartService {
	return (&chart{
		ctx: context.Background(),
		ac:  DefaultAccessControl,
	}).With(context.Background())
}

func (svc chart) With(ctx context.Context) ChartService {
	return &chart{
		ctx:       ctx,
		actionlog: DefaultActionlog,
		ac:        svc.ac,
		store:     DefaultStore,
	}
}

func (svc chart) Find(filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error) {
	var (
		aProps = &chartActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Chart) (bool, error) {
		if !svc.ac.CanReadChart(svc.ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if ns, err := loadNamespace(svc.ctx, svc.store, filter.NamespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if set, f, err = store.SearchComposeCharts(svc.ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(svc.ctx, aProps, ChartActionSearch, err)
}

func (svc chart) FindByID(namespaceID, chartID uint64) (c *types.Chart, err error) {
	return svc.lookup(namespaceID, func(aProps *chartActionProps) (*types.Chart, error) {
		if chartID == 0 {
			return nil, ChartErrInvalidID()
		}

		aProps.chart.ID = chartID
		return store.LookupComposeChartByID(svc.ctx, svc.store, chartID)
	})
}

func (svc chart) FindByHandle(namespaceID uint64, h string) (c *types.Chart, err error) {
	return svc.lookup(namespaceID, func(aProps *chartActionProps) (*types.Chart, error) {
		if !handle.IsValid(h) {
			return nil, ChartErrInvalidHandle()
		}

		aProps.chart.Handle = h
		return store.LookupComposeChartByNamespaceIDHandle(svc.ctx, svc.store, namespaceID, h)
	})
}

func (svc chart) Create(new *types.Chart) (*types.Chart, error) {
	var (
		err    error
		ns     *types.Namespace
		aProps = &chartActionProps{changed: new}
	)

	err = func() error {
		if !handle.IsValid(new.Handle) {
			return ChartErrInvalidHandle()
		}

		if ns, err = loadNamespace(svc.ctx, svc.store, new.NamespaceID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if !svc.ac.CanCreateChart(svc.ctx, ns) {
			return ChartErrNotAllowedToCreate()
		}

		new.ID = id.Next()
		new.CreatedAt = *nowPtr()
		new.UpdatedAt = nil
		new.DeletedAt = nil

		if err = svc.uniqueCheck(new); err != nil {
			return err
		}

		return store.CreateComposeChart(svc.ctx, svc.store, new)
	}()

	return new, svc.recordAction(svc.ctx, aProps, ChartActionCreate, err)
}

func (svc chart) Update(upd *types.Chart) (c *types.Chart, err error) {
	return svc.updater(upd.NamespaceID, upd.ID, ChartActionUpdate, svc.handleUpdate(upd))
}

func (svc chart) DeleteByID(namespaceID, chartID uint64) error {
	return trim1st(svc.updater(namespaceID, chartID, ChartActionDelete, svc.handleDelete))
}

func (svc chart) UndeleteByID(namespaceID, chartID uint64) error {
	return trim1st(svc.updater(namespaceID, chartID, ChartActionUndelete, svc.handleUndelete))
}

// lookup fn() orchestrates chart lookup, namespace preload and check
func (svc chart) lookup(namespaceID uint64, lookup func(*chartActionProps) (*types.Chart, error)) (c *types.Chart, err error) {
	var aProps = &chartActionProps{chart: &types.Chart{NamespaceID: namespaceID}}

	err = func() error {
		if ns, err := loadNamespace(svc.ctx, svc.store, namespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if c, err = lookup(aProps); errors.Is(err, store.ErrNotFound) {
			return ChartErrNotFound()
		} else if err != nil {
			return err
		}

		aProps.setChart(c)

		if !svc.ac.CanReadChart(svc.ctx, c) {
			return ChartErrNotAllowedToRead()
		}

		return nil
	}()

	return c, svc.recordAction(svc.ctx, aProps, ChartActionLookup, err)
}

func (svc chart) updater(namespaceID, chartID uint64, action func(...*chartActionProps) *chartAction, fn chartUpdateHandler) (*types.Chart, error) {
	var (
		changed bool
		ns      *types.Namespace
		c       *types.Chart
		aProps  = &chartActionProps{chart: &types.Chart{ID: chartID, NamespaceID: namespaceID}}
		err     error
	)

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, c, err = loadChart(ctx, s, namespaceID, chartID)
		if err != nil {
			return
		}

		aProps.setNamespace(ns)
		aProps.setChanged(c)

		if changed, err = fn(ctx, ns, c); err != nil {
			return err
		} else if !changed {
			return
		}

		return store.UpdateComposeChart(ctx, s, c)
	})

	return c, svc.recordAction(svc.ctx, aProps, action, err)
}

func (svc chart) uniqueCheck(c *types.Chart) (err error) {
	if c.Handle != "" {
		if e, _ := store.LookupComposeChartByNamespaceIDHandle(svc.ctx, svc.store, c.NamespaceID, c.Handle); e != nil && e.ID != c.ID {
			return ChartErrHandleNotUnique()
		}
	}

	return nil
}

func (svc chart) handleUpdate(upd *types.Chart) chartUpdateHandler {
	return func(ctx context.Context, ns *types.Namespace, c *types.Chart) (bool, error) {
		if isStale(upd.UpdatedAt, c.UpdatedAt, c.CreatedAt) {
			return false, ChartErrStaleData()
		}

		if upd.Handle != c.Handle && !handle.IsValid(upd.Handle) {
			return false, ChartErrInvalidHandle()
		}

		if err := svc.uniqueCheck(upd); err != nil {
			return false, err
		}

		if !svc.ac.CanUpdateChart(svc.ctx, c) {
			return false, ChartErrNotAllowedToUpdate()
		}

		c.Name = upd.Name
		c.Handle = upd.Handle
		c.Config = upd.Config
		c.UpdatedAt = nowPtr()
		return true, nil
	}
}

func (svc chart) handleDelete(ctx context.Context, ns *types.Namespace, c *types.Chart) (bool, error) {
	if !svc.ac.CanDeleteChart(ctx, c) {
		return false, ChartErrNotAllowedToDelete()
	}

	if c.DeletedAt != nil {
		// chart already deleted
		return false, nil
	}

	c.DeletedAt = nowPtr()
	return true, nil
}

func (svc chart) handleUndelete(ctx context.Context, ns *types.Namespace, c *types.Chart) (bool, error) {
	if !svc.ac.CanDeleteChart(ctx, c) {
		return false, ChartErrNotAllowedToUndelete()
	}

	if c.DeletedAt == nil {
		// chart not deleted
		return false, nil
	}

	c.DeletedAt = nil
	return true, nil
}

func loadChart(ctx context.Context, s store.Storer, namespaceID, chartID uint64) (ns *types.Namespace, c *types.Chart, err error) {
	if chartID == 0 {
		return nil, nil, ChartErrInvalidID()
	}

	if ns, err = loadNamespace(ctx, s, namespaceID); err == nil {
		if c, err = store.LookupComposeChartByID(ctx, s, chartID); errors.Is(err, store.ErrNotFound) {
			err = ChartErrNotFound()
		}
	}

	if err != nil {
		return nil, nil, err
	}

	if namespaceID != c.NamespaceID {
		// Make sure chart belongs to the right namespace
		return nil, nil, ChartErrNotFound()
	}

	return
}
