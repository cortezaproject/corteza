---
kind: folder
covers: '.'
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/router/index.js
  - client/web/unify/src/App.vue
tests: []
---

# Sections — the section contract

## Intention

Feature areas of the unified app. Each section is a self-contained module the shell
consumes via a fixed contract; adding a section = drop a folder + register it in
`index.js`. Each section folder carries its own intent doc.

## Contract (what a section's `index.js` default export declares)

- `id` — unique section id; also used as `route.meta.section` on its routes.
- `routes` — flat route list; paths prefixed with the section's base, route names namespaced (`<id>.` prefix). Uniqueness is convention-only, nothing enforces it.
- `sidebar` — optional left-sidebar component; routes may opt out with `meta.hideSidebar`.
- `topbar` — optional topbar setting overrides merged over `ui.topbar` settings.
- `profileItems(t)` — optional extra profile-menu items.
- `agentContext(ctx)` — optional agent-context enrichment.

## Map

- `index.js` — section registry: ordered `sections` list, flattened `routes`, `sectionById()`.
- `home`, `agentic`, `workflow`, `taq`, `admin`, `compose`, `chatbot` — the sections (order = app-list order).

## When changing this

- Contract changes here ripple into App.vue, the router, and every section — update all of them together.
