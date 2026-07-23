---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/workflow/components/WorkflowEditor.vue
  - lib/js
  - lib/vue
touched-by: []
tests: []
---

# Workflow editor view

## Intention

Route target for creating and editing one workflow. Owns the persistence
lifecycle (load, save, delete, undelete) via `$AutomationAPI` workflow + trigger
CRUD; the canvas/editing UX is delegated to `WorkflowEditor.vue` (see
`components/components.intent.md`).

## UX capabilities

- Create mode starts a blank enabled `automation.Workflow` (`runAs: '0'`);
  edit mode fetches workflow + triggers, syncing shared `workflowStore`.
- Save is throttled (500ms, leading) and ordered: new workflows created first
  to get a real workflowID, removed triggers deleted, remaining triggers
  created/updated (`stepID` → `workflowStepID`), then the workflow updated
  and triggers refetched. First save redirects to `workflow.edit`.
- Unsaved-changes guard while dirty and not deleted; delete returns to the
  list, undelete restores in place.

## Routes

- `workflow.create` → `/workflow/new`; `workflow.edit` →
  `/workflow/:workflowID/edit` (linked from list rows and sidebar). Instance is
  reused across param changes: a watcher reloads on `workflowID` change with a
  load-sequence guard against stale responses.

## When changing this

- Workflow and triggers props must update atomically — a stale triggers array
  (missing fresh triggerIDs) makes the next save duplicate triggers.
- Preserve create-before-triggers order; triggers created with `workflowID: '0'` are orphaned.
