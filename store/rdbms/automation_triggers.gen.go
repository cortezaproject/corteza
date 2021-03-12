package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/automation_triggers.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
)

var _ = errors.Is

// SearchAutomationTriggers returns all matching rows
//
// This function calls convertAutomationTriggerFilter with the given
// types.TriggerFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchAutomationTriggers(ctx context.Context, f types.TriggerFilter) (types.TriggerSet, types.TriggerFilter, error) {
	var (
		err error
		set []*types.Trigger
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertAutomationTriggerFilter(f)
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
		if q, err = setOrderBy(q, sort, s.sortableAutomationTriggerColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfAutomationTriggers(
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

// fetchFullPageOfAutomationTriggers collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfAutomationTriggers(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.Trigger) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.Trigger, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.Trigger

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

	set = make([]*types.Trigger, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryAutomationTriggers(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectAutomationTriggerCursorValues(set[collected-1], sort...)

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
		prev = s.collectAutomationTriggerCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectAutomationTriggerCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryAutomationTriggers queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryAutomationTriggers(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Trigger) (bool, error),
) ([]*types.Trigger, error) {
	var (
		set = make([]*types.Trigger, 0, DefaultSliceCapacity)
		res *types.Trigger

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalAutomationTriggerRowScanner(rows)
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

// LookupAutomationTriggerByID searches for trigger by ID
//
// It returns trigger even if deleted
func (s Store) LookupAutomationTriggerByID(ctx context.Context, id uint64) (*types.Trigger, error) {
	return s.execLookupAutomationTrigger(ctx, squirrel.Eq{
		s.preprocessColumn("atmt.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateAutomationTrigger creates one or more rows in automation_triggers table
func (s Store) CreateAutomationTrigger(ctx context.Context, rr ...*types.Trigger) (err error) {
	for _, res := range rr {
		err = s.checkAutomationTriggerConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateAutomationTriggers(ctx, s.internalAutomationTriggerEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateAutomationTrigger updates one or more existing rows in automation_triggers
func (s Store) UpdateAutomationTrigger(ctx context.Context, rr ...*types.Trigger) error {
	return s.partialAutomationTriggerUpdate(ctx, nil, rr...)
}

// partialAutomationTriggerUpdate updates one or more existing rows in automation_triggers
func (s Store) partialAutomationTriggerUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Trigger) (err error) {
	for _, res := range rr {
		err = s.checkAutomationTriggerConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateAutomationTriggers(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("atmt.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalAutomationTriggerEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertAutomationTrigger updates one or more existing rows in automation_triggers
func (s Store) UpsertAutomationTrigger(ctx context.Context, rr ...*types.Trigger) (err error) {
	for _, res := range rr {
		err = s.checkAutomationTriggerConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertAutomationTriggers(ctx, s.internalAutomationTriggerEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAutomationTrigger Deletes one or more rows from automation_triggers table
func (s Store) DeleteAutomationTrigger(ctx context.Context, rr ...*types.Trigger) (err error) {
	for _, res := range rr {

		err = s.execDeleteAutomationTriggers(ctx, squirrel.Eq{
			s.preprocessColumn("atmt.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAutomationTriggerByID Deletes row from the automation_triggers table
func (s Store) DeleteAutomationTriggerByID(ctx context.Context, ID uint64) error {
	return s.execDeleteAutomationTriggers(ctx, squirrel.Eq{
		s.preprocessColumn("atmt.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateAutomationTriggers Deletes all rows from the automation_triggers table
func (s Store) TruncateAutomationTriggers(ctx context.Context) error {
	return s.Truncate(ctx, s.automationTriggerTable())
}

// execLookupAutomationTrigger prepares AutomationTrigger query and executes it,
// returning types.Trigger (or error)
func (s Store) execLookupAutomationTrigger(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Trigger, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.automationTriggersSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalAutomationTriggerRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateAutomationTriggers updates all matched (by cnd) rows in automation_triggers with given data
func (s Store) execCreateAutomationTriggers(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.automationTriggerTable()).SetMap(payload))
}

// execUpdateAutomationTriggers updates all matched (by cnd) rows in automation_triggers with given data
func (s Store) execUpdateAutomationTriggers(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.automationTriggerTable("atmt")).Where(cnd).SetMap(set))
}

// execUpsertAutomationTriggers inserts new or updates matching (by-primary-key) rows in automation_triggers with given data
func (s Store) execUpsertAutomationTriggers(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.automationTriggerTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteAutomationTriggers Deletes all matched (by cnd) rows in automation_triggers with given data
func (s Store) execDeleteAutomationTriggers(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.automationTriggerTable("atmt")).Where(cnd))
}

func (s Store) internalAutomationTriggerRowScanner(row rowScanner) (res *types.Trigger, err error) {
	res = &types.Trigger{}

	if _, has := s.config.RowScanners["automationTrigger"]; has {
		scanner := s.config.RowScanners["automationTrigger"].(func(_ rowScanner, _ *types.Trigger) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.WorkflowID,
			&res.StepID,
			&res.Enabled,
			&res.ResourceType,
			&res.EventType,
			&res.Meta,
			&res.Constraints,
			&res.Input,
			&res.OwnedBy,
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
		return nil, errors.Store("could not scan automationTrigger db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryAutomationTriggers returns squirrel.SelectBuilder with set table and all columns
func (s Store) automationTriggersSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.automationTriggerTable("atmt"), s.automationTriggerColumns("atmt")...)
}

// automationTriggerTable name of the db table
func (Store) automationTriggerTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "automation_triggers" + alias
}

// AutomationTriggerColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) automationTriggerColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_workflow",
		alias + "rel_step",
		alias + "enabled",
		alias + "resource_type",
		alias + "event_type",
		alias + "meta",
		alias + "constraints",
		alias + "input",
		alias + "owned_by",
		alias + "created_by",
		alias + "updated_by",
		alias + "deleted_by",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false true true true}

// sortableAutomationTriggerColumns returns all AutomationTrigger columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableAutomationTriggerColumns() map[string]string {
	return map[string]string{
		"id": "id",
	}
}

// internalAutomationTriggerEncoder encodes fields from types.Trigger to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeAutomationTrigger
// func when rdbms.customEncoder=true
func (s Store) internalAutomationTriggerEncoder(res *types.Trigger) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"rel_workflow":  res.WorkflowID,
		"rel_step":      res.StepID,
		"enabled":       res.Enabled,
		"resource_type": res.ResourceType,
		"event_type":    res.EventType,
		"meta":          res.Meta,
		"constraints":   res.Constraints,
		"input":         res.Input,
		"owned_by":      res.OwnedBy,
		"created_by":    res.CreatedBy,
		"updated_by":    res.UpdatedBy,
		"deleted_by":    res.DeletedBy,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}

// collectAutomationTriggerCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectAutomationTriggerCursorValues(res *types.Trigger, cc ...*filter.SortExpr) *filter.PagingCursor {
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

// checkAutomationTriggerConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkAutomationTriggerConstraints(ctx context.Context, res *types.Trigger) error {
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
