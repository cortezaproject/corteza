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
- `app` — **required.** The `unify.url` of the registry application that gates entry: the router admits the user only if they hold `access` on it. `null` marks a section every authenticated user reaches (home alone, where a user with no applications lands). A section declaring neither is denied — the omission fails closed, and `sections.access.test.js` catches it.
- `routes` — flat route list; paths prefixed with the section's base, route names namespaced (`<id>.` prefix). Uniqueness is convention-only, nothing enforces it.
- `sidebar` — optional left-sidebar component; routes may opt out with `meta.hideSidebar`.
- `topbar` — optional topbar setting overrides merged over `ui.topbar` settings.
- `profileItems(t)` — optional extra profile-menu items.
- `agentContext(ctx)` — optional agent-context enrichment.

## Map

- `index.js` — section registry: ordered `sections` list, flattened `routes`, `sectionById()`.
- `home`, `agentic`, `workflow`, `taq`, `admin`, `compose`, `chatbot`, `project` — the sections (order = app-list order). `project` is WIP in parts; its locked contracts live in `project/project.intent.md`.

## Access

Two operations, two questions. `read` on the registry application decides whether
it appears in the app menu; `access` decides whether the section can be opened.
Granting `access` without `read` is a dead grant — the shell resolves the section
through the read-filtered list, so an application the user cannot see is one they
cannot enter.

The gate is the router's, not a view's: it resolves before the section's view
mounts, so a refused section issues no requests. It is a usability boundary, not
the security one — every endpoint behind it enforces its own permissions.

## When changing this

- Contract changes here ripple into App.vue, the router, and every section — update all of them together.
- A new section must declare `app`, and its registry application must exist in `server/provision/101_applications/`.
