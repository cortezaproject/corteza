package tests

// This file is auto-generated.
//
// Template:	pkg/codegen/assets/store_test_all.gen.go
// Definitions:
{{ range . }}
{{- if .Exported -}}
//  - {{ .Source }}
{{ end -}}{{- end }}
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza-server/store"
	"testing"
)

func testAllGenerated(t *testing.T, s store.Storable) {
{{- range . }}
	{{- if .Exported }}
		// Run generated tests for {{ .Types.Base }}
		t.Run({{ printf "%q" .Types.Base }}, func(t *testing.T) {
			test{{ export .Types.Base }}(t, s)
		})
	{{ end -}}
{{ end -}}
}
