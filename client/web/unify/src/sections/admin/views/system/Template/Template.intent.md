---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
  - lib/js
touched-by: []
tests: []
---

# Template

## Intention

Manage rendering templates (documents, PDFs, emails): content, type, and
composition via partial templates.

## Data touched

- `$SystemAPI.template*`; class-based resource `system.Template`.

## Map

- `List.vue` — template list with deleted filter (see sidecar).
- `Editor.vue` — content editor with partials/preview (see sidecar).

## When changing this

- Partial vs full template distinction drives both the list and the
  editor's partial picker — preserve the `partial` flag semantics.
- Editor building blocks (code editor, toolbox, preview) live in
  `sections/admin/components/Template/`.
