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
)

var _ = errors.Is

// SearchRoles returns all matching rows
//
// This function calls convertRoleFilter with the given
// types.RoleFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchRoles(ctx context.Context, f types.RoleFilter) (types.RoleSet, types.RoleFilter, error) {
	var (
		err error
		set []*types.Role
		q   squirrel.SelectBuilder
	)
	q, err = s.convertRoleFilter(f)
	if err != nil {
		return nil, f, err
	}

	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reversedCursor := f.PageCursor != nil && f.PageCursor.Reverse

	// If paging with reverse cursor, change the sorting
	// direction for all columns we're sorting by
	curSort := f.Sort.Clone()
	if reversedCursor {
		curSort.Reverse()
	}

	return set, f, s.config.ErrorHandler(func() error {
		set, err = s.fetchFullPageOfRoles(ctx, q, curSort, f.PageCursor, f.Limit, f.Check)

		if err != nil {
			return err
		}

		if f.Limit > 0 && len(set) > 0 {
			if f.PageCursor != nil && (!f.PageCursor.Reverse || uint(len(set)) == f.Limit) {
				f.PrevPage = s.collectRoleCursorValues(set[0], curSort.Columns()...)
				f.PrevPage.Reverse = true
			}

			// Less items fetched then requested by page-limit
			// not very likely there's another page
			f.NextPage = s.collectRoleCursorValues(set[len(set)-1], curSort.Columns()...)
		}

		f.PageCursor = nil
		return nil
	}())
}

// fetchFullPageOfRoles collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - sorting rules (order by ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn). Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfRoles(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*types.Role) (bool, error),
) ([]*types.Role, error) {
	var (
		set  = make([]*types.Role, 0, DefaultSliceCapacity)
		aux  []*types.Role
		last *types.Role

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedCursor = cursor != nil && cursor.Reverse

		// copy of the select builder
		tryQuery squirrel.SelectBuilder

		fetched uint
		err     error
	)

	// Make sure we always end our sort by primary keys
	if sort.Get("id") == nil {
		sort = append(sort, &filter.SortExpr{Column: "id"})
	}

	// Apply sorting expr from filter to query
	if q, err = setOrderBy(q, sort, s.sortableRoleColumns()...); err != nil {
		return nil, err
	}

	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit))
		}

		if aux, fetched, last, err = s.QueryRoles(ctx, tryQuery, check); err != nil {
			return nil, err
		}

		if limit > 0 && uint(len(aux)) >= limit {
			// we should use only as much as requested
			set = append(set, aux[0:limit]...)
			break
		} else {
			set = append(set, aux...)
		}

		// if limit is not set or we've already collected enough items
		// we can break the loop right away
		if limit == 0 || fetched == 0 || fetched < limit {
			break
		}

		// In case limit is set very low and we've missed records in the first fetch,
		// make sure next fetch limit is a bit higher
		if limit < MinEnsureFetchLimit {
			limit = MinEnsureFetchLimit
		}

		// @todo improve strategy for collecting next page with lower limit

		// Point cursor to the last fetched element
		if cursor = s.collectRoleCursorValues(last, sort.Columns()...); cursor == nil {
			break
		}
	}

	if reversedCursor {
		// Cursor for previous page was used
		// Fetched set needs to be reverseCursor because we've forced a descending order to
		// get the previous page
		for i, j := 0, len(set)-1; i < j; i, j = i+1, j-1 {
			set[i], set[j] = set[j], set[i]
		}
	}

	return set, nil
}

// QueryRoles queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryRoles(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Role) (bool, error),
) ([]*types.Role, uint, *types.Role, error) {
	var (
		set = make([]*types.Role, 0, DefaultSliceCapacity)
		res *types.Role

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
			res, err = s.internalRoleRowScanner(rows)
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

// collectRoleCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectRoleCursorValues(res *types.Role, cc ...string) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{}

		hasUnique bool

		// All known primary key columns

		pkId bool

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
				case "id":
					cursor.Set(c, res.ID, false)

					pkId = true
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
	if !hasUnique || !(pkId && true) {
		collect("id")
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
