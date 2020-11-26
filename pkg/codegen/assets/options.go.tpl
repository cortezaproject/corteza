package {{ .Package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}

{{ if $.Imports -}}
import (
{{- range .Imports }}
    {{ normalizeImport . }}
{{- end }}
){{ end }}

type (
    {{ export $.Name }}Opt struct {
        {{- range $prop := $.Properties}}
          {{ export $prop.Name }} {{ $prop.Type }} `env:"{{ toUpper $prop.Env}}"`
        {{- end }}
    }
)

// {{ export $.Name }} initializes and returns a {{ export $.Name }}Opt with default values
func {{ export $.Name }}() (o *{{ export $.Name }}Opt) {
    o = &{{ export $.Name }}Opt{
      {{- range  $prop := $.Properties }}
        {{- if $prop.Default }}
          {{ export $prop.Name }}: {{ $prop.Default }}, 
        {{- end }}
      {{- end }}
    }

    fill(o)

    // Function that allows access to custom logic inside the parent function.
    // The custom logic in the other file should be like:
    // func (o *{{ export $.Name}}) Defaults() {...}
    func(o interface{}) {
      if def, ok := o.(interface{ Defaults() }); ok {
        def.Defaults()
      }
    }(o)

    return
}
