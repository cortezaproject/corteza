---
kind: file
covers: useRightSidebarStore.ts
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/App.vue
  - lib/vue/src/composables/useRightSidebarResize.ts
tests: []
---

# useRightSidebarStore

## Intention

Coordinate the single right-sidebar slot: only one panel (notifications,
agent, or a section-specific panel) may be open at a time.

## State owned

`activePanel` — name of the open panel or null.

## API surface consumed

None.

## Consumers

Shell topbar buttons, notifications panel, agent sidebar, section panels
(e.g. workflow/compose sidebars).

## Invariants

- Opening any panel implicitly closes the previous one — panels never stack.
- `closeSectionPanels()` closes everything EXCEPT the global panels
  (`notifications`, `agent`); the shell calls it on section navigation so
  global panels survive route changes while section panels do not.
- `close(name)` is a no-op unless that panel is the active one.
