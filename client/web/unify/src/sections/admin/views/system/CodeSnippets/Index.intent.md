---
kind: file
covers: Index.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# CodeSnippets Index view

## Intention

Single page where an admin maintains the code snippets injected into the
platform — add, edit, enable/disable, and remove them.

## UX capabilities

- Table of snippets (name, enabled) with row click opening an edit dialog;
  New opens the same dialog pre-seeded with a script-tag skeleton.
- Dialog: enabled switch, name, script textarea; Save disabled until both
  name and script are set; Delete inside the dialog and via row actions.
- Every save/delete immediately persists — no separate page-level save.

## Routes

- `system.codeSnippets` → `/system/code-snippets`; single page, links
  nowhere else.

## When changing this

- All snippets live in one settings value (`code-snippets`, an array) and
  are persisted wholesale on each change — editing is index-based against
  that array.
- The table's primary key is `name`, so duplicate names confuse row
  identity; the dialog does not enforce uniqueness.
