package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/compose_records.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	composeRecordCreate struct {
		Done chan struct{}
		rec  *types.Record
		mod  *types.Module
		err  error
	}

	composeRecordUpdate struct {
		Done chan struct{}
		rec  *types.Record
		mod  *types.Module
		err  error
	}

	composeRecordRemove struct {
		Done chan struct{}
		rec  *types.Record
		mod  *types.Module
		err  error
	}

	composeRecordStore interface {
		CreateComposeRecord(context.Context, *types.Module, ...*types.Record) error
		UpdateComposeRecord(context.Context, *types.Module, ...*types.Record) error
		RemoveComposeRecord(context.Context, *types.Module, ...*types.Record) error
	}
)

// CreateComposeRecord creates a new ComposeRecord
// create job that can be pushed to store's transaction handler
func CreateComposeRecord(mod *types.Module, rec *types.Record) *composeRecordCreate {
	return &composeRecordCreate{mod: mod, rec: rec}
}

// Do Executes composeRecordCreate job
func (j *composeRecordCreate) Do(ctx context.Context, s composeRecordStore) error {
	j.err = s.CreateComposeRecord(ctx, j.mod, j.rec)
	j.Done <- struct{}{}
	return j.err
}

// UpdateComposeRecord creates a new ComposeRecord
// update job that can be pushed to store's transaction handler
func UpdateComposeRecord(mod *types.Module, rec *types.Record) *composeRecordUpdate {
	return &composeRecordUpdate{mod: mod, rec: rec}
}

// Do Executes composeRecordUpdate job
func (j *composeRecordUpdate) Do(ctx context.Context, s composeRecordStore) error {
	j.err = s.UpdateComposeRecord(ctx, j.mod, j.rec)
	j.Done <- struct{}{}
	return j.err
}

// RemoveComposeRecord creates a new ComposeRecord
// remove job that can be pushed to store's transaction handler
func RemoveComposeRecord(mod *types.Module, rec *types.Record) *composeRecordRemove {
	return &composeRecordRemove{mod: mod, rec: rec}
}

// Do Executes composeRecordRemove job
func (j *composeRecordRemove) Do(ctx context.Context, s composeRecordStore) error {
	j.err = s.RemoveComposeRecord(ctx, j.mod, j.rec)
	j.Done <- struct{}{}
	return j.err
}
