---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/config/trend.js
  - client/web/unify/src/sections/project/config/chartColors.js
  - lib/vue
touched-by:
  - client/web/unify/src/sections/project/views/dashboard
  - client/web/unify/src/sections/project/components/dashboard
tests: []
---

# Project composables

## Intention

Shared reactive logic for the project dashboards that is neither global state
(stores) nor pure functions (utils). Session-local approval/governance state
does NOT live here — it lives in `stores/projects.js` (see the WIP note
there).

## Map

- `useEventActivity.js` — project audit-event activity from the system
  actionlog, scoped by `resourceProjectID`. `loadMetrics` returns the
  events/day pulse (adaptive day/week/month buckets via config/trend) plus
  grand totals (count, distinct actors, errors) with optional
  resource/action/origin/actor filters; `loadRecent` returns the latest
  events with actors batch-resolved; `actorName` resolves an actor for
  display. Single definition shared by the Overview activity band and
  AllEventsView's metrics band — do not fork it.

## Data touched

- `$SystemAPI.actionlogReport` / `actionlogList` — needs the global
  `action-log.read` grant; callers must treat a thrown metrics call as
  "unavailable" and hide the widget for non-admins.
- lib `useUserStore` for actor-name resolution.

## When changing this

- Filter semantics must keep mirroring AllEventsView's filter popover.
- Scoping is by `resourceProjectID` (owner of the affected resource), not
  request-scope `scope.ProjectID` — keep it that way, it works unscoped.
