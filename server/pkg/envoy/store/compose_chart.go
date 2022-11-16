package store

import (
	"context"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
)

type (
	composeChart struct {
		cfg *EncoderConfig

		res *resource.ComposeChart
		chr *types.Chart

		relNS   *types.Namespace
		relMods types.ModuleSet
	}
)

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

// findComposeChart looks for the chart in the resources & the store
//
// Provided resources are prioritized.
func findComposeChart(ctx context.Context, s store.Storer, nsID uint64, rr resource.InterfaceSet, ii resource.Identifiers) (ch *types.Chart, err error) {
	ch = resource.FindComposeChart(rr, ii)
	if ch != nil {
		return ch, nil
	}

	if nsID <= 0 {
		return nil, nil
	}

	// Go in the store
	return findComposeChartStore(ctx, s, nsID, makeGenericFilter(ii))
}

// findComposeChartStore looks for the chart in the store
func findComposeChartStore(ctx context.Context, s store.Storer, nsID uint64, gf genericFilter) (ch *types.Chart, err error) {
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
