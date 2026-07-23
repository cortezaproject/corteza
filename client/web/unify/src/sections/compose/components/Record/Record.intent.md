---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/views/Pages/RecordView.vue
touched-by:
  - client/web/unify/src/sections/compose/views/Namespace/View.vue
tests: []
---

# Record modal

## Intention

Show a full record page in a modal on top of whatever page the user is on,
without losing their place — the target of record-list row clicks, chart
drill-down, and "open in modal" links.

## Map

- RecordModal.vue — Dialog hosting RecordView in `in-modal` mode; mounted once per namespace in views/Namespace/View.vue

## Data touched

Reads the injected `$pageStore` for the page title; the embedded RecordView does
all record loading.

## When changing this

The contract is URL-driven: the modal opens iff both `recordID` and
`recordPageID` query params are present, and closing strips them (plus `edit`,
`cloneFromID`) from the query. Anything in compose can open a record modal by
pushing those query params — never add a prop-driven open path that bypasses the
URL, or deep links and back-button behavior break. The escape hatch button
navigates to the real `page.record` route.
