package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/applications.yaml
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

// SearchApplications returns all matching rows
//
// This function calls convertApplicationFilter with the given
// types.ApplicationFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchApplications(ctx context.Context, f types.ApplicationFilter) (types.ApplicationSet, types.ApplicationFilter, error) {
	var (
		err error
		set []*types.Application
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertApplicationFilter(f)
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
		if q, err = setOrderBy(q, sort, s.sortableApplicationColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfApplications(
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

// fetchFullPageOfApplications collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfApplications(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.Application) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.Application, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.Application

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

	set = make([]*types.Application, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryApplications(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectApplicationCursorValues(set[collected-1], sort...)

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
		prev = s.collectApplicationCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectApplicationCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryApplications queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryApplications(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Application) (bool, error),
) ([]*types.Application, error) {
	var (
		set = make([]*types.Application, 0, DefaultSliceCapacity)
		res *types.Application

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalApplicationRowScanner(rows)
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

// LookupApplicationByID searches for application by ID
//
// It returns application even if deleted
func (s Store) LookupApplicationByID(ctx context.Context, id uint64) (*types.Application, error) {
	return s.execLookupApplication(ctx, squirrel.Eq{
		s.preprocessColumn("app.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateApplication creates one or more rows in applications table
func (s Store) CreateApplication(ctx context.Context, rr ...*types.Application) (err error) {
	for _, res := range rr {
		err = s.checkApplicationConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateApplications(ctx, s.internalApplicationEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateApplication updates one or more existing rows in applications
func (s Store) UpdateApplication(ctx context.Context, rr ...*types.Application) error {
	return s.partialApplicationUpdate(ctx, nil, rr...)
}

// partialApplicationUpdate updates one or more existing rows in applications
func (s Store) partialApplicationUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Application) (err error) {
	for _, res := range rr {
		err = s.checkApplicationConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateApplications(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("app.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalApplicationEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertApplication updates one or more existing rows in applications
func (s Store) UpsertApplication(ctx context.Context, rr ...*types.Application) (err error) {
	for _, res := range rr {
		err = s.checkApplicationConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertApplications(ctx, s.internalApplicationEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteApplication Deletes one or more rows from applications table
func (s Store) DeleteApplication(ctx context.Context, rr ...*types.Application) (err error) {
	for _, res := range rr {

		err = s.execDeleteApplications(ctx, squirrel.Eq{
			s.preprocessColumn("app.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteApplicationByID Deletes row from the applications table
func (s Store) DeleteApplicationByID(ctx context.Context, ID uint64) error {
	return s.execDeleteApplications(ctx, squirrel.Eq{
		s.preprocessColumn("app.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateApplications Deletes all rows from the applications table
func (s Store) TruncateApplications(ctx context.Context) error {
	return s.Truncate(ctx, s.applicationTable())
}

// execLookupApplication prepares Application query and executes it,
// returning types.Application (or error)
func (s Store) execLookupApplication(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Application, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.applicationsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalApplicationRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateApplications updates all matched (by cnd) rows in applications with given data
func (s Store) execCreateApplications(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.applicationTable()).SetMap(payload))
}

// execUpdateApplications updates all matched (by cnd) rows in applications with given data
func (s Store) execUpdateApplications(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.applicationTable("app")).Where(cnd).SetMap(set))
}

// execUpsertApplications inserts new or updates matching (by-primary-key) rows in applications with given data
func (s Store) execUpsertApplications(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.applicationTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteApplications Deletes all matched (by cnd) rows in applications with given data
func (s Store) execDeleteApplications(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.applicationTable("app")).Where(cnd))
}

func (s Store) internalApplicationRowScanner(row rowScanner) (res *types.Application, err error) {
	res = &types.Application{}

	if _, has := s.config.RowScanners["application"]; has {
		scanner := s.config.RowScanners["application"].(func(_ rowScanner, _ *types.Application) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Name,
			&res.OwnerID,
			&res.Enabled,
			&res.Weight,
			&res.Unify,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan application db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryApplications returns squirrel.SelectBuilder with set table and all columns
func (s Store) applicationsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.applicationTable("app"), s.applicationColumns("app")...)
}

// applicationTable name of the db table
func (Store) applicationTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "applications" + alias
}

// ApplicationColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) applicationColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "name",
		alias + "rel_owner",
		alias + "enabled",
		alias + "weight",
		alias + "unify",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false true true true}

// sortableApplicationColumns returns all Application columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableApplicationColumns() map[string]string {
	return map[string]string{
		"id": "id", "name": "name", "weight": "weight", "created_at": "created_at",
		"createdat":  "created_at",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
	}
}

// internalApplicationEncoder encodes fields from types.Application to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeApplication
// func when rdbms.customEncoder=true
func (s Store) internalApplicationEncoder(res *types.Application) store.Payload {
	return store.Payload{
		"id":         res.ID,
		"name":       res.Name,
		"rel_owner":  res.OwnerID,
		"enabled":    res.Enabled,
		"weight":     res.Weight,
		"unify":      res.Unify,
		"created_at": res.CreatedAt,
		"updated_at": res.UpdatedAt,
		"deleted_at": res.DeletedAt,
	}
}

// collectApplicationCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectApplicationCursorValues(res *types.Application, cc ...*filter.SortExpr) *filter.PagingCursor {
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

				case "weight":
					cursor.Set(c.Column, res.Weight, c.Descending)

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

// checkApplicationConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkApplicationConstraints(ctx context.Context, res *types.Application) error {
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
