---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
touched-by: []
tests: []
---

# Queue List view

## Intention

Find and manage messaging queues: browse, filter, and remove them, and reach
the editor or queue permissions.

## UX capabilities

- Search, sort, paginate via `useResourceList`; tri-state deleted filter
  (excluded/inclusive/exclusive) in a popover.
- Wildcard permissions button (`queue/*`, gated by system `grant`); per-row
  permissions (row `canGrant` OR the system-wide `grant`) and delete
  (`canDeleteQueue`) in the action menu.

## Routes

- `system.queues` → `/system/queues`.
- Header button → `system.queues.create`; row click → `system.queues.edit`
  with `queueID`.

## When changing this

- The date column intentionally shows deletedAt ?? updatedAt ?? createdAt —
  the most recent lifecycle event, not plain creation time.
- Deleted-filter values are strings '0'/'1'/'2' expected by the API; do not
  convert to booleans.
