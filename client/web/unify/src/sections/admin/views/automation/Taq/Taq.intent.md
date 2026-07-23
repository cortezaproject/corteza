---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/composables/useResourceList.ts
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Automation — TAQ admin

## Intention

Admin area for TAQ (Trigger Action Query) automations — the next-generation
(`ng-automation`) counterpart to classic workflows: listing, metadata-level
editing, lifecycle, and permissions. Graph editing happens in the separate
TAQ builder, which the editor deep-links to.

## Map

- `List.vue` — paginated TAQ list (see `List.intent.md`).
- `Editor.vue` — metadata create/edit form + builder deep-link (see `Editor.intent.md`).

## Data touched

- `$AutomationAPI`: ngAutomationList/Read/Create/Update/Delete.
- `useAutomationStore` — list kept in sync after create/update/delete.
- RBAC: `automation/` grant check plus per-item `can*` flags; permission
  resource `corteza::automation:ng-automation/<id|*>`.

## When changing this

- Display name lives in `meta.short` (not `meta.name` as in workflows); the
  builder owns `triggers`/`steps`/`paths` — see the editor sidecar.
- Keep parity with `../Workflow/` — the two areas intentionally mirror each other.
