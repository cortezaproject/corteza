package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"strings"

{{- range .imports }}
    "{{ . }}"
{{- end }}
)

// SplitResourceIdentifier takes an identifier string and splices it into path
// identifiers as defined by the resource
func SplitResourceIdentifier(ref string) (out map[string]Ref) {
	out = make(map[string]Ref, 3)

	ref = strings.TrimRight(ref, "/")
	pp := strings.Split(ref, "/")
	rt := pp[0]
	pp = pp[1:]

	gRef := func(pp []string, i int) string {
		if pp[i] == "*" {
			return ""
		}
		return pp[i]
	}

	switch rt {

{{ range .components }}
  {{ $rootCmp := . }}
  {{range .resources}}
  	{{ $a := . }}


    case "corteza::{{$rootCmp.ident}}:{{.ident}}":
    {{$res := .}}
    {{range $i, $p := .parents}}
      if gRef(pp, {{$i}}) == "" {
        return
      }

      {{- range $cmp := $rootCmp.resources }}
        {{ if or ($cmp.envoy.omit) (not $cmp.envoy.use) }}
          {{continue}}
        {{ end }}

        {{ if eq $p.handle $cmp.ident }}
          out["{{$i}}"] = Ref{
            ResourceType: "{{$cmp.fqrt}}",
            Identifiers:  MakeIdentifiers(gRef(pp, {{ $i }})),
          }
          {{break}}
        {{ end }}
      {{- end }}
    {{end}}

    if gRef(pp, {{len .parents}}) == "" {
      return
    }
    out["{{len .parents}}"] = Ref{
      ResourceType: "{{.fqrt}}",
      Identifiers:  MakeIdentifiers(gRef(pp, {{len .parents}})),
    }

  {{ end }}
{{ end }}
	}

	return
}
