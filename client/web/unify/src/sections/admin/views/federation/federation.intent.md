---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/components/Permissions/CPermissionGrid.vue
  - lib/vue/src/composables/useResourceList.ts
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Admin — Federation views

## Intention

Administer federation nodes — remote instances this one shares data with —
including the pairing handshake between instances, plus the component-wide
federation permission matrix.

## Map

- `Node/List.vue` — node list + pair dialog (see `Node/List.intent.md`).
- `Node/Editor.vue` — node create/edit + generate-URI dialog (see `Node/Editor.intent.md`).
- `Permissions/Index.vue` — `CPermissionGrid` wrapper (see `Permissions/Index.intent.md`).

## Data touched

- `$FederationAPI`: node CRUD, pairing + URI generation, permission endpoints.
- RBAC: `federation/` grant check gates the wildcard permissions button;
  permission resource `corteza::federation:node/<id|*>`.

## When changing this

- Pairing is a two-sided flow: generate a URI on one instance, paste it into
  the pair dialog on the other — keep both dialogs consistent with the
  server's pairing protocol (details in the node sidecars).
- Permission grid behavior belongs in the shared `CPermissionGrid`.
