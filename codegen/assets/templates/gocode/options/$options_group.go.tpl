package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

{{ if .imports }}
import (
{{- range .imports }}
    {{ . }}
{{- end }}
)
{{- end }}

type (
    {{ .struct }} struct {
        {{- range .options }}
          {{ .expIdent }} {{ .type }} `env:"{{ .env }}"`
        {{- end }}
    }
)

// {{ .func }} initializes and returns a {{ .struct }} with default values
//
// This function is auto-generated
func {{ .func }}() (o *{{ .struct }}) {
    o = &{{ .struct }}{
      {{- range  .options }}
        {{- if .default }}
          {{ .expIdent }}: {{ .default }},
        {{- end }}
      {{- end }}
    }

    fill(o)

    // Function that allows access to custom logic inside the parent function.
    // The custom logic in the other file should be like:
    // func (o *{{ .struct }}) Defaults() {...}
    func(o interface{}) {
      if def, ok := o.(interface{ Defaults() }); ok {
        def.Defaults()
      }
    }(o)

    return
}
