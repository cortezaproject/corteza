package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/workflow_triggers.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = errors.Is

// SearchWorkflowTriggers returns all matching rows
//
// This function calls convertWorkflowTriggerFilter with the given
// types.WorkflowTriggerFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchWorkflowTriggers(ctx context.Context, f types.WorkflowTriggerFilter) (types.WorkflowTriggerSet, types.WorkflowTriggerFilter, error) {
	var (
		err error
		set []*types.WorkflowTrigger
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q = s.workflowTriggersSelectBuilder()

		set, err = s.QueryWorkflowTriggers(ctx, q, nil)
		return err
	}()
}

// QueryWorkflowTriggers queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryWorkflowTriggers(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.WorkflowTrigger) (bool, error),
) ([]*types.WorkflowTrigger, error) {
	var (
		set = make([]*types.WorkflowTrigger, 0, DefaultSliceCapacity)
		res *types.WorkflowTrigger

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalWorkflowTriggerRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// CreateWorkflowTrigger creates one or more rows in workflow_triggers table
func (s Store) CreateWorkflowTrigger(ctx context.Context, rr ...*types.WorkflowTrigger) (err error) {
	for _, res := range rr {
		err = s.checkWorkflowTriggerConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateWorkflowTriggers(ctx, s.internalWorkflowTriggerEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateWorkflowTrigger updates one or more existing rows in workflow_triggers
func (s Store) UpdateWorkflowTrigger(ctx context.Context, rr ...*types.WorkflowTrigger) error {
	return s.partialWorkflowTriggerUpdate(ctx, nil, rr...)
}

// partialWorkflowTriggerUpdate updates one or more existing rows in workflow_triggers
func (s Store) partialWorkflowTriggerUpdate(ctx context.Context, onlyColumns []string, rr ...*types.WorkflowTrigger) (err error) {
	for _, res := range rr {
		err = s.checkWorkflowTriggerConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateWorkflowTriggers(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("wftrg.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalWorkflowTriggerEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertWorkflowTrigger updates one or more existing rows in workflow_triggers
func (s Store) UpsertWorkflowTrigger(ctx context.Context, rr ...*types.WorkflowTrigger) (err error) {
	for _, res := range rr {
		err = s.checkWorkflowTriggerConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertWorkflowTriggers(ctx, s.internalWorkflowTriggerEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteWorkflowTrigger Deletes one or more rows from workflow_triggers table
func (s Store) DeleteWorkflowTrigger(ctx context.Context, rr ...*types.WorkflowTrigger) (err error) {
	for _, res := range rr {

		err = s.execDeleteWorkflowTriggers(ctx, squirrel.Eq{
			s.preprocessColumn("wftrg.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteWorkflowTriggerByID Deletes row from the workflow_triggers table
func (s Store) DeleteWorkflowTriggerByID(ctx context.Context, ID uint64) error {
	return s.execDeleteWorkflowTriggers(ctx, squirrel.Eq{
		s.preprocessColumn("wftrg.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateWorkflowTriggers Deletes all rows from the workflow_triggers table
func (s Store) TruncateWorkflowTriggers(ctx context.Context) error {
	return s.Truncate(ctx, s.workflowTriggerTable())
}

// execLookupWorkflowTrigger prepares WorkflowTrigger query and executes it,
// returning types.WorkflowTrigger (or error)
func (s Store) execLookupWorkflowTrigger(ctx context.Context, cnd squirrel.Sqlizer) (res *types.WorkflowTrigger, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.workflowTriggersSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalWorkflowTriggerRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateWorkflowTriggers updates all matched (by cnd) rows in workflow_triggers with given data
func (s Store) execCreateWorkflowTriggers(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.workflowTriggerTable()).SetMap(payload))
}

// execUpdateWorkflowTriggers updates all matched (by cnd) rows in workflow_triggers with given data
func (s Store) execUpdateWorkflowTriggers(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.workflowTriggerTable("wftrg")).Where(cnd).SetMap(set))
}

// execUpsertWorkflowTriggers inserts new or updates matching (by-primary-key) rows in workflow_triggers with given data
func (s Store) execUpsertWorkflowTriggers(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.workflowTriggerTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteWorkflowTriggers Deletes all matched (by cnd) rows in workflow_triggers with given data
func (s Store) execDeleteWorkflowTriggers(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.workflowTriggerTable("wftrg")).Where(cnd))
}

func (s Store) internalWorkflowTriggerRowScanner(row rowScanner) (res *types.WorkflowTrigger, err error) {
	res = &types.WorkflowTrigger{}

	if _, has := s.config.RowScanners["workflowTrigger"]; has {
		scanner := s.config.RowScanners["workflowTrigger"].(func(_ rowScanner, _ *types.WorkflowTrigger) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.WorkflowID,
			&res.Meta,
			&res.Enabled,
			&res.ResourceType,
			&res.EventType,
			&res.Constraints,
			&res.Input,
			&res.OwnedBy,
			&res.CreatedBy,
			&res.UpdatedBy,
			&res.DeletedBy,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan workflowTrigger db row").Wrap(err)
	} else {
		return res, nil
	}
}

// QueryWorkflowTriggers returns squirrel.SelectBuilder with set table and all columns
func (s Store) workflowTriggersSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.workflowTriggerTable("wftrg"), s.workflowTriggerColumns("wftrg")...)
}

// workflowTriggerTable name of the db table
func (Store) workflowTriggerTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "workflow_triggers" + alias
}

// WorkflowTriggerColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) workflowTriggerColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_workflow",
		alias + "meta",
		alias + "enabled",
		alias + "resource_type",
		alias + "event_type",
		alias + "constraints",
		alias + "input",
		alias + "owned_by",
		alias + "created_by",
		alias + "updated_by",
		alias + "deleted_by",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false false false false}

// internalWorkflowTriggerEncoder encodes fields from types.WorkflowTrigger to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeWorkflowTrigger
// func when rdbms.customEncoder=true
func (s Store) internalWorkflowTriggerEncoder(res *types.WorkflowTrigger) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"rel_workflow":  res.WorkflowID,
		"meta":          res.Meta,
		"enabled":       res.Enabled,
		"resource_type": res.ResourceType,
		"event_type":    res.EventType,
		"constraints":   res.Constraints,
		"input":         res.Input,
		"owned_by":      res.OwnedBy,
		"created_by":    res.CreatedBy,
		"updated_by":    res.UpdatedBy,
		"deleted_by":    res.DeletedBy,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}

// checkWorkflowTriggerConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkWorkflowTriggerConstraints(ctx context.Context, res *types.WorkflowTrigger) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	if !valid {
		return nil
	}

	return nil
}
