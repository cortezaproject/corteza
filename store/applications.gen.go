package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/applications.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Applications interface {
		SearchApplications(ctx context.Context, f types.ApplicationFilter) (types.ApplicationSet, types.ApplicationFilter, error)
		LookupApplicationByID(ctx context.Context, id uint64) (*types.Application, error)

		CreateApplication(ctx context.Context, rr ...*types.Application) error

		UpdateApplication(ctx context.Context, rr ...*types.Application) error

		UpsertApplication(ctx context.Context, rr ...*types.Application) error

		DeleteApplication(ctx context.Context, rr ...*types.Application) error
		DeleteApplicationByID(ctx context.Context, ID uint64) error

		TruncateApplications(ctx context.Context) error

		// Additional custom functions

		// ApplicationMetrics (custom function)
		ApplicationMetrics(ctx context.Context) (*types.ApplicationMetrics, error)

		// ReorderApplications (custom function)
		ReorderApplications(ctx context.Context, _order []uint64) error
	}
)

var _ *types.Application
var _ context.Context

// SearchApplications returns all matching Applications from store
func SearchApplications(ctx context.Context, s Applications, f types.ApplicationFilter) (types.ApplicationSet, types.ApplicationFilter, error) {
	return s.SearchApplications(ctx, f)
}

// LookupApplicationByID searches for application by ID
//
// It returns application even if deleted
func LookupApplicationByID(ctx context.Context, s Applications, id uint64) (*types.Application, error) {
	return s.LookupApplicationByID(ctx, id)
}

// CreateApplication creates one or more Applications in store
func CreateApplication(ctx context.Context, s Applications, rr ...*types.Application) error {
	return s.CreateApplication(ctx, rr...)
}

// UpdateApplication updates one or more (existing) Applications in store
func UpdateApplication(ctx context.Context, s Applications, rr ...*types.Application) error {
	return s.UpdateApplication(ctx, rr...)
}

// UpsertApplication creates new or updates existing one or more Applications in store
func UpsertApplication(ctx context.Context, s Applications, rr ...*types.Application) error {
	return s.UpsertApplication(ctx, rr...)
}

// DeleteApplication Deletes one or more Applications from store
func DeleteApplication(ctx context.Context, s Applications, rr ...*types.Application) error {
	return s.DeleteApplication(ctx, rr...)
}

// DeleteApplicationByID Deletes Application from store
func DeleteApplicationByID(ctx context.Context, s Applications, ID uint64) error {
	return s.DeleteApplicationByID(ctx, ID)
}

// TruncateApplications Deletes all Applications from store
func TruncateApplications(ctx context.Context, s Applications) error {
	return s.TruncateApplications(ctx)
}

func ApplicationMetrics(ctx context.Context, s Applications) (*types.ApplicationMetrics, error) {
	return s.ApplicationMetrics(ctx)
}

func ReorderApplications(ctx context.Context, s Applications, _order []uint64) error {
	return s.ReorderApplications(ctx, _order)
}
