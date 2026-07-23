---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/sidebar/AdminSidebar.vue
  - client/web/unify/src/sections/admin/views
touched-by:
  - client/web/unify/src/sections/index.js
tests: []
---

# Admin section

> Route paths in this section's intent docs are as declared in `routes.js`;
> the section registry mounts them all under the `/admin` prefix.

## Intention

The unified app's administration area — the legacy standalone admin webapp folded
into the unify shell as one section. Covers system (identity, infrastructure,
resources), compose, automation, federation and UI administration. This doc governs
only the section entry (`index.js`) and route table (`routes.js`); `sidebar/`,
`views/` and `components/` carry their own docs.

## Contract (section values declared by index.js)

- `id: 'admin'`; every route is tagged `meta.section: 'admin'` so the shell resolves the admin sidebar.
- Route prefix `/admin` is applied to every raw path from routes.js; the legacy `root` route is renamed to `admin` (section index, redirects to `dashboard`). All other legacy route names (`system.*`, `compose.*`, `automation.*`, `federation.*`, `ui.*`) are kept verbatim — deep links and sidebar mappings depend on them.
- `sidebar: AdminSidebar`, shown on every admin route (`sidebarDisabledRoutes: []`) and auto-expanded on entering the section (`sidebarExpandedByDefault: true`).
- No `topbar` overrides, `profileItems` or `agentContext` — admin uses shell defaults.
- routes.js holds the raw, unprefixed definitions and deliberately drops the legacy catch-all — the unified router owns the global 404.

## Map

- `index.js` — section-contract entry: prefixes paths, renames `root`, tags `meta.section`, wires the sidebar.
- `routes.js` — raw route table (lazy imports into `views/**`), grouped dashboard / system / compose / automation / federation / ui.
- `sidebar/` — admin navigation tree (own INTENT.md).
- `views/` — route-target views (own INTENT.md per area).
- `components/` — reusable editor sub-panels (own INTENT.md).

## When changing this

- Route names are the API: the sidebar, cross-view `router.push` calls and external deep links reference them. Renaming one requires sweeping sidebar + views.
- New resource areas follow the List / Editor (/ Configure) route convention and must also get a sidebar entry to be discoverable.
- Route-name uniqueness across sections is convention-only — nothing enforces it (see `sections/INTENT.md`).
