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
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/jmoiron/sqlx"
)

// SearchRoleMembers returns all matching rows
//
// This function calls convertRoleMemberFilter with the given
// types.RoleMemberFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchRoleMembers(ctx context.Context, f types.RoleMemberFilter) (types.RoleMemberSet, types.RoleMemberFilter, error) {
	var scap uint
	q := s.QueryRoleMembers()

	if scap == 0 {
		scap = DefaultSliceCapacity
	}

	var (
		set = make([]*types.RoleMember, 0, scap)
		// Paging is disabled in definition yaml file
		// {search: {disablePaging:true}} and this allows
		// a much simpler row fetching logic
		fetch = func() error {
			var (
				res       *types.RoleMember
				rows, err = s.Query(ctx, q)
			)

			if err != nil {
				return err
			}

			for rows.Next() {
				if res, err = s.internalRoleMemberRowScanner(rows, rows.Err()); err != nil {
					if cerr := rows.Close(); cerr != nil {
						return fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return err
				}

				set = append(set, res)
			}

			return rows.Close()
		}
	)

	return set, f, fetch()
}

// CreateRoleMember creates one or more rows in role_members table
func (s Store) CreateRoleMember(ctx context.Context, rr ...*types.RoleMember) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.RoleMemberTable()).SetMap(s.internalRoleMemberEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateRoleMember updates one or more existing rows in role_members
func (s Store) UpdateRoleMember(ctx context.Context, rr ...*types.RoleMember) error {
	return s.PartialUpdateRoleMember(ctx, nil, rr...)
}

// PartialUpdateRoleMember updates one or more existing rows in role_members
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateRoleMember(ctx context.Context, onlyColumns []string, rr ...*types.RoleMember) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateRoleMembers(
				ctx,
				squirrel.Eq{s.preprocessColumn("rm.rel_user", ""): s.preprocessValue(res.UserID, ""),
					s.preprocessColumn("rm.rel_role", ""): s.preprocessValue(res.RoleID, ""),
				},
				s.internalRoleMemberEncoder(res).Skip("rel_user", "rel_role").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveRoleMember removes one or more rows from role_members table
func (s Store) RemoveRoleMember(ctx context.Context, rr ...*types.RoleMember) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.RoleMemberTable("rm")).Where(squirrel.Eq{s.preprocessColumn("rm.rel_user", ""): s.preprocessValue(res.UserID, ""),
				s.preprocessColumn("rm.rel_role", ""): s.preprocessValue(res.RoleID, ""),
			}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveRoleMemberByUserIDRoleID removes row from the role_members table
func (s Store) RemoveRoleMemberByUserIDRoleID(ctx context.Context, userID uint64, roleID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.RoleMemberTable("rm")).Where(squirrel.Eq{s.preprocessColumn("rm.rel_user", ""): s.preprocessValue(userID, ""),

		s.preprocessColumn("rm.rel_role", ""): s.preprocessValue(roleID, ""),
	}))
}

// TruncateRoleMembers removes all rows from the role_members table
func (s Store) TruncateRoleMembers(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.RoleMemberTable())
}

// ExecUpdateRoleMembers updates all matched (by cnd) rows in role_members with given data
func (s Store) ExecUpdateRoleMembers(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.RoleMemberTable("rm")).Where(cnd).SetMap(set))
}

// RoleMemberLookup prepares RoleMember query and executes it,
// returning types.RoleMember (or error)
func (s Store) RoleMemberLookup(ctx context.Context, cnd squirrel.Sqlizer) (*types.RoleMember, error) {
	return s.internalRoleMemberRowScanner(s.QueryRow(ctx, s.QueryRoleMembers().Where(cnd)))
}

func (s Store) internalRoleMemberRowScanner(row rowScanner, err error) (*types.RoleMember, error) {
	if err != nil {
		return nil, err
	}

	var res = &types.RoleMember{}
	if _, has := s.config.RowScanners["roleMember"]; has {
		scanner := s.config.RowScanners["roleMember"].(func(rowScanner, *types.RoleMember) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.UserID,
			&res.RoleID,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for RoleMember: %w", err)
	} else {
		return res, nil
	}
}

// QueryRoleMembers returns squirrel.SelectBuilder with set table and all columns
func (s Store) QueryRoleMembers() squirrel.SelectBuilder {
	return s.Select(s.RoleMemberTable("rm"), s.RoleMemberColumns("rm")...)
}

// RoleMemberTable name of the db table
func (Store) RoleMemberTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "role_members" + alias
}

// RoleMemberColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) RoleMemberColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "rel_user",
		alias + "rel_role",
	}
}

// {false true true false}

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
