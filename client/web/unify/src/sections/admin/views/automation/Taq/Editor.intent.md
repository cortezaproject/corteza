---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - lib/js/src/automation/types/taq.ts
  - lib/vue/src/composables/useUnsavedGuard.ts
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# TAQ Editor view

## Intention

Create a TAQ (Trigger Action Query) automation or edit its metadata (name in
`meta.short`, handle, description, enabled). Graph editing happens in the
separate TAQ builder, which this view deep-links to.

## UX capabilities

- Validated form (name required, server-aligned handle regex) plus a collapsed read-only meta panel (ID, created/updated timestamps) on edit.
- "Open in builder" opens `/taq/builder/<automationID>` in a new tab.
- Per-TAQ permissions button (`canGrant`) and delete (`canDeleteNgAutomation`, hidden once deleted); save gated by `canUpdateNgAutomation` on edit.
- Unsaved-changes guard (deep diff vs loaded copy); create redirects into edit mode.

## Routes

`automation.taq.create` at `/automation/taq/new`; `automation.taq.edit` at `/automation/taq/:automationID`. Back navigates to `automation.taq`; watches the `automationID` param and reloads (builder round-trips).

## When changing this

- Create must send empty `triggers`/`steps`/`paths`; update must never overwrite them — the builder owns those.
- Display name lives in `meta.short` (workflows use `meta.name`).
- Keep parity with `../Workflow/Editor.vue`.
