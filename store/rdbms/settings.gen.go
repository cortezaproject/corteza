package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/settings.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = errors.Is

// SearchSettings returns all matching rows
//
// This function calls convertSettingFilter with the given
// types.SettingsFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchSettings(ctx context.Context, f types.SettingsFilter) (types.SettingValueSet, types.SettingsFilter, error) {
	var (
		err error
		set []*types.SettingValue
		q   squirrel.SelectBuilder
	)
	q, err = s.convertSettingFilter(f)
	if err != nil {
		return nil, f, err
	}

	return set, f, s.config.ErrorHandler(func() error {
		set, _, _, err = s.QuerySettings(ctx, q, f.Check)
		return err

	}())
}

// QuerySettings queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QuerySettings(
	ctx context.Context,
	q squirrel.SelectBuilder,
	check func(*types.SettingValue) (bool, error),
) ([]*types.SettingValue, uint, *types.SettingValue, error) {
	var (
		set = make([]*types.SettingValue, 0, DefaultSliceCapacity)
		res *types.SettingValue

		// Query rows with
		rows, err = s.Query(ctx, q)

		fetched uint
	)

	if err != nil {
		return nil, 0, nil, err
	}

	defer rows.Close()
	for rows.Next() {
		fetched++
		if err = rows.Err(); err == nil {
			res, err = s.internalSettingRowScanner(rows)
		}

		if err != nil {
			return nil, 0, nil, err
		}

		// If check function is set, call it and act accordingly
		if check != nil {
			if chk, err := check(res); err != nil {
				return nil, 0, nil, err
			} else if !chk {
				// did not pass the check
				// go with the next row
				continue
			}
		}

		set = append(set, res)
	}

	return set, fetched, res, rows.Err()
}

// LookupSettingByNameOwnedBy searches for settings by name and owner
func (s Store) LookupSettingByNameOwnedBy(ctx context.Context, name string, owned_by uint64) (*types.SettingValue, error) {
	return s.execLookupSetting(ctx, squirrel.Eq{
		s.preprocessColumn("st.name", ""):      s.preprocessValue(name, ""),
		s.preprocessColumn("st.rel_owner", ""): s.preprocessValue(owned_by, ""),
	})
}

// CreateSetting creates one or more rows in settings table
func (s Store) CreateSetting(ctx context.Context, rr ...*types.SettingValue) (err error) {
	for _, res := range rr {
		err = s.checkSettingConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateSettings(ctx, s.internalSettingEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateSetting updates one or more existing rows in settings
func (s Store) UpdateSetting(ctx context.Context, rr ...*types.SettingValue) error {
	return s.config.ErrorHandler(s.PartialSettingUpdate(ctx, nil, rr...))
}

// PartialSettingUpdate updates one or more existing rows in settings
func (s Store) PartialSettingUpdate(ctx context.Context, onlyColumns []string, rr ...*types.SettingValue) (err error) {
	for _, res := range rr {
		err = s.checkSettingConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateSettings(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("st.name", ""): s.preprocessValue(res.Name, ""), s.preprocessColumn("st.rel_owner", ""): s.preprocessValue(res.OwnedBy, ""),
			},
			s.internalSettingEncoder(res).Skip("name", "rel_owner").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// UpsertSetting updates one or more existing rows in settings
func (s Store) UpsertSetting(ctx context.Context, rr ...*types.SettingValue) (err error) {
	for _, res := range rr {
		err = s.checkSettingConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.config.ErrorHandler(s.execUpsertSettings(ctx, s.internalSettingEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteSetting Deletes one or more rows from settings table
func (s Store) DeleteSetting(ctx context.Context, rr ...*types.SettingValue) (err error) {
	for _, res := range rr {

		err = s.execDeleteSettings(ctx, squirrel.Eq{
			s.preprocessColumn("st.name", ""): s.preprocessValue(res.Name, ""), s.preprocessColumn("st.rel_owner", ""): s.preprocessValue(res.OwnedBy, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteSettingByNameOwnedBy Deletes row from the settings table
func (s Store) DeleteSettingByNameOwnedBy(ctx context.Context, name string, ownedBy uint64) error {
	return s.execDeleteSettings(ctx, squirrel.Eq{
		s.preprocessColumn("st.name", ""):      s.preprocessValue(name, ""),
		s.preprocessColumn("st.rel_owner", ""): s.preprocessValue(ownedBy, ""),
	})
}

// TruncateSettings Deletes all rows from the settings table
func (s Store) TruncateSettings(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.settingTable()))
}

// execLookupSetting prepares Setting query and executes it,
// returning types.SettingValue (or error)
func (s Store) execLookupSetting(ctx context.Context, cnd squirrel.Sqlizer) (res *types.SettingValue, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.settingsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalSettingRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateSettings updates all matched (by cnd) rows in settings with given data
func (s Store) execCreateSettings(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.settingTable()).SetMap(payload)))
}

// execUpdateSettings updates all matched (by cnd) rows in settings with given data
func (s Store) execUpdateSettings(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.settingTable("st")).Where(cnd).SetMap(set)))
}

// execUpsertSettings inserts new or updates matching (by-primary-key) rows in settings with given data
func (s Store) execUpsertSettings(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.settingTable(),
		set,
		"name",
		"rel_owner",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteSettings Deletes all matched (by cnd) rows in settings with given data
func (s Store) execDeleteSettings(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.settingTable("st")).Where(cnd)))
}

func (s Store) internalSettingRowScanner(row rowScanner) (res *types.SettingValue, err error) {
	res = &types.SettingValue{}

	if _, has := s.config.RowScanners["setting"]; has {
		scanner := s.config.RowScanners["setting"].(func(_ rowScanner, _ *types.SettingValue) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.Name,
			&res.OwnedBy,
			&res.Value,
			&res.UpdatedBy,
			&res.UpdatedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for Setting: %w", err)
	} else {
		return res, nil
	}
}

// QuerySettings returns squirrel.SelectBuilder with set table and all columns
func (s Store) settingsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.settingTable("st"), s.settingColumns("st")...)
}

// settingTable name of the db table
func (Store) settingTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "settings" + alias
}

// SettingColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) settingColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "name",
		alias + "rel_owner",
		alias + "value",
		alias + "updated_by",
		alias + "updated_at",
	}
}

// {true true false false true}

// internalSettingEncoder encodes fields from types.SettingValue to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeSetting
// func when rdbms.customEncoder=true
func (s Store) internalSettingEncoder(res *types.SettingValue) store.Payload {
	return store.Payload{
		"name":       res.Name,
		"rel_owner":  res.OwnedBy,
		"value":      res.Value,
		"updated_by": res.UpdatedBy,
		"updated_at": res.UpdatedAt,
	}
}

func (s *Store) checkSettingConstraints(ctx context.Context, res *types.SettingValue) error {

	return nil
}
