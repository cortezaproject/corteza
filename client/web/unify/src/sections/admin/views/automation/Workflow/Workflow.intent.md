---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/components/Workflow/WorkflowTriggers.vue
  - lib/vue/src/composables/useResourceList.ts
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Automation — Workflows admin

## Intention

Administer classic automation workflows: browse them and edit their metadata
plus lifecycle (create, delete, permissions). This is not the visual workflow
builder — only the admin-facing metadata layer.

## Map

- `List.vue` — paginated workflow list (see `List.intent.md`).
- `Editor.vue` — metadata create/edit form (see `Editor.intent.md`).

## Data touched

- `$AutomationAPI`: workflowList/Read/Create/Update/Delete, triggerList.
- `useWorkflowStore` — list kept in sync after create/update/delete.
- RBAC: `automation/` grant check plus per-item `can*` flags; permission
  resource `corteza::automation:workflow/<id|*>`.

## When changing this

- Keep parity with `../Taq/` — the two areas intentionally mirror each other.
- Screen-level contracts (filters, validation, unsaved guard) live in the sidecars.
