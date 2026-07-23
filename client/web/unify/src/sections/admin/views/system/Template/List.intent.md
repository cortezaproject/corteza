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

# Template List view

## Intention

Browse rendering templates (documents, PDFs, emails), see their type and
language, and reach the editor, permissions, or delete.

## UX capabilities

- Search, sort, paginate via `useResourceList`; tri-state deleted filter in
  a popover; name column shows `meta.short` with description underneath.
- Wildcard permissions button (`template/*`) gated by system `grant`;
  per-row permissions (row `canGrant`) and delete (`canDeleteTemplate`) in
  the action menu.

## Routes

- `system.templates` → `/system/templates`.
- Header button → `system.templates.create`; row click →
  `system.templates.edit` with `templateID`.

## When changing this

- Display name is `meta.short`, not a top-level name field — keep fallbacks
  (`meta.short || handle || templateID`) consistent across list, confirm
  dialogs, and permissions titles.
- The date column shows the latest of deletedAt/updatedAt/createdAt.
