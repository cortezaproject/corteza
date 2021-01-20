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
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
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

	return set, f, func() error {
		q, err = s.convertFederationModuleMappingFilter(f)
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
		if f.Sort.Get("rel_federation_module") == nil {
			f.Sort = append(f.Sort, &filter.SortExpr{
				Column:     "rel_federation_module",
				Descending: f.Sort.LastDescending(),
			})
		}
		if f.Sort.Get("rel_compose_module") == nil {
			f.Sort = append(f.Sort, &filter.SortExpr{
				Column:     "rel_compose_module",
				Descending: f.Sort.LastDescending(),
			})
		}
		if f.Sort.Get("rel_compose_namespace") == nil {
			f.Sort = append(f.Sort, &filter.SortExpr{
				Column:     "rel_compose_namespace",
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
		if q, err = setOrderBy(q, sort, s.sortableFederationModuleMappingColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfFederationModuleMappings(
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

// fetchFullPageOfFederationModuleMappings collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfFederationModuleMappings(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.ModuleMapping) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.ModuleMapping, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.ModuleMapping

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

	set = make([]*types.ModuleMapping, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryFederationModuleMappings(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectFederationModuleMappingCursorValues(set[collected-1], sort...)

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
		prev = s.collectFederationModuleMappingCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectFederationModuleMappingCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
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
) ([]*types.ModuleMapping, error) {
	var (
		set = make([]*types.ModuleMapping, 0, DefaultSliceCapacity)
		res *types.ModuleMapping

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalFederationModuleMappingRowScanner(rows)
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

// LookupFederationModuleMappingByFederationModuleIDComposeModuleIDComposeNamespaceID searches for module mapping by federation module id and compose module id
//
// It returns module mapping
func (s Store) LookupFederationModuleMappingByFederationModuleIDComposeModuleIDComposeNamespaceID(ctx context.Context, federation_module_id uint64, compose_module_id uint64, compose_namespace_id uint64) (*types.ModuleMapping, error) {
	return s.execLookupFederationModuleMapping(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.rel_federation_module", ""): store.PreprocessValue(federation_module_id, ""),
		s.preprocessColumn("cmd.rel_compose_module", ""):    store.PreprocessValue(compose_module_id, ""),
		s.preprocessColumn("cmd.rel_compose_namespace", ""): store.PreprocessValue(compose_namespace_id, ""),
	})
}

// LookupFederationModuleMappingByFederationModuleID searches for module mapping by federation module id
//
// It returns module mapping
func (s Store) LookupFederationModuleMappingByFederationModuleID(ctx context.Context, federation_module_id uint64) (*types.ModuleMapping, error) {
	return s.execLookupFederationModuleMapping(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.rel_federation_module", ""): store.PreprocessValue(federation_module_id, ""),
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
	return s.partialFederationModuleMappingUpdate(ctx, nil, rr...)
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
				s.preprocessColumn("cmd.rel_federation_module", ""): store.PreprocessValue(res.FederationModuleID, ""), s.preprocessColumn("cmd.rel_compose_module", ""): store.PreprocessValue(res.ComposeModuleID, ""), s.preprocessColumn("cmd.rel_compose_namespace", ""): store.PreprocessValue(res.ComposeNamespaceID, ""),
			},
			s.internalFederationModuleMappingEncoder(res).Skip("rel_federation_module", "rel_compose_module", "rel_compose_namespace").Only(onlyColumns...))
		if err != nil {
			return err
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

		err = s.execUpsertFederationModuleMappings(ctx, s.internalFederationModuleMappingEncoder(res))
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
			s.preprocessColumn("cmd.rel_federation_module", ""): store.PreprocessValue(res.FederationModuleID, ""), s.preprocessColumn("cmd.rel_compose_module", ""): store.PreprocessValue(res.ComposeModuleID, ""), s.preprocessColumn("cmd.rel_compose_namespace", ""): store.PreprocessValue(res.ComposeNamespaceID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteFederationModuleMappingByFederationModuleIDComposeModuleIDComposeNamespaceID Deletes row from the federation_module_mapping table
func (s Store) DeleteFederationModuleMappingByFederationModuleIDComposeModuleIDComposeNamespaceID(ctx context.Context, federationModuleID uint64, composeModuleID uint64, composeNamespaceID uint64) error {
	return s.execDeleteFederationModuleMappings(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.rel_federation_module", ""): store.PreprocessValue(federationModuleID, ""),
		s.preprocessColumn("cmd.rel_compose_module", ""):    store.PreprocessValue(composeModuleID, ""),
		s.preprocessColumn("cmd.rel_compose_namespace", ""): store.PreprocessValue(composeNamespaceID, ""),
	})
}

// TruncateFederationModuleMappings Deletes all rows from the federation_module_mapping table
func (s Store) TruncateFederationModuleMappings(ctx context.Context) error {
	return s.Truncate(ctx, s.federationModuleMappingTable())
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
	return s.Exec(ctx, s.InsertBuilder(s.federationModuleMappingTable()).SetMap(payload))
}

// execUpdateFederationModuleMappings updates all matched (by cnd) rows in federation_module_mapping with given data
func (s Store) execUpdateFederationModuleMappings(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.federationModuleMappingTable("cmd")).Where(cnd).SetMap(set))
}

// execUpsertFederationModuleMappings inserts new or updates matching (by-primary-key) rows in federation_module_mapping with given data
func (s Store) execUpsertFederationModuleMappings(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.federationModuleMappingTable(),
		set,
		s.preprocessColumn("rel_federation_module", ""),
		s.preprocessColumn("rel_compose_module", ""),
		s.preprocessColumn("rel_compose_namespace", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteFederationModuleMappings Deletes all matched (by cnd) rows in federation_module_mapping with given data
func (s Store) execDeleteFederationModuleMappings(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.federationModuleMappingTable("cmd")).Where(cnd))
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
			&res.ComposeNamespaceID,
			&res.FieldMapping,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan federationModuleMapping db row: %s", err).Wrap(err)
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
		alias + "rel_compose_namespace",
		alias + "field_mapping",
	}
}

// {true true false true true true}

// sortableFederationModuleMappingColumns returns all FederationModuleMapping columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableFederationModuleMappingColumns() map[string]string {
	return map[string]string{
		"rel_federation_module": "rel_federation_module",
		"federationmoduleid":    "rel_federation_module",
		"rel_compose_module":    "rel_compose_module",
		"composemoduleid":       "rel_compose_module",
		"rel_compose_namespace": "rel_compose_namespace",
		"composenamespaceid":    "rel_compose_namespace",
	}
}

// internalFederationModuleMappingEncoder encodes fields from types.ModuleMapping to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeFederationModuleMapping
// func when rdbms.customEncoder=true
func (s Store) internalFederationModuleMappingEncoder(res *types.ModuleMapping) store.Payload {
	return store.Payload{
		"rel_federation_module": res.FederationModuleID,
		"rel_compose_module":    res.ComposeModuleID,
		"rel_compose_namespace": res.ComposeNamespaceID,
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
func (s Store) collectFederationModuleMappingCursorValues(res *types.ModuleMapping, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		// All known primary key columns

		pkRel_federation_module bool

		pkRel_compose_module bool

		pkRel_compose_namespace bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "rel_federation_module":
					cursor.Set(c.Column, res.FederationModuleID, c.Descending)

					pkRel_federation_module = true
				case "rel_compose_module":
					cursor.Set(c.Column, res.ComposeModuleID, c.Descending)

					pkRel_compose_module = true
				case "rel_compose_namespace":
					cursor.Set(c.Column, res.ComposeNamespaceID, c.Descending)

					pkRel_compose_namespace = true

				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !(pkRel_federation_module && pkRel_compose_module && pkRel_compose_namespace && true) {
		collect(&filter.SortExpr{Column: "rel_federation_module", Descending: false}, &filter.SortExpr{Column: "rel_compose_module", Descending: false}, &filter.SortExpr{Column: "rel_compose_namespace", Descending: false})
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
