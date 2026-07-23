---
kind: file
covers: Dashboard.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/components/chart/CChart.vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Dashboard view

## Intention

Admin landing page: at-a-glance instance health — how many users, roles,
workflows and namespaces exist, in what state, and how each grew over time.

## UX capabilities

- Four stat cards (users/roles/workflows/namespaces) with per-status breakdown pills; card totals are computed as active + each excluded status because default API listings return active items only.
- Clicking a card loads a stacked created-per-month chart with a cumulative line; the full item list is fetched lazily per card and cached for the page's lifetime.
- Errors degrade to zero counts and console errors — never block rendering. Sets the topbar title via the `#topbar-title` teleport.

## Routes

`dashboard` at `/dashboard` — target of the `admin` section-index redirect; no params.

## When changing this

- Adding a stat card means extending `statCards`, `statusClassifiers`, `statusConfig` and `listFnMap` together.
- Full-list fetching is O(all records) per card; keep it lazy (only on card activation) and cached.
- Chart colors derive from live PrimeVue theme CSS vars so theming keeps working.
