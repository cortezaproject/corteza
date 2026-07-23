---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Connection List view

## Intention

Find integration connections and route the admin to the right screen: the
definition editor for local ones, the configure page for catalog ones.

## UX capabilities

- Search, sort (default by status), paginate via `useResourceList`;
  deleted-state filter popover; status tag (suppressed for catalog drafts)
  and source tag (catalog vs local).
- "Create custom" button gated on `connection.create` RBAC; wildcard
  and per-connection permissions; delete via row actions.
- Row click: catalog → configure view; local → editor, but only when the
  admin can update or delete it (otherwise inert).

## Routes

- `system.connections` → `/system/connections`; navigates to
  `system.connections.create`, `.edit`, and `.configure` (`:connectionID`).

## When changing this

- The source-based row-click branching is the screen's core contract
  (Editor/Configure redirect between each other as a safety net).
- Permission resources use the `connection` kind, matching the `connection`
  API (historically pointed at `dal-connection` — fixed as a bug).
