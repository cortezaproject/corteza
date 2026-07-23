---
kind: file
covers: useChartStore.js
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify/src/sections/compose
tests: []
---

# useChartStore

## Intention

Compose chart cache for the active namespace, shared between chart admin
views, chart page blocks and chart pickers.

## State owned

`set` — frozen `compose.Chart` list for `namespaceID` (single-namespace,
unlike the module store).

## API surface consumed

`$ComposeAPI.chartList/Read/Create/Update/Delete`.

## Consumers

Compose chart admin views, ChartBlock, chart configurators/pickers.

## Invariants

- `loadFor(namespaceID)` is cache-first (picker use case — no refetch per
  dropdown open); `load({ clear: true })` forces a reload.
- Switching namespaces without `clear` leaves stale charts mixed in — callers
  changing namespace must clear.
- Cached charts are frozen; `findByID` and mutations return fresh
  `compose.Chart` copies.
