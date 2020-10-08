package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/federation_module_mappings.yaml
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

// SearchFederationModuleMappings returns all matching rows
//
// This function calls convertFederationModuleMappingFilter with the given
// types.ModuleMappingFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchFederationModuleMappings(ctx context.Context, f types.ModuleMappingFilter) (types.ModuleMappingSet, types.ModuleMappingFilter, error) {
	var (
		err error
		set []*types.ModuleMapping
		q   squirrel.SelectBuilder
	)
	q, err = s.convertFederationModuleMappingFilter(f)
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
		set, err = s.fetchFullPageOfFederationModuleMappings(ctx, q, curSort, f.PageCursor, f.Limit, f.Check)

		if err != nil {
			return err
		}

		if f.Limit > 0 && len(set) > 0 {
			if f.PageCursor != nil && (!f.PageCursor.Reverse || uint(len(set)) == f.Limit) {
				f.PrevPage = s.collectFederationModuleMappingCursorValues(set[0], curSort.Columns()...)
				f.PrevPage.Reverse = true
			}

			// Less items fetched then requested by page-limit
			// not very likely there's another page
			f.NextPage = s.collectFederationModuleMappingCursorValues(set[len(set)-1], curSort.Columns()...)
		}

		f.PageCursor = nil
		return nil
	}())
}

// fetchFullPageOfFederationModuleMappings collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - sorting rules (order by ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn). Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfFederationModuleMappings(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*types.ModuleMapping) (bool, error),
) ([]*types.ModuleMapping, error) {
	var (
		set  = make([]*types.ModuleMapping, 0, DefaultSliceCapacity)
		aux  []*types.ModuleMapping
		last *types.ModuleMapping

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedCursor = cursor != nil && cursor.Reverse

		// copy of the select builder
		tryQuery squirrel.SelectBuilder

		fetched uint
		err     error
	)

	// Make sure we always end our sort by primary keys
	if sort.Get("rel_federation_module") == nil {
		sort = append(sort, &filter.SortExpr{Column: "rel_federation_module"})
	}

	if sort.Get("rel_compose_module") == nil {
		sort = append(sort, &filter.SortExpr{Column: "rel_compose_module"})
	}

	// Apply sorting expr from filter to query
	if q, err = setOrderBy(q, sort, s.sortableFederationModuleMappingColumns()...); err != nil {
		return nil, err
	}

	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit))
		}

		if aux, fetched, last, err = s.QueryFederationModuleMappings(ctx, tryQuery, check); err != nil {
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
		if cursor = s.collectFederationModuleMappingCursorValues(last, sort.Columns()...); cursor == nil {
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

// QueryFederationModuleMappings queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryFederationModuleMappings(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.ModuleMapping) (bool, error),
) ([]*types.ModuleMapping, uint, *types.ModuleMapping, error) {
	var (
		set = make([]*types.ModuleMapping, 0, DefaultSliceCapacity)
		res *types.ModuleMapping

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
			res, err = s.internalFederationModuleMappingRowScanner(rows)
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

// LookupFederationModuleMappingByFederationModuleIDComposeModuleID searches for module mapping by federation module id and compose module id
//
// It returns module mapping
func (s Store) LookupFederationModuleMappingByFederationModuleIDComposeModuleID(ctx context.Context, federation_module_id uint64, compose_module_id uint64) (*types.ModuleMapping, error) {
	return s.execLookupFederationModuleMapping(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.rel_federation_module", ""): s.preprocessValue(federation_module_id, ""),
		s.preprocessColumn("cmd.rel_compose_module", ""):    s.preprocessValue(compose_module_id, ""),
	})
}

// LookupFederationModuleMappingByFederationModuleID searches for module mapping by federation module id
//
// It returns module mapping
func (s Store) LookupFederationModuleMappingByFederationModuleID(ctx context.Context, federation_module_id uint64) (*types.ModuleMapping, error) {
	return s.execLookupFederationModuleMapping(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.rel_federation_module", ""): s.preprocessValue(federation_module_id, ""),
	})
}

// CreateFederationModuleMapping creates one or more rows in federation_module_mapping table
func (s Store) CreateFederationModuleMapping(ctx context.Context, rr ...*types.ModuleMapping) (err error) {
	for _, res := range rr {
		err = s.checkFederationModuleMappingConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateFederationModuleMappings(ctx, s.internalFederationModuleMappingEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateFederationModuleMapping updates one or more existing rows in federation_module_mapping
func (s Store) UpdateFederationModuleMapping(ctx context.Context, rr ...*types.ModuleMapping) error {
	return s.config.ErrorHandler(s.partialFederationModuleMappingUpdate(ctx, nil, rr...))
}

// partialFederationModuleMappingUpdate updates one or more existing rows in federation_module_mapping
func (s Store) partialFederationModuleMappingUpdate(ctx context.Context, onlyColumns []string, rr ...*types.ModuleMapping) (err error) {
	for _, res := range rr {
		err = s.checkFederationModuleMappingConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateFederationModuleMappings(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("cmd.rel_federation_module", ""): s.preprocessValue(res.FederationModuleID, ""), s.preprocessColumn("cmd.rel_compose_module", ""): s.preprocessValue(res.ComposeModuleID, ""),
			},
			s.internalFederationModuleMappingEncoder(res).Skip("rel_federation_module", "rel_compose_module").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// UpsertFederationModuleMapping updates one or more existing rows in federation_module_mapping
func (s Store) UpsertFederationModuleMapping(ctx context.Context, rr ...*types.ModuleMapping) (err error) {
	for _, res := range rr {
		err = s.checkFederationModuleMappingConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.config.ErrorHandler(s.execUpsertFederationModuleMappings(ctx, s.internalFederationModuleMappingEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteFederationModuleMapping Deletes one or more rows from federation_module_mapping table
func (s Store) DeleteFederationModuleMapping(ctx context.Context, rr ...*types.ModuleMapping) (err error) {
	for _, res := range rr {

		err = s.execDeleteFederationModuleMappings(ctx, squirrel.Eq{
			s.preprocessColumn("cmd.rel_federation_module", ""): s.preprocessValue(res.FederationModuleID, ""), s.preprocessColumn("cmd.rel_compose_module", ""): s.preprocessValue(res.ComposeModuleID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteFederationModuleMappingByFederationModuleIDComposeModuleID Deletes row from the federation_module_mapping table
func (s Store) DeleteFederationModuleMappingByFederationModuleIDComposeModuleID(ctx context.Context, federationModuleID uint64, composeModuleID uint64) error {
	return s.execDeleteFederationModuleMappings(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.rel_federation_module", ""): s.preprocessValue(federationModuleID, ""),
		s.preprocessColumn("cmd.rel_compose_module", ""):    s.preprocessValue(composeModuleID, ""),
	})
}

// TruncateFederationModuleMappings Deletes all rows from the federation_module_mapping table
func (s Store) TruncateFederationModuleMappings(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.federationModuleMappingTable()))
}

// execLookupFederationModuleMapping prepares FederationModuleMapping query and executes it,
// returning types.ModuleMapping (or error)
func (s Store) execLookupFederationModuleMapping(ctx context.Context, cnd squirrel.Sqlizer) (res *types.ModuleMapping, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.federationModuleMappingsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalFederationModuleMappingRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateFederationModuleMappings updates all matched (by cnd) rows in federation_module_mapping with given data
func (s Store) execCreateFederationModuleMappings(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.federationModuleMappingTable()).SetMap(payload)))
}

// execUpdateFederationModuleMappings updates all matched (by cnd) rows in federation_module_mapping with given data
func (s Store) execUpdateFederationModuleMappings(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.federationModuleMappingTable("cmd")).Where(cnd).SetMap(set)))
}

// execUpsertFederationModuleMappings inserts new or updates matching (by-primary-key) rows in federation_module_mapping with given data
func (s Store) execUpsertFederationModuleMappings(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.federationModuleMappingTable(),
		set,
		"rel_federation_module",
		"rel_compose_module",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteFederationModuleMappings Deletes all matched (by cnd) rows in federation_module_mapping with given data
func (s Store) execDeleteFederationModuleMappings(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.federationModuleMappingTable("cmd")).Where(cnd)))
}

func (s Store) internalFederationModuleMappingRowScanner(row rowScanner) (res *types.ModuleMapping, err error) {
	res = &types.ModuleMapping{}

	if _, has := s.config.RowScanners["federationModuleMapping"]; has {
		scanner := s.config.RowScanners["federationModuleMapping"].(func(_ rowScanner, _ *types.ModuleMapping) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.FederationModuleID,
			&res.ComposeModuleID,
			&res.FieldMapping,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for FederationModuleMapping: %w", err)
	} else {
		return res, nil
	}
}

// QueryFederationModuleMappings returns squirrel.SelectBuilder with set table and all columns
func (s Store) federationModuleMappingsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.federationModuleMappingTable("cmd"), s.federationModuleMappingColumns("cmd")...)
}

// federationModuleMappingTable name of the db table
func (Store) federationModuleMappingTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "federation_module_mapping" + alias
}

// FederationModuleMappingColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) federationModuleMappingColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "rel_federation_module",
		alias + "rel_compose_module",
		alias + "field_mapping",
	}
}

// {true true true true true}

// sortableFederationModuleMappingColumns returns all FederationModuleMapping columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableFederationModuleMappingColumns() []string {
	return []string{}
}

// internalFederationModuleMappingEncoder encodes fields from types.ModuleMapping to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeFederationModuleMapping
// func when rdbms.customEncoder=true
func (s Store) internalFederationModuleMappingEncoder(res *types.ModuleMapping) store.Payload {
	return store.Payload{
		"rel_federation_module": res.FederationModuleID,
		"rel_compose_module":    res.ComposeModuleID,
		"field_mapping":         res.FieldMapping,
	}
}

// collectFederationModuleMappingCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectFederationModuleMappingCursorValues(res *types.ModuleMapping, cc ...string) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{}

		hasUnique bool

		// All known primary key columns

		pkRel_federation_module bool

		pkRel_compose_module bool

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
				case "rel_federation_module":
					cursor.Set(c, res.FederationModuleID, false)

					pkRel_federation_module = true
				case "rel_compose_module":
					cursor.Set(c, res.ComposeModuleID, false)

					pkRel_compose_module = true

				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !(pkRel_federation_module && pkRel_compose_module && true) {
		collect("rel_federation_module", "rel_compose_module")
	}

	return cursor
}

// checkFederationModuleMappingConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkFederationModuleMappingConstraints(ctx context.Context, res *types.ModuleMapping) error {
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
