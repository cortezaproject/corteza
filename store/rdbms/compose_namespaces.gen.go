package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_namespaces.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchComposeNamespaces returns all matching rows
//
// This function calls convertComposeNamespaceFilter with the given
// types.NamespaceFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchComposeNamespaces(ctx context.Context, f types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error) {
	var (
		err error
		set []*types.Namespace
		q   squirrel.SelectBuilder
	)
	q, err = s.convertComposeNamespaceFilter(f)
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
		set, err = s.fetchFullPageOfComposeNamespaces(ctx, q, curSort, f.PageCursor, f.Limit, f.Check)

		if err != nil {
			return err
		}

		if f.Limit > 0 && len(set) > 0 {
			if f.PageCursor != nil && (!f.PageCursor.Reverse || uint(len(set)) == f.Limit) {
				f.PrevPage = s.collectComposeNamespaceCursorValues(set[0], curSort.Columns()...)
				f.PrevPage.Reverse = true
			}

			// Less items fetched then requested by page-limit
			// not very likely there's another page
			f.NextPage = s.collectComposeNamespaceCursorValues(set[len(set)-1], curSort.Columns()...)
		}

		f.PageCursor = nil
		return nil
	}())
}

// fetchFullPageOfComposeNamespaces collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - sorting rules (order by ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn). Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfComposeNamespaces(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*types.Namespace) (bool, error),
) ([]*types.Namespace, error) {
	var (
		set  = make([]*types.Namespace, 0, DefaultSliceCapacity)
		aux  []*types.Namespace
		last *types.Namespace

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
	if q, err = setOrderBy(q, sort, s.sortableComposeNamespaceColumns()...); err != nil {
		return nil, err
	}

	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit))
		}

		if aux, fetched, last, err = s.QueryComposeNamespaces(ctx, tryQuery, check); err != nil {
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
		if cursor = s.collectComposeNamespaceCursorValues(last, sort.Columns()...); cursor == nil {
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

// QueryComposeNamespaces queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryComposeNamespaces(
	ctx context.Context,
	q squirrel.SelectBuilder,
	check func(*types.Namespace) (bool, error),
) ([]*types.Namespace, uint, *types.Namespace, error) {
	var (
		set = make([]*types.Namespace, 0, DefaultSliceCapacity)
		res *types.Namespace

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
			res, err = s.internalComposeNamespaceRowScanner(rows)
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

// LookupComposeNamespaceBySlug searches for namespace by slug (case-insensitive)
func (s Store) LookupComposeNamespaceBySlug(ctx context.Context, slug string) (*types.Namespace, error) {
	return s.execLookupComposeNamespace(ctx, squirrel.Eq{
		s.preprocessColumn("cns.slug", "lower"): s.preprocessValue(slug, "lower"),
	})
}

// LookupComposeNamespaceByID searches for compose namespace by ID
//
// It returns compose namespace even if deleted
func (s Store) LookupComposeNamespaceByID(ctx context.Context, id uint64) (*types.Namespace, error) {
	return s.execLookupComposeNamespace(ctx, squirrel.Eq{
		s.preprocessColumn("cns.id", ""): s.preprocessValue(id, ""),
	})
}

// CreateComposeNamespace creates one or more rows in compose_namespace table
func (s Store) CreateComposeNamespace(ctx context.Context, rr ...*types.Namespace) (err error) {
	for _, res := range rr {
		err = s.checkComposeNamespaceConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateComposeNamespaces(ctx, s.internalComposeNamespaceEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateComposeNamespace updates one or more existing rows in compose_namespace
func (s Store) UpdateComposeNamespace(ctx context.Context, rr ...*types.Namespace) error {
	return s.config.ErrorHandler(s.PartialComposeNamespaceUpdate(ctx, nil, rr...))
}

// PartialComposeNamespaceUpdate updates one or more existing rows in compose_namespace
func (s Store) PartialComposeNamespaceUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Namespace) (err error) {
	for _, res := range rr {
		err = s.checkComposeNamespaceConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateComposeNamespaces(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("cns.id", ""): s.preprocessValue(res.ID, ""),
			},
			s.internalComposeNamespaceEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// UpsertComposeNamespace updates one or more existing rows in compose_namespace
func (s Store) UpsertComposeNamespace(ctx context.Context, rr ...*types.Namespace) (err error) {
	for _, res := range rr {
		err = s.checkComposeNamespaceConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.config.ErrorHandler(s.execUpsertComposeNamespaces(ctx, s.internalComposeNamespaceEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteComposeNamespace Deletes one or more rows from compose_namespace table
func (s Store) DeleteComposeNamespace(ctx context.Context, rr ...*types.Namespace) (err error) {
	for _, res := range rr {

		err = s.execDeleteComposeNamespaces(ctx, squirrel.Eq{
			s.preprocessColumn("cns.id", ""): s.preprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteComposeNamespaceByID Deletes row from the compose_namespace table
func (s Store) DeleteComposeNamespaceByID(ctx context.Context, ID uint64) error {
	return s.execDeleteComposeNamespaces(ctx, squirrel.Eq{
		s.preprocessColumn("cns.id", ""): s.preprocessValue(ID, ""),
	})
}

// TruncateComposeNamespaces Deletes all rows from the compose_namespace table
func (s Store) TruncateComposeNamespaces(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.composeNamespaceTable()))
}

// execLookupComposeNamespace prepares ComposeNamespace query and executes it,
// returning types.Namespace (or error)
func (s Store) execLookupComposeNamespace(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Namespace, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.composeNamespacesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalComposeNamespaceRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateComposeNamespaces updates all matched (by cnd) rows in compose_namespace with given data
func (s Store) execCreateComposeNamespaces(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.composeNamespaceTable()).SetMap(payload)))
}

// execUpdateComposeNamespaces updates all matched (by cnd) rows in compose_namespace with given data
func (s Store) execUpdateComposeNamespaces(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.composeNamespaceTable("cns")).Where(cnd).SetMap(set)))
}

// execUpsertComposeNamespaces inserts new or updates matching (by-primary-key) rows in compose_namespace with given data
func (s Store) execUpsertComposeNamespaces(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.composeNamespaceTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteComposeNamespaces Deletes all matched (by cnd) rows in compose_namespace with given data
func (s Store) execDeleteComposeNamespaces(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.composeNamespaceTable("cns")).Where(cnd)))
}

func (s Store) internalComposeNamespaceRowScanner(row rowScanner) (res *types.Namespace, err error) {
	res = &types.Namespace{}

	if _, has := s.config.RowScanners["composeNamespace"]; has {
		scanner := s.config.RowScanners["composeNamespace"].(func(_ rowScanner, _ *types.Namespace) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Name,
			&res.Slug,
			&res.Enabled,
			&res.Meta,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for ComposeNamespace: %w", err)
	} else {
		return res, nil
	}
}

// QueryComposeNamespaces returns squirrel.SelectBuilder with set table and all columns
func (s Store) composeNamespacesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.composeNamespaceTable("cns"), s.composeNamespaceColumns("cns")...)
}

// composeNamespaceTable name of the db table
func (Store) composeNamespaceTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_namespace" + alias
}

// ComposeNamespaceColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) composeNamespaceColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "name",
		alias + "slug",
		alias + "enabled",
		alias + "meta",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true true true true}

// sortableComposeNamespaceColumns returns all ComposeNamespace columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableComposeNamespaceColumns() []string {
	return []string{
		"id",
		"name",
		"slug",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

// internalComposeNamespaceEncoder encodes fields from types.Namespace to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposeNamespace
// func when rdbms.customEncoder=true
func (s Store) internalComposeNamespaceEncoder(res *types.Namespace) store.Payload {
	return store.Payload{
		"id":         res.ID,
		"name":       res.Name,
		"slug":       res.Slug,
		"enabled":    res.Enabled,
		"meta":       res.Meta,
		"created_at": res.CreatedAt,
		"updated_at": res.UpdatedAt,
		"deleted_at": res.DeletedAt,
	}
}

// collectComposeNamespaceCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectComposeNamespaceCursorValues(res *types.Namespace, cc ...string) *filter.PagingCursor {
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
				case "name":
					cursor.Set(c, res.Name, false)

				case "slug":
					cursor.Set(c, res.Slug, false)
					hasUnique = true

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

func (s *Store) checkComposeNamespaceConstraints(ctx context.Context, res *types.Namespace) error {

	{
		ex, err := s.LookupComposeNamespaceBySlug(ctx, res.Slug)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique
		} else if !errors.Is(err, store.ErrNotFound) {
			return err
		}
	}

	return nil
}
