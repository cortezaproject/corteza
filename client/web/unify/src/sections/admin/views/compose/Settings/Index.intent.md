---
kind: file
covers: Index.vue
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Compose Settings Index view

## Intention

Edit instance-wide compose settings: attachment limits and mimetype
whitelists (page/record/icon) plus compose UI toggles (sidebar, record
toolbar).

## UX capabilities

- Max-size and mimetype-whitelist inputs per attachment kind; whitelists round-trip comma-separated strings to validated mimetype arrays (invalid entries silently dropped).
- Toggle groups for hiding compose sidebar elements and record-toolbar buttons.
- Single save writes the whole loaded setting set back via `$SystemAPI.settingsUpdate`.

## Routes

`compose.settings` at `/compose/settings`; no params.

## When changing this

- Compose settings are stored as system settings (`$SystemAPI.settingsList({ prefix: 'compose.' })`), not via the compose API.
- UI toggle groups are single object-valued settings (`compose.ui.sidebar`, `compose.ui.record-toolbar`) merged back on save — new toggles must join those objects, not become flat keys.
- Save is whole-set (every loaded key is written back) — parse values defensively (`@value` wrapper) when adding new setting types.
