package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/users.yaml
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

// SearchUsers returns all matching rows
//
// This function calls convertUserFilter with the given
// types.UserFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchUsers(ctx context.Context, f types.UserFilter) (types.UserSet, types.UserFilter, error) {
	var (
		err error
		set []*types.User
		q   squirrel.SelectBuilder
	)
	q, err = s.convertUserFilter(f)
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
		set, err = s.fetchFullPageOfUsers(ctx, q, curSort, f.PageCursor, f.Limit, f.Check)

		if err != nil {
			return err
		}

		if f.Limit > 0 && len(set) > 0 {
			if f.PageCursor != nil && (!f.PageCursor.Reverse || uint(len(set)) == f.Limit) {
				f.PrevPage = s.collectUserCursorValues(set[0], curSort.Columns()...)
				f.PrevPage.Reverse = true
			}

			// Less items fetched then requested by page-limit
			// not very likely there's another page
			f.NextPage = s.collectUserCursorValues(set[len(set)-1], curSort.Columns()...)
		}

		f.PageCursor = nil
		return nil
	}())
}

// fetchFullPageOfUsers collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - sorting rules (order by ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn). Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfUsers(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*types.User) (bool, error),
) ([]*types.User, error) {
	var (
		set  = make([]*types.User, 0, DefaultSliceCapacity)
		aux  []*types.User
		last *types.User

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
	if q, err = setOrderBy(q, sort, s.sortableUserColumns()...); err != nil {
		return nil, err
	}

	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit))
		}

		if aux, fetched, last, err = s.QueryUsers(ctx, tryQuery, check); err != nil {
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
		if cursor = s.collectUserCursorValues(last, sort.Columns()...); cursor == nil {
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

// QueryUsers queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryUsers(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.User) (bool, error),
) ([]*types.User, uint, *types.User, error) {
	var (
		set = make([]*types.User, 0, DefaultSliceCapacity)
		res *types.User

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
			res, err = s.internalUserRowScanner(rows)
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

// LookupUserByID searches for user by ID
//
// It returns user even if deleted or suspended
func (s Store) LookupUserByID(ctx context.Context, id uint64) (*types.User, error) {
	return s.execLookupUser(ctx, squirrel.Eq{
		s.preprocessColumn("usr.id", ""): store.PreprocessValue(id, ""),
	})
}

// LookupUserByEmail searches for user by their email
//
// It returns only valid users (not deleted, not suspended)
func (s Store) LookupUserByEmail(ctx context.Context, email string) (*types.User, error) {
	return s.execLookupUser(ctx, squirrel.Eq{
		s.preprocessColumn("usr.email", "lower"): store.PreprocessValue(email, "lower"),

		"usr.deleted_at":   nil,
		"usr.suspended_at": nil,
	})
}

// LookupUserByHandle searches for user by their email
//
// It returns only valid users (not deleted, not suspended)
func (s Store) LookupUserByHandle(ctx context.Context, handle string) (*types.User, error) {
	return s.execLookupUser(ctx, squirrel.Eq{
		s.preprocessColumn("usr.handle", "lower"): store.PreprocessValue(handle, "lower"),

		"usr.deleted_at":   nil,
		"usr.suspended_at": nil,
	})
}

// LookupUserByUsername searches for user by their username
//
// It returns only valid users (not deleted, not suspended)
func (s Store) LookupUserByUsername(ctx context.Context, username string) (*types.User, error) {
	return s.execLookupUser(ctx, squirrel.Eq{
		s.preprocessColumn("usr.username", "lower"): store.PreprocessValue(username, "lower"),

		"usr.deleted_at":   nil,
		"usr.suspended_at": nil,
	})
}

// CreateUser creates one or more rows in users table
func (s Store) CreateUser(ctx context.Context, rr ...*types.User) (err error) {
	for _, res := range rr {
		err = s.checkUserConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateUsers(ctx, s.internalUserEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateUser updates one or more existing rows in users
func (s Store) UpdateUser(ctx context.Context, rr ...*types.User) error {
	return s.config.ErrorHandler(s.partialUserUpdate(ctx, nil, rr...))
}

// partialUserUpdate updates one or more existing rows in users
func (s Store) partialUserUpdate(ctx context.Context, onlyColumns []string, rr ...*types.User) (err error) {
	for _, res := range rr {
		err = s.checkUserConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateUsers(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("usr.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalUserEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// UpsertUser updates one or more existing rows in users
func (s Store) UpsertUser(ctx context.Context, rr ...*types.User) (err error) {
	for _, res := range rr {
		err = s.checkUserConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.config.ErrorHandler(s.execUpsertUsers(ctx, s.internalUserEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteUser Deletes one or more rows from users table
func (s Store) DeleteUser(ctx context.Context, rr ...*types.User) (err error) {
	for _, res := range rr {

		err = s.execDeleteUsers(ctx, squirrel.Eq{
			s.preprocessColumn("usr.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteUserByID Deletes row from the users table
func (s Store) DeleteUserByID(ctx context.Context, ID uint64) error {
	return s.execDeleteUsers(ctx, squirrel.Eq{
		s.preprocessColumn("usr.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateUsers Deletes all rows from the users table
func (s Store) TruncateUsers(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.userTable()))
}

// execLookupUser prepares User query and executes it,
// returning types.User (or error)
func (s Store) execLookupUser(ctx context.Context, cnd squirrel.Sqlizer) (res *types.User, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.usersSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalUserRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateUsers updates all matched (by cnd) rows in users with given data
func (s Store) execCreateUsers(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.userTable()).SetMap(payload)))
}

// execUpdateUsers updates all matched (by cnd) rows in users with given data
func (s Store) execUpdateUsers(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.userTable("usr")).Where(cnd).SetMap(set)))
}

// execUpsertUsers inserts new or updates matching (by-primary-key) rows in users with given data
func (s Store) execUpsertUsers(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.userTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteUsers Deletes all matched (by cnd) rows in users with given data
func (s Store) execDeleteUsers(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.userTable("usr")).Where(cnd)))
}

func (s Store) internalUserRowScanner(row rowScanner) (res *types.User, err error) {
	res = &types.User{}

	if _, has := s.config.RowScanners["user"]; has {
		scanner := s.config.RowScanners["user"].(func(_ rowScanner, _ *types.User) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Email,
			&res.EmailConfirmed,
			&res.Username,
			&res.Name,
			&res.Handle,
			&res.Meta,
			&res.Kind,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.SuspendedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for User: %w", err)
	} else {
		return res, nil
	}
}

// QueryUsers returns squirrel.SelectBuilder with set table and all columns
func (s Store) usersSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.userTable("usr"), s.userColumns("usr")...)
}

// userTable name of the db table
func (Store) userTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "users" + alias
}

// UserColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) userColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "email",
		alias + "email_confirmed",
		alias + "username",
		alias + "name",
		alias + "handle",
		alias + "meta",
		alias + "kind",
		alias + "created_at",
		alias + "updated_at",
		alias + "suspended_at",
		alias + "deleted_at",
	}
}

// {true true true true true}

// sortableUserColumns returns all User columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableUserColumns() []string {
	return []string{
		"id",
		"email",
		"username",
		"name",
		"handle",
		"created_at",
		"updated_at",
		"suspended_at",
		"deleted_at",
	}
}

// internalUserEncoder encodes fields from types.User to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeUser
// func when rdbms.customEncoder=true
func (s Store) internalUserEncoder(res *types.User) store.Payload {
	return store.Payload{
		"id":              res.ID,
		"email":           res.Email,
		"email_confirmed": res.EmailConfirmed,
		"username":        res.Username,
		"name":            res.Name,
		"handle":          res.Handle,
		"meta":            res.Meta,
		"kind":            res.Kind,
		"created_at":      res.CreatedAt,
		"updated_at":      res.UpdatedAt,
		"suspended_at":    res.SuspendedAt,
		"deleted_at":      res.DeletedAt,
	}
}

// collectUserCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectUserCursorValues(res *types.User, cc ...string) *filter.PagingCursor {
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
				case "email":
					cursor.Set(c, res.Email, false)
					hasUnique = true

				case "username":
					cursor.Set(c, res.Username, false)
					hasUnique = true

				case "name":
					cursor.Set(c, res.Name, false)

				case "handle":
					cursor.Set(c, res.Handle, false)
					hasUnique = true

				case "created_at":
					cursor.Set(c, res.CreatedAt, false)

				case "updated_at":
					cursor.Set(c, res.UpdatedAt, false)

				case "suspended_at":
					cursor.Set(c, res.SuspendedAt, false)

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

// checkUserConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkUserConstraints(ctx context.Context, res *types.User) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	valid = valid && len(res.Email) > 0

	valid = valid && len(res.Handle) > 0

	valid = valid && len(res.Username) > 0

	if !valid {
		return nil
	}

	{
		ex, err := s.LookupUserByEmail(ctx, res.Email)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique
		} else if !errors.Is(err, store.ErrNotFound) {
			return err
		}
	}

	{
		ex, err := s.LookupUserByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique
		} else if !errors.Is(err, store.ErrNotFound) {
			return err
		}
	}

	{
		ex, err := s.LookupUserByUsername(ctx, res.Username)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique
		} else if !errors.Is(err, store.ErrNotFound) {
			return err
		}
	}

	return nil
}
