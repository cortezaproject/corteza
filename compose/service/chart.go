package service

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/store"
	"reflect"
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

	chartUpdateHandler func(ctx context.Context, ns *types.Namespace, c *types.Chart) (chartChanges, error)

	chartChanges uint8
)

const (
	chartUnchanged     chartChanges = 0
	chartChanged       chartChanges = 1
	chartLabelsChanged chartChanges = 2
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

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				svc.ctx,
				svc.store,
				types.Chart{}.LabelResourceKind(),
				filter.Labels,
				filter.ChartID...,
			)

			if err != nil {
				return err
			}
		}

		if set, f, err = store.SearchComposeCharts(svc.ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(svc.ctx, svc.store, toLabeledCharts(set)...); err != nil {
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

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !handle.IsValid(new.Handle) {
			return ChartErrInvalidHandle()
		}

		if ns, err = loadNamespace(ctx, s, new.NamespaceID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if !svc.ac.CanCreateChart(ctx, ns) {
			return ChartErrNotAllowedToCreate()
		}

		if err = svc.uniqueCheck(new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.UpdatedAt = nil
		new.DeletedAt = nil

		if err = store.CreateComposeChart(ctx, s, new); err != nil {
			return err
		}

		if err = label.Create(ctx, s, new); err != nil {
			return
		}

		return nil
	})

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

		if err = label.Load(svc.ctx, svc.store, c); err != nil {
			return err
		}

		return nil
	}()

	return c, svc.recordAction(svc.ctx, aProps, ChartActionLookup, err)
}

func (svc chart) updater(namespaceID, chartID uint64, action func(...*chartActionProps) *chartAction, fn chartUpdateHandler) (*types.Chart, error) {
	var (
		changes chartChanges
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

		if err = label.Load(svc.ctx, svc.store, c); err != nil {
			return err
		}

		aProps.setNamespace(ns)
		aProps.setChanged(c)

		if changes, err = fn(ctx, ns, c); err != nil {
			return err
		}

		if changes&chartChanged > 0 {
			if err = store.UpdateComposeChart(ctx, s, c); err != nil {
				return err
			}
		}

		if changes&chartLabelsChanged > 0 {
			if err = label.Update(ctx, s, c); err != nil {
				return
			}
		}

		return nil
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
	return func(ctx context.Context, ns *types.Namespace, res *types.Chart) (changes chartChanges, err error) {
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return chartUnchanged, ChartErrStaleData()
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return chartUnchanged, ChartErrInvalidHandle()
		}

		if err := svc.uniqueCheck(upd); err != nil {
			return chartUnchanged, err
		}

		if !svc.ac.CanUpdateChart(svc.ctx, res) {
			return chartUnchanged, ChartErrNotAllowedToUpdate()
		}

		if res.Name != upd.Name {
			changes |= chartChanged
			res.Name = upd.Name
		}

		if res.Handle != upd.Handle {
			changes |= chartChanged
			res.Handle = upd.Handle
		}

		if !reflect.DeepEqual(upd.Config, res.Config) {
			changes |= chartChanged
			res.Config = upd.Config
		}

		if changes&chartChanged > 0 {
			res.UpdatedAt = now()
		}

		if upd.Labels != nil {
			if label.Changed(res.Labels, upd.Labels) {
				changes |= chartLabelsChanged
				res.Labels = upd.Labels
			}
		}

		return
	}
}

func (svc chart) handleDelete(ctx context.Context, ns *types.Namespace, c *types.Chart) (chartChanges, error) {
	if !svc.ac.CanDeleteChart(ctx, c) {
		return chartUnchanged, ChartErrNotAllowedToDelete()
	}

	if c.DeletedAt != nil {
		// chart already deleted
		return chartUnchanged, nil
	}

	c.DeletedAt = now()
	return chartChanged, nil
}

func (svc chart) handleUndelete(ctx context.Context, ns *types.Namespace, c *types.Chart) (chartChanges, error) {
	if !svc.ac.CanDeleteChart(ctx, c) {
		return chartUnchanged, ChartErrNotAllowedToUndelete()
	}

	if c.DeletedAt == nil {
		// chart not deleted
		return chartUnchanged, nil
	}

	c.DeletedAt = nil
	return chartChanged, nil
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

// toLabeledCharts converts to []label.LabeledResource
//
// This function is auto-generated.
func toLabeledCharts(set []*types.Chart) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
