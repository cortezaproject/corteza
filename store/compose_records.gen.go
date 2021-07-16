package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/compose_records.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/report"
)

type (
	ComposeRecords interface {
		SearchComposeRecords(ctx context.Context, _mod *types.Module, f types.RecordFilter) (types.RecordSet, types.RecordFilter, error)
		LookupComposeRecordByID(ctx context.Context, _mod *types.Module, id uint64) (*types.Record, error)

		CreateComposeRecord(ctx context.Context, _mod *types.Module, rr ...*types.Record) error

		UpdateComposeRecord(ctx context.Context, _mod *types.Module, rr ...*types.Record) error

		UpsertComposeRecord(ctx context.Context, _mod *types.Module, rr ...*types.Record) error

		DeleteComposeRecord(ctx context.Context, _mod *types.Module, rr ...*types.Record) error
		DeleteComposeRecordByID(ctx context.Context, _mod *types.Module, ID uint64) error

		TruncateComposeRecords(ctx context.Context, _mod *types.Module) error

		// Additional custom functions

		// ComposeRecordReport (custom function)
		ComposeRecordReport(ctx context.Context, _mod *types.Module, _metrics string, _dimensions string, _filters string) ([]map[string]interface{}, error)

		// ComposeRecordDatasource (custom function)
		ComposeRecordDatasource(ctx context.Context, _mod *types.Module, _ld *report.LoadStepDefinition) (report.Datasource, error)

		// PartialComposeRecordValueUpdate (custom function)
		PartialComposeRecordValueUpdate(ctx context.Context, _mod *types.Module, _values ...*types.RecordValue) error
	}
)

var _ *types.Record
var _ context.Context

// SearchComposeRecords returns all matching ComposeRecords from store
func SearchComposeRecords(ctx context.Context, s ComposeRecords, _mod *types.Module, f types.RecordFilter) (types.RecordSet, types.RecordFilter, error) {
	return s.SearchComposeRecords(ctx, _mod, f)
}

// LookupComposeRecordByID searches for compose record by ID
// It returns compose record even if deleted
func LookupComposeRecordByID(ctx context.Context, s ComposeRecords, _mod *types.Module, id uint64) (*types.Record, error) {
	return s.LookupComposeRecordByID(ctx, _mod, id)
}

// CreateComposeRecord creates one or more ComposeRecords in store
func CreateComposeRecord(ctx context.Context, s ComposeRecords, _mod *types.Module, rr ...*types.Record) error {
	return s.CreateComposeRecord(ctx, _mod, rr...)
}

// UpdateComposeRecord updates one or more (existing) ComposeRecords in store
func UpdateComposeRecord(ctx context.Context, s ComposeRecords, _mod *types.Module, rr ...*types.Record) error {
	return s.UpdateComposeRecord(ctx, _mod, rr...)
}

// UpsertComposeRecord creates new or updates existing one or more ComposeRecords in store
func UpsertComposeRecord(ctx context.Context, s ComposeRecords, _mod *types.Module, rr ...*types.Record) error {
	return s.UpsertComposeRecord(ctx, _mod, rr...)
}

// DeleteComposeRecord Deletes one or more ComposeRecords from store
func DeleteComposeRecord(ctx context.Context, s ComposeRecords, _mod *types.Module, rr ...*types.Record) error {
	return s.DeleteComposeRecord(ctx, _mod, rr...)
}

// DeleteComposeRecordByID Deletes ComposeRecord from store
func DeleteComposeRecordByID(ctx context.Context, s ComposeRecords, _mod *types.Module, ID uint64) error {
	return s.DeleteComposeRecordByID(ctx, _mod, ID)
}

// TruncateComposeRecords Deletes all ComposeRecords from store
func TruncateComposeRecords(ctx context.Context, s ComposeRecords, _mod *types.Module) error {
	return s.TruncateComposeRecords(ctx, _mod)
}

func ComposeRecordReport(ctx context.Context, s ComposeRecords, _mod *types.Module, _metrics string, _dimensions string, _filters string) ([]map[string]interface{}, error) {
	return s.ComposeRecordReport(ctx, _mod, _metrics, _dimensions, _filters)
}

func ComposeRecordDatasource(ctx context.Context, s ComposeRecords, _mod *types.Module, _ld *report.LoadStepDefinition) (report.Datasource, error) {
	return s.ComposeRecordDatasource(ctx, _mod, _ld)
}

func PartialComposeRecordValueUpdate(ctx context.Context, s ComposeRecords, _mod *types.Module, _values ...*types.RecordValue) error {
	return s.PartialComposeRecordValueUpdate(ctx, _mod, _values...)
}
