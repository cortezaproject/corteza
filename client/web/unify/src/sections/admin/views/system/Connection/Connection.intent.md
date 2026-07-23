---
kind: folder
covers: recursive
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
  - lib/js
touched-by: []
tests: []
---

# Connection

## Intention

Manage integration connections (system `connection` resource): definitions
plus the configured instances users have set up. Distinct from DataSource,
which manages DAL/database connections.

## Data touched

- `$SystemAPI.connection*` (class-based `system.Connection`, incl.
  `connectionEnable`) and `$SystemAPI.configuredConnection*` (raw objects).
- Permission buttons target `corteza::system:connection/<ID>`.

## Map

- `List.vue` — connection list, routes rows by source (sidecar).
- `Editor.vue` — local-definition editor + enable flow (sidecar).
- `Configure.vue` — configured-instances page, catalog/active (sidecar).
- `ConfiguredConnectionsPanel.vue` — non-route sub-component listing and
  deleting one connection's configured instances; embedded by Editor and
  Configure.

## When changing this

- Connection source/status drives which screen applies: Editor and
  Configure redirect to each other — keep the pairing symmetric.
- ConfiguredConnectionsPanel changes must work in both hosts.
