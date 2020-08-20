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
	"strings"
{{- end }}
{{- range $import := $.Import }}
    {{ normalizeImport $import }}
{{- end }}
)

var _ = errors.Is

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

{{ if $.Search.Enable }}
// {{ toggleExport .Search.Export "Search" $.Types.Plural }} returns all matching rows
//
// This function calls convert{{ export $.Types.Singular }}Filter with the given
// {{ $.Types.GoFilterType }} and expects to receive a working squirrel.SelectBuilder
func (s Store) {{ toggleExport .Search.Export "Search" $.Types.Plural }}(ctx context.Context{{ template "extraArgsDef" . }}, f {{ $.Types.GoFilterType }}) ({{ $.Types.GoSetType }}, {{ $.Types.GoFilterType }}, error) {
	var scap uint

{{- if .RDBMS.CustomFilterConverter }}
	q, err := s.convert{{ export $.Types.Singular }}Filter({{ template "extraArgsCallFirst" . }}f)
	if err != nil {
	    return nil, f, err
	}
{{- else }}
	q := s.{{ unexport $.Types.Plural }}SelectBuilder()
{{- end }}

{{ if $.Search.EnablePaging }}
	scap = f.Limit

	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reverseCursor := f.PageCursor != nil && f.PageCursor.Reverse
{{ end }}

{{ if $.Search.EnableSorting }}
	if err := f.Sort.Validate(s.sortable{{ export $.Types.Singular }}Columns()...); err != nil {
		return nil, f, fmt.Errorf("could not validate sort: %v", err)
	}

	// If paging with reverse cursor, change the sorting
	// direction for all columns we're sorting by
	sort := f.Sort.Clone()
	if reverseCursor {
		sort.Reverse()
	}

	// Apply sorting expr from filter to query
	if len(sort) > 0 {
		sqlSort := make([]string, len(sort))
		for i := range sort {
			sqlSort[i] = sort[i].Column
			if sort[i].Descending {
				sqlSort[i] += " DESC"
			}
		}

		q = q.OrderBy(sqlSort...)
	}
{{ else if $.Search.EnablePaging }}
	// Sorting is disabled in definition yaml file
	// {search: {enablePaging:false}}
	//
	// We still need to sort the results by primary key for paging purposes
	sort := filter.SortExprSet{
	{{ range $.Fields }}
		{{- if or .IsPrimaryKey -}}
			&filter.SortExpr{Column: {{ printf "%q" .Column  }}, {{ if .SortDescending }}Descending: true, {{ end }}},
		{{- end }}
	{{- end }}
	}
{{ end }}

	if scap == 0 {
		scap = DefaultSliceCapacity
	}


	var (
		set = make([]*{{ $.Types.GoType }}, 0, scap)

	{{- if not $.Search.EnablePaging }}
		// Paging is disabled in definition yaml file
		// {search: {enablePaging:false}} and this allows
		// a much simpler row fetching logic
		fetch = func() error {
			var (
				res *{{ $.Types.GoType }}
				rows, err = s.Query(ctx, q)
			)

			if err != nil {
				return err
			}

			for rows.Next() {
				if rows.Err() == nil {
					res, err = s.internal{{ export $.Types.Singular }}RowScanner({{ template "extraArgsCallFirst" . }}rows)
				}

				if err != nil {
					if cerr := rows.Close(); cerr != nil {
						err = fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return err
				}

				// If check function is set, call it and act accordingly
			{{ if $.Search.EnableFilterCheckFn }}
				if f.Check != nil {
					if chk, err := f.Check(res); err != nil {
						if cerr := rows.Close(); cerr != nil {
							err = fmt.Errorf("could not close rows (%v) after check error: %w", cerr, err)
						}

						return err
					} else if !chk {
						// did not pass the check
						// go with the next row
						continue
					}
				}
			{{ end -}}

				set = append(set, res)
			}

			return rows.Close()
		}
	{{ else }}
		// fetches rows and scans them into {{ $.Types.GoType }} resource this is then passed to Check function on filter
		// to help determine if fetched resource fits or not
		//
		// Note that limit is passed explicitly and is not necessarily equal to filter's limit. We want
		// to keep that value intact.
		//
		// The value for cursor is used and set directly from/to the filter!
		//
		// It returns total number of fetched pages and modifies PageCursor value for paging
		fetchPage = func(cursor *filter.PagingCursor, limit uint) (fetched uint, err error) {
			var (
				res *{{ $.Types.GoType }}

				// Make a copy of the select query builder so that we don't change
				// the original query
				slct = q.Options()
			)

			if limit > 0 {
				slct = slct.Limit(uint64(limit))

				if cursor != nil && len(cursor.Keys()) > 0 {
					const cursorTpl = `(%s) %s (?%s)`
					op := ">"
					if cursor.Reverse {
						op = "<"
					}

					pred := fmt.Sprintf(cursorTpl, strings.Join(cursor.Keys(), ", "), op, strings.Repeat(", ?", len(cursor.Keys())-1))
					slct = slct.Where(pred, cursor.Values()...)
				}
			}

			rows, err := s.Query(ctx, slct)
			if err != nil {
				return
			}

			for rows.Next() {
				fetched++

				if rows.Err() == nil {
					res, err = s.internal{{ export $.Types.Singular }}RowScanner({{ template "extraArgsCallFirst" . }}rows)
				}

				if err != nil {
					if cerr := rows.Close(); cerr != nil {
						err = fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return
				}

				// If check function is set, call it and act accordingly
			{{ if $.Search.EnableFilterCheckFn }}
				if f.Check != nil {
					var chk bool
					if chk, err = f.Check(res); err != nil {
						if cerr := rows.Close(); cerr != nil {
							err = fmt.Errorf("could not close rows (%v) after check error: %w", cerr, err)
						}

						return
					} else if !chk {
						// did not pass the check
						// go with the next row
						continue
					}
				}
			{{ end -}}

				set = append(set, res)

				if f.Limit > 0 {
					if uint(len(set)) >= f.Limit {
						// make sure we do not fetch more than requested!
						break
					}
				}
			}

			err = rows.Close()
			return
		}

		fetch = func() error {
			var (
				// how many items were actually fetched
				fetched uint

				// starting offset & limit are from filter arg
				// note that this will have to be improved with key-based pagination
				limit = f.Limit

				// Copy cursor value
				//
				// This is where we'll start fetching and this value will be overwritten when
				// results come back
				cursor = f.PageCursor

				lastSetFull bool
			)

			for refetch := 0; refetch < MaxRefetches; refetch++ {
				if fetched, err = fetchPage(cursor, limit); err != nil {
					return err
				}

				// if limit is not set or we've already collected enough items
				// we can break the loop right away
				if limit == 0 || fetched == 0 || fetched < limit {
					break
				}

				if uint(len(set)) >= f.Limit {
					// we should return as much as requested
					set = set[0:f.Limit]
					lastSetFull = true
					break
				}

				// In case limit is set very low and we've missed records in the first fetch,
				// make sure next fetch limit is a bit higher
				if limit < MinRefetchLimit {
					limit = MinRefetchLimit
				}

				// @todo it might be good to implement different kind of strategies
				//       (beyond min-refetch-limit above) that can adjust limit on
				//       retry to more optimal number
			}

			if reverseCursor {
				// Cursor for previous page was used
				// Fetched set needs to be reverseCursor because we've forced a descending order to
				// get the previus page
				for i, j := 0, len(set)-1; i < j; i, j = i+1, j-1 {
					set[i], set[j] = set[j], set[i]
				}
			}

			if f.Limit > 0 && len(set) > 0 {
				if f.PageCursor != nil && (!f.PageCursor.Reverse || lastSetFull) {
					f.PrevPage = s.collect{{ export $.Types.Singular }}CursorValues(set[0], sort.Columns()...)
					f.PrevPage.Reverse = true
				}

				// Less items fetched then requested by page-limit
				// not very likely there's another page
				f.NextPage = s.collect{{ export $.Types.Singular }}CursorValues(set[len(set)-1], sort.Columns()...)
			}

			f.PageCursor = nil
			return nil
		}
	{{ end -}}
	)

	return set, f, s.config.ErrorHandler(fetch())
}
{{ end }}

{{- range $lookup := $.Lookups }}
// {{ toggleExport $lookup.Export "Lookup" $.Types.Singular "By" $lookup.Suffix }} {{ comment $lookup.Description true -}}
func (s Store) {{ toggleExport $lookup.Export "Lookup" $.Types.Singular "By" $lookup.Suffix }}(ctx context.Context{{ template "extraArgsDef" $ }}{{- range $field := $lookup.Fields }}, {{ cc2underscore $field }} {{ ($field | $.Fields.Find).Type  }}{{- end }}) (*{{ $.Types.GoType }}, error) {
	return s.execLookup{{ $.Types.Singular }}(ctx{{ template "extraArgsCall" $ }}, squirrel.Eq{
    {{- range $field := $lookup.Fields }}
		s.preprocessColumn({{ printf "%q" ($field | $.Fields.Find).AliasedColumn }}, {{ printf "%q" ($field | $.Fields.Find).LookupFilterPreprocess }}): s.preprocessValue({{ cc2underscore $field }}, {{ printf "%q" ($field | $.Fields.Find).LookupFilterPreprocess }}),
    {{- end }}

    {{ range $field, $value := $lookup.Filter }}
       "{{ ($field | $.Fields.Find).AliasedColumn }}": {{ $value }},
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

		// err = s.{{ unexport $.Types.Singular }}Hook(ctx, TriggerBefore{{ export $.Types.Singular }}Create{{ template "extraArgsCall" . }}, res)
		// if err != nil {
		// 	return err
		// }

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

		// err = s.{{ unexport $.Types.Singular }}Hook(ctx, TriggerBefore{{ export $.Types.Singular }}Update{{ template "extraArgsCall" . }}, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execUpdate{{ export $.Types.Plural }}(
			ctx,
			{{ template "filterByPrimaryKeys" $.Fields.PrimaryKeyFields }},
			s.internal{{ export $.Types.Singular }}Encoder(res).Skip(
				{{- range $field := $.Fields.PrimaryKeyFields -}}
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

		// err = s.{{ unexport $.Types.Singular }}Hook(ctx, TriggerBefore{{ export $.Types.Singular }}Upsert{{ template "extraArgsCall" . }}, res)
		// if err != nil {
		// 	return err
		// }

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
		// err = s.{{ unexport $.Types.Singular }}Hook(ctx, TriggerBefore{{ export $.Types.Singular }}Delete{{ template "extraArgsCall" . }}, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execDelete{{ export $.Types.Plural }}(ctx,{{ template "filterByPrimaryKeys" $.Fields.PrimaryKeyFields }})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// {{ toggleExport .Delete.Export "Delete" $.Types.Singular "By" }}{{ template "primaryKeySuffix" $.Fields }} Deletes row from the {{ $.RDBMS.Table }} table
func (s Store) {{ toggleExport .Delete.Export "Delete" $.Types.Singular "By" }}{{ template "primaryKeySuffix" $.Fields }}(ctx context.Context{{ template "extraArgsDef" . }}{{ template "primaryKeyArgsDef" $.Fields }}) error {
	return s.execDelete{{ export $.Types.Plural }}(ctx, {{ template "filterByPrimaryKeysWithArgs" $.Fields.PrimaryKeyFields }})
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
{{ range $.Fields }}
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

func (s Store) internal{{ $.Types.Singular }}RowScanner({{ template "extraArgsDefFirst" . }}row rowScanner) (res *{{ $.Types.GoType }}, err error) {
	res = &{{ $.Types.GoType }}{}

	if _, has := s.config.RowScanners[{{ printf "%q" (unexport $.Types.Singular) }}]; has {
		scanner := s.config.RowScanners[{{ printf "%q" (unexport $.Types.Singular) }}].(func({{ template "extraArgsDefFirst" . }}_ rowScanner, _ *{{ $.Types.GoType }}) error)
		err = scanner({{ template "extraArgsCallFirst" . }}row, res)
	} else {
	{{- if .RDBMS.CustomRowScanner }}
		err = s.scan{{ $.Types.Singular }}Row({{ template "extraArgsCallFirst" . }}row, res)
	{{- else }}
		err = row.Scan(
		{{- range $.Fields }}
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
	{{- range $.Fields }}
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
	{{ range $.Fields }}
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
    {{- range $.Fields }}
		"{{ .Column }}": res.{{ .Field }},
    {{- end }}
	}
{{- end }}
}

{{ if $.Search.EnablePaging }}
func (s Store) collect{{ export $.Types.Singular }}CursorValues(res *{{ $.Types.GoType }}, cc ...string) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{}

		hasUnique bool

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
			{{- range $.Fields }}
		        {{- if or .IsSortable .IsUnique .IsPrimaryKey -}}
				case "{{ .Column }}":
					cursor.Set(c, res.{{ .Field }}, false)
					{{ if .IsUnique -}}
					hasUnique = true
					{{ end }}
				{{- end }}
			{{- end }}
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique {
		collect(
		{{ range $.Fields -}}
			{{ if .IsPrimaryKey }}"{{ .Column }}",{{ end }}
		{{- end }}
		)
	}

	return cursor
}
{{ end }}

func (s *Store) check{{ export $.Types.Singular }}Constraints(ctx context.Context{{ template "extraArgsDef" $ }}, res *{{ $.Types.GoType }}) error {
{{- range $lookup := $.Lookups }}
	{{ if $lookup.UniqueConstraintCheck }}
	{
		ex, err := s.{{ toggleExport $lookup.Export "Lookup" $.Types.Singular "By" $lookup.Suffix }}(ctx{{ template "extraArgsCall" $ }}{{- range $field := $lookup.Fields }}, res.{{ $field }} {{- end }})
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

// func (s *Store) {{ unexport $.Types.Singular }}Hook(ctx context.Context, key triggerKey{{ template "extraArgsDef" . }}, res *{{ $.Types.GoType }}) error {
// 	if fn, has := s.config.TriggerHandlers[key]; has {
// 		return fn.(func (ctx context.Context, s *Store{{ template "extraArgsDef" . }}, res *{{ $.Types.GoType }}) error)(ctx, s{{ template "extraArgsCall" . }}, res)
// 	}
//
// 	return nil
// }


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
