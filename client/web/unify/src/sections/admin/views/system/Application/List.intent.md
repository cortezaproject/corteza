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

# Application List view

## Intention

Overview of registered applications so an admin can find one to edit,
create a new one, or remove/permission an existing one.

## UX capabilities

- Search, sort, paginate via `useResourceList`; deleted-state filter
  popover (excluded / inclusive / exclusive).
- Name column shows name + meta description; enabled state as a tag;
  createdAt shows the most recent of deleted/updated/created.
- New button; wildcard permissions (with system grant) and per-application
  permissions + delete via row actions (RBAC-gated per item).
- Row click opens the editor.

## Routes

- `system.applications` → `/system/applications`; navigates to
  `system.applications.create` and `.edit` (`:applicationID`).

## When changing this

- Delete goes through `useApplicationsStore().delete`, not the raw API —
  this keeps the shared app-selector/store state in sync; don't bypass it.
- Listing itself still uses `applicationListCancellable` directly (paged,
  filterable), independent of the store.
