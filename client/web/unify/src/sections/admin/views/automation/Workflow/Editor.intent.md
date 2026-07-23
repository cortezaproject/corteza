---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/components/Workflow/WorkflowTriggers.vue
  - lib/vue/src/composables/useUnsavedGuard.ts
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Workflow Editor view

## Intention

Create a workflow or edit its admin-facing metadata (name, handle,
description, enabled, trace) — not the visual workflow builder.

## UX capabilities

- Validated form: name required, handle checked against the server handle regex; invalid submit warns and scrolls to the first error.
- Edit mode adds a read-only trigger listing (`WorkflowTriggers`), per-workflow permissions button, and delete (gated by `canDeleteWorkflow`, hidden once deleted).
- Save gated by `canUpdateWorkflow` on edit; create redirects into edit mode.
- Unsaved-changes guard: dirty = deep diff vs the loaded copy.

## Routes

`automation.workflows.create` at `/automation/workflows/new`; `automation.workflows.edit` at `/automation/workflows/:workflowID`. Back navigates to `automation.workflows`; session detail links here.

## When changing this

- New editable fields must join both the update payload and the initial clone, or the dirty check breaks.
- Payload sends only handle/enabled/trace/meta — never steps or triggers.
- Keep parity with `../Taq/Editor.vue`; handle regex must stay aligned with server rules.
