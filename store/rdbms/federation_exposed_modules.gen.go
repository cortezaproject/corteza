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
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
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

	return set, f, func() error {
		q, err = s.convertFederationExposedModuleFilter(f)
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
		if q, err = setOrderBy(q, sort, s.sortableFederationExposedModuleColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfFederationExposedModules(
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

// fetchFullPageOfFederationExposedModules collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfFederationExposedModules(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.ExposedModule) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.ExposedModule, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.ExposedModule

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

	set = make([]*types.ExposedModule, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryFederationExposedModules(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectFederationExposedModuleCursorValues(set[collected-1], sort...)

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
		prev = s.collectFederationExposedModuleCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectFederationExposedModuleCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
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
) ([]*types.ExposedModule, error) {
	var (
		set = make([]*types.ExposedModule, 0, DefaultSliceCapacity)
		res *types.ExposedModule

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalFederationExposedModuleRowScanner(rows)
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

// LookupFederationExposedModuleByID searches for federation module by ID
//
// It returns federation module
func (s Store) LookupFederationExposedModuleByID(ctx context.Context, id uint64) (*types.ExposedModule, error) {
	return s.execLookupFederationExposedModule(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.id", ""): store.PreprocessValue(id, ""),
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
	return s.partialFederationExposedModuleUpdate(ctx, nil, rr...)
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
				s.preprocessColumn("cmd.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalFederationExposedModuleEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
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

		err = s.execUpsertFederationExposedModules(ctx, s.internalFederationExposedModuleEncoder(res))
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
			s.preprocessColumn("cmd.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteFederationExposedModuleByID Deletes row from the federation_module_exposed table
func (s Store) DeleteFederationExposedModuleByID(ctx context.Context, ID uint64) error {
	return s.execDeleteFederationExposedModules(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateFederationExposedModules Deletes all rows from the federation_module_exposed table
func (s Store) TruncateFederationExposedModules(ctx context.Context) error {
	return s.Truncate(ctx, s.federationExposedModuleTable())
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
	return s.Exec(ctx, s.InsertBuilder(s.federationExposedModuleTable()).SetMap(payload))
}

// execUpdateFederationExposedModules updates all matched (by cnd) rows in federation_module_exposed with given data
func (s Store) execUpdateFederationExposedModules(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.federationExposedModuleTable("cmd")).Where(cnd).SetMap(set))
}

// execUpsertFederationExposedModules inserts new or updates matching (by-primary-key) rows in federation_module_exposed with given data
func (s Store) execUpsertFederationExposedModules(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.federationExposedModuleTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteFederationExposedModules Deletes all matched (by cnd) rows in federation_module_exposed with given data
func (s Store) execDeleteFederationExposedModules(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.federationExposedModuleTable("cmd")).Where(cnd))
}

func (s Store) internalFederationExposedModuleRowScanner(row rowScanner) (res *types.ExposedModule, err error) {
	res = &types.ExposedModule{}

	if _, has := s.config.RowScanners["federationExposedModule"]; has {
		scanner := s.config.RowScanners["federationExposedModule"].(func(_ rowScanner, _ *types.ExposedModule) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Handle,
			&res.Name,
			&res.NodeID,
			&res.ComposeModuleID,
			&res.ComposeNamespaceID,
			&res.Fields,
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
		return nil, errors.Store("could not scan federationExposedModule db row: %s", err).Wrap(err)
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
		alias + "handle",
		alias + "name",
		alias + "rel_node",
		alias + "rel_compose_module",
		alias + "rel_compose_namespace",
		alias + "fields",
		alias + "created_by",
		alias + "updated_by",
		alias + "deleted_by",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false true true true}

// sortableFederationExposedModuleColumns returns all FederationExposedModule columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableFederationExposedModuleColumns() map[string]string {
	return map[string]string{
		"id": "id",
	}
}

// internalFederationExposedModuleEncoder encodes fields from types.ExposedModule to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeFederationExposedModule
// func when rdbms.customEncoder=true
func (s Store) internalFederationExposedModuleEncoder(res *types.ExposedModule) store.Payload {
	return store.Payload{
		"id":                    res.ID,
		"handle":                res.Handle,
		"name":                  res.Name,
		"rel_node":              res.NodeID,
		"rel_compose_module":    res.ComposeModuleID,
		"rel_compose_namespace": res.ComposeNamespaceID,
		"fields":                res.Fields,
		"created_by":            res.CreatedBy,
		"updated_by":            res.UpdatedBy,
		"deleted_by":            res.DeletedBy,
		"created_at":            res.CreatedAt,
		"updated_at":            res.UpdatedAt,
		"deleted_at":            res.DeletedAt,
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
func (s Store) collectFederationExposedModuleCursorValues(res *types.ExposedModule, cc ...*filter.SortExpr) *filter.PagingCursor {
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
