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
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
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

	return set, f, func() error {
		q, err = s.convertUserFilter(f)
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
		if q, err = setOrderBy(q, sort, s.sortableUserColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfUsers(
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

// fetchFullPageOfUsers collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfUsers(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.User) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.User, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.User

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

	set = make([]*types.User, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryUsers(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectUserCursorValues(set[collected-1], sort...)

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
		prev = s.collectUserCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectUserCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
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
) ([]*types.User, error) {
	var (
		set = make([]*types.User, 0, DefaultSliceCapacity)
		res *types.User

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalUserRowScanner(rows)
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

		"usr.deleted_at": nil,
	})
}

// LookupUserByHandle searches for user by their email
//
// It returns only valid users (not deleted, not suspended)
func (s Store) LookupUserByHandle(ctx context.Context, handle string) (*types.User, error) {
	return s.execLookupUser(ctx, squirrel.Eq{
		s.preprocessColumn("usr.handle", "lower"): store.PreprocessValue(handle, "lower"),

		"usr.deleted_at": nil,
	})
}

// LookupUserByUsername searches for user by their username
//
// It returns only valid users (not deleted, not suspended)
func (s Store) LookupUserByUsername(ctx context.Context, username string) (*types.User, error) {
	return s.execLookupUser(ctx, squirrel.Eq{
		s.preprocessColumn("usr.username", "lower"): store.PreprocessValue(username, "lower"),

		"usr.deleted_at": nil,
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
	return s.partialUserUpdate(ctx, nil, rr...)
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
			return err
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

		err = s.execUpsertUsers(ctx, s.internalUserEncoder(res))
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
			return err
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
	return s.Truncate(ctx, s.userTable())
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
	return s.Exec(ctx, s.InsertBuilder(s.userTable()).SetMap(payload))
}

// execUpdateUsers updates all matched (by cnd) rows in users with given data
func (s Store) execUpdateUsers(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.userTable("usr")).Where(cnd).SetMap(set))
}

// execUpsertUsers inserts new or updates matching (by-primary-key) rows in users with given data
func (s Store) execUpsertUsers(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.userTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteUsers Deletes all matched (by cnd) rows in users with given data
func (s Store) execDeleteUsers(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.userTable("usr")).Where(cnd))
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
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan user db row: %s", err).Wrap(err)
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

// {true true false true true true}

// sortableUserColumns returns all User columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableUserColumns() map[string]string {
	return map[string]string{
		"id": "id", "email": "email", "username": "username", "name": "name", "handle": "handle", "created_at": "created_at",
		"createdat":    "created_at",
		"updated_at":   "updated_at",
		"updatedat":    "updated_at",
		"suspended_at": "suspended_at",
		"suspendedat":  "suspended_at",
		"deleted_at":   "deleted_at",
		"deletedat":    "deleted_at",
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
func (s Store) collectUserCursorValues(res *types.User, cc ...*filter.SortExpr) *filter.PagingCursor {
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
				case "email":
					cursor.Set(c.Column, res.Email, c.Descending)
					hasUnique = true

				case "username":
					cursor.Set(c.Column, res.Username, c.Descending)
					hasUnique = true

				case "name":
					cursor.Set(c.Column, res.Name, c.Descending)

				case "handle":
					cursor.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true

				case "created_at":
					cursor.Set(c.Column, res.CreatedAt, c.Descending)

				case "updated_at":
					cursor.Set(c.Column, res.UpdatedAt, c.Descending)

				case "suspended_at":
					cursor.Set(c.Column, res.SuspendedAt, c.Descending)

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
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	{
		ex, err := s.LookupUserByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	{
		ex, err := s.LookupUserByUsername(ctx, res.Username)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
