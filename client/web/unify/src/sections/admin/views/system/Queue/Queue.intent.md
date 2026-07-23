---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
touched-by: []
tests: []
---

# Queue

## Intention

Manage messaging queues used by automation/eventbus: queue name, consumer,
and metadata (e.g. polling behavior).

## Data touched

- `$SystemAPI.queues*` (note the plural prefix: `queuesRead`,
  `queuesCreate`, …). Plain-object resource — normalized locally, no
  `system.*` class exists for queues.

## Map

- `List.vue` — queue list with deleted filter (see sidecar).
- `Editor.vue` — queue form with duration validation (see sidecar).

## When changing this

- Queues stay plain objects with local normalization — keep the save
  payload minimal and the editors' clone/normalize pair in sync when adding
  fields.
