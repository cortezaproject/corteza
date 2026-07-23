---
kind: file
covers: useApplicationsStore.js
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/App.vue
  - lib/vue/src/components/navigation/CAppListSidebar.vue
tests: []
---

# useApplicationsStore

## Intention

Source of truth for registered applications: what the unify app menu shows and
whether the section behind a path is usable.

## State owned

`apps` — raw application objects from the API, in weight order.

## API surface consumed

`$SystemAPI.applicationList/Read/Create/Update/Delete/Reorder`.

## Consumers

Unify shell (preload + disabled-section screen), app-list sidebar/menu, admin
application views.

## Invariants

- Menu visibility is `unify.listed`; usability is `enabled`. A listed-but-
  disabled app still appears in the menu but opens to the "disabled" screen —
  keep these concerns separate.
- `isPathEnabled(path)` / `isCurrentAppEnabled()` return true when NO app
  matches — unregistered paths (home/root) must never be blocked.
- `reorder` keeps unlisted apps at the tail; mutations sync the local list so
  no refetch is needed.
