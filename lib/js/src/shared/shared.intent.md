---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js/src/cast.ts
touched-by:
  - lib/vue
  - client/web/unify
  - lib/js/src/corredor
tests: []
---

# Shared types

## Intention

Types used by all three subsystems (compose, system, automation) that belong
to none of them.

## Map

- `types/attachment.ts` — `Attachment` class: file attachments with URL/preview handling, shared by record files, page blocks, and corredor helpers.
- `types/chart/colorschemes` — the bundled chart color scheme catalog; `types/chart/helper.ts` — `getColorschemeColors()` resolves a scheme name (or custom colors) to a color array.

## When changing this

- Color scheme names are persisted in chart/page-block configs — removing or renaming a scheme breaks saved charts; only add.
