package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/role_members.yaml
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

// SearchRoleMembers returns all matching rows
//
// This function calls convertRoleMemberFilter with the given
// types.RoleMemberFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchRoleMembers(ctx context.Context, f types.RoleMemberFilter) (types.RoleMemberSet, types.RoleMemberFilter, error) {
	var (
		err error
		set []*types.RoleMember
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertRoleMemberFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryRoleMembers(ctx, q, nil)
		return err
	}()
}

// QueryRoleMembers queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryRoleMembers(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.RoleMember) (bool, error),
) ([]*types.RoleMember, error) {
	var (
		set = make([]*types.RoleMember, 0, DefaultSliceCapacity)
		res *types.RoleMember

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalRoleMemberRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// CreateRoleMember creates one or more rows in role_members table
func (s Store) CreateRoleMember(ctx context.Context, rr ...*types.RoleMember) (err error) {
	for _, res := range rr {
		err = s.checkRoleMemberConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateRoleMembers(ctx, s.internalRoleMemberEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateRoleMember updates one or more existing rows in role_members
func (s Store) UpdateRoleMember(ctx context.Context, rr ...*types.RoleMember) error {
	return s.partialRoleMemberUpdate(ctx, nil, rr...)
}

// partialRoleMemberUpdate updates one or more existing rows in role_members
func (s Store) partialRoleMemberUpdate(ctx context.Context, onlyColumns []string, rr ...*types.RoleMember) (err error) {
	for _, res := range rr {
		err = s.checkRoleMemberConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateRoleMembers(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("rm.rel_user", ""): store.PreprocessValue(res.UserID, ""), s.preprocessColumn("rm.rel_role", ""): store.PreprocessValue(res.RoleID, ""),
			},
			s.internalRoleMemberEncoder(res).Skip("rel_user", "rel_role").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertRoleMember updates one or more existing rows in role_members
func (s Store) UpsertRoleMember(ctx context.Context, rr ...*types.RoleMember) (err error) {
	for _, res := range rr {
		err = s.checkRoleMemberConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertRoleMembers(ctx, s.internalRoleMemberEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteRoleMember Deletes one or more rows from role_members table
func (s Store) DeleteRoleMember(ctx context.Context, rr ...*types.RoleMember) (err error) {
	for _, res := range rr {

		err = s.execDeleteRoleMembers(ctx, squirrel.Eq{
			s.preprocessColumn("rm.rel_user", ""): store.PreprocessValue(res.UserID, ""), s.preprocessColumn("rm.rel_role", ""): store.PreprocessValue(res.RoleID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteRoleMemberByUserIDRoleID Deletes row from the role_members table
func (s Store) DeleteRoleMemberByUserIDRoleID(ctx context.Context, userID uint64, roleID uint64) error {
	return s.execDeleteRoleMembers(ctx, squirrel.Eq{
		s.preprocessColumn("rm.rel_user", ""): store.PreprocessValue(userID, ""),
		s.preprocessColumn("rm.rel_role", ""): store.PreprocessValue(roleID, ""),
	})
}

// TruncateRoleMembers Deletes all rows from the role_members table
func (s Store) TruncateRoleMembers(ctx context.Context) error {
	return s.Truncate(ctx, s.roleMemberTable())
}

// execLookupRoleMember prepares RoleMember query and executes it,
// returning types.RoleMember (or error)
func (s Store) execLookupRoleMember(ctx context.Context, cnd squirrel.Sqlizer) (res *types.RoleMember, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.roleMembersSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalRoleMemberRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateRoleMembers updates all matched (by cnd) rows in role_members with given data
func (s Store) execCreateRoleMembers(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.roleMemberTable()).SetMap(payload))
}

// execUpdateRoleMembers updates all matched (by cnd) rows in role_members with given data
func (s Store) execUpdateRoleMembers(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.roleMemberTable("rm")).Where(cnd).SetMap(set))
}

// execUpsertRoleMembers inserts new or updates matching (by-primary-key) rows in role_members with given data
func (s Store) execUpsertRoleMembers(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.roleMemberTable(),
		set,
		s.preprocessColumn("rel_user", ""),
		s.preprocessColumn("rel_role", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteRoleMembers Deletes all matched (by cnd) rows in role_members with given data
func (s Store) execDeleteRoleMembers(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.roleMemberTable("rm")).Where(cnd))
}

func (s Store) internalRoleMemberRowScanner(row rowScanner) (res *types.RoleMember, err error) {
	res = &types.RoleMember{}

	if _, has := s.config.RowScanners["roleMember"]; has {
		scanner := s.config.RowScanners["roleMember"].(func(_ rowScanner, _ *types.RoleMember) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.UserID,
			&res.RoleID,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan roleMember db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryRoleMembers returns squirrel.SelectBuilder with set table and all columns
func (s Store) roleMembersSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.roleMemberTable("rm"), s.roleMemberColumns("rm")...)
}

// roleMemberTable name of the db table
func (Store) roleMemberTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "role_members" + alias
}

// RoleMemberColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) roleMemberColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "rel_user",
		alias + "rel_role",
	}
}

// {true true false false false false}

// internalRoleMemberEncoder encodes fields from types.RoleMember to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeRoleMember
// func when rdbms.customEncoder=true
func (s Store) internalRoleMemberEncoder(res *types.RoleMember) store.Payload {
	return store.Payload{
		"rel_user": res.UserID,
		"rel_role": res.RoleID,
	}
}

// checkRoleMemberConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkRoleMemberConstraints(ctx context.Context, res *types.RoleMember) error {
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
