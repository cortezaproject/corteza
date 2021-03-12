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
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
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

	return set, f, func() error {
		q, err = s.convertRoleFilter(f)
		if err != nil {
			return err
		}

		// Paging enabled
		// {search: {enablePaging:true}}
		// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
		f.PrevPage, f.NextPage = nil, nil

		if f.PageCursor != nil {
			// Page cursor exists so we need to validate it against used sort
			// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
			// from the cursor.
			// This (extracted sorting info) is then returned as part of response
			if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
				return err
			}
		}

		// Make sure results are always sorted at least by primary keys
		if f.Sort.Get("id") == nil {
			f.Sort = append(f.Sort, &filter.SortExpr{
				Column:     "id",
				Descending: f.Sort.LastDescending(),
			})
		}

		// Cloned sorting instructions for the actual sorting
		// Original are passed to the fetchFullPageOfUsers fn used for cursor creation so it MUST keep the initial
		// direction information
		sort := f.Sort.Clone()

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		if f.PageCursor != nil && f.PageCursor.ROrder {
			sort.Reverse()
		}

		// Apply sorting expr from filter to query
		if q, err = setOrderBy(q, sort, s.sortableRoleColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfRoles(
			ctx,
			q, f.Sort, f.PageCursor,
			f.Limit,
			f.Check,
			func(cur *filter.PagingCursor) squirrel.Sqlizer {
				return builders.CursorCondition(cur, nil)
			},
		)

		if err != nil {
			return err
		}

		f.PageCursor = nil
		return nil
	}()
}

// fetchFullPageOfRoles collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfRoles(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.Role) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.Role, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.Role

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = cursor != nil && cursor.ROrder

		// copy of the select builder
		tryQuery squirrel.SelectBuilder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = reqItems

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = cursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool
	)

	set = make([]*types.Role, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		if cursor != nil {
			tryQuery = q.Where(cursorCond(cursor))
		} else {
			tryQuery = q
		}

		if limit > 0 {
			// fetching + 1 so we know if there are more items
			// we can fetch (next-page cursor)
			tryQuery = tryQuery.Limit(uint64(limit + 1))
		}

		if aux, err = s.QueryRoles(ctx, tryQuery, check); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 {
			// no max requested items specified, break out
			break
		}

		collected := uint(len(set))

		if reqItems > collected {
			// not enough items fetched, try again with adjusted limit
			limit = reqItems - collected

			if limit < MinEnsureFetchLimit {
				// In case limit is set very low and we've missed records in the first fetch,
				// make sure next fetch limit is a bit higher
				limit = MinEnsureFetchLimit
			}

			// Update cursor so that it points to the last item fetched
			cursor = s.collectRoleCursorValues(set[collected-1], sort...)

			// Copy reverse flag from sorting
			cursor.LThen = sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
			hasNext = true
		}

		break
	}

	collected := len(set)

	if collected == 0 {
		return nil, nil, nil, nil
	}

	if reversedOrder {
		// Fetched set needs to be reversed because we've forced a descending order to get the previous page
		for i, j := 0, collected-1; i < j; i, j = i+1, j-1 {
			set[i], set[j] = set[j], set[i]
		}

		// when in reverse-order rules on what cursor to return change
		hasPrev, hasNext = hasNext, hasPrev
	}

	if hasPrev {
		prev = s.collectRoleCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectRoleCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
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
) ([]*types.Role, error) {
	var (
		set = make([]*types.Role, 0, DefaultSliceCapacity)
		res *types.Role

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalRoleRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if check != nil {
			if chk, err := check(res); err != nil {
				return nil, err
			} else if !chk {
				continue
			}
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupRoleByID searches for role by ID
//
// It returns role even if deleted or suspended
func (s Store) LookupRoleByID(ctx context.Context, id uint64) (*types.Role, error) {
	return s.execLookupRole(ctx, squirrel.Eq{
		s.preprocessColumn("rl.id", ""): store.PreprocessValue(id, ""),
	})
}

// LookupRoleByHandle searches for role by its handle
//
// It returns only valid roles (not deleted, not archived)
func (s Store) LookupRoleByHandle(ctx context.Context, handle string) (*types.Role, error) {
	return s.execLookupRole(ctx, squirrel.Eq{
		s.preprocessColumn("rl.handle", "lower"): store.PreprocessValue(handle, "lower"),

		"rl.archived_at": nil,
		"rl.deleted_at":  nil,
	})
}

// LookupRoleByName searches for role by its name
//
// It returns only valid roles (not deleted, not archived)
func (s Store) LookupRoleByName(ctx context.Context, name string) (*types.Role, error) {
	return s.execLookupRole(ctx, squirrel.Eq{
		s.preprocessColumn("rl.name", ""): store.PreprocessValue(name, ""),

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
	return s.partialRoleUpdate(ctx, nil, rr...)
}

// partialRoleUpdate updates one or more existing rows in roles
func (s Store) partialRoleUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Role) (err error) {
	for _, res := range rr {
		err = s.checkRoleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateRoles(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("rl.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalRoleEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
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

		err = s.execUpsertRoles(ctx, s.internalRoleEncoder(res))
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
			s.preprocessColumn("rl.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteRoleByID Deletes row from the roles table
func (s Store) DeleteRoleByID(ctx context.Context, ID uint64) error {
	return s.execDeleteRoles(ctx, squirrel.Eq{
		s.preprocessColumn("rl.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateRoles Deletes all rows from the roles table
func (s Store) TruncateRoles(ctx context.Context) error {
	return s.Truncate(ctx, s.roleTable())
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
	return s.Exec(ctx, s.InsertBuilder(s.roleTable()).SetMap(payload))
}

// execUpdateRoles updates all matched (by cnd) rows in roles with given data
func (s Store) execUpdateRoles(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.roleTable("rl")).Where(cnd).SetMap(set))
}

// execUpsertRoles inserts new or updates matching (by-primary-key) rows in roles with given data
func (s Store) execUpsertRoles(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.roleTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteRoles Deletes all matched (by cnd) rows in roles with given data
func (s Store) execDeleteRoles(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.roleTable("rl")).Where(cnd))
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
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan role db row: %s", err).Wrap(err)
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

// {true true false true true true}

// sortableRoleColumns returns all Role columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableRoleColumns() map[string]string {
	return map[string]string{
		"id": "id", "name": "name", "handle": "handle", "created_at": "created_at",
		"createdat":   "created_at",
		"updated_at":  "updated_at",
		"updatedat":   "updated_at",
		"archived_at": "archived_at",
		"archivedat":  "archived_at",
		"deleted_at":  "deleted_at",
		"deletedat":   "deleted_at",
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
func (s Store) collectRoleCursorValues(res *types.Role, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		// All known primary key columns

		pkId bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cursor.Set(c.Column, res.ID, c.Descending)

					pkId = true
				case "name":
					cursor.Set(c.Column, res.Name, c.Descending)

				case "handle":
					cursor.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true

				case "created_at":
					cursor.Set(c.Column, res.CreatedAt, c.Descending)

				case "updated_at":
					cursor.Set(c.Column, res.UpdatedAt, c.Descending)

				case "archived_at":
					cursor.Set(c.Column, res.ArchivedAt, c.Descending)

				case "deleted_at":
					cursor.Set(c.Column, res.DeletedAt, c.Descending)

				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !(pkId && true) {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cursor
}

// checkRoleConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkRoleConstraints(ctx context.Context, res *types.Role) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	valid = valid && len(res.Handle) > 0

	if !valid {
		return nil
	}

	{
		ex, err := s.LookupRoleByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
