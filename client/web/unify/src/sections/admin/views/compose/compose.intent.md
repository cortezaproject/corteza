---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/components/Permissions/CPermissionGrid.vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Admin — Compose views

## Intention

Admin surface for the compose (low-code) component: instance-wide compose
settings and the component-wide compose permission matrix. Namespace/module
administration lives in the compose section itself, not here.

## Map

- `Settings/Index.vue` — compose settings editor (see `Settings/Index.intent.md`).
- `Permissions/Index.vue` — `CPermissionGrid` wrapper (see `Permissions/Index.intent.md`).

## Data touched

- Settings: `$SystemAPI.settingsList({ prefix: 'compose.' })` / `settingsUpdate`
  — compose settings are stored as system settings, not via the compose API.
- Permissions: `$ComposeAPI` permission endpoints via `CPermissionGrid`.

## When changing this

- Permission grid behavior belongs in the shared `CPermissionGrid`.
- Setting-shape contracts (whitelist round-trips, object-valued UI toggle
  groups) live in the settings sidecar.
