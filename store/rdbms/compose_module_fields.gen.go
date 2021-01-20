package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_module_fields.yaml
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

// SearchComposeModuleFields returns all matching rows
//
// This function calls convertComposeModuleFieldFilter with the given
// types.ModuleFieldFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchComposeModuleFields(ctx context.Context, f types.ModuleFieldFilter) (types.ModuleFieldSet, types.ModuleFieldFilter, error) {
	var (
		err error
		set []*types.ModuleField
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertComposeModuleFieldFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryComposeModuleFields(ctx, q, nil)
		return err
	}()
}

// QueryComposeModuleFields queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryComposeModuleFields(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.ModuleField) (bool, error),
) ([]*types.ModuleField, error) {
	var (
		set = make([]*types.ModuleField, 0, DefaultSliceCapacity)
		res *types.ModuleField

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalComposeModuleFieldRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupComposeModuleFieldByModuleIDName searches for compose module field by name (case-insensitive)
func (s Store) LookupComposeModuleFieldByModuleIDName(ctx context.Context, module_id uint64, name string) (*types.ModuleField, error) {
	return s.execLookupComposeModuleField(ctx, squirrel.Eq{
		s.preprocessColumn("cmf.rel_module", ""): store.PreprocessValue(module_id, ""),
		s.preprocessColumn("cmf.name", "lower"):  store.PreprocessValue(name, "lower"),

		"cmf.deleted_at": nil,
	})
}

// CreateComposeModuleField creates one or more rows in compose_module_field table
func (s Store) CreateComposeModuleField(ctx context.Context, rr ...*types.ModuleField) (err error) {
	for _, res := range rr {
		err = s.checkComposeModuleFieldConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateComposeModuleFields(ctx, s.internalComposeModuleFieldEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateComposeModuleField updates one or more existing rows in compose_module_field
func (s Store) UpdateComposeModuleField(ctx context.Context, rr ...*types.ModuleField) error {
	return s.partialComposeModuleFieldUpdate(ctx, nil, rr...)
}

// partialComposeModuleFieldUpdate updates one or more existing rows in compose_module_field
func (s Store) partialComposeModuleFieldUpdate(ctx context.Context, onlyColumns []string, rr ...*types.ModuleField) (err error) {
	for _, res := range rr {
		err = s.checkComposeModuleFieldConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateComposeModuleFields(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("cmf.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalComposeModuleFieldEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertComposeModuleField updates one or more existing rows in compose_module_field
func (s Store) UpsertComposeModuleField(ctx context.Context, rr ...*types.ModuleField) (err error) {
	for _, res := range rr {
		err = s.checkComposeModuleFieldConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertComposeModuleFields(ctx, s.internalComposeModuleFieldEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteComposeModuleField Deletes one or more rows from compose_module_field table
func (s Store) DeleteComposeModuleField(ctx context.Context, rr ...*types.ModuleField) (err error) {
	for _, res := range rr {

		err = s.execDeleteComposeModuleFields(ctx, squirrel.Eq{
			s.preprocessColumn("cmf.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteComposeModuleFieldByID Deletes row from the compose_module_field table
func (s Store) DeleteComposeModuleFieldByID(ctx context.Context, ID uint64) error {
	return s.execDeleteComposeModuleFields(ctx, squirrel.Eq{
		s.preprocessColumn("cmf.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateComposeModuleFields Deletes all rows from the compose_module_field table
func (s Store) TruncateComposeModuleFields(ctx context.Context) error {
	return s.Truncate(ctx, s.composeModuleFieldTable())
}

// execLookupComposeModuleField prepares ComposeModuleField query and executes it,
// returning types.ModuleField (or error)
func (s Store) execLookupComposeModuleField(ctx context.Context, cnd squirrel.Sqlizer) (res *types.ModuleField, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.composeModuleFieldsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalComposeModuleFieldRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateComposeModuleFields updates all matched (by cnd) rows in compose_module_field with given data
func (s Store) execCreateComposeModuleFields(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.composeModuleFieldTable()).SetMap(payload))
}

// execUpdateComposeModuleFields updates all matched (by cnd) rows in compose_module_field with given data
func (s Store) execUpdateComposeModuleFields(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.composeModuleFieldTable("cmf")).Where(cnd).SetMap(set))
}

// execUpsertComposeModuleFields inserts new or updates matching (by-primary-key) rows in compose_module_field with given data
func (s Store) execUpsertComposeModuleFields(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.composeModuleFieldTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteComposeModuleFields Deletes all matched (by cnd) rows in compose_module_field with given data
func (s Store) execDeleteComposeModuleFields(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.composeModuleFieldTable("cmf")).Where(cnd))
}

func (s Store) internalComposeModuleFieldRowScanner(row rowScanner) (res *types.ModuleField, err error) {
	res = &types.ModuleField{}

	if _, has := s.config.RowScanners["composeModuleField"]; has {
		scanner := s.config.RowScanners["composeModuleField"].(func(_ rowScanner, _ *types.ModuleField) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Name,
			&res.ModuleID,
			&res.Place,
			&res.Kind,
			&res.Label,
			&res.Options,
			&res.Private,
			&res.Required,
			&res.Visible,
			&res.Multi,
			&res.DefaultValue,
			&res.Expressions,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan composeModuleField db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryComposeModuleFields returns squirrel.SelectBuilder with set table and all columns
func (s Store) composeModuleFieldsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.composeModuleFieldTable("cmf"), s.composeModuleFieldColumns("cmf")...)
}

// composeModuleFieldTable name of the db table
func (Store) composeModuleFieldTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_module_field" + alias
}

// ComposeModuleFieldColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) composeModuleFieldColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "name",
		alias + "rel_module",
		alias + "place",
		alias + "kind",
		alias + "label",
		alias + "options",
		alias + "is_private",
		alias + "is_required",
		alias + "is_visible",
		alias + "is_multi",
		alias + "default_value",
		alias + "expressions",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false false false false}

// internalComposeModuleFieldEncoder encodes fields from types.ModuleField to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposeModuleField
// func when rdbms.customEncoder=true
func (s Store) internalComposeModuleFieldEncoder(res *types.ModuleField) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"name":          res.Name,
		"rel_module":    res.ModuleID,
		"place":         res.Place,
		"kind":          res.Kind,
		"label":         res.Label,
		"options":       res.Options,
		"is_private":    res.Private,
		"is_required":   res.Required,
		"is_visible":    res.Visible,
		"is_multi":      res.Multi,
		"default_value": res.DefaultValue,
		"expressions":   res.Expressions,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}

// checkComposeModuleFieldConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkComposeModuleFieldConstraints(ctx context.Context, res *types.ModuleField) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	valid = valid && res.ModuleID > 0

	valid = valid && len(res.Name) > 0

	if !valid {
		return nil
	}

	{
		ex, err := s.LookupComposeModuleFieldByModuleIDName(ctx, res.ModuleID, res.Name)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
