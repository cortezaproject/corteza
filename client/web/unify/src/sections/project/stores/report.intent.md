---
kind: file
covers: report.js
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/project/views/dashboard
tests: []
---

# report store

## Intention

Server-side aggregation for the dashboards: grouped counts over the category
resources via one endpoint, so overview metrics stay correct regardless of
volume — unlike the events store, which caps its rowset loads at 200.

## State owned

None — a stateless action wrapper (no cached refs).

## API surface consumed

`$SystemAPI.projectReportReport` (GET /project-report): `report(projectId,
resource, { dimensions, metrics, from, to })` returns the raw row set
(empty dimensions ⇒ one grand-total row); `trend(...)` shapes it into
created-over-time points `{ date, group, value }` sorted ascending,
optionally split by a `groupBy` dimension.

## Consumers

Overview and CategoryView trend/metric bands (bucketing then happens via
`config/trend.js`).

## Invariants

- Never pulls full rowsets — counts only. Anything needing rows belongs in
  the events store.
- Point dates are plain `YYYY-MM-DD` strings, matching what the trend
  bucketing helpers expect.
