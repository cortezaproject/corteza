package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: {{ .Source }}
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/pkg/errors"
{{- if $.Search.EnablePaging }}
	"github.com/cortezaproject/corteza-server/pkg/filter"
{{- if and $.Search.Enable (not $.Search.Custom) }}
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
{{- end }}
{{- end }}
{{- range .Import }}
    {{ normalizeImport . }}
{{- end }}
)

var _ = errors.Is

{{/*
const (
	{{- if .Create.Enable }}
	TriggerBefore{{ export $.Types.Singular }}Create triggerKey = "{{ unexport $.Types.Singular }}BeforeCreate"
	{{- end }}
	{{- if .Update.Enable }}
	TriggerBefore{{ export $.Types.Singular }}Update triggerKey = "{{ unexport $.Types.Singular }}BeforeUpdate"
	{{- end }}
	{{- if .Upsert.Enable }}
	TriggerBefore{{ export $.Types.Singular }}Upsert triggerKey = "{{ unexport $.Types.Singular }}BeforeUpsert"
	{{- end }}
	{{- if .Delete.Enable }}
	TriggerBefore{{ export $.Types.Singular }}Delete triggerKey = "{{ unexport $.Types.Singular }}BeforeDelete"
	{{- end }}
)
*/}}

{{ if $.Search.Enable }}
{{ if $.Search.Custom }}
// {{ toggleExport .Search.Export "Search" $.Types.Plural }} not generated
// {search: {custom:true}}
{{ else }}
// {{ toggleExport .Search.Export "Search" $.Types.Plural }} returns all matching rows
//
// This function calls convert{{ export $.Types.Singular }}Filter with the given
// {{ $.Types.GoFilterType }} and expects to receive a working squirrel.SelectBuilder
func (s Store) {{ toggleExport .Search.Export "Search" $.Types.Plural }}(ctx context.Context{{ template "extraArgsDef" . }}, f {{ $.Types.GoFilterType }}) ({{ $.Types.GoSetType }}, {{ $.Types.GoFilterType }}, error) {
	var (
		err error
		set []*{{ $.Types.GoType }}
		q squirrel.SelectBuilder
	)

	return set, f, func() error {
	{{- if .RDBMS.CustomFilterConverter }}
		q, err = s.convert{{ export $.Types.Singular }}Filter({{ template "extraArgsCallFirst" . }}f)
		if err != nil {
			return err
		}
	{{- else }}
		q = s.{{ unexport $.Types.Plural }}SelectBuilder()
	{{- end }}

	{{ if $.Search.EnablePaging }}
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
		{{- range $.RDBMS.Columns.PrimaryKeyFields }}
			if f.Sort.Get({{ printf "%q" .Column  }}) == nil {
				f.Sort = append(f.Sort, &filter.SortExpr{
					Column: {{ printf "%q" .Column  }},
					Descending: {{ if .SortDescending }}true{{ else }}f.Sort.LastDescending(){{ end }},
				})
			}
		{{- end }}

		// Cloned sorting instructions for the actual sorting
		// Original are passed to the fetchFullPageOfUsers fn used for cursor creation so it MUST keep the initial
		// direction information
		sort := f.Sort.Clone()

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		if f.PageCursor != nil && f.PageCursor.ROrder {
			sort.Reverse()
		}

	{{ end }}

	{{ if .RDBMS.CustomSortConverter }}
		if q, err = s.{{ unexport $.Types.Plural }}Sorter({{ template "extraArgsCallFirst" . }}q, sort); err != nil {
			return err
		}

	{{ else if 	$.Search.EnableSorting }}
		// Apply sorting expr from filter to query
		if q, err = setOrderBy(q, sort, s.sortable{{ export $.Types.Singular }}Columns()); err != nil {
			return err
		}
	{{ end }}

	{{- if $.Search.EnablePaging }}
		set, f.PrevPage, f.NextPage, err = s.{{ unexport "fetchFullPageOf" $.Types.Plural  }}(
			ctx{{ template "extraArgsCall" . }},
			q, f.Sort, f.PageCursor,
			f.Limit,
			{{ if $.Search.EnableFilterCheckFn }}f.Check{{ else }}nil{{ end }},
			func(cur *filter.PagingCursor) squirrel.Sqlizer {
				return builders.CursorCondition(cur, nil)
			},
		)

		if err != nil {
			return err
		}

		f.PageCursor = nil
		return nil
	{{- else }}
		set, err = s.{{ export "query" $.Types.Plural }}(ctx{{ template "extraArgsCall" . }}, q, {{ if $.Search.EnableFilterCheckFn }}f.Check{{else}}nil{{ end }})
		return err
	{{- end }}
	}()
}
{{ end }}
{{ end }}


{{ if $.Search.EnablePaging }}
// {{ unexport "fetchFullPageOf" $.Types.Plural  }} collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) {{ unexport "fetchFullPageOf" $.Types.Plural  }} (
	ctx context.Context{{ template "extraArgsDef" . }},
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*{{ $.Types.GoType }}) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*{{ $.Types.GoType }}, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*{{ $.Types.GoType }}

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

	set = make([]*{{ $.Types.GoType }}, 0, DefaultSliceCapacity)

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

		if aux, err = s.{{ export "query" $.Types.Plural }}(ctx{{ template "extraArgsCall" . }}, tryQuery, check); err != nil {
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
			cursor = s.collect{{ export $.Types.Singular }}CursorValues({{ template "extraArgsCallFirst" . }}set[collected-1], sort...)

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
		prev = s.collect{{ export $.Types.Singular }}CursorValues({{ template "extraArgsCallFirst" . }}set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collect{{ export $.Types.Singular }}CursorValues({{ template "extraArgsCallFirst" . }}set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
}
{{ end }}


// {{ export "query" $.Types.Plural }} queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) {{ export "query" $.Types.Plural }} (
	ctx context.Context{{ template "extraArgsDef" . }},
	q squirrel.Sqlizer,
	check func(*{{ $.Types.GoType }}) (bool, error),
) ([]*{{ $.Types.GoType }}, error) {
	var (
		set = make([]*{{ $.Types.GoType }}, 0, DefaultSliceCapacity)
		res  *{{ $.Types.GoType }}

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internal{{ export $.Types.Singular }}RowScanner({{ template "extraArgsCallFirst" . }}rows)
		}

		if err != nil {
			return nil, err
		}

	{{ if $.Search.EnableFilterCheckFn }}
		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if check != nil {
			if chk, err := check(res); err != nil {
				return nil, err
			} else if !chk {
				continue
			}
		}
	{{ end }}
		set = append(set, res)
	}

{{ if .RDBMS.CustomPostLoadProcessor }}
	if err = s.{{ unexport $.Types.Singular }}PostLoadProcessor(ctx{{ template "extraArgsCall" . }}, set...); err != nil {
		return nil, err
	}
{{end }}

	return set, rows.Err()
}




{{- range $.Lookups }}
// {{ toggleExport .Export "Lookup" $.Types.Singular "By" .Suffix }} {{ comment .Description true -}}
func (s Store) {{ toggleExport .Export "Lookup" $.Types.Singular "By" .Suffix }}(ctx context.Context{{ template "extraArgsDef" $ }}{{- range .RDBMSColumns }}, {{ cc2underscore .Field }} {{ .Type  }}{{- end }}) (*{{ $.Types.GoType }}, error) {
	return s.execLookup{{ $.Types.Singular }}(ctx{{ template "extraArgsCall" $ }}, squirrel.Eq{
    {{- range .RDBMSColumns }}
		s.preprocessColumn({{ printf "%q" .AliasedColumn }}, {{ printf "%q" .LookupFilterPreprocess }}): store.PreprocessValue({{ cc2underscore .Field }}, {{ printf "%q" .LookupFilterPreprocess }}),
    {{- end }}

    {{ range $field, $value := .Filter }}
       "{{ ($field | $.RDBMS.Columns.Find).AliasedColumn }}": {{ $value }},
    {{- end }}
    })
}
{{ end }}

{{ if .Create.Enable }}
// {{ toggleExport .Create.Export "Create" $.Types.Singular }} creates one or more rows in {{ $.RDBMS.Table }} table
func (s Store) {{ toggleExport .Create.Export "Create" $.Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $.Types.GoType }}) (err error) {
	for _, res := range rr {
		err = s.check{{ export $.Types.Singular }}Constraints(ctx {{ template "extraArgsCall" $ }}, res)
		if err != nil {
			return err
		}

{{/*
		err = s.{{ unexport $.Types.Singular }}Hook(ctx, TriggerBefore{{ export $.Types.Singular }}Create{{ template "extraArgsCall" . }}, res)
		if err != nil {
			return err
		}
*/}}

		err = s.execCreate{{ export $.Types.Plural }}(ctx, s.internal{{ export $.Types.Singular }}Encoder(res))
		if err != nil {
			return err
		}
	}

	return
}
{{ end }}

{{ if .Update.Enable }}
// {{ toggleExport .Update.Export "Update" $.Types.Singular }} updates one or more existing rows in {{ $.RDBMS.Table }}
func (s Store) {{ toggleExport .Update.Export "Update" $.Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $.Types.GoType }}) error {
	return s.partial{{ export $.Types.Singular "Update" }}(ctx{{ template "extraArgsCall" . }}, nil, rr...)
}

// partial{{ export $.Types.Singular "Update" }} updates one or more existing rows in {{ $.RDBMS.Table }}
func (s Store) partial{{ export $.Types.Singular "Update" }}(ctx context.Context{{ template "extraArgsDef" . }}, onlyColumns []string, rr ... *{{ $.Types.GoType }}) (err error) {
	for _, res := range rr {
		err = s.check{{ export $.Types.Singular }}Constraints(ctx {{ template "extraArgsCall" $ }}, res)
		if err != nil {
			return err
		}

{{/*
		err = s.{{ unexport $.Types.Singular }}Hook(ctx, TriggerBefore{{ export $.Types.Singular }}Update{{ template "extraArgsCall" . }}, res)
		if err != nil {
			return err
		}
*/}}
		err = s.execUpdate{{ export $.Types.Plural }}(
			ctx,
			{{ template "filterByPrimaryKeys" $.RDBMS.Columns.PrimaryKeyFields }},
			s.internal{{ export $.Types.Singular }}Encoder(res).Skip(
				{{- range $.RDBMS.Columns.PrimaryKeyFields -}}
					{{ printf "%q" .Column  }},
				{{- end -}}
		).Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}
{{ end }}

{{ if .Upsert.Enable }}
// {{ toggleExport .Upsert.Export "Upsert" $.Types.Singular }} updates one or more existing rows in {{ $.RDBMS.Table }}
func (s Store) {{ toggleExport .Upsert.Export "Upsert" $.Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $.Types.GoType }}) (err error) {
	for _, res := range rr {
		err = s.check{{ export $.Types.Singular }}Constraints(ctx {{ template "extraArgsCall" $ }}, res)
		if err != nil {
			return err
		}

{{/*
		err = s.{{ unexport $.Types.Singular }}Hook(ctx, TriggerBefore{{ export $.Types.Singular }}Upsert{{ template "extraArgsCall" . }}, res)
		if err != nil {
			return err
		}
*/}}

		err = s.execUpsert{{ export $.Types.Plural }}(ctx, s.internal{{ export $.Types.Singular }}Encoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}
{{ end }}

{{ if .Delete.Enable }}
// {{ toggleExport .Delete.Export "Delete" $.Types.Singular }} Deletes one or more rows from {{ $.RDBMS.Table }} table
func (s Store) {{ toggleExport .Delete.Export "Delete" $.Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $.Types.GoType }}) (err error) {
	for _, res := range rr {
{{/*
		err = s.{{ unexport $.Types.Singular }}Hook(ctx, TriggerBefore{{ export $.Types.Singular }}Delete{{ template "extraArgsCall" . }}, res)
		if err != nil {
			return err
		}
*/}}
		err= s.execDelete{{ export $.Types.Plural }}(ctx,{{ template "filterByPrimaryKeys" $.RDBMS.Columns.PrimaryKeyFields }})
		if err != nil {
			return err
		}
	}

	return nil
}

// {{ toggleExport .Delete.Export "Delete" $.Types.Singular "By" }}{{ template "primaryKeySuffix" $.RDBMS.Columns }} Deletes row from the {{ $.RDBMS.Table }} table
func (s Store) {{ toggleExport .Delete.Export "Delete" $.Types.Singular "By" }}{{ template "primaryKeySuffix" $.RDBMS.Columns }}(ctx context.Context{{ template "extraArgsDef" . }}{{ template "primaryKeyArgsDef" $.RDBMS.Columns }}) error {
	return s.execDelete{{ export $.Types.Plural }}(ctx, {{ template "filterByPrimaryKeysWithArgs" $.RDBMS.Columns.PrimaryKeyFields }})
}
{{ end }}

// {{ toggleExport .Truncate.Export "Truncate" $.Types.Plural }} Deletes all rows from the {{ $.RDBMS.Table }} table
func (s Store) {{ toggleExport .Truncate.Export "Truncate" $.Types.Plural }}(ctx context.Context{{ template "extraArgsDef" . }}) error {
	return s.Truncate(ctx, s.{{ unexport $.Types.Singular }}Table())
}

// execLookup{{ $.Types.Singular }} prepares {{ $.Types.Singular }} query and executes it,
// returning {{ $.Types.GoType }} (or error)
func (s Store) execLookup{{ $.Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, cnd squirrel.Sqlizer) (res *{{ $.Types.GoType }}, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.{{ unexport $.Types.Plural }}SelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internal{{ $.Types.Singular }}RowScanner({{ template "extraArgsCallFirst" . }}row)
	if err != nil {
		return
	}

{{ if .RDBMS.CustomPostLoadProcessor }}
	if err = s.{{ unexport $.Types.Singular }}PostLoadProcessor(ctx{{ template "extraArgsCall" . }}, res); err != nil {
		return nil, err
	}
{{ end }}

	return res, nil
}

{{ if .Create.Enable }}
// execCreate{{ export $.Types.Plural }} updates all matched (by cnd) rows in {{ $.RDBMS.Table }} with given data
func (s Store) execCreate{{ export $.Types.Plural }}(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.{{ unexport $.Types.Singular }}Table()).SetMap(payload))
}
{{ end }}

{{ if .Update.Enable }}
// execUpdate{{ export $.Types.Plural }} updates all matched (by cnd) rows in {{ $.RDBMS.Table }} with given data
func (s Store) execUpdate{{ export $.Types.Plural }}(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.{{ unexport $.Types.Singular }}Table({{ printf "%q" .RDBMS.Alias }})).Where(cnd).SetMap(set))
}
{{ end }}

{{ if .Upsert.Enable }}
// execUpsert{{ export $.Types.Plural }} inserts new or updates matching (by-primary-key) rows in {{ $.RDBMS.Table }} with given data
func (s Store) execUpsert{{ export $.Types.Plural }}(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.{{ unexport $.Types.Singular }}Table(),
		set,
{{ range $.RDBMS.Columns }}
	{{- if or .IsPrimaryKey -}}
		s.preprocessColumn({{ printf "%q" .Column }}, {{ printf "%q" .LookupFilterPreprocess }}),
	{{ end }}
{{- end }}
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}
{{ end }}

{{ if .Delete.Enable }}
// execDelete{{ export $.Types.Plural }} Deletes all matched (by cnd) rows in {{ $.RDBMS.Table }} with given data
func (s Store) execDelete{{ export $.Types.Plural }}(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.{{ unexport $.Types.Singular }}Table({{ printf "%q" .RDBMS.Alias }})).Where(cnd))
}
{{ end }}

func (s Store) internal{{ export $.Types.Singular }}RowScanner({{ template "extraArgsDefFirst" . }}row rowScanner) (res *{{ $.Types.GoType }}, err error) {
	res = &{{ $.Types.GoType }}{}

	if _, has := s.config.RowScanners[{{ printf "%q" (unexport $.Types.Singular) }}]; has {
		scanner := s.config.RowScanners[{{ printf "%q" (unexport $.Types.Singular) }}].(func({{ template "extraArgsDefFirst" . }}_ rowScanner, _ *{{ $.Types.GoType }}) error)
		err = scanner({{ template "extraArgsCallFirst" . }}row, res)
	} else {
	{{- if .RDBMS.CustomRowScanner }}
		err = s.scan{{ $.Types.Singular }}Row({{ template "extraArgsCallFirst" . }}row, res)
	{{- else }}
		err = row.Scan(
		{{- range $.RDBMS.Columns }}
			&res.{{ .Field }},
		{{- end }}
		)
	{{- end }}
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan {{ unexport $.Types.Singular }} db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// Query{{ export $.Types.Plural }} returns squirrel.SelectBuilder with set table and all columns
func (s Store) {{ unexport $.Types.Plural }}SelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.{{ unexport $.Types.Singular }}Table({{ printf "%q" .RDBMS.Alias }}), s.{{ unexport $.Types.Singular }}Columns({{ printf "%q" $.RDBMS.Alias }})...)
}

// {{ unexport $.Types.Singular }}Table name of the db table
func (Store) {{ unexport $.Types.Singular }}Table(aa ... string) string {
		var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "{{ $.RDBMS.Table }}" + alias
}

// {{ $.Types.Singular }}Columns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) {{ unexport $.Types.Singular }}Columns(aa ... string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
	{{- range $.RDBMS.Columns }}
		alias + "{{ .Column }}",
    {{- end }}
	}
}

// {{ printf "%v" .Search }}

{{ if $.Search.EnableSorting }}
// sortable{{ $.Types.Singular }}Columns returns all {{ $.Types.Singular }} columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortable{{ $.Types.Singular }}Columns() map[string]string {
	return map[string]string{
	{{ range $.RDBMS.Columns }}
		{{- if .IsSortable -}}
		"{{ toLower .Column }}": "{{ .Column }}",
		{{- if not (eq (.Field|toLower) (.Column|toLower)) }}
			"{{ toLower .Field }}":  "{{ .Column }}",
		{{ end -}}
		{{ end -}}
    {{- end }}
	}
}
{{ end }}

// internal{{ export $.Types.Singular }}Encoder encodes fields from {{ $.Types.GoType }} to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encode{{ export $.Types.Singular }}
// func when rdbms.customEncoder=true
func (s Store) internal{{ export $.Types.Singular }}Encoder(res *{{ $.Types.GoType }}) store.Payload {
{{- if .RDBMS.CustomEncoder }}
	return s.encode{{ export $.Types.Singular }}(res)
{{- else }}
	return store.Payload{
    {{- range $.RDBMS.Columns }}
		"{{ .Column }}": res.{{ .Field }},
    {{- end }}
	}
{{- end }}
}

{{ if and $.Search.EnablePaging (not $.RDBMS.CustomCursorCollector) }}
// collect{{ export $.Types.Singular }}CursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collect{{ export $.Types.Singular }}CursorValues({{ template "extraArgsDef" $ }}res *{{ $.Types.GoType }}, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		// All known primary key columns
		{{ range $.RDBMS.Columns.PrimaryKeyFields }}
		pk{{ export .Column }} bool
		{{ end }}

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
			{{- range $.RDBMS.Columns }}
		        {{- if or .IsSortable .IsUnique .IsPrimaryKey -}}
				case "{{ .Column }}":
					cursor.Set(c.Column, res.{{ .Field }}, c.Descending)
					{{ if .IsUnique  -}}
					hasUnique = true
					{{ end }}
					{{ if .IsPrimaryKey -}}
					pk{{ export .Column }} = true
					{{ end }}
				{{- end }}
			{{- end }}
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !({{ range $.RDBMS.Columns.PrimaryKeyFields }}pk{{ export .Column }} && {{ end }} true) {
		collect({{ range $.RDBMS.Columns.PrimaryKeyFields }}&filter.SortExpr{Column: "{{ .Column }}", Descending: {{ .SortDescending }}},{{ end }})
	}

	return cursor
}
{{ end }}

// check{{ export $.Types.Singular }}Constraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) check{{ export $.Types.Singular }}Constraints(ctx context.Context{{ template "extraArgsDef" $ }}, res *{{ $.Types.GoType }}) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true
{{- range $.Lookups }}
	{{ if .UniqueConstraintCheck }}
	{{- range .RDBMSColumns }}
		{{ if eq .Type "uint64" }}
		valid = valid && res.{{ .Field }} > 0
		{{ else if eq .Type "string" }}
		valid = valid && len(res.{{ .Field }}) > 0
		{{ else }}
		// can not check field {{ .Field }} with unsupported type: {{ .Type }}
		{{ end }}
	{{- end }}
	{{- end }}
{{- end }}

	if !valid {
		return nil
	}


{{- range $.Lookups }}
	{{ if .UniqueConstraintCheck }}
	{
		ex, err := s.{{ toggleExport .Export "Lookup" $.Types.Singular "By" .Suffix }}(ctx{{ template "extraArgsCall" $ }}{{- range .RDBMSColumns }}, res.{{ .Field }} {{- end }})
		if err == nil && ex != nil {{- range $.RDBMS.Columns.PrimaryKeyFields }} && ex.{{ .Field }} != res.{{ .Field }} {{ end }} {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}
	{{ end }}
{{ end }}

	return nil
}

{{/*
func (s *Store) {{ unexport $.Types.Singular }}Hook(ctx context.Context, key triggerKey{{ template "extraArgsDef" . }}, res *{{ $.Types.GoType }}) error {
	if fn, has := s.config.TriggerHandlers[key]; has {
		return fn.(func (ctx context.Context, s *Store{{ template "extraArgsDef" . }}, res *{{ $.Types.GoType }}) error)(ctx, s{{ template "extraArgsCall" . }}, res)
	}

	return nil
}
*/}}


{{/* ************************************************************ */}}

{{- define "filterByPrimaryKeys" -}}
    squirrel.Eq{
    {{ range . -}}
		s.preprocessColumn({{ printf "%q" .AliasedColumn }}, {{ printf "%q" .LookupFilterPreprocess }}): store.PreprocessValue(res.{{ .Field }}, {{ printf "%q" .LookupFilterPreprocess }}),
    {{- end }}
    }
{{- end -}}

{{- define "filterByPrimaryKeysWithArgs" -}}
    squirrel.Eq{
    {{ range . -}}
		s.preprocessColumn({{ printf "%q" .AliasedColumn }}, {{ printf "%q" .LookupFilterPreprocess }}): store.PreprocessValue({{ .Arg }}, {{ printf "%q" .LookupFilterPreprocess }}),
    {{ end }}
    }
{{- end -}}
