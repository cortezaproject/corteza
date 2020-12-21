package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/workflows.yaml
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

// SearchWorkflows returns all matching rows
//
// This function calls convertWorkflowFilter with the given
// types.WorkflowFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchWorkflows(ctx context.Context, f types.WorkflowFilter) (types.WorkflowSet, types.WorkflowFilter, error) {
	var (
		err error
		set []*types.Workflow
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q = s.workflowsSelectBuilder()

		set, err = s.QueryWorkflows(ctx, q, nil)
		return err
	}()
}

// QueryWorkflows queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryWorkflows(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Workflow) (bool, error),
) ([]*types.Workflow, error) {
	var (
		set = make([]*types.Workflow, 0, DefaultSliceCapacity)
		res *types.Workflow

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalWorkflowRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupWorkflowByID searches for workflow by ID
//
// It returns workflow even if deleted or suspended
func (s Store) LookupWorkflowByID(ctx context.Context, id uint64) (*types.Workflow, error) {
	return s.execLookupWorkflow(ctx, squirrel.Eq{
		s.preprocessColumn("wf.id", ""): store.PreprocessValue(id, ""),
	})
}

// LookupWorkflowByHandle searches for workflow by their handle
//
// It returns only valid workflows (not deleted, not suspended)
func (s Store) LookupWorkflowByHandle(ctx context.Context, handle string) (*types.Workflow, error) {
	return s.execLookupWorkflow(ctx, squirrel.Eq{
		s.preprocessColumn("wf.handle", "lower"): store.PreprocessValue(handle, "lower"),

		"wf.deleted_at": nil,
	})
}

// CreateWorkflow creates one or more rows in workflows table
func (s Store) CreateWorkflow(ctx context.Context, rr ...*types.Workflow) (err error) {
	for _, res := range rr {
		err = s.checkWorkflowConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateWorkflows(ctx, s.internalWorkflowEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateWorkflow updates one or more existing rows in workflows
func (s Store) UpdateWorkflow(ctx context.Context, rr ...*types.Workflow) error {
	return s.partialWorkflowUpdate(ctx, nil, rr...)
}

// partialWorkflowUpdate updates one or more existing rows in workflows
func (s Store) partialWorkflowUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Workflow) (err error) {
	for _, res := range rr {
		err = s.checkWorkflowConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateWorkflows(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("wf.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalWorkflowEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertWorkflow updates one or more existing rows in workflows
func (s Store) UpsertWorkflow(ctx context.Context, rr ...*types.Workflow) (err error) {
	for _, res := range rr {
		err = s.checkWorkflowConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertWorkflows(ctx, s.internalWorkflowEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteWorkflow Deletes one or more rows from workflows table
func (s Store) DeleteWorkflow(ctx context.Context, rr ...*types.Workflow) (err error) {
	for _, res := range rr {

		err = s.execDeleteWorkflows(ctx, squirrel.Eq{
			s.preprocessColumn("wf.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteWorkflowByID Deletes row from the workflows table
func (s Store) DeleteWorkflowByID(ctx context.Context, ID uint64) error {
	return s.execDeleteWorkflows(ctx, squirrel.Eq{
		s.preprocessColumn("wf.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateWorkflows Deletes all rows from the workflows table
func (s Store) TruncateWorkflows(ctx context.Context) error {
	return s.Truncate(ctx, s.workflowTable())
}

// execLookupWorkflow prepares Workflow query and executes it,
// returning types.Workflow (or error)
func (s Store) execLookupWorkflow(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Workflow, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.workflowsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalWorkflowRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateWorkflows updates all matched (by cnd) rows in workflows with given data
func (s Store) execCreateWorkflows(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.workflowTable()).SetMap(payload))
}

// execUpdateWorkflows updates all matched (by cnd) rows in workflows with given data
func (s Store) execUpdateWorkflows(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.workflowTable("wf")).Where(cnd).SetMap(set))
}

// execUpsertWorkflows inserts new or updates matching (by-primary-key) rows in workflows with given data
func (s Store) execUpsertWorkflows(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.workflowTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteWorkflows Deletes all matched (by cnd) rows in workflows with given data
func (s Store) execDeleteWorkflows(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.workflowTable("wf")).Where(cnd))
}

func (s Store) internalWorkflowRowScanner(row rowScanner) (res *types.Workflow, err error) {
	res = &types.Workflow{}

	if _, has := s.config.RowScanners["workflow"]; has {
		scanner := s.config.RowScanners["workflow"].(func(_ rowScanner, _ *types.Workflow) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Handle,
			&res.Meta,
			&res.Enabled,
			&res.Trace,
			&res.KeepSessions,
			&res.Scope,
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
		return nil, errors.Store("could not scan workflow db row").Wrap(err)
	} else {
		return res, nil
	}
}

// QueryWorkflows returns squirrel.SelectBuilder with set table and all columns
func (s Store) workflowsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.workflowTable("wf"), s.workflowColumns("wf")...)
}

// workflowTable name of the db table
func (Store) workflowTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "workflows" + alias
}

// WorkflowColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) workflowColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "handle",
		alias + "meta",
		alias + "enabled",
		alias + "trace",
		alias + "keep_sessions",
		alias + "scope",
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

// internalWorkflowEncoder encodes fields from types.Workflow to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeWorkflow
// func when rdbms.customEncoder=true
func (s Store) internalWorkflowEncoder(res *types.Workflow) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"handle":        res.Handle,
		"meta":          res.Meta,
		"enabled":       res.Enabled,
		"trace":         res.Trace,
		"keep_sessions": res.KeepSessions,
		"scope":         res.Scope,
		"owned_by":      res.OwnedBy,
		"created_by":    res.CreatedBy,
		"updated_by":    res.UpdatedBy,
		"deleted_by":    res.DeletedBy,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}

// checkWorkflowConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkWorkflowConstraints(ctx context.Context, res *types.Workflow) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	valid = valid && len(res.Handle) > 0

	if !valid {
		return nil
	}

	{
		ex, err := s.LookupWorkflowByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
