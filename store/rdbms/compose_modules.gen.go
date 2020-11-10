package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_modules.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchComposeModules returns all matching rows
//
// This function calls convertComposeModuleFilter with the given
// types.ModuleFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchComposeModules(ctx context.Context, f types.ModuleFilter) (types.ModuleSet, types.ModuleFilter, error) {
	var (
		err error
		set []*types.Module
		q   squirrel.SelectBuilder
	)
	q, err = s.convertComposeModuleFilter(f)
	if err != nil {
		return nil, f, err
	}

	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reversedCursor := f.PageCursor != nil && f.PageCursor.Reverse

	// Sorting and paging are both enabled in definition yaml file
	// {search: {enableSorting:true, enablePaging:true}}
	curSort := f.Sort.Clone()

	// If paging with reverse cursor, change the sorting
	// direction for all columns we're sorting by
	if reversedCursor {
		curSort.Reverse()
	}

	return set, f, func() error {
		set, err = s.fetchFullPageOfComposeModules(ctx, q, curSort, f.PageCursor, f.Limit, f.Check)

		if err != nil {
			return err
		}

		if f.Limit > 0 && len(set) > 0 {
			if f.PageCursor != nil && (!f.PageCursor.Reverse || uint(len(set)) == f.Limit) {
				f.PrevPage = s.collectComposeModuleCursorValues(set[0], curSort.Columns()...)
				f.PrevPage.Reverse = true
			}

			// Less items fetched then requested by page-limit
			// not very likely there's another page
			f.NextPage = s.collectComposeModuleCursorValues(set[len(set)-1], curSort.Columns()...)
		}

		f.PageCursor = nil
		return nil
	}()
}

// fetchFullPageOfComposeModules collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - sorting rules (order by ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn). Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfComposeModules(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*types.Module) (bool, error),
) ([]*types.Module, error) {
	var (
		set  = make([]*types.Module, 0, DefaultSliceCapacity)
		aux  []*types.Module
		last *types.Module

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
	if q, err = setOrderBy(q, sort, s.sortableComposeModuleColumns()); err != nil {
		return nil, err
	}

	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit))
		}

		if aux, fetched, last, err = s.QueryComposeModules(ctx, tryQuery, check); err != nil {
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
		if cursor = s.collectComposeModuleCursorValues(last, sort.Columns()...); cursor == nil {
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

// QueryComposeModules queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryComposeModules(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Module) (bool, error),
) ([]*types.Module, uint, *types.Module, error) {
	var (
		set = make([]*types.Module, 0, DefaultSliceCapacity)
		res *types.Module

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
			res, err = s.internalComposeModuleRowScanner(rows)
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

// LookupComposeModuleByNamespaceIDHandle searches for compose module by handle (case-insensitive)
func (s Store) LookupComposeModuleByNamespaceIDHandle(ctx context.Context, namespace_id uint64, handle string) (*types.Module, error) {
	return s.execLookupComposeModule(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.rel_namespace", ""): store.PreprocessValue(namespace_id, ""),
		s.preprocessColumn("cmd.handle", "lower"):   store.PreprocessValue(handle, "lower"),

		"cmd.deleted_at": nil,
	})
}

// LookupComposeModuleByNamespaceIDName searches for compose module by name (case-insensitive)
func (s Store) LookupComposeModuleByNamespaceIDName(ctx context.Context, namespace_id uint64, name string) (*types.Module, error) {
	return s.execLookupComposeModule(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.rel_namespace", ""): store.PreprocessValue(namespace_id, ""),
		s.preprocessColumn("cmd.name", "lower"):     store.PreprocessValue(name, "lower"),

		"cmd.deleted_at": nil,
	})
}

// LookupComposeModuleByID searches for compose module by ID
//
// It returns compose module even if deleted
func (s Store) LookupComposeModuleByID(ctx context.Context, id uint64) (*types.Module, error) {
	return s.execLookupComposeModule(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateComposeModule creates one or more rows in compose_module table
func (s Store) CreateComposeModule(ctx context.Context, rr ...*types.Module) (err error) {
	for _, res := range rr {
		err = s.checkComposeModuleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateComposeModules(ctx, s.internalComposeModuleEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateComposeModule updates one or more existing rows in compose_module
func (s Store) UpdateComposeModule(ctx context.Context, rr ...*types.Module) error {
	return s.partialComposeModuleUpdate(ctx, nil, rr...)
}

// partialComposeModuleUpdate updates one or more existing rows in compose_module
func (s Store) partialComposeModuleUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Module) (err error) {
	for _, res := range rr {
		err = s.checkComposeModuleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateComposeModules(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("cmd.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalComposeModuleEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertComposeModule updates one or more existing rows in compose_module
func (s Store) UpsertComposeModule(ctx context.Context, rr ...*types.Module) (err error) {
	for _, res := range rr {
		err = s.checkComposeModuleConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertComposeModules(ctx, s.internalComposeModuleEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteComposeModule Deletes one or more rows from compose_module table
func (s Store) DeleteComposeModule(ctx context.Context, rr ...*types.Module) (err error) {
	for _, res := range rr {

		err = s.execDeleteComposeModules(ctx, squirrel.Eq{
			s.preprocessColumn("cmd.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteComposeModuleByID Deletes row from the compose_module table
func (s Store) DeleteComposeModuleByID(ctx context.Context, ID uint64) error {
	return s.execDeleteComposeModules(ctx, squirrel.Eq{
		s.preprocessColumn("cmd.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateComposeModules Deletes all rows from the compose_module table
func (s Store) TruncateComposeModules(ctx context.Context) error {
	return s.Truncate(ctx, s.composeModuleTable())
}

// execLookupComposeModule prepares ComposeModule query and executes it,
// returning types.Module (or error)
func (s Store) execLookupComposeModule(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Module, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.composeModulesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalComposeModuleRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateComposeModules updates all matched (by cnd) rows in compose_module with given data
func (s Store) execCreateComposeModules(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.composeModuleTable()).SetMap(payload))
}

// execUpdateComposeModules updates all matched (by cnd) rows in compose_module with given data
func (s Store) execUpdateComposeModules(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.composeModuleTable("cmd")).Where(cnd).SetMap(set))
}

// execUpsertComposeModules inserts new or updates matching (by-primary-key) rows in compose_module with given data
func (s Store) execUpsertComposeModules(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.composeModuleTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteComposeModules Deletes all matched (by cnd) rows in compose_module with given data
func (s Store) execDeleteComposeModules(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.composeModuleTable("cmd")).Where(cnd))
}

func (s Store) internalComposeModuleRowScanner(row rowScanner) (res *types.Module, err error) {
	res = &types.Module{}

	if _, has := s.config.RowScanners["composeModule"]; has {
		scanner := s.config.RowScanners["composeModule"].(func(_ rowScanner, _ *types.Module) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Handle,
			&res.Name,
			&res.Meta,
			&res.NamespaceID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan composeModule db row").Wrap(err)
	} else {
		return res, nil
	}
}

// QueryComposeModules returns squirrel.SelectBuilder with set table and all columns
func (s Store) composeModulesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.composeModuleTable("cmd"), s.composeModuleColumns("cmd")...)
}

// composeModuleTable name of the db table
func (Store) composeModuleTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_module" + alias
}

// ComposeModuleColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) composeModuleColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "handle",
		alias + "name",
		alias + "meta",
		alias + "rel_namespace",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true true true true}

// sortableComposeModuleColumns returns all ComposeModule columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableComposeModuleColumns() map[string]string {
	return map[string]string{
		"id": "id", "handle": "handle", "name": "name", "created_at": "created_at",
		"createdat":  "created_at",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
	}
}

// internalComposeModuleEncoder encodes fields from types.Module to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposeModule
// func when rdbms.customEncoder=true
func (s Store) internalComposeModuleEncoder(res *types.Module) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"handle":        res.Handle,
		"name":          res.Name,
		"meta":          res.Meta,
		"rel_namespace": res.NamespaceID,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}

// collectComposeModuleCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectComposeModuleCursorValues(res *types.Module, cc ...string) *filter.PagingCursor {
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
				case "handle":
					cursor.Set(c, res.Handle, false)
					hasUnique = true

				case "name":
					cursor.Set(c, res.Name, false)

				case "created_at":
					cursor.Set(c, res.CreatedAt, false)

				case "updated_at":
					cursor.Set(c, res.UpdatedAt, false)

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

// checkComposeModuleConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkComposeModuleConstraints(ctx context.Context, res *types.Module) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	valid = valid && res.NamespaceID > 0

	valid = valid && len(res.Handle) > 0

	if !valid {
		return nil
	}

	{
		ex, err := s.LookupComposeModuleByNamespaceIDHandle(ctx, res.NamespaceID, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
