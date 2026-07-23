---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/composables/useResourceList.ts
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Workflow List view

## Intention

Browse all classic automation workflows and reach their editors; manage
lifecycle (delete) and permissions from the list.

## UX capabilities

- Paginated, searchable, sortable list (name/handle/enabled/createdAt); rows open the editor.
- Tri-state disabled/deleted filters (excluded/inclusive/exclusive) in a popover; defaults include disabled, exclude deleted.
- Wildcard permissions button gated by `automation/` grant; per-row permissions/delete gated by item `canGrant` / `canDeleteWorkflow`.
- Delete confirms, syncs `useWorkflowStore`, and refetches the list.

## Routes

`automation.workflows` at `/automation/workflows`; navigates to `.create` (new button) and `.edit` with `workflowID` (row click).

## When changing this

- Permission resource is `corteza::automation:workflow/<id|*>`.
- Keep parity with `../Taq/List.vue` — the two lists intentionally mirror each other (workflow display name is `meta.name`, TAQ uses `meta.short`).
