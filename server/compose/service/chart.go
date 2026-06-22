package service

import (
	"context"
	"reflect"
	"strconv"

	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/store"
)

type (
	chart struct {
		actionlog actionlog.Recorder
		ac        chartAccessController
		store     store.Storer
		locale    ResourceTranslationsManagerService
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
)

func Chart() *chart {
	return &chart{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		locale:    DefaultResourceTranslation,
	}
}

// onSearch is the generated Search body handler.
//
// The recordAction wrapper and aProps (filter) are owned by the generated
// chart.gen.go; this handler runs the namespace preload, access check and store
// search.
func (svc *chart) onSearch(ctx context.Context, filter types.ChartFilter, aProps *chartActionProps) (set types.ChartSet, f types.ChartFilter, err error) {
	var ns *types.Namespace

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Chart) (bool, error) {
		if !svc.ac.CanReadChart(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	ns, err = loadNamespace(ctx, svc.store, filter.NamespaceID)
	if err != nil {
		return
	}

	aProps.setNamespace(ns)
	if !svc.ac.CanSearchChartsOnNamespace(ctx, ns) {
		return nil, f, ChartErrNotAllowedToSearch()
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
			return
		}

		// labels specified but no labeled resources found
		if len(filter.LabeledIDs) == 0 {
			return
		}
	}

	if set, f, err = store.SearchComposeCharts(ctx, svc.store, filter); err != nil {
		return
	}

	if err = label.Load(ctx, svc.store, toLabeledCharts(set)...); err != nil {
		return
	}

	set.Walk(func(c *types.Chart) error {
		svc.proc(ctx, c)
		return nil
	})

	return
}

// onLookup is the generated FindByID body handler (namespace-scoped compound id).
func (svc *chart) onLookup(ctx context.Context, namespaceID, chartID uint64, aProps *chartActionProps) (c *types.Chart, err error) {
	return svc.lookup(ctx, namespaceID, aProps, func(aProps *chartActionProps) (*types.Chart, error) {
		if chartID == 0 || namespaceID == 0 {
			return nil, ChartErrInvalidID()
		}
		res, err := loadChart(ctx, svc.store, chartID)
		if err != nil {
			return nil, err
		}
		if res.NamespaceID != namespaceID {
			return nil, ChartErrNotFound()
		}
		return res, nil
	})
}

func (svc *chart) FindByHandle(ctx context.Context, namespaceID uint64, h string) (c *types.Chart, err error) {
	var aProps = &chartActionProps{chart: &types.Chart{NamespaceID: namespaceID}}

	c, err = svc.lookup(ctx, namespaceID, aProps, func(aProps *chartActionProps) (*types.Chart, error) {
		if !handle.IsValid(h) {
			return nil, ChartErrInvalidHandle()
		}

		aProps.chart.Handle = h
		return store.LookupComposeChartByNamespaceIDHandle(ctx, svc.store, namespaceID, h)
	})

	return c, svc.recordAction(ctx, aProps, ChartActionLookup, err)
}

func (svc *chart) proc(ctx context.Context, c *types.Chart) {
	if svc.locale == nil || svc.locale.Locale() == nil {
		return
	}

	tag := locale.GetAcceptLanguageFromContext(ctx)
	c.DecodeTranslations(svc.locale.Locale().ResourceTranslations(tag, c.ResourceTranslation()))
}

// onCreate is the generated Create body handler.
//
// The recordAction wrapper, aProps and res=new assignment are owned by the
// generated chart.gen.go.
func (svc *chart) onCreate(ctx context.Context, new *types.Chart) error {
	var (
		ns     *types.Namespace
		aProps = &chartActionProps{chart: new}
	)

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !handle.IsValid(new.Handle) {
			return ChartErrInvalidHandle()
		}

		if ns, err = loadNamespace(ctx, s, new.NamespaceID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

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

		// generate config element IDs
		new.Config.GenerateIDs(nextID)

		aProps.setChanged(new)

		if err = store.CreateComposeChart(ctx, s, new); err != nil {
			return err
		}

		if err = updateTranslations(ctx, svc.ac, svc.locale, new.EncodeTranslations()...); err != nil {
			return
		}

		if err = label.Create(ctx, s, new); err != nil {
			return
		}

		return nil
	})
}

// onUpdate is the generated Update body handler.
func (svc *chart) onUpdate(ctx context.Context, s store.Storer, upd, res *types.Chart, aProps *chartActionProps, _ func() error, _ func() error) error {
	if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
		return ChartErrInvalidHandle()
	}

	if err := svc.uniqueCheck(ctx, upd); err != nil {
		return err
	}

	if !svc.ac.CanUpdateChart(ctx, res) {
		return ChartErrNotAllowedToUpdate()
	}

	changed := false

	if res.Name != upd.Name {
		res.Name = upd.Name
		changed = true
	}

	if res.Handle != upd.Handle {
		res.Handle = upd.Handle
		changed = true
	}

	if !reflect.DeepEqual(upd.Config, res.Config) {
		res.Config = upd.Config
		changed = true
	}

	// Assure ReportIDs
	for i, r := range res.Config.Reports {
		if r.ReportID == 0 {
			r.ReportID = nextID()
			res.Config.Reports[i] = r
			changed = true
		}

		// Ensure chart report metric IDs
		for j, m := range r.Metrics {
			if val, ok := m["metricID"]; !ok || val == 0 {
				m["metricID"] = strconv.FormatUint(nextID(), 10)
				res.Config.Reports[i].Metrics[j] = m
				changed = true
			}
		}
	}

	if changed {
		res.UpdatedAt = now()
	}

	if upd.Labels != nil && label.Changed(res.Labels, upd.Labels) {
		res.Labels = upd.Labels
		if err := label.Update(ctx, s, res); err != nil {
			return err
		}
	}

	// generate config element IDs if missing
	res.Config.GenerateIDs(nextID)

	if err := store.UpdateComposeChart(ctx, s, res); err != nil {
		return err
	}

	if err := updateTranslations(ctx, svc.ac, svc.locale, res.EncodeTranslations()...); err != nil {
		return err
	}

	return nil
}

// onDelete is the generated DeleteByID body handler (namespace-scoped compound id).
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

	if err := updateTranslations(ctx, svc.ac, svc.locale, res.EncodeTranslations()...); err != nil {
		return err
	}

	return nil
}

// onUndelete is the generated UndeleteByID body handler (namespace-scoped compound id).
func (svc *chart) onUndelete(ctx context.Context, s store.Storer, namespaceID uint64, res *types.Chart, aProps *chartActionProps) error {
	if !svc.ac.CanDeleteChart(ctx, res) {
		return ChartErrNotAllowedToUndelete()
	}

	if res.DeletedAt == nil {
		return nil
	}

	res.DeletedAt = nil

	if err := store.UpdateComposeChart(ctx, s, res); err != nil {
		return err
	}

	if err := updateTranslations(ctx, svc.ac, svc.locale, res.EncodeTranslations()...); err != nil {
		return err
	}

	return nil
}

// lookup fn() orchestrates chart lookup, namespace preload and check
func (svc *chart) lookup(ctx context.Context, namespaceID uint64, aProps *chartActionProps, lookup func(*chartActionProps) (*types.Chart, error)) (c *types.Chart, err error) {
	if aProps.chart == nil {
		aProps.chart = &types.Chart{NamespaceID: namespaceID}
	}

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

	return c, err
}

func (svc *chart) uniqueCheck(ctx context.Context, c *types.Chart) (err error) {
	if c.Handle != "" {
		if e, _ := store.LookupComposeChartByNamespaceIDHandle(ctx, svc.store, c.NamespaceID, c.Handle); e != nil && e.ID != c.ID {
			return ChartErrHandleNotUnique()
		}
	}

	return nil
}

// loadChartScoped loads a chart by ID and validates it belongs to the given namespace.
func loadChartScoped(ctx context.Context, s store.Storer, namespaceID, chartID uint64) (res *types.Chart, err error) {
	if res, err = loadChart(ctx, s, chartID); err == nil && res.NamespaceID != namespaceID {
		return nil, ChartErrNotFound()
	}
	return
}
