package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
  "fmt"
  "github.com/cortezaproject/corteza-server/pkg/gig"

{{- range .imports }}
    {{ . }}
{{- end }}
)


{{- range .tasks }}

func (conv Gig) {{.unwrapSetFunc}}(wraps ParamWrapSet) (out gig.{{.expKind}}Set, err error) {
  for _, w := range wraps {
		var aux gig.{{.expKind}}
		aux, err = conv.{{.unwrapFunc}}(w)
		if err != nil {
			return
		}

		out = append(out, aux)
	}
	return
}

func (conv Gig) {{.unwrapFunc}}(w ParamWrap) (out gig.{{.expKind}}, err error) {
  switch w.Ref {
  {{- range .set }}
  case gig.{{.goConst}}:
    return gig.{{.expKind}}{{.expIdent}}Params(w.Params)
  {{- end }}
	}

	return nil, fmt.Errorf("unknown {{.kind}}: %s", w.Ref)
}

func (conv Gig) {{.wrapSetFunc}}(tt gig.{{.expKind}}Set) (out ParamWrapSet) {
  for _, t := range tt {
    out = append(out, conv.{{.wrapFunc}}(t))
  }

	return
}

func (conv Gig) {{.wrapFunc}}(t gig.{{.expKind}}) (out ParamWrap) {
  out.Ref = t.Ref()
	out.Params = t.Params()

  return
}


{{- end }}
