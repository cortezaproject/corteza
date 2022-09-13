package rdbms

{{ template "gocode/header-gentext.tpl" }}

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store"
{{- range $path, $alias :=  .imports }}
    {{ $alias }} {{ printf "%q" $path }}
{{- end }}
)

{{ define "extraArgs" -}}
{{/*This is temporary solution until we properly implement Compose Record Store*/}}
{{- if eq . "composeRecord" }}mod *composeType.Module, {{ end -}}
{{- end }}
{{ define "extraParams" -}}
{{/*This is temporary solution until we properly implement Compose Record Store*/}}
{{- if eq . "composeRecord" }}mod, {{ end -}}
{{- end }}

var (
{{ range .types }}
	_ store.{{ .expIdentPlural }} = &Store{}
{{- end }}
)

{{- range .types }}
// Create{{ .expIdent }} creates one or more rows in {{ .ident }} collection
//
// This function is auto-generated
func (s *Store) Create{{ .expIdent }}(ctx context.Context, {{ template "extraArgs" .ident }} rr ...*{{ .goType }}) (err error) {
	for i := range rr {
		if err = s.check{{ .expIdent }}Constraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, {{ .ident }}InsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}


	return
}

// Update{{ .expIdent }} updates one or more existing entries in {{ .ident }} collection
//
// This function is auto-generated
func (s *Store) Update{{ .expIdent }}(ctx context.Context, {{ template "extraArgs" .ident }} rr ...*{{ .goType }}) (err error) {
	for i := range rr {
		if err = s.check{{ .expIdent }}Constraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, {{ .ident }}UpdateQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// Upsert{{ .expIdent }} updates one or more existing entries in {{ .ident }} collection
//
// This function is auto-generated
func (s *Store) Upsert{{ .expIdent }}(ctx context.Context, {{ template "extraArgs" .ident }} rr ...*{{ .goType }}) (err error) {
	for i := range rr {
		if err = s.check{{ .expIdent }}Constraints(ctx, rr[i]); err != nil {
			return
		}

		if err = s.Exec(ctx, {{ .ident }}UpsertQuery(s.Dialect, rr[i])); err != nil {
			return
		}
	}

	return
}

// Delete{{ .expIdent }} Deletes one or more entries from {{ .ident }} collection
//
// This function is auto-generated
func (s *Store) Delete{{ .expIdent }}(ctx context.Context, {{ template "extraArgs" .ident }} rr ...*{{ .goType }}) (err error) {
	for i := range rr {
		if err = s.Exec(ctx, {{ .ident }}DeleteQuery(s.Dialect, {{ .ident }}PrimaryKeys(rr[i]))); err != nil {
			return
		}
	}

	return nil
}

// Delete{{ .expIdent }}ByID deletes single entry from {{ .ident }} collection
//
// This function is auto-generated
func (s *Store) {{ .api.deleteByPK.expFnIdent }}(ctx context.Context, {{ template "extraArgs" .ident }} {{ range .api.deleteByPK.attributes }}{{ .ident }} {{ .goType }},{{ end }}) error {
	return s.Exec(ctx, {{ .ident }}DeleteQuery(s.Dialect, goqu.Ex{
	{{- range .api.deleteByPK.attributes }}
		{{ printf "%q" .storeIdent }}: {{ .ident }},
	{{- end }}
	}))
}

// Truncate{{ .expIdentPlural }} Deletes all rows from the {{ .ident }} collection
func (s Store) Truncate{{ .expIdentPlural }}(ctx context.Context, {{ template "extraArgs" .ident }}) error {
	return s.Exec(ctx, {{ .ident }}TruncateQuery(s.Dialect))
}

// Search{{ .expIdentPlural }} returns (filtered) set of {{ .expIdentPlural }}
//
// This function is auto-generated
func (s *Store) Search{{ .expIdentPlural }}(ctx context.Context, {{ template "extraArgs" .ident }} f {{ .goFilterType }}) (set {{ .goSetType }}, _ {{ .goFilterType }}, err error) {
	{{ if .features.paging }}
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
	// Original are passed to the etchFullPageOf{{ .expIdentPlural }} fn used for cursor creation;
	// direction information it MUST keep the initial
	sort := f.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if f.PageCursor != nil && f.PageCursor.ROrder {
		sort.Reverse()
	}

	set, f.PrevPage, f.NextPage, err = s.fetchFullPageOf{{ .expIdentPlural }}(ctx, f, sort)

	f.PageCursor = nil
	if err != nil {
		return nil, f, err
	}
	{{ else }}
	set, _, err = s.Query{{ .expIdentPlural }}(ctx, f)
	if err != nil {
		return nil, f, err
	}

	{{ end }}
	return set, f, nil
}

{{ if .features.paging }}
// fetchFullPageOf{{ .expIdentPlural }} collects all requested results.
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
func (s *Store) fetchFullPageOf{{ .expIdentPlural }}(
	ctx context.Context,
	filter {{ .goFilterType }},
	sort filter.SortExprSet,
) (set []*{{ .goType }}, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*{{ .goType }}

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

		tryFilter {{ .goFilterType }}
	)

	set = make([]*{{ .goType }}, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		// Copy filter & apply custom sorting that might be affected by cursor
		tryFilter = filter
		tryFilter.Sort = sort

		if limit > 0 {
			// fetching + 1 to peak ahead if there are more items
			// we can fetch (next-page cursor)
			tryFilter.Limit = limit + 1
		}

		if aux, hasNext, err = s.Query{{ .expIdentPlural }}(ctx, tryFilter); err != nil {
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
			tryFilter.PageCursor = s.{{ .api.collectCursorValues.fnIdent }}(set[collected-1], filter.Sort...)

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
		prev = s.{{ .api.collectCursorValues.fnIdent }}(set[0], filter.Sort...)
		prev.ROrder = true
		prev.LThen = !filter.Sort.Reversed()
	}

	if hasNext {
		next = s.{{ .api.collectCursorValues.fnIdent }}(set[collected-1], filter.Sort...)
		next.LThen = filter.Sort.Reversed()
	}

	return set, prev, next, nil
}
{{ end }}

// Query{{ .expIdentPlural }} queries the database, converts and checks each row and returns collected set
//
// With generics, we can remove this per-resource-generated function
// and replace it with a single utility fetcher
//
// This function is auto-generated
func (s *Store) Query{{ .expIdentPlural }}(
	ctx context.Context,
	f {{ .goFilterType }},
) (_ []*{{ .goType }}, more bool, err error) {
	var (
	{{ if .features.checkFn }}
		ok          bool
	{{ end }}
		set         = make([]*{{ .goType }}, 0, DefaultSliceCapacity)
		res         *{{ .goType }}
		aux         *{{ .auxIdent }}
		rows        *sql.Rows
		count       uint
		expr, tExpr []goqu.Expression
	{{ if .features.sorting }}
		sortExpr []exp.OrderedExpression
	{{ end }}
	)

	if s.Filters.{{ .expIdent }} != nil {
		// extended filter set
		tExpr, f, err = s.Filters.{{ .expIdent }}(s, f)
	} else {
		// using generated filter
		tExpr, f, err = {{ .expIdent }}Filter(f)
	}

	if err != nil {
		err = fmt.Errorf("could generate filter expression for {{ .expIdent }}: %w", err)
		return
	}

	expr = append(expr, tExpr...)

	{{ if .features.paging }}
	// paging feature is enabled
	if f.PageCursor != nil {
        {{- if .features.sorting }}
        if tExpr, err = cursorWithSorting(f.PageCursor, s.{{ .api.sortableFields.fnIdent }}()); err != nil {
        {{- else }}
        if tExpr, err = cursor(f.PageCursor); err != nil {
        {{- end }}
            return
        } else {
            expr = append(expr, tExpr...)
        }
	}
	{{ end }}


	query := {{ .ident }}SelectQuery(s.Dialect).Where(expr...)

	{{ if .features.sorting }}
	// sorting feature is enabled
	if sortExpr, err = order(f.Sort, s.{{ .api.sortableFields.fnIdent }}()); err != nil {
		err = fmt.Errorf("could generate order expression for {{ .expIdent }}: %w", err)
		return
	}

	if len(sortExpr) > 0 {
		query = query.Order(sortExpr...)
	}
	{{ end }}

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}


	rows, err = s.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("could not query {{ .expIdent }}: %w", err)
		return
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("could not query {{ .expIdent }}: %w", err)
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
			err = fmt.Errorf("could not query {{ .expIdent }}: %w", err)
			return
		}

		aux = new({{ .auxIdent }})
		if err = aux.scan(rows); err != nil {
			err = fmt.Errorf("could not scan rows for {{ .expIdent }}: %w", err)
			return
		}

		count++
		if res, err = aux.decode(); err != nil {
			err = fmt.Errorf("could not decode {{ .expIdent }}: %w", err)
			return
		}

		{{ if .features.checkFn }}
		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if f.Check != nil {
			if ok, err = f.Check(res); err != nil {
				return
			} else if !ok {
				continue
			}
		}
		{{ end }}

		set = append(set, res)
	}

	{{ if .features.paging }}
	return set, f.Limit > 0 && count >= f.Limit, err
	{{ else }}
	return set, false, err
	{{ end }}
}

{{- range .api.lookups }}
	{{ if .description	}}{{ .description }}{{ end }}
	//
	// This function is auto-generated
	func (s *Store) {{ .expFnIdent }}(ctx context.Context, {{ template "extraArgs" .ident }} {{ range .args }}{{ .ident }} {{ .goType }}, {{ end }}) (_ *{{ .returnType }}, err error) {
		var (
			rows *sql.Rows
			aux = new({{ .auxIdent }})
			lookup = {{ .ident }}SelectQuery(s.Dialect).Where(
				{{- range .args }}
					{{- if .ignoreCase }}
					s.Functions.LOWER(goqu.I({{ printf "%q" .storeIdent }})).Eq(strings.ToLower({{ .ident }})),
					{{- else }}
					goqu.I({{ printf "%q" .storeIdent }}).Eq({{ .ident }}),
					{{- end }}
				{{- end }}
				{{- range .nullConstraint }}
					goqu.I({{ printf "%q" . }}).IsNull(),
				{{- end }}
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
{{- end }}

{{ with .api.sortableFields }}
// {{ .fnIdent }} returns all {{ .expIdent }} columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
//
// This function is auto-generated
func (Store) {{ .fnIdent }}() map[string]string {
	return map[string]string{
		{{- range $k, $v := .fields }}
			{{ printf "%q: %q" $k $v }},
		{{- end }}
	}
}
{{ end }}

{{ with .api.collectCursorValues }}
// {{ .fnIdent }} collects values from the given resource that and sets them to the cursor
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
func (s *Store) {{ .fnIdent }}(res *{{ .goType }}, cc ...*filter.SortExpr) *filter.PagingCursor {
	{{- if .fields }}
	var (
		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		{{ range .primaryKeys }}
		pk{{ .expIdent }} bool
		{{- end }}

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				{{- range .fields }}
				case {{ printf "%q" .ident }}:
					cur.Set(c.Column, res.{{ .expIdent }}, c.Descending)

					{{- if .primaryKey }}
						pk{{ .expIdent }} = true
					{{- else if .unique }}
						hasUnique = true
					{{- end }}
				{{- end }}
				}
			}
		}
	)

	collect(cc...)
	{{- range .primaryKeys }}
	if !hasUnique || !pk{{ .expIdent }} {
			collect(&filter.SortExpr{Column: {{ printf "%q" .ident }}, Descending: {{ if .descending }}true{{ else }}false{{ end }}})
	}
	{{- end }}

	return cur
	{{ else }}
	return nil
	{{- end }}
}
{{ end }}


{{ with .api.checkConstraints }}
// {{ .fnIdent }} performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant, but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
//
// This function is auto-generated
func (s *Store) {{ .fnIdent }}(ctx context.Context, res *{{ .goType }}) (err error) {
	{{- range .checks }}
	err = func() (err error) {
		{{- range .fields }}
			{{ if eq .goType "uint64" }}
			if res.{{ .expIdent }} == 0 {
				// skip check on empty values
				return nil
			}
			{{ else }}
			// handling string type as default
			if len(res.{{ .expIdent }}) == 0 {
				// skip check on empty values
				return nil
			}
			{{ end }}
		{{ end }}

		{{- range .nullConstraint }}
		if res.{{ .expIdent }} != nil {
			// skip check if value is not nil
			return nil
		}
  	{{ end }}

		ex, err := s.{{ .lookupFnIdent }}(ctx, {{ range .fields }}res.{{ .expIdent }},{{ end }})
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
	{{ end }}
	return nil
}
{{ end }}
{{ end }}
