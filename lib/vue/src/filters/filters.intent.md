---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify/src/sections
tests: []
---

# Filters

## Intention

Locale-aware display formatters ("filters" in the Vue 2 sense, now plain
importable functions) exposed through the lib's public API. The actual
formatting logic lives in `@planetcrust/human-js` `fmt` — these are thin named
wrappers so templates across sections share one date-formatting vocabulary.

## Map

- date.ts — `locFullDateTime` (long localized date+time, e.g. "Thursday, September 4, 1986 8:30 PM") and `locDate` (short localized date) wrappers over `fmt`.
- index.ts — barrel re-export; the only entry point consumers should import from.

## When changing this

- Admin/section list columns (createdAt/updatedAt cells) render through these
  — output format changes are visible everywhere at once. Datetime _field_
  viewers deliberately do not: the field's own `formatValue`, else its
  `onlyDate`/`onlyTime`/`outputRelative` options, decide the format, so
  `CFieldDateTimeViewer` formats independently.
- Keep these pure and locale-driven; new formatting logic belongs in lib/js
  `fmt`, with only a named wrapper added here.
