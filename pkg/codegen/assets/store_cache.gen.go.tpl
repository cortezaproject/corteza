package cache

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_cache.gen.go.tpl
// Definitions: {{ .Source }}
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"errors"
	"time"
	"github.com/cortezaproject/corteza-server/store"
{{- range .Import }}
    {{ normalizeImport . }}
{{- end }}
)

{{- $Types := .Types }}
{{- $Fields := .Fields }}

var _ = errors.Is

func (c Cache) cache{{ export $.Types.Singular }}(res *{{ $Types.GoType }}) {
	var (
	ttl  time.Duration = 0
	cost int64 = 1
	)

	if c.{{ unexport .Types.Plural }}.SetWithTTL(res.ID, res, cost, ttl) {
		for _, ikey := range c.{{ unexport $.Types.Singular }}Indexes(res) {
			c.{{ unexport .Types.Plural }}.SetWithTTL(ikey, res.ID, cost, ttl)
		}
	}
}

func (c Cache) getCached{{ export $.Types.Singular }}ByKey(ikey string) (interface{}, bool) {
	if val, found := c.{{ unexport .Types.Plural }}.Get(ikey); found {
		if id, ok := val.(uint64); ok {
			return c.{{ unexport .Types.Plural }}.Get(id)
		}

		c.{{ unexport .Types.Plural }}.Del(val)
	}

	return nil, false
}

func (c Cache) {{ unexport $.Types.Singular }}Indexes(res *{{ $Types.GoType }}) []string {
	return []string{
		{{ range $.Lookups -}}
		{{ if not .IsPrimary -}}
		iKey({{ printf "%q" .Suffix }},
			{{- range .Fields }}
			store.PreprocessValue(res.{{ .Field }}, {{ printf "%q" .LookupFilterPreprocess }}), {{- end }}
		),
		{{ end -}}
		{{ end -}}
	}
}

{{- range $.Lookups }}
// Lookup{{ export $.Types.Singular }}By{{ export .Suffix }} {{ comment .Description true -}}
func (c Cache) Lookup{{ export $Types.Singular }}By{{ export .Suffix }}(ctx context.Context{{ template "extraArgsDef" $ }}{{- range .Fields }}, {{ cc2underscore .Field }} {{ .Type }}{{- end }}) (*{{ $Types.GoType }}, error) {
	{{ if .IsPrimary }}
	if val, found := c.{{ unexport $Types.Plural }}.Get(id); found {
		if res, ok := val.(*{{ $Types.GoType }}); ok {
			return res, nil
		}

		c.{{ unexport $Types.Plural }}.Del(id)
	}

	{{ else }}
	key := iKey(
		{{ printf "%q" $.Types.Plural }},
		{{ printf "%q" .Suffix }},
		{{- range .Fields }}
		store.PreprocessValue({{ cc2underscore .Field }}, {{ printf "%q" .LookupFilterPreprocess }}), {{- end }}
	)

	if val, found := c.getCached{{ export $.Types.Singular }}ByKey(key); found {
		if res, ok := val.(*{{ $Types.GoType }}); ok {
			return res, nil
		}

		c.{{ unexport $Types.Plural }}.Del(key)
	}
	{{ end }}

	if res, err := c.Storer.Lookup{{ export $Types.Singular }}By{{ export .Suffix }}(ctx {{- range .Fields }}, {{ cc2underscore .Field }}{{- end }}); err != nil {
		return nil, err
	} else {
		c.cache{{ export $.Types.Singular }}(res)
		return res, nil
	}
}
{{ end }}

{{ if .Create.Enable }}
// {{ toggleExport .Create.Export "Create" $.Types.Singular }} updates cache and forwards call to next configured store
func (c Cache) {{ toggleExport .Create.Export "Create" $.Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $.Types.GoType }}) (err error) {
	for _, res := range rr {
		if err = c.Storer.{{ toggleExport .Create.Export "Create" $.Types.Singular }}(ctx, res); err != nil {
			return err
		}

		c.cache{{ export $.Types.Singular }}(res)
	}

	return nil
}
{{ end }}

{{ if .Update.Enable }}
// {{ toggleExport .Update.Export "Update" $.Types.Singular }} updates cache and forwards call to next configured store
func (c Cache) {{ toggleExport .Update.Export "Update" $.Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $.Types.GoType }}) error {
	for _, res := range rr {
		if err := c.Storer.{{ toggleExport .Update.Export "Update" $.Types.Singular }}(ctx, res); err != nil {
			return err
		}

		c.cache{{ export $.Types.Singular }}(res)
	}

	return nil
}
{{ end }}


{{ if .Upsert.Enable }}
// {{ toggleExport .Upsert.Export "Upsert" $.Types.Singular }} updates cache and forwards call to next configured store
func (c Cache) {{ toggleExport .Upsert.Export "Upsert" $.Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $.Types.GoType }}) (err error) {
	for _, res := range rr {
		if err = c.Storer.{{ toggleExport .Upsert.Export "Upsert" $.Types.Singular }}(ctx, res); err != nil {
			return err
		}

		c.cache{{ export $.Types.Singular }}(res)
	}

	return nil
}
{{ end }}

{{ if .Delete.Enable }}
// {{ toggleExport .Delete.Export "Delete" $.Types.Singular }} Deletes one or more rows from {{ $.RDBMS.Table }} table
func (c Cache) {{ toggleExport .Delete.Export "Delete" $.Types.Singular }}(ctx context.Context{{ template "extraArgsDef" . }}, rr ... *{{ $.Types.GoType }}) (err error) {
	for _, res := range rr {
		if err = c.Storer.{{ toggleExport .Delete.Export "Delete" $.Types.Singular }}(ctx, res); err != nil {
			return
		}

		c.{{ unexport $Types.Plural }}.Del(res.ID)
		for _, key := range c.{{ unexport $.Types.Singular }}Indexes(res) {
			c.{{ unexport $Types.Plural }}.Del(key)
		}
	}

	return nil
}

// {{ toggleExport .Delete.Export "Delete" $.Types.Singular "By" }}{{ template "primaryKeySuffix" $.RDBMS.Columns }} Deletes row from the {{ $.RDBMS.Table }} table
func (c Cache) {{ toggleExport .Delete.Export "Delete" $.Types.Singular "By" }}{{ template "primaryKeySuffix" $.RDBMS.Columns }}(ctx context.Context{{ template "extraArgsDef" . }}{{ template "primaryKeyArgsDef" $.Fields }}) error {
	if err := c.Storer.{{ toggleExport .Delete.Export "Delete" $.Types.Singular "By" }}{{ template "primaryKeySuffix" $.RDBMS.Columns }}(ctx{{ template "extraArgsDef" . }}{{ template "primaryKeyArgsCall" $.RDBMS.Columns.PrimaryKeyFields }}); err != nil {
		return err
	}

	c.{{ unexport $Types.Plural }}.Del(ID)
	return nil
}
{{ end }}

// {{ toggleExport .Truncate.Export "Truncate" $.Types.Plural }} Deletes all rows from the {{ $.RDBMS.Table }} table
func (c Cache) {{ toggleExport .Truncate.Export "Truncate" $.Types.Plural }}(ctx context.Context, {{ template "extraArgsDef" . }}) error {
	if err := c.Storer.{{ toggleExport .Truncate.Export "Truncate" $.Types.Plural }}(ctx); err != nil {
		return err

	}

	c.{{ unexport $Types.Plural }}.Clear()
	return nil
}
