package rdbms

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	automationType "github.com/cortezaproject/corteza-server/automation/types"
	composeType "github.com/cortezaproject/corteza-server/compose/types"
	federationType "github.com/cortezaproject/corteza-server/federation/types"
	actionlogType "github.com/cortezaproject/corteza-server/pkg/actionlog"
	discoveryType "github.com/cortezaproject/corteza-server/pkg/discovery/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	flagType "github.com/cortezaproject/corteza-server/pkg/flag/types"
	labelsType "github.com/cortezaproject/corteza-server/pkg/label/types"
	rbacType "github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	systemType "github.com/cortezaproject/corteza-server/system/types"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	_ store.Actionlogs                 = &Store{}
	_ store.ApigwFilters               = &Store{}
	_ store.ApigwRoutes                = &Store{}
	_ store.Applications               = &Store{}
	_ store.Attachments                = &Store{}
	_ store.AuthClients                = &Store{}
	_ store.AuthConfirmedClients       = &Store{}
	_ store.AuthOa2tokens              = &Store{}
	_ store.AuthSessions               = &Store{}
	_ store.AutomationSessions         = &Store{}
	_ store.AutomationTriggers         = &Store{}
	_ store.AutomationWorkflows        = &Store{}
	_ store.ComposeAttachments         = &Store{}
	_ store.ComposeCharts              = &Store{}
	_ store.ComposeModules             = &Store{}
	_ store.ComposeModuleFields        = &Store{}
	_ store.ComposeNamespaces          = &Store{}
	_ store.ComposePages               = &Store{}
	_ store.Credentials                = &Store{}
	_ store.DalConnections             = &Store{}
	_ store.DalSensitivityLevels       = &Store{}
	_ store.DataPrivacyRequests        = &Store{}
	_ store.DataPrivacyRequestComments = &Store{}
	_ store.FederationExposedModules   = &Store{}
	_ store.FederationModuleMappings   = &Store{}
	_ store.FederationNodes            = &Store{}
	_ store.FederationNodeSyncs        = &Store{}
	_ store.FederationSharedModules    = &Store{}
	_ store.Flags                      = &Store{}
	_ store.Labels                     = &Store{}
	_ store.Queues                     = &Store{}
	_ store.QueueMessages              = &Store{}
	_ store.RbacRules                  = &Store{}
	_ store.Reminders                  = &Store{}
	_ store.Reports                    = &Store{}
	_ store.ResourceActivitys          = &Store{}
	_ store.ResourceTranslations       = &Store{}
	_ store.Roles                      = &Store{}
	_ store.RoleMembers                = &Store{}
	_ store.SettingValues              = &Store{}
	_ store.Templates                  = &Store{}
	_ store.Users                      = &Store{}
)

// CreateActionlog creates one or more rows in actionlog collection
//
// This function is auto-generated
func (s *Store) CreateActionlog(ctx context.Context, rr ...*actionlogType.Action) (err error) {
	for i := range rr {
		if err = s.checkActionlogConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, actionlogInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateActionlog updates one or more existing entries in actionlog collection
//
// This function is auto-generated
func (s *Store) UpdateActionlog(ctx context.Context, rr ...*actionlogType.Action) (err error) {
	for i := range rr {
		if err = s.checkActionlogConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, actionlogUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertActionlog updates one or more existing entries in actionlog collection
//
// This function is auto-generated
func (s *Store) UpsertActionlog(ctx context.Context, rr ...*actionlogType.Action) (err error) {
	for i := range rr {
		if err = s.checkActionlogConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, actionlogUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteActionlog Deletes one or more entries from actionlog collection
//
// This function is auto-generated
func (s *Store) DeleteActionlog(ctx context.Context, rr ...*actionlogType.Action) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, actionlogDeleteQuery(s.Dialect, actionlogPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteActionlogByID deletes single entry from actionlog collection
//
// This function is auto-generated
func (s *Store) DeleteActionlogByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, actionlogDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateActionlogs Deletes all rows from the actionlog collection
func (s Store) TruncateActionlogs(ctx context.Context) error {
	return s.Exec(ctx, actionlogTruncateQuery(s.Dialect))
}

// SearchActionlogs returns (filtered) set of Actionlogs
//
// This function is auto-generated
func (s *Store) SearchActionlogs(ctx context.Context, f actionlogType.Filter) (set actionlogType.ActionSet, _ actionlogType.Filter, err error) {

	set, _, err = s.QueryActionlogs(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryActionlogs queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryActionlogs(
	ctx context.Context,
	f actionlogType.Filter,
) (_ []*actionlogType.Action, more bool, err error) {
	var (
		set         = make([]*actionlogType.Action, 0, DefaultSliceCapacity)
		res         *actionlogType.Action
		aux         *auxActionlog
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.Actionlog != nil {
		// extended filter set
		tExpr, f, err = s.Filters.Actionlog(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ActionlogFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for Actionlog: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := actionlogSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableActionlogFields()); err != nil {
		err = fmt.Errorf("could generate order expression for Actionlog: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query Actionlog: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query Actionlog: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query Actionlog: %w", err)
			return
		}

		aux = new(auxActionlog)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for Actionlog: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode Actionlog: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, false, err

}

// LookupActionlogByID searches for action log by ID
//
// This function is auto-generated
func (s *Store) LookupActionlogByID(ctx context.Context, id uint64) (_ *actionlogType.Action, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxActionlog)
		lookup = actionlogSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableActionlogFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableActionlogFields() map[string]string {
	return map[string]string{
		"id":        "id",
		"timestamp": "timestamp",
	}
}

// collectActionlogCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectActionlogCursorValues(res *actionlogType.Action, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "timestamp":
					cur.Set(c.Column, res.Timestamp, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkActionlogConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkActionlogConstraints(ctx context.Context, res *actionlogType.Action) (err error) {
	return nil
}

// CreateApigwFilter creates one or more rows in apigwFilter collection
//
// This function is auto-generated
func (s *Store) CreateApigwFilter(ctx context.Context, rr ...*systemType.ApigwFilter) (err error) {
	for i := range rr {
		if err = s.checkApigwFilterConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, apigwFilterInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateApigwFilter updates one or more existing entries in apigwFilter collection
//
// This function is auto-generated
func (s *Store) UpdateApigwFilter(ctx context.Context, rr ...*systemType.ApigwFilter) (err error) {
	for i := range rr {
		if err = s.checkApigwFilterConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, apigwFilterUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertApigwFilter updates one or more existing entries in apigwFilter collection
//
// This function is auto-generated
func (s *Store) UpsertApigwFilter(ctx context.Context, rr ...*systemType.ApigwFilter) (err error) {
	for i := range rr {
		if err = s.checkApigwFilterConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, apigwFilterUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteApigwFilter Deletes one or more entries from apigwFilter collection
//
// This function is auto-generated
func (s *Store) DeleteApigwFilter(ctx context.Context, rr ...*systemType.ApigwFilter) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, apigwFilterDeleteQuery(s.Dialect, apigwFilterPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteApigwFilterByID deletes single entry from apigwFilter collection
//
// This function is auto-generated
func (s *Store) DeleteApigwFilterByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, apigwFilterDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateApigwFilters Deletes all rows from the apigwFilter collection
func (s Store) TruncateApigwFilters(ctx context.Context) error {
	return s.Exec(ctx, apigwFilterTruncateQuery(s.Dialect))
}

// SearchApigwFilters returns (filtered) set of ApigwFilters
//
// This function is auto-generated
func (s *Store) SearchApigwFilters(ctx context.Context, f systemType.ApigwFilterFilter) (set systemType.ApigwFilterSet, _ systemType.ApigwFilterFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfApigwFilters fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfApigwFilters(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfApigwFilters collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfApigwFilters(
	ctx context.Context,
	filter systemType.ApigwFilterFilter,
	sort filter.SortExprSet,
) (set []*systemType.ApigwFilter, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.ApigwFilter

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.ApigwFilterFilter
	)

	set = make([]*systemType.ApigwFilter, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryApigwFilters(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectApigwFilterCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectApigwFilterCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectApigwFilterCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryApigwFilters queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryApigwFilters(
	ctx context.Context,
	f systemType.ApigwFilterFilter,
) (_ []*systemType.ApigwFilter, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.ApigwFilter, 0, DefaultSliceCapacity)
		res         *systemType.ApigwFilter
		aux         *auxApigwFilter
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.ApigwFilter != nil {
		// extended filter set
		tExpr, f, err = s.Filters.ApigwFilter(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ApigwFilterFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for ApigwFilter: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableApigwFilterFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := apigwFilterSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableApigwFilterFields()); err != nil {
		err = fmt.Errorf("could generate order expression for ApigwFilter: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query ApigwFilter: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query ApigwFilter: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query ApigwFilter: %w", err)
			return
		}

		aux = new(auxApigwFilter)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for ApigwFilter: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode ApigwFilter: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupApigwFilterByID searches for filter by ID
//
// This function is auto-generated
func (s *Store) LookupApigwFilterByID(ctx context.Context, id uint64) (_ *systemType.ApigwFilter, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxApigwFilter)
		lookup = apigwFilterSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupApigwFilterByRoute searches for filter by route
//
// This function is auto-generated
func (s *Store) LookupApigwFilterByRoute(ctx context.Context, route uint64) (_ *systemType.ApigwFilter, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxApigwFilter)
		lookup = apigwFilterSelectQuery(s.Dialect).Where(
			goqu.I("rel_route").Eq(route),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableApigwFilterFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableApigwFilterFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"enabled":    "enabled",
		"id":         "id",
		"kind":       "kind",
		"route":      "route",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
		"weight":     "weight",
	}
}

// collectApigwFilterCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectApigwFilterCursorValues(res *systemType.ApigwFilter, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "route":
					cur.Set(c.Column, res.Route, c.Descending)
				case "weight":
					cur.Set(c.Column, res.Weight, c.Descending)
				case "kind":
					cur.Set(c.Column, res.Kind, c.Descending)
				case "enabled":
					cur.Set(c.Column, res.Enabled, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkApigwFilterConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkApigwFilterConstraints(ctx context.Context, res *systemType.ApigwFilter) (err error) {
	return nil
}

// CreateApigwRoute creates one or more rows in apigwRoute collection
//
// This function is auto-generated
func (s *Store) CreateApigwRoute(ctx context.Context, rr ...*systemType.ApigwRoute) (err error) {
	for i := range rr {
		if err = s.checkApigwRouteConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, apigwRouteInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateApigwRoute updates one or more existing entries in apigwRoute collection
//
// This function is auto-generated
func (s *Store) UpdateApigwRoute(ctx context.Context, rr ...*systemType.ApigwRoute) (err error) {
	for i := range rr {
		if err = s.checkApigwRouteConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, apigwRouteUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertApigwRoute updates one or more existing entries in apigwRoute collection
//
// This function is auto-generated
func (s *Store) UpsertApigwRoute(ctx context.Context, rr ...*systemType.ApigwRoute) (err error) {
	for i := range rr {
		if err = s.checkApigwRouteConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, apigwRouteUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteApigwRoute Deletes one or more entries from apigwRoute collection
//
// This function is auto-generated
func (s *Store) DeleteApigwRoute(ctx context.Context, rr ...*systemType.ApigwRoute) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, apigwRouteDeleteQuery(s.Dialect, apigwRoutePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteApigwRouteByID deletes single entry from apigwRoute collection
//
// This function is auto-generated
func (s *Store) DeleteApigwRouteByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, apigwRouteDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateApigwRoutes Deletes all rows from the apigwRoute collection
func (s Store) TruncateApigwRoutes(ctx context.Context) error {
	return s.Exec(ctx, apigwRouteTruncateQuery(s.Dialect))
}

// SearchApigwRoutes returns (filtered) set of ApigwRoutes
//
// This function is auto-generated
func (s *Store) SearchApigwRoutes(ctx context.Context, f systemType.ApigwRouteFilter) (set systemType.ApigwRouteSet, _ systemType.ApigwRouteFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfApigwRoutes fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfApigwRoutes(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfApigwRoutes collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfApigwRoutes(
	ctx context.Context,
	filter systemType.ApigwRouteFilter,
	sort filter.SortExprSet,
) (set []*systemType.ApigwRoute, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.ApigwRoute

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.ApigwRouteFilter
	)

	set = make([]*systemType.ApigwRoute, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryApigwRoutes(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectApigwRouteCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectApigwRouteCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectApigwRouteCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryApigwRoutes queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryApigwRoutes(
	ctx context.Context,
	f systemType.ApigwRouteFilter,
) (_ []*systemType.ApigwRoute, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.ApigwRoute, 0, DefaultSliceCapacity)
		res         *systemType.ApigwRoute
		aux         *auxApigwRoute
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.ApigwRoute != nil {
		// extended filter set
		tExpr, f, err = s.Filters.ApigwRoute(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ApigwRouteFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for ApigwRoute: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableApigwRouteFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := apigwRouteSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableApigwRouteFields()); err != nil {
		err = fmt.Errorf("could generate order expression for ApigwRoute: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query ApigwRoute: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query ApigwRoute: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query ApigwRoute: %w", err)
			return
		}

		aux = new(auxApigwRoute)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for ApigwRoute: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode ApigwRoute: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupApigwRouteByID searches for route by ID
//
// It returns route even if deleted or suspended
//
// This function is auto-generated
func (s *Store) LookupApigwRouteByID(ctx context.Context, id uint64) (_ *systemType.ApigwRoute, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxApigwRoute)
		lookup = apigwRouteSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupApigwRouteByEndpoint searches for route by endpoint
//
// It returns route even if deleted or suspended
//
// This function is auto-generated
func (s *Store) LookupApigwRouteByEndpoint(ctx context.Context, endpoint string) (_ *systemType.ApigwRoute, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxApigwRoute)
		lookup = apigwRouteSelectQuery(s.Dialect).Where(
			goqu.I("endpoint").Eq(endpoint),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableApigwRouteFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableApigwRouteFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"enabled":    "enabled",
		"endpoint":   "endpoint",
		"group":      "group",
		"id":         "id",
		"method":     "method",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
	}
}

// collectApigwRouteCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectApigwRouteCursorValues(res *systemType.ApigwRoute, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "endpoint":
					cur.Set(c.Column, res.Endpoint, c.Descending)
				case "method":
					cur.Set(c.Column, res.Method, c.Descending)
				case "enabled":
					cur.Set(c.Column, res.Enabled, c.Descending)
				case "group":
					cur.Set(c.Column, res.Group, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkApigwRouteConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkApigwRouteConstraints(ctx context.Context, res *systemType.ApigwRoute) (err error) {
	return nil
}

// CreateApplication creates one or more rows in application collection
//
// This function is auto-generated
func (s *Store) CreateApplication(ctx context.Context, rr ...*systemType.Application) (err error) {
	for i := range rr {
		if err = s.checkApplicationConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, applicationInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateApplication updates one or more existing entries in application collection
//
// This function is auto-generated
func (s *Store) UpdateApplication(ctx context.Context, rr ...*systemType.Application) (err error) {
	for i := range rr {
		if err = s.checkApplicationConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, applicationUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertApplication updates one or more existing entries in application collection
//
// This function is auto-generated
func (s *Store) UpsertApplication(ctx context.Context, rr ...*systemType.Application) (err error) {
	for i := range rr {
		if err = s.checkApplicationConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, applicationUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteApplication Deletes one or more entries from application collection
//
// This function is auto-generated
func (s *Store) DeleteApplication(ctx context.Context, rr ...*systemType.Application) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, applicationDeleteQuery(s.Dialect, applicationPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteApplicationByID deletes single entry from application collection
//
// This function is auto-generated
func (s *Store) DeleteApplicationByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, applicationDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateApplications Deletes all rows from the application collection
func (s Store) TruncateApplications(ctx context.Context) error {
	return s.Exec(ctx, applicationTruncateQuery(s.Dialect))
}

// SearchApplications returns (filtered) set of Applications
//
// This function is auto-generated
func (s *Store) SearchApplications(ctx context.Context, f systemType.ApplicationFilter) (set systemType.ApplicationSet, _ systemType.ApplicationFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfApplications fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfApplications(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfApplications collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfApplications(
	ctx context.Context,
	filter systemType.ApplicationFilter,
	sort filter.SortExprSet,
) (set []*systemType.Application, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.Application

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.ApplicationFilter
	)

	set = make([]*systemType.Application, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryApplications(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectApplicationCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectApplicationCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectApplicationCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryApplications queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryApplications(
	ctx context.Context,
	f systemType.ApplicationFilter,
) (_ []*systemType.Application, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.Application, 0, DefaultSliceCapacity)
		res         *systemType.Application
		aux         *auxApplication
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.Application != nil {
		// extended filter set
		tExpr, f, err = s.Filters.Application(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ApplicationFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for Application: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableApplicationFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := applicationSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableApplicationFields()); err != nil {
		err = fmt.Errorf("could generate order expression for Application: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query Application: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query Application: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query Application: %w", err)
			return
		}

		aux = new(auxApplication)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for Application: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode Application: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupApplicationByID searches for role by ID
//
// It returns role even if deleted or suspended
//
// This function is auto-generated
func (s *Store) LookupApplicationByID(ctx context.Context, id uint64) (_ *systemType.Application, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxApplication)
		lookup = applicationSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableApplicationFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableApplicationFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"enabled":    "enabled",
		"id":         "id",
		"name":       "name",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
		"weight":     "weight",
	}
}

// collectApplicationCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectApplicationCursorValues(res *systemType.Application, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
				case "enabled":
					cur.Set(c.Column, res.Enabled, c.Descending)
				case "weight":
					cur.Set(c.Column, res.Weight, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkApplicationConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkApplicationConstraints(ctx context.Context, res *systemType.Application) (err error) {
	return nil
}

// CreateAttachment creates one or more rows in attachment collection
//
// This function is auto-generated
func (s *Store) CreateAttachment(ctx context.Context, rr ...*systemType.Attachment) (err error) {
	for i := range rr {
		if err = s.checkAttachmentConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, attachmentInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateAttachment updates one or more existing entries in attachment collection
//
// This function is auto-generated
func (s *Store) UpdateAttachment(ctx context.Context, rr ...*systemType.Attachment) (err error) {
	for i := range rr {
		if err = s.checkAttachmentConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, attachmentUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertAttachment updates one or more existing entries in attachment collection
//
// This function is auto-generated
func (s *Store) UpsertAttachment(ctx context.Context, rr ...*systemType.Attachment) (err error) {
	for i := range rr {
		if err = s.checkAttachmentConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, attachmentUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteAttachment Deletes one or more entries from attachment collection
//
// This function is auto-generated
func (s *Store) DeleteAttachment(ctx context.Context, rr ...*systemType.Attachment) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, attachmentDeleteQuery(s.Dialect, attachmentPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteAttachmentByID deletes single entry from attachment collection
//
// This function is auto-generated
func (s *Store) DeleteAttachmentByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, attachmentDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateAttachments Deletes all rows from the attachment collection
func (s Store) TruncateAttachments(ctx context.Context) error {
	return s.Exec(ctx, attachmentTruncateQuery(s.Dialect))
}

// SearchAttachments returns (filtered) set of Attachments
//
// This function is auto-generated
func (s *Store) SearchAttachments(ctx context.Context, f systemType.AttachmentFilter) (set systemType.AttachmentSet, _ systemType.AttachmentFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfAttachments fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfAttachments(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfAttachments collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfAttachments(
	ctx context.Context,
	filter systemType.AttachmentFilter,
	sort filter.SortExprSet,
) (set []*systemType.Attachment, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.Attachment

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.AttachmentFilter
	)

	set = make([]*systemType.Attachment, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryAttachments(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectAttachmentCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectAttachmentCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectAttachmentCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryAttachments queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryAttachments(
	ctx context.Context,
	f systemType.AttachmentFilter,
) (_ []*systemType.Attachment, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.Attachment, 0, DefaultSliceCapacity)
		res         *systemType.Attachment
		aux         *auxAttachment
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.Attachment != nil {
		// extended filter set
		tExpr, f, err = s.Filters.Attachment(s, f)
	} else {
		// using generated filter
		tExpr, f, err = AttachmentFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for Attachment: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableAttachmentFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := attachmentSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableAttachmentFields()); err != nil {
		err = fmt.Errorf("could generate order expression for Attachment: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query Attachment: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query Attachment: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query Attachment: %w", err)
			return
		}

		aux = new(auxAttachment)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for Attachment: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode Attachment: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupAttachmentByID
//
// This function is auto-generated
func (s *Store) LookupAttachmentByID(ctx context.Context, id uint64) (_ *systemType.Attachment, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAttachment)
		lookup = attachmentSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableAttachmentFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableAttachmentFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"id":         "id",
		"kind":       "kind",
		"name":       "name",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
	}
}

// collectAttachmentCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectAttachmentCursorValues(res *systemType.Attachment, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "kind":
					cur.Set(c.Column, res.Kind, c.Descending)
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkAttachmentConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkAttachmentConstraints(ctx context.Context, res *systemType.Attachment) (err error) {
	return nil
}

// CreateAuthClient creates one or more rows in authClient collection
//
// This function is auto-generated
func (s *Store) CreateAuthClient(ctx context.Context, rr ...*systemType.AuthClient) (err error) {
	for i := range rr {
		if err = s.checkAuthClientConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, authClientInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateAuthClient updates one or more existing entries in authClient collection
//
// This function is auto-generated
func (s *Store) UpdateAuthClient(ctx context.Context, rr ...*systemType.AuthClient) (err error) {
	for i := range rr {
		if err = s.checkAuthClientConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, authClientUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertAuthClient updates one or more existing entries in authClient collection
//
// This function is auto-generated
func (s *Store) UpsertAuthClient(ctx context.Context, rr ...*systemType.AuthClient) (err error) {
	for i := range rr {
		if err = s.checkAuthClientConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, authClientUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteAuthClient Deletes one or more entries from authClient collection
//
// This function is auto-generated
func (s *Store) DeleteAuthClient(ctx context.Context, rr ...*systemType.AuthClient) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, authClientDeleteQuery(s.Dialect, authClientPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteAuthClientByID deletes single entry from authClient collection
//
// This function is auto-generated
func (s *Store) DeleteAuthClientByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, authClientDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateAuthClients Deletes all rows from the authClient collection
func (s Store) TruncateAuthClients(ctx context.Context) error {
	return s.Exec(ctx, authClientTruncateQuery(s.Dialect))
}

// SearchAuthClients returns (filtered) set of AuthClients
//
// This function is auto-generated
func (s *Store) SearchAuthClients(ctx context.Context, f systemType.AuthClientFilter) (set systemType.AuthClientSet, _ systemType.AuthClientFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfAuthClients fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfAuthClients(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfAuthClients collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfAuthClients(
	ctx context.Context,
	filter systemType.AuthClientFilter,
	sort filter.SortExprSet,
) (set []*systemType.AuthClient, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.AuthClient

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.AuthClientFilter
	)

	set = make([]*systemType.AuthClient, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryAuthClients(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectAuthClientCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectAuthClientCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectAuthClientCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryAuthClients queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryAuthClients(
	ctx context.Context,
	f systemType.AuthClientFilter,
) (_ []*systemType.AuthClient, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.AuthClient, 0, DefaultSliceCapacity)
		res         *systemType.AuthClient
		aux         *auxAuthClient
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.AuthClient != nil {
		// extended filter set
		tExpr, f, err = s.Filters.AuthClient(s, f)
	} else {
		// using generated filter
		tExpr, f, err = AuthClientFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for AuthClient: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableAuthClientFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := authClientSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableAuthClientFields()); err != nil {
		err = fmt.Errorf("could generate order expression for AuthClient: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query AuthClient: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query AuthClient: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query AuthClient: %w", err)
			return
		}

		aux = new(auxAuthClient)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for AuthClient: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode AuthClient: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupAuthClientByID 	searches for auth client by ID
//
// 	It returns auth clint even if deleted
//
// This function is auto-generated
func (s *Store) LookupAuthClientByID(ctx context.Context, id uint64) (_ *systemType.AuthClient, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAuthClient)
		lookup = authClientSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupAuthClientByHandle searches for auth client by ID
//
// It returns auth clint even if deleted
//
// This function is auto-generated
func (s *Store) LookupAuthClientByHandle(ctx context.Context, handle string) (_ *systemType.AuthClient, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAuthClient)
		lookup = authClientSelectQuery(s.Dialect).Where(
			s.Functions.LOWER(goqu.I("handle")).Eq(strings.ToLower(handle)),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableAuthClientFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableAuthClientFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"enabled":    "enabled",
		"expires_at": "expires_at",
		"expiresat":  "expires_at",
		"handle":     "handle",
		"id":         "id",
		"trusted":    "trusted",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
		"valid_from": "valid_from",
		"validfrom":  "valid_from",
	}
}

// collectAuthClientCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectAuthClientCursorValues(res *systemType.AuthClient, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "enabled":
					cur.Set(c.Column, res.Enabled, c.Descending)
				case "trusted":
					cur.Set(c.Column, res.Trusted, c.Descending)
				case "validFrom":
					cur.Set(c.Column, res.ValidFrom, c.Descending)
				case "expiresAt":
					cur.Set(c.Column, res.ExpiresAt, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkAuthClientConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkAuthClientConstraints(ctx context.Context, res *systemType.AuthClient) (err error) {
	err = func() (err error) {

		// handling string type as default
		if len(res.Handle) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupAuthClientByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	return nil
}

// CreateAuthConfirmedClient creates one or more rows in authConfirmedClient collection
//
// This function is auto-generated
func (s *Store) CreateAuthConfirmedClient(ctx context.Context, rr ...*systemType.AuthConfirmedClient) (err error) {
	for i := range rr {
		if err = s.checkAuthConfirmedClientConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, authConfirmedClientInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateAuthConfirmedClient updates one or more existing entries in authConfirmedClient collection
//
// This function is auto-generated
func (s *Store) UpdateAuthConfirmedClient(ctx context.Context, rr ...*systemType.AuthConfirmedClient) (err error) {
	for i := range rr {
		if err = s.checkAuthConfirmedClientConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, authConfirmedClientUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertAuthConfirmedClient updates one or more existing entries in authConfirmedClient collection
//
// This function is auto-generated
func (s *Store) UpsertAuthConfirmedClient(ctx context.Context, rr ...*systemType.AuthConfirmedClient) (err error) {
	for i := range rr {
		if err = s.checkAuthConfirmedClientConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, authConfirmedClientUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteAuthConfirmedClient Deletes one or more entries from authConfirmedClient collection
//
// This function is auto-generated
func (s *Store) DeleteAuthConfirmedClient(ctx context.Context, rr ...*systemType.AuthConfirmedClient) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, authConfirmedClientDeleteQuery(s.Dialect, authConfirmedClientPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteAuthConfirmedClientByID deletes single entry from authConfirmedClient collection
//
// This function is auto-generated
func (s *Store) DeleteAuthConfirmedClientByUserIDClientID(ctx context.Context, userID uint64, clientID uint64) error {
	return s.Exec(ctx, authConfirmedClientDeleteQuery(s.Dialect, goqu.Ex{
		"rel_user":   userID,
		"rel_client": clientID,
	}))
}

// TruncateAuthConfirmedClients Deletes all rows from the authConfirmedClient collection
func (s Store) TruncateAuthConfirmedClients(ctx context.Context) error {
	return s.Exec(ctx, authConfirmedClientTruncateQuery(s.Dialect))
}

// SearchAuthConfirmedClients returns (filtered) set of AuthConfirmedClients
//
// This function is auto-generated
func (s *Store) SearchAuthConfirmedClients(ctx context.Context, f systemType.AuthConfirmedClientFilter) (set systemType.AuthConfirmedClientSet, _ systemType.AuthConfirmedClientFilter, err error) {

	set, _, err = s.QueryAuthConfirmedClients(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryAuthConfirmedClients queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryAuthConfirmedClients(
	ctx context.Context,
	f systemType.AuthConfirmedClientFilter,
) (_ []*systemType.AuthConfirmedClient, more bool, err error) {
	var (
		set         = make([]*systemType.AuthConfirmedClient, 0, DefaultSliceCapacity)
		res         *systemType.AuthConfirmedClient
		aux         *auxAuthConfirmedClient
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.AuthConfirmedClient != nil {
		// extended filter set
		tExpr, f, err = s.Filters.AuthConfirmedClient(s, f)
	} else {
		// using generated filter
		tExpr, f, err = AuthConfirmedClientFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for AuthConfirmedClient: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := authConfirmedClientSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query AuthConfirmedClient: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query AuthConfirmedClient: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query AuthConfirmedClient: %w", err)
			return
		}

		aux = new(auxAuthConfirmedClient)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for AuthConfirmedClient: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode AuthConfirmedClient: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, false, err

}

// LookupAuthConfirmedClientByUserIDClientID
//
// This function is auto-generated
func (s *Store) LookupAuthConfirmedClientByUserIDClientID(ctx context.Context, userID uint64, clientID uint64) (_ *systemType.AuthConfirmedClient, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAuthConfirmedClient)
		lookup = authConfirmedClientSelectQuery(s.Dialect).Where(
			goqu.I("rel_user").Eq(userID),
			goqu.I("rel_client").Eq(clientID),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableAuthConfirmedClientFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableAuthConfirmedClientFields() map[string]string {
	return map[string]string{
		"client_id":    "client_id",
		"clientid":     "client_id",
		"confirmed_at": "confirmed_at",
		"confirmedat":  "confirmed_at",
		"user_id":      "user_id",
		"userid":       "user_id",
	}
}

// collectAuthConfirmedClientCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectAuthConfirmedClientCursorValues(res *systemType.AuthConfirmedClient, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkUserID   bool
		pkClientID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "userID":
					cur.Set(c.Column, res.UserID, c.Descending)
					pkUserID = true
				case "clientID":
					cur.Set(c.Column, res.ClientID, c.Descending)
					pkClientID = true
				case "confirmedAt":
					cur.Set(c.Column, res.ConfirmedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkUserID {
		collect(&filter.SortExpr{Column: "userID", Descending: false})
	}
	if !hasUnique || !pkClientID {
		collect(&filter.SortExpr{Column: "clientID", Descending: false})
	}

	return cur

}

// checkAuthConfirmedClientConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkAuthConfirmedClientConstraints(ctx context.Context, res *systemType.AuthConfirmedClient) (err error) {
	return nil
}

// CreateAuthOa2token creates one or more rows in authOa2token collection
//
// This function is auto-generated
func (s *Store) CreateAuthOa2token(ctx context.Context, rr ...*systemType.AuthOa2token) (err error) {
	for i := range rr {
		if err = s.checkAuthOa2tokenConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, authOa2tokenInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateAuthOa2token updates one or more existing entries in authOa2token collection
//
// This function is auto-generated
func (s *Store) UpdateAuthOa2token(ctx context.Context, rr ...*systemType.AuthOa2token) (err error) {
	for i := range rr {
		if err = s.checkAuthOa2tokenConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, authOa2tokenUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertAuthOa2token updates one or more existing entries in authOa2token collection
//
// This function is auto-generated
func (s *Store) UpsertAuthOa2token(ctx context.Context, rr ...*systemType.AuthOa2token) (err error) {
	for i := range rr {
		if err = s.checkAuthOa2tokenConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, authOa2tokenUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteAuthOa2token Deletes one or more entries from authOa2token collection
//
// This function is auto-generated
func (s *Store) DeleteAuthOa2token(ctx context.Context, rr ...*systemType.AuthOa2token) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, authOa2tokenDeleteQuery(s.Dialect, authOa2tokenPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteAuthOa2tokenByID deletes single entry from authOa2token collection
//
// This function is auto-generated
func (s *Store) DeleteAuthOa2tokenByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, authOa2tokenDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateAuthOa2tokens Deletes all rows from the authOa2token collection
func (s Store) TruncateAuthOa2tokens(ctx context.Context) error {
	return s.Exec(ctx, authOa2tokenTruncateQuery(s.Dialect))
}

// SearchAuthOa2tokens returns (filtered) set of AuthOa2tokens
//
// This function is auto-generated
func (s *Store) SearchAuthOa2tokens(ctx context.Context, f systemType.AuthOa2tokenFilter) (set systemType.AuthOa2tokenSet, _ systemType.AuthOa2tokenFilter, err error) {

	set, _, err = s.QueryAuthOa2tokens(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryAuthOa2tokens queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryAuthOa2tokens(
	ctx context.Context,
	f systemType.AuthOa2tokenFilter,
) (_ []*systemType.AuthOa2token, more bool, err error) {
	var (
		set         = make([]*systemType.AuthOa2token, 0, DefaultSliceCapacity)
		res         *systemType.AuthOa2token
		aux         *auxAuthOa2token
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.AuthOa2token != nil {
		// extended filter set
		tExpr, f, err = s.Filters.AuthOa2token(s, f)
	} else {
		// using generated filter
		tExpr, f, err = AuthOa2tokenFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for AuthOa2token: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := authOa2tokenSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query AuthOa2token: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query AuthOa2token: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query AuthOa2token: %w", err)
			return
		}

		aux = new(auxAuthOa2token)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for AuthOa2token: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode AuthOa2token: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, false, err

}

// LookupAuthOa2tokenByID
//
// This function is auto-generated
func (s *Store) LookupAuthOa2tokenByID(ctx context.Context, id uint64) (_ *systemType.AuthOa2token, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAuthOa2token)
		lookup = authOa2tokenSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupAuthOa2tokenByCode
//
// This function is auto-generated
func (s *Store) LookupAuthOa2tokenByCode(ctx context.Context, code string) (_ *systemType.AuthOa2token, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAuthOa2token)
		lookup = authOa2tokenSelectQuery(s.Dialect).Where(
			goqu.I("code").Eq(code),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupAuthOa2tokenByAccess
//
// This function is auto-generated
func (s *Store) LookupAuthOa2tokenByAccess(ctx context.Context, access string) (_ *systemType.AuthOa2token, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAuthOa2token)
		lookup = authOa2tokenSelectQuery(s.Dialect).Where(
			goqu.I("access").Eq(access),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupAuthOa2tokenByRefresh
//
// This function is auto-generated
func (s *Store) LookupAuthOa2tokenByRefresh(ctx context.Context, refresh string) (_ *systemType.AuthOa2token, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAuthOa2token)
		lookup = authOa2tokenSelectQuery(s.Dialect).Where(
			goqu.I("refresh").Eq(refresh),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableAuthOa2tokenFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableAuthOa2tokenFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"expires_at": "expires_at",
		"expiresat":  "expires_at",
		"id":         "id",
	}
}

// collectAuthOa2tokenCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectAuthOa2tokenCursorValues(res *systemType.AuthOa2token, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "expiresAt":
					cur.Set(c.Column, res.ExpiresAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkAuthOa2tokenConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkAuthOa2tokenConstraints(ctx context.Context, res *systemType.AuthOa2token) (err error) {
	return nil
}

// CreateAuthSession creates one or more rows in authSession collection
//
// This function is auto-generated
func (s *Store) CreateAuthSession(ctx context.Context, rr ...*systemType.AuthSession) (err error) {
	for i := range rr {
		if err = s.checkAuthSessionConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, authSessionInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateAuthSession updates one or more existing entries in authSession collection
//
// This function is auto-generated
func (s *Store) UpdateAuthSession(ctx context.Context, rr ...*systemType.AuthSession) (err error) {
	for i := range rr {
		if err = s.checkAuthSessionConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, authSessionUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertAuthSession updates one or more existing entries in authSession collection
//
// This function is auto-generated
func (s *Store) UpsertAuthSession(ctx context.Context, rr ...*systemType.AuthSession) (err error) {
	for i := range rr {
		if err = s.checkAuthSessionConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, authSessionUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteAuthSession Deletes one or more entries from authSession collection
//
// This function is auto-generated
func (s *Store) DeleteAuthSession(ctx context.Context, rr ...*systemType.AuthSession) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, authSessionDeleteQuery(s.Dialect, authSessionPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteAuthSessionByID deletes single entry from authSession collection
//
// This function is auto-generated
func (s *Store) DeleteAuthSessionByID(ctx context.Context, id string) error {
	return s.Exec(ctx, authSessionDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateAuthSessions Deletes all rows from the authSession collection
func (s Store) TruncateAuthSessions(ctx context.Context) error {
	return s.Exec(ctx, authSessionTruncateQuery(s.Dialect))
}

// SearchAuthSessions returns (filtered) set of AuthSessions
//
// This function is auto-generated
func (s *Store) SearchAuthSessions(ctx context.Context, f systemType.AuthSessionFilter) (set systemType.AuthSessionSet, _ systemType.AuthSessionFilter, err error) {

	set, _, err = s.QueryAuthSessions(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryAuthSessions queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryAuthSessions(
	ctx context.Context,
	f systemType.AuthSessionFilter,
) (_ []*systemType.AuthSession, more bool, err error) {
	var (
		set         = make([]*systemType.AuthSession, 0, DefaultSliceCapacity)
		res         *systemType.AuthSession
		aux         *auxAuthSession
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.AuthSession != nil {
		// extended filter set
		tExpr, f, err = s.Filters.AuthSession(s, f)
	} else {
		// using generated filter
		tExpr, f, err = AuthSessionFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for AuthSession: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := authSessionSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query AuthSession: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query AuthSession: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query AuthSession: %w", err)
			return
		}

		aux = new(auxAuthSession)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for AuthSession: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode AuthSession: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, false, err

}

// LookupAuthSessionByID
//
// This function is auto-generated
func (s *Store) LookupAuthSessionByID(ctx context.Context, id string) (_ *systemType.AuthSession, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAuthSession)
		lookup = authSessionSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableAuthSessionFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableAuthSessionFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"expires_at": "expires_at",
		"expiresat":  "expires_at",
		"id":         "id",
	}
}

// collectAuthSessionCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectAuthSessionCursorValues(res *systemType.AuthSession, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "expiresAt":
					cur.Set(c.Column, res.ExpiresAt, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkAuthSessionConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkAuthSessionConstraints(ctx context.Context, res *systemType.AuthSession) (err error) {
	return nil
}

// CreateAutomationSession creates one or more rows in automationSession collection
//
// This function is auto-generated
func (s *Store) CreateAutomationSession(ctx context.Context, rr ...*automationType.Session) (err error) {
	for i := range rr {
		if err = s.checkAutomationSessionConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, automationSessionInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateAutomationSession updates one or more existing entries in automationSession collection
//
// This function is auto-generated
func (s *Store) UpdateAutomationSession(ctx context.Context, rr ...*automationType.Session) (err error) {
	for i := range rr {
		if err = s.checkAutomationSessionConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, automationSessionUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertAutomationSession updates one or more existing entries in automationSession collection
//
// This function is auto-generated
func (s *Store) UpsertAutomationSession(ctx context.Context, rr ...*automationType.Session) (err error) {
	for i := range rr {
		if err = s.checkAutomationSessionConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, automationSessionUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteAutomationSession Deletes one or more entries from automationSession collection
//
// This function is auto-generated
func (s *Store) DeleteAutomationSession(ctx context.Context, rr ...*automationType.Session) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, automationSessionDeleteQuery(s.Dialect, automationSessionPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteAutomationSessionByID deletes single entry from automationSession collection
//
// This function is auto-generated
func (s *Store) DeleteAutomationSessionByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, automationSessionDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateAutomationSessions Deletes all rows from the automationSession collection
func (s Store) TruncateAutomationSessions(ctx context.Context) error {
	return s.Exec(ctx, automationSessionTruncateQuery(s.Dialect))
}

// SearchAutomationSessions returns (filtered) set of AutomationSessions
//
// This function is auto-generated
func (s *Store) SearchAutomationSessions(ctx context.Context, f automationType.SessionFilter) (set automationType.SessionSet, _ automationType.SessionFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfAutomationSessions fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfAutomationSessions(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
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
//
// This function is auto-generated
func (s *Store) fetchFullPageOfAutomationSessions(
	ctx context.Context,
	filter automationType.SessionFilter,
	sort filter.SortExprSet,
) (set []*automationType.Session, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*automationType.Session

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter automationType.SessionFilter
	)

	set = make([]*automationType.Session, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryAutomationSessions(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectAutomationSessionCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectAutomationSessionCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectAutomationSessionCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryAutomationSessions queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryAutomationSessions(
	ctx context.Context,
	f automationType.SessionFilter,
) (_ []*automationType.Session, more bool, err error) {
	var (
		ok bool

		set         = make([]*automationType.Session, 0, DefaultSliceCapacity)
		res         *automationType.Session
		aux         *auxAutomationSession
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.AutomationSession != nil {
		// extended filter set
		tExpr, f, err = s.Filters.AutomationSession(s, f)
	} else {
		// using generated filter
		tExpr, f, err = AutomationSessionFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for AutomationSession: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableAutomationSessionFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := automationSessionSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableAutomationSessionFields()); err != nil {
		err = fmt.Errorf("could generate order expression for AutomationSession: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query AutomationSession: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query AutomationSession: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query AutomationSession: %w", err)
			return
		}

		aux = new(auxAutomationSession)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for AutomationSession: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode AutomationSession: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupAutomationSessionByID searches for session by ID
//
// It returns session even if deleted
//
// This function is auto-generated
func (s *Store) LookupAutomationSessionByID(ctx context.Context, id uint64) (_ *automationType.Session, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAutomationSession)
		lookup = automationSessionSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableAutomationSessionFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableAutomationSessionFields() map[string]string {
	return map[string]string{
		"completed_at":  "completed_at",
		"completedat":   "completed_at",
		"created_at":    "created_at",
		"createdat":     "created_at",
		"event_type":    "event_type",
		"eventtype":     "event_type",
		"id":            "id",
		"purge_at":      "purge_at",
		"purgeat":       "purge_at",
		"resource_type": "resource_type",
		"resourcetype":  "resource_type",
		"status":        "status",
		"suspended_at":  "suspended_at",
		"suspendedat":   "suspended_at",
		"workflow_id":   "workflow_id",
		"workflowid":    "workflow_id",
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
//
// This function is auto-generated
func (s *Store) collectAutomationSessionCursorValues(res *automationType.Session, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "workflowID":
					cur.Set(c.Column, res.WorkflowID, c.Descending)
				case "status":
					cur.Set(c.Column, res.Status, c.Descending)
				case "eventType":
					cur.Set(c.Column, res.EventType, c.Descending)
				case "resourceType":
					cur.Set(c.Column, res.ResourceType, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "purgeAt":
					cur.Set(c.Column, res.PurgeAt, c.Descending)
				case "suspendedAt":
					cur.Set(c.Column, res.SuspendedAt, c.Descending)
				case "completedAt":
					cur.Set(c.Column, res.CompletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkAutomationSessionConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkAutomationSessionConstraints(ctx context.Context, res *automationType.Session) (err error) {
	return nil
}

// CreateAutomationTrigger creates one or more rows in automationTrigger collection
//
// This function is auto-generated
func (s *Store) CreateAutomationTrigger(ctx context.Context, rr ...*automationType.Trigger) (err error) {
	for i := range rr {
		if err = s.checkAutomationTriggerConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, automationTriggerInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateAutomationTrigger updates one or more existing entries in automationTrigger collection
//
// This function is auto-generated
func (s *Store) UpdateAutomationTrigger(ctx context.Context, rr ...*automationType.Trigger) (err error) {
	for i := range rr {
		if err = s.checkAutomationTriggerConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, automationTriggerUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertAutomationTrigger updates one or more existing entries in automationTrigger collection
//
// This function is auto-generated
func (s *Store) UpsertAutomationTrigger(ctx context.Context, rr ...*automationType.Trigger) (err error) {
	for i := range rr {
		if err = s.checkAutomationTriggerConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, automationTriggerUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteAutomationTrigger Deletes one or more entries from automationTrigger collection
//
// This function is auto-generated
func (s *Store) DeleteAutomationTrigger(ctx context.Context, rr ...*automationType.Trigger) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, automationTriggerDeleteQuery(s.Dialect, automationTriggerPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteAutomationTriggerByID deletes single entry from automationTrigger collection
//
// This function is auto-generated
func (s *Store) DeleteAutomationTriggerByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, automationTriggerDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateAutomationTriggers Deletes all rows from the automationTrigger collection
func (s Store) TruncateAutomationTriggers(ctx context.Context) error {
	return s.Exec(ctx, automationTriggerTruncateQuery(s.Dialect))
}

// SearchAutomationTriggers returns (filtered) set of AutomationTriggers
//
// This function is auto-generated
func (s *Store) SearchAutomationTriggers(ctx context.Context, f automationType.TriggerFilter) (set automationType.TriggerSet, _ automationType.TriggerFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfAutomationTriggers fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfAutomationTriggers(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
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
//
// This function is auto-generated
func (s *Store) fetchFullPageOfAutomationTriggers(
	ctx context.Context,
	filter automationType.TriggerFilter,
	sort filter.SortExprSet,
) (set []*automationType.Trigger, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*automationType.Trigger

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter automationType.TriggerFilter
	)

	set = make([]*automationType.Trigger, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryAutomationTriggers(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectAutomationTriggerCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectAutomationTriggerCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectAutomationTriggerCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryAutomationTriggers queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryAutomationTriggers(
	ctx context.Context,
	f automationType.TriggerFilter,
) (_ []*automationType.Trigger, more bool, err error) {
	var (
		ok bool

		set         = make([]*automationType.Trigger, 0, DefaultSliceCapacity)
		res         *automationType.Trigger
		aux         *auxAutomationTrigger
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.AutomationTrigger != nil {
		// extended filter set
		tExpr, f, err = s.Filters.AutomationTrigger(s, f)
	} else {
		// using generated filter
		tExpr, f, err = AutomationTriggerFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for AutomationTrigger: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableAutomationTriggerFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := automationTriggerSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableAutomationTriggerFields()); err != nil {
		err = fmt.Errorf("could generate order expression for AutomationTrigger: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query AutomationTrigger: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query AutomationTrigger: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query AutomationTrigger: %w", err)
			return
		}

		aux = new(auxAutomationTrigger)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for AutomationTrigger: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode AutomationTrigger: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupAutomationTriggerByID searches for trigger by ID
//
// It returns trigger even if deleted
//
// This function is auto-generated
func (s *Store) LookupAutomationTriggerByID(ctx context.Context, id uint64) (_ *automationType.Trigger, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAutomationTrigger)
		lookup = automationTriggerSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableAutomationTriggerFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableAutomationTriggerFields() map[string]string {
	return map[string]string{
		"created_at":    "created_at",
		"createdat":     "created_at",
		"deleted_at":    "deleted_at",
		"deletedat":     "deleted_at",
		"enabled":       "enabled",
		"event_type":    "event_type",
		"eventtype":     "event_type",
		"id":            "id",
		"resource_type": "resource_type",
		"resourcetype":  "resource_type",
		"updated_at":    "updated_at",
		"updatedat":     "updated_at",
		"workflow_id":   "workflow_id",
		"workflowid":    "workflow_id",
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
//
// This function is auto-generated
func (s *Store) collectAutomationTriggerCursorValues(res *automationType.Trigger, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "workflowID":
					cur.Set(c.Column, res.WorkflowID, c.Descending)
				case "enabled":
					cur.Set(c.Column, res.Enabled, c.Descending)
				case "resourceType":
					cur.Set(c.Column, res.ResourceType, c.Descending)
				case "eventType":
					cur.Set(c.Column, res.EventType, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkAutomationTriggerConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkAutomationTriggerConstraints(ctx context.Context, res *automationType.Trigger) (err error) {
	return nil
}

// CreateAutomationWorkflow creates one or more rows in automationWorkflow collection
//
// This function is auto-generated
func (s *Store) CreateAutomationWorkflow(ctx context.Context, rr ...*automationType.Workflow) (err error) {
	for i := range rr {
		if err = s.checkAutomationWorkflowConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, automationWorkflowInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateAutomationWorkflow updates one or more existing entries in automationWorkflow collection
//
// This function is auto-generated
func (s *Store) UpdateAutomationWorkflow(ctx context.Context, rr ...*automationType.Workflow) (err error) {
	for i := range rr {
		if err = s.checkAutomationWorkflowConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, automationWorkflowUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertAutomationWorkflow updates one or more existing entries in automationWorkflow collection
//
// This function is auto-generated
func (s *Store) UpsertAutomationWorkflow(ctx context.Context, rr ...*automationType.Workflow) (err error) {
	for i := range rr {
		if err = s.checkAutomationWorkflowConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, automationWorkflowUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteAutomationWorkflow Deletes one or more entries from automationWorkflow collection
//
// This function is auto-generated
func (s *Store) DeleteAutomationWorkflow(ctx context.Context, rr ...*automationType.Workflow) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, automationWorkflowDeleteQuery(s.Dialect, automationWorkflowPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteAutomationWorkflowByID deletes single entry from automationWorkflow collection
//
// This function is auto-generated
func (s *Store) DeleteAutomationWorkflowByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, automationWorkflowDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateAutomationWorkflows Deletes all rows from the automationWorkflow collection
func (s Store) TruncateAutomationWorkflows(ctx context.Context) error {
	return s.Exec(ctx, automationWorkflowTruncateQuery(s.Dialect))
}

// SearchAutomationWorkflows returns (filtered) set of AutomationWorkflows
//
// This function is auto-generated
func (s *Store) SearchAutomationWorkflows(ctx context.Context, f automationType.WorkflowFilter) (set automationType.WorkflowSet, _ automationType.WorkflowFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfAutomationWorkflows fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfAutomationWorkflows(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfAutomationWorkflows collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfAutomationWorkflows(
	ctx context.Context,
	filter automationType.WorkflowFilter,
	sort filter.SortExprSet,
) (set []*automationType.Workflow, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*automationType.Workflow

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter automationType.WorkflowFilter
	)

	set = make([]*automationType.Workflow, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryAutomationWorkflows(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectAutomationWorkflowCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectAutomationWorkflowCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectAutomationWorkflowCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryAutomationWorkflows queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryAutomationWorkflows(
	ctx context.Context,
	f automationType.WorkflowFilter,
) (_ []*automationType.Workflow, more bool, err error) {
	var (
		ok bool

		set         = make([]*automationType.Workflow, 0, DefaultSliceCapacity)
		res         *automationType.Workflow
		aux         *auxAutomationWorkflow
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.AutomationWorkflow != nil {
		// extended filter set
		tExpr, f, err = s.Filters.AutomationWorkflow(s, f)
	} else {
		// using generated filter
		tExpr, f, err = AutomationWorkflowFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for AutomationWorkflow: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableAutomationWorkflowFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := automationWorkflowSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableAutomationWorkflowFields()); err != nil {
		err = fmt.Errorf("could generate order expression for AutomationWorkflow: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query AutomationWorkflow: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query AutomationWorkflow: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query AutomationWorkflow: %w", err)
			return
		}

		aux = new(auxAutomationWorkflow)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for AutomationWorkflow: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode AutomationWorkflow: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupAutomationWorkflowByID searches for workflow by ID
//
// It returns workflow even if deleted
//
// This function is auto-generated
func (s *Store) LookupAutomationWorkflowByID(ctx context.Context, id uint64) (_ *automationType.Workflow, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAutomationWorkflow)
		lookup = automationWorkflowSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupAutomationWorkflowByHandle searches for workflow by their handle
//
// It returns only valid workflows
//
// This function is auto-generated
func (s *Store) LookupAutomationWorkflowByHandle(ctx context.Context, handle string) (_ *automationType.Workflow, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxAutomationWorkflow)
		lookup = automationWorkflowSelectQuery(s.Dialect).Where(
			s.Functions.LOWER(goqu.I("handle")).Eq(strings.ToLower(handle)),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableAutomationWorkflowFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableAutomationWorkflowFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"enabled":    "enabled",
		"handle":     "handle",
		"id":         "id",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
	}
}

// collectAutomationWorkflowCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectAutomationWorkflowCursorValues(res *automationType.Workflow, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "enabled":
					cur.Set(c.Column, res.Enabled, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkAutomationWorkflowConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkAutomationWorkflowConstraints(ctx context.Context, res *automationType.Workflow) (err error) {
	err = func() (err error) {

		// handling string type as default
		if len(res.Handle) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupAutomationWorkflowByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	return nil
}

// CreateComposeAttachment creates one or more rows in composeAttachment collection
//
// This function is auto-generated
func (s *Store) CreateComposeAttachment(ctx context.Context, rr ...*composeType.Attachment) (err error) {
	for i := range rr {
		if err = s.checkComposeAttachmentConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeAttachmentInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateComposeAttachment updates one or more existing entries in composeAttachment collection
//
// This function is auto-generated
func (s *Store) UpdateComposeAttachment(ctx context.Context, rr ...*composeType.Attachment) (err error) {
	for i := range rr {
		if err = s.checkComposeAttachmentConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeAttachmentUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertComposeAttachment updates one or more existing entries in composeAttachment collection
//
// This function is auto-generated
func (s *Store) UpsertComposeAttachment(ctx context.Context, rr ...*composeType.Attachment) (err error) {
	for i := range rr {
		if err = s.checkComposeAttachmentConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeAttachmentUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteComposeAttachment Deletes one or more entries from composeAttachment collection
//
// This function is auto-generated
func (s *Store) DeleteComposeAttachment(ctx context.Context, rr ...*composeType.Attachment) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, composeAttachmentDeleteQuery(s.Dialect, composeAttachmentPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteComposeAttachmentByID deletes single entry from composeAttachment collection
//
// This function is auto-generated
func (s *Store) DeleteComposeAttachmentByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, composeAttachmentDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateComposeAttachments Deletes all rows from the composeAttachment collection
func (s Store) TruncateComposeAttachments(ctx context.Context) error {
	return s.Exec(ctx, composeAttachmentTruncateQuery(s.Dialect))
}

// SearchComposeAttachments returns (filtered) set of ComposeAttachments
//
// This function is auto-generated
func (s *Store) SearchComposeAttachments(ctx context.Context, f composeType.AttachmentFilter) (set composeType.AttachmentSet, _ composeType.AttachmentFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfComposeAttachments fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfComposeAttachments(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfComposeAttachments collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfComposeAttachments(
	ctx context.Context,
	filter composeType.AttachmentFilter,
	sort filter.SortExprSet,
) (set []*composeType.Attachment, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*composeType.Attachment

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter composeType.AttachmentFilter
	)

	set = make([]*composeType.Attachment, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryComposeAttachments(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectComposeAttachmentCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectComposeAttachmentCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectComposeAttachmentCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryComposeAttachments queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryComposeAttachments(
	ctx context.Context,
	f composeType.AttachmentFilter,
) (_ []*composeType.Attachment, more bool, err error) {
	var (
		ok bool

		set         = make([]*composeType.Attachment, 0, DefaultSliceCapacity)
		res         *composeType.Attachment
		aux         *auxComposeAttachment
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.ComposeAttachment != nil {
		// extended filter set
		tExpr, f, err = s.Filters.ComposeAttachment(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ComposeAttachmentFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for ComposeAttachment: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableComposeAttachmentFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := composeAttachmentSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableComposeAttachmentFields()); err != nil {
		err = fmt.Errorf("could generate order expression for ComposeAttachment: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query ComposeAttachment: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query ComposeAttachment: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query ComposeAttachment: %w", err)
			return
		}

		aux = new(auxComposeAttachment)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for ComposeAttachment: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode ComposeAttachment: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupComposeAttachmentByID
//
// This function is auto-generated
func (s *Store) LookupComposeAttachmentByID(ctx context.Context, id uint64) (_ *composeType.Attachment, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposeAttachment)
		lookup = composeAttachmentSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableComposeAttachmentFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableComposeAttachmentFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"id":         "id",
		"kind":       "kind",
		"name":       "name",
		"owner_id":   "owner_id",
		"ownerid":    "owner_id",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
	}
}

// collectComposeAttachmentCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectComposeAttachmentCursorValues(res *composeType.Attachment, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "ownerID":
					cur.Set(c.Column, res.OwnerID, c.Descending)
				case "kind":
					cur.Set(c.Column, res.Kind, c.Descending)
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkComposeAttachmentConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkComposeAttachmentConstraints(ctx context.Context, res *composeType.Attachment) (err error) {
	return nil
}

// CreateComposeChart creates one or more rows in composeChart collection
//
// This function is auto-generated
func (s *Store) CreateComposeChart(ctx context.Context, rr ...*composeType.Chart) (err error) {
	for i := range rr {
		if err = s.checkComposeChartConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeChartInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateComposeChart updates one or more existing entries in composeChart collection
//
// This function is auto-generated
func (s *Store) UpdateComposeChart(ctx context.Context, rr ...*composeType.Chart) (err error) {
	for i := range rr {
		if err = s.checkComposeChartConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeChartUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertComposeChart updates one or more existing entries in composeChart collection
//
// This function is auto-generated
func (s *Store) UpsertComposeChart(ctx context.Context, rr ...*composeType.Chart) (err error) {
	for i := range rr {
		if err = s.checkComposeChartConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeChartUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteComposeChart Deletes one or more entries from composeChart collection
//
// This function is auto-generated
func (s *Store) DeleteComposeChart(ctx context.Context, rr ...*composeType.Chart) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, composeChartDeleteQuery(s.Dialect, composeChartPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteComposeChartByID deletes single entry from composeChart collection
//
// This function is auto-generated
func (s *Store) DeleteComposeChartByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, composeChartDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateComposeCharts Deletes all rows from the composeChart collection
func (s Store) TruncateComposeCharts(ctx context.Context) error {
	return s.Exec(ctx, composeChartTruncateQuery(s.Dialect))
}

// SearchComposeCharts returns (filtered) set of ComposeCharts
//
// This function is auto-generated
func (s *Store) SearchComposeCharts(ctx context.Context, f composeType.ChartFilter) (set composeType.ChartSet, _ composeType.ChartFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfComposeCharts fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfComposeCharts(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfComposeCharts collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfComposeCharts(
	ctx context.Context,
	filter composeType.ChartFilter,
	sort filter.SortExprSet,
) (set []*composeType.Chart, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*composeType.Chart

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter composeType.ChartFilter
	)

	set = make([]*composeType.Chart, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryComposeCharts(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectComposeChartCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectComposeChartCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectComposeChartCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryComposeCharts queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryComposeCharts(
	ctx context.Context,
	f composeType.ChartFilter,
) (_ []*composeType.Chart, more bool, err error) {
	var (
		ok bool

		set         = make([]*composeType.Chart, 0, DefaultSliceCapacity)
		res         *composeType.Chart
		aux         *auxComposeChart
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.ComposeChart != nil {
		// extended filter set
		tExpr, f, err = s.Filters.ComposeChart(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ComposeChartFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for ComposeChart: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableComposeChartFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := composeChartSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableComposeChartFields()); err != nil {
		err = fmt.Errorf("could generate order expression for ComposeChart: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query ComposeChart: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query ComposeChart: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query ComposeChart: %w", err)
			return
		}

		aux = new(auxComposeChart)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for ComposeChart: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode ComposeChart: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupComposeChartByID searches for compose chart by ID
//
// It returns compose chart even if deleted
//
// This function is auto-generated
func (s *Store) LookupComposeChartByID(ctx context.Context, id uint64) (_ *composeType.Chart, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposeChart)
		lookup = composeChartSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupComposeChartByNamespaceIDHandle searches for compose chart by handle (case-insensitive)
//
// This function is auto-generated
func (s *Store) LookupComposeChartByNamespaceIDHandle(ctx context.Context, namespaceID uint64, handle string) (_ *composeType.Chart, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposeChart)
		lookup = composeChartSelectQuery(s.Dialect).Where(
			goqu.I("rel_namespace").Eq(namespaceID),
			s.Functions.LOWER(goqu.I("handle")).Eq(strings.ToLower(handle)),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableComposeChartFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableComposeChartFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"handle":     "handle",
		"id":         "id",
		"name":       "name",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
	}
}

// collectComposeChartCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectComposeChartCursorValues(res *composeType.Chart, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkComposeChartConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkComposeChartConstraints(ctx context.Context, res *composeType.Chart) (err error) {
	return nil
}

// CreateComposeModule creates one or more rows in composeModule collection
//
// This function is auto-generated
func (s *Store) CreateComposeModule(ctx context.Context, rr ...*composeType.Module) (err error) {
	for i := range rr {
		if err = s.checkComposeModuleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeModuleInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateComposeModule updates one or more existing entries in composeModule collection
//
// This function is auto-generated
func (s *Store) UpdateComposeModule(ctx context.Context, rr ...*composeType.Module) (err error) {
	for i := range rr {
		if err = s.checkComposeModuleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeModuleUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertComposeModule updates one or more existing entries in composeModule collection
//
// This function is auto-generated
func (s *Store) UpsertComposeModule(ctx context.Context, rr ...*composeType.Module) (err error) {
	for i := range rr {
		if err = s.checkComposeModuleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeModuleUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteComposeModule Deletes one or more entries from composeModule collection
//
// This function is auto-generated
func (s *Store) DeleteComposeModule(ctx context.Context, rr ...*composeType.Module) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, composeModuleDeleteQuery(s.Dialect, composeModulePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteComposeModuleByID deletes single entry from composeModule collection
//
// This function is auto-generated
func (s *Store) DeleteComposeModuleByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, composeModuleDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateComposeModules Deletes all rows from the composeModule collection
func (s Store) TruncateComposeModules(ctx context.Context) error {
	return s.Exec(ctx, composeModuleTruncateQuery(s.Dialect))
}

// SearchComposeModules returns (filtered) set of ComposeModules
//
// This function is auto-generated
func (s *Store) SearchComposeModules(ctx context.Context, f composeType.ModuleFilter) (set composeType.ModuleSet, _ composeType.ModuleFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfComposeModules fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfComposeModules(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfComposeModules collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfComposeModules(
	ctx context.Context,
	filter composeType.ModuleFilter,
	sort filter.SortExprSet,
) (set []*composeType.Module, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*composeType.Module

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter composeType.ModuleFilter
	)

	set = make([]*composeType.Module, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryComposeModules(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectComposeModuleCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectComposeModuleCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectComposeModuleCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryComposeModules queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryComposeModules(
	ctx context.Context,
	f composeType.ModuleFilter,
) (_ []*composeType.Module, more bool, err error) {
	var (
		ok bool

		set         = make([]*composeType.Module, 0, DefaultSliceCapacity)
		res         *composeType.Module
		aux         *auxComposeModule
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.ComposeModule != nil {
		// extended filter set
		tExpr, f, err = s.Filters.ComposeModule(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ComposeModuleFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for ComposeModule: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableComposeModuleFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := composeModuleSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableComposeModuleFields()); err != nil {
		err = fmt.Errorf("could generate order expression for ComposeModule: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query ComposeModule: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query ComposeModule: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query ComposeModule: %w", err)
			return
		}

		aux = new(auxComposeModule)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for ComposeModule: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode ComposeModule: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupComposeModuleByNamespaceIDHandle searches for compose module by handle (case-insensitive)
//
// This function is auto-generated
func (s *Store) LookupComposeModuleByNamespaceIDHandle(ctx context.Context, namespaceID uint64, handle string) (_ *composeType.Module, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposeModule)
		lookup = composeModuleSelectQuery(s.Dialect).Where(
			goqu.I("rel_namespace").Eq(namespaceID),
			s.Functions.LOWER(goqu.I("handle")).Eq(strings.ToLower(handle)),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupComposeModuleByNamespaceIDName searches for compose module by name (case-insensitive)
//
// This function is auto-generated
func (s *Store) LookupComposeModuleByNamespaceIDName(ctx context.Context, namespaceID uint64, name string) (_ *composeType.Module, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposeModule)
		lookup = composeModuleSelectQuery(s.Dialect).Where(
			goqu.I("rel_namespace").Eq(namespaceID),
			goqu.I("name").Eq(name),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupComposeModuleByID searches for compose module by ID
//
// It returns compose module even if deleted
//
// This function is auto-generated
func (s *Store) LookupComposeModuleByID(ctx context.Context, id uint64) (_ *composeType.Module, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposeModule)
		lookup = composeModuleSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableComposeModuleFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableComposeModuleFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"handle":     "handle",
		"id":         "id",
		"name":       "name",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
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
//
// This function is auto-generated
func (s *Store) collectComposeModuleCursorValues(res *composeType.Module, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkComposeModuleConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkComposeModuleConstraints(ctx context.Context, res *composeType.Module) (err error) {
	err = func() (err error) {

		if res.NamespaceID == 0 {
			// skip check on empty values
			return nil
		}

		// handling string type as default
		if len(res.Handle) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupComposeModuleByNamespaceIDHandle(ctx, res.NamespaceID, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	return nil
}

// CreateComposeModuleField creates one or more rows in composeModuleField collection
//
// This function is auto-generated
func (s *Store) CreateComposeModuleField(ctx context.Context, rr ...*composeType.ModuleField) (err error) {
	for i := range rr {
		if err = s.checkComposeModuleFieldConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeModuleFieldInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateComposeModuleField updates one or more existing entries in composeModuleField collection
//
// This function is auto-generated
func (s *Store) UpdateComposeModuleField(ctx context.Context, rr ...*composeType.ModuleField) (err error) {
	for i := range rr {
		if err = s.checkComposeModuleFieldConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeModuleFieldUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertComposeModuleField updates one or more existing entries in composeModuleField collection
//
// This function is auto-generated
func (s *Store) UpsertComposeModuleField(ctx context.Context, rr ...*composeType.ModuleField) (err error) {
	for i := range rr {
		if err = s.checkComposeModuleFieldConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeModuleFieldUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteComposeModuleField Deletes one or more entries from composeModuleField collection
//
// This function is auto-generated
func (s *Store) DeleteComposeModuleField(ctx context.Context, rr ...*composeType.ModuleField) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, composeModuleFieldDeleteQuery(s.Dialect, composeModuleFieldPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteComposeModuleFieldByID deletes single entry from composeModuleField collection
//
// This function is auto-generated
func (s *Store) DeleteComposeModuleFieldByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, composeModuleFieldDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateComposeModuleFields Deletes all rows from the composeModuleField collection
func (s Store) TruncateComposeModuleFields(ctx context.Context) error {
	return s.Exec(ctx, composeModuleFieldTruncateQuery(s.Dialect))
}

// SearchComposeModuleFields returns (filtered) set of ComposeModuleFields
//
// This function is auto-generated
func (s *Store) SearchComposeModuleFields(ctx context.Context, f composeType.ModuleFieldFilter) (set composeType.ModuleFieldSet, _ composeType.ModuleFieldFilter, err error) {

	set, _, err = s.QueryComposeModuleFields(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryComposeModuleFields queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryComposeModuleFields(
	ctx context.Context,
	f composeType.ModuleFieldFilter,
) (_ []*composeType.ModuleField, more bool, err error) {
	var (
		set         = make([]*composeType.ModuleField, 0, DefaultSliceCapacity)
		res         *composeType.ModuleField
		aux         *auxComposeModuleField
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.ComposeModuleField != nil {
		// extended filter set
		tExpr, f, err = s.Filters.ComposeModuleField(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ComposeModuleFieldFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for ComposeModuleField: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := composeModuleFieldSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query ComposeModuleField: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query ComposeModuleField: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query ComposeModuleField: %w", err)
			return
		}

		aux = new(auxComposeModuleField)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for ComposeModuleField: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode ComposeModuleField: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, false, err

}

// LookupComposeModuleFieldByModuleIDName searches for compose module field by name (case-insensitive)
//
// This function is auto-generated
func (s *Store) LookupComposeModuleFieldByModuleIDName(ctx context.Context, moduleID uint64, name string) (_ *composeType.ModuleField, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposeModuleField)
		lookup = composeModuleFieldSelectQuery(s.Dialect).Where(
			goqu.I("rel_module").Eq(moduleID),
			goqu.I("name").Eq(name),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupComposeModuleFieldByID searches for compose module field by ID
//
// This function is auto-generated
func (s *Store) LookupComposeModuleFieldByID(ctx context.Context, id uint64) (_ *composeType.ModuleField, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposeModuleField)
		lookup = composeModuleFieldSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableComposeModuleFieldFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableComposeModuleFieldFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"id":         "id",
		"kind":       "kind",
		"label":      "label",
		"name":       "name",
		"place":      "place",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
	}
}

// collectComposeModuleFieldCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectComposeModuleFieldCursorValues(res *composeType.ModuleField, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "place":
					cur.Set(c.Column, res.Place, c.Descending)
				case "kind":
					cur.Set(c.Column, res.Kind, c.Descending)
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
				case "label":
					cur.Set(c.Column, res.Label, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkComposeModuleFieldConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkComposeModuleFieldConstraints(ctx context.Context, res *composeType.ModuleField) (err error) {
	err = func() (err error) {

		if res.ModuleID == 0 {
			// skip check on empty values
			return nil
		}

		// handling string type as default
		if len(res.Name) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupComposeModuleFieldByModuleIDName(ctx, res.ModuleID, res.Name)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	return nil
}

// CreateComposeNamespace creates one or more rows in composeNamespace collection
//
// This function is auto-generated
func (s *Store) CreateComposeNamespace(ctx context.Context, rr ...*composeType.Namespace) (err error) {
	for i := range rr {
		if err = s.checkComposeNamespaceConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeNamespaceInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateComposeNamespace updates one or more existing entries in composeNamespace collection
//
// This function is auto-generated
func (s *Store) UpdateComposeNamespace(ctx context.Context, rr ...*composeType.Namespace) (err error) {
	for i := range rr {
		if err = s.checkComposeNamespaceConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeNamespaceUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertComposeNamespace updates one or more existing entries in composeNamespace collection
//
// This function is auto-generated
func (s *Store) UpsertComposeNamespace(ctx context.Context, rr ...*composeType.Namespace) (err error) {
	for i := range rr {
		if err = s.checkComposeNamespaceConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composeNamespaceUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteComposeNamespace Deletes one or more entries from composeNamespace collection
//
// This function is auto-generated
func (s *Store) DeleteComposeNamespace(ctx context.Context, rr ...*composeType.Namespace) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, composeNamespaceDeleteQuery(s.Dialect, composeNamespacePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteComposeNamespaceByID deletes single entry from composeNamespace collection
//
// This function is auto-generated
func (s *Store) DeleteComposeNamespaceByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, composeNamespaceDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateComposeNamespaces Deletes all rows from the composeNamespace collection
func (s Store) TruncateComposeNamespaces(ctx context.Context) error {
	return s.Exec(ctx, composeNamespaceTruncateQuery(s.Dialect))
}

// SearchComposeNamespaces returns (filtered) set of ComposeNamespaces
//
// This function is auto-generated
func (s *Store) SearchComposeNamespaces(ctx context.Context, f composeType.NamespaceFilter) (set composeType.NamespaceSet, _ composeType.NamespaceFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfComposeNamespaces fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfComposeNamespaces(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfComposeNamespaces collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfComposeNamespaces(
	ctx context.Context,
	filter composeType.NamespaceFilter,
	sort filter.SortExprSet,
) (set []*composeType.Namespace, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*composeType.Namespace

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter composeType.NamespaceFilter
	)

	set = make([]*composeType.Namespace, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryComposeNamespaces(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectComposeNamespaceCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectComposeNamespaceCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectComposeNamespaceCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryComposeNamespaces queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryComposeNamespaces(
	ctx context.Context,
	f composeType.NamespaceFilter,
) (_ []*composeType.Namespace, more bool, err error) {
	var (
		ok bool

		set         = make([]*composeType.Namespace, 0, DefaultSliceCapacity)
		res         *composeType.Namespace
		aux         *auxComposeNamespace
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.ComposeNamespace != nil {
		// extended filter set
		tExpr, f, err = s.Filters.ComposeNamespace(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ComposeNamespaceFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for ComposeNamespace: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableComposeNamespaceFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := composeNamespaceSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableComposeNamespaceFields()); err != nil {
		err = fmt.Errorf("could generate order expression for ComposeNamespace: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query ComposeNamespace: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query ComposeNamespace: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query ComposeNamespace: %w", err)
			return
		}

		aux = new(auxComposeNamespace)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for ComposeNamespace: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode ComposeNamespace: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupComposeNamespaceBySlug searches for namespace by slug (case-insensitive)
//
// This function is auto-generated
func (s *Store) LookupComposeNamespaceBySlug(ctx context.Context, slug string) (_ *composeType.Namespace, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposeNamespace)
		lookup = composeNamespaceSelectQuery(s.Dialect).Where(
			goqu.I("slug").Eq(slug),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupComposeNamespaceByID searches for compose namespace by ID
//
// It returns compose namespace even if deleted
//
// This function is auto-generated
func (s *Store) LookupComposeNamespaceByID(ctx context.Context, id uint64) (_ *composeType.Namespace, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposeNamespace)
		lookup = composeNamespaceSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableComposeNamespaceFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableComposeNamespaceFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"id":         "id",
		"name":       "name",
		"slug":       "slug",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
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
//
// This function is auto-generated
func (s *Store) collectComposeNamespaceCursorValues(res *composeType.Namespace, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "slug":
					cur.Set(c.Column, res.Slug, c.Descending)
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkComposeNamespaceConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkComposeNamespaceConstraints(ctx context.Context, res *composeType.Namespace) (err error) {
	err = func() (err error) {

		// handling string type as default
		if len(res.Slug) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupComposeNamespaceBySlug(ctx, res.Slug)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	return nil
}

// CreateComposePage creates one or more rows in composePage collection
//
// This function is auto-generated
func (s *Store) CreateComposePage(ctx context.Context, rr ...*composeType.Page) (err error) {
	for i := range rr {
		if err = s.checkComposePageConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composePageInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateComposePage updates one or more existing entries in composePage collection
//
// This function is auto-generated
func (s *Store) UpdateComposePage(ctx context.Context, rr ...*composeType.Page) (err error) {
	for i := range rr {
		if err = s.checkComposePageConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composePageUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertComposePage updates one or more existing entries in composePage collection
//
// This function is auto-generated
func (s *Store) UpsertComposePage(ctx context.Context, rr ...*composeType.Page) (err error) {
	for i := range rr {
		if err = s.checkComposePageConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, composePageUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteComposePage Deletes one or more entries from composePage collection
//
// This function is auto-generated
func (s *Store) DeleteComposePage(ctx context.Context, rr ...*composeType.Page) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, composePageDeleteQuery(s.Dialect, composePagePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteComposePageByID deletes single entry from composePage collection
//
// This function is auto-generated
func (s *Store) DeleteComposePageByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, composePageDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateComposePages Deletes all rows from the composePage collection
func (s Store) TruncateComposePages(ctx context.Context) error {
	return s.Exec(ctx, composePageTruncateQuery(s.Dialect))
}

// SearchComposePages returns (filtered) set of ComposePages
//
// This function is auto-generated
func (s *Store) SearchComposePages(ctx context.Context, f composeType.PageFilter) (set composeType.PageSet, _ composeType.PageFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfComposePages fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfComposePages(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
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
//
// This function is auto-generated
func (s *Store) fetchFullPageOfComposePages(
	ctx context.Context,
	filter composeType.PageFilter,
	sort filter.SortExprSet,
) (set []*composeType.Page, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*composeType.Page

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter composeType.PageFilter
	)

	set = make([]*composeType.Page, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryComposePages(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectComposePageCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectComposePageCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectComposePageCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryComposePages queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryComposePages(
	ctx context.Context,
	f composeType.PageFilter,
) (_ []*composeType.Page, more bool, err error) {
	var (
		ok bool

		set         = make([]*composeType.Page, 0, DefaultSliceCapacity)
		res         *composeType.Page
		aux         *auxComposePage
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.ComposePage != nil {
		// extended filter set
		tExpr, f, err = s.Filters.ComposePage(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ComposePageFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for ComposePage: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableComposePageFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := composePageSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableComposePageFields()); err != nil {
		err = fmt.Errorf("could generate order expression for ComposePage: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query ComposePage: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query ComposePage: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query ComposePage: %w", err)
			return
		}

		aux = new(auxComposePage)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for ComposePage: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode ComposePage: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupComposePageByNamespaceIDHandle searches for page by handle (case-insensitive)
//
// This function is auto-generated
func (s *Store) LookupComposePageByNamespaceIDHandle(ctx context.Context, namespaceID uint64, handle string) (_ *composeType.Page, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposePage)
		lookup = composePageSelectQuery(s.Dialect).Where(
			goqu.I("rel_namespace").Eq(namespaceID),
			s.Functions.LOWER(goqu.I("handle")).Eq(strings.ToLower(handle)),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupComposePageByNamespaceIDModuleID searches for page by moduleID
//
// This function is auto-generated
func (s *Store) LookupComposePageByNamespaceIDModuleID(ctx context.Context, namespaceID uint64, moduleID uint64) (_ *composeType.Page, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposePage)
		lookup = composePageSelectQuery(s.Dialect).Where(
			goqu.I("rel_namespace").Eq(namespaceID),
			goqu.I("rel_module").Eq(moduleID),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupComposePageByID searches for compose page by ID
//
// It returns compose page even if deleted
//
// This function is auto-generated
func (s *Store) LookupComposePageByID(ctx context.Context, id uint64) (_ *composeType.Page, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxComposePage)
		lookup = composePageSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableComposePageFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableComposePageFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"handle":     "handle",
		"id":         "id",
		"self_id":    "self_id",
		"selfid":     "self_id",
		"title":      "title",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
		"weight":     "weight",
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
//
// This function is auto-generated
func (s *Store) collectComposePageCursorValues(res *composeType.Page, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "title":
					cur.Set(c.Column, res.Title, c.Descending)
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "selfID":
					cur.Set(c.Column, res.SelfID, c.Descending)
				case "weight":
					cur.Set(c.Column, res.Weight, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkComposePageConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkComposePageConstraints(ctx context.Context, res *composeType.Page) (err error) {
	return nil
}

// CreateCredential creates one or more rows in credential collection
//
// This function is auto-generated
func (s *Store) CreateCredential(ctx context.Context, rr ...*systemType.Credential) (err error) {
	for i := range rr {
		if err = s.checkCredentialConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, credentialInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateCredential updates one or more existing entries in credential collection
//
// This function is auto-generated
func (s *Store) UpdateCredential(ctx context.Context, rr ...*systemType.Credential) (err error) {
	for i := range rr {
		if err = s.checkCredentialConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, credentialUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertCredential updates one or more existing entries in credential collection
//
// This function is auto-generated
func (s *Store) UpsertCredential(ctx context.Context, rr ...*systemType.Credential) (err error) {
	for i := range rr {
		if err = s.checkCredentialConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, credentialUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteCredential Deletes one or more entries from credential collection
//
// This function is auto-generated
func (s *Store) DeleteCredential(ctx context.Context, rr ...*systemType.Credential) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, credentialDeleteQuery(s.Dialect, credentialPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteCredentialByID deletes single entry from credential collection
//
// This function is auto-generated
func (s *Store) DeleteCredentialByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, credentialDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateCredentials Deletes all rows from the credential collection
func (s Store) TruncateCredentials(ctx context.Context) error {
	return s.Exec(ctx, credentialTruncateQuery(s.Dialect))
}

// SearchCredentials returns (filtered) set of Credentials
//
// This function is auto-generated
func (s *Store) SearchCredentials(ctx context.Context, f systemType.CredentialFilter) (set systemType.CredentialSet, _ systemType.CredentialFilter, err error) {

	set, _, err = s.QueryCredentials(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryCredentials queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryCredentials(
	ctx context.Context,
	f systemType.CredentialFilter,
) (_ []*systemType.Credential, more bool, err error) {
	var (
		set         = make([]*systemType.Credential, 0, DefaultSliceCapacity)
		res         *systemType.Credential
		aux         *auxCredential
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.Credential != nil {
		// extended filter set
		tExpr, f, err = s.Filters.Credential(s, f)
	} else {
		// using generated filter
		tExpr, f, err = CredentialFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for Credential: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := credentialSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query Credential: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query Credential: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query Credential: %w", err)
			return
		}

		aux = new(auxCredential)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for Credential: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode Credential: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, false, err

}

// LookupCredentialByID searches for credentials by ID
//
// It returns credentials even if deleted
//
// This function is auto-generated
func (s *Store) LookupCredentialByID(ctx context.Context, id uint64) (_ *systemType.Credential, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxCredential)
		lookup = credentialSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableCredentialFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableCredentialFields() map[string]string {
	return map[string]string{
		"created_at":   "created_at",
		"createdat":    "created_at",
		"deleted_at":   "deleted_at",
		"deletedat":    "deleted_at",
		"expires_at":   "expires_at",
		"expiresat":    "expires_at",
		"id":           "id",
		"last_used_at": "last_used_at",
		"lastusedat":   "last_used_at",
		"updated_at":   "updated_at",
		"updatedat":    "updated_at",
	}
}

// collectCredentialCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectCredentialCursorValues(res *systemType.Credential, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				case "lastUsedAt":
					cur.Set(c.Column, res.LastUsedAt, c.Descending)
				case "expiresAt":
					cur.Set(c.Column, res.ExpiresAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkCredentialConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkCredentialConstraints(ctx context.Context, res *systemType.Credential) (err error) {
	return nil
}

// CreateDalConnection creates one or more rows in dalConnection collection
//
// This function is auto-generated
func (s *Store) CreateDalConnection(ctx context.Context, rr ...*systemType.DalConnection) (err error) {
	for i := range rr {
		if err = s.checkDalConnectionConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, dalConnectionInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateDalConnection updates one or more existing entries in dalConnection collection
//
// This function is auto-generated
func (s *Store) UpdateDalConnection(ctx context.Context, rr ...*systemType.DalConnection) (err error) {
	for i := range rr {
		if err = s.checkDalConnectionConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, dalConnectionUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertDalConnection updates one or more existing entries in dalConnection collection
//
// This function is auto-generated
func (s *Store) UpsertDalConnection(ctx context.Context, rr ...*systemType.DalConnection) (err error) {
	for i := range rr {
		if err = s.checkDalConnectionConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, dalConnectionUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteDalConnection Deletes one or more entries from dalConnection collection
//
// This function is auto-generated
func (s *Store) DeleteDalConnection(ctx context.Context, rr ...*systemType.DalConnection) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, dalConnectionDeleteQuery(s.Dialect, dalConnectionPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteDalConnectionByID deletes single entry from dalConnection collection
//
// This function is auto-generated
func (s *Store) DeleteDalConnectionByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, dalConnectionDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateDalConnections Deletes all rows from the dalConnection collection
func (s Store) TruncateDalConnections(ctx context.Context) error {
	return s.Exec(ctx, dalConnectionTruncateQuery(s.Dialect))
}

// SearchDalConnections returns (filtered) set of DalConnections
//
// This function is auto-generated
func (s *Store) SearchDalConnections(ctx context.Context, f systemType.DalConnectionFilter) (set systemType.DalConnectionSet, _ systemType.DalConnectionFilter, err error) {

	set, _, err = s.QueryDalConnections(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryDalConnections queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryDalConnections(
	ctx context.Context,
	f systemType.DalConnectionFilter,
) (_ []*systemType.DalConnection, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.DalConnection, 0, DefaultSliceCapacity)
		res         *systemType.DalConnection
		aux         *auxDalConnection
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.DalConnection != nil {
		// extended filter set
		tExpr, f, err = s.Filters.DalConnection(s, f)
	} else {
		// using generated filter
		tExpr, f, err = DalConnectionFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for DalConnection: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := dalConnectionSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query DalConnection: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query DalConnection: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query DalConnection: %w", err)
			return
		}

		aux = new(auxDalConnection)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for DalConnection: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode DalConnection: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, false, err

}

// LookupDalConnectionByID searches for connection by ID
//
// It returns connection even if deleted or suspended
//
// This function is auto-generated
func (s *Store) LookupDalConnectionByID(ctx context.Context, id uint64) (_ *systemType.DalConnection, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxDalConnection)
		lookup = dalConnectionSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupDalConnectionByHandle searches for connection by handle
//
// It returns only valid connection (not deleted)
//
// This function is auto-generated
func (s *Store) LookupDalConnectionByHandle(ctx context.Context, handle string) (_ *systemType.DalConnection, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxDalConnection)
		lookup = dalConnectionSelectQuery(s.Dialect).Where(
			s.Functions.LOWER(goqu.I("handle")).Eq(strings.ToLower(handle)),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableDalConnectionFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableDalConnectionFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"handle":     "handle",
		"id":         "id",
		"type":       "type",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
	}
}

// collectDalConnectionCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectDalConnectionCursorValues(res *systemType.DalConnection, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "type":
					cur.Set(c.Column, res.Type, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkDalConnectionConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkDalConnectionConstraints(ctx context.Context, res *systemType.DalConnection) (err error) {
	err = func() (err error) {

		// handling string type as default
		if len(res.Handle) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupDalConnectionByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	return nil
}

// CreateDalSensitivityLevel creates one or more rows in dalSensitivityLevel collection
//
// This function is auto-generated
func (s *Store) CreateDalSensitivityLevel(ctx context.Context, rr ...*systemType.DalSensitivityLevel) (err error) {
	for i := range rr {
		if err = s.checkDalSensitivityLevelConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, dalSensitivityLevelInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateDalSensitivityLevel updates one or more existing entries in dalSensitivityLevel collection
//
// This function is auto-generated
func (s *Store) UpdateDalSensitivityLevel(ctx context.Context, rr ...*systemType.DalSensitivityLevel) (err error) {
	for i := range rr {
		if err = s.checkDalSensitivityLevelConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, dalSensitivityLevelUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertDalSensitivityLevel updates one or more existing entries in dalSensitivityLevel collection
//
// This function is auto-generated
func (s *Store) UpsertDalSensitivityLevel(ctx context.Context, rr ...*systemType.DalSensitivityLevel) (err error) {
	for i := range rr {
		if err = s.checkDalSensitivityLevelConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, dalSensitivityLevelUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteDalSensitivityLevel Deletes one or more entries from dalSensitivityLevel collection
//
// This function is auto-generated
func (s *Store) DeleteDalSensitivityLevel(ctx context.Context, rr ...*systemType.DalSensitivityLevel) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, dalSensitivityLevelDeleteQuery(s.Dialect, dalSensitivityLevelPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteDalSensitivityLevelByID deletes single entry from dalSensitivityLevel collection
//
// This function is auto-generated
func (s *Store) DeleteDalSensitivityLevelByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, dalSensitivityLevelDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateDalSensitivityLevels Deletes all rows from the dalSensitivityLevel collection
func (s Store) TruncateDalSensitivityLevels(ctx context.Context) error {
	return s.Exec(ctx, dalSensitivityLevelTruncateQuery(s.Dialect))
}

// SearchDalSensitivityLevels returns (filtered) set of DalSensitivityLevels
//
// This function is auto-generated
func (s *Store) SearchDalSensitivityLevels(ctx context.Context, f systemType.DalSensitivityLevelFilter) (set systemType.DalSensitivityLevelSet, _ systemType.DalSensitivityLevelFilter, err error) {

	set, _, err = s.QueryDalSensitivityLevels(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryDalSensitivityLevels queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryDalSensitivityLevels(
	ctx context.Context,
	f systemType.DalSensitivityLevelFilter,
) (_ []*systemType.DalSensitivityLevel, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.DalSensitivityLevel, 0, DefaultSliceCapacity)
		res         *systemType.DalSensitivityLevel
		aux         *auxDalSensitivityLevel
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.DalSensitivityLevel != nil {
		// extended filter set
		tExpr, f, err = s.Filters.DalSensitivityLevel(s, f)
	} else {
		// using generated filter
		tExpr, f, err = DalSensitivityLevelFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for DalSensitivityLevel: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := dalSensitivityLevelSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query DalSensitivityLevel: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query DalSensitivityLevel: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query DalSensitivityLevel: %w", err)
			return
		}

		aux = new(auxDalSensitivityLevel)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for DalSensitivityLevel: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode DalSensitivityLevel: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, false, err

}

// LookupDalSensitivityLevelByID searches for user by ID
//
// It returns user even if deleted or suspended
//
// This function is auto-generated
func (s *Store) LookupDalSensitivityLevelByID(ctx context.Context, id uint64) (_ *systemType.DalSensitivityLevel, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxDalSensitivityLevel)
		lookup = dalSensitivityLevelSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableDalSensitivityLevelFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableDalSensitivityLevelFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"handle":     "handle",
		"id":         "id",
		"level":      "level",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
	}
}

// collectDalSensitivityLevelCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectDalSensitivityLevelCursorValues(res *systemType.DalSensitivityLevel, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "level":
					cur.Set(c.Column, res.Level, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkDalSensitivityLevelConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkDalSensitivityLevelConstraints(ctx context.Context, res *systemType.DalSensitivityLevel) (err error) {
	return nil
}

// CreateDataPrivacyRequest creates one or more rows in dataPrivacyRequest collection
//
// This function is auto-generated
func (s *Store) CreateDataPrivacyRequest(ctx context.Context, rr ...*systemType.DataPrivacyRequest) (err error) {
	for i := range rr {
		if err = s.checkDataPrivacyRequestConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, dataPrivacyRequestInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateDataPrivacyRequest updates one or more existing entries in dataPrivacyRequest collection
//
// This function is auto-generated
func (s *Store) UpdateDataPrivacyRequest(ctx context.Context, rr ...*systemType.DataPrivacyRequest) (err error) {
	for i := range rr {
		if err = s.checkDataPrivacyRequestConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, dataPrivacyRequestUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertDataPrivacyRequest updates one or more existing entries in dataPrivacyRequest collection
//
// This function is auto-generated
func (s *Store) UpsertDataPrivacyRequest(ctx context.Context, rr ...*systemType.DataPrivacyRequest) (err error) {
	for i := range rr {
		if err = s.checkDataPrivacyRequestConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, dataPrivacyRequestUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteDataPrivacyRequest Deletes one or more entries from dataPrivacyRequest collection
//
// This function is auto-generated
func (s *Store) DeleteDataPrivacyRequest(ctx context.Context, rr ...*systemType.DataPrivacyRequest) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, dataPrivacyRequestDeleteQuery(s.Dialect, dataPrivacyRequestPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteDataPrivacyRequestByID deletes single entry from dataPrivacyRequest collection
//
// This function is auto-generated
func (s *Store) DeleteDataPrivacyRequestByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, dataPrivacyRequestDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateDataPrivacyRequests Deletes all rows from the dataPrivacyRequest collection
func (s Store) TruncateDataPrivacyRequests(ctx context.Context) error {
	return s.Exec(ctx, dataPrivacyRequestTruncateQuery(s.Dialect))
}

// SearchDataPrivacyRequests returns (filtered) set of DataPrivacyRequests
//
// This function is auto-generated
func (s *Store) SearchDataPrivacyRequests(ctx context.Context, f systemType.DataPrivacyRequestFilter) (set systemType.DataPrivacyRequestSet, _ systemType.DataPrivacyRequestFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfDataPrivacyRequests fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfDataPrivacyRequests(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfDataPrivacyRequests collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfDataPrivacyRequests(
	ctx context.Context,
	filter systemType.DataPrivacyRequestFilter,
	sort filter.SortExprSet,
) (set []*systemType.DataPrivacyRequest, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.DataPrivacyRequest

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.DataPrivacyRequestFilter
	)

	set = make([]*systemType.DataPrivacyRequest, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryDataPrivacyRequests(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectDataPrivacyRequestCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectDataPrivacyRequestCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectDataPrivacyRequestCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryDataPrivacyRequests queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryDataPrivacyRequests(
	ctx context.Context,
	f systemType.DataPrivacyRequestFilter,
) (_ []*systemType.DataPrivacyRequest, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.DataPrivacyRequest, 0, DefaultSliceCapacity)
		res         *systemType.DataPrivacyRequest
		aux         *auxDataPrivacyRequest
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.DataPrivacyRequest != nil {
		// extended filter set
		tExpr, f, err = s.Filters.DataPrivacyRequest(s, f)
	} else {
		// using generated filter
		tExpr, f, err = DataPrivacyRequestFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for DataPrivacyRequest: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableDataPrivacyRequestFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := dataPrivacyRequestSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableDataPrivacyRequestFields()); err != nil {
		err = fmt.Errorf("could generate order expression for DataPrivacyRequest: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query DataPrivacyRequest: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query DataPrivacyRequest: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query DataPrivacyRequest: %w", err)
			return
		}

		aux = new(auxDataPrivacyRequest)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for DataPrivacyRequest: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode DataPrivacyRequest: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupDataPrivacyRequestByID searches for data privacy request by ID
//
// It returns data privacy request even if deleted
//
// This function is auto-generated
func (s *Store) LookupDataPrivacyRequestByID(ctx context.Context, id uint64) (_ *systemType.DataPrivacyRequest, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxDataPrivacyRequest)
		lookup = dataPrivacyRequestSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableDataPrivacyRequestFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableDataPrivacyRequestFields() map[string]string {
	return map[string]string{
		"completed_at": "completed_at",
		"completedat":  "completed_at",
		"created_at":   "created_at",
		"createdat":    "created_at",
		"deleted_at":   "deleted_at",
		"deletedat":    "deleted_at",
		"id":           "id",
		"kind":         "kind",
		"requested_at": "requested_at",
		"requestedat":  "requested_at",
		"status":       "status",
		"updated_at":   "updated_at",
		"updatedat":    "updated_at",
	}
}

// collectDataPrivacyRequestCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectDataPrivacyRequestCursorValues(res *systemType.DataPrivacyRequest, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "kind":
					cur.Set(c.Column, res.Kind, c.Descending)
				case "status":
					cur.Set(c.Column, res.Status, c.Descending)
				case "requestedAt":
					cur.Set(c.Column, res.RequestedAt, c.Descending)
				case "completedAt":
					cur.Set(c.Column, res.CompletedAt, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkDataPrivacyRequestConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkDataPrivacyRequestConstraints(ctx context.Context, res *systemType.DataPrivacyRequest) (err error) {
	return nil
}

// CreateDataPrivacyRequestComment creates one or more rows in dataPrivacyRequestComment collection
//
// This function is auto-generated
func (s *Store) CreateDataPrivacyRequestComment(ctx context.Context, rr ...*systemType.DataPrivacyRequestComment) (err error) {
	for i := range rr {
		if err = s.checkDataPrivacyRequestCommentConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, dataPrivacyRequestCommentInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateDataPrivacyRequestComment updates one or more existing entries in dataPrivacyRequestComment collection
//
// This function is auto-generated
func (s *Store) UpdateDataPrivacyRequestComment(ctx context.Context, rr ...*systemType.DataPrivacyRequestComment) (err error) {
	for i := range rr {
		if err = s.checkDataPrivacyRequestCommentConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, dataPrivacyRequestCommentUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertDataPrivacyRequestComment updates one or more existing entries in dataPrivacyRequestComment collection
//
// This function is auto-generated
func (s *Store) UpsertDataPrivacyRequestComment(ctx context.Context, rr ...*systemType.DataPrivacyRequestComment) (err error) {
	for i := range rr {
		if err = s.checkDataPrivacyRequestCommentConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, dataPrivacyRequestCommentUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteDataPrivacyRequestComment Deletes one or more entries from dataPrivacyRequestComment collection
//
// This function is auto-generated
func (s *Store) DeleteDataPrivacyRequestComment(ctx context.Context, rr ...*systemType.DataPrivacyRequestComment) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, dataPrivacyRequestCommentDeleteQuery(s.Dialect, dataPrivacyRequestCommentPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteDataPrivacyRequestCommentByID deletes single entry from dataPrivacyRequestComment collection
//
// This function is auto-generated
func (s *Store) DeleteDataPrivacyRequestCommentByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, dataPrivacyRequestCommentDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateDataPrivacyRequestComments Deletes all rows from the dataPrivacyRequestComment collection
func (s Store) TruncateDataPrivacyRequestComments(ctx context.Context) error {
	return s.Exec(ctx, dataPrivacyRequestCommentTruncateQuery(s.Dialect))
}

// SearchDataPrivacyRequestComments returns (filtered) set of DataPrivacyRequestComments
//
// This function is auto-generated
func (s *Store) SearchDataPrivacyRequestComments(ctx context.Context, f systemType.DataPrivacyRequestCommentFilter) (set systemType.DataPrivacyRequestCommentSet, _ systemType.DataPrivacyRequestCommentFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfDataPrivacyRequestComments fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfDataPrivacyRequestComments(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfDataPrivacyRequestComments collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfDataPrivacyRequestComments(
	ctx context.Context,
	filter systemType.DataPrivacyRequestCommentFilter,
	sort filter.SortExprSet,
) (set []*systemType.DataPrivacyRequestComment, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.DataPrivacyRequestComment

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.DataPrivacyRequestCommentFilter
	)

	set = make([]*systemType.DataPrivacyRequestComment, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryDataPrivacyRequestComments(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectDataPrivacyRequestCommentCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectDataPrivacyRequestCommentCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectDataPrivacyRequestCommentCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryDataPrivacyRequestComments queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryDataPrivacyRequestComments(
	ctx context.Context,
	f systemType.DataPrivacyRequestCommentFilter,
) (_ []*systemType.DataPrivacyRequestComment, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.DataPrivacyRequestComment, 0, DefaultSliceCapacity)
		res         *systemType.DataPrivacyRequestComment
		aux         *auxDataPrivacyRequestComment
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.DataPrivacyRequestComment != nil {
		// extended filter set
		tExpr, f, err = s.Filters.DataPrivacyRequestComment(s, f)
	} else {
		// using generated filter
		tExpr, f, err = DataPrivacyRequestCommentFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for DataPrivacyRequestComment: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableDataPrivacyRequestCommentFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := dataPrivacyRequestCommentSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableDataPrivacyRequestCommentFields()); err != nil {
		err = fmt.Errorf("could generate order expression for DataPrivacyRequestComment: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query DataPrivacyRequestComment: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query DataPrivacyRequestComment: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query DataPrivacyRequestComment: %w", err)
			return
		}

		aux = new(auxDataPrivacyRequestComment)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for DataPrivacyRequestComment: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode DataPrivacyRequestComment: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// sortableDataPrivacyRequestCommentFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableDataPrivacyRequestCommentFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"id":         "id",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
	}
}

// collectDataPrivacyRequestCommentCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectDataPrivacyRequestCommentCursorValues(res *systemType.DataPrivacyRequestComment, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkDataPrivacyRequestCommentConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkDataPrivacyRequestCommentConstraints(ctx context.Context, res *systemType.DataPrivacyRequestComment) (err error) {
	return nil
}

// CreateFederationExposedModule creates one or more rows in federationExposedModule collection
//
// This function is auto-generated
func (s *Store) CreateFederationExposedModule(ctx context.Context, rr ...*federationType.ExposedModule) (err error) {
	for i := range rr {
		if err = s.checkFederationExposedModuleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationExposedModuleInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateFederationExposedModule updates one or more existing entries in federationExposedModule collection
//
// This function is auto-generated
func (s *Store) UpdateFederationExposedModule(ctx context.Context, rr ...*federationType.ExposedModule) (err error) {
	for i := range rr {
		if err = s.checkFederationExposedModuleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationExposedModuleUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertFederationExposedModule updates one or more existing entries in federationExposedModule collection
//
// This function is auto-generated
func (s *Store) UpsertFederationExposedModule(ctx context.Context, rr ...*federationType.ExposedModule) (err error) {
	for i := range rr {
		if err = s.checkFederationExposedModuleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationExposedModuleUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteFederationExposedModule Deletes one or more entries from federationExposedModule collection
//
// This function is auto-generated
func (s *Store) DeleteFederationExposedModule(ctx context.Context, rr ...*federationType.ExposedModule) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, federationExposedModuleDeleteQuery(s.Dialect, federationExposedModulePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteFederationExposedModuleByID deletes single entry from federationExposedModule collection
//
// This function is auto-generated
func (s *Store) DeleteFederationExposedModuleByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, federationExposedModuleDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateFederationExposedModules Deletes all rows from the federationExposedModule collection
func (s Store) TruncateFederationExposedModules(ctx context.Context) error {
	return s.Exec(ctx, federationExposedModuleTruncateQuery(s.Dialect))
}

// SearchFederationExposedModules returns (filtered) set of FederationExposedModules
//
// This function is auto-generated
func (s *Store) SearchFederationExposedModules(ctx context.Context, f federationType.ExposedModuleFilter) (set federationType.ExposedModuleSet, _ federationType.ExposedModuleFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfFederationExposedModules fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfFederationExposedModules(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
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
//
// This function is auto-generated
func (s *Store) fetchFullPageOfFederationExposedModules(
	ctx context.Context,
	filter federationType.ExposedModuleFilter,
	sort filter.SortExprSet,
) (set []*federationType.ExposedModule, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*federationType.ExposedModule

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter federationType.ExposedModuleFilter
	)

	set = make([]*federationType.ExposedModule, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryFederationExposedModules(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectFederationExposedModuleCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectFederationExposedModuleCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectFederationExposedModuleCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryFederationExposedModules queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryFederationExposedModules(
	ctx context.Context,
	f federationType.ExposedModuleFilter,
) (_ []*federationType.ExposedModule, more bool, err error) {
	var (
		ok bool

		set         = make([]*federationType.ExposedModule, 0, DefaultSliceCapacity)
		res         *federationType.ExposedModule
		aux         *auxFederationExposedModule
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.FederationExposedModule != nil {
		// extended filter set
		tExpr, f, err = s.Filters.FederationExposedModule(s, f)
	} else {
		// using generated filter
		tExpr, f, err = FederationExposedModuleFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for FederationExposedModule: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableFederationExposedModuleFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := federationExposedModuleSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableFederationExposedModuleFields()); err != nil {
		err = fmt.Errorf("could generate order expression for FederationExposedModule: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query FederationExposedModule: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query FederationExposedModule: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query FederationExposedModule: %w", err)
			return
		}

		aux = new(auxFederationExposedModule)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for FederationExposedModule: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode FederationExposedModule: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupFederationExposedModuleByID searches for federation module by ID
//
// It returns federation module
//
// This function is auto-generated
func (s *Store) LookupFederationExposedModuleByID(ctx context.Context, id uint64) (_ *federationType.ExposedModule, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxFederationExposedModule)
		lookup = federationExposedModuleSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableFederationExposedModuleFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableFederationExposedModuleFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"handle":     "handle",
		"id":         "id",
		"name":       "name",
		"node_id":    "node_id",
		"nodeid":     "node_id",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
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
//
// This function is auto-generated
func (s *Store) collectFederationExposedModuleCursorValues(res *federationType.ExposedModule, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
				case "nodeID":
					cur.Set(c.Column, res.NodeID, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkFederationExposedModuleConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkFederationExposedModuleConstraints(ctx context.Context, res *federationType.ExposedModule) (err error) {
	return nil
}

// CreateFederationModuleMapping creates one or more rows in federationModuleMapping collection
//
// This function is auto-generated
func (s *Store) CreateFederationModuleMapping(ctx context.Context, rr ...*federationType.ModuleMapping) (err error) {
	for i := range rr {
		if err = s.checkFederationModuleMappingConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationModuleMappingInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateFederationModuleMapping updates one or more existing entries in federationModuleMapping collection
//
// This function is auto-generated
func (s *Store) UpdateFederationModuleMapping(ctx context.Context, rr ...*federationType.ModuleMapping) (err error) {
	for i := range rr {
		if err = s.checkFederationModuleMappingConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationModuleMappingUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertFederationModuleMapping updates one or more existing entries in federationModuleMapping collection
//
// This function is auto-generated
func (s *Store) UpsertFederationModuleMapping(ctx context.Context, rr ...*federationType.ModuleMapping) (err error) {
	for i := range rr {
		if err = s.checkFederationModuleMappingConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationModuleMappingUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteFederationModuleMapping Deletes one or more entries from federationModuleMapping collection
//
// This function is auto-generated
func (s *Store) DeleteFederationModuleMapping(ctx context.Context, rr ...*federationType.ModuleMapping) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, federationModuleMappingDeleteQuery(s.Dialect, federationModuleMappingPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteFederationModuleMappingByID deletes single entry from federationModuleMapping collection
//
// This function is auto-generated
func (s *Store) DeleteFederationModuleMappingByNodeID(ctx context.Context, nodeID uint64) error {
	return s.Exec(ctx, federationModuleMappingDeleteQuery(s.Dialect, goqu.Ex{
		"node_id": nodeID,
	}))
}

// TruncateFederationModuleMappings Deletes all rows from the federationModuleMapping collection
func (s Store) TruncateFederationModuleMappings(ctx context.Context) error {
	return s.Exec(ctx, federationModuleMappingTruncateQuery(s.Dialect))
}

// SearchFederationModuleMappings returns (filtered) set of FederationModuleMappings
//
// This function is auto-generated
func (s *Store) SearchFederationModuleMappings(ctx context.Context, f federationType.ModuleMappingFilter) (set federationType.ModuleMappingSet, _ federationType.ModuleMappingFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfFederationModuleMappings fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfFederationModuleMappings(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
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
//
// This function is auto-generated
func (s *Store) fetchFullPageOfFederationModuleMappings(
	ctx context.Context,
	filter federationType.ModuleMappingFilter,
	sort filter.SortExprSet,
) (set []*federationType.ModuleMapping, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*federationType.ModuleMapping

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter federationType.ModuleMappingFilter
	)

	set = make([]*federationType.ModuleMapping, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryFederationModuleMappings(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectFederationModuleMappingCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectFederationModuleMappingCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectFederationModuleMappingCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryFederationModuleMappings queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryFederationModuleMappings(
	ctx context.Context,
	f federationType.ModuleMappingFilter,
) (_ []*federationType.ModuleMapping, more bool, err error) {
	var (
		ok bool

		set         = make([]*federationType.ModuleMapping, 0, DefaultSliceCapacity)
		res         *federationType.ModuleMapping
		aux         *auxFederationModuleMapping
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.FederationModuleMapping != nil {
		// extended filter set
		tExpr, f, err = s.Filters.FederationModuleMapping(s, f)
	} else {
		// using generated filter
		tExpr, f, err = FederationModuleMappingFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for FederationModuleMapping: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableFederationModuleMappingFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := federationModuleMappingSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableFederationModuleMappingFields()); err != nil {
		err = fmt.Errorf("could generate order expression for FederationModuleMapping: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query FederationModuleMapping: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query FederationModuleMapping: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query FederationModuleMapping: %w", err)
			return
		}

		aux = new(auxFederationModuleMapping)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for FederationModuleMapping: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode FederationModuleMapping: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupFederationModuleMappingByFederationModuleIDComposeModuleIDComposeNamespaceID searches for module mapping by federation module id and compose module id
//
// It returns module mapping
//
// This function is auto-generated
func (s *Store) LookupFederationModuleMappingByFederationModuleIDComposeModuleIDComposeNamespaceID(ctx context.Context, federationModuleID uint64, composeModuleID uint64, composeNamespaceID uint64) (_ *federationType.ModuleMapping, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxFederationModuleMapping)
		lookup = federationModuleMappingSelectQuery(s.Dialect).Where(
			goqu.I("federation_module_id").Eq(federationModuleID),
			goqu.I("compose_module_id").Eq(composeModuleID),
			goqu.I("compose_namespace_id").Eq(composeNamespaceID),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupFederationModuleMappingByFederationModuleID searches for module mapping by federation module id
//
// It returns module mapping
//
// This function is auto-generated
func (s *Store) LookupFederationModuleMappingByFederationModuleID(ctx context.Context, federationModuleID uint64) (_ *federationType.ModuleMapping, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxFederationModuleMapping)
		lookup = federationModuleMappingSelectQuery(s.Dialect).Where(
			goqu.I("federation_module_id").Eq(federationModuleID),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableFederationModuleMappingFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableFederationModuleMappingFields() map[string]string {
	return map[string]string{
		"compose_module_id":    "compose_module_id",
		"compose_namespace_id": "compose_namespace_id",
		"composemoduleid":      "compose_module_id",
		"composenamespaceid":   "compose_namespace_id",
		"federation_module_id": "federation_module_id",
		"federationmoduleid":   "federation_module_id",
		"node_id":              "node_id",
		"nodeid":               "node_id",
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
//
// This function is auto-generated
func (s *Store) collectFederationModuleMappingCursorValues(res *federationType.ModuleMapping, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkNodeID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "nodeID":
					cur.Set(c.Column, res.NodeID, c.Descending)
					pkNodeID = true
				case "federationModuleID":
					cur.Set(c.Column, res.FederationModuleID, c.Descending)
				case "composeModuleID":
					cur.Set(c.Column, res.ComposeModuleID, c.Descending)
				case "composeNamespaceID":
					cur.Set(c.Column, res.ComposeNamespaceID, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkNodeID {
		collect(&filter.SortExpr{Column: "nodeID", Descending: false})
	}

	return cur

}

// checkFederationModuleMappingConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkFederationModuleMappingConstraints(ctx context.Context, res *federationType.ModuleMapping) (err error) {
	return nil
}

// CreateFederationNode creates one or more rows in federationNode collection
//
// This function is auto-generated
func (s *Store) CreateFederationNode(ctx context.Context, rr ...*federationType.Node) (err error) {
	for i := range rr {
		if err = s.checkFederationNodeConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationNodeInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateFederationNode updates one or more existing entries in federationNode collection
//
// This function is auto-generated
func (s *Store) UpdateFederationNode(ctx context.Context, rr ...*federationType.Node) (err error) {
	for i := range rr {
		if err = s.checkFederationNodeConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationNodeUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertFederationNode updates one or more existing entries in federationNode collection
//
// This function is auto-generated
func (s *Store) UpsertFederationNode(ctx context.Context, rr ...*federationType.Node) (err error) {
	for i := range rr {
		if err = s.checkFederationNodeConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationNodeUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteFederationNode Deletes one or more entries from federationNode collection
//
// This function is auto-generated
func (s *Store) DeleteFederationNode(ctx context.Context, rr ...*federationType.Node) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, federationNodeDeleteQuery(s.Dialect, federationNodePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteFederationNodeByID deletes single entry from federationNode collection
//
// This function is auto-generated
func (s *Store) DeleteFederationNodeByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, federationNodeDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateFederationNodes Deletes all rows from the federationNode collection
func (s Store) TruncateFederationNodes(ctx context.Context) error {
	return s.Exec(ctx, federationNodeTruncateQuery(s.Dialect))
}

// SearchFederationNodes returns (filtered) set of FederationNodes
//
// This function is auto-generated
func (s *Store) SearchFederationNodes(ctx context.Context, f federationType.NodeFilter) (set federationType.NodeSet, _ federationType.NodeFilter, err error) {

	set, _, err = s.QueryFederationNodes(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryFederationNodes queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryFederationNodes(
	ctx context.Context,
	f federationType.NodeFilter,
) (_ []*federationType.Node, more bool, err error) {
	var (
		ok bool

		set         = make([]*federationType.Node, 0, DefaultSliceCapacity)
		res         *federationType.Node
		aux         *auxFederationNode
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.FederationNode != nil {
		// extended filter set
		tExpr, f, err = s.Filters.FederationNode(s, f)
	} else {
		// using generated filter
		tExpr, f, err = FederationNodeFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for FederationNode: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := federationNodeSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query FederationNode: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query FederationNode: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query FederationNode: %w", err)
			return
		}

		aux = new(auxFederationNode)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for FederationNode: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode FederationNode: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, false, err

}

// LookupFederationNodeByID searches for federation node by ID
//
// It returns federation node
//
// This function is auto-generated
func (s *Store) LookupFederationNodeByID(ctx context.Context, id uint64) (_ *federationType.Node, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxFederationNode)
		lookup = federationNodeSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupFederationNodeByBaseURLSharedNodeID searches for node by shared-node-id and base-url
//
// This function is auto-generated
func (s *Store) LookupFederationNodeByBaseURLSharedNodeID(ctx context.Context, baseURL string, sharedNodeID uint64) (_ *federationType.Node, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxFederationNode)
		lookup = federationNodeSelectQuery(s.Dialect).Where(
			goqu.I("base_url").Eq(baseURL),
			goqu.I("shared_node_id").Eq(sharedNodeID),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupFederationNodeBySharedNodeID searches for node by shared-node-id
//
// This function is auto-generated
func (s *Store) LookupFederationNodeBySharedNodeID(ctx context.Context, sharedNodeID uint64) (_ *federationType.Node, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxFederationNode)
		lookup = federationNodeSelectQuery(s.Dialect).Where(
			goqu.I("shared_node_id").Eq(sharedNodeID),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableFederationNodeFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableFederationNodeFields() map[string]string {
	return map[string]string{
		"base_url":       "base_url",
		"baseurl":        "base_url",
		"contact":        "contact",
		"created_at":     "created_at",
		"createdat":      "created_at",
		"deleted_at":     "deleted_at",
		"deletedat":      "deleted_at",
		"id":             "id",
		"name":           "name",
		"shared_node_id": "shared_node_id",
		"sharednodeid":   "shared_node_id",
		"status":         "status",
		"updated_at":     "updated_at",
		"updatedat":      "updated_at",
	}
}

// collectFederationNodeCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectFederationNodeCursorValues(res *federationType.Node, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "sharedNodeID":
					cur.Set(c.Column, res.SharedNodeID, c.Descending)
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
				case "baseURL":
					cur.Set(c.Column, res.BaseURL, c.Descending)
				case "status":
					cur.Set(c.Column, res.Status, c.Descending)
				case "contact":
					cur.Set(c.Column, res.Contact, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkFederationNodeConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkFederationNodeConstraints(ctx context.Context, res *federationType.Node) (err error) {
	return nil
}

// CreateFederationNodeSync creates one or more rows in federationNodeSync collection
//
// This function is auto-generated
func (s *Store) CreateFederationNodeSync(ctx context.Context, rr ...*federationType.NodeSync) (err error) {
	for i := range rr {
		if err = s.checkFederationNodeSyncConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationNodeSyncInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateFederationNodeSync updates one or more existing entries in federationNodeSync collection
//
// This function is auto-generated
func (s *Store) UpdateFederationNodeSync(ctx context.Context, rr ...*federationType.NodeSync) (err error) {
	for i := range rr {
		if err = s.checkFederationNodeSyncConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationNodeSyncUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertFederationNodeSync updates one or more existing entries in federationNodeSync collection
//
// This function is auto-generated
func (s *Store) UpsertFederationNodeSync(ctx context.Context, rr ...*federationType.NodeSync) (err error) {
	for i := range rr {
		if err = s.checkFederationNodeSyncConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationNodeSyncUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteFederationNodeSync Deletes one or more entries from federationNodeSync collection
//
// This function is auto-generated
func (s *Store) DeleteFederationNodeSync(ctx context.Context, rr ...*federationType.NodeSync) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, federationNodeSyncDeleteQuery(s.Dialect, federationNodeSyncPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteFederationNodeSyncByID deletes single entry from federationNodeSync collection
//
// This function is auto-generated
func (s *Store) DeleteFederationNodeSyncByNodeID(ctx context.Context, nodeID uint64) error {
	return s.Exec(ctx, federationNodeSyncDeleteQuery(s.Dialect, goqu.Ex{
		"node_id": nodeID,
	}))
}

// TruncateFederationNodeSyncs Deletes all rows from the federationNodeSync collection
func (s Store) TruncateFederationNodeSyncs(ctx context.Context) error {
	return s.Exec(ctx, federationNodeSyncTruncateQuery(s.Dialect))
}

// SearchFederationNodeSyncs returns (filtered) set of FederationNodeSyncs
//
// This function is auto-generated
func (s *Store) SearchFederationNodeSyncs(ctx context.Context, f federationType.NodeSyncFilter) (set federationType.NodeSyncSet, _ federationType.NodeSyncFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfFederationNodeSyncs fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfFederationNodeSyncs(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfFederationNodeSyncs collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfFederationNodeSyncs(
	ctx context.Context,
	filter federationType.NodeSyncFilter,
	sort filter.SortExprSet,
) (set []*federationType.NodeSync, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*federationType.NodeSync

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter federationType.NodeSyncFilter
	)

	set = make([]*federationType.NodeSync, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryFederationNodeSyncs(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectFederationNodeSyncCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectFederationNodeSyncCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectFederationNodeSyncCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryFederationNodeSyncs queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryFederationNodeSyncs(
	ctx context.Context,
	f federationType.NodeSyncFilter,
) (_ []*federationType.NodeSync, more bool, err error) {
	var (
		ok bool

		set         = make([]*federationType.NodeSync, 0, DefaultSliceCapacity)
		res         *federationType.NodeSync
		aux         *auxFederationNodeSync
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.FederationNodeSync != nil {
		// extended filter set
		tExpr, f, err = s.Filters.FederationNodeSync(s, f)
	} else {
		// using generated filter
		tExpr, f, err = FederationNodeSyncFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for FederationNodeSync: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableFederationNodeSyncFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := federationNodeSyncSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableFederationNodeSyncFields()); err != nil {
		err = fmt.Errorf("could generate order expression for FederationNodeSync: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query FederationNodeSync: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query FederationNodeSync: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query FederationNodeSync: %w", err)
			return
		}

		aux = new(auxFederationNodeSync)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for FederationNodeSync: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode FederationNodeSync: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupFederationNodeSyncByNodeID searches for sync activity by node ID
//
// It returns sync activity
//
// This function is auto-generated
func (s *Store) LookupFederationNodeSyncByNodeID(ctx context.Context, nodeID uint64) (_ *federationType.NodeSync, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxFederationNodeSync)
		lookup = federationNodeSyncSelectQuery(s.Dialect).Where(
			goqu.I("node_id").Eq(nodeID),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupFederationNodeSyncByNodeIDModuleIDSyncTypeSyncStatus searches for activity by node, type and status
//
// It returns sync activity
//
// This function is auto-generated
func (s *Store) LookupFederationNodeSyncByNodeIDModuleIDSyncTypeSyncStatus(ctx context.Context, nodeID uint64, moduleID uint64, syncType string, syncStatus string) (_ *federationType.NodeSync, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxFederationNodeSync)
		lookup = federationNodeSyncSelectQuery(s.Dialect).Where(
			goqu.I("node_id").Eq(nodeID),
			goqu.I("module_id").Eq(moduleID),
			goqu.I("sync_type").Eq(syncType),
			goqu.I("sync_status").Eq(syncStatus),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableFederationNodeSyncFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableFederationNodeSyncFields() map[string]string {
	return map[string]string{
		"module_id":      "module_id",
		"moduleid":       "module_id",
		"node_id":        "node_id",
		"nodeid":         "node_id",
		"sync_status":    "sync_status",
		"sync_type":      "sync_type",
		"syncstatus":     "sync_status",
		"synctype":       "sync_type",
		"time_of_action": "time_of_action",
		"timeofaction":   "time_of_action",
	}
}

// collectFederationNodeSyncCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectFederationNodeSyncCursorValues(res *federationType.NodeSync, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkNodeID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "nodeID":
					cur.Set(c.Column, res.NodeID, c.Descending)
					pkNodeID = true
				case "moduleID":
					cur.Set(c.Column, res.ModuleID, c.Descending)
				case "syncType":
					cur.Set(c.Column, res.SyncType, c.Descending)
				case "syncStatus":
					cur.Set(c.Column, res.SyncStatus, c.Descending)
				case "timeOfAction":
					cur.Set(c.Column, res.TimeOfAction, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkNodeID {
		collect(&filter.SortExpr{Column: "nodeID", Descending: false})
	}

	return cur

}

// checkFederationNodeSyncConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkFederationNodeSyncConstraints(ctx context.Context, res *federationType.NodeSync) (err error) {
	return nil
}

// CreateFederationSharedModule creates one or more rows in federationSharedModule collection
//
// This function is auto-generated
func (s *Store) CreateFederationSharedModule(ctx context.Context, rr ...*federationType.SharedModule) (err error) {
	for i := range rr {
		if err = s.checkFederationSharedModuleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationSharedModuleInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateFederationSharedModule updates one or more existing entries in federationSharedModule collection
//
// This function is auto-generated
func (s *Store) UpdateFederationSharedModule(ctx context.Context, rr ...*federationType.SharedModule) (err error) {
	for i := range rr {
		if err = s.checkFederationSharedModuleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationSharedModuleUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertFederationSharedModule updates one or more existing entries in federationSharedModule collection
//
// This function is auto-generated
func (s *Store) UpsertFederationSharedModule(ctx context.Context, rr ...*federationType.SharedModule) (err error) {
	for i := range rr {
		if err = s.checkFederationSharedModuleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, federationSharedModuleUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteFederationSharedModule Deletes one or more entries from federationSharedModule collection
//
// This function is auto-generated
func (s *Store) DeleteFederationSharedModule(ctx context.Context, rr ...*federationType.SharedModule) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, federationSharedModuleDeleteQuery(s.Dialect, federationSharedModulePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteFederationSharedModuleByID deletes single entry from federationSharedModule collection
//
// This function is auto-generated
func (s *Store) DeleteFederationSharedModuleByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, federationSharedModuleDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateFederationSharedModules Deletes all rows from the federationSharedModule collection
func (s Store) TruncateFederationSharedModules(ctx context.Context) error {
	return s.Exec(ctx, federationSharedModuleTruncateQuery(s.Dialect))
}

// SearchFederationSharedModules returns (filtered) set of FederationSharedModules
//
// This function is auto-generated
func (s *Store) SearchFederationSharedModules(ctx context.Context, f federationType.SharedModuleFilter) (set federationType.SharedModuleSet, _ federationType.SharedModuleFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfFederationSharedModules fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfFederationSharedModules(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfFederationSharedModules collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfFederationSharedModules(
	ctx context.Context,
	filter federationType.SharedModuleFilter,
	sort filter.SortExprSet,
) (set []*federationType.SharedModule, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*federationType.SharedModule

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter federationType.SharedModuleFilter
	)

	set = make([]*federationType.SharedModule, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryFederationSharedModules(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectFederationSharedModuleCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectFederationSharedModuleCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectFederationSharedModuleCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryFederationSharedModules queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryFederationSharedModules(
	ctx context.Context,
	f federationType.SharedModuleFilter,
) (_ []*federationType.SharedModule, more bool, err error) {
	var (
		ok bool

		set         = make([]*federationType.SharedModule, 0, DefaultSliceCapacity)
		res         *federationType.SharedModule
		aux         *auxFederationSharedModule
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.FederationSharedModule != nil {
		// extended filter set
		tExpr, f, err = s.Filters.FederationSharedModule(s, f)
	} else {
		// using generated filter
		tExpr, f, err = FederationSharedModuleFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for FederationSharedModule: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableFederationSharedModuleFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := federationSharedModuleSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableFederationSharedModuleFields()); err != nil {
		err = fmt.Errorf("could generate order expression for FederationSharedModule: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query FederationSharedModule: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query FederationSharedModule: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query FederationSharedModule: %w", err)
			return
		}

		aux = new(auxFederationSharedModule)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for FederationSharedModule: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode FederationSharedModule: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupFederationSharedModuleByID searches for shared federation module by ID
//
// It returns shared federation module
//
// This function is auto-generated
func (s *Store) LookupFederationSharedModuleByID(ctx context.Context, id uint64) (_ *federationType.SharedModule, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxFederationSharedModule)
		lookup = federationSharedModuleSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableFederationSharedModuleFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableFederationSharedModuleFields() map[string]string {
	return map[string]string{
		"created_at":                    "created_at",
		"createdat":                     "created_at",
		"deleted_at":                    "deleted_at",
		"deletedat":                     "deleted_at",
		"external_federation_module_id": "external_federation_module_id",
		"externalfederationmoduleid":    "external_federation_module_id",
		"handle":                        "handle",
		"id":                            "id",
		"name":                          "name",
		"node_id":                       "node_id",
		"nodeid":                        "node_id",
		"updated_at":                    "updated_at",
		"updatedat":                     "updated_at",
	}
}

// collectFederationSharedModuleCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectFederationSharedModuleCursorValues(res *federationType.SharedModule, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "nodeID":
					cur.Set(c.Column, res.NodeID, c.Descending)
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
				case "externalFederationModuleID":
					cur.Set(c.Column, res.ExternalFederationModuleID, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkFederationSharedModuleConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkFederationSharedModuleConstraints(ctx context.Context, res *federationType.SharedModule) (err error) {
	return nil
}

// CreateFlag creates one or more rows in flag collection
//
// This function is auto-generated
func (s *Store) CreateFlag(ctx context.Context, rr ...*flagType.Flag) (err error) {
	for i := range rr {
		if err = s.checkFlagConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, flagInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateFlag updates one or more existing entries in flag collection
//
// This function is auto-generated
func (s *Store) UpdateFlag(ctx context.Context, rr ...*flagType.Flag) (err error) {
	for i := range rr {
		if err = s.checkFlagConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, flagUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertFlag updates one or more existing entries in flag collection
//
// This function is auto-generated
func (s *Store) UpsertFlag(ctx context.Context, rr ...*flagType.Flag) (err error) {
	for i := range rr {
		if err = s.checkFlagConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, flagUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteFlag Deletes one or more entries from flag collection
//
// This function is auto-generated
func (s *Store) DeleteFlag(ctx context.Context, rr ...*flagType.Flag) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, flagDeleteQuery(s.Dialect, flagPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteFlagByID deletes single entry from flag collection
//
// This function is auto-generated
func (s *Store) DeleteFlagByKindResourceIDOwnedByName(ctx context.Context, kind string, resourceID uint64, ownedBy uint64, name string) error {
	return s.Exec(ctx, flagDeleteQuery(s.Dialect, goqu.Ex{
		"kind":         kind,
		"rel_resource": resourceID,
		"owned_by":     ownedBy,
		"name":         name,
	}))
}

// TruncateFlags Deletes all rows from the flag collection
func (s Store) TruncateFlags(ctx context.Context) error {
	return s.Exec(ctx, flagTruncateQuery(s.Dialect))
}

// SearchFlags returns (filtered) set of Flags
//
// This function is auto-generated
func (s *Store) SearchFlags(ctx context.Context, f flagType.FlagFilter) (set flagType.FlagSet, _ flagType.FlagFilter, err error) {

	set, _, err = s.QueryFlags(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryFlags queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryFlags(
	ctx context.Context,
	f flagType.FlagFilter,
) (_ []*flagType.Flag, more bool, err error) {
	var (
		set         = make([]*flagType.Flag, 0, DefaultSliceCapacity)
		res         *flagType.Flag
		aux         *auxFlag
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.Flag != nil {
		// extended filter set
		tExpr, f, err = s.Filters.Flag(s, f)
	} else {
		// using generated filter
		tExpr, f, err = FlagFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for Flag: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := flagSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query Flag: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query Flag: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query Flag: %w", err)
			return
		}

		aux = new(auxFlag)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for Flag: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode Flag: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, false, err

}

// LookupFlagByKindResourceIDOwnedByName searches for flag by kind, resource ID, owner and name
//
// This function is auto-generated
func (s *Store) LookupFlagByKindResourceIDOwnedByName(ctx context.Context, kind string, resourceID uint64, ownedBy uint64, name string) (_ *flagType.Flag, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxFlag)
		lookup = flagSelectQuery(s.Dialect).Where(
			goqu.I("kind").Eq(kind),
			goqu.I("rel_resource").Eq(resourceID),
			goqu.I("owned_by").Eq(ownedBy),
			s.Functions.LOWER(goqu.I("name")).Eq(strings.ToLower(name)),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableFlagFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableFlagFields() map[string]string {
	return map[string]string{
		"kind":        "kind",
		"name":        "name",
		"owned_by":    "owned_by",
		"ownedby":     "owned_by",
		"resource_id": "resource_id",
		"resourceid":  "resource_id",
	}
}

// collectFlagCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectFlagCursorValues(res *flagType.Flag, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkKind       bool
		pkResourceID bool
		pkOwnedBy    bool
		pkName       bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "kind":
					cur.Set(c.Column, res.Kind, c.Descending)
					pkKind = true
				case "resourceID":
					cur.Set(c.Column, res.ResourceID, c.Descending)
					pkResourceID = true
				case "ownedBy":
					cur.Set(c.Column, res.OwnedBy, c.Descending)
					pkOwnedBy = true
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
					pkName = true
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkKind {
		collect(&filter.SortExpr{Column: "kind", Descending: false})
	}
	if !hasUnique || !pkResourceID {
		collect(&filter.SortExpr{Column: "resourceID", Descending: false})
	}
	if !hasUnique || !pkOwnedBy {
		collect(&filter.SortExpr{Column: "ownedBy", Descending: false})
	}
	if !hasUnique || !pkName {
		collect(&filter.SortExpr{Column: "name", Descending: false})
	}

	return cur

}

// checkFlagConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkFlagConstraints(ctx context.Context, res *flagType.Flag) (err error) {
	return nil
}

// CreateLabel creates one or more rows in label collection
//
// This function is auto-generated
func (s *Store) CreateLabel(ctx context.Context, rr ...*labelsType.Label) (err error) {
	for i := range rr {
		if err = s.checkLabelConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, labelInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateLabel updates one or more existing entries in label collection
//
// This function is auto-generated
func (s *Store) UpdateLabel(ctx context.Context, rr ...*labelsType.Label) (err error) {
	for i := range rr {
		if err = s.checkLabelConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, labelUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertLabel updates one or more existing entries in label collection
//
// This function is auto-generated
func (s *Store) UpsertLabel(ctx context.Context, rr ...*labelsType.Label) (err error) {
	for i := range rr {
		if err = s.checkLabelConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, labelUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteLabel Deletes one or more entries from label collection
//
// This function is auto-generated
func (s *Store) DeleteLabel(ctx context.Context, rr ...*labelsType.Label) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, labelDeleteQuery(s.Dialect, labelPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteLabelByID deletes single entry from label collection
//
// This function is auto-generated
func (s *Store) DeleteLabelByKindResourceIDName(ctx context.Context, kind string, resourceID uint64, name string) error {
	return s.Exec(ctx, labelDeleteQuery(s.Dialect, goqu.Ex{
		"kind":         kind,
		"rel_resource": resourceID,
		"name":         name,
	}))
}

// TruncateLabels Deletes all rows from the label collection
func (s Store) TruncateLabels(ctx context.Context) error {
	return s.Exec(ctx, labelTruncateQuery(s.Dialect))
}

// SearchLabels returns (filtered) set of Labels
//
// This function is auto-generated
func (s *Store) SearchLabels(ctx context.Context, f labelsType.LabelFilter) (set labelsType.LabelSet, _ labelsType.LabelFilter, err error) {

	set, _, err = s.QueryLabels(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryLabels queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryLabels(
	ctx context.Context,
	f labelsType.LabelFilter,
) (_ []*labelsType.Label, more bool, err error) {
	var (
		set         = make([]*labelsType.Label, 0, DefaultSliceCapacity)
		res         *labelsType.Label
		aux         *auxLabel
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.Label != nil {
		// extended filter set
		tExpr, f, err = s.Filters.Label(s, f)
	} else {
		// using generated filter
		tExpr, f, err = LabelFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for Label: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := labelSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query Label: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query Label: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query Label: %w", err)
			return
		}

		aux = new(auxLabel)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for Label: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode Label: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, false, err

}

// LookupLabelByKindResourceIDName searches for label by kind, resource ID and name
//
// This function is auto-generated
func (s *Store) LookupLabelByKindResourceIDName(ctx context.Context, kind string, resourceID uint64, name string) (_ *labelsType.Label, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxLabel)
		lookup = labelSelectQuery(s.Dialect).Where(
			goqu.I("kind").Eq(kind),
			goqu.I("rel_resource").Eq(resourceID),
			s.Functions.LOWER(goqu.I("name")).Eq(strings.ToLower(name)),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableLabelFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableLabelFields() map[string]string {
	return map[string]string{
		"kind":        "kind",
		"name":        "name",
		"resource_id": "resource_id",
		"resourceid":  "resource_id",
	}
}

// collectLabelCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectLabelCursorValues(res *labelsType.Label, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkKind       bool
		pkResourceID bool
		pkName       bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "kind":
					cur.Set(c.Column, res.Kind, c.Descending)
					pkKind = true
				case "resourceID":
					cur.Set(c.Column, res.ResourceID, c.Descending)
					pkResourceID = true
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
					pkName = true
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkKind {
		collect(&filter.SortExpr{Column: "kind", Descending: false})
	}
	if !hasUnique || !pkResourceID {
		collect(&filter.SortExpr{Column: "resourceID", Descending: false})
	}
	if !hasUnique || !pkName {
		collect(&filter.SortExpr{Column: "name", Descending: false})
	}

	return cur

}

// checkLabelConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkLabelConstraints(ctx context.Context, res *labelsType.Label) (err error) {
	return nil
}

// CreateQueue creates one or more rows in queue collection
//
// This function is auto-generated
func (s *Store) CreateQueue(ctx context.Context, rr ...*systemType.Queue) (err error) {
	for i := range rr {
		if err = s.checkQueueConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, queueInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateQueue updates one or more existing entries in queue collection
//
// This function is auto-generated
func (s *Store) UpdateQueue(ctx context.Context, rr ...*systemType.Queue) (err error) {
	for i := range rr {
		if err = s.checkQueueConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, queueUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertQueue updates one or more existing entries in queue collection
//
// This function is auto-generated
func (s *Store) UpsertQueue(ctx context.Context, rr ...*systemType.Queue) (err error) {
	for i := range rr {
		if err = s.checkQueueConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, queueUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteQueue Deletes one or more entries from queue collection
//
// This function is auto-generated
func (s *Store) DeleteQueue(ctx context.Context, rr ...*systemType.Queue) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, queueDeleteQuery(s.Dialect, queuePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteQueueByID deletes single entry from queue collection
//
// This function is auto-generated
func (s *Store) DeleteQueueByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, queueDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateQueues Deletes all rows from the queue collection
func (s Store) TruncateQueues(ctx context.Context) error {
	return s.Exec(ctx, queueTruncateQuery(s.Dialect))
}

// SearchQueues returns (filtered) set of Queues
//
// This function is auto-generated
func (s *Store) SearchQueues(ctx context.Context, f systemType.QueueFilter) (set systemType.QueueSet, _ systemType.QueueFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfQueues fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfQueues(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfQueues collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfQueues(
	ctx context.Context,
	filter systemType.QueueFilter,
	sort filter.SortExprSet,
) (set []*systemType.Queue, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.Queue

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.QueueFilter
	)

	set = make([]*systemType.Queue, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryQueues(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectQueueCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectQueueCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectQueueCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryQueues queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryQueues(
	ctx context.Context,
	f systemType.QueueFilter,
) (_ []*systemType.Queue, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.Queue, 0, DefaultSliceCapacity)
		res         *systemType.Queue
		aux         *auxQueue
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.Queue != nil {
		// extended filter set
		tExpr, f, err = s.Filters.Queue(s, f)
	} else {
		// using generated filter
		tExpr, f, err = QueueFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for Queue: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableQueueFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := queueSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableQueueFields()); err != nil {
		err = fmt.Errorf("could generate order expression for Queue: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query Queue: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query Queue: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query Queue: %w", err)
			return
		}

		aux = new(auxQueue)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for Queue: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode Queue: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupQueueByID searches for queue by ID
//
// This function is auto-generated
func (s *Store) LookupQueueByID(ctx context.Context, id uint64) (_ *systemType.Queue, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxQueue)
		lookup = queueSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupQueueByQueue searches for queue by queue name
//
// This function is auto-generated
func (s *Store) LookupQueueByQueue(ctx context.Context, queue string) (_ *systemType.Queue, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxQueue)
		lookup = queueSelectQuery(s.Dialect).Where(
			goqu.I("queue").Eq(queue),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableQueueFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableQueueFields() map[string]string {
	return map[string]string{
		"consumer":   "consumer",
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"id":         "id",
		"queue":      "queue",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
	}
}

// collectQueueCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectQueueCursorValues(res *systemType.Queue, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "consumer":
					cur.Set(c.Column, res.Consumer, c.Descending)
				case "queue":
					cur.Set(c.Column, res.Queue, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkQueueConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkQueueConstraints(ctx context.Context, res *systemType.Queue) (err error) {
	return nil
}

// CreateQueueMessage creates one or more rows in queueMessage collection
//
// This function is auto-generated
func (s *Store) CreateQueueMessage(ctx context.Context, rr ...*systemType.QueueMessage) (err error) {
	for i := range rr {
		if err = s.checkQueueMessageConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, queueMessageInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateQueueMessage updates one or more existing entries in queueMessage collection
//
// This function is auto-generated
func (s *Store) UpdateQueueMessage(ctx context.Context, rr ...*systemType.QueueMessage) (err error) {
	for i := range rr {
		if err = s.checkQueueMessageConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, queueMessageUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertQueueMessage updates one or more existing entries in queueMessage collection
//
// This function is auto-generated
func (s *Store) UpsertQueueMessage(ctx context.Context, rr ...*systemType.QueueMessage) (err error) {
	for i := range rr {
		if err = s.checkQueueMessageConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, queueMessageUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteQueueMessage Deletes one or more entries from queueMessage collection
//
// This function is auto-generated
func (s *Store) DeleteQueueMessage(ctx context.Context, rr ...*systemType.QueueMessage) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, queueMessageDeleteQuery(s.Dialect, queueMessagePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteQueueMessageByID deletes single entry from queueMessage collection
//
// This function is auto-generated
func (s *Store) DeleteQueueMessageByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, queueMessageDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateQueueMessages Deletes all rows from the queueMessage collection
func (s Store) TruncateQueueMessages(ctx context.Context) error {
	return s.Exec(ctx, queueMessageTruncateQuery(s.Dialect))
}

// SearchQueueMessages returns (filtered) set of QueueMessages
//
// This function is auto-generated
func (s *Store) SearchQueueMessages(ctx context.Context, f systemType.QueueMessageFilter) (set systemType.QueueMessageSet, _ systemType.QueueMessageFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfQueueMessages fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfQueueMessages(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfQueueMessages collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfQueueMessages(
	ctx context.Context,
	filter systemType.QueueMessageFilter,
	sort filter.SortExprSet,
) (set []*systemType.QueueMessage, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.QueueMessage

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.QueueMessageFilter
	)

	set = make([]*systemType.QueueMessage, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryQueueMessages(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectQueueMessageCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectQueueMessageCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectQueueMessageCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryQueueMessages queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryQueueMessages(
	ctx context.Context,
	f systemType.QueueMessageFilter,
) (_ []*systemType.QueueMessage, more bool, err error) {
	var (
		set         = make([]*systemType.QueueMessage, 0, DefaultSliceCapacity)
		res         *systemType.QueueMessage
		aux         *auxQueueMessage
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.QueueMessage != nil {
		// extended filter set
		tExpr, f, err = s.Filters.QueueMessage(s, f)
	} else {
		// using generated filter
		tExpr, f, err = QueueMessageFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for QueueMessage: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableQueueMessageFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := queueMessageSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableQueueMessageFields()); err != nil {
		err = fmt.Errorf("could generate order expression for QueueMessage: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query QueueMessage: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query QueueMessage: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query QueueMessage: %w", err)
			return
		}

		aux = new(auxQueueMessage)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for QueueMessage: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode QueueMessage: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// sortableQueueMessageFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableQueueMessageFields() map[string]string {
	return map[string]string{
		"created":   "created",
		"id":        "id",
		"processed": "processed",
		"queue":     "queue",
	}
}

// collectQueueMessageCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectQueueMessageCursorValues(res *systemType.QueueMessage, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "queue":
					cur.Set(c.Column, res.Queue, c.Descending)
				case "created":
					cur.Set(c.Column, res.Created, c.Descending)
				case "processed":
					cur.Set(c.Column, res.Processed, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkQueueMessageConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkQueueMessageConstraints(ctx context.Context, res *systemType.QueueMessage) (err error) {
	return nil
}

// CreateRbacRule creates one or more rows in rbacRule collection
//
// This function is auto-generated
func (s *Store) CreateRbacRule(ctx context.Context, rr ...*rbacType.Rule) (err error) {
	for i := range rr {
		if err = s.checkRbacRuleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, rbacRuleInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateRbacRule updates one or more existing entries in rbacRule collection
//
// This function is auto-generated
func (s *Store) UpdateRbacRule(ctx context.Context, rr ...*rbacType.Rule) (err error) {
	for i := range rr {
		if err = s.checkRbacRuleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, rbacRuleUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertRbacRule updates one or more existing entries in rbacRule collection
//
// This function is auto-generated
func (s *Store) UpsertRbacRule(ctx context.Context, rr ...*rbacType.Rule) (err error) {
	for i := range rr {
		if err = s.checkRbacRuleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, rbacRuleUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteRbacRule Deletes one or more entries from rbacRule collection
//
// This function is auto-generated
func (s *Store) DeleteRbacRule(ctx context.Context, rr ...*rbacType.Rule) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, rbacRuleDeleteQuery(s.Dialect, rbacRulePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteRbacRuleByID deletes single entry from rbacRule collection
//
// This function is auto-generated
func (s *Store) DeleteRbacRuleByRoleIDResourceOperation(ctx context.Context, roleID uint64, resource string, operation string) error {
	return s.Exec(ctx, rbacRuleDeleteQuery(s.Dialect, goqu.Ex{
		"rel_role":  roleID,
		"resource":  resource,
		"operation": operation,
	}))
}

// TruncateRbacRules Deletes all rows from the rbacRule collection
func (s Store) TruncateRbacRules(ctx context.Context) error {
	return s.Exec(ctx, rbacRuleTruncateQuery(s.Dialect))
}

// SearchRbacRules returns (filtered) set of RbacRules
//
// This function is auto-generated
func (s *Store) SearchRbacRules(ctx context.Context, f rbacType.RuleFilter) (set rbacType.RuleSet, _ rbacType.RuleFilter, err error) {

	set, _, err = s.QueryRbacRules(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryRbacRules queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryRbacRules(
	ctx context.Context,
	f rbacType.RuleFilter,
) (_ []*rbacType.Rule, more bool, err error) {
	var (
		set         = make([]*rbacType.Rule, 0, DefaultSliceCapacity)
		res         *rbacType.Rule
		aux         *auxRbacRule
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.RbacRule != nil {
		// extended filter set
		tExpr, f, err = s.Filters.RbacRule(s, f)
	} else {
		// using generated filter
		tExpr, f, err = RbacRuleFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for RbacRule: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := rbacRuleSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query RbacRule: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query RbacRule: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query RbacRule: %w", err)
			return
		}

		aux = new(auxRbacRule)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for RbacRule: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode RbacRule: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, false, err

}

// sortableRbacRuleFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableRbacRuleFields() map[string]string {
	return map[string]string{
		"operation": "operation",
		"resource":  "resource",
		"role_id":   "role_id",
		"roleid":    "role_id",
	}
}

// collectRbacRuleCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectRbacRuleCursorValues(res *rbacType.Rule, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkRoleID    bool
		pkResource  bool
		pkOperation bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "roleID":
					cur.Set(c.Column, res.RoleID, c.Descending)
					pkRoleID = true
				case "resource":
					cur.Set(c.Column, res.Resource, c.Descending)
					pkResource = true
				case "operation":
					cur.Set(c.Column, res.Operation, c.Descending)
					pkOperation = true
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkRoleID {
		collect(&filter.SortExpr{Column: "roleID", Descending: false})
	}
	if !hasUnique || !pkResource {
		collect(&filter.SortExpr{Column: "resource", Descending: false})
	}
	if !hasUnique || !pkOperation {
		collect(&filter.SortExpr{Column: "operation", Descending: false})
	}

	return cur

}

// checkRbacRuleConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkRbacRuleConstraints(ctx context.Context, res *rbacType.Rule) (err error) {
	return nil
}

// CreateReminder creates one or more rows in reminder collection
//
// This function is auto-generated
func (s *Store) CreateReminder(ctx context.Context, rr ...*systemType.Reminder) (err error) {
	for i := range rr {
		if err = s.checkReminderConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, reminderInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateReminder updates one or more existing entries in reminder collection
//
// This function is auto-generated
func (s *Store) UpdateReminder(ctx context.Context, rr ...*systemType.Reminder) (err error) {
	for i := range rr {
		if err = s.checkReminderConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, reminderUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertReminder updates one or more existing entries in reminder collection
//
// This function is auto-generated
func (s *Store) UpsertReminder(ctx context.Context, rr ...*systemType.Reminder) (err error) {
	for i := range rr {
		if err = s.checkReminderConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, reminderUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteReminder Deletes one or more entries from reminder collection
//
// This function is auto-generated
func (s *Store) DeleteReminder(ctx context.Context, rr ...*systemType.Reminder) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, reminderDeleteQuery(s.Dialect, reminderPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteReminderByID deletes single entry from reminder collection
//
// This function is auto-generated
func (s *Store) DeleteReminderByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, reminderDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateReminders Deletes all rows from the reminder collection
func (s Store) TruncateReminders(ctx context.Context) error {
	return s.Exec(ctx, reminderTruncateQuery(s.Dialect))
}

// SearchReminders returns (filtered) set of Reminders
//
// This function is auto-generated
func (s *Store) SearchReminders(ctx context.Context, f systemType.ReminderFilter) (set systemType.ReminderSet, _ systemType.ReminderFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfReminders fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfReminders(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfReminders collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfReminders(
	ctx context.Context,
	filter systemType.ReminderFilter,
	sort filter.SortExprSet,
) (set []*systemType.Reminder, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.Reminder

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.ReminderFilter
	)

	set = make([]*systemType.Reminder, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryReminders(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectReminderCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectReminderCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectReminderCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryReminders queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryReminders(
	ctx context.Context,
	f systemType.ReminderFilter,
) (_ []*systemType.Reminder, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.Reminder, 0, DefaultSliceCapacity)
		res         *systemType.Reminder
		aux         *auxReminder
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.Reminder != nil {
		// extended filter set
		tExpr, f, err = s.Filters.Reminder(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ReminderFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for Reminder: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableReminderFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := reminderSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableReminderFields()); err != nil {
		err = fmt.Errorf("could generate order expression for Reminder: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query Reminder: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query Reminder: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query Reminder: %w", err)
			return
		}

		aux = new(auxReminder)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for Reminder: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode Reminder: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupReminderByID
//
// This function is auto-generated
func (s *Store) LookupReminderByID(ctx context.Context, id uint64) (_ *systemType.Reminder, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxReminder)
		lookup = reminderSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableReminderFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableReminderFields() map[string]string {
	return map[string]string{
		"assigned_at":  "assigned_at",
		"assignedat":   "assigned_at",
		"created_at":   "created_at",
		"createdat":    "created_at",
		"deleted_at":   "deleted_at",
		"deletedat":    "deleted_at",
		"dismissed_at": "dismissed_at",
		"dismissedat":  "dismissed_at",
		"id":           "id",
		"remind_at":    "remind_at",
		"remindat":     "remind_at",
		"resource":     "resource",
		"updated_at":   "updated_at",
		"updatedat":    "updated_at",
	}
}

// collectReminderCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectReminderCursorValues(res *systemType.Reminder, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "resource":
					cur.Set(c.Column, res.Resource, c.Descending)
				case "assignedAt":
					cur.Set(c.Column, res.AssignedAt, c.Descending)
				case "dismissedAt":
					cur.Set(c.Column, res.DismissedAt, c.Descending)
				case "remindAt":
					cur.Set(c.Column, res.RemindAt, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkReminderConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkReminderConstraints(ctx context.Context, res *systemType.Reminder) (err error) {
	return nil
}

// CreateReport creates one or more rows in report collection
//
// This function is auto-generated
func (s *Store) CreateReport(ctx context.Context, rr ...*systemType.Report) (err error) {
	for i := range rr {
		if err = s.checkReportConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, reportInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateReport updates one or more existing entries in report collection
//
// This function is auto-generated
func (s *Store) UpdateReport(ctx context.Context, rr ...*systemType.Report) (err error) {
	for i := range rr {
		if err = s.checkReportConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, reportUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertReport updates one or more existing entries in report collection
//
// This function is auto-generated
func (s *Store) UpsertReport(ctx context.Context, rr ...*systemType.Report) (err error) {
	for i := range rr {
		if err = s.checkReportConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, reportUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteReport Deletes one or more entries from report collection
//
// This function is auto-generated
func (s *Store) DeleteReport(ctx context.Context, rr ...*systemType.Report) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, reportDeleteQuery(s.Dialect, reportPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteReportByID deletes single entry from report collection
//
// This function is auto-generated
func (s *Store) DeleteReportByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, reportDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateReports Deletes all rows from the report collection
func (s Store) TruncateReports(ctx context.Context) error {
	return s.Exec(ctx, reportTruncateQuery(s.Dialect))
}

// SearchReports returns (filtered) set of Reports
//
// This function is auto-generated
func (s *Store) SearchReports(ctx context.Context, f systemType.ReportFilter) (set systemType.ReportSet, _ systemType.ReportFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfReports fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfReports(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfReports collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfReports(
	ctx context.Context,
	filter systemType.ReportFilter,
	sort filter.SortExprSet,
) (set []*systemType.Report, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.Report

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.ReportFilter
	)

	set = make([]*systemType.Report, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryReports(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectReportCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectReportCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectReportCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryReports queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryReports(
	ctx context.Context,
	f systemType.ReportFilter,
) (_ []*systemType.Report, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.Report, 0, DefaultSliceCapacity)
		res         *systemType.Report
		aux         *auxReport
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.Report != nil {
		// extended filter set
		tExpr, f, err = s.Filters.Report(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ReportFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for Report: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableReportFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := reportSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableReportFields()); err != nil {
		err = fmt.Errorf("could generate order expression for Report: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query Report: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query Report: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query Report: %w", err)
			return
		}

		aux = new(auxReport)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for Report: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode Report: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupReportByID searches for report by ID
//
// It returns report even if deleted
//
// This function is auto-generated
func (s *Store) LookupReportByID(ctx context.Context, id uint64) (_ *systemType.Report, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxReport)
		lookup = reportSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupReportByHandle searches for report by handle
//
// It returns report if deleted
//
// This function is auto-generated
func (s *Store) LookupReportByHandle(ctx context.Context, handle string) (_ *systemType.Report, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxReport)
		lookup = reportSelectQuery(s.Dialect).Where(
			s.Functions.LOWER(goqu.I("handle")).Eq(strings.ToLower(handle)),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableReportFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableReportFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"handle":     "handle",
		"id":         "id",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
	}
}

// collectReportCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectReportCursorValues(res *systemType.Report, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkReportConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkReportConstraints(ctx context.Context, res *systemType.Report) (err error) {
	err = func() (err error) {

		// handling string type as default
		if len(res.Handle) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupReportByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	return nil
}

// CreateResourceActivity creates one or more rows in resourceActivity collection
//
// This function is auto-generated
func (s *Store) CreateResourceActivity(ctx context.Context, rr ...*discoveryType.ResourceActivity) (err error) {
	for i := range rr {
		if err = s.checkResourceActivityConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, resourceActivityInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateResourceActivity updates one or more existing entries in resourceActivity collection
//
// This function is auto-generated
func (s *Store) UpdateResourceActivity(ctx context.Context, rr ...*discoveryType.ResourceActivity) (err error) {
	for i := range rr {
		if err = s.checkResourceActivityConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, resourceActivityUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertResourceActivity updates one or more existing entries in resourceActivity collection
//
// This function is auto-generated
func (s *Store) UpsertResourceActivity(ctx context.Context, rr ...*discoveryType.ResourceActivity) (err error) {
	for i := range rr {
		if err = s.checkResourceActivityConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, resourceActivityUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteResourceActivity Deletes one or more entries from resourceActivity collection
//
// This function is auto-generated
func (s *Store) DeleteResourceActivity(ctx context.Context, rr ...*discoveryType.ResourceActivity) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, resourceActivityDeleteQuery(s.Dialect, resourceActivityPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteResourceActivityByID deletes single entry from resourceActivity collection
//
// This function is auto-generated
func (s *Store) DeleteResourceActivityByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, resourceActivityDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateResourceActivitys Deletes all rows from the resourceActivity collection
func (s Store) TruncateResourceActivitys(ctx context.Context) error {
	return s.Exec(ctx, resourceActivityTruncateQuery(s.Dialect))
}

// SearchResourceActivitys returns (filtered) set of ResourceActivitys
//
// This function is auto-generated
func (s *Store) SearchResourceActivitys(ctx context.Context, f discoveryType.ResourceActivityFilter) (set discoveryType.ResourceActivitySet, _ discoveryType.ResourceActivityFilter, err error) {

	set, _, err = s.QueryResourceActivitys(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryResourceActivitys queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryResourceActivitys(
	ctx context.Context,
	f discoveryType.ResourceActivityFilter,
) (_ []*discoveryType.ResourceActivity, more bool, err error) {
	var (
		set         = make([]*discoveryType.ResourceActivity, 0, DefaultSliceCapacity)
		res         *discoveryType.ResourceActivity
		aux         *auxResourceActivity
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.ResourceActivity != nil {
		// extended filter set
		tExpr, f, err = s.Filters.ResourceActivity(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ResourceActivityFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for ResourceActivity: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := resourceActivitySelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query ResourceActivity: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query ResourceActivity: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query ResourceActivity: %w", err)
			return
		}

		aux = new(auxResourceActivity)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for ResourceActivity: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode ResourceActivity: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, false, err

}

// sortableResourceActivityFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableResourceActivityFields() map[string]string {
	return map[string]string{
		"id":        "id",
		"timestamp": "timestamp",
	}
}

// collectResourceActivityCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectResourceActivityCursorValues(res *discoveryType.ResourceActivity, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "timestamp":
					cur.Set(c.Column, res.Timestamp, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkResourceActivityConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkResourceActivityConstraints(ctx context.Context, res *discoveryType.ResourceActivity) (err error) {
	return nil
}

// CreateResourceTranslation creates one or more rows in resourceTranslation collection
//
// This function is auto-generated
func (s *Store) CreateResourceTranslation(ctx context.Context, rr ...*systemType.ResourceTranslation) (err error) {
	for i := range rr {
		if err = s.checkResourceTranslationConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, resourceTranslationInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateResourceTranslation updates one or more existing entries in resourceTranslation collection
//
// This function is auto-generated
func (s *Store) UpdateResourceTranslation(ctx context.Context, rr ...*systemType.ResourceTranslation) (err error) {
	for i := range rr {
		if err = s.checkResourceTranslationConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, resourceTranslationUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertResourceTranslation updates one or more existing entries in resourceTranslation collection
//
// This function is auto-generated
func (s *Store) UpsertResourceTranslation(ctx context.Context, rr ...*systemType.ResourceTranslation) (err error) {
	for i := range rr {
		if err = s.checkResourceTranslationConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, resourceTranslationUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteResourceTranslation Deletes one or more entries from resourceTranslation collection
//
// This function is auto-generated
func (s *Store) DeleteResourceTranslation(ctx context.Context, rr ...*systemType.ResourceTranslation) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, resourceTranslationDeleteQuery(s.Dialect, resourceTranslationPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteResourceTranslationByID deletes single entry from resourceTranslation collection
//
// This function is auto-generated
func (s *Store) DeleteResourceTranslationByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, resourceTranslationDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateResourceTranslations Deletes all rows from the resourceTranslation collection
func (s Store) TruncateResourceTranslations(ctx context.Context) error {
	return s.Exec(ctx, resourceTranslationTruncateQuery(s.Dialect))
}

// SearchResourceTranslations returns (filtered) set of ResourceTranslations
//
// This function is auto-generated
func (s *Store) SearchResourceTranslations(ctx context.Context, f systemType.ResourceTranslationFilter) (set systemType.ResourceTranslationSet, _ systemType.ResourceTranslationFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfResourceTranslations fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfResourceTranslations(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
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
//
// This function is auto-generated
func (s *Store) fetchFullPageOfResourceTranslations(
	ctx context.Context,
	filter systemType.ResourceTranslationFilter,
	sort filter.SortExprSet,
) (set []*systemType.ResourceTranslation, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.ResourceTranslation

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.ResourceTranslationFilter
	)

	set = make([]*systemType.ResourceTranslation, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryResourceTranslations(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectResourceTranslationCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectResourceTranslationCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectResourceTranslationCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryResourceTranslations queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryResourceTranslations(
	ctx context.Context,
	f systemType.ResourceTranslationFilter,
) (_ []*systemType.ResourceTranslation, more bool, err error) {
	var (
		set         = make([]*systemType.ResourceTranslation, 0, DefaultSliceCapacity)
		res         *systemType.ResourceTranslation
		aux         *auxResourceTranslation
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.ResourceTranslation != nil {
		// extended filter set
		tExpr, f, err = s.Filters.ResourceTranslation(s, f)
	} else {
		// using generated filter
		tExpr, f, err = ResourceTranslationFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for ResourceTranslation: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableResourceTranslationFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := resourceTranslationSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableResourceTranslationFields()); err != nil {
		err = fmt.Errorf("could generate order expression for ResourceTranslation: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query ResourceTranslation: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query ResourceTranslation: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query ResourceTranslation: %w", err)
			return
		}

		aux = new(auxResourceTranslation)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for ResourceTranslation: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode ResourceTranslation: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupResourceTranslationByID searches for resource translation by ID
// It also returns deleted resource translations.
//
// This function is auto-generated
func (s *Store) LookupResourceTranslationByID(ctx context.Context, id uint64) (_ *systemType.ResourceTranslation, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxResourceTranslation)
		lookup = resourceTranslationSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableResourceTranslationFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableResourceTranslationFields() map[string]string {
	return map[string]string{
		"created_at": "created_at",
		"createdat":  "created_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
		"id":         "id",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
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
//
// This function is auto-generated
func (s *Store) collectResourceTranslationCursorValues(res *systemType.ResourceTranslation, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkResourceTranslationConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkResourceTranslationConstraints(ctx context.Context, res *systemType.ResourceTranslation) (err error) {
	return nil
}

// CreateRole creates one or more rows in role collection
//
// This function is auto-generated
func (s *Store) CreateRole(ctx context.Context, rr ...*systemType.Role) (err error) {
	for i := range rr {
		if err = s.checkRoleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, roleInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateRole updates one or more existing entries in role collection
//
// This function is auto-generated
func (s *Store) UpdateRole(ctx context.Context, rr ...*systemType.Role) (err error) {
	for i := range rr {
		if err = s.checkRoleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, roleUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertRole updates one or more existing entries in role collection
//
// This function is auto-generated
func (s *Store) UpsertRole(ctx context.Context, rr ...*systemType.Role) (err error) {
	for i := range rr {
		if err = s.checkRoleConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, roleUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteRole Deletes one or more entries from role collection
//
// This function is auto-generated
func (s *Store) DeleteRole(ctx context.Context, rr ...*systemType.Role) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, roleDeleteQuery(s.Dialect, rolePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteRoleByID deletes single entry from role collection
//
// This function is auto-generated
func (s *Store) DeleteRoleByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, roleDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateRoles Deletes all rows from the role collection
func (s Store) TruncateRoles(ctx context.Context) error {
	return s.Exec(ctx, roleTruncateQuery(s.Dialect))
}

// SearchRoles returns (filtered) set of Roles
//
// This function is auto-generated
func (s *Store) SearchRoles(ctx context.Context, f systemType.RoleFilter) (set systemType.RoleSet, _ systemType.RoleFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfRoles fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfRoles(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfRoles collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfRoles(
	ctx context.Context,
	filter systemType.RoleFilter,
	sort filter.SortExprSet,
) (set []*systemType.Role, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.Role

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.RoleFilter
	)

	set = make([]*systemType.Role, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryRoles(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectRoleCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectRoleCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectRoleCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryRoles queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryRoles(
	ctx context.Context,
	f systemType.RoleFilter,
) (_ []*systemType.Role, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.Role, 0, DefaultSliceCapacity)
		res         *systemType.Role
		aux         *auxRole
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.Role != nil {
		// extended filter set
		tExpr, f, err = s.Filters.Role(s, f)
	} else {
		// using generated filter
		tExpr, f, err = RoleFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for Role: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableRoleFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := roleSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableRoleFields()); err != nil {
		err = fmt.Errorf("could generate order expression for Role: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query Role: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query Role: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query Role: %w", err)
			return
		}

		aux = new(auxRole)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for Role: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode Role: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupRoleByID searches for role by ID
//
// It returns role even if deleted or suspended
//
// This function is auto-generated
func (s *Store) LookupRoleByID(ctx context.Context, id uint64) (_ *systemType.Role, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxRole)
		lookup = roleSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupRoleByHandle searches for role by handle
//
// It returns only valid role (not deleted, not suspended)
//
// This function is auto-generated
func (s *Store) LookupRoleByHandle(ctx context.Context, handle string) (_ *systemType.Role, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxRole)
		lookup = roleSelectQuery(s.Dialect).Where(
			s.Functions.LOWER(goqu.I("handle")).Eq(strings.ToLower(handle)),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupRoleByName searches for role by name
//
// It returns only valid role (not deleted, not suspended)
//
// This function is auto-generated
func (s *Store) LookupRoleByName(ctx context.Context, name string) (_ *systemType.Role, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxRole)
		lookup = roleSelectQuery(s.Dialect).Where(
			goqu.I("name").Eq(name),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableRoleFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableRoleFields() map[string]string {
	return map[string]string{
		"archived_at": "archived_at",
		"archivedat":  "archived_at",
		"created_at":  "created_at",
		"createdat":   "created_at",
		"deleted_at":  "deleted_at",
		"deletedat":   "deleted_at",
		"handle":      "handle",
		"id":          "id",
		"name":        "name",
		"updated_at":  "updated_at",
		"updatedat":   "updated_at",
	}
}

// collectRoleCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectRoleCursorValues(res *systemType.Role, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "archivedAt":
					cur.Set(c.Column, res.ArchivedAt, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkRoleConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkRoleConstraints(ctx context.Context, res *systemType.Role) (err error) {
	err = func() (err error) {

		// handling string type as default
		if len(res.Handle) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupRoleByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	err = func() (err error) {

		// handling string type as default
		if len(res.Name) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupRoleByName(ctx, res.Name)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	return nil
}

// CreateRoleMember creates one or more rows in roleMember collection
//
// This function is auto-generated
func (s *Store) CreateRoleMember(ctx context.Context, rr ...*systemType.RoleMember) (err error) {
	for i := range rr {
		if err = s.checkRoleMemberConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, roleMemberInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateRoleMember updates one or more existing entries in roleMember collection
//
// This function is auto-generated
func (s *Store) UpdateRoleMember(ctx context.Context, rr ...*systemType.RoleMember) (err error) {
	for i := range rr {
		if err = s.checkRoleMemberConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, roleMemberUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertRoleMember updates one or more existing entries in roleMember collection
//
// This function is auto-generated
func (s *Store) UpsertRoleMember(ctx context.Context, rr ...*systemType.RoleMember) (err error) {
	for i := range rr {
		if err = s.checkRoleMemberConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, roleMemberUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteRoleMember Deletes one or more entries from roleMember collection
//
// This function is auto-generated
func (s *Store) DeleteRoleMember(ctx context.Context, rr ...*systemType.RoleMember) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, roleMemberDeleteQuery(s.Dialect, roleMemberPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteRoleMemberByID deletes single entry from roleMember collection
//
// This function is auto-generated
func (s *Store) DeleteRoleMemberByUserIDRoleID(ctx context.Context, userID uint64, roleID uint64) error {
	return s.Exec(ctx, roleMemberDeleteQuery(s.Dialect, goqu.Ex{
		"rel_user": userID,
		"rel_role": roleID,
	}))
}

// TruncateRoleMembers Deletes all rows from the roleMember collection
func (s Store) TruncateRoleMembers(ctx context.Context) error {
	return s.Exec(ctx, roleMemberTruncateQuery(s.Dialect))
}

// SearchRoleMembers returns (filtered) set of RoleMembers
//
// This function is auto-generated
func (s *Store) SearchRoleMembers(ctx context.Context, f systemType.RoleMemberFilter) (set systemType.RoleMemberSet, _ systemType.RoleMemberFilter, err error) {

	set, _, err = s.QueryRoleMembers(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QueryRoleMembers queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryRoleMembers(
	ctx context.Context,
	f systemType.RoleMemberFilter,
) (_ []*systemType.RoleMember, more bool, err error) {
	var (
		set         = make([]*systemType.RoleMember, 0, DefaultSliceCapacity)
		res         *systemType.RoleMember
		aux         *auxRoleMember
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.RoleMember != nil {
		// extended filter set
		tExpr, f, err = s.Filters.RoleMember(s, f)
	} else {
		// using generated filter
		tExpr, f, err = RoleMemberFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for RoleMember: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := roleMemberSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query RoleMember: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query RoleMember: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query RoleMember: %w", err)
			return
		}

		aux = new(auxRoleMember)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for RoleMember: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode RoleMember: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, false, err

}

// sortableRoleMemberFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableRoleMemberFields() map[string]string {
	return map[string]string{
		"role_id": "role_id",
		"roleid":  "role_id",
		"user_id": "user_id",
		"userid":  "user_id",
	}
}

// collectRoleMemberCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectRoleMemberCursorValues(res *systemType.RoleMember, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkUserID bool
		pkRoleID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "userID":
					cur.Set(c.Column, res.UserID, c.Descending)
					pkUserID = true
				case "roleID":
					cur.Set(c.Column, res.RoleID, c.Descending)
					pkRoleID = true
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkUserID {
		collect(&filter.SortExpr{Column: "userID", Descending: false})
	}
	if !hasUnique || !pkRoleID {
		collect(&filter.SortExpr{Column: "roleID", Descending: false})
	}

	return cur

}

// checkRoleMemberConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkRoleMemberConstraints(ctx context.Context, res *systemType.RoleMember) (err error) {
	return nil
}

// CreateSettingValue creates one or more rows in settingValue collection
//
// This function is auto-generated
func (s *Store) CreateSettingValue(ctx context.Context, rr ...*systemType.SettingValue) (err error) {
	for i := range rr {
		if err = s.checkSettingValueConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, settingValueInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateSettingValue updates one or more existing entries in settingValue collection
//
// This function is auto-generated
func (s *Store) UpdateSettingValue(ctx context.Context, rr ...*systemType.SettingValue) (err error) {
	for i := range rr {
		if err = s.checkSettingValueConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, settingValueUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertSettingValue updates one or more existing entries in settingValue collection
//
// This function is auto-generated
func (s *Store) UpsertSettingValue(ctx context.Context, rr ...*systemType.SettingValue) (err error) {
	for i := range rr {
		if err = s.checkSettingValueConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, settingValueUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteSettingValue Deletes one or more entries from settingValue collection
//
// This function is auto-generated
func (s *Store) DeleteSettingValue(ctx context.Context, rr ...*systemType.SettingValue) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, settingValueDeleteQuery(s.Dialect, settingValuePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteSettingValueByID deletes single entry from settingValue collection
//
// This function is auto-generated
func (s *Store) DeleteSettingValueByNameOwnedBy(ctx context.Context, name string, ownedBy uint64) error {
	return s.Exec(ctx, settingValueDeleteQuery(s.Dialect, goqu.Ex{
		"name":      name,
		"rel_owner": ownedBy,
	}))
}

// TruncateSettingValues Deletes all rows from the settingValue collection
func (s Store) TruncateSettingValues(ctx context.Context) error {
	return s.Exec(ctx, settingValueTruncateQuery(s.Dialect))
}

// SearchSettingValues returns (filtered) set of SettingValues
//
// This function is auto-generated
func (s *Store) SearchSettingValues(ctx context.Context, f systemType.SettingsFilter) (set systemType.SettingValueSet, _ systemType.SettingsFilter, err error) {

	set, _, err = s.QuerySettingValues(ctx, f)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// QuerySettingValues queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QuerySettingValues(
	ctx context.Context,
	f systemType.SettingsFilter,
) (_ []*systemType.SettingValue, more bool, err error) {
	var (
		set         = make([]*systemType.SettingValue, 0, DefaultSliceCapacity)
		res         *systemType.SettingValue
		aux         *auxSettingValue
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	)

	if s.Filters.SettingValue != nil {
		// extended filter set
		tExpr, f, err = s.Filters.SettingValue(s, f)
	} else {
		// using generated filter
		tExpr, f, err = SettingValueFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for SettingValue: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	query := settingValueSelectQuery(s.Dialect).Where(expr...)

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query SettingValue: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query SettingValue: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query SettingValue: %w", err)
			return
		}

		aux = new(auxSettingValue)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for SettingValue: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode SettingValue: %w", err)
			return
		}

		set = append(set, res)
	}

	return set, false, err

}

// LookupSettingValueByNameOwnedBy searches for settings by name and owner
//
// This function is auto-generated
func (s *Store) LookupSettingValueByNameOwnedBy(ctx context.Context, name string, ownedBy uint64) (_ *systemType.SettingValue, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxSettingValue)
		lookup = settingValueSelectQuery(s.Dialect).Where(
			s.Functions.LOWER(goqu.I("name")).Eq(strings.ToLower(name)),
			goqu.I("rel_owner").Eq(ownedBy),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableSettingValueFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableSettingValueFields() map[string]string {
	return map[string]string{
		"name":       "name",
		"owned_by":   "owned_by",
		"ownedby":    "owned_by",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
	}
}

// collectSettingValueCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectSettingValueCursorValues(res *systemType.SettingValue, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkOwnedBy bool
		pkName    bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "ownedBy":
					cur.Set(c.Column, res.OwnedBy, c.Descending)
					pkOwnedBy = true
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
					pkName = true
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkOwnedBy {
		collect(&filter.SortExpr{Column: "ownedBy", Descending: false})
	}
	if !hasUnique || !pkName {
		collect(&filter.SortExpr{Column: "name", Descending: false})
	}

	return cur

}

// checkSettingValueConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkSettingValueConstraints(ctx context.Context, res *systemType.SettingValue) (err error) {
	return nil
}

// CreateTemplate creates one or more rows in template collection
//
// This function is auto-generated
func (s *Store) CreateTemplate(ctx context.Context, rr ...*systemType.Template) (err error) {
	for i := range rr {
		if err = s.checkTemplateConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, templateInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateTemplate updates one or more existing entries in template collection
//
// This function is auto-generated
func (s *Store) UpdateTemplate(ctx context.Context, rr ...*systemType.Template) (err error) {
	for i := range rr {
		if err = s.checkTemplateConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, templateUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertTemplate updates one or more existing entries in template collection
//
// This function is auto-generated
func (s *Store) UpsertTemplate(ctx context.Context, rr ...*systemType.Template) (err error) {
	for i := range rr {
		if err = s.checkTemplateConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, templateUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteTemplate Deletes one or more entries from template collection
//
// This function is auto-generated
func (s *Store) DeleteTemplate(ctx context.Context, rr ...*systemType.Template) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, templateDeleteQuery(s.Dialect, templatePrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteTemplateByID deletes single entry from template collection
//
// This function is auto-generated
func (s *Store) DeleteTemplateByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, templateDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateTemplates Deletes all rows from the template collection
func (s Store) TruncateTemplates(ctx context.Context) error {
	return s.Exec(ctx, templateTruncateQuery(s.Dialect))
}

// SearchTemplates returns (filtered) set of Templates
//
// This function is auto-generated
func (s *Store) SearchTemplates(ctx context.Context, f systemType.TemplateFilter) (set systemType.TemplateSet, _ systemType.TemplateFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfTemplates fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfTemplates(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfTemplates collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfTemplates(
	ctx context.Context,
	filter systemType.TemplateFilter,
	sort filter.SortExprSet,
) (set []*systemType.Template, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.Template

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.TemplateFilter
	)

	set = make([]*systemType.Template, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryTemplates(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectTemplateCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectTemplateCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectTemplateCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryTemplates queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryTemplates(
	ctx context.Context,
	f systemType.TemplateFilter,
) (_ []*systemType.Template, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.Template, 0, DefaultSliceCapacity)
		res         *systemType.Template
		aux         *auxTemplate
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.Template != nil {
		// extended filter set
		tExpr, f, err = s.Filters.Template(s, f)
	} else {
		// using generated filter
		tExpr, f, err = TemplateFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for Template: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableTemplateFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := templateSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableTemplateFields()); err != nil {
		err = fmt.Errorf("could generate order expression for Template: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query Template: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query Template: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query Template: %w", err)
			return
		}

		aux = new(auxTemplate)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for Template: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode Template: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupTemplateByID searches for template by ID
//
// It also returns deleted templates.
//
// This function is auto-generated
func (s *Store) LookupTemplateByID(ctx context.Context, id uint64) (_ *systemType.Template, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxTemplate)
		lookup = templateSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupTemplateByHandle searches for template by handle
//
// It returns only valid templates (not deleted)
//
// This function is auto-generated
func (s *Store) LookupTemplateByHandle(ctx context.Context, handle string) (_ *systemType.Template, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxTemplate)
		lookup = templateSelectQuery(s.Dialect).Where(
			s.Functions.LOWER(goqu.I("handle")).Eq(strings.ToLower(handle)),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableTemplateFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableTemplateFields() map[string]string {
	return map[string]string{
		"created_at":   "created_at",
		"createdat":    "created_at",
		"deleted_at":   "deleted_at",
		"deletedat":    "deleted_at",
		"handle":       "handle",
		"id":           "id",
		"language":     "language",
		"last_used_at": "last_used_at",
		"lastusedat":   "last_used_at",
		"template":     "template",
		"type":         "type",
		"updated_at":   "updated_at",
		"updatedat":    "updated_at",
	}
}

// collectTemplateCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectTemplateCursorValues(res *systemType.Template, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "language":
					cur.Set(c.Column, res.Language, c.Descending)
				case "type":
					cur.Set(c.Column, res.Type, c.Descending)
				case "template":
					cur.Set(c.Column, res.Template, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				case "lastUsedAt":
					cur.Set(c.Column, res.LastUsedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkTemplateConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkTemplateConstraints(ctx context.Context, res *systemType.Template) (err error) {
	err = func() (err error) {

		// handling string type as default
		if len(res.Handle) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupTemplateByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	return nil
}

// CreateUser creates one or more rows in user collection
//
// This function is auto-generated
func (s *Store) CreateUser(ctx context.Context, rr ...*systemType.User) (err error) {
	for i := range rr {
		if err = s.checkUserConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, userInsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpdateUser updates one or more existing entries in user collection
//
// This function is auto-generated
func (s *Store) UpdateUser(ctx context.Context, rr ...*systemType.User) (err error) {
	for i := range rr {
		if err = s.checkUserConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, userUpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// UpsertUser updates one or more existing entries in user collection
//
// This function is auto-generated
func (s *Store) UpsertUser(ctx context.Context, rr ...*systemType.User) (err error) {
	for i := range rr {
		if err = s.checkUserConstraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, userUpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// DeleteUser Deletes one or more entries from user collection
//
// This function is auto-generated
func (s *Store) DeleteUser(ctx context.Context, rr ...*systemType.User) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, userDeleteQuery(s.Dialect, userPrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// DeleteUserByID deletes single entry from user collection
//
// This function is auto-generated
func (s *Store) DeleteUserByID(ctx context.Context, id uint64) error {
	return s.Exec(ctx, userDeleteQuery(s.Dialect, goqu.Ex{
		"id": id,
	}))
}

// TruncateUsers Deletes all rows from the user collection
func (s Store) TruncateUsers(ctx context.Context) error {
	return s.Exec(ctx, userTruncateQuery(s.Dialect))
}

// SearchUsers returns (filtered) set of Users
//
// This function is auto-generated
func (s *Store) SearchUsers(ctx context.Context, f systemType.UserFilter) (set systemType.UserSet, _ systemType.UserFilter, err error) {

	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
	f.PrevPage, f.NextPage = nil, nil

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return
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
	// Original are passed to the etchFullPageOfUsers fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfUsers(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// fetchFullPageOfUsers collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
//
// This function is auto-generated
func (s *Store) fetchFullPageOfUsers(
	ctx context.Context,
	filter systemType.UserFilter,
	sort filter.SortExprSet,
) (set []*systemType.User, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*systemType.User

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = filter.Limit

		reqItems = filter.Limit

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = filter.PageCursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool

		tryFilter systemType.UserFilter
	)

	set = make([]*systemType.User, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.QueryUsers(ctx, tryFilter); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 || !hasNext {
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
			tryFilter.PageCursor = s.collectUserCursorValues(set[collected-1], filter.Sort...)

			// Copy reverse flag from sorting
			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
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
		prev = s.collectUserCursorValues(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.collectUserCursorValues(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryUsers queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) QueryUsers(
	ctx context.Context,
	f systemType.UserFilter,
) (_ []*systemType.User, more bool, err error) {
	var (
		ok bool

		set         = make([]*systemType.User, 0, DefaultSliceCapacity)
		res         *systemType.User
		aux         *auxUser
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression

		sortExpr []exp.OrderedExpression
	)

	if s.Filters.User != nil {
		// extended filter set
		tExpr, f, err = s.Filters.User(s, f)
	} else {
		// using generated filter
		tExpr, f, err = UserFilter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for User: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	// paging feature is enabled
	if f.PageCursor != nil {
		if tExpr, err = cursorWithSorting(f.PageCursor, s.sortableUserFields()); err != nil {
			return
		} else {
			expr = append(expr, tExpr...)
		}
	}

	query := userSelectQuery(s.Dialect).Where(expr...)

	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.sortableUserFields()); err != nil {
		err = fmt.Errorf("could generate order expression for User: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query User: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query User: %w", err)
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	for rows.Next() {
		if err = rows.Err(); err != nil {
			err = fmt.Errorf("could not query User: %w", err)
			return
		}

		aux = new(auxUser)
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for User: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode User: %w", err)
			return
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}

		set = append(set, res)
	}

	return set, f.Limit > 0 && count >= f.Limit, err

}

// LookupUserByID searches for user by ID
//
// It returns user even if deleted or suspended
//
// This function is auto-generated
func (s *Store) LookupUserByID(ctx context.Context, id uint64) (_ *systemType.User, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxUser)
		lookup = userSelectQuery(s.Dialect).Where(
			goqu.I("id").Eq(id),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupUserByEmail searches for user by email
//
// It returns only valid user (not deleted, not suspended)
//
// This function is auto-generated
func (s *Store) LookupUserByEmail(ctx context.Context, email string) (_ *systemType.User, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxUser)
		lookup = userSelectQuery(s.Dialect).Where(
			s.Functions.LOWER(goqu.I("email")).Eq(strings.ToLower(email)),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupUserByHandle searches for user by handle
//
// It returns only valid user (not deleted, not suspended)
//
// This function is auto-generated
func (s *Store) LookupUserByHandle(ctx context.Context, handle string) (_ *systemType.User, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxUser)
		lookup = userSelectQuery(s.Dialect).Where(
			s.Functions.LOWER(goqu.I("handle")).Eq(strings.ToLower(handle)),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// LookupUserByUsername searches for user by username
//
// It returns only valid user (not deleted, not suspended)
//
// This function is auto-generated
func (s *Store) LookupUserByUsername(ctx context.Context, username string) (_ *systemType.User, err error) {
	var (
		rows   *sql.Rows
		aux    = new(auxUser)
		lookup = userSelectQuery(s.Dialect).Where(
			s.Functions.LOWER(goqu.I("username")).Eq(strings.ToLower(username)),
			goqu.I("deleted_at").IsNull(),
		).Limit(1)
	)

	rows, err = s.Query(ctx, lookup)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err = aux.scan(rows); err != nil {
		return
	}

	return aux.decode()
}

// sortableUserFields returns all <no value> columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) sortableUserFields() map[string]string {
	return map[string]string{
		"created_at":   "created_at",
		"createdat":    "created_at",
		"deleted_at":   "deleted_at",
		"deletedat":    "deleted_at",
		"email":        "email",
		"handle":       "handle",
		"id":           "id",
		"kind":         "kind",
		"name":         "name",
		"suspended_at": "suspended_at",
		"suspendedat":  "suspended_at",
		"updated_at":   "updated_at",
		"updatedat":    "updated_at",
		"username":     "username",
	}
}

// collectUserCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
//
// This function is auto-generated
func (s *Store) collectUserCursorValues(res *systemType.User, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		pkID bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cur.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "email":
					cur.Set(c.Column, res.Email, c.Descending)
					hasUnique = true
				case "username":
					cur.Set(c.Column, res.Username, c.Descending)
					hasUnique = true
				case "name":
					cur.Set(c.Column, res.Name, c.Descending)
				case "handle":
					cur.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true
				case "kind":
					cur.Set(c.Column, res.Kind, c.Descending)
				case "suspendedAt":
					cur.Set(c.Column, res.SuspendedAt, c.Descending)
				case "createdAt":
					cur.Set(c.Column, res.CreatedAt, c.Descending)
				case "updatedAt":
					cur.Set(c.Column, res.UpdatedAt, c.Descending)
				case "deletedAt":
					cur.Set(c.Column, res.DeletedAt, c.Descending)
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cur

}

// checkUserConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) checkUserConstraints(ctx context.Context, res *systemType.User) (err error) {
	err = func() (err error) {

		// handling string type as default
		if len(res.Email) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupUserByEmail(ctx, res.Email)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	err = func() (err error) {

		// handling string type as default
		if len(res.Handle) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupUserByHandle(ctx, res.Handle)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	err = func() (err error) {

		// handling string type as default
		if len(res.Username) == 0 {
			// skip check on empty values
			return nil
		}

		if res.DeletedAt != nil {
			// skip check if value is not nil
			return nil
		}

		ex, err := s.LookupUserByUsername(ctx, res.Username)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}

		return nil
	}()

	if err != nil {
		return
	}

	return nil
}
