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
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
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

	return set, f, func() error {
		q, err = s.convertComposePageFilter(f)
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
		if q, err = setOrderBy(q, sort, s.sortableComposePageColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfComposePages(
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

// fetchFullPageOfComposePages collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfComposePages(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.Page) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.Page, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.Page

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

	set = make([]*types.Page, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryComposePages(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectComposePageCursorValues(set[collected-1], sort...)

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
		prev = s.collectComposePageCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectComposePageCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
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
) ([]*types.Page, error) {
	var (
		set = make([]*types.Page, 0, DefaultSliceCapacity)
		res *types.Page

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalComposePageRowScanner(rows)
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
		s.preprocessColumn("id", ""),
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
		return nil, errors.Store("could not scan composePage db row: %s", err).Wrap(err)
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

// {true true false true true true}

// sortableComposePageColumns returns all ComposePage columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableComposePageColumns() map[string]string {
	return map[string]string{
		"id": "id", "weight": "weight", "created_at": "created_at",
		"createdat":  "created_at",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
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
func (s Store) collectComposePageCursorValues(res *types.Page, cc ...*filter.SortExpr) *filter.PagingCursor {
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
