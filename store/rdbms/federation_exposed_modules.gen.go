package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/federation_exposed_modules.yaml
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

// SearchFederationExposedModules returns all matching rows
//
// This function calls convertFederationExposedModuleFilter with the given
// types.ExposedModuleFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchFederationExposedModules(ctx context.Context, f types.ExposedModuleFilter) (types.ExposedModuleSet, types.ExposedModuleFilter, error) {
	var (
		err error
		set []*types.ExposedModule
		q   squirrel.SelectBuilder
	)
	q, err = s.convertFederationExposedModuleFilter(f)
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
		set, err = s.fetchFullPageOfFederationExposedModules(ctx, q, curSort, f.PageCursor, f.Limit, f.Check)

		if err != nil {
			return err
		}

		if f.Limit > 0 && len(set) > 0 {
			if f.PageCursor != nil && (!f.PageCursor.Reverse || uint(len(set)) == f.Limit) {
				f.PrevPage = s.collectFederationExposedModuleCursorValues(set[0], curSort.Columns()...)
				f.PrevPage.Reverse = true
			}

			// Less items fetched then requested by page-limit
			// not very likely there's another page
			f.NextPage = s.collectFederationExposedModuleCursorValues(set[len(set)-1], curSort.Columns()...)
		}

		f.PageCursor = nil
		return nil
	}())
}

// fetchFullPageOfFederationExposedModules collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - sorting rules (order by ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn). Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfFederationExposedModules(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*types.ExposedModule) (bool, error),
) ([]*types.ExposedModule, error) {
	var (
		set  = make([]*types.ExposedModule, 0, DefaultSliceCapacity)
		aux  []*types.ExposedModule
		last *types.ExposedModule

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
	if q, err = setOrderBy(q, sort, s.sortableFederationExposedModuleColumns()...); err != nil {
		return nil, err
	}

	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit))
		}

		if aux, fetched, last, err = s.QueryFederationExposedModules(ctx, tryQuery, check); err != nil {
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
		if cursor = s.collectFederationExposedModuleCursorValues(last, sort.Columns()...); cursor == nil {
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

// QueryFederationExposedModules queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryFederationExposedModules(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.ExposedModule) (bool, error),
) ([]*types.ExposedModule, uint, *types.ExposedModule, error) {
	var (
		set = make([]*types.ExposedModule, 0, DefaultSliceCapacity)
		res *types.ExposedModule

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
			res, err = s.internalFederationExposedModuleRowScanner(rows)
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

// LookupFederationExposedModuleByID searches for federation module by ID
//
// It returns federation module
func (s Store) LookupFederationExposedModuleByID(ctx context.Context, id uint64) (*types.ExposedModule, error) {
	return s.execLookupFederationExposedModule(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.id", ""): s.preprocessValue(id, ""),
	})
}

// CreateFederationExposedModule creates one or more rows in federation_module_exposed table
func (s Store) CreateFederationExposedModule(ctx context.Context, rr ...*types.ExposedModule) (err error) {
	for _, res := range rr {
		err = s.checkFederationExposedModuleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateFederationExposedModules(ctx, s.internalFederationExposedModuleEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateFederationExposedModule updates one or more existing rows in federation_module_exposed
func (s Store) UpdateFederationExposedModule(ctx context.Context, rr ...*types.ExposedModule) error {
	return s.config.ErrorHandler(s.partialFederationExposedModuleUpdate(ctx, nil, rr...))
}

// partialFederationExposedModuleUpdate updates one or more existing rows in federation_module_exposed
func (s Store) partialFederationExposedModuleUpdate(ctx context.Context, onlyColumns []string, rr ...*types.ExposedModule) (err error) {
	for _, res := range rr {
		err = s.checkFederationExposedModuleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateFederationExposedModules(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("cmd.id", ""): s.preprocessValue(res.ID, ""),
			},
			s.internalFederationExposedModuleEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// UpsertFederationExposedModule updates one or more existing rows in federation_module_exposed
func (s Store) UpsertFederationExposedModule(ctx context.Context, rr ...*types.ExposedModule) (err error) {
	for _, res := range rr {
		err = s.checkFederationExposedModuleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.config.ErrorHandler(s.execUpsertFederationExposedModules(ctx, s.internalFederationExposedModuleEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteFederationExposedModule Deletes one or more rows from federation_module_exposed table
func (s Store) DeleteFederationExposedModule(ctx context.Context, rr ...*types.ExposedModule) (err error) {
	for _, res := range rr {

		err = s.execDeleteFederationExposedModules(ctx, squirrel.Eq{
			s.preprocessColumn("cmd.id", ""): s.preprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteFederationExposedModuleByID Deletes row from the federation_module_exposed table
func (s Store) DeleteFederationExposedModuleByID(ctx context.Context, ID uint64) error {
	return s.execDeleteFederationExposedModules(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.id", ""): s.preprocessValue(ID, ""),
	})
}

// TruncateFederationExposedModules Deletes all rows from the federation_module_exposed table
func (s Store) TruncateFederationExposedModules(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.federationExposedModuleTable()))
}

// execLookupFederationExposedModule prepares FederationExposedModule query and executes it,
// returning types.ExposedModule (or error)
func (s Store) execLookupFederationExposedModule(ctx context.Context, cnd squirrel.Sqlizer) (res *types.ExposedModule, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.federationExposedModulesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalFederationExposedModuleRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateFederationExposedModules updates all matched (by cnd) rows in federation_module_exposed with given data
func (s Store) execCreateFederationExposedModules(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.federationExposedModuleTable()).SetMap(payload)))
}

// execUpdateFederationExposedModules updates all matched (by cnd) rows in federation_module_exposed with given data
func (s Store) execUpdateFederationExposedModules(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.federationExposedModuleTable("cmd")).Where(cnd).SetMap(set)))
}

// execUpsertFederationExposedModules inserts new or updates matching (by-primary-key) rows in federation_module_exposed with given data
func (s Store) execUpsertFederationExposedModules(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.federationExposedModuleTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteFederationExposedModules Deletes all matched (by cnd) rows in federation_module_exposed with given data
func (s Store) execDeleteFederationExposedModules(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.federationExposedModuleTable("cmd")).Where(cnd)))
}

func (s Store) internalFederationExposedModuleRowScanner(row rowScanner) (res *types.ExposedModule, err error) {
	res = &types.ExposedModule{}

	if _, has := s.config.RowScanners["federationExposedModule"]; has {
		scanner := s.config.RowScanners["federationExposedModule"].(func(_ rowScanner, _ *types.ExposedModule) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.NodeID,
			&res.ComposeModuleID,
			&res.Fields,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for FederationExposedModule: %w", err)
	} else {
		return res, nil
	}
}

// QueryFederationExposedModules returns squirrel.SelectBuilder with set table and all columns
func (s Store) federationExposedModulesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.federationExposedModuleTable("cmd"), s.federationExposedModuleColumns("cmd")...)
}

// federationExposedModuleTable name of the db table
func (Store) federationExposedModuleTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "federation_module_exposed" + alias
}

// FederationExposedModuleColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) federationExposedModuleColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_node",
		alias + "rel_compose_module",
		alias + "fields",
	}
}

// {true true true true true}

// sortableFederationExposedModuleColumns returns all FederationExposedModule columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableFederationExposedModuleColumns() []string {
	return []string{
		"id",
	}
}

// internalFederationExposedModuleEncoder encodes fields from types.ExposedModule to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeFederationExposedModule
// func when rdbms.customEncoder=true
func (s Store) internalFederationExposedModuleEncoder(res *types.ExposedModule) store.Payload {
	return store.Payload{
		"id":                 res.ID,
		"rel_node":           res.NodeID,
		"rel_compose_module": res.ComposeModuleID,
		"fields":             res.Fields,
	}
}

// collectFederationExposedModuleCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectFederationExposedModuleCursorValues(res *types.ExposedModule, cc ...string) *filter.PagingCursor {
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

// checkFederationExposedModuleConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkFederationExposedModuleConstraints(ctx context.Context, res *types.ExposedModule) error {
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
