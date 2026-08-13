---
kind: file
covers: useNamespaceStore.js
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify/src/App.vue
tests:
  - client/web/unify/src/sections/compose/stores/namespace.test.ts
---

# useNamespaceStore

## Intention

App-wide compose namespace cache; the shell preloads all namespaces once so
sections and pickers can resolve them synchronously.

## State owned

`set` — frozen `compose.Namespace` list, upserted by namespaceID.

## API surface consumed

`$ComposeAPI.namespaceList/Read/Create/Clone/Update/Delete`.

## Consumers

Unify shell (preload on mount), compose section (URL slug → namespace via
`getByUrlPart`), namespace pickers, agentic editor (resolves namespace names in
an agent's tool-access rules).

## Invariants

- `getByUrlPart` matches slug OR namespaceID — compose routes rely on both
  working.
- `findByID` is cache-first; mutations return fresh `compose.Namespace` copies
  while the cached originals stay frozen.
