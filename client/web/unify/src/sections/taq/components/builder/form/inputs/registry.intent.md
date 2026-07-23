---
kind: file
covers: registry.ts
backfilled: true
owner: fe
depends-on:
  - lib/vue
touched-by: []
tests: []
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
