{{- define "md" -}}
{{- $s := . | replace "{PRODUCT_NAME}" "Human" -}}
{{- $s = regexReplaceAll "(?m)^\t+" $s "" -}}
{{- $s = regexReplaceAll "(?m)^\\[ ?NOTE ?\\]\\n====\\n" $s "::: info\n" -}}
{{- $s = regexReplaceAll "(?m)^\\[ ?TIP ?\\]\\n====\\n" $s "::: tip\n" -}}
{{- $s = regexReplaceAll "(?m)^\\[ ?IMPORTANT ?\\]\\n====\\n" $s "::: warning\n" -}}
{{- $s = regexReplaceAll "(?m)^\\[ ?(CAUTION|WARNING) ?\\]\\n====\\n" $s "::: danger\n" -}}
{{- $s = regexReplaceAll "(?m)^====$" $s ":::" -}}
{{- $s = regexReplaceAll "<([a-zA-Z])" $s "&lt;$1" -}}
{{- trim $s -}}
{{- end -}}
---
title: Environment variables
description: Every server option, with its type and default.
outline: [2, 2]
---

<!-- This file is auto-generated from server/app/options. -->

# Environment variables

The server reads these options from the environment, or from a `.env` file
next to the binary or passed with `--env-file`. Environment variables always
win.
{{ range .groups }}
## {{ .title }}
{{- if .intro }}

{{ template "md" .intro }}
{{- end }}
{{ range .options }}
### `{{ .env }}` {#{{ .env }}}

<div class="env-meta">

**Type** `{{ .type }}`{{ if .defaultValue }} · **Default** `{{ .defaultValue }}`{{ end }}{{ if .defaultNote }} · **Default** {{ .defaultNote }}{{ end }}

</div>
{{- if .description }}

{{ template "md" .description }}
{{- end }}
{{ end }}
{{- end }}
