---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useNamespaceStore.js
  - client/web/unify/src/sections/compose/components/Namespaces
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Namespace views

## Intention

Choose, create, configure, and enter compose namespaces (workspaces). `View.vue`
doubles as the layout shell for everything inside a namespace — all other compose
screens render as its router children.

## Data touched

- `useNamespaceStore` — set of all namespaces; lookup by slug-or-ID via
  `getByUrlPart`; create/update/clone/delete should go through the store so the
  cached set stays true (List.vue's delete does not — see its DRIFT note).
- `$ComposeAPI` — namespace list (paginated table), delete, export endpoint
  (JWT-signed download URL), namespace attachment upload for logos.
- `View.vue` additionally preloads module, page, chart, and page-layout stores.
- RBAC: section-level `compose/` `namespace.create`/`grant` checks plus per-item
  `can*` flags; permission resource `corteza::compose:namespace/<id|*>`.

## Map

- `List.vue` — namespace chooser / manager (see `List.intent.md`).
- `Edit.vue` — create + edit form (see `Edit.intent.md`).
- `View.vue` — namespace context shell (see `View.intent.md`).

## When changing this

- Slug is optional: every namespace link must fall back to `namespaceID`.
- Namespace chooser routes keep the sidebar collapsed (`meta.hideSidebar`).
- Sibling components live in `../../components/Namespaces/` (importer, translator).
