---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/index.js
tests:
  - client/web/unify/e2e/sections/home/home.spec.ts
---

# Home section

## Intention

The landing page of the unified app at `/`. Gives a signed-in user one screen
with everything to start from: the app catalog, the AI agent, and their
notifications — as three side-by-side, user-resizable columns. Because those
columns are inline, the section suppresses the equivalent topbar chrome.

## Map

- `index.js` — section registration: single `home` route at `/`, no sidebar, topbar overrides hiding app selector, agent sidebar, notifications, and home button (all redundant here).
- `views/Home.vue` — the three-column screen (see `views/Home.intent.md`).
- `composables/useColumnResize.js` — drag-resize state for the left/right column widths: clamped ranges, persisted to localStorage, global mouse listeners cleaned up on unmount.

## Data touched

- `useNotificationsStore`; agent context via `useAgentRouteContextProvider('home')`;
  localStorage keys `home-column-menu-width` / `home-column-notifications-width`.

## When changing this

- Topbar overrides in `index.js` exist because the columns duplicate that
  chrome — removing a column means restoring its topbar toggle.
- Screen-level contracts (agent wiring parity, column behavior) live in the
  view sidecar.
