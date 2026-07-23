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

# Script Index view

## Intention

Read-only inventory of server-side (corredor) automation scripts so admins can
inspect what is deployed and spot broken scripts. Scripts are authored and
deployed outside the UI — no create/edit/delete here.

## UX capabilities

- Loads the script set once (`$SystemAPI.automationList`); everything after is client-side.
- Text query over name + label, plus has-errors / has-triggers / has-iterator / has-security toggles with live counts.
- Each row shows label (or a "missing label" hint), description, machine name, capability tags, and inline error messages per script.

## Routes

`automation.scripts` at `/automation/scripts`; no params, no outbound navigation.

## When changing this

- Keep filtering client-side over the fetched set unless the script count outgrows it — no server round-trips after mount.
- Script errors surface inline per row — that visibility is the main point of the view; don't hide them behind a detail screen.
