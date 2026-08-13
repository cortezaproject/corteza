---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/views/system/ActionLog/vocab.js
  - lib/vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# ActionLog List view

## Intention

Let an admin investigate the audit trail: filter, drill in by clicking row
values, and expand entries for their full context. Strictly read-only.

## UX capabilities

- Filter by from/to datetime, actor (user picker), origin, resource, and
  action; any filter change reloads from scratch.
- Click a resource/action/actor cell to filter by that value; clicking the
  active value clears it again.
- Expandable row detail; severity tag; "Load older" cursor pagination (no
  page numbers, no search box).

## Routes

- `system.actionLog` → `/system/action-log` (list only; links nowhere else).

## When changing this

- Cursor = `beforeActionID` of the last loaded row; keep append semantics
  and the `loadSeq` stale-response guard when touching `load()`.
- Actor names resolve via per-ID `userRead` into a local cache rather than
  through `useUserStore` — the log's actors are arbitrary historical IDs, not
  the working set the store preloads; keep per-user failures silent.
- Drill-down values missing from `vocab.js` become virtual Select options —
  keep that or the filter renders blank.
