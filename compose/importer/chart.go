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
	ChartImport struct {
		namespace *types.Namespace
		set       types.ChartSet

		permissions importer.PermissionImporter

		finder chartFinder
	}

	chartFinder interface {
		FindByHandle(uint64, string) (*types.Chart, error)
	}
)

func NewChartImporter(ns *types.Namespace, f chartFinder, p importer.PermissionImporter) *ChartImport {
	return &ChartImport{
		namespace:   ns,
		set:         types.ChartSet{},
		finder:      f,
		permissions: p,
	}
}

// CastSet Resolves permission rules:
// { <chart-handle>: { chart } } or [ { chart }, ... ]
func (imp *ChartImport) CastSet(set interface{}) error {
	return deinterfacer.Each(set, func(index int, handle string, def interface{}) error {
		if index > -1 {
			// Charts defined as collection
			deinterfacer.KVsetString(&handle, "handle", def)
		}

		return imp.Cast(handle, def)
	})
}

// Cast Resolves permission rules:
// { <chart-handle>: { chart } } or [ { chart }, ... ]
func (imp *ChartImport) Cast(handle string, def interface{}) (err error) {
	if !deinterfacer.IsMap(def) {
		return errors.New("expecting map of values for chart")
	}

	var chart *types.Chart

	if !importer.IsValidHandle(handle) {
		return errors.New("invalid chart handle")
	}

	handle = importer.NormalizeHandle(handle)
	if chart, err = imp.Get(handle); err != nil {
		return err
	}

	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "namespace":
			// namespace value sanity check
			if deinterfacer.ToString(val, imp.namespace.Slug) != imp.namespace.Slug {
				return fmt.Errorf("explicitly set namespace on chart %q shadows inherited namespace", imp.namespace.Slug)
			}

		case "handle":
			// handle value sanity check
			if deinterfacer.ToString(val, handle) != handle {
				return fmt.Errorf("explicitly set handle on chart %q shadows inherited handle", handle)
			}

		case "name":
			chart.Name = deinterfacer.ToString(val)

		case "config":
			// @todo Chart.Config

		case "allow", "deny":
			return imp.permissions.CastSet(types.ChartPermissionResource.String()+handle, key, val)

		default:
			return fmt.Errorf("unexpected key %q for chart %q", key, handle)
		}

		return err
	})
}

func (imp *ChartImport) Exists(handle string) bool {
	handle = importer.NormalizeHandle(handle)

	var (
		chart *types.Chart
		err   error
	)

	chart = imp.set.FindByHandle(handle)
	if chart != nil {
		return true
	}

	if imp.namespace.ID == 0 {
		// Assuming new namespace, nothing exists yet..
		return false
	}

	if imp.finder != nil {
		chart, err = imp.finder.FindByHandle(imp.namespace.ID, handle)
		if err == nil && chart != nil {
			imp.set = append(imp.set, chart)
			return true
		}
	}

	return false
}

// finds or makes new chart
func (imp *ChartImport) Get(handle string) (*types.Chart, error) {
	handle = importer.NormalizeHandle(handle)

	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid chart handle")
	}

	if !imp.Exists(handle) {
		imp.set = append(imp.set, &types.Chart{
			Handle: handle,
			Name:   handle,
		})
	}

	return imp.set.FindByHandle(handle), nil
}

func (imp *ChartImport) Store(ctx context.Context, k chartKeeper) error {
	return imp.set.Walk(func(chart *types.Chart) (err error) {
		var handle = chart.Handle

		if chart.ID == 0 {
			chart.NamespaceID = imp.namespace.ID
			chart, err = k.Create(chart)
		} else {
			chart, err = k.Update(chart)
		}

		if err != nil {
			return
		}
		// @todo update module ref for charts

		imp.permissions.UpdateResources(types.ChartPermissionResource.String(), handle, chart.ID)

		return
	})
}
