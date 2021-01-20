package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/labels.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/label/types"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchLabels returns all matching rows
//
// This function calls convertLabelFilter with the given
// types.LabelFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchLabels(ctx context.Context, f types.LabelFilter) (types.LabelSet, types.LabelFilter, error) {
	var (
		err error
		set []*types.Label
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertLabelFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryLabels(ctx, q, nil)
		return err
	}()
}

// QueryLabels queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryLabels(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Label) (bool, error),
) ([]*types.Label, error) {
	var (
		set = make([]*types.Label, 0, DefaultSliceCapacity)
		res *types.Label

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalLabelRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupLabelByKindResourceIDName Label lookup by kind, resource, name
func (s Store) LookupLabelByKindResourceIDName(ctx context.Context, kind string, resource_id uint64, name string) (*types.Label, error) {
	return s.execLookupLabel(ctx, squirrel.Eq{
		s.preprocessColumn("lbl.kind", ""):         store.PreprocessValue(kind, ""),
		s.preprocessColumn("lbl.rel_resource", ""): store.PreprocessValue(resource_id, ""),
		s.preprocessColumn("lbl.name", "lower"):    store.PreprocessValue(name, "lower"),
	})
}

// CreateLabel creates one or more rows in labels table
func (s Store) CreateLabel(ctx context.Context, rr ...*types.Label) (err error) {
	for _, res := range rr {
		err = s.checkLabelConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateLabels(ctx, s.internalLabelEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateLabel updates one or more existing rows in labels
func (s Store) UpdateLabel(ctx context.Context, rr ...*types.Label) error {
	return s.partialLabelUpdate(ctx, nil, rr...)
}

// partialLabelUpdate updates one or more existing rows in labels
func (s Store) partialLabelUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Label) (err error) {
	for _, res := range rr {
		err = s.checkLabelConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateLabels(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("lbl.kind", ""): store.PreprocessValue(res.Kind, ""), s.preprocessColumn("lbl.rel_resource", ""): store.PreprocessValue(res.ResourceID, ""), s.preprocessColumn("lbl.name", "lower"): store.PreprocessValue(res.Name, "lower"),
			},
			s.internalLabelEncoder(res).Skip("kind", "rel_resource", "name").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertLabel updates one or more existing rows in labels
func (s Store) UpsertLabel(ctx context.Context, rr ...*types.Label) (err error) {
	for _, res := range rr {
		err = s.checkLabelConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertLabels(ctx, s.internalLabelEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteLabel Deletes one or more rows from labels table
func (s Store) DeleteLabel(ctx context.Context, rr ...*types.Label) (err error) {
	for _, res := range rr {

		err = s.execDeleteLabels(ctx, squirrel.Eq{
			s.preprocessColumn("lbl.kind", ""): store.PreprocessValue(res.Kind, ""), s.preprocessColumn("lbl.rel_resource", ""): store.PreprocessValue(res.ResourceID, ""), s.preprocessColumn("lbl.name", "lower"): store.PreprocessValue(res.Name, "lower"),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteLabelByKindResourceIDName Deletes row from the labels table
func (s Store) DeleteLabelByKindResourceIDName(ctx context.Context, kind string, resourceID uint64, name string) error {
	return s.execDeleteLabels(ctx, squirrel.Eq{
		s.preprocessColumn("lbl.kind", ""):         store.PreprocessValue(kind, ""),
		s.preprocessColumn("lbl.rel_resource", ""): store.PreprocessValue(resourceID, ""),
		s.preprocessColumn("lbl.name", "lower"):    store.PreprocessValue(name, "lower"),
	})
}

// TruncateLabels Deletes all rows from the labels table
func (s Store) TruncateLabels(ctx context.Context) error {
	return s.Truncate(ctx, s.labelTable())
}

// execLookupLabel prepares Label query and executes it,
// returning types.Label (or error)
func (s Store) execLookupLabel(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Label, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.labelsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalLabelRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateLabels updates all matched (by cnd) rows in labels with given data
func (s Store) execCreateLabels(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.labelTable()).SetMap(payload))
}

// execUpdateLabels updates all matched (by cnd) rows in labels with given data
func (s Store) execUpdateLabels(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.labelTable("lbl")).Where(cnd).SetMap(set))
}

// execUpsertLabels inserts new or updates matching (by-primary-key) rows in labels with given data
func (s Store) execUpsertLabels(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.labelTable(),
		set,
		s.preprocessColumn("kind", ""),
		s.preprocessColumn("rel_resource", ""),
		s.preprocessColumn("name", "lower"),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteLabels Deletes all matched (by cnd) rows in labels with given data
func (s Store) execDeleteLabels(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.labelTable("lbl")).Where(cnd))
}

func (s Store) internalLabelRowScanner(row rowScanner) (res *types.Label, err error) {
	res = &types.Label{}

	if _, has := s.config.RowScanners["label"]; has {
		scanner := s.config.RowScanners["label"].(func(_ rowScanner, _ *types.Label) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.Kind,
			&res.ResourceID,
			&res.Name,
			&res.Value,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan label db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryLabels returns squirrel.SelectBuilder with set table and all columns
func (s Store) labelsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.labelTable("lbl"), s.labelColumns("lbl")...)
}

// labelTable name of the db table
func (Store) labelTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "labels" + alias
}

// LabelColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) labelColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "kind",
		alias + "rel_resource",
		alias + "name",
		alias + "value",
	}
}

// {true true false false false false}

// internalLabelEncoder encodes fields from types.Label to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeLabel
// func when rdbms.customEncoder=true
func (s Store) internalLabelEncoder(res *types.Label) store.Payload {
	return store.Payload{
		"kind":         res.Kind,
		"rel_resource": res.ResourceID,
		"name":         res.Name,
		"value":        res.Value,
	}
}

// checkLabelConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkLabelConstraints(ctx context.Context, res *types.Label) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	valid = valid && len(res.Kind) > 0

	valid = valid && res.ResourceID > 0

	valid = valid && len(res.Name) > 0

	if !valid {
		return nil
	}

	{
		ex, err := s.LookupLabelByKindResourceIDName(ctx, res.Kind, res.ResourceID, res.Name)
		if err == nil && ex != nil && ex.Kind != res.Kind && ex.ResourceID != res.ResourceID && ex.Name != res.Name {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
