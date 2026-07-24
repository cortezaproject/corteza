---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/projects.js
  - client/web/unify/src/sections/project/config/kinds.js
  - client/web/unify/src/sections/project/utils/kindIcons.js
touched-by:
  - client/web/unify/src/sections/project/views/Wizard.vue
tests: []
---

# Resource graph

## Intention

**Locked:** the graph is the PRIMARY canvas-centric surface of the Build tab.
It must always reflect every project resource and relation, and clicking any
node opens that resource's **editable** detail dialog (via the injected
`inspectResource(kind, id)`) — never a read-only view. Steps and panels
support the canvas, not the other way around.

## Map

- `ResourceGraph.vue` — the live surface: metric chips (one per
  `OVERVIEW_KINDS` kind — count + visibility toggle) over an ECharts force
  graph; mounted beside every step in the Wizard's resizable split.

## Data touched

Node/edge model — the backend is the single source of truth:
`$SystemAPI.projectGraph({ projectID })` (via `store.graph`) re-derives the
graph from saved state. Nodes `{ id, name, kind }`; edges
`{ sourceID, targetID, reason }` (store renames to source/target). Known
edge reasons: `module-field-ref`, `role-rbac`, `user-role`; record
references render directed (referencing → referenced module).

- **Scoping:** strictly project-scoped — always fetched by `projectID`; the
  cross-project leak (store.Search ignoring ctx scope) was fixed server-side.
  Resources join the graph by being stamped with the project (rel_project).
- **Refresh:** `store.graphVersion` (bumped by `touch()` on every persisting
  mutation) triggers a refetch; concurrent reloads coalesce.
- **Visibility:** store-owned, not local — `graphVisibleKinds` is re-seeded
  per step from `kindsThroughStep` (the build is a process: only kinds built
  so far show); role/user (`ACCESS_KINDS`) form a separate overlay set
  auto-enabled on access steps. Hiding a kind hides its edges; chips still
  show counts.

## When changing this

- A new resource kind must land in the backend graph payload, get a
  kinds.js/OVERVIEW_KINDS chip + icon, and join `INSPECTABLE_KINDS` so its
  nodes open the detail dialog — inert nodes violate the locked contract.
- Canvas colours can't read live CSS vars — resolve theme tokens via
  getComputedStyle (existing idiom).
