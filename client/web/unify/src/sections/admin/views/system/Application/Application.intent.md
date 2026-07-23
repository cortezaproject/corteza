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

# Application

## Intention

Registry of applications shown to end users (app selector / listing):
name, enabled state, unify presentation (visibility, URL), and logo.

## Data touched

- Class-based resource `system.Application` (from `@planetcrust/human-js`);
  CRUD goes through `useApplicationsStore` so the shared app-selector state
  stays in sync (listing uses `applicationListCancellable` directly).
- Logo upload via `$SystemAPI.applicationUploadEndpoint()`; logo URL is
  resolved against `$SystemAPI.baseURL`.

## Map

- `List.vue` — application list (route target, own sidecar).
- `Editor.vue` — application create/edit incl. logo (route target, sidecar).

## When changing this

- Mutations must keep going through the applications store — bypassing it
  leaves the app selector stale.
- Logo URLs are absolute (baseURL-prefixed) at save time — keep resolution
  logic consistent with wherever the app list renders them.
