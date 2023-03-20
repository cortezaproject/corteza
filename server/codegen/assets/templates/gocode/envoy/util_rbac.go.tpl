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

  // @todo this is to support res. tr. resources also.
	//       Split it into a separate function and remove this.
	if !strings.HasPrefix(rt, "corteza::") {
		rt = "corteza::" + rt
	}

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

    case "{{.fqrt}}":
      scope := Scope{}
    {{$res := .}}
    {{range $i, $p := .parents}}
      if gRef(pp, {{$i}}) == "" {
        return
      }

      {{- range $cmp := $rootCmp.resources }}
        {{ if $cmp.envoy.omit }}
          {{continue}}
        {{ end }}

        {{ if eq $p.handle $cmp.ident }}
          {{ if and (eq $rootCmp.ident "compose") (eq $i 0) }}
          aux := gRef(pp, {{ $i }})
          if aux != "" {
            scope.ResourceType = "corteza::compose:namespace"
            scope.Identifiers = MakeIdentifiers(aux)
          }
          {{ end }}
          out["Path.{{$i}}"] = Ref{
            ResourceType: "{{$cmp.fqrt}}",
            Identifiers:  MakeIdentifiers(gRef(pp, {{ $i }})),
            Scope: scope,
          }
          {{break}}
        {{ end }}
      {{- end }}
    {{end}}

    if gRef(pp, {{len .parents}}) == "" {
      return
    }
    {{if eq .ident "namespace"}}
    scope.ResourceType = "{{.fqrt}}"
    scope.Identifiers = MakeIdentifiers(gRef(pp, {{len .parents}}))
    {{end}}
    out["Path.{{len .parents}}"] = Ref{
      ResourceType: "{{.fqrt}}",
      Identifiers:  MakeIdentifiers(gRef(pp, {{len .parents}})),
      Scope: scope,
    }

  {{ end }}
{{ end }}
	}

	return
}
