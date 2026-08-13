---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
touched-by:
  - client/web/unify/src/sections/project/components/dashboard/ActivityPanel.vue
  - client/web/unify/src/sections/project/components/dashboard/EventTimelineItem.vue
tests: []
---

# ActionLog

## Intention

Read-only audit trail: lets an admin investigate who did what, to which
resource, and when. Strictly a viewer — no mutation of any kind.

## Data touched

- `$SystemAPI.actionlogList` (raw entries, no resource class);
  `userRead` to resolve actor IDs into names.
- `vocab.js` — hand-maintained resource/action label maps mirroring server
  enums; the authoritative vocabulary for the filter selects.

## Map

- `List.vue` — the single audit-log screen (route target, own sidecar).
- `vocab.js` — resource/action/origin/severity label maps and helpers.

## When changing this

- Keep `vocab.js` in sync with server resource/action enums — new server
  resources otherwise show as raw `corteza::…` strings. It is the single
  source of truth beyond this folder: the project dashboard's activity feed
  imports it rather than keeping its own copy.
- The log can be huge; never fetch unbounded (cursor semantics live in the
  List sidecar).
