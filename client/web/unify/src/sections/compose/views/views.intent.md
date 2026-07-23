---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/routes.js
touched-by: []
tests: []
---

# Compose — route-target views

## Intention

Every screen the compose section's routes point at lives here. The folder is purely
route targets: each view carries its own sidecar contract, and each resource
subfolder has a folder doc for shared notes. No loose helper files belong at this
level — shared building blocks go to `../components/` or `../composables/`.

## Data touched

All views operate on compose resources (namespace, module, page, page layout,
chart, record) through `$ComposeAPI` and the shared Pinia stores from
`@planetcrust/human-vue` (lib/vue). `Namespace/View.vue` is the context holder:
it preloads the namespace-scoped stores that every nested view then reads.

## Map

- `Namespace/` — namespace chooser, editor, and the namespace layout shell.
- `Pages/` — end-user page rendering: onboarding, page view, record view.
- `Admin/` — namespace-admin screens (modules, records, pages, builder, charts).

## When changing this

- New route targets get a sidecar (`<Name>.intent.md`) — that is the file-tier rule.
- Views nested under `namespace.view` receive `:namespace` as a prop from
  `Namespace/View.vue` and assume the stores are already loaded; do not add
  per-view store preloading without checking that contract.
