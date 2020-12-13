package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/federation_shared_modules.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchFederationSharedModules returns all matching rows
//
// This function calls convertFederationSharedModuleFilter with the given
// types.SharedModuleFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchFederationSharedModules(ctx context.Context, f types.SharedModuleFilter) (types.SharedModuleSet, types.SharedModuleFilter, error) {
	var (
		err error
		set []*types.SharedModule
		q   squirrel.SelectBuilder
	)
	q, err = s.convertFederationSharedModuleFilter(f)
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
		set, err = s.fetchFullPageOfFederationSharedModules(ctx, q, curSort, f.PageCursor, f.Limit, f.Check)

		if err != nil {
			return err
		}

		if f.Limit > 0 && len(set) > 0 {
			if f.PageCursor != nil && (!f.PageCursor.Reverse || uint(len(set)) == f.Limit) {
				f.PrevPage = s.collectFederationSharedModuleCursorValues(set[0], curSort.Columns()...)
				f.PrevPage.Reverse = true
			}

			// Less items fetched then requested by page-limit
			// not very likely there's another page
			f.NextPage = s.collectFederationSharedModuleCursorValues(set[len(set)-1], curSort.Columns()...)
		}

		f.PageCursor = nil
		return nil
	}())
}

// fetchFullPageOfFederationSharedModules collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - sorting rules (order by ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn). Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfFederationSharedModules(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*types.SharedModule) (bool, error),
) ([]*types.SharedModule, error) {
	var (
		set  = make([]*types.SharedModule, 0, DefaultSliceCapacity)
		aux  []*types.SharedModule
		last *types.SharedModule

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
	if q, err = setOrderBy(q, sort, s.sortableFederationSharedModuleColumns()...); err != nil {
		return nil, err
	}

	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit))
		}

		if aux, fetched, last, err = s.QueryFederationSharedModules(ctx, tryQuery, check); err != nil {
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
		if cursor = s.collectFederationSharedModuleCursorValues(last, sort.Columns()...); cursor == nil {
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

// QueryFederationSharedModules queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryFederationSharedModules(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.SharedModule) (bool, error),
) ([]*types.SharedModule, uint, *types.SharedModule, error) {
	var (
		set = make([]*types.SharedModule, 0, DefaultSliceCapacity)
		res *types.SharedModule

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
			res, err = s.internalFederationSharedModuleRowScanner(rows)
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

// LookupFederationSharedModuleByID searches for shared federation module by ID
//
// It returns shared federation module
func (s Store) LookupFederationSharedModuleByID(ctx context.Context, id uint64) (*types.SharedModule, error) {
	return s.execLookupFederationSharedModule(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.id", ""): s.preprocessValue(id, ""),
	})
}

// CreateFederationSharedModule creates one or more rows in federation_module_shared table
func (s Store) CreateFederationSharedModule(ctx context.Context, rr ...*types.SharedModule) (err error) {
	for _, res := range rr {
		err = s.checkFederationSharedModuleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateFederationSharedModules(ctx, s.internalFederationSharedModuleEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateFederationSharedModule updates one or more existing rows in federation_module_shared
func (s Store) UpdateFederationSharedModule(ctx context.Context, rr ...*types.SharedModule) error {
	return s.config.ErrorHandler(s.partialFederationSharedModuleUpdate(ctx, nil, rr...))
}

// partialFederationSharedModuleUpdate updates one or more existing rows in federation_module_shared
func (s Store) partialFederationSharedModuleUpdate(ctx context.Context, onlyColumns []string, rr ...*types.SharedModule) (err error) {
	for _, res := range rr {
		err = s.checkFederationSharedModuleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateFederationSharedModules(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("cmd.id", ""): s.preprocessValue(res.ID, ""),
			},
			s.internalFederationSharedModuleEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// UpsertFederationSharedModule updates one or more existing rows in federation_module_shared
func (s Store) UpsertFederationSharedModule(ctx context.Context, rr ...*types.SharedModule) (err error) {
	for _, res := range rr {
		err = s.checkFederationSharedModuleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.config.ErrorHandler(s.execUpsertFederationSharedModules(ctx, s.internalFederationSharedModuleEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteFederationSharedModule Deletes one or more rows from federation_module_shared table
func (s Store) DeleteFederationSharedModule(ctx context.Context, rr ...*types.SharedModule) (err error) {
	for _, res := range rr {

		err = s.execDeleteFederationSharedModules(ctx, squirrel.Eq{
			s.preprocessColumn("cmd.id", ""): s.preprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteFederationSharedModuleByID Deletes row from the federation_module_shared table
func (s Store) DeleteFederationSharedModuleByID(ctx context.Context, ID uint64) error {
	return s.execDeleteFederationSharedModules(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.id", ""): s.preprocessValue(ID, ""),
	})
}

// TruncateFederationSharedModules Deletes all rows from the federation_module_shared table
func (s Store) TruncateFederationSharedModules(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.federationSharedModuleTable()))
}

// execLookupFederationSharedModule prepares FederationSharedModule query and executes it,
// returning types.SharedModule (or error)
func (s Store) execLookupFederationSharedModule(ctx context.Context, cnd squirrel.Sqlizer) (res *types.SharedModule, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.federationSharedModulesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalFederationSharedModuleRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateFederationSharedModules updates all matched (by cnd) rows in federation_module_shared with given data
func (s Store) execCreateFederationSharedModules(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.federationSharedModuleTable()).SetMap(payload)))
}

// execUpdateFederationSharedModules updates all matched (by cnd) rows in federation_module_shared with given data
func (s Store) execUpdateFederationSharedModules(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.federationSharedModuleTable("cmd")).Where(cnd).SetMap(set)))
}

// execUpsertFederationSharedModules inserts new or updates matching (by-primary-key) rows in federation_module_shared with given data
func (s Store) execUpsertFederationSharedModules(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.federationSharedModuleTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteFederationSharedModules Deletes all matched (by cnd) rows in federation_module_shared with given data
func (s Store) execDeleteFederationSharedModules(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.federationSharedModuleTable("cmd")).Where(cnd)))
}

func (s Store) internalFederationSharedModuleRowScanner(row rowScanner) (res *types.SharedModule, err error) {
	res = &types.SharedModule{}

	if _, has := s.config.RowScanners["federationSharedModule"]; has {
		scanner := s.config.RowScanners["federationSharedModule"].(func(_ rowScanner, _ *types.SharedModule) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.NodeID,
			&res.Handle,
			&res.Name,
			&res.ExternalFederationModuleID,
			&res.Fields,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for FederationSharedModule: %w", err)
	} else {
		return res, nil
	}
}

// QueryFederationSharedModules returns squirrel.SelectBuilder with set table and all columns
func (s Store) federationSharedModulesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.federationSharedModuleTable("cmd"), s.federationSharedModuleColumns("cmd")...)
}

// federationSharedModuleTable name of the db table
func (Store) federationSharedModuleTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "federation_module_shared" + alias
}

// FederationSharedModuleColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) federationSharedModuleColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_node",
		alias + "handle",
		alias + "name",
		alias + "xref_module",
		alias + "fields",
	}
}

// {true true true true true}

// sortableFederationSharedModuleColumns returns all FederationSharedModule columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableFederationSharedModuleColumns() []string {
	return []string{
		"id",
	}
}

// internalFederationSharedModuleEncoder encodes fields from types.SharedModule to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeFederationSharedModule
// func when rdbms.customEncoder=true
func (s Store) internalFederationSharedModuleEncoder(res *types.SharedModule) store.Payload {
	return store.Payload{
		"id":          res.ID,
		"rel_node":    res.NodeID,
		"handle":      res.Handle,
		"name":        res.Name,
		"xref_module": res.ExternalFederationModuleID,
		"fields":      res.Fields,
	}
}

// collectFederationSharedModuleCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectFederationSharedModuleCursorValues(res *types.SharedModule, cc ...string) *filter.PagingCursor {
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

// checkFederationSharedModuleConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkFederationSharedModuleConstraints(ctx context.Context, res *types.SharedModule) error {
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
