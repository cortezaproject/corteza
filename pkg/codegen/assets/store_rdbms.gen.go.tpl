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
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
{{- if not $.Search.DisablePaging }}
	"strings"
{{- end }}
{{- range $import := $.Import }}
    {{ normalizeImport $import }}
{{- end }}
)

{{ if not $.Search.Disable }}
// Search{{ pubIdent $.Types.Plural }} returns all matching rows
//
// This function calls convert{{ pubIdent $.Types.Singular }}Filter with the given
// {{ $.Types.GoFilterType }} and expects to receive a working squirrel.SelectBuilder
func (s Store) Search{{ pubIdent $.Types.Plural }}(ctx context.Context, f {{ $.Types.GoFilterType }}) ({{ $.Types.GoSetType }}, {{ $.Types.GoFilterType }}, error) {
	var scap uint

{{- if .RDBMS.CustomFilterConverter }}
	q, err := s.convert{{ pubIdent $.Types.Singular }}Filter(f)
	if err != nil {
	    return nil, f, err
	}
{{- else }}
	q := s.Query{{ pubIdent $.Types.Plural }}()
{{- end }}

{{ if not $.Search.DisablePaging }}
	scap = f.Limit

	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reverseCursor := f.PageCursor != nil && f.PageCursor.Reverse
{{ end }}

{{ if $.Search.DisableSorting }}
	{{ if not $.Search.DisablePaging }}
	// Sorting is disabled in definition yaml file
	// {search: {disableSorting:true}}
	//
	// We still need to sort the results by primary key for paging purposes
	sort := store.SortExprSet{
	{{ range $.Fields }}
		{{- if or .IsPrimaryKey -}}
			&store.SortExpr{Column: {{ printf "%q" .Column  }}, {{ if .SortDescending }}Descending: true, {{ end }}},
		{{- end }}
	{{- end }}
	}
	{{ end }}
{{ else }}
	if err = f.Sort.Validate(s.sortable{{ pubIdent $.Types.Singular }}Columns()...); err != nil {
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
{{ end }}

	if scap == 0 {
		scap = DefaultSliceCapacity
	}


	var (
		set = make([]*{{ $.Types.GoType }}, 0, scap)

	{{- if $.Search.DisablePaging }}
		// Paging is disabled in definition yaml file
		// {search: {disablePaging:true}} and this allows
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
				if res, err = s.internal{{ pubIdent $.Types.Singular }}RowScanner(rows, rows.Err()); err != nil {
					if cerr := rows.Close(); cerr != nil {
						return fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return err
				}

				set = append(set, res)
			}

			return rows.Close()
		}
	{{ else }}
		// fetches rows and scans them into {{ pubIdent $.Types.GoType }} resource this is then passed to Check function on filter
		// to help determine if fetched resource fits or not
		//
		// Note that limit is passed explicitly and is not necessarily equal to filter's limit. We want
		// to keep that value intact.
		//
		// The value for cursor is used and set directly from/to the filter!
		//
		// It returns total number of fetched pages and modifies PageCursor value for paging
		fetchPage = func(cursor *store.PagingCursor, limit uint) (fetched uint, err error) {
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
				if res, err = s.internal{{ pubIdent $.Types.Singular }}RowScanner(rows, rows.Err()); err != nil {
					if cerr := rows.Close(); cerr != nil {
						err = fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return
				}

				// If check function is set, call it and act accordingly
			{{ if not $.Search.DisableFilterCheckFn }}
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
					f.PrevPage = s.collect{{ pubIdent $.Types.Singular }}CursorValues(set[0], sort.Columns()...)
					f.PrevPage.Reverse = true
				}

				// Less items fetched then requested by page-limit
				// not very likely there's another page
				f.NextPage = s.collect{{ pubIdent $.Types.Singular }}CursorValues(set[len(set)-1], sort.Columns()...)
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
// Lookup{{ pubIdent $.Types.Singular }}By{{ pubIdent $lookup.Suffix }} {{ comment $lookup.Description true -}}
func (s Store) Lookup{{ pubIdent $.Types.Singular }}By{{ pubIdent $lookup.Suffix }}(ctx context.Context{{- range $field := $lookup.Fields }}, {{ cc2underscore $field }} {{ ($field | $.Fields.Find).Type  }}{{- end }}) (*{{ $.Types.GoType }}, error) {
	return s.{{ $.Types.Singular }}Lookup(ctx, squirrel.Eq{
    {{- range $field := $lookup.Fields }}
       "{{ ($field | $.Fields.Find).AliasedColumn }}": {{ cc2underscore $field }},
    {{- end }}
    {{- range $field, $value := $lookup.Filter }}
       "{{ ($field | $.Fields.Find).AliasedColumn }}": {{ $value }},
    {{- end }}
    })
}
{{ end }}

// Create{{ pubIdent $.Types.Singular }} creates one or more rows in {{ $.RDBMS.Table }} table
func (s Store) Create{{ pubIdent $.Types.Singular }}(ctx context.Context, rr ... *{{ $.Types.GoType }}) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.{{ $.Types.Singular }}Table()).SetMap(s.internal{{ pubIdent $.Types.Singular }}Encoder(res)))
			if err != nil {
				return s.config.ErrorHandler(err)
			}
		}

		return nil
	})
}

// Update{{ pubIdent $.Types.Singular }} updates one or more existing rows in {{ $.RDBMS.Table }}
func (s Store) Update{{ pubIdent $.Types.Singular }}(ctx context.Context, rr ... *{{ $.Types.GoType }}) error {
	return s.config.ErrorHandler(s.PartialUpdate{{ pubIdent $.Types.Singular }}(ctx, nil, rr...))
}

// PartialUpdate{{ pubIdent $.Types.Singular }} updates one or more existing rows in {{ $.RDBMS.Table }}
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdate{{ pubIdent $.Types.Singular }}(ctx context.Context, onlyColumns []string, rr ... *{{ $.Types.GoType }}) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdate{{ pubIdent $.Types.Plural }}(
				ctx,
				{{ template "filterByPrimaryKeys" $.Fields }},
				s.internal{{ pubIdent $.Types.Singular }}Encoder(res).Skip(
					{{- range $field := $.Fields -}}
						{{- if $field.IsPrimaryKey -}}
							{{ printf "%q" $field.Column  }},
						{{- end -}}
					{{- end -}}
			).Only(onlyColumns...))
			if err != nil {
				return s.config.ErrorHandler(err)
			}
		}

		return nil
	})
}

// Remove{{ pubIdent $.Types.Singular }} removes one or more rows from {{ $.RDBMS.Table }} table
func (s Store) Remove{{ pubIdent $.Types.Singular }}(ctx context.Context, rr ... *{{ $.Types.GoType }}) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.{{ $.Types.Singular }}Table({{ printf "%q" .RDBMS.Alias }})).Where({{ template "filterByPrimaryKeys" $.Fields }},))
			if err != nil {
				return s.config.ErrorHandler(err)
			}
		}

		return nil
	})
}


// Remove{{ pubIdent $.Types.Singular }}By{{ template "primaryKeySuffix" $.Fields }} removes row from the {{ $.RDBMS.Table }} table
func (s Store) Remove{{ pubIdent $.Types.Singular }}By{{ template "primaryKeySuffix" $.Fields }}(ctx context.Context {{ template "primaryKeyArgs" $.Fields }}) error {
	return s.config.ErrorHandler(ExecuteSqlizer(ctx, s.DB(), 	s.Delete(s.{{ $.Types.Singular }}Table({{ printf "%q" .RDBMS.Alias }})).Where({{ template "filterByPrimaryKeysWithArgs" $.Fields }},)))
}


// Truncate{{ pubIdent $.Types.Plural }} removes all rows from the {{ $.RDBMS.Table }} table
func (s Store) Truncate{{ pubIdent $.Types.Plural }}(ctx context.Context) error {
	return s.config.ErrorHandler(Truncate(ctx, s.DB(), s.{{ $.Types.Singular }}Table()))
}


// ExecUpdate{{ pubIdent $.Types.Plural }} updates all matched (by cnd) rows in {{ $.RDBMS.Table }} with given data
func (s Store) ExecUpdate{{ pubIdent $.Types.Plural }}(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(ExecuteSqlizer(ctx, s.DB(), 	s.Update(s.{{ $.Types.Singular }}Table({{ printf "%q" .RDBMS.Alias }})).Where(cnd).SetMap(set)))
}

// {{ $.Types.Singular }}Lookup prepares {{ $.Types.Singular }} query and executes it,
// returning {{ $.Types.GoType }} (or error)
func (s Store) {{ $.Types.Singular }}Lookup(ctx context.Context, cnd squirrel.Sqlizer) (*{{ $.Types.GoType }}, error) {
	return s.internal{{ $.Types.Singular }}RowScanner(s.QueryRow(ctx, s.Query{{ pubIdent $.Types.Plural }}().Where(cnd)))
}

func (s Store) internal{{ $.Types.Singular }}RowScanner(row rowScanner, err error) (*{{ $.Types.GoType }}, error) {
	if err != nil {
		return nil, err
	}

	var res = &{{ $.Types.GoType }}{}
	if _, has := s.config.RowScanners[{{ printf "%q" (unpubIdent $.Types.Singular) }}]; has {
		scanner := s.config.RowScanners[{{ printf "%q" (unpubIdent $.Types.Singular) }}].(func(rowScanner, *{{ $.Types.GoType }}) error)
		err = scanner(row, res)
	} else {
	{{- if .RDBMS.CustomRowScanner }}
		err = s.scan{{ $.Types.Singular }}Row(row, res)
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

// Query{{ pubIdent $.Types.Plural }} returns squirrel.SelectBuilder with set table and all columns
func (s Store) Query{{ pubIdent $.Types.Plural }}() squirrel.SelectBuilder {
	return s.Select(s.{{ $.Types.Singular }}Table({{ printf "%q" .RDBMS.Alias }}), s.{{ $.Types.Singular }}Columns({{ printf "%q" $.RDBMS.Alias }})...)
}

// {{ $.Types.Singular }}Table name of the db table
func (Store) {{ $.Types.Singular }}Table(aa ... string) string {
		var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "{{ $.RDBMS.Table }}" + alias
}

// {{ $.Types.Singular }}Columns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) {{ $.Types.Singular }}Columns(aa ... string) []string {
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

{{ if not $.Search.DisableSorting }}
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

// internal{{ pubIdent $.Types.Singular }}Encoder encodes fields from {{ $.Types.GoType }} to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encode{{ pubIdent $.Types.Singular }}
// func when rdbms.customEncoder=true
func (s Store) internal{{ pubIdent $.Types.Singular }}Encoder(res *{{ $.Types.GoType }}) store.Payload {
{{- if .RDBMS.CustomEncoder }}
	return s.encode{{ pubIdent $.Types.Singular }}(res)
{{- else }}
	return store.Payload{
    {{- range $.Fields }}
		"{{ .Column }}": res.{{ .Field }},
    {{- end }}
	}
{{- end }}
}

{{ if not $.Search.DisablePaging }}
func (s Store) collect{{ pubIdent $.Types.Singular }}CursorValues(res *{{ $.Types.GoType }}, cc ...string) *store.PagingCursor {
	var (
		cursor = &store.PagingCursor{}

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


{{/* ************************************************************ */}}

{{- define "filterByPrimaryKeys" -}}
    squirrel.Eq{
    {{- range $field := . -}}
        {{- if $field.IsPrimaryKey -}}
            s.preprocessColumn({{ printf "%q" $field.AliasedColumn }}, {{ printf "%q" $field.LookupFilterPreprocess }}): s.preprocessValue(res.{{ $field.Field }}, {{ printf "%q" $field.LookupFilterPreprocess }}),
        {{ end }}
    {{- end -}}
    }
{{- end -}}

{{- define "filterByPrimaryKeysWithArgs" -}}
    squirrel.Eq{
    {{- range $field := . }}
        {{- if $field.IsPrimaryKey -}}
            s.preprocessColumn({{ printf "%q" $field.AliasedColumn }}, {{ printf "%q" $field.LookupFilterPreprocess }}): s.preprocessValue({{ $field.Arg }}, {{ printf "%q" $field.LookupFilterPreprocess }}),
        {{ end }}
    {{ end -}}
    }
{{- end -}}

