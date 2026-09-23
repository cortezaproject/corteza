---
title: {{ .label }} permissions
outline: [2, 2]
---

<!-- This file is auto-generated from server/app and the component definitions. -->

# {{ .label }} permissions

Every operation is denied unless a role is granted it.

## Component

| Operation | Description |
| --- | --- |
{{- range .operations }}
| `{{ .handle }}` | {{ .description }} |
{{- end }}
{{ range .resources }}
## {{ .handle | replace "-" " " | title }}

| Operation | Description |
| --- | --- |
{{- range .operations }}
| `{{ .handle }}` | {{ .description }} |
{{- end }}
{{ end -}}
