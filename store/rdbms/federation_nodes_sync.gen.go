package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/federation_nodes_sync.yaml
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

// SearchFederationNodesSyncs returns all matching rows
//
// This function calls convertFederationNodesSyncFilter with the given
// types.NodeSyncFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchFederationNodesSyncs(ctx context.Context, f types.NodeSyncFilter) (types.NodeSyncSet, types.NodeSyncFilter, error) {
	var (
		err error
		set []*types.NodeSync
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertFederationNodesSyncFilter(f)
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
		if f.Sort.Get("rel_node") == nil {
			f.Sort = append(f.Sort, &filter.SortExpr{
				Column:     "rel_node",
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
		if q, err = setOrderBy(q, sort, s.sortableFederationNodesSyncColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfFederationNodesSyncs(
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

// fetchFullPageOfFederationNodesSyncs collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfFederationNodesSyncs(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.NodeSync) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.NodeSync, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.NodeSync

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

	set = make([]*types.NodeSync, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryFederationNodesSyncs(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectFederationNodesSyncCursorValues(set[collected-1], sort...)

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
		prev = s.collectFederationNodesSyncCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectFederationNodesSyncCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryFederationNodesSyncs queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryFederationNodesSyncs(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.NodeSync) (bool, error),
) ([]*types.NodeSync, error) {
	var (
		set = make([]*types.NodeSync, 0, DefaultSliceCapacity)
		res *types.NodeSync

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalFederationNodesSyncRowScanner(rows)
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

// LookupFederationNodesSyncByNodeID searches for sync activity by node ID
//
// It returns sync activity
func (s Store) LookupFederationNodesSyncByNodeID(ctx context.Context, node_id uint64) (*types.NodeSync, error) {
	return s.execLookupFederationNodesSync(ctx, squirrel.Eq{
		s.preprocessColumn("fdns.rel_node", ""): store.PreprocessValue(node_id, ""),
	})
}

// LookupFederationNodesSyncByNodeIDModuleIDSyncTypeSyncStatus searches for activity by node, type and status
//
// It returns sync activity
func (s Store) LookupFederationNodesSyncByNodeIDModuleIDSyncTypeSyncStatus(ctx context.Context, node_id uint64, module_id uint64, sync_type string, sync_status string) (*types.NodeSync, error) {
	return s.execLookupFederationNodesSync(ctx, squirrel.Eq{
		s.preprocessColumn("fdns.rel_node", ""):    store.PreprocessValue(node_id, ""),
		s.preprocessColumn("fdns.rel_module", ""):  store.PreprocessValue(module_id, ""),
		s.preprocessColumn("fdns.sync_type", ""):   store.PreprocessValue(sync_type, ""),
		s.preprocessColumn("fdns.sync_status", ""): store.PreprocessValue(sync_status, ""),
	})
}

// CreateFederationNodesSync creates one or more rows in federation_nodes_sync table
func (s Store) CreateFederationNodesSync(ctx context.Context, rr ...*types.NodeSync) (err error) {
	for _, res := range rr {
		err = s.checkFederationNodesSyncConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateFederationNodesSyncs(ctx, s.internalFederationNodesSyncEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateFederationNodesSync updates one or more existing rows in federation_nodes_sync
func (s Store) UpdateFederationNodesSync(ctx context.Context, rr ...*types.NodeSync) error {
	return s.partialFederationNodesSyncUpdate(ctx, nil, rr...)
}

// partialFederationNodesSyncUpdate updates one or more existing rows in federation_nodes_sync
func (s Store) partialFederationNodesSyncUpdate(ctx context.Context, onlyColumns []string, rr ...*types.NodeSync) (err error) {
	for _, res := range rr {
		err = s.checkFederationNodesSyncConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateFederationNodesSyncs(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("fdns.rel_node", ""): store.PreprocessValue(res.NodeID, ""),
			},
			s.internalFederationNodesSyncEncoder(res).Skip("rel_node").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertFederationNodesSync updates one or more existing rows in federation_nodes_sync
func (s Store) UpsertFederationNodesSync(ctx context.Context, rr ...*types.NodeSync) (err error) {
	for _, res := range rr {
		err = s.checkFederationNodesSyncConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertFederationNodesSyncs(ctx, s.internalFederationNodesSyncEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteFederationNodesSync Deletes one or more rows from federation_nodes_sync table
func (s Store) DeleteFederationNodesSync(ctx context.Context, rr ...*types.NodeSync) (err error) {
	for _, res := range rr {

		err = s.execDeleteFederationNodesSyncs(ctx, squirrel.Eq{
			s.preprocessColumn("fdns.rel_node", ""): store.PreprocessValue(res.NodeID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteFederationNodesSyncByNodeID Deletes row from the federation_nodes_sync table
func (s Store) DeleteFederationNodesSyncByNodeID(ctx context.Context, nodeID uint64) error {
	return s.execDeleteFederationNodesSyncs(ctx, squirrel.Eq{
		s.preprocessColumn("fdns.rel_node", ""): store.PreprocessValue(nodeID, ""),
	})
}

// TruncateFederationNodesSyncs Deletes all rows from the federation_nodes_sync table
func (s Store) TruncateFederationNodesSyncs(ctx context.Context) error {
	return s.Truncate(ctx, s.federationNodesSyncTable())
}

// execLookupFederationNodesSync prepares FederationNodesSync query and executes it,
// returning types.NodeSync (or error)
func (s Store) execLookupFederationNodesSync(ctx context.Context, cnd squirrel.Sqlizer) (res *types.NodeSync, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.federationNodesSyncsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalFederationNodesSyncRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateFederationNodesSyncs updates all matched (by cnd) rows in federation_nodes_sync with given data
func (s Store) execCreateFederationNodesSyncs(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.federationNodesSyncTable()).SetMap(payload))
}

// execUpdateFederationNodesSyncs updates all matched (by cnd) rows in federation_nodes_sync with given data
func (s Store) execUpdateFederationNodesSyncs(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.federationNodesSyncTable("fdns")).Where(cnd).SetMap(set))
}

// execUpsertFederationNodesSyncs inserts new or updates matching (by-primary-key) rows in federation_nodes_sync with given data
func (s Store) execUpsertFederationNodesSyncs(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.federationNodesSyncTable(),
		set,
		s.preprocessColumn("rel_node", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteFederationNodesSyncs Deletes all matched (by cnd) rows in federation_nodes_sync with given data
func (s Store) execDeleteFederationNodesSyncs(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.federationNodesSyncTable("fdns")).Where(cnd))
}

func (s Store) internalFederationNodesSyncRowScanner(row rowScanner) (res *types.NodeSync, err error) {
	res = &types.NodeSync{}

	if _, has := s.config.RowScanners["federationNodesSync"]; has {
		scanner := s.config.RowScanners["federationNodesSync"].(func(_ rowScanner, _ *types.NodeSync) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.NodeID,
			&res.ModuleID,
			&res.SyncType,
			&res.SyncStatus,
			&res.TimeOfAction,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan federationNodesSync db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryFederationNodesSyncs returns squirrel.SelectBuilder with set table and all columns
func (s Store) federationNodesSyncsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.federationNodesSyncTable("fdns"), s.federationNodesSyncColumns("fdns")...)
}

// federationNodesSyncTable name of the db table
func (Store) federationNodesSyncTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "federation_nodes_sync" + alias
}

// FederationNodesSyncColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) federationNodesSyncColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "rel_node",
		alias + "rel_module",
		alias + "sync_type",
		alias + "sync_status",
		alias + "time_action",
	}
}

// {true true false true true true}

// sortableFederationNodesSyncColumns returns all FederationNodesSync columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableFederationNodesSyncColumns() map[string]string {
	return map[string]string{
		"rel_node":     "rel_node",
		"nodeid":       "rel_node",
		"time_action":  "time_action",
		"timeofaction": "time_action",
	}
}

// internalFederationNodesSyncEncoder encodes fields from types.NodeSync to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeFederationNodesSync
// func when rdbms.customEncoder=true
func (s Store) internalFederationNodesSyncEncoder(res *types.NodeSync) store.Payload {
	return store.Payload{
		"rel_node":    res.NodeID,
		"rel_module":  res.ModuleID,
		"sync_type":   res.SyncType,
		"sync_status": res.SyncStatus,
		"time_action": res.TimeOfAction,
	}
}

// collectFederationNodesSyncCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectFederationNodesSyncCursorValues(res *types.NodeSync, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		// All known primary key columns

		pkRel_node bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "rel_node":
					cursor.Set(c.Column, res.NodeID, c.Descending)

					pkRel_node = true
				case "time_action":
					cursor.Set(c.Column, res.TimeOfAction, c.Descending)

				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !(pkRel_node && true) {
		collect(&filter.SortExpr{Column: "rel_node", Descending: false})
	}

	return cursor
}

// checkFederationNodesSyncConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkFederationNodesSyncConstraints(ctx context.Context, res *types.NodeSync) error {
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
