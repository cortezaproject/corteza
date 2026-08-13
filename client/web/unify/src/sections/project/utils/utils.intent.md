---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/config/kinds.js
touched-by:
  - client/web/unify/src/sections/project/stores/projects.js
  - client/web/unify/src/sections/project/components/datamodel
  - client/web/unify/src/sections/project/components/graph
tests: []
---

# Project utils

## Intention

Pure, stateless helpers shared across the project section — anything that
would otherwise be duplicated between the store and components.

## Map

- `fields.js` — `fieldName(label)`: derives the compose field handle from a
  human label. Shared by the store (marshalling) and the data-model editors
  (predicting the machine name of unsaved fields); both sides MUST use it so
  predictions match what gets persisted.
- `kindIcons.js` — `kindIconDataUri(kind, { external, badge })`: inline-SVG
  data URIs per resource kind for ECharts node symbols; reuses the exact
  PrimeIcons artwork and kind color from `config/kinds.js` so graph nodes
  match the metrics strip. `badge` draws the missing/warning severity mark,
  `external` the out-of-project tile.

## When changing this

- `fieldName` output is a persistence contract — changing it orphans the
  label→handle prediction for existing modules.
- Kind icon changes belong in `config/kinds.js`; `kindIcons.js` must keep
  following it.
