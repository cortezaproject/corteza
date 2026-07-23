---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/compose/components/PageBlocks/Blocks/RecordListBlock.vue
tests: []
---

# Public record tools

## Intention

End-user (non-admin) record data exchange: the Import and Export dialogs offered
from a record list's toolbar. "Public" mirrors the legacy compose app's split
between public page components and admin components.

## Data touched

Both drive the Compose API record import/export endpoints for one module.
Export streams CSV/JSON of selected fields; import uploads a CSV/JSON file,
then runs a server-side import session.

## Map

- Record/Exporter/index.vue — button + dialog: field picker, range (all/selection/filter/date range), format choice; hands off to the export endpoint
- Record/Importer/index.vue — button + stepper dialog: upload (CFileDropZone) → map file columns to module fields + on-error policy → progress/result; dialog is non-closable mid-run

## When changing this

Each component is self-contained around `module`/`namespace` props from the
record list — keep them usable without any block context. Field mapping in the
importer must track module fields; exported field selection defaults should stay
aligned with the record list's visible columns. Import runs server-side: closing
guards during step 2 protect a running session.
