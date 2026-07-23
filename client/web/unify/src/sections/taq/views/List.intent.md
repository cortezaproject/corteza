---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useAutomationStore.js
  - client/web/unify/src/sections/taq/components/common/TaqConfigModal.vue
touched-by:
  - client/web/unify/src/sections/taq/index.js
tests: []
---

# TAQ List view

## Intention

Entry point of the TAQ section: browse all TAQ (Trigger Action Query)
automations and jump into the builder. Creation here is metadata-only — the
graph itself is always built in the Builder.

## UX capabilities

- Paginated, searchable `CResourceList` keyed by `automationID`: name (with description subtitle, "Untitled" fallback), enabled/disabled tag, last-change timestamp (deleted > updated > created).
- Filter popover with tri-state deleted/disabled filters ('0' exclude / '1' include / '2' only).
- "New" opens `TaqConfigModal` in create mode (creates the record, then navigates into the builder).
- Row click navigates to `/taq/builder/<automationID>`.
- Per-row actions: permissions dialog for `corteza::automation:ng-automation/<id>` (shown when the row's `canGrant` or the global grant right holds) and delete (`canDeleteNgAutomation`, confirm dialog, removes via `automationStore.remove` so the sidebar list stays in sync).
- Header wildcard-permissions button gated by `rbac.can('automation/', 'grant')`.

## Data touched

`$AutomationAPI.ngAutomationListCancellable` via `useResourceList` (default
sort `name`, limit 50); deletes through `useAutomationStore`; RBAC store for
grant checks.

## Routes

`taq` at `/taq` — the section landing route; declared with
`meta.hideSidebar` so the list opens with the drawer collapsed (legacy-app
parity). Linked from the sidebar "Automations" root item and the topbar app
navigation.

## When changing this

- Deletion must go through the automation store (not the API directly) or the
  sidebar nav list drifts.
- Keep filter values as strings '0'/'1'/'2' — the API expects that encoding.
