package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
  "testing"
	"github.com/stretchr/testify/require"

{{- range .imports }}
    {{ . }}
{{- end }}
)

type (
{{- range .tasks }}
	constructor{{ .constructor }} func(map[string]interface{}) ({{ .goType }}, error)
{{- end }}
)

func {{ .expTest }}(t *testing.T) {
{{ $test := .test }}
{{ $testWorker := .testWorker }}

{{- range .tasks }}
	{{- $hasRequired := false -}}
	{{- range .struct -}}
		{{ $hasRequired = or $hasRequired .required }}
	{{- end }}
	{{- if $hasRequired }}
	t.Run("{{ .ident }} constructor required missing", func(t *testing.T) {
		_, err := {{ .constructorParams }}(map[string]interface{}{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "required")
	})
	{{- end }}
	t.Run("{{ .ident }} constructor unknown", func(t *testing.T) {
		_, err := {{ .constructorParams }}(map[string]interface{}{
			"does not exist": "true",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown parameter")
	})
	t.Run("{{ .ident }} constructor", func(t *testing.T) {
		// func {{ $test }}_{{ .ident }}_constructor(t *testing.T, c constructor{{ .constructor }})
		{{ $test }}_{{ .ident }}_constructor(t, {{ .constructorParams }})
	})
{{ end -}}
}
