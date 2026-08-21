---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Admin views

## Intention

Route-target views of the admin section, one folder per admin subarea. Every
route-target view carries its own `*.intent.md` sidecar; subarea folders carry
their own folder docs. This doc only maps the area.

## Map

- `Dashboard.vue` — admin landing page (see `Dashboard.intent.md`).
- `system/` — per-resource List/Editor views for system resources (users, roles, connections, …), plus settings, email, code snippets, action log, permissions.
- `compose/` — compose settings + permissions.
- `automation/` — workflows, sessions, TAQ, scripts, permissions.
- `federation/` — federation nodes + permissions (group is feature-flagged in the sidebar).
- `ui/` — theming, navigation, location settings.

## Data touched

Shared data notes live in the subarea folder docs; per-screen data in the
view sidecars.

## When changing this

- Route names/paths are registered in `sections/admin/routes.js` (prefixed with `/admin` by the section index) — keep folder structure and route names in sync when adding/moving views.
- **Read-only editors**: every Editor derives one `canEdit` computed —
  `!isEdit.value || !!resource.value?.<the flag that gates Save>` — and uses it
  for three things: the `general.editor.readOnly` banner, the Save button, and
  `:disabled` on every field. A form the user cannot save must not invite them
  to fill it in, so hiding Save alone is not enough. Sub-components holding
  fields (User's Security/Roles/Avatar, Role's Members, chatbot's
  General/Scenarios/Styling) take a `disabled` prop and pass it down.
