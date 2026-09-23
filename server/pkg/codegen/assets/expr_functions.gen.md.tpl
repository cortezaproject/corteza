{{- define "exprFnMd" -}}
{{- regexReplaceAll "<([a-zA-Z])" (trim .) "&lt;$1" -}}
{{- end -}}
---
title: Expression functions
description: Every function of the expression language, with its signature and an example.
outline: [2, 3]
---

<!-- This file is auto-generated from {{ range $i, $d := .Definitions }}{{ if $i }}, {{ end }}server/{{ $d.Source }}{{ end }}. -->

# Expression functions

Every [expression](./) can call these functions: field value expressions,
value sanitization and validation, and contextual roles.
A name is a function only where it is called: `split(a, b)` is the function,
while a bare `split` reads a variable of that name.

Every example on this page is evaluated by the server's test suite, so the
result shown is the result you get.
{{- range .Definitions }}
{{- range .Groups }}

## {{ .Name }}
{{- if .Intro }}

{{ template "exprFnMd" .Intro }}
{{- end }}
{{- range .Functions }}

### `{{ .Name }}` {#{{ .Name }}}

```ts
{{ .Signature }}
```

{{ template "exprFnMd" .Description }}
{{- if .Caution }}

::: warning
{{ template "exprFnMd" .Caution }}
:::
{{- end }}

```js
{{- if .Example.Scope }}
// with {{ .Example.Scope }}
{{- end }}
{{ .Example.Expr }}
{{- if .Example.Result }}
// → {{ .Example.Result }}
{{- else }}
// → {{ .Example.Note }}
{{- end }}
```
{{- end }}
{{- end }}
{{- end }}
