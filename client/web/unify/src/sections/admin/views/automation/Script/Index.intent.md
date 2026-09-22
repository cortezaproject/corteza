---
kind: file
covers: Index.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/views/automation/Script/script-inventory.js
  - lib/vue/src/corredor/script-display.ts
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests:
  - client/web/unify/src/sections/admin/views/automation/Script/script-inventory.test.js
  - client/web/unify/e2e/sections/admin/corredor-scripts.spec.ts
---

# Script Index view

## Intention

Read-only inventory of server-side (corredor) automation scripts so admins can
inspect what is deployed and spot broken scripts. Scripts are authored and
deployed outside the UI — no create/edit/delete here.

## UX capabilities

- Loads the script set on mount (`$SystemAPI.automationList`) and again on the banner's Refresh; everything else is client-side.
- Banner from the list's `enabled`/`connected`/`refreshedAt`: Corredor off on this server, enabled but unreachable, or connected with the last script refresh (relative + absolute) and a Refresh action; no banner when the fields are absent.
- Text query over name + label, plus has-errors / has-triggers / has-iterator / has-security toggles with live counts.
- Rows are grouped by extension root (`agent-sandbox`, or "scripts outside any extension"); each shows label (or a "missing label" hint), description, machine name, a kind badge (`Server` / `Client · <bundle>`), one chip per trigger (`<event> · <resource>` plus `<name> <op> <value>` constraints), an iterator chip (event, resource, action, cron), security chips (run as / allow / deny) and inline error messages.
- Kind, bundle and extension are read off the script name (`/server-scripts/<ext>/…`, `/client-scripts/<bundle>/<ext>/…`), the wire `bundle` taking precedence; a constraint without `op` is equality.
- A Server scripts / Client scripts toggle pair with counts joins the query and the four capability toggles.

## Routes

`automation.scripts` at `/automation/scripts`; no params, no outbound navigation.

## When changing this

- Scripts are listed in memory from the server's Corredor bundle, not from a
  store, so the "Last change" column is present but not sortable.
- `script-inventory.js` holds every pure helper (kind/bundle/extension, trigger
  rows, grouping, kind filter, banner state, relative time); keep the view thin
  and test the helpers, not the template.

- Keep filtering client-side over the fetched set unless the script count outgrows it — the banner's Refresh is the only server round-trip after mount.
- Script errors surface inline per row — that visibility is the main point of the view; don't hide them behind a detail screen.
