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

{{- range .workers }}
	{{ $worker := .ident }}

	{{- range .tasks }}
	t.Run("{{ .ident }}", func(_ *testing.T) {
		g, err = svc.Create(ctx, gig.UpdatePayload{
			Worker: {{ $testWorker }}_{{ $worker }}(t, h, s),
		})
		h.a.NoError(err)

		{{ $test }}_{{$worker}}_{{ .ident }}(ctx, t, h, svc, s, g, "{{ $worker }}", "{{ .ident }}")
	})
	{{- end }}
{{- end }}

}
