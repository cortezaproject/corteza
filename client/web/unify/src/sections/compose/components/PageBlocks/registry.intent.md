---
kind: file
covers: registry.ts
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/compose/components/PageBlocks
  - client/web/unify/src/sections/compose/views/Admin/Pages/Builder.vue
tests: []
---

# Page-block registry

## Intention

Single dispatch point mapping a page block's `kind` string to its async-loaded
renderer component. A block kind exists iff it has an entry here.

## Contract

- `resolveBlock(kind)` returns the component or `null` — callers must handle
  unknown kinds gracefully (persisted pages may reference removed kinds).
- All entries are `defineAsyncComponent` — block renderers are code-split and
  must tolerate lazy mounting.
- Configurators are NOT dispatched here; the Builder imports them statically.
  Adding a kind = Block component + entry here + Configurator + Builder import.

## When changing this

- Kind strings are persisted in page/layout data — renaming one orphans every
  existing block of that kind. Additive changes only.
