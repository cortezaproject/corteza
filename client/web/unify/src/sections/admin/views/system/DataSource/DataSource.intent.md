---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
  - lib/js
touched-by: []
tests: []
---

# DataSource

## Intention

Manage DAL connections — external databases records can be stored in —
and surface the primary (built-in) connection. Named "data sources" in the
UI; the underlying resource is `dal-connection`.

## Data touched

- `$SystemAPI.dalConnection*`; class-based resource `system.DalConnection`.
- Two connection types: `corteza::system:primary-dal-connection` (the
  built-in) and `corteza::system:dal-connection` (external).

## Map

- `List.vue` — primary panel + external list (route target, own sidecar).
- `Editor.vue` — DAL connection create/edit (route target, own sidecar).

## When changing this

- The external/primary type split is the area's core contract: the primary
  connection must never appear deletable or in the paged set.
- Do not confuse with the sibling Connection area (integration
  connections); both use a `:connectionID` param but different APIs.
