---
kind: file
covers: Home.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/composables/useAgentRouteContextProvider.ts
  - client/web/unify/src/sections/home/composables/useColumnResize.js
touched-by:
  - client/web/unify/src/sections/home/index.js
tests: []
---

# Home view

## Intention

One screen with everything a signed-in user starts from: the app catalog, the
AI agent, and their notifications — three side-by-side columns.

## UX capabilities

- Left column: searchable app list (`CAppList`, query filtered locally).
- Middle column: inline agent chat (`CAgentChat`) wired identically to `CAgentSidebar` elsewhere — same `agent.sidebar.*` translation namespace and `useAgentRouteContextProvider('home')`.
- Right column: notifications panel (`CNotificationsPanel`, store activated here).
- Outer two columns are drag-resizable and persisted; the middle column flexes.

## Routes

`home` at `/` (`meta.section: 'home'`) — the app's default landing route; no params.

## When changing this

- The agent column must stay wiring-compatible with `CAgentSidebar` so behavior matches the sidebar used in other webapps.
- The section's topbar overrides (in `index.js`) exist because these columns duplicate that chrome — removing a column means restoring its topbar toggle.
- Column widths persist via `useColumnResize` localStorage keys; keep clamping when adjusting layout.
