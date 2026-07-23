---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on: []
touched-by:
  - lib/vue
  - client/web/unify
tests: []
---

# Formatting

## Intention

Locale-aware date/time and number formatting exported as the `fmt` namespace,
so all apps format values identically.

## Map

- `datetime.ts` — parse via moment, format via `Intl.DateTimeFormat` (fullDateTime/dateTime/date/time variants).
- `number.ts` — `Intl.NumberFormat` formatting + accounting style (negatives in parentheses, zero as `-`).
- `locale.ts` — `currentLanguage()`: browser-derived language used by the Intl calls (temporary until wired to user preference).

## When changing this

- Output must stay locale-driven (Intl), not hard-coded format strings — record field viewers and dashboards rely on consistent rendering.
- If `currentLanguage()` is ever wired to the user's preferred language, all formatting changes app-wide at once; verify against i18n expectations.
