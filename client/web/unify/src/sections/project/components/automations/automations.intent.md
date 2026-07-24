---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/projects.js
  - client/web/unify/src/sections/project/components/permissions
touched-by:
  - client/web/unify/src/sections/project/views/Wizard.vue
tests: []
---

# Automation dialogs (project-scoped TAQs)

## Intention

The automation resource-kind's implementation of the project resource dialog
standard. A project automation is a project-scoped TAQ (Trigger Action Query);
these dialogs manage its project-facing metadata, while the actual TAQ is
built in the TAQ builder section (`taq.builder-edit`), reached by deep link.
Rendered once by the Wizard, opened via `createResource('automation')` /
`inspectResource('automation', id)`.

## Data touched

- `useProjectsStore`: `automationsFor`, `addAutomation`, `updateAutomation`.
- Router: resolves `taq.builder-edit` for the "open builder" links.

## Map

- `AutomationCreateDialog.vue` — name + description create with two commit
  buttons: plain Create, and Create-and-open-builder (tab opened synchronously
  before the await so popup blockers don't eat it; closed if the create
  fails). Never auto-opens the detail dialog.
- `AutomationDetailDialog.vue` — staged name/description/enabled (committed on
  Save); embeds `ResourcePermissionsSection` (`kind="automation"`); footer
  deep-links to the TAQ builder (new tab).

## When changing this

- TAQ stands for Trigger Action Query — never "Task Queue"; keep naming and
  i18n consistent with that.
- In the permissions matrix an automation surfaces as a single runtime "Run"
  capability (read + execute) — nothing in these dialogs may imply build-time
  permission editing; that lives in the TAQ builder.
- Dialogs open only through the Wizard's provides (dialog standard); keep the
  synchronous `window.open` pattern in the create flow.
