package cache

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_cache_struct.gen.go.tpl
// Definitions:
{{ range . }}
{{- if .Exported -}}
//  - {{ .Source }}
{{ end -}}{{- end }}
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"github.com/cortezaproject/corteza-server/store"
	"github.com/dgraph-io/ristretto"
)

type (
	Cache struct {
		store.Storer

	{{ range . -}}
	{{ if .Cache.Enable -}}
		{{ unexport .Types.Plural }} *ristretto.Cache
	{{ end -}}
	{{ end -}}
	}
)

var _ store.Users = &Cache{}


func Connect(s store.Storer) (store.Storer, error) {
	var (
		err error
		c = &Cache{Storer: s}
	)


	{{ range . -}}
	{{ if .Cache.Enable -}}
	c.{{ unexport .Types.Plural }}, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 20, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {	return nil, err }

	{{ end -}}
	{{ end -}}

	return c, nil
}
