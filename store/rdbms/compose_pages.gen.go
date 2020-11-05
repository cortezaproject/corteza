package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_pages.yaml
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

// SearchComposePages returns all matching rows
//
// This function calls convertComposePageFilter with the given
// types.PageFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchComposePages(ctx context.Context, f types.PageFilter) (types.PageSet, types.PageFilter, error) {
	var (
		err error
		set []*types.Page
		q   squirrel.SelectBuilder
	)
	q, err = s.convertComposePageFilter(f)
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
		set, err = s.fetchFullPageOfComposePages(ctx, q, curSort, f.PageCursor, f.Limit, f.Check)

		if err != nil {
			return err
		}

		if f.Limit > 0 && len(set) > 0 {
			if f.PageCursor != nil && (!f.PageCursor.Reverse || uint(len(set)) == f.Limit) {
				f.PrevPage = s.collectComposePageCursorValues(set[0], curSort.Columns()...)
				f.PrevPage.Reverse = true
			}

			// Less items fetched then requested by page-limit
			// not very likely there's another page
			f.NextPage = s.collectComposePageCursorValues(set[len(set)-1], curSort.Columns()...)
		}

		f.PageCursor = nil
		return nil
	}()
}

// fetchFullPageOfComposePages collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - sorting rules (order by ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn). Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfComposePages(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*types.Page) (bool, error),
) ([]*types.Page, error) {
	var (
		set  = make([]*types.Page, 0, DefaultSliceCapacity)
		aux  []*types.Page
		last *types.Page

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
	if q, err = setOrderBy(q, sort, s.sortableComposePageColumns()...); err != nil {
		return nil, err
	}

	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit))
		}

		if aux, fetched, last, err = s.QueryComposePages(ctx, tryQuery, check); err != nil {
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
		if cursor = s.collectComposePageCursorValues(last, sort.Columns()...); cursor == nil {
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

// QueryComposePages queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryComposePages(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Page) (bool, error),
) ([]*types.Page, uint, *types.Page, error) {
	var (
		set = make([]*types.Page, 0, DefaultSliceCapacity)
		res *types.Page

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
			res, err = s.internalComposePageRowScanner(rows)
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

// LookupComposePageByNamespaceIDHandle searches for page by handle (case-insensitive)
func (s Store) LookupComposePageByNamespaceIDHandle(ctx context.Context, namespace_id uint64, handle string) (*types.Page, error) {
	return s.execLookupComposePage(ctx, squirrel.Eq{
		s.preprocessColumn("cpg.rel_namespace", ""): store.PreprocessValue(namespace_id, ""),
		s.preprocessColumn("cpg.handle", "lower"):   store.PreprocessValue(handle, "lower"),

		"cpg.deleted_at": nil,
	})
}

// LookupComposePageByNamespaceIDModuleID searches for page by moduleID
func (s Store) LookupComposePageByNamespaceIDModuleID(ctx context.Context, namespace_id uint64, module_id uint64) (*types.Page, error) {
	return s.execLookupComposePage(ctx, squirrel.Eq{
		s.preprocessColumn("cpg.rel_namespace", ""): store.PreprocessValue(namespace_id, ""),
		s.preprocessColumn("cpg.rel_module", ""):    store.PreprocessValue(module_id, ""),

		"cpg.deleted_at": nil,
	})
}

// LookupComposePageByID searches for compose page by ID
//
// It returns compose page even if deleted
func (s Store) LookupComposePageByID(ctx context.Context, id uint64) (*types.Page, error) {
	return s.execLookupComposePage(ctx, squirrel.Eq{
		s.preprocessColumn("cpg.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateComposePage creates one or more rows in compose_page table
func (s Store) CreateComposePage(ctx context.Context, rr ...*types.Page) (err error) {
	for _, res := range rr {
		err = s.checkComposePageConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateComposePages(ctx, s.internalComposePageEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateComposePage updates one or more existing rows in compose_page
func (s Store) UpdateComposePage(ctx context.Context, rr ...*types.Page) error {
	return s.partialComposePageUpdate(ctx, nil, rr...)
}

// partialComposePageUpdate updates one or more existing rows in compose_page
func (s Store) partialComposePageUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Page) (err error) {
	for _, res := range rr {
		err = s.checkComposePageConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateComposePages(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("cpg.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalComposePageEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertComposePage updates one or more existing rows in compose_page
func (s Store) UpsertComposePage(ctx context.Context, rr ...*types.Page) (err error) {
	for _, res := range rr {
		err = s.checkComposePageConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertComposePages(ctx, s.internalComposePageEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteComposePage Deletes one or more rows from compose_page table
func (s Store) DeleteComposePage(ctx context.Context, rr ...*types.Page) (err error) {
	for _, res := range rr {

		err = s.execDeleteComposePages(ctx, squirrel.Eq{
			s.preprocessColumn("cpg.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteComposePageByID Deletes row from the compose_page table
func (s Store) DeleteComposePageByID(ctx context.Context, ID uint64) error {
	return s.execDeleteComposePages(ctx, squirrel.Eq{
		s.preprocessColumn("cpg.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateComposePages Deletes all rows from the compose_page table
func (s Store) TruncateComposePages(ctx context.Context) error {
	return s.Truncate(ctx, s.composePageTable())
}

// execLookupComposePage prepares ComposePage query and executes it,
// returning types.Page (or error)
func (s Store) execLookupComposePage(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Page, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.composePagesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalComposePageRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateComposePages updates all matched (by cnd) rows in compose_page with given data
func (s Store) execCreateComposePages(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.composePageTable()).SetMap(payload))
}

// execUpdateComposePages updates all matched (by cnd) rows in compose_page with given data
func (s Store) execUpdateComposePages(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.composePageTable("cpg")).Where(cnd).SetMap(set))
}

// execUpsertComposePages inserts new or updates matching (by-primary-key) rows in compose_page with given data
func (s Store) execUpsertComposePages(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.composePageTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteComposePages Deletes all matched (by cnd) rows in compose_page with given data
func (s Store) execDeleteComposePages(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.composePageTable("cpg")).Where(cnd))
}

func (s Store) internalComposePageRowScanner(row rowScanner) (res *types.Page, err error) {
	res = &types.Page{}

	if _, has := s.config.RowScanners["composePage"]; has {
		scanner := s.config.RowScanners["composePage"].(func(_ rowScanner, _ *types.Page) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.SelfID,
			&res.NamespaceID,
			&res.ModuleID,
			&res.Handle,
			&res.Title,
			&res.Description,
			&res.Blocks,
			&res.Visible,
			&res.Weight,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan composePage db row").Wrap(err)
	} else {
		return res, nil
	}
}

// QueryComposePages returns squirrel.SelectBuilder with set table and all columns
func (s Store) composePagesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.composePageTable("cpg"), s.composePageColumns("cpg")...)
}

// composePageTable name of the db table
func (Store) composePageTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_page" + alias
}

// ComposePageColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) composePageColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "self_id",
		alias + "rel_namespace",
		alias + "rel_module",
		alias + "handle",
		alias + "title",
		alias + "description",
		alias + "blocks",
		alias + "visible",
		alias + "weight",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true true true true}

// sortableComposePageColumns returns all ComposePage columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableComposePageColumns() []string {
	return []string{
		"id",
		"weight",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

// internalComposePageEncoder encodes fields from types.Page to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposePage
// func when rdbms.customEncoder=true
func (s Store) internalComposePageEncoder(res *types.Page) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"self_id":       res.SelfID,
		"rel_namespace": res.NamespaceID,
		"rel_module":    res.ModuleID,
		"handle":        res.Handle,
		"title":         res.Title,
		"description":   res.Description,
		"blocks":        res.Blocks,
		"visible":       res.Visible,
		"weight":        res.Weight,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}

// collectComposePageCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectComposePageCursorValues(res *types.Page, cc ...string) *filter.PagingCursor {
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
				case "weight":
					cursor.Set(c, res.Weight, false)

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

// checkComposePageConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkComposePageConstraints(ctx context.Context, res *types.Page) error {
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
