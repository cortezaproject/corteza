package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/auth_clients.yaml
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

// SearchAuthClients returns all matching rows
//
// This function calls convertAuthClientFilter with the given
// types.AuthClientFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchAuthClients(ctx context.Context, f types.AuthClientFilter) (types.AuthClientSet, types.AuthClientFilter, error) {
	var (
		err error
		set []*types.AuthClient
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertAuthClientFilter(f)
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
		if q, err = setOrderBy(q, sort, s.sortableAuthClientColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfAuthClients(
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

// fetchFullPageOfAuthClients collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfAuthClients(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.AuthClient) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.AuthClient, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.AuthClient

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

	set = make([]*types.AuthClient, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryAuthClients(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectAuthClientCursorValues(set[collected-1], sort...)

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
		prev = s.collectAuthClientCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectAuthClientCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryAuthClients queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryAuthClients(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.AuthClient) (bool, error),
) ([]*types.AuthClient, error) {
	var (
		set = make([]*types.AuthClient, 0, DefaultSliceCapacity)
		res *types.AuthClient

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalAuthClientRowScanner(rows)
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

// LookupAuthClientByID searches for auth client by ID
//
// It returns auth client even if deleted
func (s Store) LookupAuthClientByID(ctx context.Context, id uint64) (*types.AuthClient, error) {
	return s.execLookupAuthClient(ctx, squirrel.Eq{
		s.preprocessColumn("ac.id", ""): store.PreprocessValue(id, ""),
	})
}

// LookupAuthClientByHandle searches for auth client by ID
//
// It returns auth client even if deleted
func (s Store) LookupAuthClientByHandle(ctx context.Context, handle string) (*types.AuthClient, error) {
	return s.execLookupAuthClient(ctx, squirrel.Eq{
		s.preprocessColumn("ac.handle", ""): store.PreprocessValue(handle, ""),

		"ac.deleted_at": nil,
	})
}

// CreateAuthClient creates one or more rows in auth_clients table
func (s Store) CreateAuthClient(ctx context.Context, rr ...*types.AuthClient) (err error) {
	for _, res := range rr {
		err = s.checkAuthClientConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateAuthClients(ctx, s.internalAuthClientEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateAuthClient updates one or more existing rows in auth_clients
func (s Store) UpdateAuthClient(ctx context.Context, rr ...*types.AuthClient) error {
	return s.partialAuthClientUpdate(ctx, nil, rr...)
}

// partialAuthClientUpdate updates one or more existing rows in auth_clients
func (s Store) partialAuthClientUpdate(ctx context.Context, onlyColumns []string, rr ...*types.AuthClient) (err error) {
	for _, res := range rr {
		err = s.checkAuthClientConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateAuthClients(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("ac.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalAuthClientEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertAuthClient updates one or more existing rows in auth_clients
func (s Store) UpsertAuthClient(ctx context.Context, rr ...*types.AuthClient) (err error) {
	for _, res := range rr {
		err = s.checkAuthClientConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertAuthClients(ctx, s.internalAuthClientEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAuthClient Deletes one or more rows from auth_clients table
func (s Store) DeleteAuthClient(ctx context.Context, rr ...*types.AuthClient) (err error) {
	for _, res := range rr {

		err = s.execDeleteAuthClients(ctx, squirrel.Eq{
			s.preprocessColumn("ac.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAuthClientByID Deletes row from the auth_clients table
func (s Store) DeleteAuthClientByID(ctx context.Context, ID uint64) error {
	return s.execDeleteAuthClients(ctx, squirrel.Eq{
		s.preprocessColumn("ac.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateAuthClients Deletes all rows from the auth_clients table
func (s Store) TruncateAuthClients(ctx context.Context) error {
	return s.Truncate(ctx, s.authClientTable())
}

// execLookupAuthClient prepares AuthClient query and executes it,
// returning types.AuthClient (or error)
func (s Store) execLookupAuthClient(ctx context.Context, cnd squirrel.Sqlizer) (res *types.AuthClient, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.authClientsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalAuthClientRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateAuthClients updates all matched (by cnd) rows in auth_clients with given data
func (s Store) execCreateAuthClients(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.authClientTable()).SetMap(payload))
}

// execUpdateAuthClients updates all matched (by cnd) rows in auth_clients with given data
func (s Store) execUpdateAuthClients(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.authClientTable("ac")).Where(cnd).SetMap(set))
}

// execUpsertAuthClients inserts new or updates matching (by-primary-key) rows in auth_clients with given data
func (s Store) execUpsertAuthClients(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.authClientTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteAuthClients Deletes all matched (by cnd) rows in auth_clients with given data
func (s Store) execDeleteAuthClients(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.authClientTable("ac")).Where(cnd))
}

func (s Store) internalAuthClientRowScanner(row rowScanner) (res *types.AuthClient, err error) {
	res = &types.AuthClient{}

	if _, has := s.config.RowScanners["authClient"]; has {
		scanner := s.config.RowScanners["authClient"].(func(_ rowScanner, _ *types.AuthClient) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Handle,
			&res.Meta,
			&res.Secret,
			&res.Scope,
			&res.ValidGrant,
			&res.RedirectURI,
			&res.Trusted,
			&res.Enabled,
			&res.ValidFrom,
			&res.ExpiresAt,
			&res.Security,
			&res.OwnedBy,
			&res.CreatedBy,
			&res.UpdatedBy,
			&res.DeletedBy,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan authClient db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryAuthClients returns squirrel.SelectBuilder with set table and all columns
func (s Store) authClientsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.authClientTable("ac"), s.authClientColumns("ac")...)
}

// authClientTable name of the db table
func (Store) authClientTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "auth_clients" + alias
}

// AuthClientColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) authClientColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "handle",
		alias + "meta",
		alias + "secret",
		alias + "scope",
		alias + "valid_grant",
		alias + "redirect_uri",
		alias + "trusted",
		alias + "enabled",
		alias + "valid_from",
		alias + "expires_at",
		alias + "security",
		alias + "owned_by",
		alias + "created_by",
		alias + "updated_by",
		alias + "deleted_by",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false true true true}

// sortableAuthClientColumns returns all AuthClient columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableAuthClientColumns() map[string]string {
	return map[string]string{
		"id": "id", "handle": "handle", "created_at": "created_at",
		"createdat":  "created_at",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
	}
}

// internalAuthClientEncoder encodes fields from types.AuthClient to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeAuthClient
// func when rdbms.customEncoder=true
func (s Store) internalAuthClientEncoder(res *types.AuthClient) store.Payload {
	return store.Payload{
		"id":           res.ID,
		"handle":       res.Handle,
		"meta":         res.Meta,
		"secret":       res.Secret,
		"scope":        res.Scope,
		"valid_grant":  res.ValidGrant,
		"redirect_uri": res.RedirectURI,
		"trusted":      res.Trusted,
		"enabled":      res.Enabled,
		"valid_from":   res.ValidFrom,
		"expires_at":   res.ExpiresAt,
		"security":     res.Security,
		"owned_by":     res.OwnedBy,
		"created_by":   res.CreatedBy,
		"updated_by":   res.UpdatedBy,
		"deleted_by":   res.DeletedBy,
		"created_at":   res.CreatedAt,
		"updated_at":   res.UpdatedAt,
		"deleted_at":   res.DeletedAt,
	}
}

// collectAuthClientCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectAuthClientCursorValues(res *types.AuthClient, cc ...*filter.SortExpr) *filter.PagingCursor {
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
				case "handle":
					cursor.Set(c.Column, res.Handle, c.Descending)

				case "created_at":
					cursor.Set(c.Column, res.CreatedAt, c.Descending)

				case "updated_at":
					cursor.Set(c.Column, res.UpdatedAt, c.Descending)

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

// checkAuthClientConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkAuthClientConstraints(ctx context.Context, res *types.AuthClient) error {
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
		ex, err := s.LookupAuthClientByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
