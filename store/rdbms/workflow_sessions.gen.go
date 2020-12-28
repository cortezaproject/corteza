package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/workflow_sessions.yaml
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

// SearchWorkflowSessions returns all matching rows
//
// This function calls convertWorkflowSessionFilter with the given
// types.WorkflowSessionFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchWorkflowSessions(ctx context.Context, f types.WorkflowSessionFilter) (types.WorkflowSessionSet, types.WorkflowSessionFilter, error) {
	var (
		err error
		set []*types.WorkflowSession
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q = s.workflowSessionsSelectBuilder()

		set, err = s.QueryWorkflowSessions(ctx, q, nil)
		return err
	}()
}

// QueryWorkflowSessions queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryWorkflowSessions(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.WorkflowSession) (bool, error),
) ([]*types.WorkflowSession, error) {
	var (
		set = make([]*types.WorkflowSession, 0, DefaultSliceCapacity)
		res *types.WorkflowSession

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalWorkflowSessionRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// CreateWorkflowSession creates one or more rows in workflow_sessions table
func (s Store) CreateWorkflowSession(ctx context.Context, rr ...*types.WorkflowSession) (err error) {
	for _, res := range rr {
		err = s.checkWorkflowSessionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateWorkflowSessions(ctx, s.internalWorkflowSessionEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateWorkflowSession updates one or more existing rows in workflow_sessions
func (s Store) UpdateWorkflowSession(ctx context.Context, rr ...*types.WorkflowSession) error {
	return s.partialWorkflowSessionUpdate(ctx, nil, rr...)
}

// partialWorkflowSessionUpdate updates one or more existing rows in workflow_sessions
func (s Store) partialWorkflowSessionUpdate(ctx context.Context, onlyColumns []string, rr ...*types.WorkflowSession) (err error) {
	for _, res := range rr {
		err = s.checkWorkflowSessionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateWorkflowSessions(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("wfses.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalWorkflowSessionEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertWorkflowSession updates one or more existing rows in workflow_sessions
func (s Store) UpsertWorkflowSession(ctx context.Context, rr ...*types.WorkflowSession) (err error) {
	for _, res := range rr {
		err = s.checkWorkflowSessionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertWorkflowSessions(ctx, s.internalWorkflowSessionEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteWorkflowSession Deletes one or more rows from workflow_sessions table
func (s Store) DeleteWorkflowSession(ctx context.Context, rr ...*types.WorkflowSession) (err error) {
	for _, res := range rr {

		err = s.execDeleteWorkflowSessions(ctx, squirrel.Eq{
			s.preprocessColumn("wfses.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteWorkflowSessionByID Deletes row from the workflow_sessions table
func (s Store) DeleteWorkflowSessionByID(ctx context.Context, ID uint64) error {
	return s.execDeleteWorkflowSessions(ctx, squirrel.Eq{
		s.preprocessColumn("wfses.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateWorkflowSessions Deletes all rows from the workflow_sessions table
func (s Store) TruncateWorkflowSessions(ctx context.Context) error {
	return s.Truncate(ctx, s.workflowSessionTable())
}

// execLookupWorkflowSession prepares WorkflowSession query and executes it,
// returning types.WorkflowSession (or error)
func (s Store) execLookupWorkflowSession(ctx context.Context, cnd squirrel.Sqlizer) (res *types.WorkflowSession, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.workflowSessionsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalWorkflowSessionRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateWorkflowSessions updates all matched (by cnd) rows in workflow_sessions with given data
func (s Store) execCreateWorkflowSessions(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.workflowSessionTable()).SetMap(payload))
}

// execUpdateWorkflowSessions updates all matched (by cnd) rows in workflow_sessions with given data
func (s Store) execUpdateWorkflowSessions(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.workflowSessionTable("wfses")).Where(cnd).SetMap(set))
}

// execUpsertWorkflowSessions inserts new or updates matching (by-primary-key) rows in workflow_sessions with given data
func (s Store) execUpsertWorkflowSessions(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.workflowSessionTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteWorkflowSessions Deletes all matched (by cnd) rows in workflow_sessions with given data
func (s Store) execDeleteWorkflowSessions(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.workflowSessionTable("wfses")).Where(cnd))
}

func (s Store) internalWorkflowSessionRowScanner(row rowScanner) (res *types.WorkflowSession, err error) {
	res = &types.WorkflowSession{}

	if _, has := s.config.RowScanners["workflowSession"]; has {
		scanner := s.config.RowScanners["workflowSession"].(func(_ rowScanner, _ *types.WorkflowSession) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.WorkflowID,
			&res.EventType,
			&res.EventResourceID,
			&res.ExecutedAs,
			&res.WallTime,
			&res.UserTime,
			&res.Input,
			&res.Output,
			&res.Trace,
			&res.CreatedBy,
			&res.DeletedBy,
			&res.CreatedAt,
			&res.DeletedAt,
			&res.PurgeAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan workflowSession db row").Wrap(err)
	} else {
		return res, nil
	}
}

// QueryWorkflowSessions returns squirrel.SelectBuilder with set table and all columns
func (s Store) workflowSessionsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.workflowSessionTable("wfses"), s.workflowSessionColumns("wfses")...)
}

// workflowSessionTable name of the db table
func (Store) workflowSessionTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "workflow_sessions" + alias
}

// WorkflowSessionColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) workflowSessionColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_workflow",
		alias + "event_type",
		alias + "rel_event_resource",
		alias + "executed_as",
		alias + "wall_time",
		alias + "user_time",
		alias + "input",
		alias + "output",
		alias + "trace",
		alias + "created_by",
		alias + "deleted_by",
		alias + "created_at",
		alias + "deleted_at",
		alias + "purge_at",
	}
}

// {true true false false false false}

// internalWorkflowSessionEncoder encodes fields from types.WorkflowSession to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeWorkflowSession
// func when rdbms.customEncoder=true
func (s Store) internalWorkflowSessionEncoder(res *types.WorkflowSession) store.Payload {
	return store.Payload{
		"id":                 res.ID,
		"rel_workflow":       res.WorkflowID,
		"event_type":         res.EventType,
		"rel_event_resource": res.EventResourceID,
		"executed_as":        res.ExecutedAs,
		"wall_time":          res.WallTime,
		"user_time":          res.UserTime,
		"input":              res.Input,
		"output":             res.Output,
		"trace":              res.Trace,
		"created_by":         res.CreatedBy,
		"deleted_by":         res.DeletedBy,
		"created_at":         res.CreatedAt,
		"deleted_at":         res.DeletedAt,
		"purge_at":           res.PurgeAt,
	}
}

// checkWorkflowSessionConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkWorkflowSessionConstraints(ctx context.Context, res *types.WorkflowSession) error {
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
