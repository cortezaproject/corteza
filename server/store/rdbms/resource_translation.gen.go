package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/resource_translation.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = errors.Is

// SearchResourceTranslations returns all matching rows
//
// This function calls convertResourceTranslationFilter with the given
// types.ResourceTranslationFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchResourceTranslations(ctx context.Context, f types.ResourceTranslationFilter) (types.ResourceTranslationSet, types.ResourceTranslationFilter, error) {
	var (
		err error
		set []*types.ResourceTranslation
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertResourceTranslationFilter(f)
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
		if q, err = setOrderBy(q, sort, s.sortableResourceTranslationColumns(), s.Config().SqlSortHandler); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfResourceTranslations(
			ctx,
			q, f.Sort, f.PageCursor,
			f.Limit,
			nil,
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

// fetchFullPageOfResourceTranslations collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfResourceTranslations(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.ResourceTranslation) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.ResourceTranslation, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.ResourceTranslation

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

	set = make([]*types.ResourceTranslation, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryResourceTranslations(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectResourceTranslationCursorValues(set[collected-1], sort...)

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
		prev = s.collectResourceTranslationCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectResourceTranslationCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryResourceTranslations queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryResourceTranslations(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.ResourceTranslation) (bool, error),
) ([]*types.ResourceTranslation, error) {
	var (
		tmp = make([]*types.ResourceTranslation, 0, DefaultSliceCapacity)
		set = make([]*types.ResourceTranslation, 0, DefaultSliceCapacity)
		res *types.ResourceTranslation

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalResourceTranslationRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		tmp = append(tmp, res)
	}

	for _, res = range tmp {

		set = append(set, res)
	}

	return set, nil
}

// LookupResourceTranslationByID searches for resource translation by ID
// It also returns deleted resource translations.
func (s Store) LookupResourceTranslationByID(ctx context.Context, id uint64) (*types.ResourceTranslation, error) {
	return s.execLookupResourceTranslation(ctx, squirrel.Eq{
		s.preprocessColumn("rt.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateResourceTranslation creates one or more rows in resource_translations table
func (s Store) CreateResourceTranslation(ctx context.Context, rr ...*types.ResourceTranslation) (err error) {
	for _, res := range rr {
		err = s.checkResourceTranslationConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateResourceTranslations(ctx, s.internalResourceTranslationEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateResourceTranslation updates one or more existing rows in resource_translations
func (s Store) UpdateResourceTranslation(ctx context.Context, rr ...*types.ResourceTranslation) error {
	return s.partialResourceTranslationUpdate(ctx, nil, rr...)
}

// partialResourceTranslationUpdate updates one or more existing rows in resource_translations
func (s Store) partialResourceTranslationUpdate(ctx context.Context, onlyColumns []string, rr ...*types.ResourceTranslation) (err error) {
	for _, res := range rr {
		err = s.checkResourceTranslationConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateResourceTranslations(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("rt.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalResourceTranslationEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertResourceTranslation updates one or more existing rows in resource_translations
func (s Store) UpsertResourceTranslation(ctx context.Context, rr ...*types.ResourceTranslation) (err error) {
	for _, res := range rr {
		err = s.checkResourceTranslationConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertResourceTranslations(ctx, s.internalResourceTranslationEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteResourceTranslation Deletes one or more rows from resource_translations table
func (s Store) DeleteResourceTranslation(ctx context.Context, rr ...*types.ResourceTranslation) (err error) {
	for _, res := range rr {

		err = s.execDeleteResourceTranslations(ctx, squirrel.Eq{
			s.preprocessColumn("rt.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteResourceTranslationByID Deletes row from the resource_translations table
func (s Store) DeleteResourceTranslationByID(ctx context.Context, ID uint64) error {
	return s.execDeleteResourceTranslations(ctx, squirrel.Eq{
		s.preprocessColumn("rt.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateResourceTranslations Deletes all rows from the resource_translations table
func (s Store) TruncateResourceTranslations(ctx context.Context) error {
	return s.Truncate(ctx, s.resourceTranslationTable())
}

// execLookupResourceTranslation prepares ResourceTranslation query and executes it,
// returning types.ResourceTranslation (or error)
func (s Store) execLookupResourceTranslation(ctx context.Context, cnd squirrel.Sqlizer) (res *types.ResourceTranslation, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.resourceTranslationsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalResourceTranslationRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateResourceTranslations updates all matched (by cnd) rows in resource_translations with given data
func (s Store) execCreateResourceTranslations(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.resourceTranslationTable()).SetMap(payload))
}

// execUpdateResourceTranslations updates all matched (by cnd) rows in resource_translations with given data
func (s Store) execUpdateResourceTranslations(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.resourceTranslationTable("rt")).Where(cnd).SetMap(set))
}

// execUpsertResourceTranslations inserts new or updates matching (by-primary-key) rows in resource_translations with given data
func (s Store) execUpsertResourceTranslations(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.resourceTranslationTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteResourceTranslations Deletes all matched (by cnd) rows in resource_translations with given data
func (s Store) execDeleteResourceTranslations(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.resourceTranslationTable("rt")).Where(cnd))
}

func (s Store) internalResourceTranslationRowScanner(row rowScanner) (res *types.ResourceTranslation, err error) {
	res = &types.ResourceTranslation{}

	err = row.Scan(
		&res.ID,
		&res.Lang,
		&res.Resource,
		&res.K,
		&res.Message,
		&res.OwnedBy,
		&res.CreatedBy,
		&res.UpdatedBy,
		&res.DeletedBy,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan resourceTranslation db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryResourceTranslations returns squirrel.SelectBuilder with set table and all columns
func (s Store) resourceTranslationsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.resourceTranslationTable("rt"), s.resourceTranslationColumns("rt")...)
}

// resourceTranslationTable name of the db table
func (Store) resourceTranslationTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "resource_translations" + alias
}

// ResourceTranslationColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) resourceTranslationColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "lang",
		alias + "resource",
		alias + "k",
		alias + "message",
		alias + "owned_by",
		alias + "created_by",
		alias + "updated_by",
		alias + "deleted_by",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false true true false}

// sortableResourceTranslationColumns returns all ResourceTranslation columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableResourceTranslationColumns() map[string]string {
	return map[string]string{
		"id": "id", "owned_by": "owned_by",
		"ownedby":    "owned_by",
		"created_by": "created_by",
		"createdby":  "created_by",
		"updated_by": "updated_by",
		"updatedby":  "updated_by",
		"deleted_by": "deleted_by",
		"deletedby":  "deleted_by",
		"created_at": "created_at",
		"createdat":  "created_at",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
	}
}

// internalResourceTranslationEncoder encodes fields from types.ResourceTranslation to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeResourceTranslation
// func when rdbms.customEncoder=true
func (s Store) internalResourceTranslationEncoder(res *types.ResourceTranslation) store.Payload {
	return store.Payload{
		"id":         res.ID,
		"lang":       res.Lang,
		"resource":   res.Resource,
		"k":          res.K,
		"message":    res.Message,
		"owned_by":   res.OwnedBy,
		"created_by": res.CreatedBy,
		"updated_by": res.UpdatedBy,
		"deleted_by": res.DeletedBy,
		"created_at": res.CreatedAt,
		"updated_at": res.UpdatedAt,
		"deleted_at": res.DeletedAt,
	}
}

// collectResourceTranslationCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectResourceTranslationCursorValues(res *types.ResourceTranslation, cc ...*filter.SortExpr) *filter.PagingCursor {
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
				case "owned_by":
					cursor.Set(c.Column, res.OwnedBy, c.Descending)

				case "created_by":
					cursor.Set(c.Column, res.CreatedBy, c.Descending)

				case "updated_by":
					cursor.Set(c.Column, res.UpdatedBy, c.Descending)

				case "deleted_by":
					cursor.Set(c.Column, res.DeletedBy, c.Descending)

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

// checkResourceTranslationConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkResourceTranslationConstraints(ctx context.Context, res *types.ResourceTranslation) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	if !valid {
		return nil
	}

	var checks = make([]func() error, 0)

	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}

	return nil
}
