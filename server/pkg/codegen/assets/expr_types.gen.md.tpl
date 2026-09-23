---
title: Expression types
description: The value types of the expression language and the fields of each structured type.
outline: [2, 3]
---

<!-- This file is auto-generated from {{ range $i, $d := .Definitions }}{{ if $i }}, {{ end }}server/{{ $d.Source }}{{ end }}. -->

# Expression types

Every value in an expression, workflow variable or function argument has one
of these types. Structured types have fields, read with a dot:
`record.values.name`, `user.email`.

## Value types

{{ range $i, $t := .ValueTypes }}{{ if $i }}, {{ end }}`{{ $t }}`{{ end }}

## Structured types
{{- range .StructTypes }}

### `{{ .Name }}` {#{{ toLower .Name }}}

| Field | Type | Notes |
| --- | --- | --- |
{{- range .Def.Struct }}
| `{{ .Name }}` | {{ if index $.StructNames .ExprType }}[`{{ .ExprType }}`](#{{ toLower .ExprType }}){{ else }}`{{ .ExprType }}`{{ end }} | {{ if .Readonly }}read-only{{ end }}{{ if and .Readonly .Alias }}; {{ end }}{{ if .Alias }}also `{{ .Alias }}`{{ end }} |
{{- end }}
{{- end }}
