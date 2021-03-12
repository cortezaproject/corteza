package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/actionlog.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchActionlogs returns all matching rows
//
// This function calls convertActionlogFilter with the given
// actionlog.Filter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchActionlogs(ctx context.Context, f actionlog.Filter) (actionlog.ActionSet, actionlog.Filter, error) {
	var (
		err error
		set []*actionlog.Action
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertActionlogFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryActionlogs(ctx, q, nil)
		return err
	}()
}

// QueryActionlogs queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryActionlogs(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*actionlog.Action) (bool, error),
) ([]*actionlog.Action, error) {
	var (
		set = make([]*actionlog.Action, 0, DefaultSliceCapacity)
		res *actionlog.Action

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalActionlogRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// CreateActionlog creates one or more rows in actionlog table
func (s Store) CreateActionlog(ctx context.Context, rr ...*actionlog.Action) (err error) {
	for _, res := range rr {
		err = s.checkActionlogConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateActionlogs(ctx, s.internalActionlogEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// TruncateActionlogs Deletes all rows from the actionlog table
func (s Store) TruncateActionlogs(ctx context.Context) error {
	return s.Truncate(ctx, s.actionlogTable())
}

// execLookupActionlog prepares Actionlog query and executes it,
// returning actionlog.Action (or error)
func (s Store) execLookupActionlog(ctx context.Context, cnd squirrel.Sqlizer) (res *actionlog.Action, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.actionlogsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalActionlogRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateActionlogs updates all matched (by cnd) rows in actionlog with given data
func (s Store) execCreateActionlogs(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.actionlogTable()).SetMap(payload))
}

func (s Store) internalActionlogRowScanner(row rowScanner) (res *actionlog.Action, err error) {
	res = &actionlog.Action{}

	if _, has := s.config.RowScanners["actionlog"]; has {
		scanner := s.config.RowScanners["actionlog"].(func(_ rowScanner, _ *actionlog.Action) error)
		err = scanner(row, res)
	} else {
		err = s.scanActionlogRow(row, res)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan actionlog db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryActionlogs returns squirrel.SelectBuilder with set table and all columns
func (s Store) actionlogsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.actionlogTable("alg"), s.actionlogColumns("alg")...)
}

// actionlogTable name of the db table
func (Store) actionlogTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "actionlog" + alias
}

// ActionlogColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) actionlogColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "ts",
		alias + "request_origin",
		alias + "request_id",
		alias + "actor_ip_addr",
		alias + "actor_id",
		alias + "resource",
		alias + "action",
		alias + "error",
		alias + "severity",
		alias + "description",
		alias + "meta",
	}
}

// {true true false false false false}

// internalActionlogEncoder encodes fields from actionlog.Action to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeActionlog
// func when rdbms.customEncoder=true
func (s Store) internalActionlogEncoder(res *actionlog.Action) store.Payload {
	return s.encodeActionlog(res)
}

// checkActionlogConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkActionlogConstraints(ctx context.Context, res *actionlog.Action) error {
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
