---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/events.js
  - client/web/unify/src/sections/project/stores/backlogItems.js
  - client/web/unify/src/sections/project/config/categories.js
  - client/web/unify/src/sections/project/config/chartColors.js
  - client/web/unify/src/sections/project/composables/useEventActivity.js
  - client/web/unify/src/sections/project/components/wizard/GovernanceForm.vue
touched-by:
  - client/web/unify/src/sections/project/views/dashboard
tests:
  - client/web/unify/src/sections/project/components/dashboard/EventTimelineItem.test.js
---

# Dashboard building blocks

## Intention

Parts for the Manage & Monitor dashboards. The **view set is locked and
lives in `views/dashboard/`** (Overview, Category, Backlog, All Events);
this folder holds only their reusable pieces — navigation, charts, badges,
and the event/backlog dialogs and drawers.

## Map

- `DashboardNav.vue` — rail; category entries own a `:category` route param
  and reuse the same coloured badge as their view.
- Charts (presentational — parents map store data in): `CategoryKpiRow`,
  `CategoryDonutChart`, `CategoryTrendChart`, `CategoryRankBar`,
  `ChartLegend`; colours come from `config/chartColors`.
- Event display: `EventsActivityPanel` (overview pulse + feed, via
  `useEventActivity`), `EventTimelineItem` (sentence-shaped log row, has a
  unit test), `EventDiff` (action-log `delta` old→new), `EventBadge`,
  `RiskPips` (5-pip ordinal meter), `UserCell`, `TimeRangeSelect`.
- Dialogs/drawers: `EventDetailDialog`/`EventDetailDrawer`,
  `NewEventDialog`, `BacklogItemDialog`/`BacklogItemDrawer` — all render
  per-category field schemas (config/eventForm, config/categories) through
  the wizard folder's `GovernanceForm` (deliberate cross-folder reuse).

## Data touched

Governance-event records via `stores/events.js` and linked backlog items via
`stores/backlogItems.js` (both hand store-mapped rows to these components);
user options for owner/assignee fields; config: `categories`, `chartColors`,
`dashboard`, `eventForm`, `eventKinds`, `trend`.

## When changing this

- Keep chart components presentational — data shaping stays in views/stores.
- The ordinal colour ramps (chartColors ↔ EventBadge tints ↔ RiskPips) are
  one system; change them together or not at all.
- New dashboard screens touch the locked view set — needs a ruling first;
  new widgets for existing views belong here.
- EventTimelineItem changes must keep its test green.
