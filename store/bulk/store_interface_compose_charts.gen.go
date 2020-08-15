package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/compose_charts.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	composeChartsStore interface {
		SearchComposeCharts(ctx context.Context, f types.ChartFilter) (types.ChartSet, types.ChartFilter, error)
		LookupComposeChartByID(ctx context.Context, id uint64) (*types.Chart, error)
		LookupComposeChartByHandle(ctx context.Context, handle string) (*types.Chart, error)
		CreateComposeChart(ctx context.Context, rr ...*types.Chart) error
		UpdateComposeChart(ctx context.Context, rr ...*types.Chart) error
		PartialUpdateComposeChart(ctx context.Context, onlyColumns []string, rr ...*types.Chart) error
		RemoveComposeChart(ctx context.Context, rr ...*types.Chart) error
		RemoveComposeChartByID(ctx context.Context, ID uint64) error

		TruncateComposeCharts(ctx context.Context) error
	}
)
