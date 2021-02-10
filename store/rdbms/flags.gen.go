package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/flags.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/flag/types"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchFlags returns all matching rows
//
// This function calls convertFlagFilter with the given
// types.FlagFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchFlags(ctx context.Context, f types.FlagFilter) (types.FlagSet, types.FlagFilter, error) {
	var (
		err error
		set []*types.Flag
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertFlagFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryFlags(ctx, q, nil)
		return err
	}()
}

// QueryFlags queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryFlags(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Flag) (bool, error),
) ([]*types.Flag, error) {
	var (
		set = make([]*types.Flag, 0, DefaultSliceCapacity)
		res *types.Flag

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalFlagRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupFlagByKindResourceIDName Flag lookup by kind, resource, name
func (s Store) LookupFlagByKindResourceIDName(ctx context.Context, kind string, resource_id uint64, name string) (*types.Flag, error) {
	return s.execLookupFlag(ctx, squirrel.Eq{
		s.preprocessColumn("flg.kind", ""):         store.PreprocessValue(kind, ""),
		s.preprocessColumn("flg.rel_resource", ""): store.PreprocessValue(resource_id, ""),
		s.preprocessColumn("flg.name", "lower"):    store.PreprocessValue(name, "lower"),
	})
}

// LookupFlagByKindResourceID Flag lookup by kind, resource
func (s Store) LookupFlagByKindResourceID(ctx context.Context, kind string, resource_id uint64) (*types.Flag, error) {
	return s.execLookupFlag(ctx, squirrel.Eq{
		s.preprocessColumn("flg.kind", ""):         store.PreprocessValue(kind, ""),
		s.preprocessColumn("flg.rel_resource", ""): store.PreprocessValue(resource_id, ""),
	})
}

// LookupFlagByKindResourceIDOwnedBy Flag lookup by kind, resource, owner
func (s Store) LookupFlagByKindResourceIDOwnedBy(ctx context.Context, kind string, resource_id uint64, owned_by uint64) (*types.Flag, error) {
	return s.execLookupFlag(ctx, squirrel.Eq{
		s.preprocessColumn("flg.kind", ""):         store.PreprocessValue(kind, ""),
		s.preprocessColumn("flg.rel_resource", ""): store.PreprocessValue(resource_id, ""),
		s.preprocessColumn("flg.owned_by", ""):     store.PreprocessValue(owned_by, ""),
	})
}

// LookupFlagByKindResourceIDOwnedByName Flag lookup by kind, resource, owner, name
func (s Store) LookupFlagByKindResourceIDOwnedByName(ctx context.Context, kind string, resource_id uint64, owned_by uint64, name string) (*types.Flag, error) {
	return s.execLookupFlag(ctx, squirrel.Eq{
		s.preprocessColumn("flg.kind", ""):         store.PreprocessValue(kind, ""),
		s.preprocessColumn("flg.rel_resource", ""): store.PreprocessValue(resource_id, ""),
		s.preprocessColumn("flg.owned_by", ""):     store.PreprocessValue(owned_by, ""),
		s.preprocessColumn("flg.name", "lower"):    store.PreprocessValue(name, "lower"),
	})
}

// CreateFlag creates one or more rows in flags table
func (s Store) CreateFlag(ctx context.Context, rr ...*types.Flag) (err error) {
	for _, res := range rr {
		err = s.checkFlagConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateFlags(ctx, s.internalFlagEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateFlag updates one or more existing rows in flags
func (s Store) UpdateFlag(ctx context.Context, rr ...*types.Flag) error {
	return s.partialFlagUpdate(ctx, nil, rr...)
}

// partialFlagUpdate updates one or more existing rows in flags
func (s Store) partialFlagUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Flag) (err error) {
	for _, res := range rr {
		err = s.checkFlagConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateFlags(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("flg.kind", ""): store.PreprocessValue(res.Kind, ""), s.preprocessColumn("flg.rel_resource", ""): store.PreprocessValue(res.ResourceID, ""), s.preprocessColumn("flg.owned_by", ""): store.PreprocessValue(res.OwnedBy, ""), s.preprocessColumn("flg.name", "lower"): store.PreprocessValue(res.Name, "lower"),
			},
			s.internalFlagEncoder(res).Skip("kind", "rel_resource", "owned_by", "name").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertFlag updates one or more existing rows in flags
func (s Store) UpsertFlag(ctx context.Context, rr ...*types.Flag) (err error) {
	for _, res := range rr {
		err = s.checkFlagConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertFlags(ctx, s.internalFlagEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteFlag Deletes one or more rows from flags table
func (s Store) DeleteFlag(ctx context.Context, rr ...*types.Flag) (err error) {
	for _, res := range rr {

		err = s.execDeleteFlags(ctx, squirrel.Eq{
			s.preprocessColumn("flg.kind", ""): store.PreprocessValue(res.Kind, ""), s.preprocessColumn("flg.rel_resource", ""): store.PreprocessValue(res.ResourceID, ""), s.preprocessColumn("flg.owned_by", ""): store.PreprocessValue(res.OwnedBy, ""), s.preprocessColumn("flg.name", "lower"): store.PreprocessValue(res.Name, "lower"),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteFlagByKindResourceIDOwnedByName Deletes row from the flags table
func (s Store) DeleteFlagByKindResourceIDOwnedByName(ctx context.Context, kind string, resourceID uint64, ownedBy uint64, name string) error {
	return s.execDeleteFlags(ctx, squirrel.Eq{
		s.preprocessColumn("flg.kind", ""):         store.PreprocessValue(kind, ""),
		s.preprocessColumn("flg.rel_resource", ""): store.PreprocessValue(resourceID, ""),
		s.preprocessColumn("flg.owned_by", ""):     store.PreprocessValue(ownedBy, ""),
		s.preprocessColumn("flg.name", "lower"):    store.PreprocessValue(name, "lower"),
	})
}

// TruncateFlags Deletes all rows from the flags table
func (s Store) TruncateFlags(ctx context.Context) error {
	return s.Truncate(ctx, s.flagTable())
}

// execLookupFlag prepares Flag query and executes it,
// returning types.Flag (or error)
func (s Store) execLookupFlag(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Flag, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.flagsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalFlagRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateFlags updates all matched (by cnd) rows in flags with given data
func (s Store) execCreateFlags(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.flagTable()).SetMap(payload))
}

// execUpdateFlags updates all matched (by cnd) rows in flags with given data
func (s Store) execUpdateFlags(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.flagTable("flg")).Where(cnd).SetMap(set))
}

// execUpsertFlags inserts new or updates matching (by-primary-key) rows in flags with given data
func (s Store) execUpsertFlags(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.flagTable(),
		set,
		s.preprocessColumn("kind", ""),
		s.preprocessColumn("rel_resource", ""),
		s.preprocessColumn("owned_by", ""),
		s.preprocessColumn("name", "lower"),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteFlags Deletes all matched (by cnd) rows in flags with given data
func (s Store) execDeleteFlags(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.flagTable("flg")).Where(cnd))
}

func (s Store) internalFlagRowScanner(row rowScanner) (res *types.Flag, err error) {
	res = &types.Flag{}

	if _, has := s.config.RowScanners["flag"]; has {
		scanner := s.config.RowScanners["flag"].(func(_ rowScanner, _ *types.Flag) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.Kind,
			&res.ResourceID,
			&res.OwnedBy,
			&res.Name,
			&res.Active,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan flag db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryFlags returns squirrel.SelectBuilder with set table and all columns
func (s Store) flagsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.flagTable("flg"), s.flagColumns("flg")...)
}

// flagTable name of the db table
func (Store) flagTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "flags" + alias
}

// FlagColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) flagColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "kind",
		alias + "rel_resource",
		alias + "owned_by",
		alias + "name",
		alias + "active",
	}
}

// {true true false false false false}

// internalFlagEncoder encodes fields from types.Flag to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeFlag
// func when rdbms.customEncoder=true
func (s Store) internalFlagEncoder(res *types.Flag) store.Payload {
	return store.Payload{
		"kind":         res.Kind,
		"rel_resource": res.ResourceID,
		"owned_by":     res.OwnedBy,
		"name":         res.Name,
		"active":       res.Active,
	}
}

// checkFlagConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkFlagConstraints(ctx context.Context, res *types.Flag) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	valid = valid && len(res.Kind) > 0

	valid = valid && res.ResourceID > 0

	valid = valid && len(res.Name) > 0

	valid = valid && len(res.Kind) > 0

	valid = valid && res.ResourceID > 0

	valid = valid && len(res.Kind) > 0

	valid = valid && res.ResourceID > 0

	valid = valid && res.OwnedBy > 0

	valid = valid && len(res.Kind) > 0

	valid = valid && res.ResourceID > 0

	valid = valid && res.OwnedBy > 0

	valid = valid && len(res.Name) > 0

	if !valid {
		return nil
	}

	{
		ex, err := s.LookupFlagByKindResourceIDName(ctx, res.Kind, res.ResourceID, res.Name)
		if err == nil && ex != nil && ex.Kind != res.Kind && ex.ResourceID != res.ResourceID && ex.OwnedBy != res.OwnedBy && ex.Name != res.Name {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	{
		ex, err := s.LookupFlagByKindResourceID(ctx, res.Kind, res.ResourceID)
		if err == nil && ex != nil && ex.Kind != res.Kind && ex.ResourceID != res.ResourceID && ex.OwnedBy != res.OwnedBy && ex.Name != res.Name {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	{
		ex, err := s.LookupFlagByKindResourceIDOwnedBy(ctx, res.Kind, res.ResourceID, res.OwnedBy)
		if err == nil && ex != nil && ex.Kind != res.Kind && ex.ResourceID != res.ResourceID && ex.OwnedBy != res.OwnedBy && ex.Name != res.Name {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	{
		ex, err := s.LookupFlagByKindResourceIDOwnedByName(ctx, res.Kind, res.ResourceID, res.OwnedBy, res.Name)
		if err == nil && ex != nil && ex.Kind != res.Kind && ex.ResourceID != res.ResourceID && ex.OwnedBy != res.OwnedBy && ex.Name != res.Name {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
