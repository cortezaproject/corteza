---
kind: file
covers: Editor.vue
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/views/system/Connection/ConfiguredConnectionsPanel.vue
  - lib/vue
  - lib/js
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Connection Editor view

## Intention

Author a local integration connection definition: metadata plus raw JSON
service/resources/operations configuration, and promote a draft to active.

## UX capabilities

- Tabs: General (name, handle, description, service JSON), Configuration
  (resources + operations JSON), Configured (embedded
  `ConfiguredConnectionsPanel`, only when the connection is active).
- Enable button (local draft, confirm) flips status to active via
  `connectionEnable`; delete, permissions, unsaved guard incl. the raw
  JSON text; create redirects to edit after first save.

## Routes

- `system.connections.create` → `/system/connections/new`; `.edit` →
  `/system/connections/:connectionID`; a catalog connection loaded here
  redirects to `.configure`; back lands on `system.connections`.

## When changing this

- The raw JSON strings are the source of truth at save (resolver-validated,
  parsed last-minute); invalid JSON must keep blocking submit.
- Keep the redirect pairing with `Configure.vue` intact; saves are also
  blocked without `canUpdateConnection`.
