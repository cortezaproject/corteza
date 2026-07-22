package service

import (
	"context"
	"reflect"
	"strconv"

	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/store"
)

type (
	chartServices struct {
		locale ResourceTranslationsManagerService
	}

	chartAccessController interface {
		CanManageResourceTranslations(ctx context.Context) bool
		CanSearchChartsOnNamespace(context.Context, *types.Namespace) bool
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanCreateChartOnNamespace(context.Context, *types.Namespace) bool
		CanReadChart(context.Context, *types.Chart) bool
		CanUpdateChart(context.Context, *types.Chart) bool
		CanDeleteChart(context.Context, *types.Chart) bool
	}

	chartUpdateHandler func(ctx context.Context, ns *types.Namespace, c *types.Chart) (chartChanges, error)

	chartChanges uint8
)

const (
	chartUnchanged     chartChanges = 0
	chartChanged       chartChanges = 1
	chartLabelsChanged chartChanges = 2
)

func Chart() *chart {
	return &chart{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		services: &chartServices{
			locale: DefaultResourceTranslation,
		},
	}
}

func (svc *chart) guard(ctx context.Context, res *types.Chart) error {
	return guardNamespaceWritable(ctx, svc.store, res.NamespaceID)
}

func (svc *chart) onLookup(ctx context.Context, namespaceID, ID uint64, aProps *chartActionProps) (*types.Chart, error) {
	return svc.lookup(ctx, namespaceID, func(p *chartActionProps) (*types.Chart, error) {
		if ID == 0 {
			return nil, ChartErrInvalidID()
		}
		p.chart.ID = ID
		return store.LookupComposeChartByID(ctx, svc.store, ID)
	})
}

func (svc *chart) onSearch(ctx context.Context, filter types.ChartFilter, aProps *chartActionProps) (types.ChartSet, types.ChartFilter, error) {
	return svc.Find(ctx, filter)
}

func (svc *chart) onCreate(ctx context.Context, new *types.Chart) error {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !handle.IsValid(new.Handle) {
			return ChartErrInvalidHandle()
		}

		var ns *types.Namespace
		if ns, err = loadNamespace(ctx, s, new.NamespaceID); err != nil {
			return err
		}

		if !svc.ac.CanCreateChartOnNamespace(ctx, ns) {
			return ChartErrNotAllowedToCreate()
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.UpdatedAt = nil
		new.DeletedAt = nil
		new.Config.GenerateIDs(nextID)

		// Inherit the owning project from the namespace; a chart is addressed via
		// its namespace, so the project never comes from the request payload.
		new.ProjectID = ns.ProjectID

		if err = store.CreateComposeChart(ctx, s, new); err != nil {
			return err
		}

		if err = updateTranslations(ctx, svc.ac, svc.services.locale, new.EncodeTranslations()...); err != nil {
			return
		}

		return label.Create(ctx, s, new)
	})
}

func (svc *chart) onUpdate(ctx context.Context, s store.Storer, upd, res *types.Chart, aProps *chartActionProps, before, after func() error) error {
	if !svc.ac.CanUpdateChart(ctx, res) {
		return ChartErrNotAllowedToUpdate()
	}

	if err := svc.uniqueCheck(ctx, upd); err != nil {
		return err
	}

	res.Name = upd.Name
	res.Handle = upd.Handle
	res.Config = upd.Config
	res.Config.GenerateIDs(nextID)

	if err := updateTranslations(ctx, svc.ac, svc.services.locale, res.EncodeTranslations()...); err != nil {
		return err
	}

	return nil
}

func (svc *chart) onDelete(ctx context.Context, s store.Storer, namespaceID uint64, res *types.Chart, aProps *chartActionProps) error {
	if !svc.ac.CanDeleteChart(ctx, res) {
		return ChartErrNotAllowedToDelete()
	}

	if res.DeletedAt != nil {
		return nil
	}

	res.DeletedAt = now()

	if err := store.UpdateComposeChart(ctx, s, res); err != nil {
		return err
	}

	return updateTranslations(ctx, svc.ac, svc.services.locale, res.EncodeTranslations()...)
}

func (svc *chart) onUndelete(ctx context.Context, s store.Storer, namespaceID uint64, res *types.Chart, aProps *chartActionProps) error {
	if !svc.ac.CanDeleteChart(ctx, res) {
		return ChartErrNotAllowedToUndelete()
	}
	res.DeletedAt = nil
	return store.UpdateComposeChart(ctx, s, res)
}

func (svc chart) Find(ctx context.Context, filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error) {
	var (
		aProps = &chartActionProps{filter: &filter}
		ns     *types.Namespace
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Chart) (bool, error) {
		if !svc.ac.CanReadChart(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		ns, err = loadNamespace(ctx, svc.store, filter.NamespaceID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)
		if !svc.ac.CanSearchChartsOnNamespace(ctx, ns) {
			return ChartErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Chart{}.LabelResourceKind(),
				filter.Labels,
				id.Uints(filter.ChartID...)...,
			)

			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if set, f, err = store.SearchComposeCharts(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledCharts(set)...); err != nil {
			return err
		}

		set.Walk(func(c *types.Chart) error {
			svc.proc(ctx, c)
			return nil
		})

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ChartActionSearch, err)
}

func (svc chart) FindByHandle(ctx context.Context, namespaceID uint64, h string) (c *types.Chart, err error) {
	return svc.lookup(ctx, namespaceID, func(aProps *chartActionProps) (*types.Chart, error) {
		if !handle.IsValid(h) {
			return nil, ChartErrInvalidHandle()
		}

		aProps.chart.Handle = h
		return store.LookupComposeChartByNamespaceIDHandle(ctx, svc.store, namespaceID, h)
	})
}

func (svc chart) proc(ctx context.Context, c *types.Chart) {
	if svc.services.locale == nil || svc.services.locale.Locale() == nil {
		return
	}

	tag := locale.GetAcceptLanguageFromContext(ctx)
	c.DecodeTranslations(svc.services.locale.Locale().ResourceTranslations(tag, c.ResourceTranslation()))
}

// lookup fn() orchestrates chart lookup, namespace preload and check
func (svc chart) lookup(ctx context.Context, namespaceID uint64, lookup func(*chartActionProps) (*types.Chart, error)) (c *types.Chart, err error) {
	var aProps = &chartActionProps{chart: &types.Chart{NamespaceID: namespaceID}}

	err = func() error {
		if ns, err := loadNamespace(ctx, svc.store, namespaceID); err != nil {
			return err
		} else {
			aProps.setNamespace(ns)
		}

		if c, err = lookup(aProps); errors.IsNotFound(err) {
			return ChartErrNotFound()
		} else if err != nil {
			return err
		}

		aProps.setChart(c)

		if !svc.ac.CanReadChart(ctx, c) {
			return ChartErrNotAllowedToRead()
		}

		svc.proc(ctx, c)

		if err = label.Load(ctx, svc.store, c); err != nil {
			return err
		}

		return nil
	}()

	return c, svc.recordAction(ctx, aProps, ChartActionLookup, err)
}

func (svc chart) updater(ctx context.Context, namespaceID, chartID uint64, action func(...*chartActionProps) *chartAction, fn chartUpdateHandler) (*types.Chart, error) {
	var (
		changes chartChanges
		ns      *types.Namespace
		c       *types.Chart
		aProps  = &chartActionProps{chart: &types.Chart{ID: chartID, NamespaceID: namespaceID}}
		err     error
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		ns, c, err = loadChartCombo(ctx, s, namespaceID, chartID)
		if err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, c); err != nil {
			return err
		}

		aProps.setNamespace(ns)
		aProps.setChanged(c)

		if changes, err = fn(ctx, ns, c); err != nil {
			return err
		}

		// generate config element IDs if missing
		c.Config.GenerateIDs(nextID)

		if changes&chartChanged > 0 {
			if err = store.UpdateComposeChart(ctx, s, c); err != nil {
				return err
			}
		}

		if err = updateTranslations(ctx, svc.ac, svc.services.locale, c.EncodeTranslations()...); err != nil {
			return
		}

		if changes&chartLabelsChanged > 0 {
			if err = label.Update(ctx, s, c); err != nil {
				return
			}
		}

		return nil
	})

	return c, svc.recordAction(ctx, aProps, action, err)
}

func (svc chart) uniqueCheck(ctx context.Context, c *types.Chart) (err error) {
	if c.Handle != "" {
		if e, _ := store.LookupComposeChartByNamespaceIDHandle(ctx, svc.store, c.NamespaceID, c.Handle); e != nil && e.ID != c.ID {
			return ChartErrHandleNotUnique()
		}
	}

	return nil
}

func (svc chart) handleUpdate(ctx context.Context, upd *types.Chart) chartUpdateHandler {
	return func(ctx context.Context, ns *types.Namespace, res *types.Chart) (changes chartChanges, err error) {
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return chartUnchanged, ChartErrStaleData()
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return chartUnchanged, ChartErrInvalidHandle()
		}

		if err := svc.uniqueCheck(ctx, upd); err != nil {
			return chartUnchanged, err
		}

		if !svc.ac.CanUpdateChart(ctx, res) {
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

		// Assure ReportIDs
		for i, r := range res.Config.Reports {
			if r.ReportID == 0 {
				r.ReportID = nextID()
				res.Config.Reports[i] = r

				changes |= chartChanged
			}

			// Ensure chart report metric IDs
			for j, m := range r.Metrics {
				if val, ok := m["metricID"]; !ok || val == 0 {
					m["metricID"] = strconv.FormatUint(nextID(), 10)
					res.Config.Reports[i].Metrics[j] = m

					changes |= chartChanged
				}
			}
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

func loadChartCombo(ctx context.Context, s interface {
	store.ComposeCharts
	store.ComposeNamespaces
}, namespaceID, chartID uint64) (ns *types.Namespace, c *types.Chart, err error) {
	ns, err = loadNamespace(ctx, s, namespaceID)
	if err != nil {
		return
	}

	c, err = loadChartScoped(ctx, s, namespaceID, chartID)
	return
}
