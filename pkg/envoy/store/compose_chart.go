package store

import (
	"context"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composeChartState struct {
		cfg *EncoderConfig

		res *resource.ComposeChart
		chr *types.Chart

		relNS   *types.Namespace
		relMods types.ModuleSet
	}
)

func NewComposeChartState(res *resource.ComposeChart, cfg *EncoderConfig) resourceState {
	return &composeChartState{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,
	}
}

func (n *composeChartState) Prepare(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	// Initial values
	if n.res.Res.CreatedAt.IsZero() {
		n.res.Res.CreatedAt = time.Now()
	}

	// Get relate namespace
	n.relNS, err = findComposeNamespaceRS(ctx, s, state.ParentResources, n.res.NsRef.Identifiers)
	if err != nil {
		return err
	}
	if n.relNS == nil {
		return composeNamespaceErrUnresolved(n.res.NsRef.Identifiers)
	}

	// Get related modules
	n.relMods = make(types.ModuleSet, len(n.res.ModRef))
	for i, mRef := range n.res.ModRef {
		n.relMods[i], err = findComposeModuleRS(ctx, s, n.relNS.ID, state.ParentResources, mRef.Identifiers)
		if err != nil {
			return err
		}
		if n.relMods[i] == nil {
			return composeModuleErrUnresolved(mRef.Identifiers)
		}
	}

	// Try to get the original chart
	n.chr, err = findComposeChartS(ctx, s, n.relNS.ID, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.chr != nil {
		n.res.Res.ID = n.chr.ID
		n.res.Res.NamespaceID = n.chr.NamespaceID
	}
	return nil
}

func (n *composeChartState) Encode(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	res := n.res.Res
	exists := n.chr != nil && n.chr.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.chr.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	// This is not possible, but let's do it anyway
	if state.Conflicting {
		return nil
	}

	// Timestamps
	ts := n.res.Timestamps()
	if ts != nil {
		if ts.CreatedAt != "" {
			t := toTime(ts.CreatedAt)
			if t != nil {
				res.CreatedAt = *t
			}
		}
		if ts.UpdatedAt != "" {
			res.UpdatedAt = toTime(ts.UpdatedAt)
		}
		if ts.DeletedAt != "" {
			res.DeletedAt = toTime(ts.DeletedAt)
		}
	}

	// Namespace
	res.NamespaceID = n.relNS.ID
	if res.NamespaceID <= 0 {
		ns := findComposeNamespaceR(state.ParentResources, n.res.NsRef.Identifiers)
		res.NamespaceID = ns.ID
	}
	if res.NamespaceID <= 0 {
		return composeNamespaceErrUnresolved(n.res.NsRef.Identifiers)
	}

	// Report modules
	for i, r := range res.Config.Reports {
		mod := n.relMods[i]
		if mod == nil {
			mod = findComposeModuleR(state.ParentResources, n.res.ModRef[i].Identifiers)
		}
		if mod == nil || mod.ID <= 0 {
			return composeModuleErrUnresolved(n.res.ModRef[i].Identifiers)
		}

		r.ModuleID = mod.ID
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to compose/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create fresh chart
	if !exists {
		return store.CreateComposeChart(ctx, s, res)
	}

	// Update existing chart
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeComposeChart(n.chr, res)

	case resource.MergeRight:
		res = mergeComposeChart(res, n.chr)
	}

	err = store.UpdateComposeChart(ctx, s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}

// mergeComposeChart merges b into a, prioritising a
func mergeComposeChart(a, b *types.Chart) *types.Chart {
	c := *a

	if c.Handle == "" {
		c.Handle = b.Handle
	}
	if c.Name == "" {
		c.Name = b.Name
	}
	c.NamespaceID = b.NamespaceID
	if c.CreatedAt.IsZero() {
		c.CreatedAt = b.CreatedAt
	}
	if c.UpdatedAt == nil {
		c.UpdatedAt = b.UpdatedAt
	}
	if c.DeletedAt == nil {
		c.DeletedAt = b.DeletedAt
	}

	if len(c.Config.Reports) <= 0 {
		c.Config.Reports = b.Config.Reports
	}
	if c.Config.ColorScheme == "" {
		c.Config.ColorScheme = b.Config.ColorScheme
	}

	return &c
}

// findComposeChartRS looks for the chart in the resources & the store
//
// Provided resources are prioritized.
func findComposeChartRS(ctx context.Context, s store.Storer, nsID uint64, rr resource.InterfaceSet, ii resource.Identifiers) (ch *types.Chart, err error) {
	ch = findComposeChartR(rr, ii)
	if ch != nil {
		return ch, nil
	}

	if nsID <= 0 {
		return nil, nil
	}

	// Go in the store
	return findComposeChartS(ctx, s, nsID, makeGenericFilter(ii))
}

// findComposeChartS looks for the chart in the store
func findComposeChartS(ctx context.Context, s store.Storer, nsID uint64, gf genericFilter) (ch *types.Chart, err error) {
	if nsID == 0 {
		return nil, nil
	}

	if gf.id > 0 {
		ch, err = store.LookupComposeChartByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if ch != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		ch, err = store.LookupComposeChartByNamespaceIDHandle(ctx, s, nsID, i)
		if err == store.ErrNotFound {
			var cc types.ChartSet
			cc, _, err = store.SearchComposeCharts(ctx, s, types.ChartFilter{
				NamespaceID: nsID,
				Name:        i,
				Paging: filter.Paging{
					Limit: 2,
				},
			})
			if len(cc) > 1 {
				return nil, resourceErrIdentifierNotUnique(i)
			}
			if len(cc) == 1 {
				ch = cc[0]
			}
		}

		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if ch != nil {
			return
		}
	}

	return nil, nil
}

// findComposeChartR looks for the chart in the resources
func findComposeChartR(rr resource.InterfaceSet, ii resource.Identifiers) (ch *types.Chart) {
	var chRes *resource.ComposeChart

	rr.Walk(func(r resource.Interface) error {
		cr, ok := r.(*resource.ComposeChart)
		if !ok {
			return nil
		}

		if cr.Identifiers().HasAny(ii) {
			chRes = cr
		}
		return nil
	})

	// Found it
	if chRes != nil {
		return chRes.Res
	}

	return nil
}

func composeChartErrUnresolved(ii resource.Identifiers) error {
	return fmt.Errorf("compose chart unresolved %v", ii.StringSlice())
}
