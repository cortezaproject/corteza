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

# AuthClient List view

## Intention

Overview of OAuth2 clients so an admin can find, create, permission, or
remove a client.

## UX capabilities

- Search, sort, paginate via `useResourceList`; deleted-state filter
  popover.
- Columns: name (meta.name + description), handle, enabled tag, grant type,
  createdAt (most recent of deleted/updated/created).
- New button; wildcard permissions (system grant) and per-client
  permissions + delete via row actions (RBAC-gated per item).
- Row click opens the editor.

## Routes

- `system.authClients` → `/system/auth-clients`; navigates to
  `system.authClients.create` and `.edit` (`:authClientID`).

## When changing this

- The delete API parameter is `clientID`, while the row/route identifier is
  `authClientID` — keep the mapping when touching CRUD calls.
- Display name falls back meta.name → handle → ID everywhere (headers,
  permission dialogs); keep that chain consistent.
