---
kind: file
covers: useModuleStore.js
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - lib/vue/src/stores/useRecordStore.js
  - lib/vue/src/components/input/CInputModule.vue
tests:
  - client/web/unify/src/sections/compose/stores/module.test.ts
---

# useModuleStore

## Intention

Shared compose module cache serving both the compose section (one "active"
namespace) and cross-namespace pickers, without either disturbing the other.

## State owned

`byNamespace` map (namespaceID → frozen Module[]) plus the active
`namespaceID`; `set` is a computed view of the active namespace's modules
(back-compat for compose nav/views).

## API surface consumed

`$ComposeAPI.moduleList/Read/Create/Update/Delete`.

## Consumers

Compose nav/admin/record views, record store (record construction requires the
module), module pickers, TAQ builder.

## Invariants

- `getByID` resolves across EVERY cached namespace (moduleID is globally
  unique) and re-runs when the cache changes — a component rendered before its
  namespace loaded will re-resolve.
- `load()` switches the active namespace; `loadFor()` caches a namespace
  WITHOUT switching (cache-first, no refetch per picker open).
- Cached modules are frozen; replace, never mutate.
