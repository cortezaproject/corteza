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

# Admin — UI settings views

## Intention

Instance-wide look-and-feel and chrome configuration for the whole webapp:
branding/theming, topbar/navigation visibility, and map/location provider
setup. Everything here edits `ui.*` system settings that the app shell and
map components consume at runtime.

## Map

- `Theming/Index.vue` — branding + theme studio, `ui.studio` + logo settings (see `Theming/Index.intent.md`).
- `Navigation/Index.vue` — topbar chrome toggles + custom links/buttons, `ui.topbar` (see `Navigation/Index.intent.md`).
- `Location/Index.vue` — geosearch provider selection, `ui.location` (see `Location/Index.intent.md`).

## Data touched

- `$SystemAPI.settingsList` / `settingsUpdate` with prefixes `ui.studio`,
  `ui.topbar`, `ui.location`; logo uploads stored as `attachment:`-prefixed
  setting values resolved via the `$Settings` plugin.

## When changing this

- These settings are consumed elsewhere (app shell topbar, theme runtime,
  section topbar overrides merge over `ui.topbar`) — renaming a setting key is
  a cross-app contract change. Per-screen storage-format traps live in the sidecars.
