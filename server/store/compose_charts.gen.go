package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/compose_charts.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeCharts interface {
		SearchComposeCharts(ctx context.Context, f types.ChartFilter) (types.ChartSet, types.ChartFilter, error)
		LookupComposeChartByID(ctx context.Context, id uint64) (*types.Chart, error)
		LookupComposeChartByNamespaceIDHandle(ctx context.Context, namespace_id uint64, handle string) (*types.Chart, error)

		CreateComposeChart(ctx context.Context, rr ...*types.Chart) error

		UpdateComposeChart(ctx context.Context, rr ...*types.Chart) error

		UpsertComposeChart(ctx context.Context, rr ...*types.Chart) error

		DeleteComposeChart(ctx context.Context, rr ...*types.Chart) error
		DeleteComposeChartByID(ctx context.Context, ID uint64) error

		TruncateComposeCharts(ctx context.Context) error
	}
)

var _ *types.Chart
var _ context.Context

// SearchComposeCharts returns all matching ComposeCharts from store
func SearchComposeCharts(ctx context.Context, s ComposeCharts, f types.ChartFilter) (types.ChartSet, types.ChartFilter, error) {
	return s.SearchComposeCharts(ctx, f)
}

// LookupComposeChartByID searches for compose chart by ID
//
// It returns compose chart even if deleted
func LookupComposeChartByID(ctx context.Context, s ComposeCharts, id uint64) (*types.Chart, error) {
	return s.LookupComposeChartByID(ctx, id)
}

// LookupComposeChartByNamespaceIDHandle searches for compose chart by handle (case-insensitive)
func LookupComposeChartByNamespaceIDHandle(ctx context.Context, s ComposeCharts, namespace_id uint64, handle string) (*types.Chart, error) {
	return s.LookupComposeChartByNamespaceIDHandle(ctx, namespace_id, handle)
}

// CreateComposeChart creates one or more ComposeCharts in store
func CreateComposeChart(ctx context.Context, s ComposeCharts, rr ...*types.Chart) error {
	return s.CreateComposeChart(ctx, rr...)
}

// UpdateComposeChart updates one or more (existing) ComposeCharts in store
func UpdateComposeChart(ctx context.Context, s ComposeCharts, rr ...*types.Chart) error {
	return s.UpdateComposeChart(ctx, rr...)
}

// UpsertComposeChart creates new or updates existing one or more ComposeCharts in store
func UpsertComposeChart(ctx context.Context, s ComposeCharts, rr ...*types.Chart) error {
	return s.UpsertComposeChart(ctx, rr...)
}

// DeleteComposeChart Deletes one or more ComposeCharts from store
func DeleteComposeChart(ctx context.Context, s ComposeCharts, rr ...*types.Chart) error {
	return s.DeleteComposeChart(ctx, rr...)
}

// DeleteComposeChartByID Deletes ComposeChart from store
func DeleteComposeChartByID(ctx context.Context, s ComposeCharts, ID uint64) error {
	return s.DeleteComposeChartByID(ctx, ID)
}

// TruncateComposeCharts Deletes all ComposeCharts from store
func TruncateComposeCharts(ctx context.Context, s ComposeCharts) error {
	return s.TruncateComposeCharts(ctx)
}
