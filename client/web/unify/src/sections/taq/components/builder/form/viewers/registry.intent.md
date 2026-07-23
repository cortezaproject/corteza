---
kind: file
covers: registry.ts
backfilled: true
owner: fe
depends-on: []
touched-by: []
tests: []
---

# TAQ step-config viewer registry

## Intention

Read-only counterpart of `../inputs/registry.ts`: maps an automation
function's `input.type` to the component that displays a configured value in
step previews and the reference panel.

## Contract

- Must mirror the inputs registry — an editable type without a viewer entry
  degrades to raw text display.
- Viewers are display-only: they never mutate the bound value.

## When changing this

- Add entries in the same change that extends the inputs registry.
