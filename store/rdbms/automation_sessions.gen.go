package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/automation_sessions.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchAutomationSessions returns all matching rows
//
// This function calls convertAutomationSessionFilter with the given
// types.SessionFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchAutomationSessions(ctx context.Context, f types.SessionFilter) (types.SessionSet, types.SessionFilter, error) {
	var (
		err error
		set []*types.Session
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q = s.automationSessionsSelectBuilder()

		set, err = s.QueryAutomationSessions(ctx, q, nil)
		return err
	}()
}

// QueryAutomationSessions queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryAutomationSessions(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Session) (bool, error),
) ([]*types.Session, error) {
	var (
		set = make([]*types.Session, 0, DefaultSliceCapacity)
		res *types.Session

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalAutomationSessionRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupAutomationSessionByID searches for session by ID
//
// It returns session even if deleted
func (s Store) LookupAutomationSessionByID(ctx context.Context, id uint64) (*types.Session, error) {
	return s.execLookupAutomationSession(ctx, squirrel.Eq{
		s.preprocessColumn("atms.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateAutomationSession creates one or more rows in automation_sessions table
func (s Store) CreateAutomationSession(ctx context.Context, rr ...*types.Session) (err error) {
	for _, res := range rr {
		err = s.checkAutomationSessionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateAutomationSessions(ctx, s.internalAutomationSessionEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateAutomationSession updates one or more existing rows in automation_sessions
func (s Store) UpdateAutomationSession(ctx context.Context, rr ...*types.Session) error {
	return s.partialAutomationSessionUpdate(ctx, nil, rr...)
}

// partialAutomationSessionUpdate updates one or more existing rows in automation_sessions
func (s Store) partialAutomationSessionUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Session) (err error) {
	for _, res := range rr {
		err = s.checkAutomationSessionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateAutomationSessions(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("atms.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalAutomationSessionEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertAutomationSession updates one or more existing rows in automation_sessions
func (s Store) UpsertAutomationSession(ctx context.Context, rr ...*types.Session) (err error) {
	for _, res := range rr {
		err = s.checkAutomationSessionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertAutomationSessions(ctx, s.internalAutomationSessionEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAutomationSession Deletes one or more rows from automation_sessions table
func (s Store) DeleteAutomationSession(ctx context.Context, rr ...*types.Session) (err error) {
	for _, res := range rr {

		err = s.execDeleteAutomationSessions(ctx, squirrel.Eq{
			s.preprocessColumn("atms.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAutomationSessionByID Deletes row from the automation_sessions table
func (s Store) DeleteAutomationSessionByID(ctx context.Context, ID uint64) error {
	return s.execDeleteAutomationSessions(ctx, squirrel.Eq{
		s.preprocessColumn("atms.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateAutomationSessions Deletes all rows from the automation_sessions table
func (s Store) TruncateAutomationSessions(ctx context.Context) error {
	return s.Truncate(ctx, s.automationSessionTable())
}

// execLookupAutomationSession prepares AutomationSession query and executes it,
// returning types.Session (or error)
func (s Store) execLookupAutomationSession(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Session, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.automationSessionsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalAutomationSessionRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateAutomationSessions updates all matched (by cnd) rows in automation_sessions with given data
func (s Store) execCreateAutomationSessions(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.automationSessionTable()).SetMap(payload))
}

// execUpdateAutomationSessions updates all matched (by cnd) rows in automation_sessions with given data
func (s Store) execUpdateAutomationSessions(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.automationSessionTable("atms")).Where(cnd).SetMap(set))
}

// execUpsertAutomationSessions inserts new or updates matching (by-primary-key) rows in automation_sessions with given data
func (s Store) execUpsertAutomationSessions(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.automationSessionTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteAutomationSessions Deletes all matched (by cnd) rows in automation_sessions with given data
func (s Store) execDeleteAutomationSessions(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.automationSessionTable("atms")).Where(cnd))
}

func (s Store) internalAutomationSessionRowScanner(row rowScanner) (res *types.Session, err error) {
	res = &types.Session{}

	if _, has := s.config.RowScanners["automationSession"]; has {
		scanner := s.config.RowScanners["automationSession"].(func(_ rowScanner, _ *types.Session) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.WorkflowID,
			&res.EventType,
			&res.ResourceType,
			&res.Status,
			&res.WallTime,
			&res.UserTime,
			&res.Input,
			&res.Output,
			&res.Trace,
			&res.CreatedBy,
			&res.CreatedAt,
			&res.PurgeAt,
			&res.CompletedAt,
			&res.SuspendedAt,
			&res.Error,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan automationSession db row").Wrap(err)
	} else {
		return res, nil
	}
}

// QueryAutomationSessions returns squirrel.SelectBuilder with set table and all columns
func (s Store) automationSessionsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.automationSessionTable("atms"), s.automationSessionColumns("atms")...)
}

// automationSessionTable name of the db table
func (Store) automationSessionTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "automation_sessions" + alias
}

// AutomationSessionColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) automationSessionColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_workflow",
		alias + "event_type",
		alias + "resource_type",
		alias + "status",
		alias + "wall_time",
		alias + "user_time",
		alias + "input",
		alias + "output",
		alias + "trace",
		alias + "created_by",
		alias + "created_at",
		alias + "purge_at",
		alias + "completed_at",
		alias + "suspended_at",
		alias + "error",
	}
}

// {true true false false false false}

// internalAutomationSessionEncoder encodes fields from types.Session to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeAutomationSession
// func when rdbms.customEncoder=true
func (s Store) internalAutomationSessionEncoder(res *types.Session) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"rel_workflow":  res.WorkflowID,
		"event_type":    res.EventType,
		"resource_type": res.ResourceType,
		"status":        res.Status,
		"wall_time":     res.WallTime,
		"user_time":     res.UserTime,
		"input":         res.Input,
		"output":        res.Output,
		"trace":         res.Trace,
		"created_by":    res.CreatedBy,
		"created_at":    res.CreatedAt,
		"purge_at":      res.PurgeAt,
		"completed_at":  res.CompletedAt,
		"suspended_at":  res.SuspendedAt,
		"error":         res.Error,
	}
}

// checkAutomationSessionConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkAutomationSessionConstraints(ctx context.Context, res *types.Session) error {
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
