package importer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

type (
	Chart struct {
		imp       *Importer
		namespace *types.Namespace
		set       types.ChartSet
		dirty     map[uint64]bool
		modRefs   []chartModuleRef
	}

	chartModuleRef struct {
		// chart handle, report index, module handle
		ch string
		ri int
		mh string
	}

	// @todo remove finder strategy, directly provide set of items
	chartFinder interface {
		Find(filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error)
	}
)

func NewChartImporter(imp *Importer, ns *types.Namespace) *Chart {
	out := &Chart{
		imp:       imp,
		namespace: ns,
		set:       types.ChartSet{},
		dirty:     make(map[uint64]bool),
	}

	if imp.chartFinder != nil && ns.ID > 0 {
		out.set, _, _ = imp.chartFinder.Find(types.ChartFilter{NamespaceID: ns.ID})
	}

	return out
}

func (pImp *Chart) getModule(handle string) (*types.Module, error) {
	if g, ok := pImp.imp.namespaces.modules[pImp.namespace.Slug]; !ok {
		return nil, errors.Errorf("could not get modules %q from non existing namespace %q", handle, pImp.namespace.Slug)
	} else {
		return g.Get(handle)
	}
}

// CastSet Resolves permission rules:
// { <chart-handle>: { chart } } or [ { chart }, ... ]
func (cImp *Chart) CastSet(set interface{}) error {
	return deinterfacer.Each(set, func(index int, handle string, def interface{}) error {
		if index > -1 {
			// Charts defined as collection
			deinterfacer.KVsetString(&handle, "handle", def)
		}

		return cImp.Cast(handle, def)
	})
}

// Cast Resolves permission rules:
// { <chart-handle>: { chart } } or [ { chart }, ... ]
func (cImp *Chart) Cast(handle string, def interface{}) (err error) {
	if !deinterfacer.IsMap(def) {
		return errors.New("expecting map of values for chart")
	}

	var chart *types.Chart

	if !importer.IsValidHandle(handle) {
		return errors.New("invalid chart handle")
	}

	handle = importer.NormalizeHandle(handle)
	if chart, err = cImp.Get(handle); err != nil {
		return err
	} else if chart == nil {
		chart = &types.Chart{
			Handle: handle,
			Name:   handle,
		}

		cImp.set = append(cImp.set, chart)
	} else if chart.ID == 0 {
		return errors.Errorf("chart handle %q already defined in this import session", chart.Handle)
	} else {
		cImp.dirty[chart.ID] = true
	}

	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "handle":
			// handle value sanity check
			if deinterfacer.ToString(val, handle) != handle {
				return fmt.Errorf("explicitly set handle on chart %q shadows inherited handle", handle)
			}

		case "name", "title", "label":
			chart.Name = deinterfacer.ToString(val)

		case "config":
			chart.Config, err = cImp.castConfig(chart, val)

		case "allow", "deny":
			return cImp.imp.permissions.CastSet(types.ChartPermissionResource.String()+handle, key, val)

		default:
			return fmt.Errorf("unexpected key %q for chart %q", key, handle)
		}

		return
	})
}

func (cImp *Chart) castConfig(chart *types.Chart, def interface{}) (types.ChartConfig, error) {
	var cfg = types.ChartConfig{}

	return cfg, deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "reports":
			cfg.Reports, err = cImp.castConfigReports(chart, val)

		default:
			return fmt.Errorf("unexpected key %q for chart %q config", key, chart.Handle)

		}
		return
	})
}

func (cImp *Chart) castConfigReports(chart *types.Chart, def interface{}) ([]*types.ChartConfigReport, error) {
	var rr = make([]*types.ChartConfigReport, 0)

	return rr, deinterfacer.Each(def, func(_ int, _ string, report interface{}) (err error) {
		var r = &types.ChartConfigReport{}
		err = deinterfacer.Each(report, func(_ int, key string, val interface{}) (err error) {
			switch key {
			case "filter":
				r.Filter = deinterfacer.ToString(val)
			case "module":
				module := deinterfacer.ToString(val)
				if m, err := cImp.getModule(module); err != nil || m == nil {
					return fmt.Errorf("unknown module %q referenced from chart %q report config", module, chart.Handle)
				}
				cImp.modRefs = append(cImp.modRefs, chartModuleRef{chart.Handle, len(rr), module})
			case "metrics":
				r.Metrics = deinterfacer.ToSliceOfStringToInterfaceMap(val)
			case "dimensions":
				r.Dimensions = deinterfacer.ToSliceOfStringToInterfaceMap(val)
			case "renderer":
				// @todo implement renderer decoding
				// r.Renderer.Version
			default:
				return fmt.Errorf("unexpected key %q for chart %q report config", key, chart.Handle)

			}

			return
		})

		if err != nil {
			return
		}

		rr = append(rr, r)
		return
	})
}

// Get existing charts
func (cImp *Chart) Get(handle string) (*types.Chart, error) {
	handle = importer.NormalizeHandle(handle)
	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid chart handle")
	}

	return cImp.set.FindByHandle(handle), nil
}

func (cImp *Chart) Store(ctx context.Context, k chartKeeper) (err error) {
	if err = cImp.resolveRefs(); err != nil {
		return
	}

	return cImp.set.Walk(func(chart *types.Chart) (err error) {
		var handle = chart.Handle

		if chart.ID == 0 {
			chart.NamespaceID = cImp.namespace.ID
			chart, err = k.Create(chart)
		} else if cImp.dirty[chart.ID] {
			chart, err = k.Update(chart)
		}

		if err != nil {
			return
		}

		cImp.dirty[chart.ID] = false
		cImp.imp.permissions.UpdateResources(types.ChartPermissionResource.String(), handle, chart.ID)

		return
	})
}

// Resolve all refs for this page (page module, inside block)
func (cImp *Chart) resolveRefs() error {

	for _, ref := range cImp.modRefs {
		chart := cImp.set.FindByHandle(ref.ch)
		if chart == nil {
			return errors.Errorf("invalid reference, unknown chart (%v)", ref)
		}

		if ref.ri > len(chart.Config.Reports) {
			return errors.Errorf("invalid reference, report index out of range (%v)", ref)
		}

		if module, err := cImp.getModule(ref.mh); err != nil {
			return errors.Errorf("invalid reference, module loading error: %v", err)
		} else if module == nil {
			return errors.Errorf("invalid reference, unknown module (%v)", ref)
		} else {
			chart.Config.Reports[ref.ri].ModuleID = module.ID
		}
	}

	return nil
}
