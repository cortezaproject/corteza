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
	"errors"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
{{- if $.Search.EnablePaging }}
	"github.com/cortezaproject/corteza-server/pkg/filter"
{{- end }}
{{- range $import := $.Import }}
    {{ normalizeImport $import }}
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


{{- if .RDBMS.CustomFilterConverter }}
	q, err = s.convert{{ export $.Types.Singular }}Filter({{ template "extraArgsCallFirst" . }}f)
	if err != nil {
	    return nil, f, err
	}
{{- else }}
	q = s.{{ unexport $.Types.Plural }}SelectBuilder()
{{- end }}

{{ if $.Search.EnablePaging }}
	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reversedCursor := f.PageCursor != nil && f.PageCursor.Reverse
{{ end }}

{{ if $.Search.EnableSorting }}
	// If paging with reverse cursor, change the sorting
	// direction for all columns we're sorting by
	curSort := f.Sort.Clone()
	if reversedCursor {
		curSort.Reverse()
	}
{{ else if $.Search.EnablePaging }}
	// Sorting is disabled in definition yaml file
	// {search: {enableSorting:false}}
	//
	// We still need to sort the results by primary key for paging purposes
	curSort := filter.SortExprSet{
	{{- range $.RDBMS.Columns.PrimaryKeyFields }}
		&filter.SortExpr{Column: {{ printf "%q" .Column  }}, {{ if .SortDescending }}Descending: !reversedCursor, {{ end }}},
	{{- end }}
	}
{{ end }}

	return set, f, s.config.ErrorHandler(func() error {
	{{- if $.Search.EnablePaging }}
		set, err = s.{{ unexport "fetchFullPageOf" $.Types.Plural  }}(ctx{{ template "extraArgsCall" . }}, q, curSort, f.PageCursor, f.Limit, {{ if $.Search.EnableFilterCheckFn }}f.Check{{ else }}nil{{ end }},)

		if err != nil {
			return err
		}

		if f.Limit > 0 && len(set) > 0 {
			if f.PageCursor != nil && (!f.PageCursor.Reverse || uint(len(set)) == f.Limit) {
				f.PrevPage = s.collect{{ export $.Types.Singular }}CursorValues(set[0], curSort.Columns()...)
				f.PrevPage.Reverse = true
			}

			// Less items fetched then requested by page-limit
			// not very likely there's another page
			f.NextPage = s.collect{{ export $.Types.Singular }}CursorValues(set[len(set)-1], curSort.Columns()...)
		}

		f.PageCursor = nil
		return nil
	{{- else }}
		set, _, _, err = s.{{ export "query" $.Types.Plural }}(ctx{{ template "extraArgsCall" . }}, q, {{ if $.Search.EnableFilterCheckFn }}f.Check{{else}}nil{{ end }})
		return err
	{{ end }}
	}())
}
{{ end }}


{{ if $.Search.EnablePaging }}
// {{ unexport "fetchFullPageOf" $.Types.Plural  }} collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - sorting rules (order by ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn). Function then moves cursor to the last item fetched
func (s Store) {{ unexport "fetchFullPageOf" $.Types.Plural  }} (
	ctx context.Context{{ template "extraArgsDef" . }},
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*{{ $.Types.GoType }}) (bool, error),
) ([]*{{ $.Types.GoType }}, error) {
	var (
		set = make([]*{{ $.Types.GoType }}, 0, DefaultSliceCapacity)
		aux []*{{ $.Types.GoType }}
		last *{{ $.Types.GoType }}

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedCursor = cursor != nil && cursor.Reverse

		// copy of the select builder
		tryQuery squirrel.SelectBuilder

		fetched uint
		err error
	)


{{ if .RDBMS.Columns.PrimaryKeyFields }}
	// Make sure we always end our sort by primary keys
	{{- range .RDBMS.Columns.PrimaryKeyFields }}
	if sort.Get({{ printf "%q" .Column }}) == nil {
		sort = append(sort, &filter.SortExpr{Column: {{ printf "%q" .Column }}})
	}
	{{ end }}
{{ end }}

{{ if .RDBMS.CustomSortConverter }}
	if q, err = s.{{ unexport $.Types.Plural }}Sorter({{ template "extraArgsCallFirst" . }}q, sort); err != nil {
		return nil, err
	}
{{ else if .Search.EnableSorting }}
	// Apply sorting expr from filter to query
	if q, err = setOrderBy(q, sort, s.sortable{{ export $.Types.Singular }}Columns()...); err != nil {
		return nil, err
	}
{{ else if .RDBMS.Columns.PrimaryKeyFields }}
	// Sort by primary keys by default
	if q, err = setOrderBy(q, sort, {{ range .RDBMS.Columns.PrimaryKeyFields }}"{{ .Column }}",{{ end }}); err != nil {
		return nil, err
	}
{{ end }}


	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit))
		}

		if aux, fetched, last, err = s.{{ export "query" $.Types.Plural }}(ctx{{ template "extraArgsCall" . }}, tryQuery, check); err != nil {
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
		if cursor = s.collect{{ export $.Types.Singular }}CursorValues(last, sort.Columns()...); cursor == nil {
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
) ([]*{{ $.Types.GoType }}, uint, *{{ $.Types.GoType }}, error) {
	var (
		set = make([]*{{ $.Types.GoType }}, 0, DefaultSliceCapacity)
		res  *{{ $.Types.GoType }}

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
			res, err = s.internal{{ export $.Types.Singular }}RowScanner({{ template "extraArgsCallFirst" . }}rows)
		}

		if err != nil {
			return nil, 0, nil, err
		}

	{{ if $.Search.EnableFilterCheckFn }}
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
	{{ end }}
		set = append(set, res)
	}

{{ if .RDBMS.CustomPostLoadProcessor }}
	if err = s.{{ unexport $.Types.Singular }}PostLoadProcessor(ctx{{ template "extraArgsCall" . }}, set...); err != nil {
		return nil, 0, nil, err
	}
{{end }}

	return set, fetched, res, rows.Err()
}




{{- range $lookup := $.Lookups }}
// {{ toggleExport $lookup.Export "Lookup" $.Types.Singular "By" $lookup.Suffix }} {{ comment $lookup.Description true -}}
func (s Store) {{ toggleExport $lookup.Export "Lookup" $.Types.Singular "By" $lookup.Suffix }}(ctx context.Context{{ template "extraArgsDef" $ }}{{- range $lookup.RDBMSColumns }}, {{ cc2underscore .Field }} {{ .Type  }}{{- end }}) (*{{ $.Types.GoType }}, error) {
	return s.execLookup{{ $.Types.Singular }}(ctx{{ template "extraArgsCall" $ }}, squirrel.Eq{
    {{- range $lookup.RDBMSColumns }}
		s.preprocessColumn({{ printf "%q" .AliasedColumn }}, {{ printf "%q" .LookupFilterPreprocess }}): s.preprocessValue({{ cc2underscore .Field }}, {{ printf "%q" .LookupFilterPreprocess }}),
    {{- end }}

    {{ range $field, $value := $lookup.Filter }}
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
	return s.config.ErrorHandler(s.{{ toggleExport .Update.Export "Partial" $.Types.Singular "Update" }}(ctx{{ template "extraArgsCall" . }}, nil, rr...))
}

// {{ toggleExport .Update.Export "Partial" $.Types.Singular "Update" }} updates one or more existing rows in {{ $.RDBMS.Table }}
func (s Store) {{ toggleExport .Update.Export "Partial" $.Types.Singular "Update" }}(ctx context.Context{{ template "extraArgsDef" . }}, onlyColumns []string, rr ... *{{ $.Types.GoType }}) (err error) {
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
				{{- range $field := $.RDBMS.Columns.PrimaryKeyFields -}}
					{{ printf "%q" $field.Column  }},
				{{- end -}}
		).Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}
{{ end }}

{{ if .Upsert.Enable }}
// {{ toggleExport .Delete.Export "Upsert" $.Types.Singular }} updates one or more existing rows in {{ $.RDBMS.Table }}
func (s Store) {{ toggleExport .Delete.Export "Upsert" $.Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $.Types.GoType }}) (err error) {
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

		err = s.config.ErrorHandler(s.execUpsert{{ export $.Types.Plural }}(ctx, s.internal{{ export $.Types.Singular }}Encoder(res)))
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
		err = s.execDelete{{ export $.Types.Plural }}(ctx,{{ template "filterByPrimaryKeys" $.RDBMS.Columns.PrimaryKeyFields }})
		if err != nil {
			return s.config.ErrorHandler(err)
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
	return s.config.ErrorHandler(s.Truncate(ctx, s.{{ unexport $.Types.Singular }}Table()))
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
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.{{ unexport $.Types.Singular }}Table()).SetMap(payload)))
}
{{ end }}

{{ if .Update.Enable }}
// execUpdate{{ export $.Types.Plural }} updates all matched (by cnd) rows in {{ $.RDBMS.Table }} with given data
func (s Store) execUpdate{{ export $.Types.Plural }}(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.{{ unexport $.Types.Singular }}Table({{ printf "%q" .RDBMS.Alias }})).Where(cnd).SetMap(set)))
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
		{{ printf "%q" .Column }},
	{{ end }}
{{- end }}
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}
{{ end }}

{{ if .Delete.Enable }}
// execDelete{{ export $.Types.Plural }} Deletes all matched (by cnd) rows in {{ $.RDBMS.Table }} with given data
func (s Store) execDelete{{ export $.Types.Plural }}(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.{{ unexport $.Types.Singular }}Table({{ printf "%q" .RDBMS.Alias }})).Where(cnd)))
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
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for {{ $.Types.Singular }}: %w", err)
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
func (Store) sortable{{ $.Types.Singular }}Columns() []string {
	return []string{
	{{ range $.RDBMS.Columns }}
		{{- if .IsSortable -}}
		"{{ .Column }}",
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
func (s Store) collect{{ export $.Types.Singular }}CursorValues(res *{{ $.Types.GoType }}, cc ...string) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{}

		hasUnique bool

		// All known primary key columns
		{{ range $.RDBMS.Columns.PrimaryKeyFields }}
		pk{{ export .Column }} bool
		{{ end }}

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
			{{- range $.RDBMS.Columns }}
		        {{- if or .IsSortable .IsUnique .IsPrimaryKey -}}
				case "{{ .Column }}":
					cursor.Set(c, res.{{ .Field }}, false)
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
		collect({{ range $.RDBMS.Columns.PrimaryKeyFields }}"{{ .Column }}",{{ end }})
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
{{- range $lookup := $.Lookups }}
	{{ if $lookup.UniqueConstraintCheck }}
	{{- range $lookup.RDBMSColumns }}
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


{{- range $lookup := $.Lookups }}
	{{ if $lookup.UniqueConstraintCheck }}
	{
		ex, err := s.{{ toggleExport $lookup.Export "Lookup" $.Types.Singular "By" $lookup.Suffix }}(ctx{{ template "extraArgsCall" $ }}{{- range $lookup.RDBMSColumns }}, res.{{ .Field }} {{- end }})
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique
		} else if !errors.Is(err, store.ErrNotFound) {
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
    {{ range $field := . -}}
		s.preprocessColumn({{ printf "%q" $field.AliasedColumn }}, {{ printf "%q" $field.LookupFilterPreprocess }}): s.preprocessValue(res.{{ $field.Field }}, {{ printf "%q" $field.LookupFilterPreprocess }}),
    {{- end }}
    }
{{- end -}}

{{- define "filterByPrimaryKeysWithArgs" -}}
    squirrel.Eq{
    {{ range $field := . -}}
		s.preprocessColumn({{ printf "%q" $field.AliasedColumn }}, {{ printf "%q" $field.LookupFilterPreprocess }}): s.preprocessValue({{ $field.Arg }}, {{ printf "%q" $field.LookupFilterPreprocess }}),
    {{ end }}
    }
{{- end -}}
