package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/automation_sessions.yaml
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

// SearchAutomationSessions returns all matching rows
//
// This function calls convertAutomationSessionFilter with the given
// types.SessionFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchAutomationSessions(ctx context.Context, f types.SessionFilter) (types.SessionSet, types.SessionFilter, error) {
	var (
		err error
		set []*types.Session
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertAutomationSessionFilter(f)
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
		if q, err = setOrderBy(q, sort, s.sortableAutomationSessionColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfAutomationSessions(
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

// fetchFullPageOfAutomationSessions collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfAutomationSessions(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.Session) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.Session, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.Session

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

	set = make([]*types.Session, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryAutomationSessions(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectAutomationSessionCursorValues(set[collected-1], sort...)

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
		prev = s.collectAutomationSessionCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectAutomationSessionCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryAutomationSessions queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryAutomationSessions(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Session) (bool, error),
) ([]*types.Session, error) {
	var (
		set = make([]*types.Session, 0, DefaultSliceCapacity)
		res *types.Session

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalAutomationSessionRowScanner(rows)
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

// LookupAutomationSessionByID searches for session by ID
//
// It returns session even if deleted
func (s Store) LookupAutomationSessionByID(ctx context.Context, id uint64) (*types.Session, error) {
	return s.execLookupAutomationSession(ctx, squirrel.Eq{
		s.preprocessColumn("atms.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateAutomationSession creates one or more rows in automation_sessions table
func (s Store) CreateAutomationSession(ctx context.Context, rr ...*types.Session) (err error) {
	for _, res := range rr {
		err = s.checkAutomationSessionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateAutomationSessions(ctx, s.internalAutomationSessionEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateAutomationSession updates one or more existing rows in automation_sessions
func (s Store) UpdateAutomationSession(ctx context.Context, rr ...*types.Session) error {
	return s.partialAutomationSessionUpdate(ctx, nil, rr...)
}

// partialAutomationSessionUpdate updates one or more existing rows in automation_sessions
func (s Store) partialAutomationSessionUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Session) (err error) {
	for _, res := range rr {
		err = s.checkAutomationSessionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateAutomationSessions(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("atms.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalAutomationSessionEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertAutomationSession updates one or more existing rows in automation_sessions
func (s Store) UpsertAutomationSession(ctx context.Context, rr ...*types.Session) (err error) {
	for _, res := range rr {
		err = s.checkAutomationSessionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertAutomationSessions(ctx, s.internalAutomationSessionEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAutomationSession Deletes one or more rows from automation_sessions table
func (s Store) DeleteAutomationSession(ctx context.Context, rr ...*types.Session) (err error) {
	for _, res := range rr {

		err = s.execDeleteAutomationSessions(ctx, squirrel.Eq{
			s.preprocessColumn("atms.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAutomationSessionByID Deletes row from the automation_sessions table
func (s Store) DeleteAutomationSessionByID(ctx context.Context, ID uint64) error {
	return s.execDeleteAutomationSessions(ctx, squirrel.Eq{
		s.preprocessColumn("atms.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateAutomationSessions Deletes all rows from the automation_sessions table
func (s Store) TruncateAutomationSessions(ctx context.Context) error {
	return s.Truncate(ctx, s.automationSessionTable())
}

// execLookupAutomationSession prepares AutomationSession query and executes it,
// returning types.Session (or error)
func (s Store) execLookupAutomationSession(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Session, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.automationSessionsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalAutomationSessionRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateAutomationSessions updates all matched (by cnd) rows in automation_sessions with given data
func (s Store) execCreateAutomationSessions(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.automationSessionTable()).SetMap(payload))
}

// execUpdateAutomationSessions updates all matched (by cnd) rows in automation_sessions with given data
func (s Store) execUpdateAutomationSessions(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.automationSessionTable("atms")).Where(cnd).SetMap(set))
}

// execUpsertAutomationSessions inserts new or updates matching (by-primary-key) rows in automation_sessions with given data
func (s Store) execUpsertAutomationSessions(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.automationSessionTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteAutomationSessions Deletes all matched (by cnd) rows in automation_sessions with given data
func (s Store) execDeleteAutomationSessions(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.automationSessionTable("atms")).Where(cnd))
}

func (s Store) internalAutomationSessionRowScanner(row rowScanner) (res *types.Session, err error) {
	res = &types.Session{}

	if _, has := s.config.RowScanners["automationSession"]; has {
		scanner := s.config.RowScanners["automationSession"].(func(_ rowScanner, _ *types.Session) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.WorkflowID,
			&res.EventType,
			&res.ResourceType,
			&res.Status,
			&res.Input,
			&res.Output,
			&res.Stacktrace,
			&res.CreatedBy,
			&res.CreatedAt,
			&res.PurgeAt,
			&res.CompletedAt,
			&res.SuspendedAt,
			&res.Error,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan automationSession db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryAutomationSessions returns squirrel.SelectBuilder with set table and all columns
func (s Store) automationSessionsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.automationSessionTable("atms"), s.automationSessionColumns("atms")...)
}

// automationSessionTable name of the db table
func (Store) automationSessionTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "automation_sessions" + alias
}

// AutomationSessionColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) automationSessionColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_workflow",
		alias + "event_type",
		alias + "resource_type",
		alias + "status",
		alias + "input",
		alias + "output",
		alias + "stacktrace",
		alias + "created_by",
		alias + "created_at",
		alias + "purge_at",
		alias + "completed_at",
		alias + "suspended_at",
		alias + "error",
	}
}

// {true true false true true true}

// sortableAutomationSessionColumns returns all AutomationSession columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableAutomationSessionColumns() map[string]string {
	return map[string]string{
		"id": "id",
	}
}

// internalAutomationSessionEncoder encodes fields from types.Session to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeAutomationSession
// func when rdbms.customEncoder=true
func (s Store) internalAutomationSessionEncoder(res *types.Session) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"rel_workflow":  res.WorkflowID,
		"event_type":    res.EventType,
		"resource_type": res.ResourceType,
		"status":        res.Status,
		"input":         res.Input,
		"output":        res.Output,
		"stacktrace":    res.Stacktrace,
		"created_by":    res.CreatedBy,
		"created_at":    res.CreatedAt,
		"purge_at":      res.PurgeAt,
		"completed_at":  res.CompletedAt,
		"suspended_at":  res.SuspendedAt,
		"error":         res.Error,
	}
}

// collectAutomationSessionCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectAutomationSessionCursorValues(res *types.Session, cc ...*filter.SortExpr) *filter.PagingCursor {
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

// checkAutomationSessionConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkAutomationSessionConstraints(ctx context.Context, res *types.Session) error {
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
