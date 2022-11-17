package tests

{{ template "gocode/header-gentext.tpl" }}

import (
	"testing"

	"github.com/cortezaproject/corteza/server/store"
)

func testAllGenerated(t *testing.T, s store.Storer) {
{{ range .types }}
	t.Run({{ printf "%q" .ident }}, func(t *testing.T) {
		test{{ .expIdentPlural }}(t, s)
	})
{{- end }}
}
