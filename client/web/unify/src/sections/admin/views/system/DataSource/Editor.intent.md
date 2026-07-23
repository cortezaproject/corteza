---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue
  - lib/js
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# DataSource Editor view

## Intention

Create or edit a DAL connection: identity, location, data-protection
properties, and (for privileged admins) the actual DAL wiring.

## UX capabilities

- Name required, handle format-validated; location as a free-text name plus
  a map-picked point (`CInputLocation`); ownership text field.
- Properties panel (edit-only): encryption/protection/restoration toggles
  with notes. DAL panel only with `canManageDalConfig` (warning otherwise):
  model ident, type, params JSON; server-reported issues shown as errors.
- Delete, permissions, unsaved guard (incl. raw params text); create
  redirects to edit after first save.

## Routes

- `system.dataSources.create` → `/system/data-sources/new`; `.edit` →
  `/system/data-sources/:connectionID`; back lands on `system.dataSources`.

## When changing this

- DAL params must parse to a plain JSON object (resolver rejects
  arrays/null); parsed on blur and again before save — raw text wins.
- `ensureLocationShape` keeps `meta.location` a valid GeoJSON Feature —
  never write coordinates without it; primary connections are edited here
  too, so don't assume the external type.
