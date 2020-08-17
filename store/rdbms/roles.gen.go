package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/roles.yaml
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

// SearchRoles returns all matching rows
//
// This function calls convertRoleFilter with the given
// types.RoleFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchRoles(ctx context.Context, f types.RoleFilter) (types.RoleSet, types.RoleFilter, error) {
	q, err := s.convertRoleFilter(f)
	if err != nil {
		return nil, f, err
	}

	q = ApplyPaging(q, f.PageFilter)

	scap := f.PerPage
	if scap == 0 {
		scap = DefaultSliceCapacity
	}

	var (
		set = make([]*types.Role, 0, scap)
		res *types.Role
	)

	return set, f, func() error {
		if f.Count, err = Count(ctx, s.db, q); err != nil || f.Count == 0 {
			return err
		}
		rows, err := s.Query(ctx, q)
		if err != nil {
			return err
		}

		for rows.Next() {
			if res, err = s.internalRoleRowScanner(rows, rows.Err()); err != nil {
				if cerr := rows.Close(); cerr != nil {
					return fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
				}

				return err
			}

			set = append(set, res)
		}

		return rows.Close()
	}()
}

// LookupRoleByID searches for role by ID
//
// It returns role even if deleted or suspended
func (s Store) LookupRoleByID(ctx context.Context, id uint64) (*types.Role, error) {
	return s.RoleLookup(ctx, squirrel.Eq{
		"rl.id": id,
	})
}

// LookupRoleByHandle searches for role by its handle
//
// It returns only valid roles (not deleted, not archived)
func (s Store) LookupRoleByHandle(ctx context.Context, handle string) (*types.Role, error) {
	return s.RoleLookup(ctx, squirrel.Eq{
		"rl.handle":      handle,
		"rl.archived_at": nil,
		"rl.deleted_at":  nil,
	})
}

// LookupRoleByName searches for role by its name
//
// It returns only valid roles (not deleted, not archived)
func (s Store) LookupRoleByName(ctx context.Context, name string) (*types.Role, error) {
	return s.RoleLookup(ctx, squirrel.Eq{
		"rl.name":        name,
		"rl.archived_at": nil,
		"rl.deleted_at":  nil,
	})
}

// CreateRole creates one or more rows in roles table
func (s Store) CreateRole(ctx context.Context, rr ...*types.Role) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.RoleTable()).SetMap(s.internalRoleEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateRole updates one or more existing rows in roles
func (s Store) UpdateRole(ctx context.Context, rr ...*types.Role) error {
	return s.PartialUpdateRole(ctx, nil, rr...)
}

// PartialUpdateRole updates one or more existing rows in roles
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateRole(ctx context.Context, onlyColumns []string, rr ...*types.Role) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateRoles(
				ctx,
				squirrel.Eq{s.preprocessColumn("rl.id", ""): s.preprocessValue(res.ID, "")},
				s.internalRoleEncoder(res).Skip("id").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveRole removes one or more rows from roles table
func (s Store) RemoveRole(ctx context.Context, rr ...*types.Role) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.RoleTable("rl")).Where(squirrel.Eq{s.preprocessColumn("rl.id", ""): s.preprocessValue(res.ID, "")}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveRoleByID removes row from the roles table
func (s Store) RemoveRoleByID(ctx context.Context, ID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.RoleTable("rl")).Where(squirrel.Eq{s.preprocessColumn("rl.id", ""): s.preprocessValue(ID, "")}))
}

// TruncateRoles removes all rows from the roles table
func (s Store) TruncateRoles(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.RoleTable())
}

// ExecUpdateRoles updates all matched (by cnd) rows in roles with given data
func (s Store) ExecUpdateRoles(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.RoleTable("rl")).Where(cnd).SetMap(set))
}

// RoleLookup prepares Role query and executes it,
// returning types.Role (or error)
func (s Store) RoleLookup(ctx context.Context, cnd squirrel.Sqlizer) (*types.Role, error) {
	return s.internalRoleRowScanner(s.QueryRow(ctx, s.QueryRoles().Where(cnd)))
}

func (s Store) internalRoleRowScanner(row rowScanner, err error) (*types.Role, error) {
	if err != nil {
		return nil, err
	}

	var res = &types.Role{}
	if _, has := s.config.RowScanners["role"]; has {
		scanner := s.config.RowScanners["role"].(func(rowScanner, *types.Role) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Name,
			&res.Handle,
			&res.ArchivedAt,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for Role: %w", err)
	} else {
		return res, nil
	}
}

// QueryRoles returns squirrel.SelectBuilder with set table and all columns
func (s Store) QueryRoles() squirrel.SelectBuilder {
	return s.Select(s.RoleTable("rl"), s.RoleColumns("rl")...)
}

// RoleTable name of the db table
func (Store) RoleTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "roles" + alias
}

// RoleColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) RoleColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "name",
		alias + "handle",
		alias + "archived_at",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// internalRoleEncoder encodes fields from types.Role to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeRole
// func when rdbms.customEncoder=true
func (s Store) internalRoleEncoder(res *types.Role) store.Payload {
	return store.Payload{
		"id":          res.ID,
		"name":        res.Name,
		"handle":      res.Handle,
		"archived_at": res.ArchivedAt,
		"created_at":  res.CreatedAt,
		"updated_at":  res.UpdatedAt,
		"deleted_at":  res.DeletedAt,
	}
}
