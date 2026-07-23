---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useNamespaceStore.js
  - lib/vue/src/composables/useResourceList.ts
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Namespace List view

## Intention

Entry point of the compose section: pick a namespace to enter, or manage the
namespace lifecycle (create, import, export, permissions, delete) from one place.

## UX capabilities

- Two switchable modes: card grid (default; renders the whole store set with client-side name/slug search) and a server-backed paginated table.
- Clicking a namespace enters it at the `pages` route (slug, falling back to namespaceID).
- Create and import buttons gated by `compose/` `namespace.create`; wildcard permissions button gated by `compose/` `grant`.
- Per-namespace actions menu (card hover ⋮ / table row): edit, export (opens JWT-signed download URL), permissions dialog, delete — each gated by that item's `can*` flags.
- Delete confirms first; success refreshes both the store set and the table.

## Routes

`namespace.list` at `/namespaces` (sidebar hidden); `/namespaces/manage` redirects here. Links out to `namespace.create`, `namespace.edit` (`:slug`), and `pages`.

## When changing this

- Card mode reads `namespaceStore.set` while table mode fetches independently via `useResourceList` — any mutation must refresh both (`refreshAll`).
- Permission resource is `corteza::compose:namespace/<id|*>`.
- Import success must reload the store with `force: true` before refiltering.
