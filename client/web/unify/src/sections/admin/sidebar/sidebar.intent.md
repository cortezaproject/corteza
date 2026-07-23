---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue/src/components/navigation/CSidebarNav.vue
touched-by:
  - client/web/unify/src/sections/admin/index.js
tests: []
---

# Admin sidebar

## Intention

The admin section's navigation tree — the single place that decides what an
administrator sees in the left sidebar and how the flat admin route table is
grouped into a mental model.

## Contract

- `AdminSidebar.vue` builds a flat item list that `CSidebarNav` renders as a tree via `_id`/`_parentId` (root parent `'0'`), with `_label`/`_icon`/`_route`; all groups render expanded (`expand-all`).
- Top-level entries: Dashboard (leaf), then groups System, Compose, Automation, Federation, UI. Group nodes carry no route; leaves navigate by route _name_ only — names must match routes.js.
- The System group is ordered by concern: identity & access (users, roles, user groups, labels, LLM providers), infrastructure (connections, data sources), resources (applications, auth clients, templates, code snippets, queues, API gateway, email), monitoring (action log), configuration (settings, permissions).
- Labels come from i18n `navigation.*` keys (the Labels entry is currently hardcoded, pending i18n).
- Conditional visibility: the Federation group appears only when the `federation.enabled` setting is on (via injected `$Settings`). There is no per-item permission/RBAC filtering — every other item shows to anyone who can open /admin; the views enforce access themselves.

## Data touched

- Reads injected `$Settings` for `federation.enabled`. No API calls, no stores.

## When changing this

- Adding a route without a sidebar entry makes it reachable only by URL — decide that deliberately.
- Item `_route` names must stay in sync with routes.js; there is no dead-link check.
- If per-permission filtering is ever introduced, gate items the way federation is gated (inside the computed `navItems`), so the tree stays reactive.
