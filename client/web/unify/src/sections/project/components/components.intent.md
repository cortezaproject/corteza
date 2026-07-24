---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/config/kinds.js
touched-by:
  - client/web/unify/src/sections/project/views/Wizard.vue
tests: []
---

# Project components

## Intention

All non-view building blocks of the project section: the wizard machinery, the
graph, one dialog family per resource kind (the resource dialog standard), the
permissions matrix, members UI, and dashboards. Each subfolder carries its own
intent doc (nearest doc wins); this doc governs the loose shared primitives
that keep the section visually and behaviorally uniform.

## Data touched

- Loose primitives are presentational only; kind styling resolves through
  `config/kinds.js` (`kindConfig`). No stores, no APIs.

## Map

- `KindIcon.vue` — the one owner of the kind icon-square recipe (colored bg +
  ring + glyph, sizes sm–xl); accepts a `kind` or a pre-resolved `config`.
- `KindBadge.vue` — bordered chip of KindIcon + label; `interactive` makes it
  a keyboard-accessible toggle (used e.g. for role chips), `disabled`/`muted`
  are its de-emphasized states.
- `ResourceBadge.vue` — chip for a concrete resource (`{ kind, name }`),
  optionally clickable/removable.
- `DialogEyebrow.vue` — tiny uppercase kind label above every dialog title.
- `ValidationMessage.vue` — the single inline validation-error style; renders
  nothing without a message.
- Subfolders (own docs): `wizard/` (steps + toolbar + governance form),
  `graph/`, `project/` (project-level dialogs incl. members), `dashboard/`,
  `datamodel/`, `pages/`, `automations/`, `agents/`, `chatbots/`,
  `connections/`, `users/`, `roles/`, `permissions/`.

## When changing this

- New resource kinds must reuse these primitives (icon/badge/eyebrow/
  validation) — per-kind bespoke styling is drift from the dialog standard.
- Keep KindBadge a `<span>` with button semantics; a native `<button>` picks
  up the app's default button background.
