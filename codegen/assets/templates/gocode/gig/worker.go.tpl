package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
    "context"
    "fmt"

{{- range .imports }}
    {{ . }}
{{- end }}
)

const (
    {{- range .workers }}
    {{ .goConst }} = "{{ .ident }}"
    {{- end }}
)

{{ range .workers }}
func (w *{{ .goType }}) Preprocess(ctx context.Context, tasks ...Preprocessor) (err error) {
    for _, t := range tasks {
        switch tc := t.(type) {
{{- range .tasks }}
        case {{ .goType }}:
            err = w.{{ .ident }}(ctx, tc)
{{- end }}
        default:
            err = fmt.Errorf("unknown preprocessor: %s", w.Ref())
        }

        if err != nil {
            return
        }
    }

    return nil
}
{{ end }}

// ------------------------------------------------------------------------
// Worker registry

func workerDefinitions() WorkerDefSet {
  return WorkerDefSet{
  {{- range .workers }}
    {
      Ref: {{ .goConst }},
    {{- if .description }}
      Description: "{{ .description }}",
    {{- end }}
    },
  {{- end }}
  }
}
