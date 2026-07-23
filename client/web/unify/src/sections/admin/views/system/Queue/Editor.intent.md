---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
touched-by: []
tests: []
---

# Queue Editor view

## Intention

Create or edit a single queue: its name, which consumer processes it, and
the poll delay for polling consumers.

## UX capabilities

- Required queue name and consumer (store / eventbus / corteza / redis).
- Optional poll delay validated as a Go-style duration (`1h`, `1m15s`, …).
- Per-queue permissions button in edit mode; delete gated by
  `canDeleteQueue`; unsaved-changes guard on leave.

## Routes

- `system.queues.create` → `/system/queues/new`;
  `system.queues.edit` → `/system/queues/:queueID`.
- Successful create redirects to `.edit`; back action → `system.queues`;
  load failure bounces to the list.

## When changing this

- Save payload is exactly `{ queue, consumer, meta }` (+ `queueID` on
  update) — queues are plain objects normalized locally; when adding fields
  extend `normalizeQueue`/`newQueue` too or dirty-tracking breaks.
- Poll delay lives at `meta.poll_delay` (snake_case), matching the server.
