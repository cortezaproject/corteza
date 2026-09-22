---
kind: file
covers: registry.ts
backfilled: true
owner: fe
depends-on:
  - lib/vue
touched-by: []
tests:
  - client/web/unify/src/sections/taq/components/builder/form/inputs/registry.corredorScript.test.ts
  - client/web/unify/e2e/sections/taq/corredor-step.spec.ts
---

# TAQ step-config input registry

## Intention

Maps an automation function's declared `input.type` (from API function
definition segments) to the Vue input component that edits it in the TAQ
builder's step-config forms.

## Contract

- Consumed by DynamicInput: unknown types must fall back to a plain text
  input, never fail to render a step form.
- Type strings come from the server's function definitions — the server is
  the source of the vocabulary; this registry only maps it.
- Mirrored by `../viewers/registry.ts` (read-only rendering) — every editable
  type should have a viewer counterpart, or previews degrade to raw text.

## When changing this

- Keep inputs and viewers registries in sync when adding a type.
- Shared inputs come from lib/vue; only TAQ-specific editors live beside this
  file.
- `CorredorScriptSelector` / `CorredorScript` (the `corredorExec` step's script
  argument) map to lib/vue `CInputCorredorScript`; the viewer registry shows the
  script's label over the stored name.
- `Expression` is the catalog's free-typed value, used for anything from a
  script's arguments to a user's email, so it maps to a plain text box rather
  than to the expression editor the compose configurators use: the form writes
  what is typed as a literal, and nothing evaluates it.
- `isTypedValueInput(type)` answers whether a type is typed in rather than
  picked, by asking this registry which component it resolves to. The step form
  words its placeholder from that, so an unmapped type reads as typed without
  anyone listing it twice.

> **DRIFT:** nine of this registry's types have no viewer counterpart — see
> the DRIFT note in `../../builder.intent.md` for the list.
