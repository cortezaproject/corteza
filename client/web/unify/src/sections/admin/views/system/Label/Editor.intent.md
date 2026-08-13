---
kind: file
covers: Editor.vue
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
touched-by: []
tests: []
---

# Label Editor view

## Intention

Manage everything carrying one label, across component boundaries: compose
namespaces, system agents, and NG (TAQ) automations, each in its own panel.

## UX capabilities

- Per kind: table of tagged resources, create-with-label-pre-applied dialog,
  delete with confirm, and permission dialogs (wildcard header button gated
  by that component's `grant`; per-row gated by the row's `canGrant` — except
  the agent panel, which also opens on the component-wide `system/` grant).
- Row click opens the resource in its owning app (compose/agentic/taq) in a
  new tab — this view never edits the resources themselves.

## Routes

- `system.labels.edit` → `/system/labels/:labelID`; `labelID` is the
  URI-encoded label name, decoded before use. Back action → `system.labels`.

## When changing this

- All three fetches filter with `labels: "<name>="` — keep the trailing `=`.
- New automations are created disabled with empty triggers/steps/paths so
  they never fire before being built out.
- Cross-app links are built from `window.location.origin` + app base path;
  they bypass the router on purpose.
