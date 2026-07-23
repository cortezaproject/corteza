---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
touched-by: []
tests: []
---

# Label List view

## Intention

Browse all labels in use across the platform with their resource counts,
jump into a label to manage its tagged resources, and "create" new labels.

## UX capabilities

- Search labels by name; per-row resource-count tag.
- Row click opens the label's editor.
- Create-label dialog: validates the handle pattern, then navigates straight
  to the editor — no create API call; a label only materializes once a
  resource is saved carrying it.

## Routes

- `system.labels` → `/system/labels`. No dedicated create route.
- Row click / dialog submit → `system.labels.edit` with
  `labelID = encodeURIComponent(label name)`.

## When changing this

- Rows have no server ID — a synthetic `_rowKey` is built from kind+name+idx;
  keep it unique or the table misbehaves.
- The search `query` is remapped to the API's `name` param.
- Keep the create-dialog handle validation in sync with the shared handle
  pattern used by other admin editors.
