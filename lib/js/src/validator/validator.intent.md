---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on:
  - lib/js/src/guards.ts
touched-by:
  - lib/vue
  - client/web/unify
tests:
  - lib/js/src/validator/validator.test.ts
---

# Validator

## Intention

Generic validation primitives used for record-value validation: a `Validated`
result set that aggregates `ValidatorError`s and a `Validator` that runs
registered validator functions over values.

## Contract

- `ValidatorError.kind` is a translation key; `meta` carries interpolation/grouping data — UI layers translate by `kind`, `message` is only the untranslated fallback.
- Validator functions may return nothing (valid), an error, an error set, or a promise of these; `Validated` normalizes all shapes.
- Errors compare/merge by `kind` + `meta`, so duplicate registration of the same failure collapses.

## When changing this

- Compose record editing in the apps renders these errors per field; changing error shape or `kind` semantics breaks field-level error display and its translations.
