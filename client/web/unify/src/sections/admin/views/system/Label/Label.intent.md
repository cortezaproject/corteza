---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
touched-by: []
tests: []
---

# Label

## Intention

Cross-application label management: one place to see and administer every
resource carrying a given label, spanning compose namespaces, system agents,
and NG automations. Labels group resources across component boundaries.

## Data touched

- `$SystemAPI.labelListCancellable`, `agentList/Create/Delete`;
  `$ComposeAPI.namespace*`; `$AutomationAPI.ngAutomation*`.
- All label queries filter by `labels: "<name>="`; raw objects, no classes.

## Map

- `List.vue` — label list with counts + create dialog (see sidecar).
- `Editor.vue` — per-label tagged-resource panels (see sidecar).

## When changing this

- The label name doubles as the route param and the query key; keep
  encode/decode symmetric or labels with special characters break.
- The rare admin area talking to three APIs — grant checks are per
  component (`system/`, `compose/`, `automation/`), keep them separate.
