---
kind: file
covers: AllEventsView.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/views/system/ActionLog/vocab.js
  - client/web/unify/src/sections/project/composables/useEventActivity.js
touched-by: []
tests: []
---

# AllEventsView view

## Intention

Route target for **Activity** (renamed from All Events, 2026-07-28, route
`project.overview.activity` — renamed outright, no redirect). A thin wrapper
now: the screen itself is `components/dashboard/ActivityPanel.vue`, shared
with the wizard's Manage & Monitor Activity section. The project's real
audit-event log as a timeline: who did what to which project resource,
day-grouped, newest first, with a metrics band over the effective window.

## UX capabilities

- Filters (from/to, actor, origin, resource, action) fold into a popover and
  surface as removable chips; quick time-range presets window list + metrics
  in one click; clicking a value in a row drills down into a filter.
- Free-text search is client-side over loaded rows only (backend has no
  search param); infinite scroll with an explicit "load older" fallback.
- Metrics band (volume trend, event/actor/error counts) fails soft — it hides
  entirely when the viewer lacks `action-log.read`.
- Deliberate empty-state explanation: only create/update/delete are recorded
  and nothing is backfilled, so emptiness is common and explained.

## Routes

- `project.overview.activity` at `activity` under the dashboard layout.

## When changing this

- The resource/action/origin vocabulary is reused from the admin ActionLog
  section — never duplicate it here.
- Project scoping filters on `resourceProjectID` (the project owning the
  affected resource), not the request scope; list/search actions carry no
  resource and never appear. Reading is admin-only until RBAC is scope-aware.
