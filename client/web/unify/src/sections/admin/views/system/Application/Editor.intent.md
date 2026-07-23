---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/utils/appIcons.js
  - lib/vue
  - lib/js
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Application Editor view

## Intention

Create or edit one application: its name/enabled state plus how it presents
in the unified shell (display name, URL, listed flag, logo).

## UX capabilities

- Name is required; enabled toggle; unify panel with display name, URL, and
  "listed" visibility toggle.
- Logo via drag-and-drop upload to the application attachment endpoint;
  preview resolves through `resolveAppLogoUrl` + the built-in icon map;
  only custom logos are clearable — built-in default icons are not.
- Delete, permissions, unsaved guard; create redirects to edit.

## Routes

- `system.applications.create` → `/system/applications/new`; `.edit` →
  `/system/applications/:applicationID`; back lands on
  `system.applications`.

## When changing this

- All CRUD goes through `useApplicationsStore` (class `system.Application`)
  — keep the store in the loop so the app selector stays fresh.
- Uploaded logo URL is stored absolute (`$SystemAPI.baseURL` + path) with
  `logoID`; clearing resets logo to `''` and logoID to `'0'`.
