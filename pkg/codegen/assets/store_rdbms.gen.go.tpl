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
{{- if .RDBMS.CustomFilterConverter }}
	q, err := s.convert{{ pubIdent $.Types.Singular }}Filter(f)
	if err != nil {
	    return nil, f, err
	}
{{- else }}
	q := s.Query{{ pubIdent $.Types.Plural }}()
{{- end }}

    {{ if $.Search.DisablePaging }}
    scap := DefaultSliceCapacity
	{{ else }}
	q = ApplyPaging(q, f.PageFilter)

	scap := f.PerPage
	if scap == 0 {
		scap = DefaultSliceCapacity
	}
	{{ end }}

	var (
		set = make([]*{{ $.Types.GoType }}, 0, scap)
		res *{{ $.Types.GoType }}
	)

	return set, f, func() error {
    {{- if not $.Search.DisablePaging }}
		if f.Count, err = Count(ctx, s.db, q); err != nil || f.Count == 0 {
			return err
		}
	{{- end }}
		rows, err := s.Query(ctx, q)
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
	}()
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
				return err
			}
		}

		return nil
	})
}

// Update{{ pubIdent $.Types.Singular }} updates one or more existing rows in {{ $.RDBMS.Table }}
func (s Store) Update{{ pubIdent $.Types.Singular }}(ctx context.Context, rr ... *{{ $.Types.GoType }}) error {
	return s.PartialUpdate{{ pubIdent $.Types.Singular }}(ctx, nil, rr...)
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
				return err
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
				return err
			}
		}

		return nil
	})
}


// Remove{{ pubIdent $.Types.Singular }}By{{ template "primaryKeySuffix" $.Fields }} removes row from the {{ $.RDBMS.Table }} table
func (s Store) Remove{{ pubIdent $.Types.Singular }}By{{ template "primaryKeySuffix" $.Fields }}(ctx context.Context {{ template "primaryKeyArgs" $.Fields }}) error {
	return ExecuteSqlizer(ctx, s.DB(), 	s.Delete(s.{{ $.Types.Singular }}Table({{ printf "%q" .RDBMS.Alias }})).Where({{ template "filterByPrimaryKeysWithArgs" $.Fields }},))
}


// Truncate{{ pubIdent $.Types.Plural }} removes all rows from the {{ $.RDBMS.Table }} table
func (s Store) Truncate{{ pubIdent $.Types.Plural }}(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.{{ $.Types.Singular }}Table())
}


// ExecUpdate{{ pubIdent $.Types.Plural }} updates all matched (by cnd) rows in {{ $.RDBMS.Table }} with given data
func (s Store) ExecUpdate{{ pubIdent $.Types.Plural }}(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), 	s.Update(s.{{ $.Types.Singular }}Table({{ printf "%q" .RDBMS.Alias }})).Where(cnd).SetMap(set))
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

