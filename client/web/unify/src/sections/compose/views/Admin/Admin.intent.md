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

# Compose namespace administration

## Intention

Admin screens for building a namespace: modules (data model), their records,
pages with layouts and blocks, and charts. All routes live under
`/namespace/:slug/admin/*`, inside the `namespace.view` shell, so admins work in
the same context end users see.

## Data touched

Each subfolder documents its own stores/APIs. Shared traits: everything is
namespace-scoped through the `:namespace` prop; mutations go through the shared
lib/vue stores so the sidebar and public views stay in sync; permission dialogs
use `corteza::compose:<resource>/<namespaceID>/<id|*>` resources.

## Map

- `Modules/` — module list + editor, with `Records/` for admin record CRUD.
- `Pages/` — page tree, page/layout editor, and the block builder.
- `Charts/` — chart list + editor with live preview.

## When changing this

- Disabled namespaces are still reachable on `admin.*` routes for users with
  update rights (`Namespace/View.vue` contract) — do not add per-view guards
  that break that.
- Cross-links between the areas (module ↔ record page ↔ builder ↔ records) are
  part of the workflow; renaming routes breaks buttons in sibling folders.
