---
kind: file
covers: useAutomationStore.js
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify/src/sections/taq
tests: []
---

# useAutomationStore

## Intention

Shared cache for TAQ (next-gen) automations plus the construct-library catalog
of available functions and triggers the builder composes from.

## State owned

`list` — `NgAutomation` instances; `functions` / `triggers` catalog with a
`catalogReady` flag; `loading` / `error`.

## API surface consumed

`$AutomationAPI.ngAutomationList/Create/Delete`,
`constructLibraryFunctions/Triggers`.

## Consumers

TAQ section list views and builder; admin automation TAQ list/editor.

## Invariants

- Default `fetchList` filter includes disabled automations (`disabled: 1`).
- `loadCatalog()` gates the builder: it must not render steps until
  `catalogReady` is true; catalog load failures log but do not throw.
- Save/rename flows use `updateInList` (upsert) / `removeFromList` to keep the
  list fresh without refetching; `remove` is the API-backed delete.
