package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/reports.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Reports interface {
		SearchReports(ctx context.Context, f types.ReportFilter) (types.ReportSet, types.ReportFilter, error)
		LookupReportByID(ctx context.Context, id uint64) (*types.Report, error)
		LookupReportByHandle(ctx context.Context, handle string) (*types.Report, error)

		CreateReport(ctx context.Context, rr ...*types.Report) error

		UpdateReport(ctx context.Context, rr ...*types.Report) error

		UpsertReport(ctx context.Context, rr ...*types.Report) error

		DeleteReport(ctx context.Context, rr ...*types.Report) error
		DeleteReportByID(ctx context.Context, ID uint64) error

		TruncateReports(ctx context.Context) error
	}
)

var _ *types.Report
var _ context.Context

// SearchReports returns all matching Reports from store
func SearchReports(ctx context.Context, s Reports, f types.ReportFilter) (types.ReportSet, types.ReportFilter, error) {
	return s.SearchReports(ctx, f)
}

// LookupReportByID searches for report by ID
//
// It returns report even if deleted
func LookupReportByID(ctx context.Context, s Reports, id uint64) (*types.Report, error) {
	return s.LookupReportByID(ctx, id)
}

// LookupReportByHandle searches for report by Handle
//
// It returns report even if deleted
func LookupReportByHandle(ctx context.Context, s Reports, handle string) (*types.Report, error) {
	return s.LookupReportByHandle(ctx, handle)
}

// CreateReport creates one or more Reports in store
func CreateReport(ctx context.Context, s Reports, rr ...*types.Report) error {
	return s.CreateReport(ctx, rr...)
}

// UpdateReport updates one or more (existing) Reports in store
func UpdateReport(ctx context.Context, s Reports, rr ...*types.Report) error {
	return s.UpdateReport(ctx, rr...)
}

// UpsertReport creates new or updates existing one or more Reports in store
func UpsertReport(ctx context.Context, s Reports, rr ...*types.Report) error {
	return s.UpsertReport(ctx, rr...)
}

// DeleteReport Deletes one or more Reports from store
func DeleteReport(ctx context.Context, s Reports, rr ...*types.Report) error {
	return s.DeleteReport(ctx, rr...)
}

// DeleteReportByID Deletes Report from store
func DeleteReportByID(ctx context.Context, s Reports, ID uint64) error {
	return s.DeleteReportByID(ctx, ID)
}

// TruncateReports Deletes all Reports from store
func TruncateReports(ctx context.Context, s Reports) error {
	return s.TruncateReports(ctx)
}
