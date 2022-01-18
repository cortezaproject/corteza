package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
  "github.com/spf13/cast"
	"fmt"
{{- range .imports }}
    {{ . }}
{{- end }}
)

type (
  {{- range .tasks }}
  {{.goType}} struct {
    {{- range .struct }}
    {{- .ident }} {{ .goType }}
    {{ end -}}
  }

  {{- end }}
)

const (
  {{- range .tasks }}
  {{ .goConst }} = "{{ .ident }}"
  {{- end }}
)

{{ $taskConst := .taskConst }}

// ------------------------------------------------------------------------
// Constructors and utils

{{- range .tasks -}}

{{ $taskGoType := .goType }}
{{ $taskGoInterface := .goInterface }}
{{ $taskTransformer := .transformer }}
{{ $taskIdent := .ident }}

{{ $hasRequired := false }}

// {{ .constructorParams }} returns a new {{ $taskGoType }} from the params
func {{ .constructorParams }}(params map[string]interface{}) ({{ $taskGoInterface }}, error) {
  var (
    out = {{ .goType }}{}
    err error
  )

  // Param validation
  // - supported params
  index := map[string]bool{
  {{- range .struct }}
    "{{- .ident }}": true,
  {{- end }}
  }
  for p := range params {
    if !index[p] {
      return nil, fmt.Errorf("unknown parameter provided to {{ $taskIdent }}: %s", p)
    }
  }

  // Fill and check requirements
  {{- range .struct }}
    {{- if .required }}
  if _, ok := params["{{ .ident }}"]; !ok {
    return nil, fmt.Errorf("required parameter not provided: {{ .ident }}")
  }
    {{- end }}
  out.{{ .ident }} = {{ .castFunc }}(params["{{ .ident }}"])
    {{- $hasRequired = or $hasRequired .required -}}
  {{- end }}

  {{- if $taskTransformer }}
  out, err = {{ $taskTransformer }}(out)
  {{- end }}
  return out, err
}

{{ $struct := .struct }}

  {{- range $struct }}
    {{- if not .required }}
// {{ .constructor }} returns a new {{ $taskGoType }} from the required fields and {{ .ident }}
func {{ .constructor }}(
  {{- range $struct -}}
    {{- if .required }}
    {{- .ident }} {{ .goType -}},
    {{- end }}
  {{- end -}}
  {{ .ident }} {{ .goType -}}) ({{ $taskGoInterface }}, error) {
  var (
    err error
    out {{ $taskGoType }}
  )
  out = {{ $taskGoType }}{
  {{- range $struct }}
    {{- if .required }}
    {{ .ident }}: {{ .ident }},
    {{- end }}
  {{- end }}
  {{ .ident }}: {{ .ident }},
  }

  {{- if $taskTransformer }}
  out, err = {{ $taskTransformer }}(out)
  {{- end }}

  return out, err
}

    {{- end }}
  {{- end }}

  {{- if $hasRequired }}
  func {{ .constructor }}(
    {{- range .struct -}}
      {{- if .required }}
      {{- .ident }} {{ .goType -}},
      {{- end }}
    {{- end -}}
  ) ({{ $taskGoInterface }}, error) {
  var (
    err error
    out {{ $taskGoType }}
  )

  out = {{ $taskGoType }}{
  {{- range $struct }}
    {{- if .required }}
    {{ .ident }}: {{ .ident }},
    {{- end }}
  {{- end }}
  }

  {{- if $taskTransformer }}
  out, err = {{ $taskTransformer }}(out)
  {{- end }}

  return out, err
  }
  {{- end }}


func (t {{ $taskGoType }}) Ref() string {
	return {{ .goConst }}
}

func (t {{ $taskGoType }}) Params() map[string]interface{} {
{{- $length := len .struct }} {{- if eq $length 0 }}
	return nil
{{- else }}
  return map[string]interface{}{
  {{- range .struct }}
    {{- if .param }}
    "{{ .ident }}": t.{{ .ident }},
    {{- end }}
  {{- end }}
  }
{{- end }}
}

{{- end}}


// ------------------------------------------------------------------------
// Task registry

func {{.taskKind}}Definitions() TaskDefSet {
  return TaskDefSet{
  
  {{- range .tasks }}
    {
      Ref: {{ .goConst }},
      Kind: {{ $taskConst }},
    {{- if .description }}
      Description: "{{ .description }}",
    {{- end }}

    {{- if .struct }}
      Params: []taskDefParam{
      {{- range .struct }}
        {{- if .param }}
        {
          Name: "{{ .ident }}",
          Kind: "{{ .exprType }}",
          Required: {{ .required }},

        },
        {{- end }}
      {{- end }}
      },

    {{- end }}
    },
  {{- end }}
  
  }
}
