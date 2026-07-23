---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Automation — Server scripts admin

## Intention

Read-only inventory of server-side (corredor) automation scripts so admins can
inspect what is deployed and spot broken scripts. Scripts are authored and
deployed outside the UI — there is no create/edit/delete here.

## Map

- `Index.vue` — the whole area (see `Index.intent.md`).

## Data touched

- `$SystemAPI.automationList` — one fetch on mount; everything after is local.

## When changing this

- This area deliberately stays a single read-only view; if write operations
  ever appear, split into List/Editor like the sibling areas.
- Screen-level contracts (client-side filtering, inline error visibility) live
  in the sidecar.
