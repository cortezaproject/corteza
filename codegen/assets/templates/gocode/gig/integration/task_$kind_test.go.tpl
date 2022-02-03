package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
  "testing"

	"github.com/cortezaproject/corteza-server/pkg/gig"

{{- range .imports }}
    {{ . }}
{{- end }}
)


func {{ .expTest }}(t *testing.T) {
	var (
		ctx, svc, h, s = setup(t)
		err            error
    g gig.Gig
	)
	_ = s
	_ = svc
	_ = err

{{ $test := .test }}
{{ $testWorker := .testWorker }}

{{- range .tasks }}
	t.Run("{{ .ident }}", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: {{ $testWorker }}(t, h, s),
		})
		h.a.NoError(err)

		{{ $test }}_{{ .ident }}(ctx, t, h, svc, s, g, "{{ .ident }}")
	})
{{- end }}
}
