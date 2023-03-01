package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
{{- range .imports }}
    "{{ . }}"
{{- end }}
)

var (
  // needyResources is a list of resources that require a parent resource
  //
  // This list is primarily used when figuring out what nodes the dep. graph
  // should return when traversing.
  needyResources = map[string]bool{
{{- range .components -}}
{{- range .resources -}}
  {{- if not .parents }}
    {{continue}}
  {{- end }}
  "{{ .fqrt }}": true,
{{- end }}
{{- end }}
    "corteza::compose:record-datasource": true,
  }

  // superNeedyResources is the second level of filtering in case the first
  // pass removes everything
  superNeedyResources = map[string]bool{
    "corteza::compose:module-field": true,
  }
)
