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
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"strings"
)

var _ = errors.Is

const (
	TriggerBeforeRoleCreate triggerKey = "roleBeforeCreate"
	TriggerBeforeRoleUpdate triggerKey = "roleBeforeUpdate"
	TriggerBeforeRoleUpsert triggerKey = "roleBeforeUpsert"
	TriggerBeforeRoleDelete triggerKey = "roleBeforeDelete"
)

// SearchRoles returns all matching rows
//
// This function calls convertRoleFilter with the given
// types.RoleFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchRoles(ctx context.Context, f types.RoleFilter) (types.RoleSet, types.RoleFilter, error) {
	var scap uint
	q, err := s.convertRoleFilter(f)
	if err != nil {
		return nil, f, err
	}

	scap = f.Limit

	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reverseCursor := f.PageCursor != nil && f.PageCursor.Reverse

	if err := f.Sort.Validate(s.sortableRoleColumns()...); err != nil {
		return nil, f, fmt.Errorf("could not validate sort: %v", err)
	}

	// If paging with reverse cursor, change the sorting
	// direction for all columns we're sorting by
	sort := f.Sort.Clone()
	if reverseCursor {
		sort.Reverse()
	}

	// Apply sorting expr from filter to query
	if len(sort) > 0 {
		sqlSort := make([]string, len(sort))
		for i := range sort {
			sqlSort[i] = sort[i].Column
			if sort[i].Descending {
				sqlSort[i] += " DESC"
			}
		}

		q = q.OrderBy(sqlSort...)
	}

	if scap == 0 {
		scap = DefaultSliceCapacity
	}

	var (
		set = make([]*types.Role, 0, scap)
		// fetches rows and scans them into types.Role resource this is then passed to Check function on filter
		// to help determine if fetched resource fits or not
		//
		// Note that limit is passed explicitly and is not necessarily equal to filter's limit. We want
		// to keep that value intact.
		//
		// The value for cursor is used and set directly from/to the filter!
		//
		// It returns total number of fetched pages and modifies PageCursor value for paging
		fetchPage = func(cursor *filter.PagingCursor, limit uint) (fetched uint, err error) {
			var (
				res *types.Role

				// Make a copy of the select query builder so that we don't change
				// the original query
				slct = q.Options()
			)

			if limit > 0 {
				slct = slct.Limit(uint64(limit))

				if cursor != nil && len(cursor.Keys()) > 0 {
					const cursorTpl = `(%s) %s (?%s)`
					op := ">"
					if cursor.Reverse {
						op = "<"
					}

					pred := fmt.Sprintf(cursorTpl, strings.Join(cursor.Keys(), ", "), op, strings.Repeat(", ?", len(cursor.Keys())-1))
					slct = slct.Where(pred, cursor.Values()...)
				}
			}

			rows, err := s.Query(ctx, slct)
			if err != nil {
				return
			}

			for rows.Next() {
				fetched++

				if rows.Err() == nil {
					res, err = s.internalRoleRowScanner(rows)
				}

				if err != nil {
					if cerr := rows.Close(); cerr != nil {
						err = fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return
				}

				// If check function is set, call it and act accordingly

				if f.Check != nil {
					var chk bool
					if chk, err = f.Check(res); err != nil {
						if cerr := rows.Close(); cerr != nil {
							err = fmt.Errorf("could not close rows (%v) after check error: %w", cerr, err)
						}

						return
					} else if !chk {
						// did not pass the check
						// go with the next row
						continue
					}
				}
				set = append(set, res)

				if f.Limit > 0 {
					if uint(len(set)) >= f.Limit {
						// make sure we do not fetch more than requested!
						break
					}
				}
			}

			err = rows.Close()
			return
		}

		fetch = func() error {
			var (
				// how many items were actually fetched
				fetched uint

				// starting offset & limit are from filter arg
				// note that this will have to be improved with key-based pagination
				limit = f.Limit

				// Copy cursor value
				//
				// This is where we'll start fetching and this value will be overwritten when
				// results come back
				cursor = f.PageCursor

				lastSetFull bool
			)

			for refetch := 0; refetch < MaxRefetches; refetch++ {
				if fetched, err = fetchPage(cursor, limit); err != nil {
					return err
				}

				// if limit is not set or we've already collected enough items
				// we can break the loop right away
				if limit == 0 || fetched == 0 || fetched < limit {
					break
				}

				if uint(len(set)) >= f.Limit {
					// we should return as much as requested
					set = set[0:f.Limit]
					lastSetFull = true
					break
				}

				// In case limit is set very low and we've missed records in the first fetch,
				// make sure next fetch limit is a bit higher
				if limit < MinRefetchLimit {
					limit = MinRefetchLimit
				}

				// @todo it might be good to implement different kind of strategies
				//       (beyond min-refetch-limit above) that can adjust limit on
				//       retry to more optimal number
			}

			if reverseCursor {
				// Cursor for previous page was used
				// Fetched set needs to be reverseCursor because we've forced a descending order to
				// get the previus page
				for i, j := 0, len(set)-1; i < j; i, j = i+1, j-1 {
					set[i], set[j] = set[j], set[i]
				}
			}

			if f.Limit > 0 && len(set) > 0 {
				if f.PageCursor != nil && (!f.PageCursor.Reverse || lastSetFull) {
					f.PrevPage = s.collectRoleCursorValues(set[0], sort.Columns()...)
					f.PrevPage.Reverse = true
				}

				// Less items fetched then requested by page-limit
				// not very likely there's another page
				f.NextPage = s.collectRoleCursorValues(set[len(set)-1], sort.Columns()...)
			}

			f.PageCursor = nil
			return nil
		}
	)

	return set, f, s.config.ErrorHandler(fetch())
}

// LookupRoleByID searches for role by ID
//
// It returns role even if deleted or suspended
func (s Store) LookupRoleByID(ctx context.Context, id uint64) (*types.Role, error) {
	return s.execLookupRole(ctx, squirrel.Eq{
		s.preprocessColumn("rl.id", ""): s.preprocessValue(id, ""),
	})
}

// LookupRoleByHandle searches for role by its handle
//
// It returns only valid roles (not deleted, not archived)
func (s Store) LookupRoleByHandle(ctx context.Context, handle string) (*types.Role, error) {
	return s.execLookupRole(ctx, squirrel.Eq{
		s.preprocessColumn("rl.handle", "lower"): s.preprocessValue(handle, "lower"),

		"rl.archived_at": nil,
		"rl.deleted_at":  nil,
	})
}

// LookupRoleByName searches for role by its name
//
// It returns only valid roles (not deleted, not archived)
func (s Store) LookupRoleByName(ctx context.Context, name string) (*types.Role, error) {
	return s.execLookupRole(ctx, squirrel.Eq{
		s.preprocessColumn("rl.name", ""): s.preprocessValue(name, ""),

		"rl.archived_at": nil,
		"rl.deleted_at":  nil,
	})
}

// CreateRole creates one or more rows in roles table
func (s Store) CreateRole(ctx context.Context, rr ...*types.Role) (err error) {
	for _, res := range rr {
		err = s.checkRoleConstraints(ctx, res)
		if err != nil {
			return err
		}

		// err = s.roleHook(ctx, TriggerBeforeRoleCreate, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execCreateRoles(ctx, s.internalRoleEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateRole updates one or more existing rows in roles
func (s Store) UpdateRole(ctx context.Context, rr ...*types.Role) error {
	return s.config.ErrorHandler(s.PartialRoleUpdate(ctx, nil, rr...))
}

// PartialRoleUpdate updates one or more existing rows in roles
func (s Store) PartialRoleUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Role) (err error) {
	for _, res := range rr {
		err = s.checkRoleConstraints(ctx, res)
		if err != nil {
			return err
		}

		// err = s.roleHook(ctx, TriggerBeforeRoleUpdate, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execUpdateRoles(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("rl.id", ""): s.preprocessValue(res.ID, ""),
			},
			s.internalRoleEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// UpsertRole updates one or more existing rows in roles
func (s Store) UpsertRole(ctx context.Context, rr ...*types.Role) (err error) {
	for _, res := range rr {
		err = s.checkRoleConstraints(ctx, res)
		if err != nil {
			return err
		}

		// err = s.roleHook(ctx, TriggerBeforeRoleUpsert, res)
		// if err != nil {
		// 	return err
		// }

		err = s.config.ErrorHandler(s.execUpsertRoles(ctx, s.internalRoleEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteRole Deletes one or more rows from roles table
func (s Store) DeleteRole(ctx context.Context, rr ...*types.Role) (err error) {
	for _, res := range rr {
		// err = s.roleHook(ctx, TriggerBeforeRoleDelete, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execDeleteRoles(ctx, squirrel.Eq{
			s.preprocessColumn("rl.id", ""): s.preprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteRoleByID Deletes row from the roles table
func (s Store) DeleteRoleByID(ctx context.Context, ID uint64) error {
	return s.execDeleteRoles(ctx, squirrel.Eq{
		s.preprocessColumn("rl.id", ""): s.preprocessValue(ID, ""),
	})
}

// TruncateRoles Deletes all rows from the roles table
func (s Store) TruncateRoles(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.roleTable()))
}

// execLookupRole prepares Role query and executes it,
// returning types.Role (or error)
func (s Store) execLookupRole(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Role, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.rolesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalRoleRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateRoles updates all matched (by cnd) rows in roles with given data
func (s Store) execCreateRoles(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.roleTable()).SetMap(payload)))
}

// execUpdateRoles updates all matched (by cnd) rows in roles with given data
func (s Store) execUpdateRoles(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.roleTable("rl")).Where(cnd).SetMap(set)))
}

// execUpsertRoles inserts new or updates matching (by-primary-key) rows in roles with given data
func (s Store) execUpsertRoles(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.roleTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteRoles Deletes all matched (by cnd) rows in roles with given data
func (s Store) execDeleteRoles(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.roleTable("rl")).Where(cnd)))
}

func (s Store) internalRoleRowScanner(row rowScanner) (res *types.Role, err error) {
	res = &types.Role{}

	if _, has := s.config.RowScanners["role"]; has {
		scanner := s.config.RowScanners["role"].(func(_ rowScanner, _ *types.Role) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Name,
			&res.Handle,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.ArchivedAt,
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
func (s Store) rolesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.roleTable("rl"), s.roleColumns("rl")...)
}

// roleTable name of the db table
func (Store) roleTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "roles" + alias
}

// RoleColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) roleColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "name",
		alias + "handle",
		alias + "created_at",
		alias + "updated_at",
		alias + "archived_at",
		alias + "deleted_at",
	}
}

// {true true true true true}

// sortableRoleColumns returns all Role columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableRoleColumns() []string {
	return []string{
		"id",
		"name",
		"handle",
		"created_at",
		"updated_at",
		"archived_at",
		"deleted_at",
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
		"created_at":  res.CreatedAt,
		"updated_at":  res.UpdatedAt,
		"archived_at": res.ArchivedAt,
		"deleted_at":  res.DeletedAt,
	}
}

func (s Store) collectRoleCursorValues(res *types.Role, cc ...string) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{}

		hasUnique bool

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
				case "id":
					cursor.Set(c, res.ID, false)
				case "name":
					cursor.Set(c, res.Name, false)
				case "handle":
					cursor.Set(c, res.Handle, false)
					hasUnique = true
				case "created_at":
					cursor.Set(c, res.CreatedAt, false)
				case "updated_at":
					cursor.Set(c, res.UpdatedAt, false)
				case "archived_at":
					cursor.Set(c, res.ArchivedAt, false)
				case "deleted_at":
					cursor.Set(c, res.DeletedAt, false)

				}
			}
		}
	)

	collect(cc...)
	if !hasUnique {
		collect(
			"id",
		)
	}

	return cursor
}

func (s *Store) checkRoleConstraints(ctx context.Context, res *types.Role) error {

	{
		ex, err := s.LookupRoleByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique
		} else if !errors.Is(err, store.ErrNotFound) {
			return err
		}
	}

	return nil
}

// func (s *Store) roleHook(ctx context.Context, key triggerKey, res *types.Role) error {
// 	if fn, has := s.config.TriggerHandlers[key]; has {
// 		return fn.(func (ctx context.Context, s *Store, res *types.Role) error)(ctx, s, res)
// 	}
//
// 	return nil
// }
