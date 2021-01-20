package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_record_values.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// searchComposeRecordValues returns all matching rows
//
// This function calls convertComposeRecordValueFilter with the given
// types.RecordValueFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) searchComposeRecordValues(ctx context.Context, _mod *types.Module, f types.RecordValueFilter) (types.RecordValueSet, types.RecordValueFilter, error) {
	var (
		err error
		set []*types.RecordValue
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertComposeRecordValueFilter(_mod, f)
		if err != nil {
			return err
		}

		set, err = s.QueryComposeRecordValues(ctx, _mod, q, nil)
		return err
	}()
}

// QueryComposeRecordValues queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryComposeRecordValues(
	ctx context.Context, _mod *types.Module,
	q squirrel.Sqlizer,
	check func(*types.RecordValue) (bool, error),
) ([]*types.RecordValue, error) {
	var (
		set = make([]*types.RecordValue, 0, DefaultSliceCapacity)
		res *types.RecordValue

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalComposeRecordValueRowScanner(_mod, rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// createComposeRecordValue creates one or more rows in compose_record_value table
func (s Store) createComposeRecordValue(ctx context.Context, _mod *types.Module, rr ...*types.RecordValue) (err error) {
	for _, res := range rr {
		err = s.checkComposeRecordValueConstraints(ctx, _mod, res)
		if err != nil {
			return err
		}

		err = s.execCreateComposeRecordValues(ctx, s.internalComposeRecordValueEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// updateComposeRecordValue updates one or more existing rows in compose_record_value
func (s Store) updateComposeRecordValue(ctx context.Context, _mod *types.Module, rr ...*types.RecordValue) error {
	return s.partialComposeRecordValueUpdate(ctx, _mod, nil, rr...)
}

// partialComposeRecordValueUpdate updates one or more existing rows in compose_record_value
func (s Store) partialComposeRecordValueUpdate(ctx context.Context, _mod *types.Module, onlyColumns []string, rr ...*types.RecordValue) (err error) {
	for _, res := range rr {
		err = s.checkComposeRecordValueConstraints(ctx, _mod, res)
		if err != nil {
			return err
		}

		err = s.execUpdateComposeRecordValues(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("crv.record_id", ""): store.PreprocessValue(res.RecordID, ""), s.preprocessColumn("crv.name", ""): store.PreprocessValue(res.Name, ""), s.preprocessColumn("crv.place", ""): store.PreprocessValue(res.Place, ""),
			},
			s.internalComposeRecordValueEncoder(res).Skip("record_id", "name", "place").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// upsertComposeRecordValue updates one or more existing rows in compose_record_value
func (s Store) upsertComposeRecordValue(ctx context.Context, _mod *types.Module, rr ...*types.RecordValue) (err error) {
	for _, res := range rr {
		err = s.checkComposeRecordValueConstraints(ctx, _mod, res)
		if err != nil {
			return err
		}

		err = s.execUpsertComposeRecordValues(ctx, s.internalComposeRecordValueEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// deleteComposeRecordValue Deletes one or more rows from compose_record_value table
func (s Store) deleteComposeRecordValue(ctx context.Context, _mod *types.Module, rr ...*types.RecordValue) (err error) {
	for _, res := range rr {

		err = s.execDeleteComposeRecordValues(ctx, squirrel.Eq{
			s.preprocessColumn("crv.record_id", ""): store.PreprocessValue(res.RecordID, ""), s.preprocessColumn("crv.name", ""): store.PreprocessValue(res.Name, ""), s.preprocessColumn("crv.place", ""): store.PreprocessValue(res.Place, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// deleteComposeRecordValueByRecordIDNamePlace Deletes row from the compose_record_value table
func (s Store) deleteComposeRecordValueByRecordIDNamePlace(ctx context.Context, _mod *types.Module, recordID uint64, name string, place uint) error {
	return s.execDeleteComposeRecordValues(ctx, squirrel.Eq{
		s.preprocessColumn("crv.record_id", ""): store.PreprocessValue(recordID, ""),
		s.preprocessColumn("crv.name", ""):      store.PreprocessValue(name, ""),
		s.preprocessColumn("crv.place", ""):     store.PreprocessValue(place, ""),
	})
}

// truncateComposeRecordValues Deletes all rows from the compose_record_value table
func (s Store) truncateComposeRecordValues(ctx context.Context, _mod *types.Module) error {
	return s.Truncate(ctx, s.composeRecordValueTable())
}

// execLookupComposeRecordValue prepares ComposeRecordValue query and executes it,
// returning types.RecordValue (or error)
func (s Store) execLookupComposeRecordValue(ctx context.Context, _mod *types.Module, cnd squirrel.Sqlizer) (res *types.RecordValue, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.composeRecordValuesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalComposeRecordValueRowScanner(_mod, row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateComposeRecordValues updates all matched (by cnd) rows in compose_record_value with given data
func (s Store) execCreateComposeRecordValues(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.composeRecordValueTable()).SetMap(payload))
}

// execUpdateComposeRecordValues updates all matched (by cnd) rows in compose_record_value with given data
func (s Store) execUpdateComposeRecordValues(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.composeRecordValueTable("crv")).Where(cnd).SetMap(set))
}

// execUpsertComposeRecordValues inserts new or updates matching (by-primary-key) rows in compose_record_value with given data
func (s Store) execUpsertComposeRecordValues(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.composeRecordValueTable(),
		set,
		s.preprocessColumn("record_id", ""),
		s.preprocessColumn("name", ""),
		s.preprocessColumn("place", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteComposeRecordValues Deletes all matched (by cnd) rows in compose_record_value with given data
func (s Store) execDeleteComposeRecordValues(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.composeRecordValueTable("crv")).Where(cnd))
}

func (s Store) internalComposeRecordValueRowScanner(_mod *types.Module, row rowScanner) (res *types.RecordValue, err error) {
	res = &types.RecordValue{}

	if _, has := s.config.RowScanners["composeRecordValue"]; has {
		scanner := s.config.RowScanners["composeRecordValue"].(func(_mod *types.Module, _ rowScanner, _ *types.RecordValue) error)
		err = scanner(_mod, row, res)
	} else {
		err = row.Scan(
			&res.RecordID,
			&res.Name,
			&res.Place,
			&res.Value,
			&res.Ref,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan composeRecordValue db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryComposeRecordValues returns squirrel.SelectBuilder with set table and all columns
func (s Store) composeRecordValuesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.composeRecordValueTable("crv"), s.composeRecordValueColumns("crv")...)
}

// composeRecordValueTable name of the db table
func (Store) composeRecordValueTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_record_value" + alias
}

// ComposeRecordValueColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) composeRecordValueColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "record_id",
		alias + "name",
		alias + "place",
		alias + "value",
		alias + "ref",
		alias + "deleted_at",
	}
}

// {true false false false false false}

// internalComposeRecordValueEncoder encodes fields from types.RecordValue to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposeRecordValue
// func when rdbms.customEncoder=true
func (s Store) internalComposeRecordValueEncoder(res *types.RecordValue) store.Payload {
	return store.Payload{
		"record_id":  res.RecordID,
		"name":       res.Name,
		"place":      res.Place,
		"value":      res.Value,
		"ref":        res.Ref,
		"deleted_at": res.DeletedAt,
	}
}

// checkComposeRecordValueConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkComposeRecordValueConstraints(ctx context.Context, _mod *types.Module, res *types.RecordValue) error {
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
