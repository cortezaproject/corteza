---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/projects.js
  - client/web/unify/src/sections/project/config/connectors.js
  - client/web/unify/src/sections/project/config/kinds.js
touched-by:
  - client/web/unify/src/sections/project/views/Wizard.vue
  - client/web/unify/src/sections/project/components/wizard/steps/ResourceManagementStep.vue
tests: []
---

# Connection dialogs & connector picker

## Intention

The connection resource-kind's implementation of the project resource dialog
standard, shaped by connections being _configured instances of connectors_: a
two-phase CreateDialog (pick a connector from the live library, then name +
auth-param form) and a DetailDialog that re-edits that configuration. Opened
via the Wizard's `createResource('connection')` /
`inspectResource('connection', id)`.

## Data touched

- `useProjectsStore`: `connectionsFor`, `connectionLibrary` /
  `loadConnectionLibrary`, `allowedConnectorIds` (Resource Management
  whitelist; unset = full library), `prepareConnection` (imports the base
  connection to learn its `derivedParams` auth-field schema), `saveConnection`
  (create and update — update passes `configuredConnectionID`).
- `config/connectors.js`: static catalog for icons/labels; the live library is
  the source of truth for what is offered.

## Map

- `ConnectionCreateDialog.vue` — phase `pick` (searchable connector grid from
  the live library, whitelist-filtered) → phase `configure` (spinner while
  `prepareConnection` imports; then name + declared params, seeded with
  defaults). Never auto-opens the detail dialog.
- `ConnectionDetailDialog.vue` — re-imports the base connection on open to
  learn the param schema, seeds values from the stored config without
  clobbering in-progress edits, saves via `saveConnection`. Store entries are
  keyed by `configuredConnectionID` (`configurationID` exists only on the raw
  API response). **Deliberately has no `ResourcePermissionsSection`** — the
  locked permissions contract omits connections (build-time infra, no
  end-user runtime op) — and no builder deep link (none exists yet).
- `ConnectorPicker.vue` — standalone generic picker dialog (defaults to the
  static catalog); used by ResourceManagementStep, not by the create dialog,
  which embeds its own whitelist-aware grid.

## When changing this

- Do not "fix" the missing permissions section — its absence is contract.
- Param values are only persisted when non-empty; the schema always comes from
  a fresh `prepareConnection`, never from stored config alone.
