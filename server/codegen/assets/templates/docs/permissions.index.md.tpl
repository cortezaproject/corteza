---
title: Permissions
---

<!-- This file is auto-generated from server/app. -->

# Permissions

Access in Human is granted to roles, per operation. Anything not granted is
denied. Each component lists the operations it checks.
{{ range .components }}
- [{{ .label }}](./{{ .handle }})
{{- end }}
