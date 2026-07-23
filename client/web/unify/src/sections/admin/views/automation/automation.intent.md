---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
touched-by: []
tests: []
---

# Admin — Automation views

## Intention

Admin area for everything automation-related on the server's automation component:
classic workflows, the newer TAQ (Trigger Action Query) automations, corredor server
scripts, execution/session monitoring, and component-wide automation permissions.
Each resource lives in its own subfolder with its own intent doc.

## Map

- `Workflow/` — workflow list + metadata editor (classic automation workflows).
- `Taq/` — TAQ (Trigger Action Query) automation list + metadata editor; deep-links to the TAQ builder.
- `Script/` — read-only registry of server (corredor) scripts.
- `Session/` — execution monitoring: TAQ executions + workflow sessions, session detail.
- `Permissions/` — component-wide automation permission grid.

## Data touched

All subfolders talk to `$AutomationAPI` (workflows, triggers, sessions,
ng-automation/TAQ, permissions); `Script/` and user lookups use `$SystemAPI`.

## When changing this

- Route names are `automation.*` and registered in `sections/admin/routes.js` —
  keep names and folder structure in sync when adding/moving views.
- Workflow and TAQ views are deliberate near-twins (list + editor pattern);
  changes to one usually belong in the other too.
