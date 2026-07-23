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

# TAQ List view

## Intention

Browse all TAQ (Trigger Action Query) automations and reach their editors;
manage lifecycle (delete) and permissions from the list.

## UX capabilities

- Paginated, searchable, sortable list (name/handle/enabled/createdAt); rows open the editor. Display name comes from `meta.short`.
- Tri-state disabled/deleted filters (excluded/inclusive/exclusive) in a popover; defaults include disabled, exclude deleted.
- Wildcard permissions button gated by `automation/` grant; per-row permissions/delete gated by item `canGrant` / `canDeleteNgAutomation`.
- Delete confirms, syncs `useAutomationStore`, and refetches the list.

## Routes

`automation.taq` at `/automation/taq`; navigates to `.create` (new button) and `.edit` with `automationID` (row click).

## When changing this

- API surface is `ngAutomation*`; permission resource is `corteza::automation:ng-automation/<id|*>`.
- Keep parity with `../Workflow/List.vue` — the two lists intentionally mirror each other.
