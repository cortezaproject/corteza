---
kind: file
covers: dateUtils.js
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/project/stores/events.js
  - client/web/unify/src/sections/project/stores/backlogItems.js
tests: []
---

# dateUtils helper

## Intention

Not a store: the shared date-normalization helper for the dashboard stores'
payloads. PrimeVue date pickers emit JS Date objects; the backend stores
dates as plain `YYYY-MM-DD` strings (the report endpoint's overdue metric
parses them as ISO), so every create/update payload normalizes here instead
of each store serializing its own.

## State owned

None — pure functions.

## API surface consumed

None.

## Consumers

`events.js` and `backlogItems.js` (their DATE_KEYS on every create/update).

## Invariants

- `toISODate` converts Dates to local-timezone `YYYY-MM-DD` (invalid Dates
  become `''`); non-Date values pass through untouched, so already-ISO
  strings round-trip.
- `normalizeDates(body, keys)` mutates and returns `body`; keys absent from
  `body` are left untouched.
- Any new date field a dashboard store persists must go through this helper.
